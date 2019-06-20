package googleAssistant

import (
	"go-domotique/devices"
	"go-domotique/google_talk"
	"go-domotique/logger"
	"go-domotique/models"
	"go-domotique/utils"
	"net/http"
	"strings"
)

/**
Method that searches for the ip(s) concerned by a room.
When an instruction is used, it will always be linked to a room
 */
func findIpOfGoogleHome(config *models.Configuration, concernedRoom string) ([]string) {
	ips := []string{}
	for _, google := range config.GoogleAssistant.GoogleBoxes {
		if google.Name == concernedRoom {
			ips = append(ips, google.Ip)
		}
	}
	return ips
}

/**
Method that splits the instruction into an action and a instruction
 */
func getActionAndInstruction(config *models.Configuration, instruction string) (action string, newInstruction string) {
	instruction = utils.ConvertInstruction(instruction)
	logger.Info(config, "getActionAndInstruction","instructions: <%s> ", instruction)
	mainAction := strings.Split(instruction, " ")[0]
	instruction = strings.Replace(instruction, mainAction, "", 1)
	instruction = strings.Trim(instruction, " ")
	return mainAction, instruction
}

/**
Method that checks if the action demanded exists and retrieves the information linked to this action.
 */
func checkActionValidity(config *models.Configuration, mainAction string) (string) {
	found := ""
	for _, action := range config.GoogleAssistant.GoogleTranslatedActionTypes {
		if action.ActionWord == mainAction {
			found = action.Action
			break
		}
	}
	return found
}

/**
Method that searches throw the list of possible commands for the
command sent by google home.
It first tries to find the corresponding "sentence" in its database.
IF it is found, it will check if the action is autorized in that room
If so, it will execute the command
 */
func RunDomoticCommand(config *models.Configuration, instruction string, concernedRoom string, mainAction string) (bool) {
	found := false
	for _, word := range config.GoogleAssistant.GoogleWords {
		if utils.CompareWords(config, word.Words, instruction) {
			for _, ListInstructions := range config.GoogleAssistant.GoogleTranslatedInstructions {
				if strings.ToUpper(ListInstructions.Room) == strings.ToUpper(concernedRoom) &&
					strings.ToUpper(ListInstructions.Type) == strings.ToUpper(mainAction) &&
					word.ActionNameId == ListInstructions.ActionNameId {
					logger.Info(config, "RunDomoticCommand", "Found instruction %v", ListInstructions)
					devices.ExecuteAction(config, ListInstructions)
					found = true
				}
			}
		}
	}
	return found
}

func AnalyseRequest(w http.ResponseWriter, r *http.Request, urlParams []string, config *models.Configuration) {
	concernedRoom := urlParams[1]
	ips := findIpOfGoogleHome(config, concernedRoom)
	if len(ips) == 0 {
		logger.Info(config, "AnalyseRequest", "No google home ips found")
		w.WriteHeader(500)
	} else {
		mainAction, instruction := getActionAndInstruction(config, urlParams[2])
		mainAction = checkActionValidity(config, mainAction)
		logger.Info(config, "AnalyseRequest", "Checked instructions: <%s> <%s>", mainAction, instruction)
		if mainAction == "" {
			logger.Error(config, "AnalyseRequest","not found action <%s>, room <%s>, command <%s>", mainAction, concernedRoom, instruction)
			google_talk.Talk(config, ips, "Action introuvable")
			w.WriteHeader(500)
		} else {
			found := false
			logger.Info(config, "AnalyseRequest", "Running action <%s>, room <%s>, command <%v>, level <%v>", mainAction, concernedRoom, instruction)
			found = RunDomoticCommand(config, instruction, concernedRoom, mainAction)
			if found {
				w.WriteHeader(200)
			} else {
				logger.Error(config, "AnalyseRequest","not found action <%s>, room <%s>, command <%s>", mainAction, concernedRoom, instruction)
				google_talk.Talk(config, ips, "Instruction introuvable")
				w.WriteHeader(500)
			}
		}
	}
}

