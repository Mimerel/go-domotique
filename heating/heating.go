package heating

import (
	"go-domotique/models"
	"go-domotique/logger"
	"html/template"
	"net/http"
)


func StatusPage(w http.ResponseWriter, r *http.Request, config *models.Configuration) {
	t := template.New("status.html")
	t, err := t.ParseFiles("./heating/templates/status.html")
	if err != nil {
		logger.Error(config, "StatusPage", "Error Parsing template%+v", err)
	}
	data, err := HeatingStatus(config)
	err = t.Execute(w, data)
	if err != nil {
		logger.Error(config, "StatusPage","Error Execution %+v", err)
	}
}
