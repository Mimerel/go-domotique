package daemon

import (
	"go-domotique/extractZway"
	"go-domotique/configuration"
	"go-domotique/heating"
	"go-domotique/models"
	"go-domotique/devices"
	"go-domotique/logger"
	"go-domotique/prowl"
	"time"
)

func Daemon(config *models.Configuration, updateConfig chan bool) {
	logger.Info(config, "Daemon", "Daemon Started")

	for {

		select {
		case <-updateConfig:
			prowl.SendProwlNotification(config, "Domotique", "Application", "updating deamon configuration")
			config = configuration.ReadConfiguration()
		default:
			hour := time.Now().In(config.Location).Hour()
			minute := time.Now().In(config.Location).Minute()
			go extractZway.ExtractZWayMetrics(config)
			if config.Heating.HeatingSettings.Activated {
				go func() {
					err := heating.UpdateHeatingExecute(config)
					if err != nil {
						config.Logger.Error("unable to update heating information")
					}
				}()
			}
			for _, v := range config.Daemon.CronTab {
				if v.Hour == int64(hour) && v.Minute == int64(minute) {
					for _, k := range config.Devices.DevicesTranslated {
						if k.DomotiqueId == v.DomotiqueId {
							go func(){
								err := devices.ExecuteRequest(config, k.ZwaveUrl, k.DeviceId, k.Instance, k.CommandClass, v.Value)
								if err != nil {
									config.Logger.Error("unable to apply cron request device <%s> in value <%v>", k.Name, v.Value)
								}
							}()
							logger.Debug(config, "Daemon", "Putting device <%s> in value <%v> ", k.Name, v.Value)
						}
					}
				}
			}
			time.Sleep(1 * time.Minute)
		}
	}
}