package models

type DeviceDetails struct {
	Name         string `csv:"name"`
	DeviceId     int64  `csv:"deviceId"`
	DomotiqueId  int64  `csv:"domotiqueId"`
	RoomId       int64  `csv:"roomId"`
	TypeId       int64  `csv:"typeId"`
	Zwave        int64  `csv:"boxId"`
	Instance     int64  `csv:"instance"`
	CommandClass int64  `csv:"commandClass"`
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
