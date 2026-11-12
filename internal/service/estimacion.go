package service

import (
	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// EstimarItem asigna los story points al ítem identificado por itemID dentro
// del proyecto identificado por proyectoID.
func (s *ServicioBacklog) EstimarItem(proyectoID, itemID string, storyPoints int) (domain.ItemBacklog, error) {
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

	estimado, err := actual.Estimar(storyPoints)
	if err != nil {
		return domain.ItemBacklog{}, err
	}

	if err := s.backlog.ActualizarItem(proyectoID, estimado); err != nil {
		return domain.ItemBacklog{}, envolverErrorRepositorio("actualizar ítem del backlog", err)
	}

	return estimado, nil
}
