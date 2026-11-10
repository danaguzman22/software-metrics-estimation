package service

import (
	"fmt"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func (s *ServicioBacklog) EstimarItem(proyectoID, itemID string, storyPoints int) (domain.ItemBacklog, error) {
	if err := verificarProyectoExistente(s.almacen, proyectoID); err != nil {
		return domain.ItemBacklog{}, err
	}

	backlog, err := cargarColeccionBacklog(s.almacen)
	if err != nil {
		return domain.ItemBacklog{}, err
	}

	indice := -1
	for i, registro := range backlog {
		if registro.ProyectoID == proyectoID && registro.Item.ID == itemID {
			indice = i
			break
		}
	}
	if indice == -1 {
		return domain.ItemBacklog{}, domain.ErrItemNoEncontrado
	}

	estimado, err := backlog[indice].Item.Estimar(storyPoints)
	if err != nil {
		return domain.ItemBacklog{}, err
	}

	backlog[indice] = itemRegistrado{ProyectoID: proyectoID, Item: estimado}
	if err := s.almacen.Save("backlog", backlog); err != nil {
		return domain.ItemBacklog{}, fmt.Errorf("service: guardar backlog: %w", err)
	}

	return estimado, nil
}
