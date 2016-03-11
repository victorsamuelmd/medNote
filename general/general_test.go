package general

import "testing"

func TestNew(t *testing.T) {
	jsonString := `{
		"motivo-consulta":"un motivo",
		"enfermedad-actual":"enfermo",
		"antecedentes":{
			"Patológicos":"No refiere",
			"Ginecobstétricos":"No refiere"
		},
		"revision-sistemas":"No refiere",
		"examen-fisico":{
			"signos-vitales": {
				"Sistólica":"180",
				"Diastólica":"90"
			},
			"contenido": "Buenas condiciones generales"
		},
		"diagnosticos":{},
		"analisis":"Esta enfermo"}`
	hc, err := New(jsonString)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	if hc.Efisico.Contenido != "Buenas condiciones generales" {
		t.Fail()
	}
}
