package service

import (
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
func (s *ServicioBacklog) CrearItem(
	proyectoID, titulo, descripcion, prioridad string, criterios []string,
) (domain.ItemBacklog, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return domain.ItemBacklog{}, envolverErrorRepositorio("verificar proyecto", err)
	}
	if !existe {
		return domain.ItemBacklog{}, domain.ErrProyectoNoEncontrado
	}

	item, err := domain.NuevoItemBacklog(titulo, descripcion, prioridad, criterios)
	if err != nil {
		return domain.ItemBacklog{}, err
	}

	creado, err := s.backlog.AgregarItem(proyectoID, item)
	if err != nil {
		return domain.ItemBacklog{}, envolverErrorRepositorio("agregar ítem al backlog", err)
	}

	return creado, nil
}
