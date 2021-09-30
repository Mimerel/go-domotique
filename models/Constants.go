package models

const (
	Prefix                    = "shellies/device_"
	ShellyPower               = "/relay/0/power"
	ShellyEnergy              = "/relay/0/energy"
	ShellyOnOff_0             = "/relay/0"
	ShellyOnOff_0_ison        = "/relay/0/ison"
	ShellyOnOff_1             = "/relay/1"
	ShellyOnOff_1_ison        = "/relay/1/ison"
	ShellyOnline              = "/online"
	ShellyTemperature0        = "/ext_temperature/0"
	ShellyStatus              = "/status"
	ShellyRollerState         = "/roller/0"
	ShellyCurrentPos          = "/roller/0/pos" //command/pos pour modifiier
	ShellyRollerLastDirection = "/roller/0/last_direction"
	ShellyRollerStopReason    = "/roller/0/stop_reason"
	ShellyRollerPower         = "/roller/0/power"
	ShellyRollerEnergy        = "/roller/0/energy"
	ShellyAnnounce            = "/announce"
)

