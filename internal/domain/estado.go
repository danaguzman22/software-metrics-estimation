package domain

import (
	"errors"
	"math"
	"time"
)

// Estados calculados por fechas (US-04, D2).
const (
	// EstadoEnCurso indica que "hoy" esta entre la fecha de inicio y la de fin.
	EstadoEnCurso Estado = "En curso"
	// EstadoFinalizado indica que "hoy" es posterior a la fecha de fin.
	EstadoFinalizado Estado = "Finalizado"
)

// ErrHistoriasInvalidas indica que las cantidades de historias son invalidas.
var ErrHistoriasInvalidas = errors.New("las cantidades de historias son inválidas")

// soloDia descarta la hora de una fecha y conserva año, mes y dia (D7).
func soloDia(t time.Time) time.Time {
	anio, mes, dia := t.Date()
	return time.Date(anio, mes, dia, 0, 0, 0, 0, time.UTC)
}

// EstadoPorFechas calcula el estado de un proyecto comparando "hoy" con sus
// fechas, solo por dia. No usa el estado guardado del proyecto (D1, D7).
//
// Devuelve ErrFechaFinAnteriorAInicio si la fecha de fin es anterior a la de
// inicio (D9). Las fechas cero no se validan.
func EstadoPorFechas(inicio, fin, hoy time.Time) (Estado, error) {
	diaInicio := soloDia(inicio)
	diaFin := soloDia(fin)
	diaHoy := soloDia(hoy)

	if diaFin.Before(diaInicio) {
		return "", ErrFechaFinAnteriorAInicio
	}
	if diaHoy.Before(diaInicio) {
		return EstadoPlanificado, nil
	}
	if diaHoy.After(diaFin) {
		return EstadoFinalizado, nil
	}
	return EstadoEnCurso, nil
}

// PorcentajeCompletadas calcula el porcentaje de historias completadas,
// entre 0 y 100, redondeado a 2 decimales (D4, D6).
//
// Con total 0 y completadas 0 devuelve 0 sin error. Devuelve
// ErrHistoriasInvalidas si algun valor es negativo o si completadas > total.
func PorcentajeCompletadas(total, completadas int) (float64, error) {
	if total < 0 || completadas < 0 || completadas > total {
		return 0, ErrHistoriasInvalidas
	}
	if total == 0 {
		return 0, nil
	}
	porcentaje := float64(completadas) / float64(total) * 100
	return math.Round(porcentaje*100) / 100, nil
}
