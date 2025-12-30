package heating

import (
	"encoding/json"
	"go-domotique/TemplateGlobal"
	"go-domotique/models"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

var statusTemplate *template.Template
var confirmationTemplate *template.Template

func InitTemplates() error {
	files := []string{
		"./heating/templates/Roller.html",
		"./heating/templates/Switch.html",
		"./heating/templates/Radiator.html",
		"./heating/templates/Shelly4PM.html",
		"./heating/templates/ShellyHT.html",
		"./heating/templates/ShellyI3.html",
		"./heating/templates/ShellyUni.html",
		"./heating/templates/status.html",
	}
	t := template.New("status.html")
	t = template.Must(t.Funcs(TemplateGlobal.GetUIDict()).Funcs(TemplateGlobal.GetUIFunctions()).ParseGlob("./heating/templates/*.html"))
	var err error
	statusTemplate, err = t.ParseFiles(files...)
	if err != nil {
		return err
	}
	confirmationTemplate, err = template.New("confirmation.html").ParseFiles("./heating/templates/confirmation.html")
	return err
}

func GetConfirmationTemplate() *template.Template {
	return confirmationTemplate
}

func StatusPage(w http.ResponseWriter, r *http.Request, config *models.Configuration) {
	if statusTemplate == nil {
		err := InitTemplates()
		if err != nil {
			config.Logger.Error("Error Parsing template: %+v", err)
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}
	}
	config.Logger.DebugPlus("Getting Heating Status")

	data, err := HeatingStatus(config)
	if err != nil {
		config.Logger.Error("Collected Heating status info failed : %v", err)
	}
	data.Devices = config.Devices.DevicesToggle

	config.Channels.MqttGetArray <- true
	data.DevicesNew = <-config.Channels.MqttArray

	for _, v := range data.DevicesNew {
		data.Totals.Watts += v.Power
	}
	data.Totals.Watts = math.Round(data.Totals.Watts*100) / 100
	err = statusTemplate.Execute(w, data)
	if err != nil {
		config.Logger.Error("Error Execution %+v", err)
	}
}

func Update(w http.ResponseWriter, r *http.Request, config *models.Configuration) {
	config.Channels.MqttGetArray <- true
	//fmt.Printf("request sent Asked...\n")
	deviceList := <-config.Channels.MqttArray
	//fmt.Printf("response received Asked...%v \n", deviceList)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(deviceList)
	if err != nil {
		config.Logger.Error("Unable to convert response to json %v", err)
	}

	//config.Logger.Info("Finished sending updates")
	return
}

func RunAction(config *models.Configuration, idString string, action string, payload string) {
	id, err := strconv.Atoi(idString)
	if err != nil {
		config.Logger.Error("Error Converting Id to int64 <%v> : %v", idString, err)
		return
	}
	actionParams := models.MqttSendMessage{
		DomotiqueId: int64(id),
		Topic:       action,
		Payload:     payload,
	}
	config.Logger.Warn("request ; %v", actionParams)
	config.Channels.MqttSend <- actionParams
}
