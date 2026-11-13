package domain

import (
	"errors"
	"strings"
	"time"
)

// EstadoSprintActivo es el estado inicial de todo Sprint creado (D2).
const EstadoSprintActivo = "Activo"

// Errores de dominio de Sprints (D3, D4, D5, D6).
var (
	// ErrNumeroSprintInvalido indica que el número del Sprint no es mayor que cero.
	ErrNumeroSprintInvalido = errors.New("el número del sprint debe ser mayor que cero")
	// ErrObjetivoSprintVacio indica que el objetivo del Sprint está vacío o
	// compuesto únicamente por espacios.
	ErrObjetivoSprintVacio = errors.New("el objetivo del sprint no puede estar vacío")
	// ErrFechasSprintInvalidas indica que la fecha de fin es anterior a la de inicio.
	ErrFechasSprintInvalidas = errors.New("la fecha de fin del sprint no puede ser anterior a la fecha de inicio")
	// ErrSprintActivoExistente indica que el proyecto ya tiene un Sprint activo.
	ErrSprintActivoExistente = errors.New("el proyecto ya tiene un sprint activo")
	// ErrSprintDuplicado indica que el número del Sprint ya existe en el proyecto.
	ErrSprintDuplicado = errors.New("el número del sprint ya existe en el proyecto")
)

// Sprint representa una iteración de un proyecto con número, objetivo, fechas
// y estado (D2).
type Sprint struct {
	Numero      int
	Objetivo    string
	FechaInicio time.Time
	FechaFin    time.Time
	Estado      string
}

// NuevoSprint crea un Sprint a partir de su número, su objetivo y sus fechas.
//
// Recorta los espacios externos del objetivo. Valida en este orden: número
// mayor que cero, objetivo no vacío y fecha de fin no anterior a la fecha de
// inicio. Las fechas se comparan por año, mes y día con la misma forma que usa
// Proyecto (soloDia); fechas iguales son válidas. El estado inicial es siempre
// EstadoSprintActivo. Ante cualquier error devuelve Sprint{} y no persiste.
func NuevoSprint(numero int, objetivo string, fechaInicio, fechaFin time.Time) (Sprint, error) {
	if numero <= 0 {
		return Sprint{}, ErrNumeroSprintInvalido
	}

	objetivo = strings.TrimSpace(objetivo)
	if objetivo == "" {
		return Sprint{}, ErrObjetivoSprintVacio
	}

	if soloDia(fechaFin).Before(soloDia(fechaInicio)) {
		return Sprint{}, ErrFechasSprintInvalidas
	}

	return Sprint{
		Numero:      numero,
		Objetivo:    objetivo,
		FechaInicio: fechaInicio,
		FechaFin:    fechaFin,
		Estado:      EstadoSprintActivo,
	}, nil
}
