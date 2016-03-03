package general

import "testing"

func TestNew(t *testing.T) {
	jsonString := `{
		"Mconsulta":"un motivo",
		"Eactual":"enfermo",
		"Antecedentes":{
			"Patológicos":"No refiere",
			"Ginecobstétricos":"No refiere"
		},
		"Rsistemas":"No refiere",
		"Efisico":{
			"SVitales": {
				"Sistólica":"180",
				"Diastólica":"90"
			},
			"Contenido": "Buenas condiciones generales"
		},
		"Diagnostico":{},
		"Analisis":"Esta enfermo"}`
	hc, err := New(jsonString)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	if hc.Efisico.Contenido != "Buenas condiciones generales" {
		t.Fail()
	}
}
