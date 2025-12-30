package models

type Zwave struct {
	Id   int64  `csv:"id"`
	Name string `csv:"name"`
	Ip   string `csv:"ip"`
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
