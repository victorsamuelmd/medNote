package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/victorsamuelmd/mednote/general"
)

var translator, _ = gofpdf.UnicodeTranslatorFromFile("iso-8859-1.map")

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/pdf", unPDF)
	mux.HandleFunc("/json", consultaJson)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.ListenAndServe(":8000", mux)
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

func unPDF(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	pdf := gofpdf.New("P", "pt", "letter", "")
	pdf.AddFont("Open Sans", "", "OpenSans-Regular.json")
	pdf.AddFont("Open Sans", "B", "OpenSans-Bold.json")
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

func WriteItem(t, d string, pdf *gofpdf.Fpdf) {
	pdf.SetFont("Open Sans", "B", 10)
	pdf.Write(18.8, translator(t))
	pdf.SetFont("Open Sans", "", 10)
	pdf.Write(18.8, translator(d))
	pdf.Ln(18.8)
}
