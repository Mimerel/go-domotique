package events

import (
	"go-domotique/devices"
	"go-domotique/logger"
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
		logger.Error(config, true,"CatchEvent", "unable to convert recevied device Id in int")
	}
	zwaveId := devices.GetZwaveIdFromZwaveName(config, eventZwave).Id
	if zwaveId == 0 {
		logger.Error(config, true,"Unable to find corresponding domotiqueId for Id <%v>, eventValue <%v>, eventValue <%v>", eventId, eventZwave, eventValue)
	}
	domotique := devices.GetDeviceFromId(config, deviceId)
	logger.Info(config, false, "CatchEvent", "Received event from %s %v %v", domotique.Name, domotique.DomotiqueId, domotique.DeviceId)
	prowl.SendProwlNotification(config, "Event", domotique.Name, eventValue)
	saveEvent(config, domotique, eventValue)
}

func saveEvent(config *models.Configuration, domotique models.DeviceTranslated, eventValue string) {
	db := utils.CreateDbConnection(config)
	db.Database = utils.DatabaseStats
	db.Debug = true
	const createdFormat = "2006-01-02 15:04:05"
	timestamp := time.Unix(time.Now().In(config.Location).Unix(), 0).Format(createdFormat)
	logs := []models.Event{models.Event{DomotiqueId:domotique.DomotiqueId, Value:eventValue, DeviceName:domotique.Name, Timestamp:timestamp}}

	col, val, err := db.DecryptStructureAndData(logs)
	if err != nil {
		logger.Error(config, true,"saveEvent", "col %s", col)
		logger.Error(config, true,"saveEvent", "val %s", val)
	}
	err = db.Insert(false, utils.TableEvents, col, val)

	if err != nil {
		logger.Error(config, true,"saveEvent", "err %v", err)
		logger.Error(config, true,"saveEvent", "table %s", utils.TableEvents)
		logger.Error(config, true,"saveEvent", "col %s", col)
		values := strings.Split(val, "),(")
		for k, v := range values {
			logger.Error(config, true,"saveEvent", "row %v - %s", k, v)
		}
	}
}