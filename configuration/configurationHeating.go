package configuration

import (
	"fmt"
	"github.com/Mimerel/go-utils"
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

func getHeatingProgram(config *models.Configuration) (err error) {
	db := utils.CreateDbConnection(config)
	db.Table = models.TableHeatingProgram
	db.FullRequest = "SELECT heatingProgram.day as dayId, days.name as day, moment, heatingLevels.name as levelName, heatingLevels.value as levelValue FROM `heatingProgram` join heatingInstructions on heatingProgram.modelId = heatingInstructions.modelId join days on heatingProgram.day = days.id join heatingLevels on heatingInstructions.levelId=heatingLevels.id order by dayId, moment"
	db.Debug = false
	db.DataType = new([]models.HeatingProgram)
	res, err := go_utils.SearchInTable2(db)
	if err != nil {
		config.Logger.Error("getHeatingProgram Unable to request database : %v", err)
		return err
	}
	if len(*res.(*[]models.HeatingProgram)) > 0 {
		config.Heating.HeatingProgram = *res.(*[]models.HeatingProgram)
		return nil
	}
	return fmt.Errorf("Unable to find heating program")
}

func getHeatingGlobals(config *models.Configuration) (err error) {
	query := `
		Select 
		    IFNULL(id, 0),
		    IFNULL(module, ''),
		    IFNULL(domotiqueId, 0)
		    from heating
`

	err = config.MariaDB.DB.Check(config)
	if err != nil {
		config.Logger.Warn("Error on dbase connexion")
	}

	result, err := config.MariaDB.DB.DB.Query(query)
	if err != nil {
		config.Logger.Error("Request : %v", query)
		return err
	}
	for result.Next() {
		var device models.HeatingSettings
		// for each row, scan the result into our tag composite object
		err = result.Scan(
			&device.Id,
			&device.Module,
			&device.DomotiqueId,
		)
		if err == nil {
			config.Heating.HeatingSettings = append(config.Heating.HeatingSettings, device)
		}
	}
	return nil

}

func getHeatingLevels(config *models.Configuration) (err error) {
	db := utils.CreateDbConnection(config)
	db.Table = models.TableHeatingLevels
	db.FullRequest = "SELECT * from " + models.TableHeatingLevels
	db.Debug = false
	db.DataType = new([]models.HeatingLevels)
	res, err := go_utils.SearchInTable2(db)
	if err != nil {
		config.Logger.Error("getHeatingLevels", "Unable to request database : %v", err)
		return err
	}
	if len(*res.(*[]models.HeatingLevels)) > 0 {
		config.Heating.HeatingLevels = *res.(*[]models.HeatingLevels)
		return nil
	}
	return fmt.Errorf("Unable to find heating levels")
}
