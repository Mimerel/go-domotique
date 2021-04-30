package models

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type WifiStatus struct {
	Relays          []WifiRelays `json:"relays"`
	Meters          []WifiMeters `json:"meters"`
	Rollers         []WifiRoller `json:"rollers"`
	Temperature     float64      `json:"temperature"`
	OverTemperature bool         `json:"overtemperature"`
	IsValid         bool         `json:"is_valid"`
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

func GetStatusWifi(config *Configuration, id int64) (status bool, power float64, currentPos int64) {
	config.Logger.Info("GetStatusWifi", "Préparing url")
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	postingUrl := "http://" + config.Ip[:12] + strconv.Itoa(int(id)) + "/status"
	config.Logger.Info("GetStatusWifi", "Request posted : %s", postingUrl)

	res, err := client.Get(postingUrl)
	if err != nil {
		config.Logger.Info("GetStatusWifi", "Failed to execute request %s ", postingUrl, err)
		return false, 0, -1
	}
	temp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		config.Logger.Info("GetStatusWifi", "Failed to read body:", err)
		return false, 0, -1
	}

	defer res.Body.Close()
	data := WifiStatus{}
	err = json.Unmarshal(temp, &data)
	if err != nil {
		config.Logger.Warn("GetStatusWifi", "received response %v", string(temp))
		config.Logger.Error("GetStatusWifi", "error decoding wifi status response: %v", err)
	}
	config.Logger.DebugPlus("GetStatusWifi", "Status: %+v", data)

	if len(data.Relays) > 0 && data.Relays[0].IsOn {
		status = true
	}
	if len(data.Meters) > 0 {
		power = data.Meters[0].Power
	}
	if len(data.Rollers) > 0 {
		currentPos = data.Rollers[0].CurrentPos
	}
	return status, power, currentPos
}
