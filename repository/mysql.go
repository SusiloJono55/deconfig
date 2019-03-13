package repository

type MySQLInfo struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Port     int    `json:"port,omitempty"`
	Schema   string `json:"schema,omitempty"`
}
