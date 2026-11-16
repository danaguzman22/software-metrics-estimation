package web

import "testing"

// D15: la escala de los gráficos se calcula en la capa web; el tope del eje Y
// es el máximo de las barras, con mínimo 1 (sin división por cero).
func TestUS19_Fase2_EscalaConMinimoUno_TodoEnCero(t *testing.T) {
	if got := escalaConMinimoUno(); got != 1 {
		t.Errorf("escalaConMinimoUno() = %v, se esperaba 1 (sin barras, eje mínimo 1)", got)
	}
	if got := escalaConMinimoUno(0, 0, 0); got != 1 {
		t.Errorf("escalaConMinimoUno(0, 0, 0) = %v, se esperaba 1 (todo en cero, eje mínimo 1)", got)
	}
}

func TestUS19_Fase2_EscalaConMinimoUno_MaximoPositivo(t *testing.T) {
	casos := []struct {
		nombre   string
		valores  []float64
		esperado float64
	}{
		{"un valor", []float64{13}, 13},
		{"varios valores enteros", []float64{8, 13}, 13},
		{"valores con decimales", []float64{40.5, 35, 0}, 40.5},
		{"máximo al final", []float64{3, 1, 7}, 7},
		{"horas estimadas y reales", []float64{111.1, 222.2, 40, 35}, 222.2},
	}
	for _, caso := range casos {
		if got := escalaConMinimoUno(caso.valores...); got != caso.esperado {
			t.Errorf("%s: escalaConMinimoUno(%v) = %v, se esperaba %v", caso.nombre, caso.valores, got, caso.esperado)
		}
	}
}

func TestUS19_Fase2_EscalaConMinimoUno_ValoresGrandes(t *testing.T) {
	if got := escalaConMinimoUno(1000000, 2000000); got != 2000000 {
		t.Errorf("escalaConMinimoUno(1000000, 2000000) = %v, se esperaba 2000000", got)
	}
	if got := escalaConMinimoUno(1e12); got != 1e12 {
		t.Errorf("escalaConMinimoUno(1e12) = %v, se esperaba 1e12", got)
	}
}
