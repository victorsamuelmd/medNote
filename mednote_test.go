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
		t.Errorf("Fallon con error %s", err.Error())
	}

	result := Usuario{}
	err = db.C("usuario").Find(bson.M{"primerNombre": "Victor"}).One(&result)
	if err != nil {
		t.Errorf("Fallo con error %s", err.Error())
	}

	if result.Identificacion != "1087998004" {
		t.Error("No se encontro el usuario con Identificacion 1087998004")
	}

}

func TestGuardarUsuarioRepetido(t *testing.T) {
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

	err = GuardarUsuario(
		&Usuario{"Victor", "Samuel", "Mosquera", "Artamonov",
			"1087998004", "cc", "m", time.Now(),
			"victorsamuelmd", "natanata", "medico",
		}, db)
	if err == nil {
		t.Fail()
	}
}

func TestCrearToken(t *testing.T) {
	tokenInvalido := `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJncnVwbyI6Im1lZGljb3MiLCJub21icmVfdXN1YXJpbyI6InZpY3RvcnNhbXVlbG1kIn0.Y32PB5Ij1iwcPwB7ER6NEPaDrQzM0_hS-osi6c2AXfDGXWWkSBk5ojYeKbq_udSso2u0Bhi4Xj7Fe7Gz1koZeMGq8N9T8isnyYMUNoxw5sv6hADoLYfzHj6U3FIVE_cvdb4xVZ9LFqWm7fvyjWQ_LbkVLrm5tH1PrLbWq5oceis`
	token, _ := crearToken("victorsamuelmd", "medicos")
	if !validarToken(token) {
		t.Fail()
	}
	if validarToken(tokenInvalido) {
		t.Fail()
	}
}

func TestAutenticarUsuario(t *testing.T) {
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

	token, err := AutenticarUsuario("victorsamuelmd", "natanata", db)
	if len(token) == 0 || err != nil {
		t.Fail()
	}
}
