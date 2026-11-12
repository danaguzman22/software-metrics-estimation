package service

import (
	"errors"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// itemRegistrado es el registro persistido de un ítem del backlog: el ID del
// proyecto junto al ítem del dominio (D2).
//
// Se conserva temporalmente durante la Fase 2 para el armado de los tests de
// servicio; se elimina en la Fase 3 de limpieza.
type itemRegistrado struct {
	ProyectoID string             `json:"proyecto_id"`
	Item       domain.ItemBacklog `json:"item"`
}

// ServicioBacklog orquesta los casos de uso del Product Backlog sobre el
// dominio (D11).
type ServicioBacklog struct {
	proyectos RepositorioProyectos
	backlog   RepositorioBacklog
}

// NuevoServicioBacklog crea el servicio con los repositorios indicados.
func NuevoServicioBacklog(proyectos RepositorioProyectos, backlog RepositorioBacklog) *ServicioBacklog {
	return &ServicioBacklog{proyectos: proyectos, backlog: backlog}
}

// CrearItem registra un ítem del backlog en el proyecto identificado por
// proyectoID (D11).
//
// Fase RED de T-09: la firma pública se conserva y la implementación se
// completa en la fase GREEN.
func (s *ServicioBacklog) CrearItem(
	proyectoID, titulo, descripcion, prioridad string, criterios []string,
) (domain.ItemBacklog, error) {
	return domain.ItemBacklog{}, errors.New("no implementado")
}
