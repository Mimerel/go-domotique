package devices

import (
	"go-domotique/logger"
	"go-domotique/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)


func GetDeviceFromId(config *models.Configuration, id int64) (models.DeviceTranslated) {
	for _, v := range config.Devices.DevicesTranslated {
		if v.DomotiqueId == id {
			return v
		}
	}
	return models.DeviceTranslated{}
}

func GetDomotiqueIdFromDeviceIdAndBoxId(config *models.Configuration, deviceId int64, ZwaveId int64) (models.DeviceTranslated) {
	for _, v := range config.Devices.DevicesTranslated {
		if v.DeviceId == deviceId && v.Zwave == ZwaveId{
			return v
		}
	}
	return models.DeviceTranslated{}
}

func GetZwaveIdFromZwaveName(config *models.Configuration, name string) (models.Zwave) {
	for _, v := range config.Zwaves {
		if strings.ToUpper(v.Name) == strings.ToUpper(name) {
			return v
		}
	}
	return models.Zwave{}
}



func ExecuteAction(config *models.Configuration, instruction models.GoogleTranslatedInstruction) (hasError bool) {
	err := ExecuteRequest(config, instruction.ZwaveUrl, instruction.DeviceId, instruction.Instance, instruction.CommandClass, instruction.Value)
	if err != nil {
		return true
	}
	return false
}

func ExecuteActionDomotiqueId(config *models.Configuration, domotiqueId int64, value int64) (err error) {
	device := GetDeviceFromId(config, domotiqueId)
	err = ExecuteRequest(config, device.ZwaveUrl, device.DeviceId, device.Instance, device.CommandClass, value)
	return nil
}

/**
Method that sends a request to a domotic zwave server to run an instruction
 */
func ExecuteRequest(config *models.Configuration, url string, id int64, instance int64, commandClass int64, level int64) (err error) {
	logger.Info(config, false, "ExecuteRequest", "Préparing post")
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	postingUrl := "http://" + url + ":8083/ZWaveAPI/Run/devices[" + strconv.FormatInt(id, 10) + "].instances[" + strconv.FormatInt(instance, 10) + "].commandClasses[" + strconv.FormatInt(commandClass, 10) + "].Set(" + strconv.FormatInt(level, 10) + ")"
	logger.Info(config, false, "ExecuteRequest", "Request posted : %s", postingUrl)

	_, err = client.Get(postingUrl)
	if err != nil {
		logger.Error(config, true,"ExecuteRequest", "Failed to execute request %s ", postingUrl, err)
		return err
	}
	logger.Info(config, false, "ExecuteRequest", "Request successful...")
	return nil
}


