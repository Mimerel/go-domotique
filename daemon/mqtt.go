package daemon

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go-domotique/logger"
	"go-domotique/models"
	"strconv"
	"strings"
	"time"
)

const (
	Prefix                    = "shellies/device_"
	ShellyPower               = "/relay/0/power"
	ShellyEnergy              = "/relay/0/energy"
	ShellyOnOff_0             = "/relay/0"
	ShellyOnOff_0_ison        = "/relay/0/ison"
	ShellyOnOff_1             = "/relay/1"
	ShellyOnOff_1_ison        = "/relay/1/ison"
	ShellyOnline              = "/online"
	ShellyTemperature0        = "/ext_temperature/0"
	ShellyStatus              = "/status"
	ShellyRollerState         = "/roller/0"
	ShellyCurrentPos          = "/roller/0/pos" //command/pos pour modifiier
	ShellyRollerLastDirection = "/roller/0/last_direction"
	ShellyRollerStopReason    = "/roller/0/stop_reason"
	ShellyRollerPower         = "/roller/0/power"
	ShellyRollerEnergy        = "/roller/0/energy"
	ShellyAnnounce            = "/announce"
)

var mqttConfig *models.Configuration
var Client mqtt.Client
var Devices models.MqqtData
var DataTypes = []string{ShellyPower, ShellyEnergy, ShellyOnOff_0, ShellyOnOff_1, ShellyOnOff_0_ison, ShellyOnOff_1_ison, ShellyOnline, ShellyTemperature0, ShellyStatus,
	ShellyCurrentPos, ShellyRollerLastDirection, ShellyRollerStopReason, ShellyRollerEnergy, ShellyRollerState, ShellyRollerPower, ShellyAnnounce}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var err error
	id, datatype := getIdFromMessage(msg.Topic())
	//if datatype == ShellyAnnounce {
	//	logger.Debug(mqttConfig, false, "messagePubHandler", "Id %v, DataType %v, %v", id, datatype, string(msg.Payload()))
	//}
	CurrentDevice := Devices.Id[id]
	switch datatype {
	case ShellyEnergy, ShellyRollerEnergy:
		CurrentDevice.Energy, err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error(mqttConfig, false, "messagePubHandler", "Unable to convert Payload Energy %v to float", msg.Payload())
		}
		break
	case ShellyPower, ShellyRollerPower:
		CurrentDevice.Power, err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error(mqttConfig, false, "messagePubHandler", "Unable to convert Payload Float %v to float", msg.Payload())
		}
		break
	case ShellyTemperature0:
		CurrentDevice.Temperature, err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error(mqttConfig, false, "messagePubHandler", "Unable to convert Payload Float %v to float", msg.Payload())
		}
		break
	case ShellyOnline:
		if string(msg.Payload()) == "true" {
			CurrentDevice.Online = true
		} else {
			CurrentDevice.Online = false
		}
	case ShellyOnOff_0, ShellyOnOff_0_ison:
		CurrentDevice.Status = string(msg.Payload())
		if CurrentDevice.Status == "on" {
			CurrentDevice.StatusOn = "green"
			CurrentDevice.StatusOff = ""
		} else {
			CurrentDevice.StatusOn = ""
			CurrentDevice.StatusOff = "red"
		}
		break
	case ShellyOnOff_1, ShellyOnOff_1_ison:
		CurrentDevice.Status = string(msg.Payload())
		if CurrentDevice.Status == "on" {
			CurrentDevice.StatusOn = "green"
			CurrentDevice.StatusOff = ""
		} else {
			CurrentDevice.StatusOn = ""
			CurrentDevice.StatusOff = "red"
		}
		break
	case ShellyRollerState:
		CurrentDevice.Status = string(msg.Payload())
		if CurrentDevice.Status == "open" {
			CurrentDevice.StatusOn = "green"
			CurrentDevice.StatusStop = ""
			CurrentDevice.StatusOff = ""
		}
		if CurrentDevice.Status == "close" {
			CurrentDevice.StatusOn = ""
			CurrentDevice.StatusStop = ""
			CurrentDevice.StatusOff = "green"
		}
		if CurrentDevice.Status == "stop" {
			CurrentDevice.StatusOn = ""
			CurrentDevice.StatusStop = "green"
			CurrentDevice.StatusOff = ""
		}
		break
	case ShellyStatus:
		var resultStatus struct {
			Motion    bool
			Timestamp int
			Active    bool
			Vibration bool
			Lux       float64
			Bat       float64
		}
		err = json.Unmarshal(msg.Payload(), &resultStatus)
		if err != nil {
			logger.Debug(mqttConfig, false, "messagePubHandler", "Payload %v", string(msg.Payload()))
			logger.Debug(mqttConfig, false, "messagePubHandler", "Unable to convert response body to json : %+v", err)
			break
		}
		CurrentDevice.Active = resultStatus.Active
		CurrentDevice.Motion = resultStatus.Motion
		CurrentDevice.Lux = resultStatus.Lux
		CurrentDevice.Battery = resultStatus.Bat
		CurrentDevice.Vibration = resultStatus.Vibration
		break
	case ShellyAnnounce:
		var result struct {
			Id     string
			Mode   string
			Model  string
			Mac    string
			Ip     string
			New_fw bool
			Fw_ver string
		}

		err = json.Unmarshal(msg.Payload(), &result)
		if err != nil {
			logger.Debug(mqttConfig, false, "messagePubHandler", "Payload %v", string(msg.Payload()))
			logger.Debug(mqttConfig, false, "messagePubHandler", "Unable to convert response body to json : %+v", err)
			break
		}
		CurrentDevice.Id = result.Id
		CurrentDevice.Mode = result.Mode
		CurrentDevice.Model = result.Model
		CurrentDevice.Mac = result.Mac
		CurrentDevice.Ip = result.Ip
		CurrentDevice.NewFirmware = result.New_fw
		break
	case ShellyCurrentPos:
		CurrentDevice.CurrentPos, err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error(mqttConfig, false, "messagePubHandler", "Unable to convert Payload CurrentPos Float %v to float", msg.Payload())
		}
		break
	case ShellyRollerLastDirection:
		CurrentDevice.LastDirection = string(msg.Payload())
		break
	case ShellyRollerStopReason:
		CurrentDevice.StopReason = string(msg.Payload())
		break
	}

	Devices.Id[id] = CurrentDevice

	//logger.Debug(mqttConfig, false, "messagePubHandler", "Message %s received for topic %s", msg.Payload(), msg.Topic())
	//for _, v := range Devices.Id {
	//	if v.Power >= 0 {
	//		logger.Debug(mqttConfig, false, "messagePubHandler", "Status %+v", v)
	//	}
	//}

}

func getIdFromMessage(topic string) (id int64, datatype string) {
	var err error
	topic = strings.Replace(topic, Prefix, "", -1)
	//logger.Debug(mqttConfig, false, "getIdFromMessage", "topic %v", topic)
	topicArray := strings.Split(topic, "/")

	if len(topicArray) > 1 {
		id, err = strconv.ParseInt(topicArray[0], 10, 64)
		if err != nil {
			logger.Error(mqttConfig, false, "getIdFromMessage", "Unable to get id from message")
		}
		datatype = strings.Replace(topic, topicArray[0], "", -1)
		//logger.Debug(mqttConfig, false, "getIdFromMessage", "dataType %v", datatype)

	}
	return
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	mqttConfig.Logger.Error("Connected")
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	mqttConfig.Logger.Error("Connection Lost: %s\n", err.Error())
}

func Mqtt_Deamon(c *models.Configuration) {
	mqttConfig = c
	var broker = "tcp://192.168.222.55:1883"
	options := mqtt.NewClientOptions()
	options.AddBroker(broker)
	options.SetClientID("mimerel_mqtt")
	options.SetDefaultPublishHandler(messagePubHandler)
	options.OnConnect = connectHandler
	options.OnConnectionLost = connectionLostHandler

	Client = mqtt.NewClient(options)
	token := Client.Connect()
	if token.Wait() && token.Error() != nil {
		logger.Error(mqttConfig, false, "getIdFromMessage", "%v", token.Error())
	}
	Devices.Id = make(map[int64]models.MqqtDataDetails)

	for _, temp := range c.Devices.DevicesTranslated {
		if temp.BoxId == 100 {
			Devices.Id[temp.DomotiqueId] = models.MqqtDataDetails{
				DeviceId:    temp.DeviceId,
				BoxId:       temp.BoxId,
				DeviceUrl:   temp.ZwaveUrl + strconv.Itoa(int(temp.DeviceId)),
				DomotiqueId: temp.DomotiqueId,
				Room:        temp.Room,
				Name:        temp.Name,
				Type:        temp.Type,
				Status:      "initial",
				Power:       0,
			}
		}
	}

	for _, device := range Devices.Id {
		if device.BoxId == 100 {
			domotiqueId := device.DomotiqueId
			c.Logger.Info("device %v", device.Name)
			for _, event := range DataTypes {
				validevent := event
				go func() {
					topic := Prefix + strconv.FormatInt(domotiqueId, 10) + validevent
					token := Client.Subscribe(topic, 1, nil)
					token.Wait()
					logger.Debug(mqttConfig, false, "getIdFromMessage", "Subscribed to topic %s", topic)
				}()
			}
		}
	}
	go func() {
		for {
			for _, device := range Devices.Id {
				if device.BoxId == 100 {
					domotiqueId := device.DomotiqueId
					c.Logger.Info("device %v", device.Name)
					updateAnnounce(domotiqueId)
					time.Sleep(time.Second)
				}
			}
			time.Sleep(23 * time.Hour)
		}
	}()

	//}
	//
	//client.Disconnect(100)

	defer Client.Disconnect(100)

	for {
		select {
		case <-mqttConfig.Channels.MqttCall:
			mqttConfig.Channels.MqttReceive <- Devices
			break
		case mqttAction := <-mqttConfig.Channels.MqttSend:
			actionToDo := Prefix + strconv.FormatInt(mqttAction.DomotiqueId, 10) + mqttAction.Topic
			logger.Debug(mqttConfig, false, "ActionToDo", "Action : %s - payload %v", actionToDo, mqttAction.Payload)
			token := Client.Publish(actionToDo, 0, false, mqttAction.Payload)
			switch mqttAction.Topic + fmt.Sprint("%v", mqttAction.Payload) {
			case "/commandupdate_fw":
				go updateAnnounce(mqttAction.DomotiqueId)
				break
			}
			token.Wait()
			break
		}
	}
}

func updateAnnounce(domotiqueId int64) {
	time.Sleep(5 * time.Minute)
	token := Client.Publish(Prefix+strconv.FormatInt(domotiqueId, 10)+"/command", 0, false, "announce")
	token.Wait()

}
