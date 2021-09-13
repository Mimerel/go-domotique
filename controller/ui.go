package controller

import (
	"fmt"
	"go-domotique/logger"
	"go-domotique/models"
	"net/http"
	"strconv"
)

func getActions(config *models.Configuration) {
	http.HandleFunc("/runAction", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		config.Logger.Info("Called :  %s - %s - %+v", r.Method, r.URL.Path, r.URL.RawQuery)
		switch r.Method {
		case http.MethodOptions:
			w.WriteHeader(200)
			return
		case http.MethodGet:
			eventId := r.URL.Query().Get("id")
			action := r.URL.Query().Get("action")
			payload := r.URL.Query().Get("payload")
			logger.Info(config, false, "getAction", "Request to do <%v> on device <%v> with payload %v)", action, eventId, payload)
			if eventId != "" && action != "" {
				go runAction(config, eventId, action, payload)
				logger.Info(config, false, "getActions", "Request succeeded")
				//go events.CatchEvent(config, eventId, eventValue, eventZwave)
				w.WriteHeader(200)
			} else {
				logger.Error(config, false, "getActions", "Request failed")
				w.WriteHeader(500)
			}
			return
		default:
			config.Logger.Error("Method <%s> not supported for /runAction", r.Method)
			w.WriteHeader(500)
			_, _ = fmt.Fprintf(w, "Method <%s> not supported for /runAction", r.Method)
			return
		}

	})
}

func runAction(config *models.Configuration, idString string, action string, payload string) {
	id, err := strconv.Atoi(idString)
	if err != nil {
		logger.Error(config, false, "runAction", "Error Converting Id to int64 %v", err)
		return
	}
	actionParams := models.MqttSendMessage{
		DomotiqueId: int64(id),
		Topic:       action,
		Payload:     payload,
	}
	config.Channels.MqttSend <- actionParams
}
