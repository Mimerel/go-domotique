package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const (
	ActionReplaceLowPriority = "replaceLowPriority"
	ActionReplace            = "replace"
	ActionInsertIgnore       = "insertIgnore"
	ActionInsert             = "insert"
	DatabaseLogger           = "logs"
	DatabaseStats            = "domotiqueStats"
	LoggerDomotique          = "domotique"

	TableDomotiqueBox = "domotiqueBox"

	TableDevices           = "devices"
	TableDevicesTranslated = "devicesTranslated"
	TableDevicesLastValues = "devicesLastValues"
	TableDeviceTypes       = "deviceTypes"
	TableDeviceActions     = "deviceActions"

	TableRooms = "rooms"

	TableEvents = "events"

	TableHeatingProgram = "heatingProgram"
	TableHeating        = "heating"
	TableHeatingLevels  = "heatingLevels"

	TableGoogleWords                  = "googleWords"
	TableGoogleInstructions           = "googleInstructions"
	TableGoogleActionNames            = "googleActionNames"
	TableGoogleBox                    = "googleBox"
	TableGoogleActionTypes            = "googleActionTypes"
	TableGoogleActionTypesWords       = "googleActionTypesWords"
	TableGoogleTranslatedInstructions = "googleTranslatedInstructions"

	TableCronTab = "cronTab"
)

type Maria struct {
	DB *sql.DB
}

func (i *Maria) connect(c *Configuration) (err error) {
	_ = i.Close(c)
	i.DB, err = sql.Open("mysql", c.MariaDB.User+":"+c.MariaDB.Password+"@tcp("+c.MariaDB.IP+":"+c.MariaDB.Port+")/"+c.MariaDB.Database)
	if err != nil {
		c.Logger.Error("Unable to connect to database")
		return err
	}
	i.DB.SetMaxIdleConns(100)
	i.DB.SetMaxIdleConns(100)
	i.DB.SetConnMaxLifetime(time.Hour * 24)
	return nil
}

func (i *Maria) Close(c *Configuration) (err error) {
	if i.DB == nil {
		return nil
	}
	err = i.DB.Close()
	if err != nil {
		c.Logger.Error("Unable to close database")
		return err
	}
	return nil
}

func (i *Maria) Check(c *Configuration) (err error) {
	if i.DB == nil || i.DB.Ping() != nil {
		err = i.connect(c)
		if err != nil {
			return err
		}
	}
	return nil
}
