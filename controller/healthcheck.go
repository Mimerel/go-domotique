package controller

import (
	"go-domotique/heating"
	"go-domotique/models"
	"go-domotique/logger"
	"net/http"
	"html/template"
)

func heatingController(config *models.Configuration) {
	http.HandleFunc("/heating/update", func(w http.ResponseWriter, r *http.Request) {
		err := heating.UpdateHeating(w, r, config)
		if err != nil {
			logger.Error(config, "heatingController", "Unable to update heating %+v ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(200)
		}
	})

	http.HandleFunc("/heating/status", func(w http.ResponseWriter, r *http.Request) {
		heating.StatusPage(w, r, config)
	})

	http.HandleFunc("/heating/temporary/", func(w http.ResponseWriter, r *http.Request) {
		err := heating.SettingTemporaryValues(config, r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			t := template.New("confirmation.html")
			t, err := t.ParseFiles( "./heating/templates/confirmation.html")
			if err != nil {
				config.Logger.Error("Error Parsing template%+v", err)
			}
			err = t.Execute(w, models.HeatingConfirmation{
				IpPort: config.Ip + ":" + config.Port,
			} )
			if err != nil {
				config.Logger.Error("Error Execution %+v", err)
			}
		}
	})

}