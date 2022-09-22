package logger

import (
	"fmt"
	"time"
)

type LogLevel int

const (
	NoticeColor    = "\033[1;36m%s\033[0m"
	WarnColor      = "\033[1;33m%s\033[0m"
	ErrorColor     = "\033[1;31m%s\033[0m"
	DebugColor     = "\033[0;36m%s\033[0m"
	DebugPlusColor = "\033[1;34m%s\033[0m"
)

type LogParams struct {
	level LogLevel
}

func NewLogger(level LogLevel) LogParams {
	l := LogParams{level: level}
	return l
}

func (l LogParams) Info(message string, args ...interface{}) {
	computedMessage := fmt.Sprintf(message, args...)
	fmt.Printf(NoticeColor, time.Now().Format(time.RFC3339)+" Info  "+computedMessage+" \n")
}

func (l LogParams) Debug(message string, args ...interface{}) {
	computedMessage := fmt.Sprintf(message, args...)
	fmt.Printf(DebugColor, time.Now().Format(time.RFC3339)+" Debug "+computedMessage+" \n")
}

func (l LogParams) DebugPlus(message string, args ...interface{}) {
	computedMessage := fmt.Sprintf(message, args...)
	fmt.Printf(DebugPlusColor, time.Now().Format(time.RFC3339)+" Debug "+computedMessage+" \n")
}

func (l LogParams) Warn(message string, args ...interface{}) {
	computedMessage := fmt.Sprintf(message, args...)
	fmt.Printf(WarnColor, time.Now().Format(time.RFC3339)+" Warn  "+computedMessage+" \n")
}

func (l LogParams) Error(message string, args ...interface{}) {
	computedMessage := fmt.Sprintf(message, args...)
	fmt.Printf(ErrorColor, time.Now().Format(time.RFC3339)+" Error "+computedMessage+" \n")
}
