package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/bearbin/go-age"
	"github.com/jung-kurt/gofpdf"
)

// CrearPDFRemision crea la remisión en formato pdf
func CrearPDFRemision(rem io.Reader) (*gofpdf.Fpdf, error) {
	d := json.NewDecoder(rem)
	var r Remision
	err := d.Decode(&r)

	pdf := gofpdf.New("P", "pt", "letter", "")
	if err != nil {
		return pdf, err
	}

	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 14)
	pdf.Text(64, 156, "Hospital San Juan de Dios")
	pdf.Text(327, 156, utf8toIso(r.Receptor))
	pdf.Text(464, 218, utf8toIso(r.Paciente.Identificacion))
	pdf.Text(46, 194, utf8toIso(r.Paciente.PrimerApellido))
	pdf.Text(180, 194, utf8toIso(r.Paciente.SegundoApellido))
	pdf.Text(320, 194, utf8toIso(r.Paciente.Nombres()))
	pdf.Text(43, 222, r.Paciente.Edad(time.Now()))
	if r.Paciente.Genero == "m" {
		pdf.Text(138, 228, "X")
	} else {
		pdf.Text(170, 228, "X")
	}
	pdf.Text(162, 328, r.Medico.NombreCompleto())
	pdf.Text(366, 328, utf8toIso(r.Servicio))
	tstring := r.Fecha.Format(time.RFC3339)
	pdf.SetFont("Helvetica", "", 10)
	pdf.Text(50, 322, tstring[8:10])
	pdf.Text(87, 322, tstring[5:7])
	pdf.Text(126, 322, tstring[2:4])
	pdf.SetLeftMargin(50)
	pdf.SetRightMargin(50)
	pdf.SetY(416)
	pdf.Write(20, utf8toIso(r.Contenido))

	return pdf, nil
}

// Remision contiene la estructura de la información necesaria en caso de una
// remisión
type Remision struct {
	Paciente         Usuario   `json:"paciente" bson:"paciente"`
	Medico           Usuario   `json:"medico" bson:"medico"`
	Receptor         string    `json:"receptor" bson:"receptor"`
	Fecha            time.Time `json:"fecha" bson:"fecha"`
	Servicio         string    `json:"servicio" bson:"servicio"`
	Contenido        string    `json:"contenido" bson:"contenido"`
	Diagnostico      string    `json:"diagnostico" bson:"diagnostico"`
	TelefonoPaciente string    `json:"telefonoPaciente" bson:"telefonoPaciente"`
}

// Usuario contiente la estructura de le información sobre un usuario tanto del
// sistema como de un paciente
type Usuario struct {
	PrimerNombre    string `json:"primerNombre" bson:"primerNombre"`
	SegundoNombre   string `json:"segundoNombre" bson:"segundoNombre"`
	PrimerApellido  string `json:"primerApellido" bson:"primerApellido"`
	SegundoApellido string `json:"segundoApellido" bson:"segundoApellido"`

	Identificacion string `json:"identificacion" bson:"identificacion"`
	TipoID         string `json:"tipoIdentificacion" bson:"tipoIdentificacion"`

	Genero          string    `json:"genero" bson:"genero"`
	FechaNacimiento time.Time `json:"fechaNacimiento" bson:"fechaNacimiento"`

	NombreUsuario string `bson:"nombreUsuario" json:"nombreUsuario"`
	Contraseña    string `bson:"contrasena" json:"contrasena,omitempty"`
	Grupo         string `json:"grupo" bson:"grupo"`
}

// NombreCompleto devuelve el nombre completo del usuario siendo la union de
// ambos nombres seguido de ambos apellidos
func (u *Usuario) NombreCompleto() string {
	return fmt.Sprintf("%s %s %s %s", u.PrimerNombre,
		u.SegundoNombre, u.PrimerApellido, u.SegundoApellido)
}

// Nombres devuelve ambos nombres el usuario en un solo texto
func (u *Usuario) Nombres() string {
	return fmt.Sprintf("%s %s", u.PrimerNombre, u.SegundoNombre)
}

// Edad devuelve la edad del usuario en años cumplidos
func (u *Usuario) Edad(t time.Time) string {
	return fmt.Sprint(age.AgeAt(u.FechaNacimiento, t))
}
