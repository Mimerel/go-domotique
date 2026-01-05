package daemon

import (
	"go-domotique/models"
	"go-domotique/utils"
	"net/http"
	"strconv"
	"time"
)

// rebootTrvDevices runs daily at 1am and sends a reboot command to all TRV devices
func rebootTrvDevices(config *models.Configuration) {
	for {
		hour := time.Now().In(config.Location).Hour()
		minute := time.Now().In(config.Location).Minute()

		if hour == 1 && minute == 0 {
			config.Logger.Info("Starting daily TRV reboot cycle")

			for _, v := range config.Heating.HeatingSettings {
				if v.Module == "radiator" {
					config.Channels.MqttDomotiqueIdGet <- v.DomotiqueId
					device := <-config.Channels.MqttDomotiqueDeviceGet

					ip := device.Ip
					if ip == "" {
						// Fallback to constructed IP if not available from MQTT
						ip = "192.168.222." + strconv.FormatInt(device.DeviceId, 10)
					}

					config.Logger.Info("Rebooting TRV device %v (%s) at IP %s", v.DomotiqueId, device.Name, ip)

					httpParams := new(utils.HttpRequestParams)
					httpParams.Debug = false
					httpParams.Method = http.MethodGet
					httpParams.Url = "http://" + ip + "/reboot"
					httpParams.Timeout = 30
					httpParams.Retry = 1
					httpParams.DelayBetweenRetry = 5

					go func(domotiqueId int64, deviceName string) {
						err, _ := utils.HttpExecuteRequest(config, httpParams)
						if err != nil {
							config.Logger.Error("Failed to reboot TRV device %v (%s): %v", domotiqueId, deviceName, err)
						} else {
							config.Logger.Info("Successfully sent reboot command to TRV device %v (%s)", domotiqueId, deviceName)
						}
					}(v.DomotiqueId, device.Name)

					// Small delay between reboots to avoid overwhelming the network
					time.Sleep(5 * time.Second)
				}
			}

			// Sleep for 2 minutes to avoid re-triggering during the same minute window
			time.Sleep(2 * time.Minute)
		}

		time.Sleep(30 * time.Second)
	}
}
