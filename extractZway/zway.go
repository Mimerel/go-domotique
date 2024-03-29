package extractZway

import (
	"encoding/json"
	"go-domotique/devices"
	"go-domotique/models"
	"go-domotique/utils"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func getDataFromZWay(config *models.Configuration, url string) (data models.ZwaveExtractionData) {
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	res, err := client.Get("http://" + url + ":8083/ZWaveAPI/Data")
	if err != nil {
		config.Logger.Error("There was a get site error:", err)
	} else {

		temp, err := ioutil.ReadAll(res.Body)
		if err != nil {
			config.Logger.Error("There was a error while reading the body of zway request error:", err)
		}

		res.Body.Close()

		err = json.Unmarshal(temp, &data.Json)
		if err != nil {
			config.Logger.Error("error decoding zway response: %v", err)
		}
	}
	return data
}

func ExtractZWayMetrics(config *models.Configuration) {
	lastValues := new([]models.ElementDetails)
	for _, v := range config.Zwaves {
		if v.Id == 100 {
			continue
		}
		//logger.Debug(config, false, "ExtractZWayMetrics", "Requesting data from %v", v.Ip)
		data := getDataFromZWay(config, v.Ip)
		elements := extractElements(config, data, v.Id)
		*lastValues = append(*lastValues, elements...)
		//logger.Debug(config, false, "ExtractZWayMetrics", "found %v - total %v", len(elements), len(*lastValues))
	}
	config.Devices.LastValues = *lastValues
	if len(*lastValues) > 0 {
		saveExtractZwaveDataToDataBase(config)
	}
}

func saveOnlyLastValues(config *models.Configuration, col string, val string) {
	db := utils.CreateDbConnection(config)
	db.Debug = false
	err := db.Request("delete from " + models.TableDevicesLastValues)
	if err != nil {
		config.Logger.Error("Unable to clear last device values DB")
	}

	err = db.Insert(false, models.TableDevicesLastValues, col, val)

	if err != nil {
		config.Logger.Error("table %s", models.LoggerDomotique)
		config.Logger.Error("col %s", col)
		values := strings.Split(val, "),(")
		for k, v := range values {
			config.Logger.Error("row %v - %s", k, v)
		}
	}
}

func saveExtractZwaveDataToDataBase(config *models.Configuration) {
	db := utils.CreateDbConnection(config)
	db.Database = models.DatabaseStats
	db.Debug = false
	col, val, err := db.DecryptStructureAndData(config.Devices.LastValues)
	if err != nil {
		config.Logger.Error("col %s", col)
		config.Logger.Error("val %s", val)
	}
	err = db.Insert(false, models.TableDevicesLastValues, col, val)

	if err != nil {
		config.Logger.Error("table %s", models.LoggerDomotique)
		config.Logger.Error("col %s", col)
		values := strings.Split(val, "),(")
		for k, v := range values {
			config.Logger.Error("row %v - %s", k, v)
		}
		return
	}
	saveOnlyLastValues(config, col, val)
}

func extractElements(config *models.Configuration, data models.ZwaveExtractionData, boxId int64) (elements []models.ElementDetails) {
	const createdFormat = "2006-01-02 15:04:05"
	for deviceId, v := range data.Json.Devices {
		deviceIdInInt, err := strconv.ParseInt(deviceId, 10, 64)
		if err != nil {
			config.Logger.Error("unable to convert deviceId to int")
			continue
		}
		domotiqueId := devices.GetDomotiqueIdFromDeviceIdAndBoxId(config, deviceIdInInt, boxId).DomotiqueId
		if domotiqueId == 0 {
			continue
		}
		for instanceKey, instanceContent := range v.Instances {
			var useClass50 bool
			if (instanceContent.CommandClasses.Class50.Data.Data2 != (models.CommandClass50DataVal{}) &&
				instanceContent.CommandClasses.Class50.Data.Data2.Val.Value > 0) ||
				instanceContent.CommandClasses.Class49.Data.Data4 == (models.CommandClass49DataVal{}) {
				useClass50 = true
			} else {
				useClass50 = false
			}
			if instanceContent.CommandClasses.Class50 != (models.CommandClass50{}) {
				if instanceContent.CommandClasses.Class50.Data.Data2 != (models.CommandClass50DataVal{}) {
					element := new(models.ElementDetails)
					element.Unit = "Watt"
					element.Value = instanceContent.CommandClasses.Class50.Data.Data2.Val.Value
					element.DeviceId = deviceIdInInt
					element.InstanceId, err = strconv.ParseInt(instanceKey, 10, 64)
					element.Timestamp = time.Unix(time.Now().In(config.Location).Unix(), 0).Format(createdFormat)
					if err != nil {
						config.Logger.Error("unable to convert instance from string to int")
					}
					element.DomotiqueId = domotiqueId
					if useClass50 == true {
						elements = append(elements, *element)
					}
				}
				if instanceContent.CommandClasses.Class50.Data.Data4 != (models.CommandClass50DataVal{}) {
					element := new(models.ElementDetails)
					element.Unit = "Volts"
					element.Value = instanceContent.CommandClasses.Class50.Data.Data4.Val.Value
					element.DeviceId = deviceIdInInt
					element.InstanceId, err = strconv.ParseInt(instanceKey, 10, 64)
					if err != nil {
						config.Logger.Error("unable to convert instance from string to int")
					}
					element.Timestamp = time.Unix(time.Now().In(config.Location).Unix(), 0).Format(createdFormat)
					element.DomotiqueId = domotiqueId
					elements = append(elements, *element)
				}
				if instanceContent.CommandClasses.Class50.Data.Data5 != (models.CommandClass50DataVal{}) {
					element := new(models.ElementDetails)
					element.Unit = "Ampères"
					element.Value = instanceContent.CommandClasses.Class50.Data.Data5.Val.Value
					element.InstanceId, err = strconv.ParseInt(instanceKey, 10, 64)
					if err != nil {
						config.Logger.Error("unable to convert instance from string to int")
					}
					element.DeviceId = deviceIdInInt
					element.Timestamp = time.Unix(time.Now().In(config.Location).Unix(), 0).Format(createdFormat)
					element.DomotiqueId = domotiqueId
					elements = append(elements, *element)
				}
			}
			if instanceContent.CommandClasses.Class49 != (models.CommandClass49{}) {
				if instanceContent.CommandClasses.Class49.Data.Data1 != (models.CommandClass49DataVal{}) {
					element := new(models.ElementDetails)
					element.Unit = "Degré"
					element.Value = instanceContent.CommandClasses.Class49.Data.Data1.Val.Value
					element.InstanceId, err = strconv.ParseInt(instanceKey, 10, 64)
					if err != nil {
						config.Logger.Error("unable to convert instance from string to int")
					}
					element.DeviceId = deviceIdInInt
					element.Timestamp = time.Unix(time.Now().In(config.Location).Unix(), 0).Format(createdFormat)
					element.DomotiqueId = domotiqueId
					elements = append(elements, *element)
				}
				if instanceContent.CommandClasses.Class49.Data.Data3 != (models.CommandClass49DataVal{}) {
					element := new(models.ElementDetails)
					element.Unit = "Lux"
					element.Value = instanceContent.CommandClasses.Class49.Data.Data3.Val.Value
					element.InstanceId, err = strconv.ParseInt(instanceKey, 10, 64)
					if err != nil {
						config.Logger.Error("unable to convert instance from string to int")
					}
					element.DeviceId = deviceIdInInt
					element.Timestamp = time.Unix(time.Now().In(config.Location).Unix(), 0).Format(createdFormat)
					element.DomotiqueId = domotiqueId
					elements = append(elements, *element)
				}
				if instanceContent.CommandClasses.Class49.Data.Data4 != (models.CommandClass49DataVal{}) {
					element := new(models.ElementDetails)
					element.Unit = "Watt"
					element.Value = instanceContent.CommandClasses.Class49.Data.Data4.Val.Value
					element.InstanceId, err = strconv.ParseInt(instanceKey, 10, 64)
					if err != nil {
						config.Logger.Error("unable to convert instance from string to int")
					}
					element.DeviceId = deviceIdInInt
					element.Timestamp = time.Unix(time.Now().In(config.Location).Unix(), 0).Format(createdFormat)
					element.DomotiqueId = domotiqueId
					if useClass50 == false {
						elements = append(elements, *element)
					}
				}
				if instanceContent.CommandClasses.Class49.Data.Data5 != (models.CommandClass49DataVal{}) {
					element := new(models.ElementDetails)
					element.Unit = "Humidity"
					element.Value = instanceContent.CommandClasses.Class49.Data.Data5.Val.Value
					element.InstanceId, err = strconv.ParseInt(instanceKey, 10, 64)
					if err != nil {
						config.Logger.Error("unable to convert instance from string to int")
					}
					element.DeviceId = deviceIdInInt
					element.Timestamp = time.Unix(time.Now().In(config.Location).Unix(), 0).Format(createdFormat)
					element.DomotiqueId = domotiqueId
					elements = append(elements, *element)
				}
			}
			if instanceContent.CommandClasses.Class48 != (models.CommandClass48{}) {
				if instanceContent.CommandClasses.Class48.Data.Data1 != (models.CommandClass48DataValBool{}) {
					element := new(models.ElementDetails)
					element.Unit = "Alarm"
					element.Value = BoolToIntensity(instanceContent.CommandClasses.Class48.Data.Data1.Level.Value)
					element.InstanceId, err = strconv.ParseInt(instanceKey, 10, 64)
					if err != nil {
						config.Logger.Error("unable to convert instance from string to int")
					}
					element.DeviceId = deviceIdInInt
					element.Timestamp = time.Unix(time.Now().In(config.Location).Unix(), 0).Format(createdFormat)
					element.DomotiqueId = domotiqueId
					elements = append(elements, *element)
				}
				if instanceContent.CommandClasses.Class48.Data.Data6 != (models.CommandClass48DataValBool{}) {
					element := new(models.ElementDetails)
					element.Unit = "Flood"
					element.Value = BoolToIntensity(instanceContent.CommandClasses.Class48.Data.Data6.Level.Value)
					element.InstanceId, err = strconv.ParseInt(instanceKey, 10, 64)
					if err != nil {
						config.Logger.Error("unable to convert instance from string to int")
					}
					element.DeviceId = deviceIdInInt
					element.Timestamp = time.Unix(time.Now().In(config.Location).Unix(), 0).Format(createdFormat)
					element.DomotiqueId = domotiqueId
					elements = append(elements, *element)
				}
				if instanceContent.CommandClasses.Class48.Data.Data8 != (models.CommandClass48DataValBool{}) {
					element := new(models.ElementDetails)
					element.Unit = "Tempered"
					element.Value = BoolToIntensity(instanceContent.CommandClasses.Class48.Data.Data8.Level.Value)
					element.InstanceId, err = strconv.ParseInt(instanceKey, 10, 64)
					if err != nil {
						config.Logger.Error("unable to convert instance from string to int")
					}
					element.DeviceId = deviceIdInInt
					element.Timestamp = time.Unix(time.Now().In(config.Location).Unix(), 0).Format(createdFormat)
					element.DomotiqueId = domotiqueId
					elements = append(elements, *element)
				}
				if instanceContent.CommandClasses.Class48.Data.Data12 != (models.CommandClass48DataValBool{}) {
					element := new(models.ElementDetails)
					element.Unit = "Tempered"
					element.Value = BoolToIntensity(instanceContent.CommandClasses.Class48.Data.Data12.Level.Value)
					element.InstanceId, err = strconv.ParseInt(instanceKey, 10, 64)
					if err != nil {
						config.Logger.Error("unable to convert instance from string to int")
					}
					element.DeviceId = deviceIdInInt
					element.Timestamp = time.Unix(time.Now().In(config.Location).Unix(), 0).Format(createdFormat)
					element.DomotiqueId = domotiqueId
					elements = append(elements, *element)
				}
			}
			if instanceContent.CommandClasses.Class37 != (models.CommandClass37{}) {
				element := new(models.ElementDetails)
				element.Unit = "Level"
				element.Value = BoolToIntensity(instanceContent.CommandClasses.Class37.Data.Level.Value)

				element.Switch = "fix"
				element.InstanceId, err = strconv.ParseInt(instanceKey, 10, 64)
				if err != nil {
					config.Logger.Error("unable to convert instance from string to int")
				}
				element.DeviceId = deviceIdInInt
				element.Timestamp = time.Unix(time.Now().In(config.Location).Unix(), 0).Format(createdFormat)
				element.DomotiqueId = domotiqueId
				elements = append(elements, *element)
			}
			if instanceContent.CommandClasses.Class38 != (models.CommandClass38{}) {
				element := new(models.ElementDetails)
				element.Unit = "Level"
				element.Value = instanceContent.CommandClasses.Class38.Data.Level.Value
				element.Switch = "variable"
				element.InstanceId, err = strconv.ParseInt(instanceKey, 10, 64)
				if err != nil {
					config.Logger.Error("unable to convert instance from string to int")
				}
				element.DeviceId = deviceIdInInt
				element.Timestamp = time.Unix(time.Now().Unix(), 0).Format(createdFormat)
				element.DomotiqueId = domotiqueId
				elements = append(elements, *element)
			}
		}
	}
	return elements
}
