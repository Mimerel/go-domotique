package controller

import (
	"go-domotique/configuration"
	"go-domotique/daemon"
	"go-domotique/logger"
	"go-domotique/prowl"
	"net/http"
)

func Controller() {
	config := configuration.ReadConfiguration()

	logger.Info(config, "Controller", "Application Starting")
	prowl.SendProwlNotification(config, "Domotique", "Application", "Starting")

	go daemon.Daemon(config)

	heatingController(config)
	getControllerEvents(config)
	getControllerGoogleAssistant(config)

	http.HandleFunc("/configuration/update", func(w http.ResponseWriter, r *http.Request) {
		logger.Info(config, "Controller", "Request to update Configuration")
		configuration.ReadConfiguration()
		w.WriteHeader(200)
	})


	err := http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		logger.Error(config, "Controller", "error %+v", err)
	}
}
