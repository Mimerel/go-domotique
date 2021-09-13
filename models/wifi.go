package models

type WifiStatus struct {
	Relays              []WifiRelays `json:"relays"`
	Meters              []WifiMeters `json:"meters"`
	Rollers             []WifiRoller `json:"rollers"`
	Temperature         float64      `json:"temperature"`
	OverTemperature     bool         `json:"overtemperature"`
	IsValid             bool         `json:"is_valid"`
	ExternalTemperature WifiExtTemp  `json:"ext_temperature"`
}

type WifiExtTemp struct {
	Relay1 WifiRelay `json:"0"`
}

type WifiRelay struct {
	Id         string  `json:"hwID"`
	Celcius    float64 `json:"tC"`
	Fahrenheit float64 `json:"tF"`
}

type WifiRelays struct {
	IsOn           bool   `json:"ison"`
	HasTimer       bool   `json:"has_timer"`
	TimerStarted   int64  `json:"timer_started"`
	TimerDuration  int64  `json:"timer_duration"`
	TimerRemaining int64  `json:"timer_remaining"`
	OverPower      bool   `json:"overpower"`
	Source         string `json:"source"`
}

type WifiRoller struct {
	State           string  `json:"state"`
	Source          string  `json:"source"`
	Power           float64 `json:"power"`
	IsValid         bool    `json:"is_valid"`
	SatefySwitch    bool    `json:"safety_switch"`
	Overtemperature bool    `json:"overtemperature"`
	StopReason      string  `json:"stop_reason"`
	LastDirection   string  `json:"last_direction"`
	CurrentPos      int64   `json:"current_pos"`
	Calibrating     bool    `json:"calibrating"`
	Positioning     bool    `json:"positioning"`
}

type WifiMeters struct {
	Power      float64   `json:"power"`
	OverPrower float64   `json:"overpower"`
	IsValid    bool      `json:"is_valid"`
	Timestamp  int64     `json:"timestamp"`
	Counters   []float64 `json:"counters"`
	Total      int64     `json:"total"`
}
