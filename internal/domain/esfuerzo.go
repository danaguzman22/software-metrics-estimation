package domain

import (
	"errors"
	"math"
	"time"
)

var (
	ErrIntegranteNoPerteneceAlProyecto = errors.New("el integrante no pertenece al proyecto")
	ErrHorasEsfuerzoInvalidas          = errors.New("las horas de esfuerzo son inválidas")
	ErrFechaEsfuerzoFutura             = errors.New("la fecha del esfuerzo es futura")
	ErrHorasEstimadasInvalidas         = errors.New("las horas estimadas son inválidas")
	ErrMetricasSprintNoDisponibles     = errors.New("las métricas históricas del sprint no están disponibles")
)

// RegistroEsfuerzo representa las horas trabajadas por un integrante en una
// historia.
type RegistroEsfuerzo struct {
	Integrante string
	Fecha      time.Time
	Actividad  string
	Horas      float64
}

// MetricasEsfuerzo agrupa estimaciones, horas registradas y desviación.
type MetricasEsfuerzo struct {
	HorasEstimadas        float64
	EstimacionesCompletas bool
	HorasReales           float64
	Desviacion            *float64
}

// ItemSprintEsfuerzo conserva la estimación de una historia en un Sprint cerrado.
type ItemSprintEsfuerzo struct {
	ItemNumero     int
	HorasEstimadas *float64
}

// ValidarHorasEsfuerzo acepta valores de (0, 24] con hasta dos decimales.
func ValidarHorasEsfuerzo(horas float64) error {
	if math.IsNaN(horas) || math.IsInf(horas, 0) || horas <= 0 || horas > 24 {
		return ErrHorasEsfuerzoInvalidas
	}
	centecimas := horas * 100
	if math.Abs(centecimas-math.Round(centecimas)) >= 1e-9 {
		return ErrHorasEsfuerzoInvalidas
	}
	return nil
}

// ValidarFechaEsfuerzo impide registrar esfuerzo en un día posterior a hoy.
func ValidarFechaEsfuerzo(fecha, hoy time.Time) error {
	if soloDia(fecha).After(soloDia(hoy)) {
		return ErrFechaEsfuerzoFutura
	}
	return nil
}

// ValidarHorasEstimadas acepta estimaciones finitas no negativas con dos
// decimales como máximo.
func ValidarHorasEstimadas(horas float64) error {
	if math.IsNaN(horas) || math.IsInf(horas, 0) || horas < 0 {
		return ErrHorasEstimadasInvalidas
	}
	centecimas := horas * 100
	if math.IsInf(centecimas, 0) || math.Abs(centecimas-math.Round(centecimas)) >= 1e-9 {
		return ErrHorasEstimadasInvalidas
	}
	return nil
}

// CalcularDesviacion devuelve el porcentaje de diferencia o nil si no está
// definido porque la estimación es cero.
func CalcularDesviacion(estimadas, reales float64) *float64 {
	if estimadas == 0 {
		return nil
	}
	desviacion := (reales - estimadas) / estimadas * 100
	return &desviacion
}
