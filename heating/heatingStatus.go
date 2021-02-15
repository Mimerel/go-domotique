package heating

import (
	"go-domotique/models"
	"go-domotique/utils"
	"go-domotique/logger"
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
	data.IpPort = config.Ip + ":" + config.Port
	data.UpdateTime = config.Heating.LastUpdate
	data.NormalValues = config.Heating.HeatingProgram
	return data, nil
}

func collectMetrics(config *models.Configuration) (heater float64, temperature float64) {
	found := 0
	err := utils.GetLastDeviceValues(config)
	if err != nil {
		logger.Error(config, true,"collectMetrics", "unable to read device values", err)
		return
	}
	for _, v := range config.Devices.LastValues {
		if v.DomotiqueId == config.Heating.HeatingSettings.HeaterId {
			heater = v.Value
			found += 1
		}
		if v.DomotiqueId == config.Heating.HeatingSettings.SensorId &&
			v.Unit == "Degré" &&
			v.InstanceId == 0 {
			temperature = v.Value
			found += 1
		}
	}
	logger.Info(config, false, "collectMetrics", "Metrics retrieved, heater %f , temperature %f", heater, temperature)
	return heater, temperature
}


