package controller

import (
	"go-domotique/events"
	"go-domotique/logger"
	"go-domotique/models"
	"net/http"
	"strings"
)

func getControllerEvents(config *models.Configuration) {
	http.HandleFunc("/event/new", func(w http.ResponseWriter, r *http.Request) {
		eventId := strings.ToUpper(r.URL.Query().Get("id"))
		eventValue := strings.ToUpper(r.URL.Query().Get("value"))
		eventZwave := strings.ToUpper(r.URL.Query().Get("zwave"))
		logger.Info(config, false, "getControllerEvents", "Request received events %v (%s / %s / %s)", r.URL.Query(), eventId, eventValue, eventZwave)

		if eventId != "" && eventValue !="" {
			logger.Info(config, false, "getControllerEvents", "Request succeeded")
			w.WriteHeader(200)
			go events.CatchEvent(config, eventId, eventValue, eventZwave)
		} else {
			logger.Error(config, false,"getControllerEvents", "Request failed")
			w.WriteHeader(500)
		}
	})
}
