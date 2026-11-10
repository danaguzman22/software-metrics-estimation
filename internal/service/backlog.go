package service

import (
	"errors"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// itemRegistrado es el registro persistido de un ítem del backlog: el ID del
// proyecto junto al ítem del dominio (D2).
type itemRegistrado struct {
	ProyectoID string             `json:"proyecto_id"`
	Item       domain.ItemBacklog `json:"item"`
}

// ServicioBacklog orquesta los casos de uso del Product Backlog sobre el
// dominio (D11).
type ServicioBacklog struct {
	almacen almacen
}

// NuevoServicioBacklog crea el servicio con el almacenamiento indicado.
func NuevoServicioBacklog(almacen almacen) *ServicioBacklog {
	return &ServicioBacklog{almacen: almacen}
}

// CrearItem registra un ítem del backlog en el proyecto identificado por
// proyectoID (D11).
//
// Stub de la fase RED: la lógica real se implementa en la fase GREEN.
func (s *ServicioBacklog) CrearItem(
	proyectoID, titulo, descripcion, prioridad string, criterios []string,
) (domain.ItemBacklog, error) {
	return domain.ItemBacklog{}, errors.New("no implementado")
}
