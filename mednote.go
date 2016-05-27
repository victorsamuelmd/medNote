package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
)

var (
	utf8toIso, _ = gofpdf.UnicodeTranslatorFromFile("iso-8859-1.map")
)

func main() {

	router := http.NewServeMux()

	router.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static"))))

	router.HandleFunc("/pdf", ageneralPDF)
	router.HandleFunc("/remision", remisionPDF)
	router.HandleFunc("/urgencia", urgenciaPDF)
	router.HandleFunc("/formula", formulaPDF)

	fmt.Println("Listening on localhost:8000")
	if err := http.ListenAndServe(":8000", router); err != nil {
		fmt.Print(err.Error())
	}
}

func remisionPDF(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Disposition",
		"inline; filename=\"remision.pdf\"")
	var b bytes.Buffer
	pdf, err := CrearPDFRemision(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pdf.Output(&b)
	fmt.Fprint(w, base64.RawStdEncoding.EncodeToString(b.Bytes()))
}

func ageneralPDF(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	t := time.Now()
	name := fmt.Sprintf("%s %s %s %s %s",
		r.FormValue("pnombre"),
		r.FormValue("snombre"),
		r.FormValue("papellido"),
		t.Format(time.RFC3339))

	w.Header().Set("Content-Disposition",
		fmt.Sprintf(`filename="%s.pdf"`, name))

	pdf := gofpdf.New("P", "pt", "letter", "")
	pdf.AddPage()

	pdf.SetFont("Helvetica", "", 14)
	pdf.Text(440, 195, utf8toIso(r.FormValue("cedula")))
	pdf.Text(40, 185, utf8toIso(r.FormValue("papellido")))
	pdf.Text(180, 185, utf8toIso(r.FormValue("sapellido")))
	pdf.Text(300, 185, utf8toIso(fmt.Sprintf("%s %s",
		r.FormValue("pnombre"), r.FormValue("snombre"))))
	pdf.Text(40, 215, r.FormValue("edad"))
	if r.FormValue("genero") == "m" {
		pdf.Text(145, 225, "X")
	} else {
		pdf.Text(180, 225, "X")
	}
	pdf.SetTitle(name, true)
	pdf.SetFont("Helvetica", "", 10)
	pdf.Text(37, 378, t.Format("02"))
	pdf.Text(62, 378, t.Format("01"))
	pdf.Text(85, 378, t.Format("06"))
	pdf.Text(109, 378, t.Format("15:04"))
	pdf.SetLeftMargin(140)
	pdf.SetRightMargin(52)
	pdf.SetTopMargin(80)
	pdf.SetY(365)

	WriteItem("Motivo de Consulta: ", r.FormValue("mconsulta"), pdf)
	WriteItem("Enfermedad Actual: ", r.FormValue("eactual"), pdf)
	WriteItem("Antecedentes: ", r.FormValue("antecedentes"), pdf)
	WriteItem("Revisión por Sistemas: ", r.FormValue("rsistemas"), pdf)
	WriteItem("Signos Vitales: ",
		fmt.Sprintf(
			"TA: %v/%v FC: %v FR: %v T: %v SPO2: %v Peso: %v Talla: %v IMC: %.2f",
			r.FormValue("tsistolica"),
			r.FormValue("tdiastolica"),
			r.FormValue("fcardiaca"),
			r.FormValue("frespiratoria"),
			r.FormValue("temperatura"),
			r.FormValue("saturacion"),
			r.FormValue("peso"),
			r.FormValue("talla"),
			func() float64 {
				v, _ := strconv.ParseFloat(
					r.FormValue("imc"), 64)
				return v
			}()),
		pdf)
	WriteItem("Exámen Físico: ", r.FormValue("efisico"), pdf)
	WriteItem("Análisis: ", r.FormValue("analisis"), pdf)
	WriteItem("Diagnósticos: ", r.FormValue("diagnostico"), pdf)
	WriteItem("Conducta: ", r.FormValue("conducta"), pdf)

	if err := pdf.Output(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func urgenciaPDF(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	t, err := time.Parse(time.RFC3339, r.FormValue("date"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	usrName := fmt.Sprintf("%s %s %s %s %s",
		r.FormValue("pnombre"),
		r.FormValue("snombre"),
		r.FormValue("papellido"),
		r.FormValue("sapellido"), t.Format(time.RFC3339))

	w.Header().Set("Content-Disposition",
		fmt.Sprintf("filename=\"%s.pdf\"", usrName))

	pdf := gofpdf.NewCustom(
		&gofpdf.InitType{
			"P", "pt", "pt",
			gofpdf.SizeType{625.5, 922.5}, "Helvetica"})
	pdf.SetTitle(usrName, true)
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 14)
	pdf.Text(444, 152, utf8toIso(r.FormValue("cedula")))
	pdf.Text(322, 197, utf8toIso(r.FormValue("cedula")))
	pdf.Text(35, 150, utf8toIso(r.FormValue("papellido")))
	pdf.Text(180, 150, utf8toIso(r.FormValue("sapellido")))
	pdf.Text(300, 150, utf8toIso(fmt.Sprintf("%s %s",
		r.FormValue("pnombre"),
		r.FormValue("snombre"))))
	pdf.Text(68, 214, r.FormValue("edad"))
	if r.FormValue("genero") == "m" {
		pdf.Text(296, 214, "X")
	} else {
		pdf.Text(321, 214, "X")
	}
	pdf.SetFont("Helvetica", "", 10)
	pdf.Text(39, 340, t.Format("02"))
	pdf.Text(77, 340, t.Format("01"))
	pdf.Text(118, 340, t.Format("2006"))
	pdf.Text(213, 340, t.Format("15   06"))

	pdf.SetLeftMargin(45)
	pdf.SetRightMargin(45)
	pdf.SetY(380)

	WriteItemMargin("Motivo de Consulta: ",
		r.FormValue("mconsulta"), pdf, 15.748)
	WriteItemMargin("Enfermedad Actual: ",
		r.FormValue("eactual"), pdf, 15.748)
	WriteItemMargin("Antecedentes: ",
		r.FormValue("antecedentes"), pdf, 15.748)
	WriteItemMargin("Revisión por Sistemas: ",
		r.FormValue("rsistemas"), pdf, 15.748)
	WriteItemMargin("Exámen Físico: ",
		r.FormValue("efisico"), pdf, 15.748)
	WriteItemMargin("Análisis: ",
		r.FormValue("analisis"), pdf, 15.748)
	WriteItemMargin("Diagnósticos: ",
		r.FormValue("diagnostico"), pdf, 15.748)
	WriteItemMargin("Conducta: ",
		r.FormValue("conducta"), pdf, 15.748)

	if err := pdf.Output(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func formulaPDF(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	f := &Formula{
		utf8toIso(fmt.Sprintf("%s %s %s %s", r.FormValue("pnombre"),
			r.FormValue("snombre"),
			r.FormValue("papellido"),
			r.FormValue("sapellido"))),
		utf8toIso(r.FormValue("cedula")),
		utf8toIso(r.FormValue("centro-salud")),
		utf8toIso(r.FormValue("eps")),
		utf8toIso(r.FormValue("conducta")),
	}

	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr:    "pt",
		Size:       gofpdf.SizeType{396, 612},
		FontDirStr: "",
	})

	pdf.SetFont("Helvetica", "", 14)
	pdf.AddPage()
	pdf.Text(220, 84, f.Id)
	t := time.Now()
	pdf.Text(60, 106, t.Format("2006  01  02"))
	pdf.Text(240, 106, f.CentroSalud)
	pdf.Text(55, 128, f.Nombre)
	pdf.Text(65, 170, f.EPS)
	pdf.SetY(200)
	pdf.Write(16, f.Receta)
	pdf.Output(w)
}

func WriteItem(t, d string, pdf *gofpdf.Fpdf) {
	pdf.SetFont("Helvetica", "B", 10)
	pdf.Write(18.8, utf8toIso(t))
	pdf.SetFont("Helvetica", "", 10)
	pdf.Write(18.8, utf8toIso(d))
	pdf.Ln(18.8)
}

func WriteItemMargin(t, d string, pdf *gofpdf.Fpdf, m float64) {
	pdf.SetFont("Helvetica", "B", 10)
	pdf.Write(m, utf8toIso(t))
	pdf.SetFont("Helvetica", "", 10)
	pdf.Write(m, utf8toIso(d))
	pdf.Ln(m)
}

type Formula struct {
	Nombre      string `json:"nombre"`
	Id          string `json:"id"`
	CentroSalud string `json:"centro-salud"`
	EPS         string `json:"eps"`
	Receta      string `json:"receta"`
}
