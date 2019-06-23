package healthCheck

import (
	"comparator/models"
	"encoding/json"
	"net/http"
)

var healthStructure []HealthDetail

func HealthInfo(w http.ResponseWriter, req *http.Request, c *models.Configuration) {
	defer func() {
		// recover from panic if one occurred. Set err to nil otherwise.
		if (recover() != nil) {
			c.Logger.Error("An unexpected panic occurred when getting http Scenario Exporter HealthCkeck...")
		}
	}()

	health := append(healthStructure, HealthDetail{
		Name: "HttpScenarioExporter", Health: "true",
	})

	//c.Logger.Info("HealthCheck demanded %+v",health)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(health)
	if err != nil {
		c.Logger.Error("Could not interpret health info : %s", err.Error())
	}
}
