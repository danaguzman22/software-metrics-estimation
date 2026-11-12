package service

import (
	"errors"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// EstimarItem asigna los story points al ítem identificado por itemID dentro
// del proyecto identificado por proyectoID.
//
// Fase RED de T-09: la firma pública se conserva y la implementación se
// completa en la fase GREEN.
func (s *ServicioBacklog) EstimarItem(proyectoID, itemID string, storyPoints int) (domain.ItemBacklog, error) {
	return domain.ItemBacklog{}, errors.New("no implementado")
}
