package heating

import (
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
	data.IpPort = config.Ip + ":" + config.Port
	data.UpdateTime = config.Heating.LastUpdate
	data.NormalValues = config.Heating.HeatingProgram
	return data, nil
}

func collectMetrics(config *models.Configuration) (heater float64, temperature float64) {
	found := 0
	for _, v := range config.Devices.LastValues {
		if v.DomotiqueId == config.Heating.HeaterId {
			heater = v.Value
			found += 1
		}
		if v.DomotiqueId == config.Heating.SensorId &&
			v.Unit == "Degré" &&
			v.InstanceId == 0 {
			temperature = v.Value
			found += 1
		}
	}
	config.Logger.Info("Metrics retrieved, heater %f , temperature %f", heater, temperature)
	return heater, temperature
}
