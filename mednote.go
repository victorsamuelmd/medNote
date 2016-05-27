package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	"github.com/justinas/alice"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/victorsamuelmd/mednote/general"
)

const (
	MgoHostStr             = "localhost:27017"
	NombreBaseDatos        = "mednote"
	NombreCollecionUsuario = "usuario"
	NombreBaseDatosTest    = "test"
	NumeroPuertoAplicacion = ":8000"
	TipoContenidoJSONAPI   = "application/vnd.api+json"
	OrigenesConfiables     = "http://localhost:4200"
)

var (
	utf8toIso, _ = gofpdf.UnicodeTranslatorFromFile("iso-8859-1.map")
)

func main() {

	var ds = NewDataStore()

	m := mux.NewRouter()
	router := http.NewServeMux()

	router.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static"))))
	router.Handle("/dist/", http.StripPrefix("/dist/",
		http.FileServer(http.Dir("./dist"))))

	router.Handle("/pdf", alice.New(logger).ThenFunc(ageneralPDF))
	router.HandleFunc("/remision", remisionPDF)
	router.HandleFunc("/urgencia", urgenciaPDF)
	router.HandleFunc("/formula", formulaPDF)
	m.HandleFunc("/token", ds.token)
	m.Handle("/api/usuarios/{id:[0-9]{4,14}}", alice.New(proteger).
		ThenFunc(ds.usuariosGET)).
		Methods(http.MethodGet)

	m.Handle("/api/usuarios", alice.New(proteger).
		ThenFunc(ds.usuariosPOST)).
		Methods(http.MethodPost)

	fmt.Println("Listening on localhost:8000")
	if err := http.ListenAndServe(NumeroPuertoAplicacion, router); err != nil {
		fmt.Print(err.Error())
	}
}

func (store *DataStore) usuariosGET(w http.ResponseWriter, r *http.Request) {
	ds := store.Copy()
	defer ds.Session.Close()

	usuarioId := mux.Vars(r)["id"]

	usr := &Usuario{}
	err := ds.Session.DB(NombreBaseDatos).
		C(NombreCollecionUsuario).
		Find(bson.M{"identificacion": usuarioId}).
		Select(bson.M{"contrasena": 0}).
		One(&usr)

	if err != nil {
		http.Error(w, fmt.Sprintf(
			`{"errors": [{"title": "No encontrado", "detail": "%s"}]}`,
			err.Error()), http.StatusNotFound)
		return
	}

	j, err := jsonapi.Marshal(usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprintf(w, "%s", j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
	} else if r.Method == http.MethodGet {
		var u []Usuario
		err := store.session.DB(NombreBaseDatos).
			C(NombreCollecionUsuario).Find(nil).
			Select(bson.M{"contrasena": 0}).One(&u)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		j, err := jsonapi.Marshal(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = fmt.Fprintf(w, "%s", j)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
*/

func (store *DataStore) usuariosPOST(w http.ResponseWriter, r *http.Request) {
	ds := store.Copy()
	defer ds.Session.Close()

	var j []byte
	b := bytes.NewBuffer(j)
	reader := bufio.NewReader(r.Body)
	_, err := reader.WriteTo(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u := &Usuario{}
	err = jsonapi.Unmarshal(b.Bytes(), u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ds.Session.DB(NombreBaseDatos).
		C(NombreCollecionUsuario).Insert(u)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", u)
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\033[33m[%s] %s\033[0m\n",
			time.Now().Format(time.Stamp), r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func proteger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", TipoContenidoJSONAPI)
		token := r.Header.Get("Authorization")
		if len(token) == 0 {
			http.Error(w, `{"errors": [{"title": "Fallo autorizacion",
			"detail": "No fue provisto un token de autorizacion"}]}`,
				http.StatusUnauthorized)
			return
		}
		if ok := validarToken(strings.TrimPrefix(token, "Bearer ")); !ok {
			http.Error(w, `{"errors": [{"title": "Fallo autorizacion",
			"detail": "El token es inválido"}]}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// activarCORS configura los encabezados para que se puedan realizar peticiones
// desde un origen diferente. En caso de que el navegador solicite opciones
// se comporta adecuadamente.
func activarCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", OrigenesConfiables)
		w.Header().Add("Access-Control-Allow-Methods", http.MethodPost)
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			return
		}

		w.Header().Add("content-type", TipoContenidoJSONAPI)
		next.ServeHTTP(w, r)
	})
}

// token genera un JWT(javascript web token) usado como autorizacion para
// acceder a la API del servidor, la petición del cliente debe ser un objeto en
// JSON de la forma {"username": "", "password": ""} y devuelve un el token
// codificado dentro de un objeto {"token": "<token>"}
func (store *DataStore) token(w http.ResponseWriter, r *http.Request) {
	ds := store.Copy()
	defer ds.Session.Close()

	db := ds.Session.DB(NombreBaseDatos)

	var rd struct {
		NombreUsuario string `json:"username"`
		Contraseña    string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&rd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := AutenticarUsuario(rd.NombreUsuario, rd.Contraseña, db)
	if err != nil {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	if err := json.NewEncoder(w).Encode(
		map[string]string{"token": token}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
	tstring := t.Format(time.RFC3339)
	name := fmt.Sprintf("%s %s %s %s %s",
		r.FormValue("pnombre"),
		r.FormValue("snombre"),
		r.FormValue("papellido"),
		r.FormValue("sapellido"),
		tstring)

	w.Header().Set("Content-Disposition",
		fmt.Sprintf("filename=\"%s.pdf\"", name))

	pdf := gofpdf.New("P", "pt", "letter", "")
	pdf.AddPage()
	//	pdf.Image("frenteFormatoHistoria.png",
	//		0, 0, 612, 792, false, "png", 0, "")
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
	pdf.Text(37, 378, tstring[8:10])
	pdf.Text(62, 378, tstring[5:7])
	pdf.Text(85, 378, tstring[2:4])
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

	if pdf.PageNo() == 2 {
		//		pdf.Image("reversoFormatoHistoria.png",
		//			0, 0, 612, 792, false, "png", 0, "")
	}
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

	w.Header().Set("Content-Disposition",
		fmt.Sprintf("filename=\"%s.pdf\"", r.FormValue("pnombre")))

	pdf := gofpdf.NewCustom(
		&gofpdf.InitType{
			"P", "pt", "pt",
			gofpdf.SizeType{625.5, 922.5}, "Helvetica"})
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
	tstring := t.Format(time.RFC3339)
	pdf.SetTitle(fmt.Sprintf("%s %s %s %s %s",
		r.FormValue("pnombre"),
		r.FormValue("snombre"),
		r.FormValue("papellido"),
		r.FormValue("sapellido"), tstring), true)
	pdf.SetFont("Helvetica", "", 10)
	pdf.Text(39, 340, tstring[8:10])
	pdf.Text(77, 340, tstring[5:7])
	pdf.Text(118, 340, fmt.Sprint(t.Year()))
	pdf.Text(213, 340, fmt.Sprintf("%v   %v", t.Hour(), t.Minute()))

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
