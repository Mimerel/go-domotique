package models

import (
	"github.com/Mimerel/go-utils"
	"sort"
	"strconv"
	"strings"
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
	UpdateConfig        chan bool
	MqttSend            chan MqttSendMessage
	MqttReconnect       chan bool
	MqttDomotiqueIdGet  chan int64
	MqttDomotiqueDevice chan MqqtDataDetails
	MqttGetArray        chan bool
	MqttArray           chan []MqqtDataDetails
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
	ParentId              int64
	Instance              int64
	DeviceType            string
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
	NameRegistered        string
	NameUnique            string
	Power                 float64
	Energy                float64
	Temperature           float64
	TemperatureTarget     float64
	TemperatureStatus     string
	Humidity              float64
	Motion                bool
	Timestamp             int
	Active                bool
	Vibration             bool
	Lux                   float64
	Battery               float64
	CurrentPos            float64
	LastDirection         string
	StopReason            string // a integrer dans Reasons
	DeviceTemperature     float64
	DeviceOverTemperature float64
	Voltage               float64
	Reasons               []string
	Valid                 bool
	StatusBool            bool
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

func (i *MqqtData) GetInstanceId(value string) int64 {
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
