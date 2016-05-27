// Package main provides ...
package main

import (
	"encoding/json"

	"gopkg.in/mgo.v2"

	"io/ioutil"
)

type DataStore struct {
	Session *mgo.Session
}

func NewDataStore() *DataStore {
	jsonConf, err := ioutil.ReadFile("config.json")

	// Si no hay archivo de configuración no debe continuar la aplicación
	if err != nil {
		panic(err)
	}

	var conf struct {
		Pwd string `json:"pwd"`
	}
	json.Unmarshal(jsonConf, &conf)

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{MgoHostStr},
		Database: NombreBaseDatos,
		Username: "golangApplication",
		Password: conf.Pwd,
	})

	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	return &DataStore{session}
}

func (ds *DataStore) Copy() *DataStore {
	return &DataStore{ds.Session.Copy()}
}
