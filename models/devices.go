package models

import (
	"strconv"
)

type DeviceDetails struct {
	Name         string `csv:"name"`
	DeviceId     int64  `csv:"deviceId"`
	DomotiqueId  int64  `csv:"domotiqueId"`
	RoomId       int64  `csv:"roomId"`
	TypeId       int64  `csv:"typeId"`
	Zwave        int64  `csv:"boxId"`
	Instance     int64  `csv:"instance"`
	CommandClass int64  `csv:"commandClass"`
	OnUi         int64  `csv:"onUi"`
}

type DeviceTranslated struct {
	Name         string `csv:"name"`
	DeviceId     int64  `csv:"deviceId"`
	DomotiqueId  int64  `csv:"domotiqueId"`
	Room         string `csv:"room"`
	Type         string `csv:"type"`
	Zwave        int64  `csv:"boxId"`
	ZwaveName    string `csv:"zwaveName"`
	ZwaveUrl     string `csv:"zwaveUrl"`
	Instance     int64  `csv:"instance"`
	CommandClass int64  `csv:"commandClass"`
}

type DeviceToggle struct {
	DomotiqueId int64
	Type        string
	Name        string
	UrlOn       string
	UrlOff      string
}

func (i *DeviceTranslated) CollectDeviceToggleDetails(config *Configuration) (deviceToggle DeviceToggle) {
	deviceToggle.Name = i.Name
	deviceToggle.Type = i.Type
	deviceToggle.DomotiqueId = i.DomotiqueId
	deviceToggle.UrlOn = GetRequest(config, i.ZwaveUrl, i.DeviceId, i.Instance, i.CommandClass, 255)
	deviceToggle.UrlOff = GetRequest(config, i.ZwaveUrl, i.DeviceId, i.Instance, i.CommandClass, 0)
	return deviceToggle
}

func GetRequest(config *Configuration, url string, id int64, instance int64, commandClass int64, level int64) string {
	config.Logger.Info("GetRequest", "Préparing url")
	postingUrl := "http://" + url + ":8083/ZWaveAPI/Run/devices[" + strconv.FormatInt(id, 10) + "].instances[" + strconv.FormatInt(instance, 10) + "].commandClasses[" + strconv.FormatInt(commandClass, 10) + "].Set(" + strconv.FormatInt(level, 10) + ")"
	return postingUrl
}
