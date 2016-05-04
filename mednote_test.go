package main

import (
	"log"
	"testing"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func TestGuardarUsuario(t *testing.T) {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	db := session.DB("test")

	err = GuardarUsuario(
		&Usuario{"Victor", "Samuel", "Mosquera", "Artamonov",
			"1087998004", "cc", "m", time.Now(),
			"victorsamuelmd", "natanata", "medico",
		}, db)
	defer db.C("usuario").DropCollection()

	if err != nil {
		log.Fatal(err)
	}

	result := Usuario{}
	err = db.C("usuario").Find(bson.M{"primer_nombre": "Victor"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	if result.Identificacion != "1087998004" {
		t.Fail()
	}

	if !UsuarioAutentico("victorsamuelmd", "natanata", db) {
		t.Fail()
	}

}
