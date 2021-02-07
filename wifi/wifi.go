package wifi

import (
	"go-domotique/models"
	"go-domotique/logger"
	"go-domotique/devices"
	"net/http"
	"strconv"
	"time"
)


func AnalyseRequest(w http.ResponseWriter, r *http.Request, urlParams []string, config *models.Configuration) {
	concernedAction := urlParams[2]
	actionStatus := urlParams[3]
	concernedActionInt , err := strconv.Atoi(concernedAction)
	if err != nil {
		config.Logger.Error("Unable to convert %v to int", concernedAction)
		return
	}
	actionStatusInt , err := strconv.Atoi(actionStatus)
	if err != nil {
		config.Logger.Error("Unable to convert %v to int", concernedAction)
		return
	}

	device := devices.GetDomotiqueIdFromDeviceIdAndBoxId(config, int64(concernedActionInt) , 100)
	for _, k := range config.Devices.DevicesActions {
		if k.DomotiqueId == device.DomotiqueId {
			for _, googleAction := range config.GoogleAssistant.GoogleInstructions {
				if googleAction.ActionNameId == k.ActionNameId && googleAction.TypeId == int64(actionStatusInt){
					ExecuteRequestRelay(strconv.Itoa(int(googleAction.DomotiqueId)), actionStatus, config)
				}
			}
		}
	}

	//ExecuteRequest(concernedDevice, action, config)
}

func ExecuteRequestRelay (concernedDevice string, action string, config *models.Configuration ){
	logger.Info(config, "ExecuteRequest", "Préparing post")
	timeout := time.Duration(20 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	postingUrl := "http://" + config.Ip[:12] + concernedDevice + "/relay/0?turn=" + action
	logger.Info(config, "ExecuteRequest", "Request posted : %s", postingUrl)

	_, err := client.Get(postingUrl)
	if err != nil {
		logger.Error(config, "ExecuteRequest", "Failed to execute request %s ", postingUrl, err)
		return
	}
	logger.Info(config, "ExecuteRequest", "Request successful...")
}

