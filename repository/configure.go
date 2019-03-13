package repository

import (
	"fmt"
	"log"

	"github.com/dgraph-io/badger"
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

type info interface {
	Set(types DBType) interface{}
	Get(types DBType) interface{}
}

type Config struct {
	Index string
	Key   string
	Type  string
	Id    string
	Value interface{}
}

func NewConfiguration(_name, _category, _version string, _type DBType) *Config {
	conf := &Config{}
	return conf
}

func (conf *Config) Create() error {
	options := badgerhold.DefaultOptions
	options.Dir = "./badgerhold-conf"
	options.ValueDir = "./badgerhold-conf"
	// store := conf.Db

	store, err := badgerhold.Open(options)
	defer store.Close()

	if err != nil {
		// handle error
		log.Fatal(err)
	}

	// insert the data in one transaction

	err = store.Badger().Update(func(tx *badger.Txn) error {

		err := store.TxInsert(tx, conf.Key, conf)
		if err != nil {
			return err

		}
		return nil
	})

	return nil
}

func (conf *Config) Update() *Config {
	return conf
}

func (conf *Config) Delete() *Config {
	return conf
}

func GetAll(category string) ([]Config, error) {
	options := badgerhold.DefaultOptions
	options.Dir = "./repository/badgerhold-conf"
	options.ValueDir = "./repository/badgerhold-conf"
	// store := conf.Db

	store, err := badgerhold.Open(options)
	defer store.Close()

	if err != nil {
		// handle error
		log.Fatal(err)
		return nil, err
	}

	var result []Config

	err = store.Find(&result, badgerhold.Where("Category").Eq(category).Index("Category"))
	if err != nil {
		fmt.Printf("Query on alternate tag index failed: %s", err)
		return nil, err
	}

	// get all config data
	return result, nil
}
