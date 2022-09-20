package daemon

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go-domotique/models"
	"strconv"
	"strings"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var err error
	var messageTopic string
	logger := mqttConfig.Logger
	announce := ShellyAnnounce{}
	messageTopic = msg.Topic()
	if msg.Topic() == "shellies/announce" {
		err = json.Unmarshal(msg.Payload(), &announce)
		if err != nil {
			logger.Debug("Payload %v", string(msg.Payload()))
			logger.Debug("Unable to convert response body to json : %+v", err)
		}
		messageTopic = "shellies/" + announce.Id + "/announce"
	}
	id, datatype := getIdFromMessage(messageTopic)
	if datatype == "/command" {
		//	logger.Debug(mqttConfig, false, "messagePubHandler", "Id %v, DataType %v, %v", id, datatype, string(msg.Payload()))
		return
	}
	if id == 0 {
		logger.Debug("Topic %v, %v", msg.Topic(), string(msg.Payload()))
		return
	}

	go updateDeviceValuesFromMessage(id, datatype, msg)
}

func updateDeviceValuesFromMessage(id int64, datatype string, msg mqtt.Message) {
	var err error
	instance := getInstanceId(datatype)

	logger := mqttConfig.Logger
	mqttConfig.Channels.MqttGetArray <- true
	deviceList := <-mqttConfig.Channels.MqttArray

	if instance >= 0 {
		found := false
		for _, v := range deviceList {
			if v.ParentId != 0 {
				logger.Debug("Id %v instance %v, parent %v", v.Id, v.Instance, v.ParentId)
				if v.ParentId == id && v.Instance == instance {
					id = v.DomotiqueId
					found = true
				}
			}
		}
		if found == false {
			return
		}
		//logger.Debug(mqttConfig, false, "messagePubHandler", "INSTANCE ID=%v, Instance=%v, %v", id, instance, string(msg.Payload()))
	}

	mqttConfig.Channels.MqttDomotiqueIdGet <- id
	CurrentDevice := <-mqttConfig.Channels.MqttDomotiqueDeviceGet
	switch datatype {
	case models.ShellyInfo:
		info := ShellyInfo{}
		err = json.Unmarshal(msg.Payload(), &info)
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
			logger.Error("Unable to convert Payload Energy %v to float", msg.Payload())
		}
		break
	case models.ShellyPower, models.ShellyPower2, models.ShellyRollerPower:
		CurrentDevice.Power, err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error("Unable to convert Payload Float %v to float", msg.Payload())
		}
		break
	case models.ShellyTemperature0:
		CurrentDevice.Temperature, err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error("Unable to convert Payload Float %v to float", msg.Payload())
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
			logger.Debug("Payload %v", string(msg.Payload()))
			logger.Debug("Unable to convert response body to json : %+v", err)
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
			logger.Debug("Payload %v", string(msg.Payload()))
			logger.Debug("Unable to convert response body to json : %+v", err)
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
			logger.Debug("Payload %v", string(msg.Payload()))
			logger.Debug("Unable to convert response body to json : %+v", err)
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
			logger.Error("Unable to convert Payload CurrentPos Float %v to float", string(msg.Payload()))
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
			logger.Error("Unable to convert Payload Float %v to float", msg.Payload())
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
			logger.Debug("Payload %v", string(msg.Payload()))
			logger.Debug("Unable to convert response body to json : %+v", err)
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
			logger.Debug("Payload %v", string(msg.Payload()))
			logger.Debug("Unable to convert response body to json : %+v", err)
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
			logger.Debug("Payload %v", string(msg.Payload()))
			logger.Debug("Unable to convert response body to json : %+v", err)
			break
		}
		CurrentDevice.Temperature = result.TC

	case models.ShellyTemperatureDevice:
		CurrentDevice.DeviceTemperature, err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error("Unable to convert Payload Float %v to float", msg.Payload())
		}
		break
	case models.ShellyOverTemperatureDevice:
		CurrentDevice.DeviceOverTemperature, err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error("Unable to convert Payload Float %v to float", msg.Payload())
		}
		break
	case models.ShellyVoltage:
		CurrentDevice.Voltage, err = strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			logger.Error("Unable to convert Payload Float %v to float", msg.Payload())
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
			logger.Debug("Payload %v", string(msg.Payload()))
			logger.Debug("Unable to convert response body to json : %+v", err)
			break
		}
		CurrentDevice.NameUnique = result.Src
		break
	default:
		logger.Debug("Id %v, DataType %v, Payload %v", id, datatype, string(msg.Payload()))
		CurrentDevice.Online = true
	}

	mqttConfig.Channels.MqttDomotiqueDevicePost <- CurrentDevice

}

func getInstanceId(value string) int64 {
	instance := ""
	instances := []string{"/0", "/1", "/2", "/3", "/4", ":0", ":1", ":2", ":3"}
	for _, instanceFound := range instances {
		if strings.Contains(value, instanceFound) {
			instanceFound = strings.Replace(instanceFound, "/", "", -1)
			instance = strings.Replace(instanceFound, ":", "", -1)
		}
	}
	if instance == "" {
		return -1
	}
	found, err := strconv.ParseInt(instance, 10, 64)
	if err != nil {
		return -1
	}
	return found
}
