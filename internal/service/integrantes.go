package service

import (
	"errors"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// integranteRegistrado es el registro persistido de un integrante: el ID del
// proyecto junto al integrante del dominio (D3). El ID no forma parte de
// domain.Integrante.
type integranteRegistrado struct {
	ProyectoID string            `json:"proyecto_id"`
	Integrante domain.Integrante `json:"integrante"`
}

// AgregarIntegrante registra un integrante con su nombre y su rol en el
// proyecto identificado por proyectoID (D5).
//
// Stub de la fase RED de US-03: sin lógica; devuelve domain.Integrante{} y un
// error "no implementado" para que los tests fallen.
func (s *ServicioProyectos) AgregarIntegrante(
	proyectoID, nombre, rol string,
) (domain.Integrante, error) {
	return domain.Integrante{}, errors.New("no implementado")
}
