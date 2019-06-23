package prowl

import (
	"go-domotique/logger"
	"go-domotique/models"
	"net/http"
	"time"
)

/**
Sends Prowl notification
 */
func SendProwlNotification(config *models.Configuration, AppName string, Event string, Description string) {
	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	postingUrl := "https://api.prowlapp.com/publicapi/add?apikey=" + config.Token + "&application=" + AppName + "&event=" + Event + "&description=" + Description + "&priority=1"
	_, err := client.Get(postingUrl)
	if err != nil {
		logger.Error(config, "SendProwlNotification", "Unable to post prown notification %s - %s - %s", AppName, Event, Description)
	} else {
		logger.Info(config, "SendProwlNotification", "Prowl notification sent")
	}
}
