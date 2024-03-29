package daemon

import (
	"fmt"
	"go-domotique/models"
	"strconv"
	"time"
)

var Devices []models.MqqtDataDetails

func Mqtt_Deamon(c *models.Configuration) {
	mqttConfig = c
	logger := mqttConfig.Logger
	Devices = []models.MqqtDataDetails{}
	queue := []models.MqqtDataDetails{}
	reconnect(true)
	var initial bool
	defer Client.Disconnect(100)
	for {
		select {

		case <-mqttConfig.Channels.MqttGetArray:
			mqttConfig.Channels.MqttArray <- Devices
			break

		case id := <-mqttConfig.Channels.MqttDomotiqueIdGet:
			mqttConfig.Channels.MqttDomotiqueDeviceGet <- getDevice(id)
			break

		case deviceUpdate := <-mqttConfig.Channels.MqttDomotiqueDevicePost:
			queue = append(queue, deviceUpdate)
			break

		case value := <-mqttConfig.Channels.MqttReconnect:
			reconnectUpdate(value)
			break

		case mqttAction := <-mqttConfig.Channels.MqttSend:
			actionToDo := models.Prefix + strconv.FormatInt(mqttAction.DomotiqueId, 10) + mqttAction.Topic
			logger.Debug("Action : %s - payload %v", actionToDo, mqttAction.Payload)
			token := Client.Publish(actionToDo, 0, false, mqttAction.Payload)
			switch mqttAction.Topic + fmt.Sprint("%v", mqttAction.Payload) {
			case "/commandupdate_fw":
				go updateAnnounce(mqttAction.DomotiqueId)
				break
			}
			token.Wait()
			break
		default:
			if initial == false {
				initial = true
				reconnectUpdate(false)
			}
			if len(queue) > 0 {
				list := make(map[int64]int)
				for k, v := range Devices {
					list[v.DomotiqueId] = k
				}
				if len(queue) > 10000 {
					queue = []models.MqqtDataDetails{}
				}
				processing := []models.MqqtDataDetails{}
				if len(queue) > 100 {
					logger.Debug("Queue length to process : %v\n", len(queue))
					processing = queue[0:99]
					queue = queue[100:]
				} else {
					processing = queue
					queue = []models.MqqtDataDetails{}
				}
				for _, v := range processing {
					if index, ok := list[v.DomotiqueId]; ok {
						setDeviceByIndex(index, v)
					} else {
						index := addDevice(v)
						list[v.DomotiqueId] = index
					}
				}
			} else {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}

func getDevice(id int64) models.MqqtDataDetails {
	for _, v := range Devices {
		if v.DomotiqueId == id {
			return v
		}
	}
	return models.MqqtDataDetails{}
}

func setDevice(device models.MqqtDataDetails) {
	for k, v := range Devices {
		if v.DomotiqueId == device.DomotiqueId {
			Devices[k] = device
			return
		}
	}
	Devices = append(Devices, device)
}
func addDevice(device models.MqqtDataDetails) (index int) {
	Devices = append(Devices, device)
	return len(Devices) - 1
}

func setDeviceByIndex(index int, device models.MqqtDataDetails) {
	if Devices[index].DomotiqueId == device.DomotiqueId {
		Devices[index] = device
	} else {
		fmt.Printf("Device index : %v has id %v whereas update device has id : %v\n", index, Devices[index], device.DomotiqueId)
	}
}
