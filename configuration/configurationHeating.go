package configuration

import (
	"fmt"
	"go-domotique/models"
	"go-domotique/utils"
	"go-domotique/logger"
	"github.com/Mimerel/go-utils"
)

//

func executeHeatingConfiguration(config *models.Configuration) {
	getHeatingProgram(config)
}

func getHeatingProgram(config *models.Configuration) (err error){
	db := utils.CreateDbConnection(config)
	db.Table = utils.TableHeatingProgram
	db.FullRequest = "SELECT heatingProgram.day as dayId, days.name as day, moment, heatingLevels.name as levelName, heatingLevels.value as levelValue FROM `heatingProgram` join heatingInstructions on heatingProgram.modelId = heatingInstructions.modelId join days on heatingProgram.day = days.id join heatingLevels on heatingInstructions.levelId=heatingLevels.id order by dayId, moment"
	db.Debug = false
	db.DataType = new([]models.HeatingProgram)
	res, err := go_utils.SearchInTable(db)
	if err != nil {
		logger.Error(config, "getHeatingProgram", "Unable to request database : %v", err)
		return err
	}
	if len(*res.(*[]models.HeatingProgram)) > 0 {
		config.Heating.HeatingProgram = *res.(*[]models.HeatingProgram)
		return nil
	}
	return fmt.Errorf("Unable to find heating program")
}

func getHeatingGlobals(config *models.Configuration) (err error){
	db := utils.CreateDbConnection(config)
	db.Table = utils.TableHeating
	db.FullRequest = "SELECT * from " + utils.TableHeating
	db.Debug = false
	db.DataType = new([]models.HeatingSettings)
	res, err := go_utils.SearchInTable(db)
	if err != nil {
		logger.Error(config, "getHeatingGlobals", "Unable to request database : %v", err)
		return err
	}
	if len(*res.(*[]models.HeatingSettings)) > 0 {
		config.Heating.HeatingSettings = (*res.(*[]models.HeatingSettings))[0]
		return nil
	}
	return fmt.Errorf("Unable to find heating program")
}