package heating

import (
	"go-domotique/logger"
	"go-domotique/models"
	"go-domotique/TemplateGlobal"
	"html/template"
	"net/http"
)

func StatusPage(w http.ResponseWriter, r *http.Request, config *models.Configuration) {
	t := template.New("status.html")
	t = template.Must(t.Funcs(TemplateGlobal.GetUIDict()).Funcs(TemplateGlobal.GetUIFunctions()).ParseGlob("./heating/templates/*.html"))
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
	data.DevicesNew = deviceData.ToArray()
	data.Totals.Watts = deviceData.TotalWatts

	err = t.Execute(w, data)
	if err != nil {
		logger.Error(config, false,"StatusPage", "Error Execution %+v", err)
	}
}

