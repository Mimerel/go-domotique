package heating

import (
	"go-domotique/devices"
	"go-domotique/logger"
	"go-domotique/models"
	"go-domotique/utils"
	"go-domotique/wifi"
	"net/http"
)

func UpdateHeating(w http.ResponseWriter, r *http.Request, config *models.Configuration) (error) {
	err := UpdateHeatingExecute(config)
	if err != nil {
		return err
	}
	return nil
}

func UpdateHeatingExecute(config *models.Configuration) (err error) {
	utils.GetTimeAndDay(config)
	config.Heating.LastUpdate = config.Heating.HeatingMoment.Moment
	floatLevel, err := GetInitialHeaterParams(config)
	if err != nil {
		floatLevel = 15
	}
	heater, temperature := collectMetrics(config)

	activateHeating := CheckIfHeatingNeedsActivating(config, floatLevel, temperature)
	logger.Info(config,false, "UpdateHeatingExecute", "Heating should be activated, %t", activateHeating)
	if heater == 0 && activateHeating {
		wifi.ExecuteRequestRelay( devices.GetDeviceFromId(config, config.Heating.HeatingSettings.HeaterId) ,255, config)
	}
	if heater == 255 && !activateHeating {
		wifi.ExecuteRequestRelay( devices.GetDeviceFromId(config, config.Heating.HeatingSettings.HeaterId) ,0, config)
	}
	return nil
}