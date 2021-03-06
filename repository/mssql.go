package repository

type MSSqlConn struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Database string `json:"database,omitempty"`
	Instance string `json:"instance,omitempty"`
}

func (info MSSqlConn) Set(types DBType) interface{} {
	return info
}

func (info MSSqlConn) Get(types DBType) interface{} {
	return configInfo[TypeStr[types]]
}
