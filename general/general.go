package general

import "encoding/json"

type ConsultaGeneral struct {
	ID           string            `json:"id-consulta"`
	MID          string            `json:"id-medico"`
	PID          string            `json:"id-paciente"`
	Mconsulta    string            `json:"motivo-consulta"`
	Eactual      string            `json:"enfermedad-actual"`
	Antecedentes map[string]string `json:"antecedentes"`
	Rsistemas    string            `json:"revision-sistemas"`

	Efisico struct {
		SVitales  map[string]string `json:"signos-vitales"`
		Contenido string            `json:"contenido"`
	} `json:"examen-fisico"`

	Diagnosticos map[string]string `json:"diagnosticos"`
	Analisis     string            `json:"analisis"`
	HoraInicio   string            `json:"hora-inicio"`
	HoraFinal    string            `json:"hora-final"`
}

func New(j string) (*ConsultaGeneral, error) {
	cg := &ConsultaGeneral{}
	err := json.Unmarshal([]byte(j), cg)
	return cg, err
}
