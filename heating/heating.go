package heating

import (
	"go-domotique/logger"
	"go-domotique/models"
	"html/template"
	"net/http"
)

func StatusPage(w http.ResponseWriter, r *http.Request, config *models.Configuration) {
	t := template.New("status.html")
	t, err := t.ParseFiles("./heating/templates/status.html")
	if err != nil {
		logger.Error(config, false,"StatusPage", "Error Parsing template%+v", err)
	}
	logger.DebugPlus(config, false,"StatusPage", "Getting Heating Status")

	data, err := HeatingStatus(config)
	if err != nil {
		config.Logger.Error("Collected Heating status info failed : %v", err)
	}
	data.Devices = config.Devices.DevicesToggle
	config.Channels.MqttCall <- true
	deviceData := <- config.Channels.MqttReceive
	data.DevicesNew = deviceData.Id

	err = t.Execute(w, data)
	if err != nil {
		logger.Error(config, false,"StatusPage", "Error Execution %+v", err)
	}
}

