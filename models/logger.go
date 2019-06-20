package models

type Log struct {
	Type string `csv:"type"`
	Module   string `csv:"module"`
	Message  string `csv:"message"`
}
