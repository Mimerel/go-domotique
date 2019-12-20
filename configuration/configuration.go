package configuration

import (
	"fmt"
	"go-domotique/logger"
	"go-domotique/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

/**
Method that reads the configuration file.
If a environment variable is set, the program will read the configuration
file from the path provided otherwize it will use the path coded in hard
 */
func ReadConfiguration() (*models.Configuration) {
	pathToFile := os.Getenv("LOGGER_CONFIGURATION_FILE")
	if _, err := os.Stat("/home/pi/go/src/go-domotique/configuration.yaml"); !os.IsNotExist(err) {
		pathToFile = "/home/pi/go/src/go-domotique/configuration.yaml"
	} else {
		pathToFile = "./configuration.yaml"
	}
	yamlFile, err := ioutil.ReadFile(pathToFile)

	if err != nil {
		panic(err)
	}

	var config *models.Configuration
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	} else {
		config.Location, err = time.LoadLocation("Europe/Paris")
		if err != nil {
			fmt.Println(err)
		}
		getListDevices(config)
		executeGoogleAssistantConfiguration(config)
		err := getCronTab(config)
		if err != nil {
			config.Logger.Error("Error getting cron elements")
		}
		logger.Info(config, "ReadConfiguration","Checking configuration")
		CheckConfigurationDevices(config)
		SaveDevicesToDataBase(config)
		CheckGoogleConfiguration(config)
		SaveGoogleConfigToDataBase(config)
		executeHeatingConfiguration(config)
		logger.Info(config, "ReadConfiguration","Configuration Loaded : %+v ", config)
	}
	return config
}
