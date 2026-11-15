package domain_test

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// fechaEsfuerzoUS14 arma una fecha UTC para los escenarios de US-14.
func fechaEsfuerzoUS14(anio int, mes time.Month, dia, hora, minuto int) time.Time {
	return time.Date(anio, mes, dia, hora, minuto, 0, 0, time.UTC)
}

// CA2 / R5 / V2 / V8: las horas deben estar en (0, 24] y tener como máximo
// dos decimales, validados con tolerancia sobre los centésimos.
func TestUS14_ValidarHorasEsfuerzo(t *testing.T) {
	casos := []struct {
		descripcion string
		horas       float64
		valida      bool
	}{
		{
			descripcion: "0.01 horas es válido",
			horas:       0.01,
			valida:      true,
		},
		{
			descripcion: "1.5 horas es válido",
			horas:       1.5,
			valida:      true,
		},
		{
			descripcion: "24 horas es válido",
			horas:       24,
			valida:      true,
		},
		{
			descripcion: "0 horas es inválido",
			horas:       0,
			valida:      false,
		},
		{
			descripcion: "horas negativas son inválidas",
			horas:       -1,
			valida:      false,
		},
		{
			descripcion: "24.01 horas supera el máximo",
			horas:       24.01,
			valida:      false,
		},
		{
			descripcion: "25 horas supera el máximo",
			horas:       25,
			valida:      false,
		},
		{
			descripcion: "1.234 horas tiene más de dos decimales y no se redondea",
			horas:       1.234,
			valida:      false,
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			err := domain.ValidarHorasEsfuerzo(caso.horas)
			if caso.valida {
				if err != nil {
					t.Fatalf("ValidarHorasEsfuerzo(%v) devolvió un error inesperado: %v", caso.horas, err)
				}
				return
			}
			if !errors.Is(err, domain.ErrHorasEsfuerzoInvalidas) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrHorasEsfuerzoInvalidas)
			}
		})
	}
}

// CA1 / R8 / V3: la fecha no puede ser futura y se compara por día de
// calendario, no por instante; el día de hoy se permite a cualquier hora.
func TestUS14_ValidarFechaEsfuerzo(t *testing.T) {
	hoy := fechaEsfuerzoUS14(2026, time.November, 15, 9, 0)

	casos := []struct {
		descripcion string
		fecha       time.Time
		valida      bool
	}{
		{
			descripcion: "hoy a las 00:00 es válido",
			fecha:       fechaEsfuerzoUS14(2026, time.November, 15, 0, 0),
			valida:      true,
		},
		{
			descripcion: "hoy a las 12:30 es válido",
			fecha:       fechaEsfuerzoUS14(2026, time.November, 15, 12, 30),
			valida:      true,
		},
		{
			descripcion: "hoy a las 23:59 es válido",
			fecha:       fechaEsfuerzoUS14(2026, time.November, 15, 23, 59),
			valida:      true,
		},
		{
			descripcion: "ayer a las 23:59 es válido",
			fecha:       fechaEsfuerzoUS14(2026, time.November, 14, 23, 59),
			valida:      true,
		},
		{
			descripcion: "mañana a las 00:00 es futuro",
			fecha:       fechaEsfuerzoUS14(2026, time.November, 16, 0, 0),
			valida:      false,
		},
		{
			descripcion: "mañana a las 23:59 es futuro",
			fecha:       fechaEsfuerzoUS14(2026, time.November, 16, 23, 59),
			valida:      false,
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			err := domain.ValidarFechaEsfuerzo(caso.fecha, hoy)
			if caso.valida {
				if err != nil {
					t.Fatalf("ValidarFechaEsfuerzo(%v, %v) devolvió un error inesperado: %v", caso.fecha, hoy, err)
				}
				return
			}
			if !errors.Is(err, domain.ErrFechaEsfuerzoFutura) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrFechaEsfuerzoFutura)
			}
		})
	}
}

// D6 / D29: RegistroEsfuerzo conserva los cuatro datos de CA1.
func TestUS14_RegistroEsfuerzo_Campos(t *testing.T) {
	fecha := fechaEsfuerzoUS14(2026, time.November, 15, 9, 0)
	registro := domain.RegistroEsfuerzo{
		Integrante: "Ana Pérez",
		Fecha:      fecha,
		Actividad:  "Análisis",
		Horas:      1.5,
	}

	if registro.Integrante != "Ana Pérez" {
		t.Errorf("Integrante = %q, se esperaba %q", registro.Integrante, "Ana Pérez")
	}
	if !registro.Fecha.Equal(fecha) {
		t.Errorf("Fecha = %v, se esperaba %v", registro.Fecha, fecha)
	}
	if registro.Actividad != "Análisis" {
		t.Errorf("Actividad = %q, se esperaba %q", registro.Actividad, "Análisis")
	}
	if registro.Horas != 1.5 {
		t.Errorf("Horas = %v, se esperaba 1.5", registro.Horas)
	}
}

// D7 / D24 / D28: los tres errores de esfuerzo son sentinelas distintos.
func TestUS14_ErroresEsfuerzo_SonSentinelasDistintas(t *testing.T) {
	errores := map[string]error{
		"ErrIntegranteNoPerteneceAlProyecto": domain.ErrIntegranteNoPerteneceAlProyecto,
		"ErrHorasEsfuerzoInvalidas":          domain.ErrHorasEsfuerzoInvalidas,
		"ErrFechaEsfuerzoFutura":             domain.ErrFechaEsfuerzoFutura,
	}
	for nombre, err := range errores {
		if err == nil {
			t.Fatalf("%s no está definido", nombre)
		}
	}
	for nombre, err := range errores {
		for otroNombre, otro := range errores {
			if nombre == otroNombre {
				continue
			}
			if errors.Is(err, otro) {
				t.Errorf("%s no debe coincidir con %s", nombre, otroNombre)
			}
		}
	}

}

func TestUS17_ValidarHorasEstimadas(t *testing.T) {
	casos := []struct {
		nombre string
		horas  float64
		valida bool
	}{
		{nombre: "cero", horas: 0, valida: true},
		{nombre: "decimal", horas: 1.25, valida: true},
		{nombre: "horas mayores a un día", horas: 30, valida: true},
		{nombre: "negativa", horas: -0.01},
		{nombre: "más de dos decimales", horas: 1.234},
		{nombre: "NaN", horas: math.NaN()},
		{nombre: "infinito", horas: math.Inf(1)},
		{nombre: "desborda al convertir a centésimos", horas: math.MaxFloat64},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			err := domain.ValidarHorasEstimadas(caso.horas)
			if caso.valida && err != nil {
				t.Fatalf("ValidarHorasEstimadas(%v) devolvió error: %v", caso.horas, err)
			}
			if !caso.valida && !errors.Is(err, domain.ErrHorasEstimadasInvalidas) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrHorasEstimadasInvalidas)
			}
		})
	}
}

func TestUS17_CalcularDesviacion(t *testing.T) {
	casos := []struct {
		nombre    string
		estimadas float64
		reales    float64
		esperada  *float64
	}{
		{nombre: "sobreestimación real", estimadas: 10, reales: 12, esperada: float64Puntero(20)},
		{nombre: "subestimación real", estimadas: 10, reales: 8, esperada: float64Puntero(-20)},
		{nombre: "estimación cero", estimadas: 0, reales: 2},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			obtenida := domain.CalcularDesviacion(caso.estimadas, caso.reales)
			if caso.esperada == nil {
				if obtenida != nil {
					t.Fatalf("desviación = %v, se esperaba nil", *obtenida)
				}
				return
			}
			if obtenida == nil || *obtenida != *caso.esperada {
				t.Fatalf("desviación = %v, se esperaba %v", obtenida, *caso.esperada)
			}
		})
	}
}

func float64Puntero(valor float64) *float64 {
	return &valor
}
