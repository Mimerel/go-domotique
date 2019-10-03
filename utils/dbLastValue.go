package utils

import (
	"fmt"
	"go-domotique/models"
	"github.com/Mimerel/go-utils"
)

func GetLastDeviceValues(config *models.Configuration) (err error) {
	db := CreateDbConnection(config)
	db.Table = TableDevicesLastValues
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
