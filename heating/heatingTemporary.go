package heating

import (
	"fmt"
	"github.com/Mimerel/go-utils"
	"go-domotique/logger"
	"go-domotique/models"
	"go-domotique/prowl"
	"go-domotique/utils"
	"strconv"
	"strings"
	"time"
)

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
		prowl.SendProwlNotification(config, "Domotique", "Application", fmt.Sprintf("Updated Temporary settings till %v, to level %v", config.Heating.TemporaryValues.Moment.Format(time.RFC3339), config.Heating.TemporaryValues.Level))

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
