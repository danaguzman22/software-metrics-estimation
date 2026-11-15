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
)

// RegistroEsfuerzo representa las horas trabajadas por un integrante en una
// historia.
type RegistroEsfuerzo struct {
	Integrante string
	Fecha      time.Time
	Actividad  string
	Horas      float64
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
