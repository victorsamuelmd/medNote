package general

import "encoding/json"

type ConsultaGeneral struct {
	ID           string            `json:"id_consulta"`
	MID          string            `json:"id_medico"`
	PID          string            `json:"id_paciente"`
	Mconsulta    string            `json:"motivo_consulta"`
	Eactual      string            `json:"enfermedad_actual"`
	Antecedentes map[string]string `json:"antecedentes"`
	Rsistemas    string            `json:"revision_sistemas"`

	Efisico struct {
		SVitales  map[string]string `json:"signos_vitales"`
		Contenido string            `json:"contenido"`
	} `json:"examen_fisico"`

	Diagnosticos map[string]string `json:"diagnosticos"`
	Analisis     string            `json:"analisis"`
	HoraInicio   string            `json:"hora_inicio"`
	HoraFinal    string            `json:"hora_final"`
}

func New(j string) (*ConsultaGeneral, error) {
	cg := &ConsultaGeneral{}
	err := json.Unmarshal([]byte(j), cg)
	return cg, err
}
