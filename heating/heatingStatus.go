package heating

import (
	"go-domotique/devices"
	"go-domotique/logger"
	"go-domotique/models"
	"go-domotique/utils"
)

func HeatingStatus(config *models.Configuration) (data models.HeatingStatus, err error) {
	utils.GetTimeAndDay(config)

	floatLevel, err := GetInitialHeaterParams(config)
	if err != nil {
		floatLevel = 15
	}
	heater, temperature := collectMetrics(config)

	data.Until = config.Heating.TemporaryValues.Moment
	data.Temperature_Actual = temperature
	data.Temperature_Requested = floatLevel
	data.Heater_Level = heater
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

func collectMetrics(config *models.Configuration) (heater float64, temperature float64) {
	err := utils.GetLastDeviceValues(config)
	if err != nil {
		logger.Error(config, true,"collectMetrics", "unable to read device values", err)
		return
	}
/*	for _, v := range config.Devices.LastValues {
		//if v.DomotiqueId == config.Heating.HeatingSettings.HeaterId {
		//	heater = v.Value
		//}
		if v.DomotiqueId == config.Heating.HeatingSettings.SensorId &&
			v.Unit == "Degré" &&
			v.InstanceId == 0 {
			temperature = v.Value
		}
	}
*/
	logger.Info(config, false, "collectMetrics","Heater id %v", config.Heating.HeatingSettings.HeaterId)
	logger.Info(config, false, "collectMetrics","Found device %v", devices.GetDeviceFromId(config, config.Heating.HeatingSettings.HeaterId))
	status :=  models.GetStatusWifi(config, devices.GetDeviceFromId(config, config.Heating.HeatingSettings.HeaterId).DeviceId)
	logger.Info(config, false, "collectMetrics","Status %v", status)


	logger.Info(config, false, "collectMetrics","Sensor id %v", config.Heating.HeatingSettings.SensorId)
	logger.Info(config, false, "collectMetrics","Found device %v", devices.GetDeviceFromId(config, config.Heating.HeatingSettings.SensorId))
	statusTemperature :=  models.GetStatusWifi(config, devices.GetDeviceFromId(config, config.Heating.HeatingSettings.SensorId).DeviceId)
	logger.Info(config, false, "collectMetrics","Status %v", statusTemperature)



	logger.Info(config, false, "collectMetrics", "Heating wifi status : %v - %v" , status.Power, statusTemperature.Temperature)

	heater = 0
	heaterstate := "Off"
	if status.Value {
		heater = 255
		heaterstate = "On"
	}
	logger.Info(config, false, "collectMetrics", "Metrics retrieved, heater %v , temperature %f", heaterstate, statusTemperature.Temperature)
	return heater, statusTemperature.Temperature
}


