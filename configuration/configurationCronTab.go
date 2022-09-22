package configuration

import (
	"fmt"
	"github.com/Mimerel/go-utils"
	"go-domotique/models"
	"go-domotique/utils"
)

func getCronTab(config *models.Configuration) (err error) {
	db := utils.CreateDbConnection(config)
	db.Table = models.TableCronTab
	db.WhereClause = ""
	db.Seperator = ","
	db.Debug = false
	db.DataType = new([]models.CronTab)
	res, err := go_utils.SearchInTable2(db)
	if err != nil {
		config.Logger.Info("ReadConfiguration Unable to request database for daemon CronTab: %v", err)
		return err
	}
	if len(*res.(*[]models.CronTab)) > 0 {
		config.Daemon.CronTab = *res.(*[]models.CronTab)
		return nil
	}
	return fmt.Errorf("Unable to find list daemon CronTab")
}
