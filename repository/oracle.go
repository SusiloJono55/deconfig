package repository

type OracleConn struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Sid      string `json:"sid,omitempty"`
}

func (info OracleConn) Set(types DBType) interface{} {
	return info
}

func (info OracleConn) Get(types DBType) interface{} {
	return configInfo[TypeStr[types]]
}
