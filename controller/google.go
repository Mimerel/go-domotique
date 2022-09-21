package controller

import (
	"go-domotique/googleAssistant"
	"go-domotique/logger"
	"go-domotique/models"
	"net/http"
	"strings"
)

func getControllerGoogleAssistant(config *models.Configuration) {
	http.HandleFunc("/google/", func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		urlParams := strings.Split(urlPath, "/")
		logger.Info(config, false, "getControllerGoogleAssistant", "Request received question %s / %d", urlPath, len(urlParams))
		if len(urlParams) == 3 {
			logger.Info(config, false, "getControllerGoogleAssistant", "Request succeeded")
			googleAssistant.AnalyseRequest(w, r, urlParams, config)
			return
		}
		logger.Error(config, false, "getControllerGoogleAssistant", "Request failed")
		w.WriteHeader(500)

	})
}
