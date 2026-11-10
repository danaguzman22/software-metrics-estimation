package domain

import (
	"errors"
	"strings"
	"time"
)

// Estado representa el estado de un proyecto.
type Estado string

// EstadoPlanificado es el estado inicial de todo proyecto creado.
const EstadoPlanificado Estado = "Planificado"

// Errores de validación de un proyecto.
var (
	// ErrNombreVacio indica que el nombre del proyecto está vacío o compuesto
	// únicamente por espacios.
	ErrNombreVacio = errors.New("el nombre del proyecto no puede estar vacío")
	// ErrFechaFinAnteriorAInicio indica que la fecha de fin es anterior a la de inicio.
	ErrFechaFinAnteriorAInicio = errors.New("la fecha de fin no puede ser anterior a la fecha de inicio")
	// ErrProyectoNoEncontrado indica que no existe un proyecto con el ID indicado.
	ErrProyectoNoEncontrado = errors.New("no se encontró el proyecto")
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
// Devuelve un error de validación si el nombre está vacío o compuesto
// únicamente por espacios, o si la fecha de fin es anterior a la fecha de
// inicio. Cuando el nombre es válido se conserva tal como se ingresó, sin
// recortar espacios externos. Todo proyecto creado queda en estado
// Planificado.
//
// Limitación: el tratamiento de fechas cero o faltantes queda fuera del
// alcance de US-01. La función recibe time.Time y no valida fechas cero.
func NuevoProyecto(nombre string, fechaInicio, fechaFin time.Time) (Proyecto, error) {
	if strings.TrimSpace(nombre) == "" {
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

// ModificarProyecto devuelve una copia del proyecto existente con el nombre y
// las fechas indicados, conservando su estado actual.
//
// Reutiliza las validaciones de NuevoProyecto: el nombre no puede estar vacío
// ni compuesto únicamente por espacios, y la fecha de fin no puede ser anterior
// a la de inicio. El nombre se valida antes que las fechas. Si el nombre es
// válido se conserva tal como se ingresó, sin recortar espacios externos.
//
// Ante cualquier error se devuelve Proyecto{} y el proyecto recibido no se
// muta. El estado nunca se modifica.
//
// Limitación: el tratamiento de fechas cero o faltantes queda fuera del
// alcance de US-02. La función recibe time.Time y no valida fechas cero.
func ModificarProyecto(
	proyecto Proyecto,
	nombre string,
	fechaInicio, fechaFin time.Time,
) (Proyecto, error) {
	modificado, err := NuevoProyecto(nombre, fechaInicio, fechaFin)
	if err != nil {
		return Proyecto{}, err
	}

	modificado.Estado = proyecto.Estado
	return modificado, nil
}
