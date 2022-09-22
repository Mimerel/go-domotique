package google_talk

import (
	"github.com/evalphobia/google-home-client-go/googlehome"
	"go-domotique/models"
	"time"
)

/**
Given a list of ips, and a message, this method will
loop through the list and run the method to send the message to the different
google homes.
*/
func Talk(config *models.Configuration, ips []string, message string) {
	for _, ip := range ips {
		config.Logger.Debug("talk message sent to ip : %s ", ip)
		talkIndividual(config, ip, message)
	}
}

/**
Method that send a message to the google home for the
message to be read out loud
*/
func talkIndividual(config *models.Configuration, ip string, message string) {
	cli, err := googlehome.NewClientWithConfig(googlehome.Config{
		Hostname: ip,
		Lang:     "fr",
		Accent:   "FR",
	})
	if err != nil {
		config.Logger.Error("unable to send message")
	}
	cli.SetLang("fr")
	cli.Notify(message)
	time.Sleep(3 * time.Second)
}
