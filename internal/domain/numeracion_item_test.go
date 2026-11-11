package domain_test

import (
	"reflect"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

const sufijoOverflowT05 = "HU-123456789012345678901234567890"

func TestT05_NumeroDeItem_TablaDeCasos(t *testing.T) {
	casos := []struct {
		nombre   string
		id       string
		esperado int
		ok       bool
	}{
		{"id_valido_HU01", "HU-01", 1, true},
		{"id_valido_HU04", "HU-04", 4, true},
		{"id_valido_HU100", "HU-100", 100, true},
		{"sin_cero_a_la_izquierda", "HU-9", 9, true},
		{"ceros_a_la_izquierda", "HU-001", 1, true},
		{"prefijo_incorrecto", "XX-01", 0, false},
		{"capitalizacion_incorrecta", "hu-01", 0, false},
		{"vacio", "", 0, false},
		{"prefijo_sin_numero", "HU-", 0, false},
		{"contenido_no_numerico", "HU-abc", 0, false},
		{"caracteres_extra", "HU-01x", 0, false},
		{"espacio_externo_izquierda", " HU-01", 0, false},
		{"espacio_externo_derecha", "HU-01 ", 0, false},
		{"cero", "HU-0", 0, false},
		{"cero_doble", "HU-00", 0, false},
		{"signo_negativo", "HU--1", 0, false},
		{"signo_positivo", "HU-+1", 0, false},
		{"overflow", sufijoOverflowT05, 0, false},
	}

	for _, caso := range casos {
		caso := caso
		t.Run(caso.nombre, func(t *testing.T) {
			numero, ok := domain.NumeroDeItem(caso.id)
			if numero != caso.esperado {
				t.Errorf("NumeroDeItem(%q) número = %d, se esperaba %d", caso.id, numero, caso.esperado)
			}
			if ok != caso.ok {
				t.Errorf("NumeroDeItem(%q) ok = %v, se esperaba %v", caso.id, ok, caso.ok)
			}
		})
	}
}

func TestT05_SiguienteIDItem_TablaDeCasos(t *testing.T) {
	casos := []struct {
		nombre   string
		ids      []string
		esperado string
	}{
		{"lista_vacia", []string{}, "HU-01"},
		{"lista_nil", nil, "HU-01"},
		{"solo_malformados", []string{"XX-01", "HU-abc", ""}, "HU-01"},
		{"un_solo_id", []string{"HU-01"}, "HU-02"},
		{"secuencia", []string{"HU-01", "HU-02", "HU-03"}, "HU-04"},
		{"formato_dos_digitos", []string{"HU-03"}, "HU-04"},
		{"hueco", []string{"HU-01", "HU-03"}, "HU-04"},
		{"orden_irrelevante", []string{"HU-03", "HU-01"}, "HU-04"},
		{"ids_duplicados", []string{"HU-01", "HU-01"}, "HU-02"},
		{"comparacion_numerica", []string{"HU-9", "HU-10"}, "HU-11"},
		{"malformados_ignorados", []string{"HU-01", "XX-99", "HU-abc", "", "HU-03"}, "HU-04"},
		{"overflow_ignorado", []string{"HU-01", sufijoOverflowT05}, "HU-02"},
		{"mas_de_99", []string{"HU-99"}, "HU-100"},
		{"mas_de_99_siguiente", []string{"HU-100"}, "HU-101"},
	}

	for _, caso := range casos {
		caso := caso
		t.Run(caso.nombre, func(t *testing.T) {
			obtenido := domain.SiguienteIDItem(caso.ids)
			if obtenido != caso.esperado {
				t.Errorf("SiguienteIDItem(%v) = %q, se esperaba %q", caso.ids, obtenido, caso.esperado)
			}
		})
	}
}

func TestT05_SiguienteIDItem_HuecoNoRellena(t *testing.T) {
	obtenido := domain.SiguienteIDItem([]string{"HU-01", "HU-03"})
	if obtenido != "HU-04" {
		t.Errorf("SiguienteIDItem(hueco) = %q, se esperaba %q", obtenido, "HU-04")
	}
	if obtenido == "HU-02" {
		t.Errorf("SiguienteIDItem(hueco) = %q, no debe rellenar el hueco con HU-02", obtenido)
	}
}

func TestT05_SiguienteIDItem_ComparacionNumericaNoTextual(t *testing.T) {
	obtenido := domain.SiguienteIDItem([]string{"HU-9", "HU-10"})
	if obtenido != "HU-11" {
		t.Errorf("SiguienteIDItem([HU-9 HU-10]) = %q, se esperaba %q", obtenido, "HU-11")
	}
	if obtenido == "HU-10" {
		t.Errorf("SiguienteIDItem([HU-9 HU-10]) = %q, comparó textualmente en vez de numéricamente", obtenido)
	}
}

func TestT05_SiguienteIDItem_FormatoDosDigitos(t *testing.T) {
	obtenido := domain.SiguienteIDItem([]string{"HU-03"})
	if obtenido != "HU-04" {
		t.Errorf("SiguienteIDItem([HU-03]) = %q, se esperaba %q", obtenido, "HU-04")
	}
	if obtenido == "HU-4" {
		t.Errorf("SiguienteIDItem([HU-03]) = %q, se esperaba el formato de dos dígitos HU-04", obtenido)
	}
}

func TestT05_SiguienteIDItem_NoMutaLaEntrada(t *testing.T) {
	entrada := []string{"HU-01", "XX-99", "HU-03"}
	copia := append([]string(nil), entrada...)

	domain.SiguienteIDItem(entrada)

	if !reflect.DeepEqual(entrada, copia) {
		t.Errorf("SiguienteIDItem mutó la entrada: obtenido %v, se esperaba %v", entrada, copia)
	}
}

func TestT05_SiguienteIDItem_Determinista(t *testing.T) {
	entrada := []string{"HU-09", "HU-10"}

	primera := domain.SiguienteIDItem(entrada)
	segunda := domain.SiguienteIDItem(entrada)

	if primera != segunda {
		t.Errorf("SiguienteIDItem no es determinista: primera = %q, segunda = %q", primera, segunda)
	}
}
