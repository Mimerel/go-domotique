package daemon

import (
	"go-domotique/configuration"
	"go-domotique/devices"
	"go-domotique/extractZway"
	"go-domotique/heating"
	"go-domotique/logger"
	"go-domotique/models"
	"go-domotique/prowl"
	"go-domotique/utils"
	"go-domotique/wifi"
	"os"
	"strconv"
	"time"
)

func Daemon(config *models.Configuration, updateConfig chan bool) {
	logger.Info(config, false, "Daemon", "Daemon Started")

	for {

		select {
		case <-updateConfig:
			prowl.SendProwlNotification(config, "Domotique", "Application", "updating deamon configuration")
			config = configuration.ReadConfiguration()
			os.Exit(0)
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
				//if skipCronInstruction(v, config) == true {
				//	continue
				//}
				if v.Hour == int64(hour) && v.Minute == int64(minute) {
					for _, k := range config.Devices.DevicesTranslated {
						if k.DomotiqueId == v.DomotiqueId {
							go func() {
								if v.ProwlIt {
									prowl.SendProwlNotification(config, "Domotique", "Application", "Starting")
								}
								switch k.Zwave {
								case 100:
									logger.Info(config, false, "RunDomoticCommand", "CRON Running Wifi instruction : %+v, %+v", k.DeviceId, k.Type)
									wifi.ExecuteRequestRelay(strconv.Itoa(int(k.DeviceId)), wifi.WifiTranslateValue(v.Value), config)
								default:
									logger.Info(config, false, "RunDomoticCommand", "CRON Running Zwave instruction")
									err := devices.ExecuteRequest(config, k.ZwaveUrl, k.DeviceId, k.Instance, k.CommandClass, v.Value)
									if err != nil {
										config.Logger.Error("unable to apply cron request device <%s> in value <%v>", k.Name, v.Value)
									}
								}
							}()
						}
					}
				}
			}
			time.Sleep(1 * time.Minute)
		}
	}
}

func skipCronInstruction(v models.CronTab, config *models.Configuration) bool {
	if (v.NotOnAway == 1 && config.Heating.TemporaryValues.Level == 3) ||
		(v.NotOnAlarmTotal == 1 && utils.GetLastDeviceValue(config, 74, 1).Value == 255) {
		return true
	}
	return false
}
