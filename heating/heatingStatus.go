package heating

import (
	"fmt"
	"github.com/Mimerel/go-utils"
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
	config.Channels.MqttGetArray <- true
	data.DevicesNew = <-config.Channels.MqttArray

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
	UpdateRadiatorTarget(config, data.Temperature_Requested)
	data.IpPort = config.Ip + ":" + config.Port
	data.UpdateTime = config.Heating.LastUpdate
	data.NormalValues = config.Heating.HeatingProgram
	data.Rooms = config.Rooms
	return data, nil
}

func UpdateRadiatorTarget(config *models.Configuration, temp_requested float64) {
	radiators := map[int64]models.MqqtDataDetails{}
	sensor := map[int64]models.MqqtDataDetails{}
	// Collecting Sensor Values
	for _, v := range config.Heating.HeatingSettings {
		if v.Module == "radiator" {
			config.Channels.MqttDomotiqueIdGet <- v.DomotiqueId
			radiators[v.DomotiqueId] = <-config.Channels.MqttDomotiqueDeviceGet
		}
		if v.Module == "sensor" {
			config.Channels.MqttDomotiqueIdGet <- v.DomotiqueId
			sensor[v.DomotiqueId] = <-config.Channels.MqttDomotiqueDeviceGet
		}
	}
	// Setting radiator temperatures
	for _, v := range config.Heating.HeatingSettings {
		if v.Module == "sensor" && v.RadiatorId != 0 {
			if sensor[v.DomotiqueId].Temperature != 0 {
				swapper := radiators[v.RadiatorId]
				swapper.Temperature = sensor[v.DomotiqueId].Temperature
				radiators[v.RadiatorId] = swapper
				httpParams := new(utils.HttpRequestParams)
				httpParams.Debug = false
				httpParams.Method = http.MethodGet
				httpParams.Headers = make(map[string]string)
				httpParams.Headers["Accept"] = "application/json"
				httpParams.Url = "http://192.168.222." + strconv.FormatInt(radiators[v.RadiatorId].DeviceId, 10) + "/ext_t?temp=" + strconv.FormatFloat(sensor[v.DomotiqueId].Temperature, 'f', 2, 32)
				httpParams.Timeout = 120
				httpParams.Retry = 0
				httpParams.DelayBetweenRetry = 5
				//httpParams.Proxy = proxy
				go func(idRadiator int64) {
					err, response := utils.HttpExecuteRequest(config, httpParams)
					if err != nil {
						config.Logger.Error("Unable to execute request for updates : %v , %+v", httpParams.Url, err)
						//return data, err
					} else {
						config.Logger.DebugPlus("Response code for device %v = %v", idRadiator, response.StatusCode)
					}
				}(v.RadiatorId)
			}
		}
	}
	// managing radiator valves
	for id, v := range radiators {
		config.Logger.Info("Device %v (%v) is at %v temperature for requested %v", v.Name, v.DeviceId, v.Temperature, temp_requested)
		if v.Temperature < temp_requested && v.CurrentPos <= 99 {
			go RunAction(config, strconv.FormatInt(id, 10), "/thermostat/0/command", "valve_pos="+strconv.FormatFloat(100, 'f', 2, 32))
			config.Logger.Info("Changing valve position to fully open (100) ")
		}
		if v.Temperature >= temp_requested && v.CurrentPos > 10 {
			go RunAction(config, strconv.FormatInt(id, 10), "/thermostat/0/command", "valve_pos="+strconv.FormatFloat(0, 'f', 2, 32))
			config.Logger.Info("Changing valve position to fully closed (0)")
		}
	}
}

func CollectHeatingStatus(config *models.Configuration) (Heater_Level float64, Temperature_Actual float64) {
	var heaterDevice models.MqqtDataDetails
	Temperature_Actual = 999
	for _, v := range config.Heating.HeatingSettings {
		config.Channels.MqttDomotiqueIdGet <- v.DomotiqueId
		devTemp := <-config.Channels.MqttDomotiqueDeviceGet

		if v.Module == "heater" {
			heaterDevice = devTemp
		}
		if v.Module == "sensor" || v.Module == "radiator" {
			devTemp := devTemp.Temperature
			if devTemp == 0 {
				devTemp = 999
			}

			if devTemp < Temperature_Actual {
				Temperature_Actual = devTemp
			}
		}
	}
	Heater_Level = heaterDevice.GetStatus()
	//deviceData.Unlock()

	config.Logger.Info("Heating Sensor value : temp : %v", Temperature_Actual)
	return Heater_Level, Temperature_Actual
}

func UpdateHeating(w http.ResponseWriter, r *http.Request, config *models.Configuration) error {
	err := UpdateHeatingExecute(config)
	if err != nil {
		return err
	}
	return nil
}

func UpdateHeatingExecute(config *models.Configuration) (err error) {
	var heaterDevice models.MqqtDataDetails
	for _, v := range config.Heating.HeatingSettings {
		config.Channels.MqttDomotiqueIdGet <- v.DomotiqueId
		devTemp := <-config.Channels.MqttDomotiqueDeviceGet
		if v.Module == "heater" {
			heaterDevice = devTemp
		}
	}
	utils.GetTimeAndDay(config)
	config.Heating.LastUpdate = config.Heating.HeatingMoment.Moment
	floatLevel, err := GetInitialHeaterParams(config)
	if err != nil {
		floatLevel = 15
	}
	heater, temperature := CollectHeatingStatus(config)

	activateHeating := CheckIfHeatingNeedsActivating(config, floatLevel, temperature)
	config.Logger.Info("UpdateHeatingExecute Heating should be activated, %t (%v)", activateHeating, heaterDevice.DomotiqueId)
	if heater == 0 && activateHeating {
		go RunAction(config, strconv.FormatInt(heaterDevice.DomotiqueId, 10), models.ShellyOnOff_0+"/command", "on")
		config.Logger.Info("getActions Request succeeded")
	}
	if heater == 255 && !activateHeating {
		go RunAction(config, strconv.FormatInt(heaterDevice.DomotiqueId, 10), models.ShellyOnOff_0+"/command", "off")
	}
	UpdateRadiatorTarget(config, floatLevel)
	return nil
}

func SettingTemporaryValues(config *models.Configuration, urlPath string) (err error) {
	utils.GetTimeAndDay(config)
	urlParams := strings.Split(urlPath, "/")
	for k, v := range urlParams {
		config.Logger.DebugPlus("SettingTemporaryValues UrlParams %v => %v", k, v)
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
		config.Logger.Info("SettingTemporaryValues Updated Temporary settings till %v, to level %v", config.Heating.TemporaryValues.Moment.Format(time.RFC3339), config.Heating.TemporaryValues.Level)
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

func getLevel(config *models.Configuration) float64 {
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
	if temperature < floatLevel {
		return true
	}
	return false
}

func GetInitialHeaterParams(config *models.Configuration) (floatLevel float64, err error) {
	setLevel := getLevel(config)
	config.Logger.Info("GetInitialHeaterParams Retreived heating level, %v", setLevel)
	if config.Heating.TemporaryValues.Moment.After(config.Heating.HeatingMoment.Moment) {
		setLevel = config.Heating.TemporaryValues.Level
		config.Logger.Info("GetInitialHeaterParams Temporary heating override, %v", setLevel)
	} else if config.Heating.TemporaryValues.Moment.Before(config.Heating.TemporaryValues.Moment) {
		config.Heating.TemporaryValues = models.HeatingMoment{}
		config.Logger.Info("GetInitialHeaterParams Clearing old temporary settings")
	}
	return setLevel, nil
}
