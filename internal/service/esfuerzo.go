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
	var sprintAtribuido *int
	if asignado {
		sprintAtribuido = &sprintNumero
	}
	if err := s.esfuerzo.Agregar(proyectoID, itemNumero, sprintAtribuido, registro); err != nil {
		return domain.RegistroEsfuerzo{}, envolverErrorRepositorio("agregar esfuerzo", err)
	}
	return registro, nil
}

// EstimarHoras establece o reemplaza las horas estimadas de una historia.
func (s *ServicioEsfuerzo) EstimarHoras(
	proyectoID, itemID string,
	horas float64,
) (domain.ItemBacklog, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return domain.ItemBacklog{}, envolverErrorRepositorio("verificar proyecto", err)
	}
	if !existe {
		return domain.ItemBacklog{}, domain.ErrProyectoNoEncontrado
	}
	if err := domain.ValidarHorasEstimadas(horas); err != nil {
		return domain.ItemBacklog{}, err
	}
	item, err := s.backlog.ObtenerItem(proyectoID, itemID)
	if err != nil {
		return domain.ItemBacklog{}, envolverErrorRepositorio("obtener ítem del backlog", err)
	}
	item.HorasEstimadas = &horas
	if err := s.backlog.ActualizarItem(proyectoID, item); err != nil {
		return domain.ItemBacklog{}, envolverErrorRepositorio("actualizar horas estimadas", err)
	}
	return item, nil
}

// MetricasHistoria devuelve las métricas de una historia del proyecto.
func (s *ServicioEsfuerzo) MetricasHistoria(
	proyectoID, itemID string,
) (domain.MetricasEsfuerzo, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return domain.MetricasEsfuerzo{}, envolverErrorRepositorio("verificar proyecto", err)
	}
	if !existe {
		return domain.MetricasEsfuerzo{}, domain.ErrProyectoNoEncontrado
	}
	item, err := s.backlog.ObtenerItem(proyectoID, itemID)
	if err != nil {
		return domain.MetricasEsfuerzo{}, envolverErrorRepositorio("obtener ítem del backlog", err)
	}
	numero, ok := domain.NumeroDeItem(item.ID)
	if !ok {
		return domain.MetricasEsfuerzo{}, domain.ErrItemNoEncontrado
	}
	reales, err := s.esfuerzo.HorasRealesPorItem(proyectoID, numero)
	if err != nil {
		return domain.MetricasEsfuerzo{}, envolverErrorRepositorio("sumar horas reales de la historia", err)
	}
	metricas := domain.MetricasEsfuerzo{HorasReales: reales, EstimacionesCompletas: item.HorasEstimadas != nil}
	if item.HorasEstimadas != nil {
		metricas.HorasEstimadas = *item.HorasEstimadas
		metricas.Desviacion = domain.CalcularDesviacion(metricas.HorasEstimadas, reales)
	}
	return metricas, nil
}

// MetricasSprint devuelve las métricas del Sprint. Para Sprints cerrados usa
// exclusivamente el snapshot persistido al momento del cierre.
func (s *ServicioEsfuerzo) MetricasSprint(
	proyectoID string,
	sprintNumero int,
) (domain.MetricasEsfuerzo, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return domain.MetricasEsfuerzo{}, envolverErrorRepositorio("verificar proyecto", err)
	}
	if !existe {
		return domain.MetricasEsfuerzo{}, domain.ErrProyectoNoEncontrado
	}
	sprint, err := s.sprints.Obtener(proyectoID, sprintNumero)
	if err != nil {
		return domain.MetricasEsfuerzo{}, envolverErrorRepositorio("obtener sprint", err)
	}
	metricas := domain.MetricasEsfuerzo{EstimacionesCompletas: true}
	if sprint.Estado == domain.EstadoSprintCerrado {
		items, disponible, err := s.esfuerzo.SnapshotSprint(proyectoID, sprintNumero)
		if err != nil {
			return domain.MetricasEsfuerzo{}, envolverErrorRepositorio("obtener snapshot del sprint", err)
		}
		if !disponible {
			return domain.MetricasEsfuerzo{}, domain.ErrMetricasSprintNoDisponibles
		}
		sumarEstimacionesSprint(&metricas, items)
	} else {
		numeros, err := s.asignaciones.ListarNumerosItem(proyectoID, sprintNumero)
		if err != nil {
			return domain.MetricasEsfuerzo{}, envolverErrorRepositorio("listar historias del sprint", err)
		}
		items, err := s.backlog.ListarPorProyecto(proyectoID)
		if err != nil {
			return domain.MetricasEsfuerzo{}, envolverErrorRepositorio("listar historias del proyecto", err)
		}
		asignados := make(map[int]struct{}, len(numeros))
		for _, numero := range numeros {
			asignados[numero] = struct{}{}
		}
		for _, item := range items {
			numero, ok := domain.NumeroDeItem(item.ID)
			if !ok {
				continue
			}
			if _, asignado := asignados[numero]; asignado {
				sumarEstimacionHistoria(&metricas, item.HorasEstimadas)
			}
		}
	}
	reales, err := s.esfuerzo.HorasRealesPorSprint(proyectoID, sprintNumero)
	if err != nil {
		return domain.MetricasEsfuerzo{}, envolverErrorRepositorio("sumar horas reales del sprint", err)
	}
	metricas.HorasReales = reales
	if metricas.EstimacionesCompletas {
		metricas.Desviacion = domain.CalcularDesviacion(metricas.HorasEstimadas, reales)
	}
	return metricas, nil
}

// MetricasProyecto devuelve las métricas agregadas de todas las historias y
// registros del proyecto.
func (s *ServicioEsfuerzo) MetricasProyecto(proyectoID string) (domain.MetricasEsfuerzo, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return domain.MetricasEsfuerzo{}, envolverErrorRepositorio("verificar proyecto", err)
	}
	if !existe {
		return domain.MetricasEsfuerzo{}, domain.ErrProyectoNoEncontrado
	}
	items, err := s.backlog.ListarPorProyecto(proyectoID)
	if err != nil {
		return domain.MetricasEsfuerzo{}, envolverErrorRepositorio("listar historias del proyecto", err)
	}
	metricas := domain.MetricasEsfuerzo{EstimacionesCompletas: true}
	for _, item := range items {
		sumarEstimacionHistoria(&metricas, item.HorasEstimadas)
	}
	reales, err := s.esfuerzo.HorasRealesPorProyecto(proyectoID)
	if err != nil {
		return domain.MetricasEsfuerzo{}, envolverErrorRepositorio("sumar horas reales del proyecto", err)
	}
	metricas.HorasReales = reales
	if metricas.EstimacionesCompletas {
		metricas.Desviacion = domain.CalcularDesviacion(metricas.HorasEstimadas, reales)
	}
	return metricas, nil
}

func sumarEstimacionesSprint(metricas *domain.MetricasEsfuerzo, items []domain.ItemSprintEsfuerzo) {
	for _, item := range items {
		sumarEstimacionHistoria(metricas, item.HorasEstimadas)
	}
}

func sumarEstimacionHistoria(metricas *domain.MetricasEsfuerzo, estimadas *float64) {
	if estimadas == nil {
		metricas.EstimacionesCompletas = false
		return
	}
	metricas.HorasEstimadas += *estimadas
}
