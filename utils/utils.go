package utils

import (
	"go-domotique/models"
	"strings"
	"time"
)

/**
On the instruction string received by google home it :
Removes the escape characters possibility sent by google Home
Remove beginning and end spaces
And creates a new string to rewrite properly the string sent
 */
func ConvertInstruction (value string) string {
	instruction := strings.Replace(value, "<<", "", 1)
	instruction = strings.Replace(instruction, ">>", "", 1)
	instruction = strings.Trim(instruction, " ")
	newValue := ""
	for i := 0; i < len(instruction); i++ {
		newValue = newValue + string(instruction[i])
	}
	return newValue
}

/**
Method that removes "s" and "x" from command to manage pluirals that would not
have been indicated in the list of commands in the configuration file
 */
func RemoveEndletter( value string) string {
	if strings.HasSuffix(value, "s") {
		value = strings.TrimSuffix(value, "s")
	} else if strings.HasSuffix(value, "x"){
		value = strings.TrimSuffix(value, "x")
	}
	return value
}

/**
Method that compares the string sent to a given string
by removing spaces and converting to lower case.
Google home puts spaces before and after '
This method solves that problem
 */
func CompareWords( word string, instruction string ) (bool) {
	if word != instruction {
		return false
	}
	return true
}

/**
Method that checks if the command is applicable for the google Home that
has sent the command (A google home is in a room)
 */
func CompareRooms(rooms []string, requestedFrom string) (bool) {
	same := false;
	for _, value := range rooms {
		if value == requestedFrom {
			same = true
		}
	}
	return same
}

/**
Checks if the command that is will be executed is autorised in the list
of authorized commands for that instruction
 */
func CompareActions(actions []string, RequestedAction string) (bool) {
	same := false;
	for _, value := range actions {
		if value == RequestedAction {
			same = true
		}
	}
	return same
}

/**
Generic method that could replace above but makes code
less readable.
 */
func IsInArray(list []string, value string) (bool) {
	exists := false;
	for _, v := range list {
		if v == value {
			exists = true
		}
	}
	return exists
}

func GetTimeAndDay(config *models.Configuration) {
	// getting time
	config.Heating.HeatingMoment.Moment = time.Now().In(config.Location)
	hour := config.Heating.HeatingMoment.Moment.Hour() * 100

	config.Heating.HeatingMoment.Time = hour + config.Heating.HeatingMoment.Moment.Minute()
	// getting weekday
	config.Heating.HeatingMoment.Weekday = config.Heating.HeatingMoment.Moment.Weekday()
	// creatingDate
	config.Heating.HeatingMoment.Date = CreateDate(config.Heating.HeatingMoment.Moment)
}

func CreateDate(moment time.Time) int {
	return moment.Year() * 10000 + int(moment.Month()) * 100 + moment.Day()
}
