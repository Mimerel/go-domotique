package controller

import (
	"fmt"
	"go-goole-home-requests/configuration"
	"go-goole-home-requests/daemon"
	"go-goole-home-requests/logger"
	"net/http"
)

func Controller() {
	config := configuration.ReadConfiguration()

	logger.Info(config, "Application", "Application Starting")

	go daemon.Daemon(config)

	getControllerGoogleAssistant(config)

	err := http.ListenAndServe(":9998", nil)
	if err != nil {
		fmt.Printf("error %+v", err)
	}
}