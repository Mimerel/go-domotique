package models

type CronTab struct {
	Id              int64 `csv:"id"`
	DomotiqueId     int64 `csv:"domotiqueId"`
	Hour            int64 `csv:"hour"`
	Minute          int64 `csv:"minute"`
	Value           int64 `csv:"value"`
	NotOnAway       int64 `csv:"notOnAway"`
	NotOnAlarmTotal int64 `csv:"notOnAlarmTotal"`
}
