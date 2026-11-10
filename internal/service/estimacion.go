package service

import (
	"errors"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func (s *ServicioBacklog) EstimarItem(proyectoID, itemID string, storyPoints int) (domain.ItemBacklog, error) {
	return domain.ItemBacklog{}, errors.New("no implementado")
}
