package service

import (
	"errors"
	"fmt"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage"
)

// almacen es el contrato mínimo de persistencia que ServicioBacklog sigue usando
// en Fase 1/Fase 2 (compatibilidad temporal). Se mantiene aquí para que el
// código productivo del paquete service no dependa del paquete de storage; la
// migración funcional de ServicioBacklog a repositorios ocurre en Fase 2.
type almacen interface {
	Save(name string, v any) error
	Load(name string, v any) error
}

// verificarProyectoExistente carga la colección "proyectos" y comprueba que
// exista el proyecto identificado por proyectoID. Conserva el comportamiento
// actual del backlog en Fase 1 (compatibilidad temporal).
func verificarProyectoExistente(a almacen, proyectoID string) error {
	var proyectos []ProyectoRegistrado
	if err := a.Load("proyectos", &proyectos); err != nil {
		if !errors.Is(err, storage.ErrNotFound) {
			return fmt.Errorf("service: cargar proyectos: %w", err)
		}
		proyectos = nil
	}

	for _, registro := range proyectos {
		if registro.ID == proyectoID {
			return nil
		}
	}
	return domain.ErrProyectoNoEncontrado
}

// cargarColeccionBacklog carga la colección "backlog" desde el almacenamiento.
// Un archivo inexistente se interpreta como colección vacía (comportamiento actual).
func cargarColeccionBacklog(a almacen) ([]itemRegistrado, error) {
	var backlog []itemRegistrado
	if err := a.Load("backlog", &backlog); err != nil {
		if !errors.Is(err, storage.ErrNotFound) {
			return nil, fmt.Errorf("service: cargar backlog: %w", err)
		}
		return nil, nil
	}
	return backlog, nil
}
