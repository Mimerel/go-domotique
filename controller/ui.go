package controller

import (
	"fmt"
	"go-domotique/heating"
	"go-domotique/logger"
	"go-domotique/models"
	"net/http"
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
				go heating.RunAction(config, eventId, action, payload)
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
	http.HandleFunc("/reconnect", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		config.Logger.Info("Called :  %s - %s - %+v", r.Method, r.URL.Path, r.URL.RawQuery)
		switch r.Method {
		case http.MethodOptions:
			w.WriteHeader(200)
			return
		case http.MethodGet:
			config.Channels.MqttReconnect <- true
			w.WriteHeader(200)
			return
		default:
			config.Logger.Error("Method <%s> not supported for /reconnect", r.Method)
			w.WriteHeader(500)
			_, _ = fmt.Fprintf(w, "Method <%s> not supported for /reconnect", r.Method)
			return
		}

	})
}

