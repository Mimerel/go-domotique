package controller

import (
	"go-domotique/models"
	"go-domotique/logger"
	"go-domotique/wifi"
	"net/http"
	"strings"
)

func getControllerWifiCdes(config *models.Configuration) {
	http.HandleFunc("/wifi/", func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		urlParams := strings.Split(urlPath, "/")
		logger.Info(config, false, "getControllerWifiCdes", "Request received for wifi device %s / %d", urlPath, len(urlParams))
		if len(urlParams) == 4 {
			logger.Info(config,false,  "getControllerWifiCdes", "Request succeeded")
			wifi.AnalyseRequest(w, r, urlParams, config)
		} else {
			logger.Error(config, true, "getControllerWifiCdes", "Request failed")
			w.WriteHeader(500)
		}
	})

}