package domain

import "errors"

// Errores de dominio de ítems del backlog (D10).
var (
	// ErrTituloItemVacio indica que el título del ítem está vacío o compuesto
	// únicamente por espacios.
	ErrTituloItemVacio = errors.New("el título del ítem no puede estar vacío")
	// ErrPrioridadInvalida indica que la prioridad del ítem no es Alta, Media
	// ni Baja.
	ErrPrioridadInvalida = errors.New("la prioridad del ítem no es válida: debe ser Alta, Media o Baja")
)

// ItemBacklog representa un ítem del Product Backlog de un proyecto, con los
// campos de CA1 (D10).
type ItemBacklog struct {
	ID          string
	Titulo      string
	Descripcion string
	Prioridad   string
	Estado      string
	StoryPoints int
	Criterios   []string
}

// NuevoItemBacklog crea un ítem del backlog a partir de su título, descripción,
// prioridad y criterios de aceptación (D10).
//
// Stub de la fase RED: la lógica real se implementa en la fase GREEN.
func NuevoItemBacklog(titulo, descripcion, prioridad string, criterios []string) (ItemBacklog, error) {
	return ItemBacklog{}, errors.New("no implementado")
}
