package configuration

import (
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
	query := `
		Select 
		    IFNULL(domotiqueid, 0),
		    IFNULL(name, ''),
		    IFNULL(boxid, 0),
		    IFNULL(roomid, 0),
		    IFNULL(typeid, 0),
		    IFNULL(deviceid, 0),
		    IFNULL(instance, 0),
		    IFNULL(commandClass, 0),
		    IFNULL(onUi, 0),
		    IFNULL(typeWifi, ''),
		    IFNULL(model, ''),
		    IFNULL(parentId, 0)
		    from devices
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
		var device models.DeviceDetails
		// for each row, scan the result into our tag composite object
		err = result.Scan(
			&device.DomotiqueId,
			&device.Name,
			&device.BoxId,
			&device.RoomId,
			&device.TypeId,
			&device.DeviceId,
			&device.Instance,
			&device.CommandClass,
			&device.OnUi,
			&device.TypeWifi,
			&device.DeviceType,
			&device.ParentId,
		)
		if err == nil {
			if device.OnUi == 0 {
				continue
			}
			if device.ParentId == 0 {
				device.ParentId = device.DomotiqueId
			}
			config.Devices.Devices = append(config.Devices.Devices, device)
		}
	}
	return nil
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
		translated.DeviceType = device.DeviceType
		translated.ParentId = device.ParentId
		config.Devices.DevicesTranslated = append(config.Devices.DevicesTranslated, *translated)
		if device.OnUi == 1 {
			config.Devices.DevicesToggle = append(config.Devices.DevicesToggle, translated.CollectDeviceToggleDetails(config))
		}
	}
	sort.Slice(config.Devices.DevicesToggle, func(a, b int) bool {
		return config.Devices.DevicesToggle[a].Room+config.Devices.DevicesToggle[a].Name > config.Devices.DevicesToggle[b].Room+config.Devices.DevicesToggle[b].Name
	})

}

func SaveDevicesToDataBase(config *models.Configuration) {
	db := utils.CreateDbConnection(config)
	db.Debug = false
	logger.Info(config, false, "SaveDevicesToDataBase", "Emptied devicestranslated")
	_ = db.Request("delete from " + models.TableDevicesTranslated)
	logger.Info(config, false, "SaveDevicesToDataBase", "saving Devices")
	err := utils.ActionInMariaDB(config, config.Devices.DevicesTranslated, models.TableDevicesTranslated, models.ActionInsertIgnore)
	if err != nil {
		logger.Error(config, true, "SaveDevicesToDataBase", "Unable to store request model in MariaDB : %+v", err)
	}
}

func getDeviceActions(config *models.Configuration) {
	db := utils.CreateDbConnection(config)
	db.Table = models.TableDeviceActions
	db.WhereClause = ""
	db.Seperator = ","
	db.Debug = false
	db.DataType = new([]models.DeviceActions)
	res, err := go_utils.SearchInTable2(db)
	if err != nil {
		logger.Error(config, true, "getDeviceActions", "Unable to request database for device Actions: %v", err)
		return
	}
	if len(*res.(*[]models.DeviceActions)) > 0 {
		config.Devices.DevicesActions = *res.(*[]models.DeviceActions)
	}
}
