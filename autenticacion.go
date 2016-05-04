package main

import (
	"crypto/sha256"
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Remision struct {
	Paciente    Usuario   `json:"paciente"`
	Medico      Usuario   `json:"medico"`
	Receptor    string    `json:"receptor"`
	Fecha       time.Time `json:"fecha"`
	Servicio    string    `json:"servicio"`
	Contenido   string    `json:"contenido"`
	Diagnostico string    `json:"diagnostico"`
}

type Usuario struct {
	PrimerNombre    string `bson:"primer_nombre"`
	SegundoNombre   string `json:"segundo_nombre"`
	PrimerApellido  string `json:"primer_apellido"`
	SegundoApellido string `json:"segundo_apellido"`

	Identificacion string `json:"identificacion"`
	TipoId         string `json:"tipo_identificacion"`

	Genero          string    `json:"genero"`
	FechaNacimiento time.Time `json:"fecha_nacimiento"`

	NombreUsuario string `bson:"nombre_usuario" json:"nombre_usuario"`
	Contraseña    string `bson:"contrasena" json:"contrasena"`
	Grupo         string `json:"grupo"`
}

func GuardarUsuario(usr *Usuario, db *mgo.Database) error {
	usr.Contraseña = fmt.Sprintf("%x", sha256.Sum256([]byte(usr.Contraseña)))
	return db.C("usuario").Insert(usr)
}

func UsuarioAutentico(username, password string, db *mgo.Database) bool {
	usr := &Usuario{}
	pwd := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	err := db.C("usuario").Find(bson.M{"nombre_usuario": username,
		"contrasena": pwd}).One(usr)
	if err != nil {
		return false
	}
	return true
}

func (u *Usuario) NombreCompleto() string {
	return fmt.Sprintf("%s %s %s %s", u.PrimerNombre,
		u.SegundoNombre, u.PrimerApellido, u.SegundoApellido)
}

func (u *Usuario) Nombres() string {
	return fmt.Sprintf("%s %s", u.PrimerNombre, u.SegundoNombre)
}

func (u *Usuario) Edad(t time.Time) string {
	return fmt.Sprint(AgeAt(u.FechaNacimiento, t))
}

func AgeAt(birthDate time.Time, now time.Time) int {
	// Get the year number change since the player's birth.
	years := now.Year() - birthDate.Year()

	// If the date is before the date of birth, then not that many years
	// have elapsed.
	birthDay := getAdjustedBirthDay(birthDate, now)
	if now.YearDay() < birthDay {
		years -= 1
	}

	return years
}

// Age is shorthand for AgeAt(birthDate, time.Now()), and carries the same usage
// and limitations.
func Age(birthDate time.Time) int {
	return AgeAt(birthDate, time.Now())
}

// Gets the adjusted date of birth to work around leap year differences.
func getAdjustedBirthDay(birthDate time.Time, now time.Time) int {
	birthDay := birthDate.YearDay()
	currentDay := now.YearDay()
	if isLeap(birthDate) && !isLeap(now) && birthDay >= 60 {
		return birthDay - 1
	}
	if isLeap(now) && !isLeap(birthDate) && currentDay >= 60 {
		return birthDay + 1
	}
	return birthDay
}

// Works out if a time.Time is in a leap year.
func isLeap(date time.Time) bool {
	year := date.Year()
	if year%400 == 0 {
		return true
	} else if year%100 == 0 {
		return false
	} else if year%4 == 0 {
		return true
	}
	return false
}
