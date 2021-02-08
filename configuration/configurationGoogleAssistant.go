package configuration

import (
	"fmt"
	"github.com/Mimerel/go-utils"
	"go-domotique/devices"
	"go-domotique/logger"
	"go-domotique/models"
	"go-domotique/utils"
	"strings"
)

func executeGoogleAssistantConfiguration(config *models.Configuration) {
	logger.Info(config, "executeGoogleAssistantConfiguration", "Collecting Database Data")
	err := getBoxes(config)
	if err != nil {
		panic(err)
	}
	logger.Info(config, "executeGoogleAssistantConfiguration", "Collecting Google Words information")
	err = getWords(config)
	if err != nil {
		panic(err)
	}
	logger.Info(config, "executeGoogleAssistantConfiguration", "Collecting Google Instructions")
	err = getGoogleInstructions(config)
	if err != nil {
		panic(err)
	}
	logger.Info(config, "executeGoogleAssistantConfiguration", "Collecting Google Action Names")
	err = getActionNames(config)
	if err != nil {
		panic(err)
	}
	logger.Info(config, "executeGoogleAssistantConfiguration", "Collecting Rooms")
	err = getRooms(config)
	if err != nil {
		panic(err)
	}
	logger.Info(config, "executeGoogleAssistantConfiguration", "Collecting Device Types")
	err = getDeviceTypes(config)
	if err != nil {
		panic(err)
	}
	logger.Info(config, "executeGoogleAssistantConfiguration", "Collecting Google Boxes")
	err = getGoogleBox(config)
	if err != nil {
		panic(err)
	}
	logger.Info(config, "executeGoogleAssistantConfiguration", "Collecting Google Action Types")
	err = getGoogleActionTypes(config)
	if err != nil {
		panic(err)
	}
	logger.Info(config, "executeGoogleAssistantConfiguration", "Collecting Google Type words")
	err = getGoogleActionTypesWords(config)
	if err != nil {
		panic(err)
	}

}

func SaveGoogleConfigToDataBase(config *models.Configuration) {
	db := utils.CreateDbConnection(config)
	db.Debug = false
	logger.Info(config, "SaveGoogleConfigToDataBase", "Emptied instructions")
	db.Request("delete from " + utils.TableGoogleTranslatedInstructions)
	logger.Info(config, "SaveGoogleConfigToDataBase", "saving instructions")
	err := utils.ActionInMariaDB(config, config.GoogleAssistant.GoogleTranslatedInstructions, utils.TableGoogleTranslatedInstructions, utils.ActionInsertIgnore)
	if err != nil {
		logger.Error(config, "SaveGoogleConfigToDataBase", "Unable to store request model in MariaDB : %+v", err)
	}
}

func getBoxes(config *models.Configuration) (err error) {
	db := utils.CreateDbConnection(config)
	db.Table = utils.TableDomotiqueBox
	db.WhereClause = ""
	db.Debug = false
	db.DataType = new([]models.Zwave)
	res, err := go_utils.SearchInTable2(db)
	if err != nil {
		logger.Error(config, "getBoxes", "Unable to request database : %v", err)
		return err
	}
	if len(*res.(*[]models.Zwave)) > 0 {
		config.Zwaves = *res.(*[]models.Zwave)
		return nil
	}
	return fmt.Errorf("Unable to find list of Zwave Boxes")
}

func getGoogleActionTypes(config *models.Configuration) (err error) {
	db := utils.CreateDbConnection(config)
	db.Table = utils.TableGoogleActionTypes
	db.WhereClause = ""
	db.Debug = false
	db.DataType = new([]models.GoogleActionTypes)
	res, err := go_utils.SearchInTable2(db)
	if err != nil {
		logger.Error(config, "getGoogleActionTypes", "Unable to request database : %v", err)
		return err
	}
	if len(*res.(*[]models.GoogleActionTypes)) > 0 {
		config.GoogleAssistant.GoogleActionTypes = *res.(*[]models.GoogleActionTypes)
		return nil
	}
	return fmt.Errorf("Unable to find list of Zwave Boxes")
}

func getGoogleActionTypesWords(config *models.Configuration) (err error) {
	db := utils.CreateDbConnection(config)
	db.Table = utils.TableGoogleActionTypesWords
	db.WhereClause = ""
	db.Debug = false
	db.DataType = new([]models.GoogleActionTypesWords)
	res, err := go_utils.SearchInTable2(db)
	if err != nil {
		logger.Error(config, "getGoogleActionTypesWords", "Unable to request database : %v", err)
		return err
	}
	if len(*res.(*[]models.GoogleActionTypesWords)) > 0 {
		config.GoogleAssistant.GoogleActionTypesWords = *res.(*[]models.GoogleActionTypesWords)
		return nil
	}
	return fmt.Errorf("Unable to find list of Zwave Boxes")
}

func getWords(config *models.Configuration) (err error) {
	db := utils.CreateDbConnection(config)
	db.Table = utils.TableGoogleWords
	db.WhereClause = ""
	db.Seperator = ","
	db.Debug = false
	db.DataType = new([]models.GoogleWords)
	res, err := go_utils.SearchInTable2(db)
	if err != nil {
		logger.Error(config, "getWords", "Unable to request database for words: %v", err)
		return err
	}
	if len(*res.(*[]models.GoogleWords)) > 0 {
		config.GoogleAssistant.GoogleWords = *res.(*[]models.GoogleWords)
		return nil
	}
	for k, _ := range config.GoogleAssistant.GoogleWords {
		config.GoogleAssistant.GoogleWords[k].WordsConverted = strings.ToLower(strings.Replace(config.GoogleAssistant.GoogleWords[k].Words, " ", "", -1))
		for _, v := range config.CharsToReplace {
			config.GoogleAssistant.GoogleWords[k].WordsConverted = strings.Replace(config.GoogleAssistant.GoogleWords[k].WordsConverted, v.From, v.To, -1)
		}
	}

	return fmt.Errorf("Unable to find list of words")
}

func getGoogleInstructions(config *models.Configuration) (err error) {
	db := utils.CreateDbConnection(config)
	db.Table = utils.TableGoogleInstructions
	db.WhereClause = ""
	db.Seperator = ","
	db.Debug = false
	db.DataType = new([]models.GoogleInstruction)
	res, err := go_utils.SearchInTable2(db)
	if err != nil {
		logger.Error(config, "getGoogleInstructions", "Unable to request database for words: %v", err)
		return err
	}
	if len(*res.(*[]models.GoogleInstruction)) > 0 {
		config.GoogleAssistant.GoogleInstructions = *res.(*[]models.GoogleInstruction)
		return nil
	}
	return fmt.Errorf("Unable to find list of words")
}

func getActionNames(config *models.Configuration) (err error) {
	db := utils.CreateDbConnection(config)
	db.Table = utils.TableGoogleActionNames
	db.WhereClause = ""
	db.Seperator = ","
	db.Debug = false
	db.DataType = new([]models.GoogleActionNames)
	res, err := go_utils.SearchInTable2(db)
	if err != nil {
		logger.Error(config, "getActionNames", "Unable to request database for words: %v", err)
		return err
	}
	if len(*res.(*[]models.GoogleActionNames)) > 0 {
		config.GoogleAssistant.GoogleActionNames = *res.(*[]models.GoogleActionNames)
		return nil
	}
	return fmt.Errorf("Unable to find list of words")
}

func getRooms(config *models.Configuration) (err error) {
	db := utils.CreateDbConnection(config)
	db.Table = utils.TableRooms
	db.WhereClause = ""
	db.Seperator = ","
	db.Debug = false
	db.DataType = new([]models.Room)
	res, err := go_utils.SearchInTable2(db)
	if err != nil {
		logger.Error(config, "getRooms", "Unable to request database for words: %v", err)
		return err
	}
	if len(*res.(*[]models.Room)) > 0 {
		config.Rooms = *res.(*[]models.Room)
		return nil
	}
	return fmt.Errorf("Unable to find list of words")
}

func getGoogleBox(config *models.Configuration) (err error) {
	db := utils.CreateDbConnection(config)
	db.Table = utils.TableGoogleBox
	db.WhereClause = ""
	db.Seperator = ","
	db.Debug = false
	db.DataType = new([]models.GoogleBox)
	res, err := go_utils.SearchInTable2(db)
	if err != nil {
		logger.Error(config, "getGoogleBox", "Unable to request database for words: %v", err)
		return err
	}
	if len(*res.(*[]models.GoogleBox)) > 0 {
		config.GoogleAssistant.GoogleBoxes = *res.(*[]models.GoogleBox)
		return nil
	}
	return fmt.Errorf("Unable to find list of words")
}

func getDeviceTypes(config *models.Configuration) (err error) {
	db := utils.CreateDbConnection(config)
	db.Table = utils.TableDeviceTypes
	db.WhereClause = ""
	db.Seperator = ","
	db.Debug = false
	db.DataType = new([]models.DeviceType)
	res, err := go_utils.SearchInTable2(db)
	if err != nil {
		logger.Error(config, "getDeviceTypes", "Unable to request database for words: %v", err)
		return err
	}
	if len(*res.(*[]models.DeviceType)) > 0 {
		config.DeviceTypes = *res.(*[]models.DeviceType)
		return nil
	}
	return fmt.Errorf("Unable to find list of words")
}

/**
Method that checks that the configuration file is consistent.
If a device name that does not exist in the device list
or if a zwave device that does not exist in the zwave list
are used, error message will be displayed and the program will stop
*/
func CheckGoogleConfiguration(config *models.Configuration) {
	// Check if devices that are used in commands are in the device list
	//[]GoogleTranslatedInstruction

	for _, v := range config.GoogleAssistant.GoogleActionTypesWords {
		translated := new(models.GoogleTranslatedActionTypes)
		translated.ActionWord = v.Action
		translated.Action = getActionTypeFromId(config, v.ActionTypeId).Name
		config.GoogleAssistant.GoogleTranslatedActionTypes = append(config.GoogleAssistant.GoogleTranslatedActionTypes, *translated)
	}

	for _, instruction := range config.GoogleAssistant.GoogleInstructions {
		translated := new(models.GoogleTranslatedInstruction)
		translated.Value = instruction.Value
		translated.Id = instruction.Id
		translated.ActionName = getActionNameFromId(config, instruction.ActionNameId).Name
		translated.ActionNameId = instruction.ActionNameId
		translated.CommandClass = devices.GetDeviceFromId(config, instruction.DomotiqueId).CommandClass
		translated.Instance = devices.GetDeviceFromId(config, instruction.DomotiqueId).Instance
		translated.DeviceName = devices.GetDeviceFromId(config, instruction.DomotiqueId).Name
		translated.DeviceId = devices.GetDeviceFromId(config, instruction.DomotiqueId).DeviceId
		translated.Room = devices.GetDeviceFromId(config, instruction.DomotiqueId).Room
		translated.ZwaveId = devices.GetDeviceFromId(config, instruction.DomotiqueId).Zwave
		translated.ZwaveUrl = devices.GetDeviceFromId(config, instruction.DomotiqueId).ZwaveUrl
		translated.GoogleBox = getGoogleBoxFromId(config, instruction.GoogleBoxId).Name
		translated.TypeDevice = devices.GetDeviceFromId(config, instruction.DomotiqueId).Type
		translated.Type = getActionTypeFromId(config, instruction.TypeId).Name
		config.GoogleAssistant.GoogleTranslatedInstructions = append(config.GoogleAssistant.GoogleTranslatedInstructions, *translated)
	}

}

func getRoomFromId(config *models.Configuration, id int64) models.Room {
	for _, v := range config.Rooms {
		if v.Id == id {
			return v
		}
	}
	return models.Room{}
}

func getZwaveFromId(config *models.Configuration, id int64) models.Zwave {
	for _, v := range config.Zwaves {
		if v.Id == id {
			return v
		}
	}
	return models.Zwave{}
}

func getGoogleBoxFromId(config *models.Configuration, id int64) models.GoogleBox {
	for _, v := range config.GoogleAssistant.GoogleBoxes {
		if v.Id == id {
			return v
		}
	}
	return models.GoogleBox{}
}

func getTypeFromId(config *models.Configuration, id int64) models.DeviceType {
	for _, v := range config.DeviceTypes {
		if v.Id == id {
			return v
		}
	}
	return models.DeviceType{}
}

func getActionNameFromId(config *models.Configuration, id int64) models.GoogleActionNames {
	for _, v := range config.GoogleAssistant.GoogleActionNames {
		if v.Id == id {
			return v
		}
	}
	return models.GoogleActionNames{}
}

func getActionTypeFromId(config *models.Configuration, id int64) models.GoogleActionTypes {
	for _, v := range config.GoogleAssistant.GoogleActionTypes {
		if v.Id == id {
			return v
		}
	}
	return models.GoogleActionTypes{}
}
