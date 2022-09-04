package models

const (
	Prefix                      = "shellies/device_"
	ShellyInfo                  = "/info"
	ShellySettings              = "/settings"
	ShellyPower                 = "/relay/0/power"
	ShellyPower2                = "/relay/power"
	ShellyEnergy                = "/relay/0/energy"
	ShellyEnergy2               = "/relay/energy"
	ShellyOnOff_0               = "/relay/0"
	ShellyOnOff_0_ison          = "/relay/0/ison"
	ShellyOnOff_1               = "/relay/1"
	ShellyOnOff_1_ison          = "/relay/1/ison"
	ShellyOnline                = "/online"
	ShellyTemperature0          = "/ext_temperature/0"
	ShellyTemperature0F         = "/ext_temperature_f/0"
	ShellyTemperatures          = "/ext_temperatures"
	ShellyTemperaturesF         = "/ext_temperatures_f"
	ShellyStatus                = "/status"
	ShellyStatusSwitch0         = "/status/switch:0"
	ShellyStatusSwitch1         = "/status/switch:1"
	ShellyStatusSwitch2         = "/status/switch:2"
	ShellyStatusSwitch3         = "/status/switch:3"
	ShellyRollerState           = "/roller/0"
	ShellyCurrentPos            = "/roller/0/pos"
	ShellyRollerLastDirection   = "/roller/0/last_direction"
	ShellyRollerStopReason      = "/roller/0/stop_reason"
	ShellyRollerPower           = "/roller/0/power"
	ShellyRollerEnergy          = "/roller/0/energy"
	ShellyTemperatureStatus     = "/temperature_status"
	ShellyTemperatureDevice     = "/temperature"
	ShellyTemperatureFDevice    = "/temperature_f"
	ShellyInputO                = "/input/0"
	ShellyInput1                = "/input/1"
	ShellyInput2                = "/input/2"
	ShellyInput3                = "/input/3"
	ShellyOverTemperatureDevice = "/overtemperature"
	ShellyVoltage               = "/voltage"
	ShellyAnnounce              = "/announce"
	ShellyReasons               = "/sensor/act_reasons"
	ShellySensorBattery         = "/sensor/battery"
	ShellyFlood                 = "/sensor/flood"
)
