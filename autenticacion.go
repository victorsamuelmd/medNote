package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/bearbin/go-age"
	"github.com/dgrijalva/jwt-go"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var secretKey, _ = rsa.GenerateKey(rand.Reader, 1024)

type Remision struct {
	Paciente         Usuario   `json:"paciente" bson:"paciente"`
	Medico           Usuario   `json:"medico" bson:"medico"`
	Receptor         string    `json:"receptor bson:"receptor""`
	Fecha            time.Time `json:"fecha" bson:"fecha"`
	Servicio         string    `json:"servicio" bson:"servicio"`
	Contenido        string    `json:"contenido" bson:"contenido"`
	Diagnostico      string    `json:"diagnostico" bson:"diagnostico"`
	TelefonoPaciente string    `json:"telefonoPaciente" bson:"telefonoPaciente"`
}

type Usuario struct {
	PrimerNombre    string `json:"primerNombre" bson:"primerNombre"`
	SegundoNombre   string `json:"segundoNombre" bson:"segundoNombre"`
	PrimerApellido  string `json:"primerApellido" bson:"primerApellido"`
	SegundoApellido string `json:"segundoApellido" bson:"segundoApellido"`

	Identificacion string `json:"identificacion" bson:"identificacion"`
	TipoId         string `json:"tipoIdentificacion" bson:"tipoIdentificacion"`

	Genero          string    `json:"genero" bson:"genero"`
	FechaNacimiento time.Time `json:"fechaNacimiento" bson:"fechaNacimiento"`

	NombreUsuario string `bson:"nombreUsuario" json:"nombreUsuario"`
	Contraseña    string `bson:"contrasena" json:"contrasena"`
	Grupo         string `json:"grupo" bson:"grupo"`
}

func (u *Usuario) NombreCompleto() string {
	return fmt.Sprintf("%s %s %s %s", u.PrimerNombre,
		u.SegundoNombre, u.PrimerApellido, u.SegundoApellido)
}

func (u *Usuario) Nombres() string {
	return fmt.Sprintf("%s %s", u.PrimerNombre, u.SegundoNombre)
}

func (u *Usuario) Edad(t time.Time) string {
	return fmt.Sprint(age.AgeAt(u.FechaNacimiento, t))
}

func (u Usuario) EntityNamer() string {
	return "usuario"
}

func (u Usuario) GetID() string {
	return crearHashSHA256(fmt.Sprintf("%s%s%s", u.PrimerNombre, u.PrimerApellido, u.Identificacion))
}

func (u Usuario) SetID(string) error {
	return nil
}

func GuardarUsuario(usr *Usuario, db *mgo.Database) error {
	c := db.C("usuario")
	if count, _ := c.Find(
		bson.M{"nombreUsuario": usr.NombreUsuario}).
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
func usuarioAutentico(username, password string, db *mgo.Database) bool {
	count, err := db.C("usuario").Find(bson.M{"nombreUsuario": username,
		"contrasena": crearHashSHA256(password)}).Count()

	if err != nil || count < 1 {
		return false
	}
	return true
}

func AutenticarUsuario(usr, pwd string, db *mgo.Database) (string, error) {
	if usuarioAutentico(usr, pwd, db) {
		return crearToken(usr, "")
	} else {
		return "", errors.New("Nombre de usuario o contraseña invalida")
	}
}

func crearHashSHA256(pwd string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(pwd)))
}

func crearToken(username, authLevel string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims["nombreUsuario"] = username
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
