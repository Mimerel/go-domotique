package prowl

import (
	"go-domotique/models"
	"net/http"
	"net/url"
	"time"
)

/**
Sends Prowl notification
*/
func SendProwlNotification(config *models.Configuration, AppName string, Event string, Description string) {
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	params := url.Values{}
	params.Add("apikey", config.Token)
	params.Add("application", AppName)
	params.Add("event", Event)
	params.Add("description", Description)
	params.Add("priority", "1")
	postingUrl := "https://api.prowlapp.com/publicapi/add?" + params.Encode()
	_, err := client.Get(postingUrl)
	if err != nil {
		config.Logger.Error("Unable to post prown notification %s - %s - %s", AppName, Event, Description)
	} else {
		config.Logger.Info("Prowl notification sent %s", postingUrl)
	}
}
