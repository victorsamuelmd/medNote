package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/victorsamuelmd/mednote/general"
)

var translator, _ = gofpdf.UnicodeTranslatorFromFile("iso-8859-1.map")

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/pdf", ageneralPDF)
	mux.HandleFunc("/json", consultaJson)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("./dist"))))
	mux.HandleFunc("/remision", remisionPDF)
	mux.HandleFunc("/urgencia", urgenciaPDF)
	fmt.Println("Listening on localhost:8000, Hola mari")
	if err := http.ListenAndServe(":8000", mux); err != nil {
		fmt.Print(err.Error())
	} else {
		fmt.Println("Listening on localhost:8000, Hola mari")
	}
}

func consultaJson(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	hc := &general.ConsultaGeneral{}
	err := dec.Decode(hc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprint(w, hc)

}

func remisionPDF(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	pdf := gofpdf.New("P", "pt", "letter", "")
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 14)
	pdf.Text(64, 156, "Hospital San Juan de Dios")
	pdf.Text(327, 156, translator(r.FormValue("eps")))
	pdf.Text(464, 218, translator(r.FormValue("cedula")))
	pdf.Text(46, 194, translator(r.FormValue("papellido")))
	pdf.Text(180, 194, translator(r.FormValue("sapellido")))
	pdf.Text(320, 194, translator(fmt.Sprintf("%s %s", r.FormValue("pnombre"), r.FormValue("snombre"))))
	pdf.Text(43, 222, r.FormValue("edad"))
	if r.FormValue("genero") == "m" {
		pdf.Text(138, 228, "X")
	} else {
		pdf.Text(170, 228, "X")
	}
	pdf.Text(162, 328, "Victor Samuel Mosquera A.")
	pdf.Text(366, 328, translator(r.FormValue("servicio")))
	t := time.Now()
	tstring := t.Format(time.RFC3339)
	pdf.SetFont("Helvetica", "", 10)
	pdf.Text(50, 322, tstring[8:10])
	pdf.Text(87, 322, tstring[5:7])
	pdf.Text(126, 322, tstring[2:4])
	pdf.SetLeftMargin(50)
	pdf.SetRightMargin(50)
	pdf.SetY(416)
	pdf.Write(20, translator(r.FormValue("contenido")))

	if err := pdf.Output(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ageneralPDF(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	pdf := gofpdf.New("P", "pt", "letter", "")
	pdf.AddPage()
	pdf.Image("frenteFormatoHistoria.png", 0, 0, 612, 792, false, "png", 0, "")
	pdf.SetFont("Helvetica", "", 14)
	pdf.Text(440, 195, translator(r.FormValue("cedula")))
	pdf.Text(40, 185, translator(r.FormValue("papellido")))
	pdf.Text(180, 185, translator(r.FormValue("sapellido")))
	pdf.Text(300, 185, translator(fmt.Sprintf("%s %s", r.FormValue("pnombre"), r.FormValue("snombre"))))
	pdf.Text(40, 215, r.FormValue("edad"))
	if r.FormValue("genero") == "m" {
		pdf.Text(145, 225, "X")
	} else {
		pdf.Text(180, 225, "X")
	}
	t := time.Now()
	tstring := t.Format(time.RFC3339)
	pdf.SetTitle(fmt.Sprintf("%s %s %s %s %s", r.FormValue("pnombre"), r.FormValue("snombre"), r.FormValue("papellido"), r.FormValue("sapellido"), tstring), true)

	pdf.SetFont("Helvetica", "", 10)
	pdf.Text(37, 378, tstring[8:10])
	pdf.Text(62, 378, tstring[5:7])
	pdf.Text(85, 378, tstring[2:4])
	pdf.Text(109, 378, fmt.Sprintf("%v:%v", t.Hour(), t.Minute()))
	pdf.SetLeftMargin(140)
	pdf.SetRightMargin(52)
	pdf.SetTopMargin(80)
	pdf.SetY(365)

	WriteItem("Motivo de Consulta: ", r.FormValue("mconsulta"), pdf)
	WriteItem("Enfermedad Actual: ", r.FormValue("eactual"), pdf)
	WriteItem("Antecedentes: ", r.FormValue("antecedentes"), pdf)
	WriteItem("Revisión por Sistemas: ", r.FormValue("rsistemas"), pdf)
	WriteItem("Signos Vitales: ", fmt.Sprintf("TA: %v/%v FC: %v FR: %v T: %v SPO2: %v Peso: %v Talla: %v IMC: %.2f",
		r.FormValue("tsistolica"),
		r.FormValue("tdiastolica"),
		r.FormValue("fcardiaca"),
		r.FormValue("frespiratoria"),
		r.FormValue("temperatura"),
		r.FormValue("saturacion"),
		r.FormValue("peso"),
		r.FormValue("talla"),
		func() float64 { v, _ := strconv.ParseFloat(r.FormValue("imc"), 64); return v }()),
		pdf)
	WriteItem("Exámen Físico: ", r.FormValue("efisico"), pdf)
	WriteItem("Análisis: ", r.FormValue("analisis"), pdf)
	WriteItem("Diagnósticos: ", r.FormValue("diagnostico"), pdf)
	WriteItem("Conducta: ", r.FormValue("conducta"), pdf)

	if pdf.PageNo() == 2 {
		pdf.Image("reversoFormatoHistoria.png", 0, 0, 612, 792, false, "png", 0, "")
	}
	if err := pdf.Output(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func urgenciaPDF(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	pdf := gofpdf.New("P", "pt", "legal", "")
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 14)
	pdf.Text(444, 152, translator(r.FormValue("cedula")))
	pdf.Text(322, 197, translator(r.FormValue("cedula")))
	pdf.Text(35, 150, translator(r.FormValue("papellido")))
	pdf.Text(180, 150, translator(r.FormValue("sapellido")))
	pdf.Text(300, 150, translator(fmt.Sprintf("%s %s", r.FormValue("pnombre"), r.FormValue("snombre"))))
	pdf.Text(68, 214, r.FormValue("edad"))
	if r.FormValue("genero") == "m" {
		pdf.Text(296, 214, "X")
	} else {
		pdf.Text(321, 214, "X")
	}
	t := time.Now()
	tstring := t.Format(time.RFC3339)
	pdf.SetTitle(fmt.Sprintf("%s %s %s %s %s", r.FormValue("pnombre"), r.FormValue("snombre"), r.FormValue("papellido"), r.FormValue("sapellido"), tstring), true)
	pdf.SetFont("Helvetica", "", 10)
	pdf.Text(39, 340, tstring[8:10])
	pdf.Text(77, 340, tstring[5:7])
	pdf.Text(118, 340, fmt.Sprint(t.Year()))
	pdf.Text(213, 340, fmt.Sprintf("%v   %v", t.Hour(), t.Minute()))

	pdf.SetLeftMargin(45)
	pdf.SetRightMargin(45)
	pdf.SetY(380)

	WriteItemMargin("Motivo de Consulta: ", r.FormValue("mconsulta"), pdf, 15.748)
	WriteItemMargin("Enfermedad Actual: ", r.FormValue("eactual"), pdf, 15.748)
	WriteItemMargin("Antecedentes: ", r.FormValue("antecedentes"), pdf, 15.748)
	WriteItemMargin("Revisión por Sistemas: ", r.FormValue("rsistemas"), pdf, 15.748)
	WriteItemMargin("Exámen Físico: ", r.FormValue("efisico"), pdf, 15.748)
	WriteItemMargin("Análisis: ", r.FormValue("analisis"), pdf, 15.748)
	WriteItemMargin("Diagnósticos: ", r.FormValue("diagnostico"), pdf, 15.748)
	WriteItemMargin("Conducta: ", r.FormValue("conducta"), pdf, 15.748)

	if err := pdf.Output(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func WriteItem(t, d string, pdf *gofpdf.Fpdf) {
	pdf.SetFont("Helvetica", "B", 10)
	pdf.Write(18.8, translator(t))
	pdf.SetFont("Helvetica", "", 10)
	pdf.Write(18.8, translator(d))
	pdf.Ln(18.8)
}

func WriteItemMargin(t, d string, pdf *gofpdf.Fpdf, m float64) {
	pdf.SetFont("Helvetica", "B", 10)
	pdf.Write(m, translator(t))
	pdf.SetFont("Helvetica", "", 10)
	pdf.Write(m, translator(d))
	pdf.Ln(m)
}
