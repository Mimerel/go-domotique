package wifi

import (
	"go-domotique/configuration"
	"go-domotique/devices"
	"go-domotique/models"
	"net/http"
	"strconv"
)

func AnalyseRequest(w http.ResponseWriter, r *http.Request, urlParams []string, config *models.Configuration) {
	//logger.Info(config, false, "AnalyseRequest", "Analysing wifi request")
	emittingDevice := urlParams[2]
	actionStatus := urlParams[3]
	emittingDeviceInt, err := strconv.Atoi(emittingDevice)
	if err != nil {
		config.Logger.Error("Unable to convert %v to int", emittingDevice)
		return
	}
	//logger.Info(config, false, "AnalyseRequest", "Recevied action from domotiqueId %+v", emittingDevice)

	for _, k := range config.Devices.DevicesActions {
		//logger.Info(config, false, "AnalyseRequest", "found wifi actions %v", k)
		if k.DomotiqueId == int64(emittingDeviceInt) {
			//logger.Info(config, false, "AnalyseRequest", "Found wifi action %+v", k)
			for _, googleAction := range config.GoogleAssistant.GoogleTranslatedInstructions {
				if googleAction.ActionNameId == k.ActionNameId && googleAction.Type == actionStatus {
					value := int64(0)
					if actionStatus == configuration.ALLUME {
						value = 255
					}
					go ExecuteRequestRelay(devices.GetDeviceFromId(config, k.DomotiqueId), value, config)
					//logger.Info(config, false, "AnalyseRequest", "Updating device %v", googleAction.DeviceId)
				}
			}
		}
	}
}

func ExecuteRequestRelay(concernedDevice models.DeviceTranslated, value int64, config *models.Configuration) {
	topic, payload := concernedDevice.GetUrlForValue(config, value)
	actionParams := models.MqttSendMessage{
		DomotiqueId: int64(concernedDevice.DomotiqueId),
		Topic: topic,
		Payload: payload,
	}
	config.Channels.MqttSend <- actionParams
}
