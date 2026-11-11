package domain_test

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// CA1 / R1 / D1: con "hoy" anterior a la fecha de inicio el estado es Planificado.
func TestUS04_EstadoPorFechas_Planificado(t *testing.T) {
	casos := []struct {
		descripcion string
		hoy         time.Time
	}{
		{
			descripcion: "día anterior al inicio",
			hoy:         fechaUS01(2026, time.February, 28),
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			estado, err := domain.EstadoPorFechas(
				fechaUS01(2026, time.March, 1),
				fechaUS01(2026, time.June, 30),
				caso.hoy,
			)
			if err != nil {
				t.Fatalf("EstadoPorFechas() devolvió un error inesperado: %v", err)
			}
			if estado != domain.EstadoPlanificado {
				t.Errorf("Estado = %q, se esperaba %q", estado, domain.EstadoPlanificado)
			}
		})
	}
}

// CA1 / R1 / D1: los extremos inicio y fin están incluidos como En curso, igual
// que una fecha intermedia.
func TestUS04_EstadoPorFechas_EnCursoExtremos(t *testing.T) {
	casos := []struct {
		descripcion string
		hoy         time.Time
	}{
		{
			descripcion: "hoy igual a la fecha de inicio",
			hoy:         fechaUS01(2026, time.March, 1),
		},
		{
			descripcion: "hoy en una fecha intermedia",
			hoy:         fechaUS01(2026, time.April, 15),
		},
		{
			descripcion: "hoy igual a la fecha de fin",
			hoy:         fechaUS01(2026, time.June, 30),
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			estado, err := domain.EstadoPorFechas(
				fechaUS01(2026, time.March, 1),
				fechaUS01(2026, time.June, 30),
				caso.hoy,
			)
			if err != nil {
				t.Fatalf("EstadoPorFechas() devolvió un error inesperado: %v", err)
			}
			if estado != domain.EstadoEnCurso {
				t.Errorf("Estado = %q, se esperaba %q", estado, domain.EstadoEnCurso)
			}
		})
	}
}

// CA1 / R1 / D1: con "hoy" posterior a la fecha de fin el estado es Finalizado.
func TestUS04_EstadoPorFechas_Finalizado(t *testing.T) {
	casos := []struct {
		descripcion string
		hoy         time.Time
	}{
		{
			descripcion: "día siguiente al fin",
			hoy:         fechaUS01(2026, time.July, 1),
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			estado, err := domain.EstadoPorFechas(
				fechaUS01(2026, time.March, 1),
				fechaUS01(2026, time.June, 30),
				caso.hoy,
			)
			if err != nil {
				t.Fatalf("EstadoPorFechas() devolvió un error inesperado: %v", err)
			}
			if estado != domain.EstadoFinalizado {
				t.Errorf("Estado = %q, se esperaba %q", estado, domain.EstadoFinalizado)
			}
		})
	}
}

// CA1 / R1 / D7: la comparación de "hoy" usa solo el día, ignorando la hora.
func TestUS04_EstadoPorFechas_ComparacionPorDia(t *testing.T) {
	inicio := fechaUS01(2026, time.March, 1)
	fin := fechaUS01(2026, time.June, 30)

	casos := []struct {
		descripcion string
		hoy         time.Time
	}{
		{
			descripcion: "mismo día a las 8:00",
			hoy:         time.Date(2026, time.March, 1, 8, 0, 0, 0, time.UTC),
		},
		{
			descripcion: "mismo día a las 20:00",
			hoy:         time.Date(2026, time.March, 1, 20, 0, 0, 0, time.UTC),
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			estado, err := domain.EstadoPorFechas(inicio, fin, caso.hoy)
			if err != nil {
				t.Fatalf("EstadoPorFechas() devolvió un error inesperado: %v", err)
			}
			if estado != domain.EstadoEnCurso {
				t.Errorf("Estado = %q, se esperaba %q", estado, domain.EstadoEnCurso)
			}
		})
	}
}

// CA1 / R1 / D9: una fecha de fin anterior al inicio devuelve
// ErrFechaFinAnteriorAInicio.
func TestUS04_EstadoPorFechas_FechaFinAnteriorAInicio(t *testing.T) {
	casos := []struct {
		descripcion string
		inicio      time.Time
		fin         time.Time
	}{
		{
			descripcion: "fin anterior al inicio",
			inicio:      fechaUS01(2026, time.June, 30),
			fin:         fechaUS01(2026, time.March, 1),
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			_, err := domain.EstadoPorFechas(caso.inicio, caso.fin, fechaUS01(2026, time.April, 15))
			if err == nil {
				t.Fatal("EstadoPorFechas() con fecha de fin anterior al inicio debería devolver un error")
			}
			if !errors.Is(err, domain.ErrFechaFinAnteriorAInicio) {
				t.Errorf("err = %v, se esperaba %v", err, domain.ErrFechaFinAnteriorAInicio)
			}
		})
	}
}

// CA1 / R5 / D1: el estado calculado puede diferir del guardado y no lo modifica.
func TestUS04_EstadoCalculadoNoModificaEstadoGuardado(t *testing.T) {
	proyecto := domain.Proyecto{
		Nombre:      "Sistema de Métricas",
		FechaInicio: fechaUS01(2026, time.March, 1),
		FechaFin:    fechaUS01(2026, time.June, 30),
		Estado:      domain.EstadoPlanificado,
	}

	estado, err := domain.EstadoPorFechas(
		proyecto.FechaInicio,
		proyecto.FechaFin,
		fechaUS01(2026, time.July, 1),
	)
	if err != nil {
		t.Fatalf("EstadoPorFechas() devolvió un error inesperado: %v", err)
	}
	if estado != domain.EstadoFinalizado {
		t.Errorf("Estado = %q, se esperaba %q", estado, domain.EstadoFinalizado)
	}
	if proyecto.Estado != domain.EstadoPlanificado {
		t.Errorf("Estado guardado = %q, no debería modificarse de %q", proyecto.Estado, domain.EstadoPlanificado)
	}
	if estado == proyecto.Estado {
		t.Errorf("Estado calculado = %q, debería diferir del guardado %q", estado, proyecto.Estado)
	}
}

// CA1 / R3 / D4: el porcentaje se calcula y redondea a 2 decimales.
func TestUS04_PorcentajeHistoriasCompletadas(t *testing.T) {
	const tolerancia = 1e-9

	casos := []struct {
		descripcion string
		total       int
		completadas int
		esperado    float64
	}{
		{
			descripcion: "0 de 3",
			total:       3,
			completadas: 0,
			esperado:    0,
		},
		{
			descripcion: "1 de 3",
			total:       3,
			completadas: 1,
			esperado:    33.33,
		},
		{
			descripcion: "2 de 3",
			total:       3,
			completadas: 2,
			esperado:    66.67,
		},
		{
			descripcion: "3 de 3",
			total:       3,
			completadas: 3,
			esperado:    100,
		},
		{
			descripcion: "5 de 8",
			total:       8,
			completadas: 5,
			esperado:    62.5,
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			obtenido, err := domain.PorcentajeCompletadas(caso.total, caso.completadas)
			if err != nil {
				t.Fatalf("PorcentajeCompletadas() devolvió un error inesperado: %v", err)
			}
			if math.Abs(obtenido-caso.esperado) > tolerancia {
				t.Errorf("PorcentajeCompletadas(%d, %d) = %v, se esperaba %v", caso.total, caso.completadas, obtenido, caso.esperado)
			}
		})
	}
}

// CA2 / R4 / D4: un proyecto sin historias devuelve 0 sin error ni división por cero.
func TestUS04_PorcentajeSinHistorias(t *testing.T) {
	obtenido, err := domain.PorcentajeCompletadas(0, 0)
	if err != nil {
		t.Fatalf("PorcentajeCompletadas(0, 0) devolvió un error inesperado: %v", err)
	}
	if obtenido != 0 {
		t.Errorf("PorcentajeCompletadas(0, 0) = %v, se esperaba 0", obtenido)
	}
}

// CA1 / R3 / R4 / D6: cantidades negativas o completadas mayores que el total
// devuelven ErrHistoriasInvalidas.
func TestUS04_PorcentajeCompletadas_ValoresInvalidos(t *testing.T) {
	casos := []struct {
		descripcion string
		total       int
		completadas int
	}{
		{
			descripcion: "total negativo",
			total:       -1,
			completadas: 0,
		},
		{
			descripcion: "completadas negativas",
			total:       3,
			completadas: -1,
		},
		{
			descripcion: "completadas mayor que el total",
			total:       3,
			completadas: 4,
		},
		{
			descripcion: "total cero con completadas mayor que cero",
			total:       0,
			completadas: 1,
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			_, err := domain.PorcentajeCompletadas(caso.total, caso.completadas)
			if err == nil {
				t.Fatal("PorcentajeCompletadas() con valores inválidos debería devolver un error")
			}
			if !errors.Is(err, domain.ErrHistoriasInvalidas) {
				t.Errorf("err = %v, se esperaba %v", err, domain.ErrHistoriasInvalidas)
			}
		})
	}
}

// CA1 / R1 / D7: el ultimo dia del proyecto, con una hora posterior a la de fin,
// sigue siendo En curso porque la comparacion usa solo el dia.
func TestUS04_EstadoPorFechas_ComparacionPorDia_UltimoDia(t *testing.T) {
	estado, err := domain.EstadoPorFechas(
		fechaUS01(2026, time.March, 1),
		fechaUS01(2026, time.June, 30),
		time.Date(2026, time.June, 30, 20, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("EstadoPorFechas() devolvio un error inesperado: %v", err)
	}
	if estado != domain.EstadoEnCurso {
		t.Errorf("Estado = %q, se esperaba %q", estado, domain.EstadoEnCurso)
	}
}
