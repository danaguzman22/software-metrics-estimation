package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage"
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
// Carga la colección "proyectos"; un archivo inexistente se interpreta como
// colección vacía (D2). Busca el proyecto por ID: si no existe devuelve
// domain.ErrProyectoNoEncontrado antes de validar título y prioridad (D1, D9).
// Delega las validaciones en domain.NuevoItemBacklog; si falla, devuelve ese
// error sin escribir nada. Carga la colección "backlog" (D2), calcula el
// siguiente ID del proyecto según R2 ("HU-01" si no hay ítems; si no, el mayor
// existente más uno), agrega el registro {ProyectoID, Item} y guarda la
// colección completa.
//
// Ante cualquier error devuelve domain.ItemBacklog{} sin invocar Save. Si la
// persistencia falla durante Save, también devuelve domain.ItemBacklog{} junto
// con el error envuelto. Los errores de almacenamiento distintos de "no
// encontrado" se propagan envueltos con %w y no se mapean a
// domain.ErrProyectoNoEncontrado (D2).
func (s *ServicioBacklog) CrearItem(
	proyectoID, titulo, descripcion, prioridad string, criterios []string,
) (domain.ItemBacklog, error) {
	if err := verificarProyectoExistente(s.almacen, proyectoID); err != nil {
		return domain.ItemBacklog{}, err
	}

	nuevo, err := domain.NuevoItemBacklog(titulo, descripcion, prioridad, criterios)
	if err != nil {
		return domain.ItemBacklog{}, err
	}

	var backlog []itemRegistrado
	if err := s.almacen.Load("backlog", &backlog); err != nil {
		if !errors.Is(err, storage.ErrNotFound) {
			return domain.ItemBacklog{}, fmt.Errorf("service: cargar backlog: %w", err)
		}
		// D2: un archivo inexistente se interpreta como backlog vacío.
		backlog = nil
	}

	maximo := 0
	for _, registro := range backlog {
		if registro.ProyectoID != proyectoID {
			continue
		}
		if numero, ok := numeroDeItem(registro.Item.ID); ok && numero > maximo {
			maximo = numero
		}
	}
	nuevo.ID = fmt.Sprintf("HU-%02d", maximo+1)

	backlog = append(backlog, itemRegistrado{
		ProyectoID: proyectoID,
		Item:       nuevo,
	})
	if err := s.almacen.Save("backlog", backlog); err != nil {
		return domain.ItemBacklog{}, fmt.Errorf("service: guardar backlog: %w", err)
	}

	return nuevo, nil
}

// numeroDeItem extrae el número de un ID con el patrón "HU-NN". Devuelve false
// si el ID no tiene ese patrón.
func numeroDeItem(id string) (int, bool) {
	const prefijo = "HU-"
	if !strings.HasPrefix(id, prefijo) {
		return 0, false
	}
	numero, err := strconv.Atoi(strings.TrimPrefix(id, prefijo))
	if err != nil {
		return 0, false
	}
	return numero, true
}
