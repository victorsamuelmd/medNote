// Package main provides ...
package main

import (
	"encoding/json"

	"gopkg.in/mgo.v2"

	"io/ioutil"
)

type datastore struct {
	session *mgo.Session
}

func NewDataStore() *datastore {
	jsonConf, _ := ioutil.ReadFile("config.json")
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
	return &datastore{session}
}

func (ds *datastore) Copy() *datastore {
	return &datastore{ds.session.Copy()}
}
