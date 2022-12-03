package models

import (
	"time"
)

type Heating struct {
	HeatingMoment   HeatingMoment
	HeatingProgram  []HeatingProgram
	LastUpdate      time.Time
	TemporaryValues HeatingMoment
	HeatingSettings []HeatingSettings
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
	Id          int64 `csv:"id"`
	Module      string
	DomotiqueId int64
	RadiatorId  int64
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
	DevicesNew            []MqqtDataDetails
	Rooms                 []Room
	Totals                Totals
}

type Totals struct {
	Watts float64
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
