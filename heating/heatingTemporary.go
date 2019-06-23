package heating

import (
	"fmt"
	"github.com/Mimerel/go-utils"
	"go-domotique/models"
	"go-domotique/utils"
	"go-domotique/logger"
	"strconv"
	"strings"
	"time"
)

func SettingTemporaryValues(config *models.Configuration, urlPath string) (err error) {
	utils.GetTimeAndDay(config)
	urlParams := strings.Split(urlPath, "/")
	if len(urlParams) == 3 && strings.ToLower(urlParams[2]) == "reset" {
		config.Heating.TemporaryValues = models.HeatingMoment{}
	} else if len(urlParams) == 4 {
		hours, err := strconv.ParseInt(urlParams[3], 10, 64)
		if err != nil {
			return fmt.Errorf("unable to convert duration string to int64")
		}
		if !go_utils.StringInArray(urlParams[2], []string{"away", "low", "high", "max"}) {
			return fmt.Errorf("Level requested does not exist %s", urlParams[2])
		}
		config.Heating.TemporaryValues.Moment = config.Heating.HeatingMoment.Moment.Local().Add(time.Hour * time.Duration(hours))
		value, err := getValueCorrespondingToLevel(config, urlParams[2])
		config.Heating.TemporaryValues.Level = value
		logger.Info(config, "SettingTemporaryValues", "Updated Temporary settings till %v, to level %v", config.Heating.TemporaryValues.Moment.Format(time.RFC3339), config.Heating.TemporaryValues.Level)
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
