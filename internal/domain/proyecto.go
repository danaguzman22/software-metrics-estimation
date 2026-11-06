package domain

import "time"

// Estado representa el estado de un proyecto.
type Estado string

// EstadoPlanificado es el estado inicial de todo proyecto creado.
const EstadoPlanificado Estado = "Planificado"

// Proyecto representa un proyecto con nombre, fechas y estado.
type Proyecto struct {
	Nombre      string
	FechaInicio time.Time
	FechaFin    time.Time
	Estado      Estado
}

// NuevoProyecto crea un proyecto a partir de su nombre y sus fechas.
//
// RED: firma mínima para que los tests compilen. La implementación real
// corresponde a la fase GREEN.
func NuevoProyecto(nombre string, fechaInicio, fechaFin time.Time) (Proyecto, error) {
	panic("no implementado")
}
