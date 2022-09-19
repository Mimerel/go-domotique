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

type ShellyStatus struct {
	Motion      bool
	Timestamp   int
	Active      bool
	Vibration   bool
	Lux         float64
	Bat         float64
	Temperature ShellyInfoThermostatsTemperature `json:"tmp"`
	Target      ShellyInfoThermostatsTemperature `json:"target_t"`
}

type ShellyInfoThermostats struct {
	Position          int                              `json:"pos"`
	TemperatureTarget ShellyInfoThermostatsTemperature `json:"target_t"`
	Temperature       ShellyInfoThermostatsTemperature `json:"tmp"`
	Schedule          bool                             `json:"schedule"`
	ScheduleProfile   int                              `json:"schedule_profile"`
	BoostMinutes      float64                          `json:"boost_minutes"`
	WindowOpen        bool                             `json:"window_open"`
}

type ShellyInfoThermostatsTemperature struct {
	Enabled bool    `json:"enabled"`
	Value   float64 `json:"value"`
	ValueOp float64 `json:"value_op"`
	Units   string  `json:"units"`
	IsValid bool    `json:"is_valid"`
}

type ShellyBattery struct {
	Value   float64 `json:"value"`
	Voltage float64 `json:"voltage"`
}
type ShellyUpdate struct {
	HasUpdate bool `json:"has_update"`
}

type ShellyInfo struct {
	HasUpdate   bool                    `json:"has_update"`
	MACAddress  string                  `json:"mac"`
	Thermostats []ShellyInfoThermostats `json:"thermostats"`
	Calibrated  bool                    `json:"calibrated"`
	Battery     ShellyBattery           `json:"bat"`
	Charger     bool                    `json:"charger"`
	Update      ShellyUpdate            `json:"update"`
	Firmware    ShellyFirmware          `json:"fw_info"`
}

type ShellyFirmware struct {
	Device   string `json:"device"`
	Firmware string `json:"fw"`
}

type ShellySettings struct {
	Device ShellySettingsDevice `json:"device"`
	Name   string               `json:"name"`
}

type ShellySettingsDevice struct {
	DeviceType string `json:"type"`
}

type ShellyAnnounce struct {
	Id              string `json:"id"`
	Model           string `json:"model"`
	Mac             string `json:"mac"`
	IP              string `josn:"ip"`
	NewFirmware     bool   `json:"new_fw"`
	FirewareVersion string `json:"fw_ver"`
}

var mqttConfig *models.Configuration
var Client mqtt.Client
var broker = "tcp://192.168.222.55:1883"
var Devices models.MqqtData
var DataTypes = []string{models.ShellyPower, models.ShellyEnergy, models.ShellyOnOff_0, models.ShellyOnOff_1, models.ShellyOnOff_0_ison, models.ShellyOnOff_1_ison,
	models.ShellyOnline, models.ShellyTemperature0, models.ShellyStatus,
	models.ShellyCurrentPos, models.ShellyRollerLastDirection, models.ShellyRollerStopReason, models.ShellyRollerEnergy, models.ShellyRollerState, models.ShellyRollerPower,
	models.ShellyAnnounce, models.ShellyEventRpc,
	models.ShellyTemperatureDevice, models.ShellyOverTemperatureDevice, models.ShellyReasons, models.ShellySensorBattery, models.ShellyFlood, models.ShellySettings,
	models.ShellyInfo, models.ShellyStatusSwitch0, models.ShellyStatusSwitch1, models.ShellyStatusSwitch2, models.ShellyStatusSwitch3,
	models.ShellyInput1, models.ShellyInput2, models.ShellyInputO, models.ShellyInput3, models.ShellyStatusHumidity0, models.ShellyStatusHumidity1, models.ShellyStatusHumidity2,
	models.ShellyStatusDevicePower0, models.ShellyStatusDevicePower1, models.ShellyStatusDevicePower2,
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var err error
	var messageTopic string

	announce := ShellyAnnounce{}
	messageTopic = msg.Topic()
	if msg.Topic() == "shellies/announce" {
		err = json.Unmarshal(msg.Payload(), &announce)
		if err != nil {
			logger.Debug(mqttConfig, false, "messagePubHandler", "Payload %v", string(msg.Payload()))
			logger.Debug(mqttConfig, false, "messagePubHandler", "Unable to convert response body to json : %+v", err)
		}
		messageTopic = "shellies/" + announce.Id + "/announce"
	}
	id, datatype := getIdFromMessage(messageTopic)
	if datatype == "/command" {
		//	logger.Debug(mqttConfig, false, "messagePubHandler", "Id %v, DataType %v, %v", id, datatype, string(msg.Payload()))
		return
	}
	if id == 0 {
		logger.Debug(mqttConfig, false, "Failed To Find ID", "Topic %v, %v", msg.Topic(), string(msg.Payload()))
	}
	if id == 191 {
		logger.Debug(mqttConfig, false, "messagePubHandler", "Id %v, DataType %v, %v", id, datatype, string(msg.Payload()))
	}
	Devices.Lock()
	instance := Devices.GetInstanceId(datatype)

	if instance >= 0 {
		found := false
		for _, v := range Devices.Id {
			if v.ParentId != 0 {
				//logger.Debug(mqttConfig, false, "messagePubHandler", "Id %v instance %v, parent %v", v.Id, v.Instance, v.ParentId)
				if v.ParentId == id && v.Instance == instance {
					id = v.DomotiqueId
					found = true
				}
			}
		}
		if found == false {
			Devices.Unlock()
			return
		}
		//logger.Debug(mqttConfig, false, "messagePubHandler", "INSTANCE ID=%v, Instance=%v, %v", id, instance, string(msg.Payload()))
	}
	CurrentDevice := Devices.Id[id]
	switch datatype {
	case models.ShellyInfo:
		info := ShellyInfo{}
		err = json.Unmarshal(msg.Payload(), &info)
		if id == 50 {
			logger.Error(mqttConfig, false, "messagePubHandler", "Values Info: %+v", info)
		}
		CurrentDevice.Battery = info.Battery.Value
		CurrentDevice.Voltage = info.Battery.Voltage
		CurrentDevice.NewFirmware = info.HasUpdate
		if len(info.Thermostats) > 0 {
			CurrentDevice.Temperature = info.Thermostats[0].Temperature.Value
			CurrentDevice.TemperatureTarget = info.Thermostats[0].TemperatureTarget.Value
		}
		break
	case models.ShellySettings:
		settings := ShellySettings{}
		err = json.Unmarshal(msg.Payload(), &settings)
		CurrentDevice.NameRegistered = settings.Name
		break
	case models.ShellyEnergy, models.ShellyEnergy2, models.ShellyRollerEnergy:
		CurrentDevice.Energy, err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error(mqttConfig, false, "messagePubHandler", "Unable to convert Payload Energy %v to float", msg.Payload())
		}
		break
	case models.ShellyPower, models.ShellyPower2, models.ShellyRollerPower:
		CurrentDevice.Power, err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error(mqttConfig, false, "messagePubHandler", "Unable to convert Payload Float %v to float", msg.Payload())
		}
		Devices.CalculateTotalWatts()
		break
	case models.ShellyTemperature0:
		CurrentDevice.Temperature, err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error(mqttConfig, false, "messagePubHandler", "Unable to convert Payload Float %v to float", msg.Payload())
		}
		break
	case models.ShellyOnline:
		if len(string(msg.Payload())) > 0 {
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
		break
	case models.ShellyStatus:

		resultStatus := ShellyStatus{}
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
		CurrentDevice.Temperature = resultStatus.Temperature.Value
		CurrentDevice.TemperatureTarget = resultStatus.Target.Value
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
	case models.ShellyStatusSwitch0, models.ShellyStatusSwitch1, models.ShellyStatusSwitch2, models.ShellyStatusSwitch3:
		type temperature struct {
			TC float64 `json:"tC"`
			TF float64 `json:"tF"`
		}
		var result struct {
			Id          int         `json:"id"`
			Source      string      `json:"source"`
			Output      bool        `json:"output"`
			Apower      float64     `json:"apower"`
			Voltage     float64     `json:"voltage"`
			Current     float64     `json:"current"`
			Pf          float64     `json:"pf"`
			Temperature temperature `json:"temperature"`
		}

		err = json.Unmarshal(msg.Payload(), &result)
		if err != nil {
			logger.Debug(mqttConfig, false, "messagePubHandler", "Payload %v", string(msg.Payload()))
			logger.Debug(mqttConfig, false, "messagePubHandler", "Unable to convert response body to json : %+v", err)
			break
		}

		CurrentDevice.Power = result.Apower
		CurrentDevice.Voltage = result.Voltage
		CurrentDevice.DeviceTemperature = result.Temperature.TC
		if result.Output == true {
			CurrentDevice.Status = "on"
			CurrentDevice.StatusOn = "green"
			CurrentDevice.StatusOff = ""
		} else {
			CurrentDevice.Status = "off"
			CurrentDevice.StatusOn = ""
			CurrentDevice.StatusOff = "red"
		}
		break
	case models.ShellyCurrentPos:
		CurrentDevice.CurrentPos, err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error(mqttConfig, false, "messagePubHandler", "Unable to convert Payload CurrentPos Float %v to float", string(msg.Payload()))
		}
		break
	case models.ShellyRollerLastDirection:
		CurrentDevice.LastDirection = string(msg.Payload())
		break
	case models.ShellyRollerStopReason:
		CurrentDevice.Reasons = []string{string(msg.Payload())}
		break
	case models.ShellyReasons:
		CurrentDevice.Reasons = strings.Split(string(msg.Payload()), ",")
		break
	case models.ShellySensorBattery:
		CurrentDevice.Battery, err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error(mqttConfig, false, "messagePubHandler", "Unable to convert Payload Float %v to float", msg.Payload())
		}
		break
	case models.ShellyStatusDevicePower0, models.ShellyStatusDevicePower1, models.ShellyStatusDevicePower2:
		type battery struct {
			Percent float64 `json:"percent"`
			Present bool    `json:"present"`
		}
		var result struct {
			Id       int     `json:"id"`
			Battery  battery `json:"battery"`
			External battery `json:"external"`
		}

		err = json.Unmarshal(msg.Payload(), &result)
		if err != nil {
			logger.Debug(mqttConfig, false, "messagePubHandler", "Payload %v", string(msg.Payload()))
			logger.Debug(mqttConfig, false, "messagePubHandler", "Unable to convert response body to json : %+v", err)
			break
		}
		if result.External.Present {
			CurrentDevice.Battery = 100
		} else {
			CurrentDevice.Battery = result.Battery.Percent
		}
		break
	case models.ShellyTemperatureStatus:
		CurrentDevice.TemperatureStatus = string(msg.Payload())
		break
	case models.ShellyStatusHumidity0, models.ShellyStatusHumidity1, models.ShellyStatusHumidity2:
		type humidity struct {
			RH float64 `json:"rh"`
		}
		result := humidity{}
		err = json.Unmarshal(msg.Payload(), &result)
		if err != nil {
			logger.Debug(mqttConfig, false, "messagePubHandler", "Payload %v", string(msg.Payload()))
			logger.Debug(mqttConfig, false, "messagePubHandler", "Unable to convert response body to json : %+v", err)
			break
		}
		CurrentDevice.Humidity = result.RH
		break
	case models.ShellyStatusTemperature0:
		type temperature struct {
			TC float64 `json:"tC"`
			TF float64 `json:"tF"`
		}
		result := temperature{}
		err = json.Unmarshal(msg.Payload(), &result)
		if err != nil {
			logger.Debug(mqttConfig, false, "messagePubHandler", "Payload %v", string(msg.Payload()))
			logger.Debug(mqttConfig, false, "messagePubHandler", "Unable to convert response body to json : %+v", err)
			break
		}
		CurrentDevice.Temperature = result.TC

	case models.ShellyTemperatureDevice:
		CurrentDevice.DeviceTemperature, err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error(mqttConfig, false, "messagePubHandler", "Unable to convert Payload Float %v to float", msg.Payload())
		}
		break
	case models.ShellyOverTemperatureDevice:
		CurrentDevice.DeviceOverTemperature, err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error(mqttConfig, false, "messagePubHandler", "Unable to convert Payload Float %v to float", msg.Payload())
		}
		break
	case models.ShellyVoltage:
		CurrentDevice.Voltage, err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error(mqttConfig, false, "messagePubHandler", "Unable to convert Payload Float %v to float", msg.Payload())
		}
		break
	case models.ShellyFlood:
		CurrentDevice.StatusBool = ConvertStringToBool(string(msg.Payload()))
		break
	case models.ShellyTemperatureFDevice, models.ShellyInputO, models.ShellyInput1, models.ShellyTemperatures, models.ShellyTemperaturesF, models.ShellyTemperature0F:
		CurrentDevice.Online = true
		break
	case models.ShellyEventRpc:
		type aenergy struct {
			ByMinute []float64 `json:"by_minute"`
			MinuteTS float64   `json:"minute_ts"`
			Total    float64   `json:"total"`
		}
		type switchInstance struct {
			Id      int     `json:"ts"`
			AEnergy aenergy `json:"aenergy"`
		}
		type params struct {
			TS      float64        `json:"ts"`
			Switch0 switchInstance `json:"switch:0"`
			Switch1 switchInstance `json:"switch:1"`
			Switch2 switchInstance `json:"switch:2"`
			Switch3 switchInstance `json:"switch:3"`
		}
		type resultStruct struct {
			Src    string `json:"src"`
			Dst    string `json:"dst"`
			Method string `json:"method"`
			Params params `json:"params"`
		}
		result := resultStruct{}
		err = json.Unmarshal(msg.Payload(), &result)
		if err != nil {
			logger.Debug(mqttConfig, false, "messagePubHandler", "Payload %v", string(msg.Payload()))
			logger.Debug(mqttConfig, false, "messagePubHandler", "Unable to convert response body to json : %+v", err)
			break
		}
		CurrentDevice.NameUnique = result.Src
		break
	default:
		logger.Debug(mqttConfig, false, "No Corresponding Topic", "Id %v, DataType %v, Payload %v", id, datatype, string(msg.Payload()))
		CurrentDevice.Online = true
	}

	Devices.Id[id] = CurrentDevice
	if id == 191 {
		logger.Error(mqttConfig, false, "messagePubHandler", "Values : %+v", CurrentDevice)
	}
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

	if len(topicArray) > 0 {
		id, err = strconv.ParseInt(topicArray[0], 10, 64)
		if err != nil && topic != "shellies/announce" {
			logger.Error(mqttConfig, false, "getIdFromMessage", "Unable to get id from message "+topic)
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
	options.SetClientID("mimerel_mqtt" + time.Now().String())
	options.SetDefaultPublishHandler(messagePubHandler)
	options.OnConnect = connectHandler
	options.OnConnectionLost = connectionLostHandler

	Client = mqtt.NewClient(options)
	token := Client.Connect()
	if token.Wait() && token.Error() != nil {
		logger.Error(mqttConfig, false, "getIdFromMessage", "%v", token.Error())
	}
	Devices.Lock()
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
				ParentId:    temp.ParentId,
				Instance:    temp.Instance,
				Status:      "initial",
				Power:       0,
				DeviceType:  temp.DeviceType,
			}
		}
	}
	Devices.Unlock()

	token = Client.Subscribe("#", 1, nil)
	token.Wait()
	logger.Debug(mqttConfig, false, "getIdFromMessage", "Subscribed to topic %s", "# => All topics")
	go func() {
		for {
			Devices.Lock()
			for _, device := range Devices.Id {
				if device.BoxId == 100 {
					domotiqueId := device.DomotiqueId
					mqttConfig.Logger.Info("device %v", device.Name)
					updateAnnounce(domotiqueId)
					//time.Sleep(time.Second)
				}
			}
			Devices.Unlock()
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
			//Devices.Lock()
			mqttConfig.Channels.MqttReceive <- Devices
			//Devices.Unlock()
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
	//time.Sleep(5 * time.Minute)
	token := Client.Publish(models.Prefix+strconv.FormatInt(domotiqueId, 10)+"/command", 0, false, "announce")
	token.Wait()

}

func ConvertStringToBool(val string) bool {
	if val == "true" {
		return true
	}
	return false
}
