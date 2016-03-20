package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// CrearPDFRemision acepta un string que contiene la informacion sobre la
// remision en formato json y devuelve un objeto gofdf.Fdf con la remisión en
// formato pdf.
func CrearPDFRemision(rem io.Reader) *gofpdf.Fpdf {
	d := json.NewDecoder(rem)
	var r Remision
	err := d.Decode(&r)
	// TODO: Handle the error
	if err != nil {
		panic(err.Error())
	}

	pdf := gofpdf.New("P", "pt", "letter", "")
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 14)
	pdf.Text(64, 156, "Hospital San Juan de Dios")
	pdf.Text(327, 156, translator(r.Receptor))
	pdf.Text(464, 218, translator(r.Paciente.Identificacion))
	pdf.Text(46, 194, translator(r.Paciente.PrimerApellido))
	pdf.Text(180, 194, translator(r.Paciente.SegundoApellido))
	pdf.Text(320, 194, translator(r.Paciente.Nombres()))
	pdf.Text(43, 222, r.Paciente.Edad(time.Now()))
	if r.Paciente.Genero == "m" {
		pdf.Text(138, 228, "X")
	} else {
		pdf.Text(170, 228, "X")
	}
	pdf.Text(162, 328, r.Medico.NombreCompleto())
	pdf.Text(366, 328, translator(r.Servicio))
	tstring := r.Fecha.Format(time.RFC3339)
	pdf.SetFont("Helvetica", "", 10)
	pdf.Text(50, 322, tstring[8:10])
	pdf.Text(87, 322, tstring[5:7])
	pdf.Text(126, 322, tstring[2:4])
	pdf.SetLeftMargin(50)
	pdf.SetRightMargin(50)
	pdf.SetY(416)
	pdf.Write(20, translator(r.Contenido))
	return pdf
}

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
	PrimerNombre    string    `json:"primer-nombre"`
	SegundoNombre   string    `json:"segundo-nombre"`
	PrimerApellido  string    `json:"primer-apellido"`
	SegundoApellido string    `json:"segundo-apellido"`
	Identificacion  string    `json:"identificacion"`
	TipoId          string    `json:"tipo-identificacion"`
	Genero          string    `json:"genero"`
	FechaNacimiento time.Time `json:"fecha-nacimiento"`
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

	// If the date is before the date of birth, then not that many years have elapsed.
	birthDay := getAdjustedBirthDay(birthDate, now)
	if now.YearDay() < birthDay {
		years -= 1
	}

	return years
}

// Age is shorthand for AgeAt(birthDate, time.Now()), and carries the same usage and limitations.
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
