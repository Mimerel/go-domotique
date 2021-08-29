package models

import (
	"math"
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
	IsHeating             bool
	IsCorrectTemperature  bool
	IpPort                string
	UpdateTime            time.Time
	NormalValues          []HeatingProgram
	Devices               []DeviceToggle
}

type HeatingConfirmation struct {
	IpPort string `yaml:"ipPort,omitempty"`
}

type Status struct {
	k           int
	Value       bool
	Power       float64
	CurrentPos  int64
	Temperature float64
}

func (i *HeatingStatus) GetLastValuesForDevice(config *Configuration) {

	var status Status
	amount := len(i.Devices)
	StatusChan := make(chan Status, amount)
	for k, device := range i.Devices {
		go GetStatus(config, StatusChan, device.DomotiqueId, device, k)
	}

	for j := 0; j < amount; j++ {
		status = <-StatusChan

		if status.Value {
			i.Devices[status.k].StatusOn = "green"
			i.Devices[status.k].StatusOff = ""
		} else {
			i.Devices[status.k].StatusOn = ""
			i.Devices[status.k].StatusOff = "red"
		}
		i.Devices[status.k].Power = math.Round(status.Power)
		i.Devices[status.k].CurrentPos = status.CurrentPos
		i.Devices[status.k].Temperature = status.Temperature
	}

}

func GetStatus(config *Configuration, StatusChan chan Status, domotiqueId int64, device DeviceToggle, k int) {
	var status Status
	if device.Source == 100 {
		status = GetStatusWifi(config, device.DeviceId)
		status.k = k
		StatusChan <- status
		return
	}
	status.k = k
	for _, v := range config.Devices.LastValues {
		if v.DomotiqueId == domotiqueId {
			switch v.Unit {
			case "Level":
				if v.Value == 0 {
					status.Value = false
				} else {
					status.Value = true
				}
			case "Watt":
				status.Power = math.Round(v.Value)
			}
			StatusChan <- status
			return
		}
	}
	StatusChan <- status
}
