package controller

import (
	"go-domotique/models"
	"go-domotique/healthCheck"
	"net/http"
)

func healthcheckController(config *models.Configuration) {
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		healthCheck.HealthInfo(w, r, config)
	})

}