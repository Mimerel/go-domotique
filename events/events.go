package events

import (
	"go-domotique/devices"
	"go-domotique/models"
	"go-domotique/prowl"
	"go-domotique/utils"
	"strconv"
	"strings"
	"time"
)

func CatchEvent(config *models.Configuration, eventId string, eventValue string, eventZwave string) {
	deviceId, err := strconv.ParseInt(eventId, 10, 64)
	if err != nil {
		config.Logger.Error("unable to convert recevied device Id in int")
	}
	domotique := devices.GetDeviceFromId(config, deviceId)
	config.Logger.Info("Received event from %s %v %v", domotique.Name, domotique.DomotiqueId, domotique.DeviceId)
	go prowl.SendProwlNotification(config, "Event", domotique.Name, eventValue)
	saveEvent(config, domotique, eventValue)
}

func saveEvent(config *models.Configuration, domotique models.DeviceTranslated, eventValue string) {
	db := utils.CreateDbConnection(config)
	db.Database = models.DatabaseStats
	db.Debug = true
	const createdFormat = "2006-01-02 15:04:05"
	timestamp := time.Unix(time.Now().In(config.Location).Unix(), 0).Format(createdFormat)
	logs := []models.Event{models.Event{DomotiqueId: domotique.DomotiqueId, Value: eventValue, DeviceName: domotique.Name, Timestamp: timestamp}}

	col, val, err := db.DecryptStructureAndData(logs)
	if err != nil {
		config.Logger.Error("saveEvent "+"col %s", col)
		config.Logger.Error("saveEvent "+"val %s", val)
	}
	err = db.Insert(false, models.TableEvents, col, val)

	if err != nil {
		config.Logger.Error("saveEvent "+"err %v", err)
		config.Logger.Error("saveEvent "+"table %s", models.TableEvents)
		config.Logger.Error("saveEvent "+"col %s", col)
		values := strings.Split(val, "),(")
		for k, v := range values {
			config.Logger.Error("saveEvent "+"row %v - %s", k, v)
		}
	}
}
