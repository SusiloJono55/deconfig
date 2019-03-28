package repository

import (
	b64 "encoding/base64"
	"fmt"
	"log"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"github.com/timshannon/badgerhold"
)

type DBType int

const (
	MSSql DBType = iota //MSSql Type
	Oracle
	Postgresql
	MySQL
	DB2
)

var TypeStr = map[DBType]string{
	MSSql:      "MSSql",
	Oracle:     "Oracle",
	Postgresql: "Postgresql",
	MySQL:      "MySQL",
	DB2:        "DB2",
}

var MapType = map[string]DBType{
	"MSSql":      MSSql,
	"Oracle":     Oracle,
	"Postgresql": Postgresql,
	"MySQL":      MySQL,
	"DB2":        DB2,
}

var (
	store      *badgerhold.Store
	configInfo map[string]interface{}
)

type Info interface {
	Set(types DBType) interface{}
	Get(types DBType) interface{}
}

type ReqBody struct {
	CType    string `json:"c_type,omitempty"`
	CName    string `json:"c_name,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Database string `json:"database,omitempty"`
	SID      string `json:"sid,omitempty"`
	Instance string `json:"instance,omitempty"`
}

type Config struct {
	Key       string                 `badgerhold:"key" json:"key,omitempty"`
	Type      DBType                 `badgerhold:"index"`
	Name      string                 `json:"name,omitempty"`
	Owner     string                 `badgerhold:"index" json:"owner,omitempty"` //groupName : "" as administrator
	Value     map[string]interface{} `json:"value,omitempty"`
	CreatedAt time.Time              `json:"created_at,omitempty"`
	UpdatedAt time.Time              `json:"updated_at,omitempty"`
}

func OpenStore() error {
	var err error
	options := badgerhold.DefaultOptions
	options.Dir = "./badgerhold-conf"
	options.ValueDir = "./badgerhold-conf"
	options.ReadOnly = false
	options.Truncate = true

	store, err = badgerhold.Open(options)
	return err
}

func CloseStore() {
	store.Close()
}

func LoadFileConfiguration() {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	configInfo = viper.GetStringMap("info")
}

func NewConfiguration(_name, _owner string, _type DBType, _value ReqBody) *Config {
	conf := &Config{}
	id, _ := uuid.NewV4()
	conf.Key = id.String()
	conf.Name = _name
	conf.Type = _type
	conf.Owner = _owner
	conf.CreatedAt = time.Now()

	conf.MapConfig(_value)
	return conf
}

func (conf *Config) MapConfig(_value ReqBody) {
	value := map[string]interface{}{}
	reqValue := _value
	var password string
	oldPass, oldOk := conf.Value["password"]
	if oldOk && oldPass == reqValue.Password {
		password = reqValue.Password
	} else {
		password = b64.StdEncoding.EncodeToString([]byte(reqValue.Password))
	}
	value["username"] = reqValue.Username
	value["password"] = password
	value["host"] = reqValue.Host
	value["port"] = reqValue.Port
	switch conf.Type {
	case MSSql:
		value["database"] = reqValue.Database
		value["instance"] = reqValue.Instance
	case Oracle:
		value["sid"] = reqValue.SID
	case Postgresql, MySQL, DB2:
		value["database"] = reqValue.Database
	}
	conf.Value = value
	conf.UpdatedAt = time.Now()
	return
}

func GetInfoAll() interface{} {
	return viper.GetStringMap("info")
}

func GetInfo(_type DBType) interface{} {
	return viper.GetStringMap("info." + TypeStr[_type])
}

func (conf *Config) Create() error {
	err := store.Insert(conf.Key, conf)

	return err
}

func (conf *Config) GetOne(key string, owner string) error {
	var datas []Config
	err := store.Find(&datas, badgerhold.Where("Key").Eq(key).And("Owner").Eq(owner).Index("Owner"))
	if err != nil {
		return err
	}
	if len(datas) != 0 {
		data := datas[0]
		*conf = data
		return nil
	}
	return fmt.Errorf("No data found")
}

func (conf *Config) GetByName(name string, owner string) error {
	var datas []Config
	err := store.Find(&datas, badgerhold.Where("Name").Eq(name).And("Owner").Eq(owner).Index("Owner"))
	if err != nil {
		return err
	}
	if len(datas) != 0 {
		data := datas[0]
		*conf = data
		return nil
	}
	return fmt.Errorf("No data found")
}

func (conf *Config) Update() error {
	return store.Update(conf.Key, conf)
}

func (conf *Config) Delete() error {
	return store.Delete(conf.Key, conf)
}

func GetAll(owner string) ([]Config, error) {
	var result []Config

	err := store.Find(&result, badgerhold.Where("Owner").Eq(owner).Index("Owner"))
	if err != nil {
		fmt.Printf("Query on alternate tag index failed: %s", err)
		return nil, err
	}

	// get all config data
	return result, nil
}

func GetAllByType(dbType DBType, owner string) ([]Config, error) {
	var result []Config

	err := store.Find(&result, badgerhold.Where("Type").Eq(dbType).Index("Type").And("Owner").Eq(owner).Index("Owner"))
	if err != nil {
		fmt.Printf("Query on alternate tag index failed: %s", err)
		return nil, err
	}

	// get all config data
	return result, nil
}
