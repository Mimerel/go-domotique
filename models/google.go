package models

type Elasticsearch struct {
	Url string `yaml:"url,omitempty"`
}

type GoogleWords struct {
	Words          string `csv:"words"`
	Id             int64  `csv:"id"`
	ActionNameId   int64  `csv:"actionNameId"`
	WordsConverted string
}

type GoogleActionNames struct {
	Id   int64  `csv:"id"`
	Name string `csv:"name"`
}

type GoogleInstruction struct {
	Id           int64 `csv:"id"`
	ActionNameId int64 `csv:"actionNameId"`
	TypeId       int64 `csv:"typeId"`
	DomotiqueId  int64 `csv:"domotiqueId"`
	GoogleBoxId  int64 `csv:"googleBoxId"`
	Value        int64 `csv:"value"`
}

type GoogleTranslatedInstruction struct {
	Id           int64  `csv:"id"`
	ActionName   string `csv:"actionName"`
	ActionNameId int64  `csv:"actionNameId"`
	Type         string `csv:"type"`
	DeviceName   string `csv:"deviceName"`
	DeviceId     int64  `csv:"deviceId"`
	ZwaveId      int64  `csv:"zwaveId"`
	ZwaveUrl     string `csv:"zwaveUrl"`
	GoogleBox    string `csv:"zwaveName"`
	Room         string `csv:"room"`
	TypeDevice   string `csv:"typeDevice"`
	Instance     int64  `csv:"instance"`
	CommandClass int64  `csv:"commandClass"`
	Value        int64  `csv:"value"`
}

type Zwave struct {
	Id   int64  `csv:"id"`
	Name string `csv:"name"`
	Ip   string `csv:"ip"`
}

type GoogleActionTypes struct {
	Id   int64  `csv:"id"`
	Name string `csv:"name"`
}

type GoogleActionTypesWords struct {
	Id           int64  `csv:"id"`
	ActionTypeId int64  `csv:"remplacementId"`
	Action       string `csv:"name"`
}

type GoogleTranslatedActionTypes struct {
	ActionWord string
	Action     string
}

type Room struct {
	Id    int64  `csv:"id"`
	Name  string `csv:"room"`
	Floor string `csv:"floor"`
}

type DeviceType struct {
	Id   int64  `csv:"id"`
	Name string `csv:"name"`
}

type GoogleBox struct {
	Id   int64  `csv:"id"`
	Name string `csv:"name"`
	Ip   string `csv:"ip"`
}

type MariaDB struct {
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
	IP       string `yaml:"ip,omitempty"`
	Port     string `yaml:"port,omitempty"`
	Database string `yaml:"database,omitempty"`
	Logs     string `yaml:"logs,omitempty"`
	Stats    string `yaml:"stats,omitempty"`
	DB       Maria
}

type ConfigurationGoogleAssistant struct {
	GoogleWords                  []GoogleWords
	GoogleBoxes                  []GoogleBox
	GoogleInstructions           []GoogleInstruction
	GoogleActionNames            []GoogleActionNames
	GoogleTranslatedInstructions []GoogleTranslatedInstruction
	GoogleActionTypesWords       []GoogleActionTypesWords
	GoogleActionTypes            []GoogleActionTypes
	GoogleTranslatedActionTypes  []GoogleTranslatedActionTypes
}
