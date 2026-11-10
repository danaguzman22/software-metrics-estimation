package domain

import (
	"errors"
	"strings"
)

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
// Recorta los espacios externos del título y la prioridad. Devuelve
// ErrTituloItemVacio si el título queda vacío y ErrPrioridadInvalida si la
// prioridad no es exactamente "Alta", "Media" o "Baja" (D4, D5). El estado
// inicial es siempre "Pendiente" (D6) y los story points quedan en 0 (D7). La
// descripción y los criterios se guardan tal como se reciben (D8).
func NuevoItemBacklog(titulo, descripcion, prioridad string, criterios []string) (ItemBacklog, error) {
	titulo = strings.TrimSpace(titulo)
	if titulo == "" {
		return ItemBacklog{}, ErrTituloItemVacio
	}

	prioridad = strings.TrimSpace(prioridad)
	switch prioridad {
	case "Alta", "Media", "Baja":
	default:
		return ItemBacklog{}, ErrPrioridadInvalida
	}

	return ItemBacklog{
		Titulo:      titulo,
		Descripcion: descripcion,
		Prioridad:   prioridad,
		Estado:      "Pendiente",
		StoryPoints: 0,
		Criterios:   criterios,
	}, nil
}
