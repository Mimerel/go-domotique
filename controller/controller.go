package controller

import (
	"go-domotique/configuration"
	"go-domotique/daemon"
	"go-domotique/models"
	"go-domotique/prowl"
	"net/http"
	"os"
	"time"
)

func Controller() {
	config := configuration.ReadConfiguration()
	var err error

	config.Logger.Info("Controller Application Starting (%v - %v)", time.Now().In(config.Location), time.Now())
	go prowl.SendProwlNotification(config, "Domotique", "Application", "Starting")

	config.Channels.UpdateConfig = make(chan bool)
	//config.Channels.MqttCall = make(chan bool)
	//config.Channels.MqttReceive = make(chan models.MqqtData)
	config.Channels.MqttSend = make(chan models.MqttSendMessage)
	config.Channels.MqttReconnect = make(chan bool)
	config.Channels.MqttDomotiqueIdGet = make(chan int64)
	config.Channels.MqttDomotiqueDeviceGet = make(chan models.MqqtDataDetails)
	config.Channels.MqttDomotiqueDevicePost = make(chan models.MqqtDataDetails, 10000)
	config.Channels.MqttGetArray = make(chan bool)
	config.Channels.MqttArray = make(chan []models.MqqtDataDetails)

	go daemon.Daemon(config)
	//config.Channels.MqttReconnect <- false

	heatingController(config)
	getControllerEvents(config)
	getActions(config)
	getControllerWifiCdes(config)
	getControllerGoogleAssistant(config)
	healthcheckController(config)
	getControllerWebHooks(config)

	http.HandleFunc("/configuration/update", func(w http.ResponseWriter, r *http.Request) {
		config.Logger.Info("Controller Request to update Configuration")
		go prowl.SendProwlNotification(config, "Domotique", "Configuration", "Reloaded")
		w.WriteHeader(200)
		os.Exit(0)
	})

	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		config.Logger.Error("Controller error %+v", err)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Add("Vary", "Origin")
	(*w).Header().Add("Vary", "Access-Control-Request-Method")
	(*w).Header().Add("Vary", "Access-Control-Request-Headers")
	(*w).Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token, Content-Encoding, x-csrf-token")
	(*w).Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
}
