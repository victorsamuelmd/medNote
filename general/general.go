package general

import "encoding/json"

type ConsultaGeneral struct {
	Mconsulta    string
	Eactual      string
	Antecedentes map[string]string
	Rsistemas    string
	Efisico      Exfisico
	Diagnostico  map[string]string
	Analisis     string
}

type Exfisico struct {
	SVitales  map[string]string
	Contenido string
}

func New(j string) (*ConsultaGeneral, error) {
	cg := &ConsultaGeneral{}
	err := json.Unmarshal([]byte(j), cg)
	return cg, err
}
