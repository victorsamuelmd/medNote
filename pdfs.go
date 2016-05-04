package main

import (
	"encoding/json"
	"io"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// CrearPDFRemision acepta un string que contiene la informacion sobre la
// remision en formato json y devuelve un objeto gofdf.Fdf con la remisión en
// formato pdf.
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
