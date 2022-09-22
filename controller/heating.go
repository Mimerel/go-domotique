package controller

import (
	"go-domotique/heating"
	"go-domotique/models"
	"html/template"
	"net/http"
)

func heatingController(config *models.Configuration) {
	http.HandleFunc("/heating/update", func(w http.ResponseWriter, r *http.Request) {
		err := heating.UpdateHeating(w, r, config)
		if err != nil {
			config.Logger.Error("heatingController Unable to update heating %+v ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(200)
		}
	})

	http.HandleFunc("/heating/status", func(w http.ResponseWriter, r *http.Request) {
		heating.StatusPage(w, r, config)
	})

	http.HandleFunc("/heating/updateValues", func(w http.ResponseWriter, r *http.Request) {
		heating.Update(w, r, config)
	})
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/heating/temporary/", func(w http.ResponseWriter, r *http.Request) {
		config.Logger.Debug("heatingController In temporary")

		err := heating.SettingTemporaryValues(config, r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			t := template.New("confirmation.html")
			t, err := t.ParseFiles("./heating/templates/confirmation.html")
			if err != nil {
				config.Logger.Error("heatingController Error Parsing template%+v", err)
			}
			err = t.Execute(w, models.HeatingConfirmation{
				IpPort: config.Ip + ":" + config.Port,
			})
			if err != nil {
				config.Logger.Error("heatingController Error Execution %+v", err)
			}
		}
	})
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./heating/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./heating/js"))))
}
