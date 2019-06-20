package devices

import (
	"go-domotique/logger"
	"go-domotique/models"
	"net/http"
	"strconv"
	"time"
)



func GetDeviceFromId(c *models.Configuration, id int64) (models.DeviceTranslated) {
	for _, v := range c.Devices.DevicesTranslated {
		if v.DomotiqueId == id {
			return v
		}
	}
	return models.DeviceTranslated{}
}

func GetDomotiqueIdFromDeviceIdAndBoxId(c *models.Configuration, deviceId int64, boxId int64) (models.DeviceTranslated) {
	for _, v := range c.Devices.DevicesTranslated {
		if v.DeviceId == deviceId && v.Zwave == boxId{
			return v
		}
	}
	return models.DeviceTranslated{}
}


func ExecuteAction(config *models.Configuration, instruction models.GoogleTranslatedInstruction) (hasError bool) {
	err := ExecuteRequest(config, instruction.ZwaveUrl, instruction.DeviceId, instruction.Instance, instruction.CommandClass, instruction.Value)
	if err != nil {
		return true
	}
	return false
}


/**
Method that sends a request to a domotic zwave server to run an instruction
 */
func ExecuteRequest(config *models.Configuration, url string, id int64, instance int64, commandClass int64, level int64) (err error) {
	logger.Info(config, "ExecuteRequest", "Préparing post")
	timeout := time.Duration(20 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	postingUrl := "http://" + url + ":8083/ZWaveAPI/Run/devices[" + strconv.FormatInt(id, 10) + "].instances[" + strconv.FormatInt(instance, 10) + "].commandClasses[" + strconv.FormatInt(commandClass, 10) + "].Set(" + strconv.FormatInt(level, 10) + ")"
	logger.Info(config, "ExecuteRequest", "Request posted : %s", postingUrl)

	_, err = client.Get(postingUrl)
	if err != nil {
		logger.Error(config, "ExecuteRequest", "Failed to execute request %s ", postingUrl, err)
		return err
	}
	logger.Info(config, "ExecuteRequest", "Request successful...")
	return nil
}

