package repository

type DB2Conn struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Database string `json:"database,omitempty"`
}

func (info DB2Conn) Set(types DBType) interface{} {
	return info
}

func (info DB2Conn) Get(types DBType) interface{} {
	return configInfo[TypeStr[types]]
}
