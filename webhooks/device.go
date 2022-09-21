package webhooks

import (
	"go-domotique/models"
	"net/http"
	"strings"
)

func DeviceWebhookEndpoint(w http.ResponseWriter, r *http.Request, c *models.Configuration) {
	urlPath := r.URL.Path
	urlParams := strings.Split(urlPath, "/")
	for k, v := range urlParams {
		c.Logger.Info("URL params %v - %+v", k, v)
	}

}
