package heating

import (
	"fmt"
	"github.com/Mimerel/go-utils"
	"go-domotique/models"
	"go-domotique/utils"
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
		// TODO : get corresponding value to demanded value => directly in mysql table
		config.Heating.TemporaryValues.Level = 15 // replace by equivalent of urlParams[2]
		config.Logger.Info("Updated Temporary settings till %v, to level %v", config.Heating.TemporaryValues.Moment.Format(time.RFC3339), config.Heating.TemporaryValues.Level)
	} else {
		return fmt.Errorf("Wrong number of parameters sent")
	}
	return nil
}
