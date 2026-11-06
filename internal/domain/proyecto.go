package domain

import (
	"errors"
	"time"
)

// Estado representa el estado de un proyecto.
type Estado string

// EstadoPlanificado es el estado inicial de todo proyecto creado.
const EstadoPlanificado Estado = "Planificado"

// Errores de validación de un proyecto.
var (
	// ErrNombreVacio indica que el nombre del proyecto está vacío.
	ErrNombreVacio = errors.New("el nombre del proyecto no puede estar vacío")
	// ErrFechaFinAnteriorAInicio indica que la fecha de fin es anterior a la de inicio.
	ErrFechaFinAnteriorAInicio = errors.New("la fecha de fin no puede ser anterior a la fecha de inicio")
)

// Proyecto representa un proyecto con nombre, fechas y estado.
type Proyecto struct {
	Nombre      string
	FechaInicio time.Time
	FechaFin    time.Time
	Estado      Estado
}

// NuevoProyecto crea un proyecto a partir de su nombre y sus fechas.
//
// Devuelve un error de validación si el nombre está vacío o si la fecha de fin
// es anterior a la fecha de inicio. Todo proyecto creado queda en estado
// Planificado.
func NuevoProyecto(nombre string, fechaInicio, fechaFin time.Time) (Proyecto, error) {
	if nombre == "" {
		return Proyecto{}, ErrNombreVacio
	}
	if fechaFin.Before(fechaInicio) {
		return Proyecto{}, ErrFechaFinAnteriorAInicio
	}

	return Proyecto{
		Nombre:      nombre,
		FechaInicio: fechaInicio,
		FechaFin:    fechaFin,
		Estado:      EstadoPlanificado,
	}, nil
}
