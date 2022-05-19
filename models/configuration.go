package models

import (
	"github.com/Mimerel/go-utils"
	"sort"
	"sync"
	"time"
)

type Configuration struct {
	MariaDB         MariaDB           `yaml:"mariaDB,omitempty"`
	CharsToReplace  []CharsConversion `yaml:"charsToRemove,omitempty"`
	Token           string            `yaml:"token,omitempty"`
	Heating         Heating           `yaml:"heating,omitempty"`
	Ip              string            `yaml:"ip,omitempty"`
	Port            string            `yaml:"port,omitempty"`
	GoogleAssistant ConfigurationGoogleAssistant
	Daemon          Daemon
	Devices         Devices
	Zwaves          []Zwave
	Rooms           []Room
	DeviceTypes     []DeviceType
	Logger          go_utils.LogParams
	Location        *time.Location
	Channels        Channels
}

type Channels struct {
	UpdateConfig  chan bool
	MqttCall      chan bool
	MqttReceive   chan MqqtData
	MqttSend      chan MqttSendMessage
	MqttReconnect chan bool
}

type CharsConversion struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

type Daemon struct {
	CronTab []CronTab
}

type Devices struct {
	Devices           []DeviceDetails
	DevicesTranslated []DeviceTranslated
	DevicesToggle     []DeviceToggle
	DevicesActions    []DeviceActions
	LastValues        []ElementDetails
}

type MqqtData struct {
	sync.RWMutex
	Id         map[int64]MqqtDataDetails
	TotalWatts float64
}

type MqttSendMessage struct {
	DomotiqueId int64
	Topic       string
	Payload     interface{}
}

type MqqtDataDetails struct {
	DeviceId              int64
	Online                bool
	BoxId                 int64
	DeviceUrl             string
	DomotiqueId           int64
	Room                  string
	Type                  string
	Id                    string
	Mode                  string
	Model                 string
	Mac                   string
	Ip                    string
	NewFirmware           bool
	Status                string
	StatusOn              string
	StatusStop            string
	StatusOff             string
	Name                  string
	Power                 float64
	Energy                float64
	Temperature           float64
	TemperatureStatus     string
	Motion                bool
	Timestamp             int
	Active                bool
	Vibration             bool
	Lux                   float64
	Battery               float64
	CurrentPos            float64
	LastDirection         string
	StopReason            string
	DeviceTemperature     float64
	DeviceOverTemperature float64
	Voltage float64
}

func (i *MqqtDataDetails) GetStatus() float64 {
	if i.Status == "on" {
		return 255
	}
	return 0
}

func (i *MqqtData) ToArray() (result []MqqtDataDetails) {
	for _, v := range i.Id {
		result = append(result, v)
	}
	sort.Slice(result, func(a, b int) bool {
		return result[a].Room+result[a].Name < result[b].Room+result[b].Name
	})
	return result
}

func (i *MqqtData) CalculateTotalWatts() {
	temp := float64(0)
	for _, v := range i.Id {
		temp += v.Power
	}
	i.TotalWatts = temp
}
