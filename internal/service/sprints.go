package service

import (
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// RepositorioSprints persiste los Sprints por proyecto (D8).
type RepositorioSprints interface {
	// Agregar registra el Sprint en el proyecto.
	// Devuelve domain.ErrSprintDuplicado si ya existe un Sprint con ese número
	// en el proyecto y domain.ErrSprintActivoExistente si ya hay un Sprint
	// activo (defensa ante concurrencia, R8).
	Agregar(proyectoID string, sprint domain.Sprint) error

	// ExisteNumero informa si el proyecto ya tiene un Sprint con ese número.
	ExisteNumero(proyectoID string, numero int) (bool, error)

	// ExisteActivo informa si el proyecto tiene un Sprint en estado Activo.
	ExisteActivo(proyectoID string) (bool, error)
}

// ServicioSprints orquesta los casos de uso de Sprints sobre el dominio (D8).
type ServicioSprints struct {
	proyectos RepositorioProyectos
	sprints   RepositorioSprints
}

// NuevoServicioSprints crea el servicio con los repositorios indicados.
func NuevoServicioSprints(proyectos RepositorioProyectos, sprints RepositorioSprints) *ServicioSprints {
	return &ServicioSprints{proyectos: proyectos, sprints: sprints}
}

// CrearSprint crea un Sprint en el proyecto identificado por proyectoID.
//
// Orden de evaluación: proyecto, número, objetivo, fechas, Sprint activo y
// número duplicado. Devuelve el Sprint creado con estado Activo.
func (s *ServicioSprints) CrearSprint(
	proyectoID string,
	numero int,
	objetivo string,
	fechaInicio, fechaFin time.Time,
) (domain.Sprint, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return domain.Sprint{}, envolverErrorRepositorio("verificar proyecto", err)
	}
	if !existe {
		return domain.Sprint{}, domain.ErrProyectoNoEncontrado
	}

	sprint, err := domain.NuevoSprint(numero, objetivo, fechaInicio, fechaFin)
	if err != nil {
		return domain.Sprint{}, err
	}

	activo, err := s.sprints.ExisteActivo(proyectoID)
	if err != nil {
		return domain.Sprint{}, envolverErrorRepositorio("verificar sprint activo", err)
	}
	if activo {
		return domain.Sprint{}, domain.ErrSprintActivoExistente
	}

	existeNumero, err := s.sprints.ExisteNumero(proyectoID, numero)
	if err != nil {
		return domain.Sprint{}, envolverErrorRepositorio("verificar número de sprint", err)
	}
	if existeNumero {
		return domain.Sprint{}, domain.ErrSprintDuplicado
	}

	if err := s.sprints.Agregar(proyectoID, sprint); err != nil {
		return domain.Sprint{}, envolverErrorRepositorio("agregar sprint", err)
	}

	return sprint, nil
}
