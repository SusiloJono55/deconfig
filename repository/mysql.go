package repository

type MySQLConn struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Database string `json:"database,omitempty"`
}

func (info MySQLConn) Set(types DBType) interface{} {
	return info
}

func (info MySQLConn) Get(types DBType) interface{} {
	return configInfo[TypeStr[types]]
}
