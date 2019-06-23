package models

type Event struct {
	Id          int64  `csv:"id"`
	DomotiqueId int64  `csv:"domotiqueId"`
	DeviceName  string `csv:"deviceName"`
	Timestamp   string `csv:"timestamp"`
	Value       string `csv:"value"`
}
