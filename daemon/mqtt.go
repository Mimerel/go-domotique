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


var mqttConfig *models.Configuration
var Client mqtt.Client
var broker = "tcp://192.168.222.55:1883"
var Devices models.MqqtData
var DataTypes = []string{models.ShellyPower, models.ShellyEnergy, models.ShellyOnOff_0, models.ShellyOnOff_1, models.ShellyOnOff_0_ison, models.ShellyOnOff_1_ison, models.ShellyOnline, models.ShellyTemperature0, models.ShellyStatus,
	models.ShellyCurrentPos, models.ShellyRollerLastDirection, models.ShellyRollerStopReason, models.ShellyRollerEnergy, models.ShellyRollerState, models.ShellyRollerPower, models.ShellyAnnounce,
	models.ShellyTemperatureDevice, models.ShellyOverTemperatureDevice}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var err error
	id, datatype := getIdFromMessage(msg.Topic())
	//if datatype == ShellyAnnounce {
	//	logger.Debug(mqttConfig, false, "messagePubHandler", "Id %v, DataType %v, %v", id, datatype, string(msg.Payload()))
	//}
	//if id == 33 || id== 34 || id== 32 || id== 30 {
	//	logger.Debug(mqttConfig, false, "messagePubHandler", "Id %v, DataType %v, %v", id, datatype, string(msg.Payload()))
	//}
	Devices.Lock()
	CurrentDevice := Devices.Id[id]
	switch datatype {
	case models.ShellyEnergy,models.ShellyEnergy2, models.ShellyRollerEnergy:
		CurrentDevice.Energy, err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error(mqttConfig, false, "messagePubHandler", "Unable to convert Payload Energy %v to float", msg.Payload())
		}
		CurrentDevice.Online = true
		break
	case models.ShellyPower,models.ShellyPower2, models.ShellyRollerPower:
		CurrentDevice.Power, err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error(mqttConfig, false, "messagePubHandler", "Unable to convert Payload Float %v to float", msg.Payload())
		}
		Devices.CalculateTotalWatts()
		CurrentDevice.Online = true
		break
	case models.ShellyTemperature0:
		CurrentDevice.Temperature, err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error(mqttConfig, false, "messagePubHandler", "Unable to convert Payload Float %v to float", msg.Payload())
		}
		CurrentDevice.Online = true
		break
	case models.ShellyOnline:
		if len(string(msg.Payload()))>0 {
			CurrentDevice.Online = true
		} else {
			CurrentDevice.Online = false
		}
		break
	case models.ShellyOnOff_0, models.ShellyOnOff_0_ison:
		CurrentDevice.Status = string(msg.Payload())
		if CurrentDevice.Status == "on" {
			CurrentDevice.StatusOn = "green"
			CurrentDevice.StatusOff = ""
		} else {
			CurrentDevice.StatusOn = ""
			CurrentDevice.StatusOff = "red"
		}
		CurrentDevice.Online = true
		break
	case models.ShellyOnOff_1, models.ShellyOnOff_1_ison:
		CurrentDevice.Status = string(msg.Payload())
		if CurrentDevice.Status == "on" {
			CurrentDevice.StatusOn = "green"
			CurrentDevice.StatusOff = ""
		} else {
			CurrentDevice.StatusOn = ""
			CurrentDevice.StatusOff = "red"
		}
		CurrentDevice.Online = true
		break
	case models.ShellyRollerState:
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
		CurrentDevice.Online = true
		break
	case models.ShellyStatus:
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
	case models.ShellyAnnounce:
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
	case models.ShellyCurrentPos:
		CurrentDevice.CurrentPos, err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error(mqttConfig, false, "messagePubHandler", "Unable to convert Payload CurrentPos Float %v to float", msg.Payload())
		}
		break
	case models.ShellyRollerLastDirection:
		CurrentDevice.LastDirection = string(msg.Payload())
		break
	case models.ShellyRollerStopReason:
		CurrentDevice.StopReason = string(msg.Payload())
		break
	case models.ShellyTemperatureStatus:
		CurrentDevice.TemperatureStatus = string(msg.Payload())
		break
	case models.ShellyTemperatureDevice:
		CurrentDevice.DeviceTemperature , err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error(mqttConfig, false, "messagePubHandler", "Unable to convert Payload Float %v to float", msg.Payload())
		}
		break
	case models.ShellyOverTemperatureDevice:
		CurrentDevice.DeviceOverTemperature , err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error(mqttConfig, false, "messagePubHandler", "Unable to convert Payload Float %v to float", msg.Payload())
		}
		break
	case models.ShellyVoltage:
		CurrentDevice.Voltage , err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error(mqttConfig, false, "messagePubHandler", "Unable to convert Payload Float %v to float", msg.Payload())
		}
		break
	case models.ShellyTemperatureFDevice, models.ShellyInputO, models.ShellyInput1, models.ShellyTemperatures, models.ShellyTemperaturesF,models.ShellyTemperature0F:
		break
	default :
		logger.Debug(mqttConfig, false, "messagePubHandler", "Id %v, DataType %v, %v", id, datatype, string(msg.Payload()))
	}

	Devices.Id[id] = CurrentDevice
	Devices.Unlock()
	//logger.Debug(mqttConfig, false, "messagePubHandler", "Message %s received for topic %s", msg.Payload(), msg.Topic())
	//for _, v := range Devices.Id {
	//	if v.Power >= 0 {
	//		logger.Debug(mqttConfig, false, "messagePubHandler", "Status %+v", v)
	//	}
	//}

}

func getIdFromMessage(topic string) (id int64, datatype string) {
	var err error
	topic = strings.Replace(topic, models.Prefix, "", -1)
	//logger.Debug(mqttConfig, false, "getIdFromMessage", "topic %v", topic)
	topicArray := strings.Split(topic, "/")

	if len(topicArray) > 1 {
		id, err = strconv.ParseInt(topicArray[0], 10, 64)
		if err != nil {
			logger.Error(mqttConfig, false, "getIdFromMessage", "Unable to get id from message " + topic)
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

func reconnect(initial bool) {
	if !initial {
		Client.Disconnect(1)
	}
	options := mqtt.NewClientOptions()
	options.AddBroker(broker)
	options.SetClientID("mimerel_mqtt"+time.Now().String())
	options.SetDefaultPublishHandler(messagePubHandler)
	options.OnConnect = connectHandler
	options.OnConnectionLost = connectionLostHandler

	Client = mqtt.NewClient(options)
	token := Client.Connect()
	if token.Wait() && token.Error() != nil {
		logger.Error(mqttConfig, false, "getIdFromMessage", "%v", token.Error())
	}
	Devices.Id = make(map[int64]models.MqqtDataDetails)

	for _, temp := range mqttConfig.Devices.DevicesTranslated {
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

/*	for _, device := range Devices.Id {
		if device.BoxId == 100 {
			domotiqueId := device.DomotiqueId
			mqttConfig.Logger.Info("device %v", device.Name)
			for _, event := range DataTypes {
				validevent := event
				go func() {
					topic := models.Prefix + strconv.FormatInt(domotiqueId, 10) + validevent
					token := Client.Subscribe(topic, 1, nil)
					token.Wait()
					logger.Debug(mqttConfig, false, "getIdFromMessage", "Subscribed to topic %s", topic)
				}()
			}
		}
	}*/
	token = Client.Subscribe("#", 1, nil)
	token.Wait()
	logger.Debug(mqttConfig, false, "getIdFromMessage", "Subscribed to topic %s", "# => All topics")
	go func() {
		for {
			for _, device := range Devices.Id {
				if device.BoxId == 100 {
					domotiqueId := device.DomotiqueId
					mqttConfig.Logger.Info("device %v", device.Name)
					updateAnnounce(domotiqueId)
					time.Sleep(time.Second)
				}
			}
			time.Sleep(23 * time.Hour)
		}
	}()

}

func Mqtt_Deamon(c *models.Configuration) {
	mqttConfig = c

	reconnect(true)

	defer Client.Disconnect(100)

	for {
		select {
		case <-mqttConfig.Channels.MqttCall:
			mqttConfig.Channels.MqttReceive <- Devices
			break
		case <-mqttConfig.Channels.MqttReconnect:
			reconnect(false)
			break
		case mqttAction := <-mqttConfig.Channels.MqttSend:
			actionToDo := models.Prefix + strconv.FormatInt(mqttAction.DomotiqueId, 10) + mqttAction.Topic
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
	token := Client.Publish(models.Prefix+strconv.FormatInt(domotiqueId, 10)+"/command", 0, false, "announce")
	token.Wait()

}


