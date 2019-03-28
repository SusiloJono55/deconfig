package repository

type PostgresqlConn struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Database string `json:"database,omitempty"`
}

func (info PostgresqlConn) Set(types DBType) interface{} {
	return info
}

func (info PostgresqlConn) Get(types DBType) interface{} {
	return configInfo[TypeStr[types]]
}
