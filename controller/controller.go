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

	getControllerGoogleAssistant(config)
	heatingController(config)

	err := http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		logger.Error(config, "Controller", "error %+v", err)
	}
}
