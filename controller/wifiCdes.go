package controller

import (
	"go-domotique/models"
	"go-domotique/wifi"
	"net/http"
	"strings"
)

func getControllerWifiCdes(config *models.Configuration) {
	http.HandleFunc("/wifi/", func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		urlParams := strings.Split(urlPath, "/")
		config.Logger.Info("getControllerWifiCdes Request received for wifi device %s / %d", urlPath, len(urlParams))
		if len(urlParams) == 4 {
			config.Logger.Info("getControllerWifiCdes Request succeeded")
			wifi.AnalyseRequest(w, r, urlParams, config)
		} else {
			config.Logger.Error("getControllerWifiCdes Request failed")
			w.WriteHeader(500)
		}
	})

}
