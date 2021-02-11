package configuration

import (
	"fmt"
	"go-domotique/models"
	"go-domotique/utils"
	"go-domotique/logger"
	"github.com/Mimerel/go-utils"
)

func getCronTab(config *models.Configuration) (err error) {
	db := utils.CreateDbConnection(config)
	db.Table = utils.TableCronTab
	db.WhereClause = ""
	db.Seperator = ","
	db.Debug = true
	db.DataType = new([]models.CronTab)
	res, err := go_utils.SearchInTable2(db)
	if err != nil {
		logger.Info(config, true,"ReadConfiguration","Unable to request database for daemon CronTab: %v", err)
		return err
	}
	if len(*res.(*[]models.CronTab)) > 0 {
		config.Daemon.CronTab = *res.(*[]models.CronTab)
		return nil
	}
	return fmt.Errorf("Unable to find list daemon CronTab")
}

