package models

import "time"

type Heating struct {
	SensorId        int64 `yaml:"sensorId"`
	HeaterId        int64 `yaml:"heaterId"`
	HeatingMoment   HeatingMoment
	HeatingProgram  []HeatingProgram
	LastUpdate      time.Time
	TemporaryValues HeatingMoment
	HeatingSettings  HeatingSettings
}

type HeatingMoment struct {
	Moment  time.Time
	Time    int
	Weekday time.Weekday
	Date    int
	Level   float64
}

type HeatingSettings struct {
	Activated bool `csv:"activated"`
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
}

type HeatingConfirmation struct {
	IpPort string `yaml:"ipPort,omitempty"`
}
