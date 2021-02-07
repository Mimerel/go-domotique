package models

import (
	"github.com/Mimerel/go-utils"
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
