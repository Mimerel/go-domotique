package controller

import (
	"fmt"
	"go-domotique/models"
	"go-domotique/webhooks"
	"net/http"
)

func getControllerWebHooks(config *models.Configuration) {
	http.HandleFunc("/device/webhook", func(w http.ResponseWriter, r *http.Request) {

		go webhooks.DeviceWebhookEndpoint(w, r, config)

		w.WriteHeader(200)
		_, _ = fmt.Fprintf(w, "Successfully received webhook information")
		return
	})
}
