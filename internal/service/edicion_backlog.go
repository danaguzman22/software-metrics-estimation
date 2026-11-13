package service

import (
	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// EditarItem edita los campos editables del ítem identificado por itemID
// dentro del proyecto identificado por proyectoID. Conserva el ID y el estado.
func (s *ServicioBacklog) EditarItem(
	proyectoID, itemID string,
	titulo, descripcion, prioridad string,
	criterios []string,
	storyPoints int,
) (domain.ItemBacklog, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return domain.ItemBacklog{}, envolverErrorRepositorio("verificar proyecto", err)
	}
	if !existe {
		return domain.ItemBacklog{}, domain.ErrProyectoNoEncontrado
	}

	actual, err := s.backlog.ObtenerItem(proyectoID, itemID)
	if err != nil {
		return domain.ItemBacklog{}, envolverErrorRepositorio("obtener ítem del backlog", err)
	}

	editado, err := actual.Editar(titulo, descripcion, prioridad, criterios, storyPoints)
	if err != nil {
		return domain.ItemBacklog{}, err
	}

	if err := s.backlog.ActualizarItem(proyectoID, editado); err != nil {
		return domain.ItemBacklog{}, envolverErrorRepositorio("actualizar ítem del backlog", err)
	}

	return editado, nil
}

// ListarBacklog devuelve los ítems del proyecto ordenados por prioridad
// Alta, Media, Baja; las prioridades desconocidas quedan después de Baja. A
// igual prioridad, por número ascendente; los IDs malformados quedan después
// de los válidos y entre sí por ID como texto ascendente.
func (s *ServicioBacklog) ListarBacklog(proyectoID string) ([]domain.ItemBacklog, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return nil, envolverErrorRepositorio("verificar proyecto", err)
	}
	if !existe {
		return nil, domain.ErrProyectoNoEncontrado
	}

	items, err := s.backlog.ListarPorProyecto(proyectoID)
	if err != nil {
		return nil, envolverErrorRepositorio("listar ítems del backlog", err)
	}

	return domain.OrdenarBacklog(items), nil
}
