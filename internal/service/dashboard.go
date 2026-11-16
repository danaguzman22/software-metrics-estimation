package service

import (
	"errors"
	"fmt"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// PuntoVelocidadSprint es la velocidad (Story Points completados) de un
// Sprint cerrado del proyecto (CA2).
type PuntoVelocidadSprint struct {
	SprintNumero int
	Velocidad    int
}

// PuntoHorasSprint son las horas estimadas y reales de un Sprint cerrado.
// MetricasDisponibles es false cuando el Sprint no tiene snapshot de esfuerzo
// (cerrado antes de US-17): en ese caso HorasEstimadas y HorasReales son 0 y
// el dashboard se muestra igualmente (D12).
// Desviacion es el porcentaje de diferencia del Sprint, calculado por el
// servicio con domain.CalcularDesviacion y con la misma política de nil que
// R4 (nil si MetricasDisponibles = false, si EstimacionesCompletas = false o
// si la suma de estimaciones del Sprint es cero — domain.CalcularDesviacion
// devuelve nil en ese caso). La web no la calcula: solo la formatea (R23).
type PuntoHorasSprint struct {
	SprintNumero          int
	HorasEstimadas        float64
	HorasReales           float64
	EstimacionesCompletas bool
	MetricasDisponibles   bool
	Desviacion            *float64
}

// DashboardProyecto es el "modelo de lectura" del dashboard: lo arma
// ServicioDashboard reutilizando los servicios existentes. Crece de forma
// append-only entre historias (D19).
type DashboardProyecto struct {
	SPPlanificadosTotal            int
	SPCompletadosTotal             int
	VelocidadPromedio              float64
	PorcentajeHistoriasCompletadas float64
	HorasEstimadasTotal            float64
	HorasRealesTotal               float64
	EstimacionesCompletas          bool
	DesviacionHorasTotal           *float64
	SprintsCerrados                int
	VelocidadPorSprint             []PuntoVelocidadSprint
	HorasPorSprint                 []PuntoHorasSprint
}

// ServicioDashboard arma los datos del dashboard del proyecto.
type ServicioDashboard struct {
	proyectos RepositorioProyectos
	sprints   *ServicioSprints
	esfuerzo  *ServicioEsfuerzo
}

// NuevoServicioDashboard crea el servicio con los servicios existentes.
func NuevoServicioDashboard(
	proyectos RepositorioProyectos,
	sprints *ServicioSprints,
	esfuerzo *ServicioEsfuerzo,
) *ServicioDashboard {
	return &ServicioDashboard{proyectos: proyectos, sprints: sprints, esfuerzo: esfuerzo}
}

// ObtenerDashboard devuelve las métricas clave y las listas por Sprint del
// dashboard. Devuelve domain.ErrProyectoNoEncontrado sin envolver si el
// proyecto no existe; el resto de los errores de repositorio se envuelven con
// %w.
func (s *ServicioDashboard) ObtenerDashboard(proyectoID string) (DashboardProyecto, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return DashboardProyecto{}, envolverErrorRepositorio("verificar proyecto", err)
	}
	if !existe {
		return DashboardProyecto{}, domain.ErrProyectoNoEncontrado
	}

	sprints, err := s.sprints.ListarSprints(proyectoID)
	if err != nil {
		return DashboardProyecto{}, err
	}

	cerrados := make([]domain.Sprint, 0, len(sprints))
	for _, sprint := range sprints {
		if sprint.Estado == domain.EstadoSprintCerrado {
			cerrados = append(cerrados, sprint)
		}
	}

	dashboard := DashboardProyecto{SprintsCerrados: len(cerrados)}
	for _, sprint := range sprints {
		planificados, err := s.sprints.StoryPointsPlanificados(proyectoID, sprint.Numero)
		if err != nil {
			return DashboardProyecto{}, envolverErrorRepositorio(
				fmt.Sprintf("sumar story points planificados del sprint %d", sprint.Numero), err)
		}
		completados, err := s.sprints.StoryPointsCompletados(proyectoID, sprint.Numero)
		if err != nil {
			return DashboardProyecto{}, envolverErrorRepositorio(
				fmt.Sprintf("sumar story points completados del sprint %d", sprint.Numero), err)
		}
		dashboard.SPPlanificadosTotal += planificados
		dashboard.SPCompletadosTotal += completados
	}

	promedio, err := s.sprints.VelocidadPromedio(proyectoID)
	if err != nil {
		return DashboardProyecto{}, err
	}
	dashboard.VelocidadPromedio = promedio

	porcentaje, err := s.sprints.PorcentajeHistoriasCompletadas(proyectoID)
	if err != nil {
		return DashboardProyecto{}, err
	}
	dashboard.PorcentajeHistoriasCompletadas = porcentaje

	metricas, err := s.esfuerzo.MetricasProyecto(proyectoID)
	if err != nil {
		return DashboardProyecto{}, err
	}
	dashboard.HorasEstimadasTotal = metricas.HorasEstimadas
	dashboard.HorasRealesTotal = metricas.HorasReales
	dashboard.EstimacionesCompletas = metricas.EstimacionesCompletas
	dashboard.DesviacionHorasTotal = metricas.Desviacion

	for _, sprint := range cerrados {
		velocidad, err := s.sprints.VelocidadSprint(proyectoID, sprint.Numero)
		if err != nil {
			return DashboardProyecto{}, envolverErrorRepositorio(
				fmt.Sprintf("sumar velocidad del sprint %d", sprint.Numero), err)
		}
		dashboard.VelocidadPorSprint = append(dashboard.VelocidadPorSprint,
			PuntoVelocidadSprint{SprintNumero: sprint.Numero, Velocidad: velocidad})

		m, err := s.esfuerzo.MetricasSprint(proyectoID, sprint.Numero)
		if errors.Is(err, domain.ErrMetricasSprintNoDisponibles) {
			dashboard.HorasPorSprint = append(dashboard.HorasPorSprint,
				PuntoHorasSprint{SprintNumero: sprint.Numero, MetricasDisponibles: false})
			continue
		}
		if err != nil {
			return DashboardProyecto{}, envolverErrorRepositorio(
				fmt.Sprintf("obtener métricas del sprint %d", sprint.Numero), err)
		}
		dashboard.HorasPorSprint = append(dashboard.HorasPorSprint, PuntoHorasSprint{
			SprintNumero:          sprint.Numero,
			HorasEstimadas:        m.HorasEstimadas,
			HorasReales:           m.HorasReales,
			EstimacionesCompletas: m.EstimacionesCompletas,
			MetricasDisponibles:   true,
		})
	}

	return dashboard, nil
}
