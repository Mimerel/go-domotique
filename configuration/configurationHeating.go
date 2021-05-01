package configuration

import (
	"fmt"
	"github.com/Mimerel/go-utils"
	"go-domotique/logger"
	"go-domotique/models"
	"go-domotique/utils"
)

//

func executeHeatingConfiguration(config *models.Configuration) {
	err := getHeatingProgram(config)
	if err != nil {
		config.Logger.Error("Error getting heating program %v", err)
	}
	err = getHeatingGlobals(config)
	if err != nil {
		config.Logger.Error("Error getting heating program %v", err)
	}
	err = getHeatingLevels(config)
	if err != nil {
		config.Logger.Error("Error getting heating program %v", err)
	}
}

func getHeatingProgram(config *models.Configuration) (err error){
	db := utils.CreateDbConnection(config)
	db.Table = utils.TableHeatingProgram
	db.FullRequest = "SELECT heatingProgram.day as dayId, days.name as day, moment, heatingLevels.name as levelName, heatingLevels.value as levelValue FROM `heatingProgram` join heatingInstructions on heatingProgram.modelId = heatingInstructions.modelId join days on heatingProgram.day = days.id join heatingLevels on heatingInstructions.levelId=heatingLevels.id order by dayId, moment"
	db.Debug = false
	db.DataType = new([]models.HeatingProgram)
	res, err := go_utils.SearchInTable2(db)
	if err != nil {
		logger.Error(config, true,"getHeatingProgram", "Unable to request database : %v", err)
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
	res, err := go_utils.SearchInTable2(db)
	if err != nil {
		logger.Error(config, true,"getHeatingGlobals", "Unable to request database : %v", err)
		return err
	}
	if len(*res.(*[]models.HeatingSettings)) > 0 {
		config.Heating.HeatingSettings = (*res.(*[]models.HeatingSettings))[0]
		return nil
	}
	return fmt.Errorf("Unable to find heating global variables")
}

func getHeatingLevels(config *models.Configuration) (err error){
	db := utils.CreateDbConnection(config)
	db.Table = utils.TableHeatingLevels
	db.FullRequest = "SELECT * from " + utils.TableHeatingLevels
	db.Debug = false
	db.DataType = new([]models.HeatingLevels)
	res, err := go_utils.SearchInTable2(db)
	if err != nil {
		logger.Error(config, true,"getHeatingLevels", "Unable to request database : %v", err)
		return err
	}
	if len(*res.(*[]models.HeatingLevels)) > 0 {
		config.Heating.HeatingLevels = *res.(*[]models.HeatingLevels)
		return nil
	}
	return fmt.Errorf("Unable to find heating levels")
}
