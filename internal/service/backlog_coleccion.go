package service

import (
	"fmt"
	"strings"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// almacen es el contrato mínimo de persistencia que necesita ServicioBacklog
// (compatibilidad Fase 2). No debe usarse en ServicioProyectos.
type almacen interface {
	Save(name string, v any) error
	Load(name string, v any) error
}

func esNoEncontradoStorage(err error) bool {
	return err != nil && strings.Contains(err.Error(), "archivo no encontrado")
}

// verificarProyectoExistente carga la colección "proyectos" y comprueba que
// exista el proyecto identificado por proyectoID.
func verificarProyectoExistente(a almacen, proyectoID string) error {
	var proyectos []ProyectoRegistrado
	if err := a.Load("proyectos", &proyectos); err != nil {
		if !esNoEncontradoStorage(err) {
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
func cargarColeccionBacklog(a almacen) ([]itemRegistrado, error) {
	var backlog []itemRegistrado
	if err := a.Load("backlog", &backlog); err != nil {
		if !esNoEncontradoStorage(err) {
			return nil, fmt.Errorf("service: cargar backlog: %w", err)
		}
		return nil, nil
	}
	return backlog, nil
}
