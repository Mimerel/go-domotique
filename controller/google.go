package controller

import (
	"go-domotique/googleAssistant"
	"go-domotique/models"
	"net/http"
	"strings"
)

func getControllerGoogleAssistant(config *models.Configuration) {
	http.HandleFunc("/google/", func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		urlParams := strings.Split(urlPath, "/")
		config.Logger.Info("getControllerGoogleAssistant Request received question %s / %d", urlPath, len(urlParams))
		if len(urlParams) == 3 {
			config.Logger.Info("getControllerGoogleAssistant Request succeeded")
			googleAssistant.AnalyseRequest(w, r, urlParams, config)
			return
		}
		config.Logger.Error("getControllerGoogleAssistant Request failed")
		w.WriteHeader(500)

	})
}
