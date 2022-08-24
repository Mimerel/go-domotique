package utils

import (
	"fmt"
	"github.com/Mimerel/go-utils"
	"go-domotique/models"
	"strconv"
)

func GetLastDeviceValues(config *models.Configuration) (err error) {
	db := CreateDbConnection(config)
	db.Table = models.TableDevicesLastValues
	db.WhereClause = ""
	db.Seperator = ","
	db.Debug = false
	db.DataType = new([]models.ElementDetails)
	res, err := go_utils.SearchInTable2(db)
	if err != nil {
		return err
	}
	if len(*res.(*[]models.ElementDetails)) > 0 {
		config.Devices.LastValues = *res.(*[]models.ElementDetails)
		return nil
	}
	return fmt.Errorf("Unable to find Devices values")
}

func GetLastDeviceValue(config *models.Configuration, domotiqueId int64, instanceId int64) (device models.ElementDetails) {
	db := CreateDbConnection(config)
	db.Table = models.TableDevicesLastValues
	db.WhereClause = " domotiqueId = " + strconv.FormatInt(domotiqueId, 10) + " and instanceId = " + strconv.FormatInt(instanceId, 10)
	db.Seperator = ","
	db.Debug = false
	db.DataType = new([]models.ElementDetails)
	res, err := go_utils.SearchInTable2(db)
	if err != nil {
		return models.ElementDetails{}
	}
	if len(*res.(*[]models.ElementDetails)) > 0 {
		values := *res.(*[]models.ElementDetails)
		return values[0]
	}
	return models.ElementDetails{}
}
