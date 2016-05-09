package main

import (
	"log"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func usuarioDePrueba() *Usuario {
	fechaNac, _ := time.Parse("2006-01-02", "1989-01-28")
	return &Usuario{
		"Victor",
		"Samuel",
		"Mosquera",
		"Artamonov",
		"1087998004",
		"cc",
		"m",
		fechaNac,
		"victorsamuelmd",
		"natanata",
		"medico",
	}
}

func TestGuardarUsuario(t *testing.T) {
	store := ds.Copy()
	defer store.session.Close()
	db := ds.session.DB(NombreBaseDatosTest)

	err := GuardarUsuario(usuarioDePrueba(), db)
	defer db.C(NombreCollecionUsuario).DropCollection()

	if err != nil {
		t.Errorf("Fallo con error %s", err.Error())
	}

	result := Usuario{}
	err = db.C(NombreCollecionUsuario).Find(bson.M{"primerNombre": "Victor"}).One(&result)
	if err != nil {
		t.Errorf("Fallo con error %s", err.Error())
	}

	if result.Identificacion != "1087998004" {
		t.Error("No se encontro el usuario con Identificacion 1087998004")
	}

}

func TestGuardarUsuarioRepetido(t *testing.T) {
	store := ds.Copy()
	defer store.session.Close()
	db := ds.session.DB(NombreBaseDatosTest)

	err := GuardarUsuario(usuarioDePrueba(), db)
	defer db.C(NombreCollecionUsuario).DropCollection()

	if err != nil {
		t.Errorf("Fallo con la conexion a la base de datos: %s", err.Error())
	}

	err = GuardarUsuario(usuarioDePrueba(), db)
	if err == nil {
		t.Errorf("Guardó un usuario con igual Identificacion: %s", err.Error())
	}
}

func TestCrearToken(t *testing.T) {
	tokenInvalido := `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJncnVwbyI6Im1lZGljb3MiLCJub21icmVfdXN1YXJpbyI6InZpY3RvcnNhbXVlbG1kIn0.Y32PB5Ij1iwcPwB7ER6NEPaDrQzM0_hS-osi6c2AXfDGXWWkSBk5ojYeKbq_udSso2u0Bhi4Xj7Fe7Gz1koZeMGq8N9T8isnyYMUNoxw5sv6hADoLYfzHj6U3FIVE_cvdb4xVZ9LFqWm7fvyjWQ_LbkVLrm5tH1PrLbWq5oceis`
	token, _ := crearToken("victorsamuelmd", "medicos")
	if !validarToken(token) {
		t.Error("Fallo al validar el token producido")
	}
	if validarToken(tokenInvalido) {
		t.Error("Fallo al invalidar el token inválido")
	}
}

func TestAutenticarUsuario(t *testing.T) {
	store := ds.Copy()
	defer store.session.Close()
	db := ds.session.DB(NombreBaseDatosTest)

	err := GuardarUsuario(usuarioDePrueba(), db)
	defer db.C(NombreCollecionUsuario).DropCollection()

	if err != nil {
		log.Fatal(err)
	}

	token, err := AutenticarUsuario("victorsamuelmd", "natanata", db)
	if len(token) == 0 || err != nil {
		t.Error("Falló autenticacion y el token no fue producido")
	}
}

func BenchmarkMgoConnectionConcurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		go insertarAlgoEnLaBaseDatos()
	}
	ds.session.DB(NombreBaseDatosTest).C("duration").DropCollection()
}

func BenchmarkMgoConnection(b *testing.B) {
	for i := 0; i < b.N; i++ {
		insertarAlgoEnLaBaseDatos()
	}
	ds.session.DB(NombreBaseDatosTest).C("duration").DropCollection()
}

func insertarAlgoEnLaBaseDatos() {
	store := ds.Copy()
	db := ds.session.DB(NombreBaseDatosTest)
	db.C("duration").Insert(bson.M{"name": "Some Name"})
	store.session.Close()
}
