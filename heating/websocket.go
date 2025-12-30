package heating

import (
	"encoding/json"
	"go-domotique/models"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

type WebSocketHub struct {
	clients    map[*websocket.Conn]bool
	mutex      sync.RWMutex
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
}

var hub = &WebSocketHub{
	clients:    make(map[*websocket.Conn]bool),
	broadcast:  make(chan []byte),
	register:   make(chan *websocket.Conn),
	unregister: make(chan *websocket.Conn),
}

func (h *WebSocketHub) run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mutex.Unlock()
		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					client.Close()
					delete(h.clients, client)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

func StartWebSocketHub(config *models.Configuration) {
	go hub.run()
	go broadcastUpdates(config)
}

func broadcastUpdates(config *models.Configuration) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		config.Channels.MqttGetArray <- true
		deviceList := <-config.Channels.MqttArray

		data, err := json.Marshal(deviceList)
		if err != nil {
			config.Logger.Error("WebSocket: Unable to marshal device list: %v", err)
			continue
		}

		hub.mutex.RLock()
		clientCount := len(hub.clients)
		hub.mutex.RUnlock()

		if clientCount > 0 {
			hub.broadcast <- data
		}
	}
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request, config *models.Configuration) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		config.Logger.Error("WebSocket upgrade failed: %v", err)
		return
	}

	hub.register <- conn

	// Send initial data
	config.Channels.MqttGetArray <- true
	deviceList := <-config.Channels.MqttArray
	data, err := json.Marshal(deviceList)
	if err == nil {
		conn.WriteMessage(websocket.TextMessage, data)
	}

	// Keep connection alive and handle disconnection
	go func() {
		defer func() {
			hub.unregister <- conn
		}()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}
