package controller

import (
	"go-domotique/configuration"
	"go-domotique/daemon"
	"go-domotique/logger"
	"go-domotique/prowl"
	"net/http"
	"os"
	"time"
)

func Controller() {
	config := configuration.ReadConfiguration()
	var err error

	logger.Info(config,false,  "Controller", "Application Starting (%v - %v)", time.Now().In(config.Location), time.Now() )
	prowl.SendProwlNotification(config, "Domotique", "Application", "Starting")

	var updateConfig chan bool

	go daemon.Daemon(config, updateConfig)

	heatingController(config)
	getControllerEvents(config)
	getControllerWifiCdes(config)
	getControllerGoogleAssistant(config)
	healthcheckController(config)

	http.HandleFunc("/configuration/update", func(w http.ResponseWriter, r *http.Request) {
		logger.Info(config,false,  "Controller", "Request to update Configuration")
		go prowl.SendProwlNotification(config, "Domotique", "Configuration", "Reloaded")
		w.WriteHeader(200)
		os.Exit(0)
	})


	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		logger.Error(config, true,"Controller", "error %+v", err)
	}
}
