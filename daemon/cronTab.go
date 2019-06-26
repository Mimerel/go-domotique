package daemon

import (
	"go-domotique/extractZway"
	"go-domotique/heating"
	"go-domotique/models"
	"go-domotique/devices"
	"go-domotique/logger"
	"time"
)

func Daemon(config *models.Configuration) {
	logger.Info(config, "Daemon", "Daemon Started")

	for {
		hour := time.Now().In(config.Location).Hour()
		minute := time.Now().In(config.Location).Minute()
		go extractZway.ExtractZWayMetrics(config)
		if config.Heating.HeatingSettings.Activated {
			go heating.UpdateHeatingExecute(config)
		}
		for _, v := range config.Daemon.CronTab {
			if v.Hour == int64(hour) && v.Minute == int64(minute) {
				for _,k := range config.Devices.DevicesTranslated {
					if k.DomotiqueId == v.DomotiqueId {
						go devices.ExecuteRequest(config, k.ZwaveUrl, k.DeviceId, k.Instance, k.CommandClass, v.Value)
						logger.Debug(config, "Daemon", "Putting device <%s> in value <%v> ", k.Name, v.Value)
					}
				}
			}
		}
		time.Sleep(1 * time.Minute)
	}
}