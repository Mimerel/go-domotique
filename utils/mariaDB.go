package utils

import (
	"github.com/Mimerel/go-utils"
	"go-domotique/models"
	"strings"
)

const (
	ActionReplaceLowPriority = "replaceLowPriority"
	ActionReplace            = "replace"
	ActionInsertIgnore       = "insertIgnore"
	ActionInsert             = "insert"

	DatabaseLogger  = "logs"
	DatabaseStats  = "domotiqueStats"
	LoggerDomotique = "domotique"

	TableDomotiqueBox                 = "domotiqueBox"
	TableDevices                      = "devices"
	TableDevicesTranslated            = "devicesTranslated"
	TableDevicesLastValues            = "devicesLastValues"
	TableRooms                        = "rooms"
	TableDeviceTypes                  = "deviceTypes"
	TableGoogleWords                  = "googleWords"
	TableGoogleInstructions           = "googleInstructions"
	TableGoogleActionNames            = "googleActionNames"
	TableGoogleBox                    = "googleBox"
	TableGoogleActionTypes            = "googleActionTypes"
	TableGoogleActionTypesWords       = "googleActionTypesWords"
	TableGoogleTranslatedInstructions = "googleTranslatedInstructions"
	TableCronTab                      = "cronTab"
)

func CreateDbConnection(c *models.Configuration) (db *go_utils.MariaDBConfiguration) {
	db = go_utils.NewMariaDB()
	db.Database = c.MariaDB.Database
	db.User = c.MariaDB.User
	db.Password = c.MariaDB.Password
	db.IP = c.MariaDB.IP
	db.Port = c.MariaDB.Port
	db.LoggerError = c.Logger.Error
	db.LoggerInfo = c.Logger.Info
	return db
}

/**
Method that stores the content of a []structure to mariaDB
*/
func ActionInMariaDB(c *models.Configuration, data interface{}, table string, action string) (err error) {
	db := CreateDbConnection(c)
	col, val, err := db.DecryptStructureAndData(data)
	if err != nil {
		c.Logger.Error("col %s", col)
		c.Logger.Error("val %s", val)
		return err
	}
	switch action {
	case ActionReplaceLowPriority:
		err = db.Replace(go_utils.Low_Priority, table, col, val)
	case ActionInsertIgnore:
		err = db.Insert(true, table, col, val)
	case ActionReplace:
		err = db.Replace("", table, col, val)
	case ActionInsert:
		err = db.Insert(false, table, col, val)
	}

	if err != nil {
		c.Logger.Error("table %s", table)
		c.Logger.Error("col %s", col)
		values := strings.Split(val, "),(")
		for k, v := range values {
			c.Logger.Error("row %v - %s", k, v)
		}
		return err
	}
	return nil
}
