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
		logger.Error(config, true,"StatusPage", "Error Parsing template%+v", err)
	}
	data, err := HeatingStatus(config)
	if err != nil {
		config.Logger.Error("Collected Heating status info failed : %v", err)
	}
	data.Devices = config.Devices.DevicesToggle
	data.GetLastValuesForDevice(config)
	err = t.Execute(w, data)
	if err != nil {
		logger.Error(config, true,"StatusPage", "Error Execution %+v", err)
	}
}

