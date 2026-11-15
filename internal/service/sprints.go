package service

import (
	"sort"
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

	// Obtener devuelve el Sprint del proyecto o domain.ErrSprintNoEncontrado.
	Obtener(proyectoID string, numero int) (domain.Sprint, error)

	// ListarPorProyecto devuelve los Sprints del proyecto ordenados por número.
	ListarPorProyecto(proyectoID string) ([]domain.Sprint, error)

	// Cerrar cierra un Sprint activo y elimina las asignaciones no completadas
	// en una transacción. Si el Sprint no está activo devuelve
	// domain.ErrSprintCerrado.
	Cerrar(proyectoID string, numero, spPlanificados int) error
}

// ServicioSprints orquesta los casos de uso de Sprints sobre el dominio (D8).
type ServicioSprints struct {
	proyectos    RepositorioProyectos
	sprints      RepositorioSprints
	backlog      RepositorioBacklog
	asignaciones RepositorioAsignaciones
}

// NuevoServicioSprints crea el servicio con los repositorios indicados.
func NuevoServicioSprints(
	proyectos RepositorioProyectos,
	sprints RepositorioSprints,
	backlog RepositorioBacklog,
	asignaciones RepositorioAsignaciones,
) *ServicioSprints {
	return &ServicioSprints{
		proyectos:    proyectos,
		sprints:      sprints,
		backlog:      backlog,
		asignaciones: asignaciones,
	}
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

// ListarSprints devuelve los Sprints del proyecto ordenados por número.
//
// Verifica primero que el proyecto exista y devuelve
// domain.ErrProyectoNoEncontrado sin envolver si no existe. Reutiliza
// RepositorioSprints.ListarPorProyecto y no recalcula Story Points.
func (s *ServicioSprints) ListarSprints(proyectoID string) ([]domain.Sprint, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return nil, envolverErrorRepositorio("verificar proyecto", err)
	}
	if !existe {
		return nil, domain.ErrProyectoNoEncontrado
	}

	sprints, err := s.sprints.ListarPorProyecto(proyectoID)
	if err != nil {
		return nil, envolverErrorRepositorio("listar sprints del proyecto", err)
	}

	return sprints, nil
}

// CerrarSprint cierra el Sprint y conserva el total de Story Points
// planificados antes de retirar las asignaciones no completadas.
func (s *ServicioSprints) CerrarSprint(proyectoID string, sprintNumero int) error {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return envolverErrorRepositorio("verificar proyecto", err)
	}
	if !existe {
		return domain.ErrProyectoNoEncontrado
	}

	sprint, err := s.sprints.Obtener(proyectoID, sprintNumero)
	if err != nil {
		return envolverErrorRepositorio("obtener sprint", err)
	}
	if sprint.Estado != domain.EstadoSprintActivo {
		return domain.ErrSprintCerrado
	}

	spPlanificados, err := s.StoryPointsPlanificados(proyectoID, sprintNumero)
	if err != nil {
		return err
	}
	if err := s.sprints.Cerrar(proyectoID, sprintNumero, spPlanificados); err != nil {
		return envolverErrorRepositorio("cerrar sprint", err)
	}
	return nil
}

// ListarItemsDelSprint devuelve las historias del backlog asignadas al Sprint,
// ordenadas por número de historia ascendente.
//
// Verifica primero que el proyecto exista y devuelve
// domain.ErrProyectoNoEncontrado sin envolver si no existe. Devuelve
// domain.ErrSprintNoEncontrado sin envolver si el Sprint no existe. Para un
// Sprint sin historias devuelve un slice vacío y nil.
func (s *ServicioSprints) ListarItemsDelSprint(proyectoID string, sprintNumero int) ([]domain.ItemBacklog, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return nil, envolverErrorRepositorio("verificar proyecto", err)
	}
	if !existe {
		return nil, domain.ErrProyectoNoEncontrado
	}

	if _, err := s.sprints.Obtener(proyectoID, sprintNumero); err != nil {
		return nil, envolverErrorRepositorio("obtener sprint", err)
	}

	numeros, err := s.asignaciones.ListarNumerosItem(proyectoID, sprintNumero)
	if err != nil {
		return nil, envolverErrorRepositorio("listar ítems del sprint", err)
	}

	items, err := s.backlog.ListarPorProyecto(proyectoID)
	if err != nil {
		return nil, envolverErrorRepositorio("listar ítems del proyecto", err)
	}

	porNumero := make(map[int]domain.ItemBacklog, len(items))
	for _, item := range items {
		numero, ok := domain.NumeroDeItem(item.ID)
		if !ok {
			continue
		}
		porNumero[numero] = item
	}

	ordenados := append([]int(nil), numeros...)
	sort.Ints(ordenados)

	asignados := make([]domain.ItemBacklog, 0, len(ordenados))
	for _, numero := range ordenados {
		item, ok := porNumero[numero]
		if !ok {
			continue
		}
		asignados = append(asignados, item)
	}
	return asignados, nil
}
