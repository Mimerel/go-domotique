package daemon

import (
	"go-domotique/models"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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
	Position          float64                          `json:"pos"`
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

var DataTypes = []string{models.ShellyPower, models.ShellyEnergy, models.ShellyOnOff_0, models.ShellyOnOff_1, models.ShellyOnOff_0_ison, models.ShellyOnOff_1_ison,
	models.ShellyOnline, models.ShellyTemperature0, models.ShellyStatus,
	models.ShellyCurrentPos, models.ShellyRollerLastDirection, models.ShellyRollerStopReason, models.ShellyRollerEnergy, models.ShellyRollerState, models.ShellyRollerPower,
	models.ShellyAnnounce, models.ShellyEventRpc,
	models.ShellyTemperatureDevice, models.ShellyOverTemperatureDevice, models.ShellyReasons, models.ShellySensorBattery, models.ShellyFlood, models.ShellySettings,
	models.ShellyInfo, models.ShellyStatusSwitch0, models.ShellyStatusSwitch1, models.ShellyStatusSwitch2, models.ShellyStatusSwitch3,
	models.ShellyInput1, models.ShellyInput2, models.ShellyInputO, models.ShellyInput3, models.ShellyStatusHumidity0, models.ShellyStatusHumidity1, models.ShellyStatusHumidity2,
	models.ShellyStatusDevicePower0, models.ShellyStatusDevicePower1, models.ShellyStatusDevicePower2,
	models.ShellySensorState, models.ShellySensorTilt, models.ShellySensorVibration, models.ShellySensorTemperature, models.ShellySensorLux, models.ShellySensorIllumination,
}

func getIdFromMessage(topic string) (id int64, datatype string) {
	var err error
	topic = strings.Replace(topic, models.Prefix, "", -1)
	//logger.Debug(mqttConfig, false, "getIdFromMessage", "topic %v", topic)
	topicArray := strings.Split(topic, "/")

	if len(topicArray) > 0 {
		id, err = strconv.ParseInt(topicArray[0], 10, 64)
		if err != nil && topic != "shellies/announce" {
			mqttConfig.Logger.Error("getIdFromMessage Unable to get id from message " + topic)
		}
		datatype = strings.Replace(topic, topicArray[0], "", -1)
		//mqttConfig.Logger.Error("dataType %v", datatype)
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
}

func reconnectUpdate(initial bool) {

	reconnect(initial)

	for _, temp := range mqttConfig.Devices.DevicesTranslated {
		if temp.BoxId == 100 {
			setDevice(models.MqqtDataDetails{
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
			})
		}
	}
	//Devices.Unlock()
	token := Client.Connect()
	if token.Wait() && token.Error() != nil {
		mqttConfig.Logger.Error("getIdFromMessage %v", token.Error())
	}
	token = Client.Subscribe("#", 1, nil)
	token.Wait()
	mqttConfig.Logger.Error("getIdFromMessage Subscribed to topic %s", "# => All topics")
	go func() {
		for {
			devices := <-mqttConfig.Channels.MqttArray
			for _, device := range devices {
				if device.BoxId == 100 {
					domotiqueId := device.DomotiqueId
					mqttConfig.Logger.Info("device %v", device.Name)
					updateAnnounce(domotiqueId)
					//time.Sleep(time.Second)
				}
			}

			time.Sleep(23 * time.Hour)
		}
	}()

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
