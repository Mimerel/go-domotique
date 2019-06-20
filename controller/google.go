package controller

import (
	"go-goole-home-requests/googleAssistant"
	"go-goole-home-requests/models"
	"net/http"
	"strings"
)

func getControllerGoogleAssistant(config *models.Configuration) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		urlParams := strings.Split(urlPath, "/")
		config.Logger.Info("Request received question %s / %d", urlPath, len(urlParams))

		if len(urlParams) == 3 {
			config.Logger.Info("Request succeeded")
			googleAssistant.AnalyseRequest(w, r, urlParams, config)
		} else {
			config.Logger.Info("Request failed")
			w.WriteHeader(500)
		}
	})

}