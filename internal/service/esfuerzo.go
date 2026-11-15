package service

import (
	"strings"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// ServicioEsfuerzo registra el esfuerzo realizado en historias de un proyecto.
type ServicioEsfuerzo struct {
	proyectos    RepositorioProyectos
	backlog      RepositorioBacklog
	integrantes  RepositorioIntegrantes
	sprints      RepositorioSprints
	asignaciones RepositorioAsignaciones
	esfuerzo     RepositorioEsfuerzo
}

// NuevoServicioEsfuerzo crea el servicio con sus repositorios.
func NuevoServicioEsfuerzo(
	proyectos RepositorioProyectos,
	backlog RepositorioBacklog,
	integrantes RepositorioIntegrantes,
	sprints RepositorioSprints,
	asignaciones RepositorioAsignaciones,
	esfuerzo RepositorioEsfuerzo,
) *ServicioEsfuerzo {
	return &ServicioEsfuerzo{
		proyectos:    proyectos,
		backlog:      backlog,
		integrantes:  integrantes,
		sprints:      sprints,
		asignaciones: asignaciones,
		esfuerzo:     esfuerzo,
	}
}

// RegistrarEsfuerzo valida y persiste un registro de esfuerzo.
func (s *ServicioEsfuerzo) RegistrarEsfuerzo(
	proyectoID, itemID, integrante string,
	fecha time.Time,
	actividad string,
	horas float64,
	hoy time.Time,
) (domain.RegistroEsfuerzo, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return domain.RegistroEsfuerzo{}, envolverErrorRepositorio("verificar proyecto", err)
	}
	if !existe {
		return domain.RegistroEsfuerzo{}, domain.ErrProyectoNoEncontrado
	}

	if err := domain.ValidarHorasEsfuerzo(horas); err != nil {
		return domain.RegistroEsfuerzo{}, err
	}
	if err := domain.ValidarFechaEsfuerzo(fecha, hoy); err != nil {
		return domain.RegistroEsfuerzo{}, err
	}

	item, err := s.backlog.ObtenerItem(proyectoID, itemID)
	if err != nil {
		return domain.RegistroEsfuerzo{}, envolverErrorRepositorio("obtener ítem del backlog", err)
	}
	itemNumero, ok := domain.NumeroDeItem(item.ID)
	if !ok {
		return domain.RegistroEsfuerzo{}, domain.ErrItemNoEncontrado
	}

	sprintNumero, asignado, err := s.asignaciones.SprintDeItem(proyectoID, itemNumero)
	if err != nil {
		return domain.RegistroEsfuerzo{}, envolverErrorRepositorio("obtener sprint de la historia", err)
	}
	if asignado {
		sprint, err := s.sprints.Obtener(proyectoID, sprintNumero)
		if err != nil {
			return domain.RegistroEsfuerzo{}, envolverErrorRepositorio("obtener sprint", err)
		}
		if sprint.Estado == domain.EstadoSprintCerrado {
			return domain.RegistroEsfuerzo{}, domain.ErrSprintCerrado
		}
	}

	integrantes, err := s.integrantes.ListarPorProyecto(proyectoID)
	if err != nil {
		return domain.RegistroEsfuerzo{}, envolverErrorRepositorio("listar integrantes", err)
	}
	integranteRegistrado, encontrado := domain.BuscarIntegrante(integrantes, integrante)
	if !encontrado {
		return domain.RegistroEsfuerzo{}, domain.ErrIntegranteNoPerteneceAlProyecto
	}

	registro := domain.RegistroEsfuerzo{
		Integrante: integranteRegistrado.Nombre,
		Fecha:      fecha,
		Actividad:  strings.TrimSpace(actividad),
		Horas:      horas,
	}
	if err := s.esfuerzo.Agregar(proyectoID, itemNumero, registro); err != nil {
		return domain.RegistroEsfuerzo{}, envolverErrorRepositorio("agregar esfuerzo", err)
	}
	return registro, nil
}
