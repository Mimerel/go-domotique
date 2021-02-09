package googleAssistant

import (
	"go-domotique/devices"
	"go-domotique/google_talk"
	"go-domotique/logger"
	"go-domotique/models"
	"go-domotique/utils"
	"go-domotique/wifi"
	"net/http"
	"strconv"
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
func RunDomoticCommand(config *models.Configuration, instruction string, mainAction string) (bool) {
	found := false

	instruction = strings.ToLower(strings.Replace(instruction, " ", "", -1))
	for _,v := range config.CharsToReplace {
		instruction = strings.Replace(instruction, v.From, v.To, -1)
	}

	for _, word := range config.GoogleAssistant.GoogleWords {
		if utils.CompareWords(word.WordsConverted, instruction) {
			for _, ListInstructions := range config.GoogleAssistant.GoogleTranslatedInstructions {
				go runDomotiqueInstruction(config, mainAction, word, ListInstructions)
			}
			found = true
		}
	}
	return found
}

func runDomotiqueInstruction (config *models.Configuration , mainAction string, word models.GoogleWords, ListInstructions models.GoogleTranslatedInstruction) {
		if strings.ToUpper(ListInstructions.Type) == strings.ToUpper(mainAction) &&
			word.ActionNameId == ListInstructions.ActionNameId {
			logger.Info(config, "RunDomoticCommand", "Found instruction %+v", ListInstructions)
			device := devices.GetDomotiqueIdFromDeviceIdAndBoxId(config, ListInstructions.DeviceId, ListInstructions.ZwaveId)
			logger.Info(config, "RunDomoticCommand", "Found Device %+v", device)
			switch device.Zwave {
			case 100:
				logger.Info(config, "RunDomoticCommand", "Running Wifi instruction : %+v, %+v",ListInstructions.DeviceId , ListInstructions.Type )
				wifi.ExecuteRequestRelay(strconv.Itoa(int(ListInstructions.DeviceId)), ListInstructions.Type, config)
			default:
				logger.Info(config, "RunDomoticCommand", "Running Zwave instruction")
				devices.ExecuteAction(config, ListInstructions)
			}
		}
}


func AnalyseRequest(w http.ResponseWriter, r *http.Request, urlParams []string, config *models.Configuration) {
	ips := findIpOfGoogleHome(config, "Salon")
	if len(ips) == 0 {
		logger.Info(config, "AnalyseRequest", "No google home ips found")
		w.WriteHeader(500)
	} else {
		mainAction, instruction := getActionAndInstruction(config, urlParams[2])
		mainAction = checkActionValidity(config, mainAction)
		//logger.Info(config, "AnalyseRequest", "Checked instructions: <%s> <%s>", mainAction, instruction)
		if mainAction == "" {
			logger.Error(config, "AnalyseRequest","not found action <%s>, room <%s>, command <%s>", mainAction, instruction)
			google_talk.Talk(config, ips, "Action introuvable")
			w.WriteHeader(500)
		} else {
			found := false
			//logger.Info(config, "AnalyseRequest", "Running action <%s>, room <%s>, command <%v>, level <%v>", mainAction, concernedRoom, instruction)
			found = RunDomoticCommand(config, instruction, mainAction)
			if found {
				w.WriteHeader(200)
			} else {
				logger.Error(config, "AnalyseRequest","not found action <%s>, room <%s>, command <%s>", mainAction, instruction)
				google_talk.Talk(config, ips, "Instruction introuvable")
				w.WriteHeader(500)
			}
		}
	}
}

