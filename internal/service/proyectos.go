package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage"
)

// almacen es el contrato mínimo de persistencia que necesita el servicio
// (D16). Lo satisface *storage.JSONStore y cualquier doble con Save y Load.
type almacen interface {
	Save(name string, v any) error
	Load(name string, v any) error
}

// proyectoRegistrado es el registro persistido de un proyecto: el ID junto al
// proyecto del dominio (D13). El ID no forma parte de domain.Proyecto (D11).
type proyectoRegistrado struct {
	ID       string          `json:"id"`
	Proyecto domain.Proyecto `json:"proyecto"`
}

// ServicioProyectos orquesta los casos de uso de proyectos sobre el dominio.
type ServicioProyectos struct {
	almacen almacen
}

// NuevoServicioProyectos crea el servicio con el almacenamiento indicado.
func NuevoServicioProyectos(almacen almacen) *ServicioProyectos {
	return &ServicioProyectos{almacen: almacen}
}

// ModificarProyecto modifica el nombre, la fecha de inicio y la fecha de fin
// del proyecto identificado por id, conservando su estado actual (D17).
//
// Carga la colección "proyectos" (D13); un archivo inexistente se interpreta
// como colección vacía (D14). Busca el registro por ID: si no existe devuelve
// domain.ErrProyectoNoEncontrado, antes de validar nombre y fechas (D4, D12).
// Delega las validaciones y la copia en domain.ModificarProyecto; si falla,
// devuelve ese error sin escribir nada. En caso de éxito reemplaza únicamente
// el registro encontrado y guarda la colección completa (D13).
//
// En caso de error en las validaciones o la carga, devuelve
// domain.Proyecto{} sin invocar Save. Si la persistencia falla durante Save,
// también devuelve domain.Proyecto{} junto con el error envuelto.
// Los errores de almacenamiento distintos de "no encontrado" se propagan
// envueltos con %w y no se mapean a domain.ErrProyectoNoEncontrado (D15).
func (s *ServicioProyectos) ModificarProyecto(
	id, nombre string,
	fechaInicio, fechaFin time.Time,
) (domain.Proyecto, error) {
	var coleccion []proyectoRegistrado
	if err := s.almacen.Load("proyectos", &coleccion); err != nil {
		if !errors.Is(err, storage.ErrNotFound) {
			return domain.Proyecto{}, fmt.Errorf("service: cargar proyectos: %w", err)
		}
		// D14: un archivo inexistente se interpreta como colección vacía.
		coleccion = nil
	}

	indice := -1
	for i, registro := range coleccion {
		if registro.ID == id {
			indice = i
			break
		}
	}
	if indice == -1 {
		return domain.Proyecto{}, domain.ErrProyectoNoEncontrado
	}

	modificado, err := domain.ModificarProyecto(coleccion[indice].Proyecto, nombre, fechaInicio, fechaFin)
	if err != nil {
		return domain.Proyecto{}, err
	}

	coleccion[indice].Proyecto = modificado
	if err := s.almacen.Save("proyectos", coleccion); err != nil {
		return domain.Proyecto{}, fmt.Errorf("service: guardar proyectos: %w", err)
	}

	return modificado, nil
}
