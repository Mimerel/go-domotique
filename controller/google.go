package controller

import (
	"go-domotique/googleAssistant"
	"go-domotique/models"
	"go-domotique/logger"
	"net/http"
	"strings"
)

func getControllerGoogleAssistant(config *models.Configuration) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		urlParams := strings.Split(urlPath, "/")
		logger.Info(config, "getControllerGoogleAssistant", "Request received question %s / %d", urlPath, len(urlParams))
		if len(urlParams) == 3 {
			logger.Info(config, "getControllerGoogleAssistant", "Request succeeded")
			googleAssistant.AnalyseRequest(w, r, urlParams, config)
		} else {
			logger.Error(config, "getControllerGoogleAssistant", "Request failed")
			w.WriteHeader(500)
		}
	})

}