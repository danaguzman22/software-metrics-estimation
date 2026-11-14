package domain

import "errors"

// EstadoItemPendiente es el único estado persistido de una historia hasta
// US-09. Una historia es asignable a un Sprint apenas está pendiente (CA1).
const EstadoItemPendiente = "Pendiente"

// Errores de dominio de la asignación de historias a Sprints (D1, D2, CA1).
var (
	// ErrSprintNoEncontrado indica que el Sprint no existe en el proyecto.
	ErrSprintNoEncontrado = errors.New("el sprint no existe en el proyecto")
	// ErrItemYaAsignado indica que la historia ya pertenece a un Sprint.
	ErrItemYaAsignado = errors.New("el ítem ya está asignado a un sprint")
	// ErrItemNoEstimado indica que la historia no tiene story points válidos.
	ErrItemNoEstimado = errors.New("el ítem no está estimado")
	// ErrItemNoPendiente indica que la historia no está en estado Pendiente.
	ErrItemNoPendiente = errors.New("el ítem no está pendiente")
)

// ValidarItemAsignable valida que una historia pueda asignarse a un Sprint
// (CA1): debe estar estimada con story points válidos y en estado Pendiente.
// Devuelve ErrItemNoEstimado y luego ErrItemNoPendiente, en ese orden.
func ValidarItemAsignable(item ItemBacklog) error {
	if !EsStoryPointsValido(item.StoryPoints) {
		return ErrItemNoEstimado
	}
	if item.Estado != EstadoItemPendiente {
		return ErrItemNoPendiente
	}
	return nil
}

// StoryPointsPlanificados suma los story points de las historias asignadas a
// un Sprint (CA3). No valida la estimación: las historias asignadas ya están
// estimadas por ValidarItemAsignable. Con una lista vacía devuelve 0.
func StoryPointsPlanificados(items []ItemBacklog) int {
	total := 0
	for _, item := range items {
		total += item.StoryPoints
	}
	return total
}
