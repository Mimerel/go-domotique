package configuration

import (
	"fmt"
	"github.com/Mimerel/go-utils"
	"go-domotique/logger"
	"go-domotique/models"
	"go-domotique/utils"
	"sort"
	"strconv"
)

func getListDevices(config *models.Configuration) {
	logger.Info(config, false, "getListDevices", "Collecting Devices")
	err := getDevices(config)
	if err != nil {
		panic(err)
	}
}

func getDevices(config *models.Configuration) (err error) {
	db := utils.CreateDbConnection(config)
	db.Table = utils.TableDevices
	db.WhereClause = ""
	db.Seperator = ","
	db.Debug = false
	db.DataType = new([]models.DeviceDetails)
	res, err := go_utils.SearchInTable2(db)
	if err != nil {
		logger.Error(config, false,"getDevices", "Unable to request database for devices: %v", err)
		return err
	}
	if len(*res.(*[]models.DeviceDetails)) > 0 {
		config.Devices.Devices = *res.(*[]models.DeviceDetails)
		return nil
	}
	return fmt.Errorf("Unable to find list of Devices")
}

func CheckConfigurationDevices(config *models.Configuration) {
	// Check if devices that are used in commands are in the device list
	//[]GoogleTranslatedInstruction

	for _, device := range config.Devices.Devices {
		translated := new(models.DeviceTranslated)
		translated.DomotiqueId = device.DomotiqueId
		translated.Instance = device.Instance
		translated.InstanceString = strconv.Itoa(int(device.Instance))
		translated.CommandClass = device.CommandClass
		translated.DeviceId = device.DeviceId
		translated.DeviceIdString = strconv.Itoa(int(device.DeviceId))
		translated.Room = getRoomFromId(config, device.RoomId).Name
		translated.Type = getTypeFromId(config, device.TypeId).Name
		translated.Name = device.Name
		translated.BoxId = device.BoxId
		translated.ZwaveName = getZwaveFromId(config, device.BoxId).Name
		translated.ZwaveUrl = getZwaveFromId(config, device.BoxId).Ip
		translated.TypeWifi = device.TypeWifi
		config.Devices.DevicesTranslated = append(config.Devices.DevicesTranslated, *translated)
		if device.OnUi == 1 {
			config.Devices.DevicesToggle = append(config.Devices.DevicesToggle, translated.CollectDeviceToggleDetails(config))
		}
	}
	sort.Slice(config.Devices.DevicesToggle, func(a, b int) bool {
		return config.Devices.DevicesToggle[a].Room + config.Devices.DevicesToggle[a].Name > config.Devices.DevicesToggle[b].Room + config.Devices.DevicesToggle[b].Name
	})

}

func SaveDevicesToDataBase(config *models.Configuration) {
	db := utils.CreateDbConnection(config)
	db.Debug = false
	logger.Info(config, false, "SaveDevicesToDataBase", "Emptied devicestranslated")
	_ = db.Request("delete from " + utils.TableDevicesTranslated)
	logger.Info(config, false, "SaveDevicesToDataBase", "saving Devices")
	err := utils.ActionInMariaDB(config, config.Devices.DevicesTranslated, utils.TableDevicesTranslated, utils.ActionInsertIgnore)
	if err != nil {
		logger.Error(config, true, "SaveDevicesToDataBase", "Unable to store request model in MariaDB : %+v", err)
	}
}

func getDeviceActions(config *models.Configuration) {
	db := utils.CreateDbConnection(config)
	db.Table = utils.TableDeviceActions
	db.WhereClause = ""
	db.Seperator = ","
	db.Debug = false
	db.DataType = new([]models.DeviceActions)
	res, err := go_utils.SearchInTable2(db)
	if err != nil {
		logger.Error(config, true,"getDeviceActions", "Unable to request database for device Actions: %v", err)
		return
	}
	if len(*res.(*[]models.DeviceActions)) > 0 {
		config.Devices.DevicesActions = *res.(*[]models.DeviceActions)
	}
}