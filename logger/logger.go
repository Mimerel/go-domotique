package logger

import (
	"fmt"
	"go-domotique/models"
	"go-domotique/utils"
	"strings"
	"time"
)

func Info(config *models.Configuration, module string, message string, args ...interface{}) {
	computedMessage := fmt.Sprintf(message, args...)
	fmt.Printf(time.Now().In(config.Location).Format(time.RFC3339)+" - Info (%s): %s \n", module, computedMessage)
	sendLogToDB(config, "Info", module, computedMessage)
}

func Debug(config *models.Configuration, module string, message string, args ...interface{}) {
	computedMessage := fmt.Sprintf(message, args...)
	fmt.Printf(time.Now().In(config.Location).Format(time.RFC3339)+" - Debug (%s): %s \n", module, computedMessage)
	sendLogToDB(config, "Debug" , module, computedMessage)
}

func Error(config *models.Configuration, module string, message string, args ...interface{}) {
	computedMessage := fmt.Sprintf(message, args...)
	fmt.Printf(time.Now().In(config.Location).Format(time.RFC3339)+" - Error (%s): %s \n", module, computedMessage)
	sendLogToDB(config, "Error" , module, computedMessage)
}

func sendLogToDB(c *models.Configuration, messageType string, module string, computedMessage string) {
	db := utils.CreateDbConnection(c)
	db.Database = utils.DatabaseLogger
	db.Debug = true
	logs := []models.Log{models.Log{Module: module, Message: computedMessage, Type: messageType}}

	col, val, err := db.DecryptStructureAndData(logs)
	if err != nil {
		c.Logger.Error("col %s", col)
		c.Logger.Error("val %s", val)
	}
	err = db.Insert(false, utils.LoggerDomotique, col, val)

	if err != nil {
		c.Logger.Error("table %s", utils.LoggerDomotique)
		c.Logger.Error("col %s", col)
		values := strings.Split(val, "),(")
		for k, v := range values {
			c.Logger.Error("row %v - %s", k, v)
		}
	}
}
