package controller

import (
	"go-domotique/configuration"
	"go-domotique/daemon"
	"go-domotique/logger"
	"net/http"
)

func Controller() {
	config := configuration.ReadConfiguration()

	logger.Info(config, "Controller", "Application Starting")

	go daemon.Daemon(config)

	getControllerGoogleAssistant(config)

	err := http.ListenAndServe(":9998", nil)
	if err != nil {
		logger.Error(config, "Controller", "error %+v", err)
	}
}