package models

import (
	"time"
)

type Heating struct {
	HeatingMoment   HeatingMoment
	HeatingProgram  []HeatingProgram
	LastUpdate      time.Time
	TemporaryValues HeatingMoment
	HeatingSettings HeatingSettings
	HeatingLevels   []HeatingLevels
}

type HeatingMoment struct {
	Moment  time.Time
	Time    int
	Weekday time.Weekday
	Date    int
	Level   float64
}

type HeatingLevels struct {
	Id    int64   `csv:"id"`
	Name  string  `csv:"name"`
	Value float64 `csv:"value"`
}

type HeatingSettings struct {
	Id        int64 `csv:"id"`
	Activated bool  `csv:"activated"`
	SensorId  int64 `csv:"sensorId"`
	HeaterId  int64 `csv:"heaterId"`
}

type HeatingProgram struct {
	DayId      int64   `csv:"dayId"`
	Day        string  `csv:"day"`
	Moment     int64   `csv:"moment"`
	LevelName  string  `csv:"levelName"`
	LevelValue float64 `csv:"levelValue"`
}

type HeatingStatus struct {
	Heater_Level          float64
	Temperature_Requested float64
	Temperature_Actual    float64
	Until                 time.Time
	TemporaryLevel        float64
	IsTemporary           bool
	IpPort                string
	UpdateTime            time.Time
	NormalValues          []HeatingProgram
	Devices               []DeviceToggle
}

type HeatingConfirmation struct {
	IpPort string `yaml:"ipPort,omitempty"`
}

func (i *HeatingStatus) GetLastValuesForDevice(config *Configuration) {
	for k, device := range i.Devices {
		if device.Source == 100 {
			value, power := GetStatusWifi(config, device.DeviceId)
			if value {
				i.Devices[k].StatusOn = "green"
				i.Devices[k].StatusOff = ""
			} else {
				i.Devices[k].StatusOn = ""
				i.Devices[k].StatusOff = "red"
			}
			i.Devices[k].Power = power
			continue
		}
		for _, v := range config.Devices.LastValues {
			if v.DomotiqueId == i.Devices[k].DomotiqueId {
				switch v.Unit {
				case "Level":
					if v.Value == 0 {
						i.Devices[k].StatusOn = ""
						i.Devices[k].StatusOff = "red"
					} else {
						i.Devices[k].StatusOn = "green"
						i.Devices[k].StatusOff = ""
					}
				case "Watt":
					i.Devices[k].Power = v.Value
				}
			}
		}
	}
}
