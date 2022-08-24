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
	BoxId        int64  `csv:"boxId"`
	Instance     int64  `csv:"instance"`
	CommandClass int64  `csv:"commandClass"`
	OnUi         int64  `csv:"onUi"`
	TypeWifi     string `csv:"typeWifi"`
	DeviceType   string `csv:"model"`
}

type DeviceTranslated struct {
	Name           string `csv:"name"`
	DeviceId       int64  `csv:"deviceId"`
	DeviceIdString string `csv:"deviceIdString"`
	DomotiqueId    int64  `csv:"domotiqueId"`
	Room           string `csv:"room"`
	Type           string `csv:"type"`
	BoxId          int64  `csv:"boxId"`
	ZwaveName      string `csv:"zwaveName"`
	ZwaveUrl       string `csv:"zwaveUrl"`
	Instance       int64  `csv:"instance"`
	InstanceString string `csv:"instanceString"`
	CommandClass   int64  `csv:"commandClass"`
	TypeWifi       string `csv:"typeWifi"`
	DeviceType     string `csv:"model"`
}

type DeviceToggle struct {
	DomotiqueId int64
	DeviceId    int64
	DeviceType  string
	Source      int64
	Type        string
	Name        string
	Room        string
	UrlOn       string
	UrlStop     string
	UrlOff      string
	UrlSlide    string
	StatusOn    string
	StatusOff   string
	Power       float64
	CurrentPos  int64
	Temperature float64
}

type DeviceActions struct {
	Id           int64 `csv:"id"`
	DomotiqueId  int64 `csv:"domotiqueId"`
	ActionNameId int64 `csv:"actionNameId"`
}

func (i *DeviceTranslated) CollectDeviceToggleDetails(config *Configuration) (deviceToggle DeviceToggle) {
	deviceToggle.Name = i.Name
	deviceToggle.DeviceId = i.DeviceId
	deviceToggle.Type = i.Type
	deviceToggle.Room = i.Room
	deviceToggle.Source = i.BoxId
	deviceToggle.DomotiqueId = i.DomotiqueId
	deviceToggle.DeviceType = i.DeviceType
	switch i.BoxId {
	case 100:
		break
	default:
		deviceToggle.UrlOn = GetRequest(config, i.ZwaveUrl, i.DeviceId, i.Instance, i.CommandClass, 255)
		deviceToggle.UrlOff = GetRequest(config, i.ZwaveUrl, i.DeviceId, i.Instance, i.CommandClass, 0)
	}
	return deviceToggle
}

func GetRequest(config *Configuration, url string, id int64, instance int64, commandClass int64, level int64) string {
	config.Logger.Info("GetRequest", "Préparing url")
	postingUrl := "http://" + url + ":8083/ZWaveAPI/Run/devices[" + strconv.FormatInt(id, 10) + "].instances[" + strconv.FormatInt(instance, 10) + "].commandClasses[" + strconv.FormatInt(commandClass, 10) + "].Set(" + strconv.FormatInt(level, 10) + ")"
	return postingUrl
}

func (i *DeviceTranslated) GetUrlForValue(config *Configuration, value int64) (postingUrl string, command string) {
	switch i.TypeWifi {
	case "relay":
		switch value {
		case 0:
			return "/relay/0/command", "off"
		case 255:
			return "/relay/0/command", "on"
		}
	case "roller":
		switch value {
		case -1:
			return "/roller/0/command", "stop"
		case 0:
			return "/roller/0/command", "close"
		case 255:
			return "/roller/0/command", "open"
		default:
			return "/roller/0/command/pos", strconv.Itoa(int(value))
		}
	}
	return "", ""
}
