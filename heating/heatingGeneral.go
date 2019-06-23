package heating

import (
	"go-domotique/models"
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
	config.Logger.Info("Retreived heating level, %v", setLevel)
	if config.Heating.TemporaryValues.Moment.After(config.Heating.HeatingMoment.Moment) {
		setLevel = config.Heating.TemporaryValues.Level
		config.Logger.Info("Temporary heating override, %v", setLevel)
	} else if config.Heating.TemporaryValues.Moment.Before(config.Heating.TemporaryValues.Moment) {
		config.Heating.TemporaryValues = models.HeatingMoment{}
		config.Logger.Info("Clearing old temporary settings")
	}
	return setLevel, nil
}
