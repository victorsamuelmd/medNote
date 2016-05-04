package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var secretKey, _ = rsa.GenerateKey(rand.Reader, 1024)

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
	c := db.C("usuario")
	if count, _ := c.Find(
		bson.M{"nombre_usuario": usr.NombreUsuario}).
		Count(); count > 0 {
		return errors.New("El usuario ya existe")
	}
	if count, _ := c.Find(
		bson.M{"identificacion": usr.Identificacion}).
		Count(); count > 0 {
		return errors.New("Número de identificacion ya fue utilizado")
	}
	usr.Contraseña = crearHashSHA256(usr.Contraseña)
	return c.Insert(usr)
}

//TODO: Falta agregar la funcion para determinar grupos o roles.
func UsuarioAutentico(username, password string, db *mgo.Database) bool {
	count, err := db.C("usuario").Find(bson.M{"nombre_usuario": username,
		"contrasena": crearHashSHA256(password)}).Count()

	if err != nil || count < 1 {
		return false
	}
	return true
}

func AutenticarUsuario(usr, pwd string, db *mgo.Database) (string, error) {
	if UsuarioAutentico(usr, pwd, db) {
		return crearToken(usr, "")
	} else {
		return "", errors.New("Nombre de usuario o contraseña invalida")
	}
}

func crearHashSHA256(pwd string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(pwd)))
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

func crearToken(username, authLevel string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims["nombre_usuario"] = username
	token.Claims["grupo"] = authLevel

	return token.SignedString(secretKey)
}

func validarToken(token string) bool {
	if token, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return secretKey.Public(), nil
	}); err == nil && token.Valid {
		return true
	}
	return false
}
