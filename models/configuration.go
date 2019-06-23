package models

import (
	"github.com/Mimerel/go-utils"
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
	LastValues        []ElementDetails
}

type ConfigurationGoogleAssistant struct {
	GoogleWords                  []GoogleWords
	GoogleBoxes                  []GoogleBox
	GoogleInstructions           []GoogleInstruction
	GoogleActionNames            []GoogleActionNames
	GoogleTranslatedInstructions []GoogleTranslatedInstruction
	GoogleActionTypesWords       []GoogleActionTypesWords
	GoogleActionTypes            []GoogleActionTypes
	GoogleTranslatedActionTypes  []GoogleTranslatedActionTypes
}
