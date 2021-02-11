package heating

import (
	"go-domotique/models"
	"go-domotique/logger"
)

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
	if temperature < floatLevel {
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
