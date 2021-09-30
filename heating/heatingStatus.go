package heating

import (
	"fmt"
	"github.com/Mimerel/go-utils"
	"go-domotique/logger"
	"go-domotique/models"
	"go-domotique/prowl"
	"go-domotique/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func HeatingStatus(config *models.Configuration) (data models.HeatingStatus, err error) {
	utils.GetTimeAndDay(config)

	floatLevel, err := GetInitialHeaterParams(config)
	if err != nil {
		floatLevel = 15
	}

	config.Channels.MqttCall <- true
	deviceData := <- config.Channels.MqttReceive
	data.DevicesNew = deviceData.ToArray()

	data.Heater_Level, data.Temperature_Actual = CollectHeatingStatus(config)

	data.Until = config.Heating.TemporaryValues.Moment
	data.Temperature_Requested = floatLevel

	data.TemporaryLevel = config.Heating.TemporaryValues.Level
	if config.Heating.TemporaryValues.Level != 0 {
		data.IsTemporary = true
	} else {
		data.IsTemporary = false
	}
	if data.Heater_Level > 0 {
		data.IsHeating = true
	} else {
		data.IsHeating = false
	}
	if data.Temperature_Actual >= data.Temperature_Requested {
		data.IsCorrectTemperature = true
	} else {
		data.IsCorrectTemperature = false
	}
	data.IpPort = config.Ip + ":" + config.Port
	data.UpdateTime = config.Heating.LastUpdate
	data.NormalValues = config.Heating.HeatingProgram
	return data, nil
}


func CollectHeatingStatus(config *models.Configuration) (Heater_Level float64, Temperature_Actual float64) {
	config.Channels.MqttCall <- true
	deviceData := <- config.Channels.MqttReceive
	DevicesNew := deviceData.Id

	heaterDevice := DevicesNew[config.Heating.HeatingSettings.HeaterId]
	Heater_Level = heaterDevice.GetStatus()
	Temperature_Actual = DevicesNew[config.Heating.HeatingSettings.SensorId].Temperature
	return Heater_Level, Temperature_Actual
}


func UpdateHeating(w http.ResponseWriter, r *http.Request, config *models.Configuration) (error) {
	err := UpdateHeatingExecute(config)
	if err != nil {
		return err
	}
	return nil
}

func UpdateHeatingExecute(config *models.Configuration) (err error) {
	config.Channels.MqttCall <- true
	deviceData := <- config.Channels.MqttReceive
	DevicesNew := deviceData.Id
	heaterDevice := DevicesNew[config.Heating.HeatingSettings.HeaterId]

	utils.GetTimeAndDay(config)
	config.Heating.LastUpdate = config.Heating.HeatingMoment.Moment
	floatLevel, err := GetInitialHeaterParams(config)
	if err != nil {
		floatLevel = 15
	}
	heater, temperature := CollectHeatingStatus(config)

	activateHeating := CheckIfHeatingNeedsActivating(config, floatLevel, temperature)
	logger.Info(config,false, "UpdateHeatingExecute", "Heating should be activated, %t (%v)", activateHeating, heaterDevice.DomotiqueId)
	if heater == 0 && activateHeating {
		//go wifi.ExecuteRequestRelay( devices.GetDeviceFromId(config, config.Heating.HeatingSettings.HeaterId) ,255, config)
		go RunAction(config, strconv.FormatInt(heaterDevice.DomotiqueId,10), models.ShellyOnOff_0 + "/command", "on")
		logger.Info(config, false, "getActions", "Request succeeded")
	}
	if heater == 255 && !activateHeating {
		//go wifi.ExecuteRequestRelay( devices.GetDeviceFromId(config, config.Heating.HeatingSettings.HeaterId) ,0, config)
		go RunAction(config, strconv.FormatInt(heaterDevice.DomotiqueId,10), models.ShellyOnOff_0 + "/command", "off")
	}
	return nil
}

func SettingTemporaryValues(config *models.Configuration, urlPath string) (err error) {
	utils.GetTimeAndDay(config)
	urlParams := strings.Split(urlPath, "/")
	for k, v := range urlParams {
		logger.Debug(config, false, "SettingTemporaryValues", "UrlParams %v => %v", k, v)
	}
	if len(urlParams) >= 4 && strings.ToLower(urlParams[3]) == "reset" {
		config.Heating.TemporaryValues = models.HeatingMoment{}
	} else if len(urlParams) == 5 {
		hours, err := strconv.ParseInt(urlParams[4], 10, 64)
		if err != nil {
			return fmt.Errorf("unable to convert duration string to int64")
		}
		if !go_utils.StringInArray(urlParams[3], []string{"away", "low", "high", "max"}) {
			return fmt.Errorf("Level requested does not exist %s", urlParams[3])
		}
		config.Heating.TemporaryValues.Moment = config.Heating.HeatingMoment.Moment.In(config.Location).Add(time.Hour * time.Duration(hours))
		value, err := getValueCorrespondingToLevel(config, urlParams[3])
		config.Heating.TemporaryValues.Level = value
		logger.Info(config, false, "SettingTemporaryValues", "Updated Temporary settings till %v, to level %v", config.Heating.TemporaryValues.Moment.Format(time.RFC3339), config.Heating.TemporaryValues.Level)
		go prowl.SendProwlNotification(config, "Domotique", "Application", fmt.Sprintf("Updated Temporary settings till %v, to level %v", config.Heating.TemporaryValues.Moment.Format(time.RFC3339), config.Heating.TemporaryValues.Level))

	} else {
		return fmt.Errorf("Wrong number of parameters sent")
	}
	return nil
}

func getValueCorrespondingToLevel(config *models.Configuration, value string) (result float64, err error) {
	for _, v := range config.Heating.HeatingLevels {
		if v.Name == value {
			return v.Value, nil
		}
	}
	return result, fmt.Errorf("Unable to find corresponding value to heating level demanded")
}

func getLevel(config *models.Configuration) (float64) {
	setLevel := 15.0
	for _, v := range config.Heating.HeatingProgram {
		if v.DayId == int64(config.Heating.HeatingMoment.Weekday) &&
			int(v.Moment) < config.Heating.HeatingMoment.Time {
			setLevel = v.LevelValue
		}
	}
	return setLevel
}

func CheckIfHeatingNeedsActivating(config *models.Configuration, floatLevel float64, temperature float64) bool {
	if temperature <= floatLevel {
		return true
	}
	return false
}

func GetInitialHeaterParams(config *models.Configuration) (floatLevel float64, err error) {
	setLevel := getLevel(config)
	logger.Info(config, false, "GetInitialHeaterParams", "Retreived heating level, %v", setLevel)
	if config.Heating.TemporaryValues.Moment.After(config.Heating.HeatingMoment.Moment) {
		setLevel = config.Heating.TemporaryValues.Level
		logger.Info(config, false, "GetInitialHeaterParams","Temporary heating override, %v", setLevel)
	} else if config.Heating.TemporaryValues.Moment.Before(config.Heating.TemporaryValues.Moment) {
		config.Heating.TemporaryValues = models.HeatingMoment{}
		logger.Info(config, false, "GetInitialHeaterParams", "Clearing old temporary settings")
	}
	return setLevel, nil
}
