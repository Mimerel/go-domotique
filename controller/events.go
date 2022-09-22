package controller

import (
	"go-domotique/events"
	"go-domotique/models"
	"net/http"
	"strings"
)

func getControllerEvents(config *models.Configuration) {
	http.HandleFunc("/event/new", func(w http.ResponseWriter, r *http.Request) {
		eventId := strings.ToUpper(r.URL.Query().Get("id"))
		eventValue := strings.ToUpper(r.URL.Query().Get("value"))
		eventZwave := strings.ToUpper(r.URL.Query().Get("zwave"))
		config.Logger.Info("getControllerEvents Request received events %v (%s / %s / %s)", r.URL.Query(), eventId, eventValue, eventZwave)

		if eventId != "" && eventValue != "" {
			config.Logger.Info("getControllerEvents Request succeeded")
			w.WriteHeader(200)
			go events.CatchEvent(config, eventId, eventValue, eventZwave)
		} else {
			config.Logger.Info("getControllerEvents Request failed")
			w.WriteHeader(500)
		}
	})
}
