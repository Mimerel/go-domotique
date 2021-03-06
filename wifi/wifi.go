package wifi

import (
	"go-domotique/configuration"
	"go-domotique/devices"
	"go-domotique/logger"
	"go-domotique/models"
	"net/http"
	"strconv"
	"time"
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

	//ExecuteRequest(concernedDevice, action, config)
}




func ExecuteRequestRelay(concernedDevice models.DeviceTranslated, value int64, config *models.Configuration) {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	postingUrl := concernedDevice.GetUrlForValue(config, value)

	logger.Info(config, false, "ExecuteRequest", "Request posted : %s", postingUrl)

	_, err := client.Get(postingUrl)
	if err != nil {
		logger.Error(config, true, "ExecuteRequest", "Failed to execute request %s ", postingUrl, err)
		return
	}
	//logger.Info(config, "ExecuteRequest", "Request successful...")
}

func WifiTranslateValue(value int64) string {
	if value == 0 {
		return configuration.ETEINS
	}
	return configuration.ALLUME
}
