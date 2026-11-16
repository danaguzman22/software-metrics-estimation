package service

import (
	"errors"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// HistoriaReporte es una historia del backlog o asignada a un Sprint, lista
// para presentarse en el reporte. HorasEstimadas es nil cuando la historia
// no está estimada en horas ("sin estimar"). HorasReales es la suma
// histórica de esfuerzo de la historia en todo el proyecto (US-17, D13); en
// la página figura bajo el encabezado "Horas reales (acumuladas)".
type HistoriaReporte struct {
	ID             string
	Titulo         string
	Estado         string
	StoryPoints    int
	HorasEstimadas *float64
	HorasReales    float64
}

// ResumenDefectos combina dos fuentes (D6, D14): Detectados y Resueltos son
// reutilizados tal cual de US-18 (MetricasDefectosPorProyecto en el reporte
// del proyecto; MetricasDefectosPorSprint en el del Sprint); Alta..Resuelto
// son los conteos por severidad y por estado armados localmente iterando la
// lista devuelta por ServicioDefectos.ListarDefectos (US-18 no los expone).
// En el reporte de Sprint, Alta..Total y Defectos contabilizan los defectos
// detectados en el Sprint; DefectosResueltos la lista de resueltos en él.
type ResumenDefectos struct {
	Alta     int
	Media    int
	Baja     int
	Abierto  int
	EnCurso  int
	Resuelto int
	Total    int
	// Solo se completan en el reporte de Sprint (D14):
	Detectados        int
	Resueltos         int
	Defectos          []domain.Defecto // proyecto: todos; sprint: detectados en el Sprint
	DefectosResueltos []domain.Defecto // solo sprint: resueltos en el Sprint
}

// ReporteProyecto es el "modelo de lectura" del reporte del proyecto.
// Los totales se copian de ObtenerDashboard (D2, R2).
type ReporteProyecto struct {
	// Totales del dashboard (reutilizados, no recalculados):
	SPPlanificadosTotal            int
	SPCompletadosTotal             int
	VelocidadPromedio              float64
	PorcentajeHistoriasCompletadas float64
	HorasEstimadasTotal            float64
	HorasRealesTotal               float64
	EstimacionesCompletas          bool
	DesviacionHorasTotal           *float64
	SprintsCerrados                int
	SprintsCerradosVelocidad       []PuntoVelocidadSprint // reutiliza el tipo de US-19
	// Historias (D12, R3):
	HistoriasPlanificadas []HistoriaReporte
	HistoriasCompletadas  []HistoriaReporte
	// Defectos (D14, R5):
	Defectos ResumenDefectos
}

// ReporteSprint es el "modelo de lectura" del reporte de un Sprint.
type ReporteSprint struct {
	Numero                int
	Objetivo              string
	FechaInicio           time.Time
	FechaFin              time.Time
	Estado                string
	HistoriasAsignadas    []HistoriaReporte
	SPPlanificados        int
	SPCompletados         int
	Velocidad             int
	HorasEstimadas        float64
	HorasReales           float64
	EstimacionesCompletas bool
	MetricasDisponibles   bool // false cuando MetricasSprint devuelve ErrMetricasSprintNoDisponibles (D16)
	DesviacionHoras       *float64
	Defectos              ResumenDefectos
}

// ServicioReportes arma los datos de los reportes del proyecto y del Sprint.
type ServicioReportes struct {
	proyectos RepositorioProyectos
	backlog   *ServicioBacklog
	sprints   *ServicioSprints
	esfuerzo  *ServicioEsfuerzo
	dashboard *ServicioDashboard
	defectos  *ServicioDefectos
}

// NuevoServicioReportes crea el servicio con los servicios existentes.
func NuevoServicioReportes(
	proyectos RepositorioProyectos,
	backlog *ServicioBacklog,
	sprints *ServicioSprints,
	esfuerzo *ServicioEsfuerzo,
	dashboard *ServicioDashboard,
	defectos *ServicioDefectos,
) *ServicioReportes {
	return &ServicioReportes{
		proyectos: proyectos,
		backlog:   backlog,
		sprints:   sprints,
		esfuerzo:  esfuerzo,
		dashboard: dashboard,
		defectos:  defectos,
	}
}

// historiaReporteDesde arma la HistoriaReporte desde el ítem del backlog y las
// métricas de esfuerzo de la historia (su "suma histórica" de horas reales).
func historiaReporteDesde(item domain.ItemBacklog, metricas domain.MetricasEsfuerzo) HistoriaReporte {
	historia := HistoriaReporte{
		ID:          item.ID,
		Titulo:      item.Titulo,
		Estado:      item.Estado,
		StoryPoints: item.StoryPoints,
		HorasReales: metricas.HorasReales,
	}
	if item.HorasEstimadas != nil {
		horas := *item.HorasEstimadas
		historia.HorasEstimadas = &horas
	}
	return historia
}

// contarDefectosPorSeveridadYEstado acumula los conteos por severidad y por
// estado sobre la lista dada (US-18 no los expone; D14).
func contarDefectosPorSeveridadYEstado(resumen *ResumenDefectos, defectos []domain.Defecto) {
	for _, defecto := range defectos {
		switch defecto.Severidad {
		case domain.SeveridadDefectoAlta:
			resumen.Alta++
		case domain.SeveridadDefectoMedia:
			resumen.Media++
		case domain.SeveridadDefectoBaja:
			resumen.Baja++
		}
		switch defecto.Estado {
		case domain.EstadoDefectoAbierto:
			resumen.Abierto++
		case domain.EstadoDefectoEnCurso:
			resumen.EnCurso++
		case domain.EstadoDefectoResuelto:
			resumen.Resuelto++
		}
	}
}

// ReporteProyecto devuelve el reporte completo del proyecto.
// Devuelve domain.ErrProyectoNoEncontrado sin envolver si el proyecto no
// existe; el resto de los errores de repositorio se envuelven con %w.
func (s *ServicioReportes) ReporteProyecto(proyectoID string) (ReporteProyecto, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return ReporteProyecto{}, envolverErrorRepositorio("verificar proyecto", err)
	}
	if !existe {
		return ReporteProyecto{}, domain.ErrProyectoNoEncontrado
	}

	historias, err := s.backlog.ListarBacklog(proyectoID)
	if err != nil {
		return ReporteProyecto{}, err
	}

	dash, err := s.dashboard.ObtenerDashboard(proyectoID)
	if err != nil {
		return ReporteProyecto{}, err
	}

	var planificadas, completadas []HistoriaReporte
	for _, item := range historias {
		metricas, err := s.esfuerzo.MetricasHistoria(proyectoID, item.ID)
		if err != nil {
			return ReporteProyecto{}, err
		}
		historia := historiaReporteDesde(item, metricas)
		if item.Estado == domain.EstadoCompletada {
			completadas = append(completadas, historia)
			continue
		}
		planificadas = append(planificadas, historia)
	}

	metricasDefectos, err := s.defectos.MetricasDefectosPorProyecto(proyectoID)
	if err != nil {
		return ReporteProyecto{}, err
	}
	listaDefectos, err := s.defectos.ListarDefectos(proyectoID)
	if err != nil {
		return ReporteProyecto{}, err
	}
	resumenDefectos := ResumenDefectos{
		Detectados: metricasDefectos.Detectados,
		Resueltos:  metricasDefectos.Resueltos,
		Total:      metricasDefectos.Detectados,
	}
	contarDefectosPorSeveridadYEstado(&resumenDefectos, listaDefectos)
	if len(listaDefectos) > 0 {
		resumenDefectos.Defectos = listaDefectos
	}

	return ReporteProyecto{
		SPPlanificadosTotal:            dash.SPPlanificadosTotal,
		SPCompletadosTotal:             dash.SPCompletadosTotal,
		VelocidadPromedio:              dash.VelocidadPromedio,
		PorcentajeHistoriasCompletadas: dash.PorcentajeHistoriasCompletadas,
		HorasEstimadasTotal:            dash.HorasEstimadasTotal,
		HorasRealesTotal:               dash.HorasRealesTotal,
		EstimacionesCompletas:          dash.EstimacionesCompletas,
		DesviacionHorasTotal:           dash.DesviacionHorasTotal,
		SprintsCerrados:                dash.SprintsCerrados,
		SprintsCerradosVelocidad:       dash.VelocidadPorSprint,
		HistoriasPlanificadas:          planificadas,
		HistoriasCompletadas:           completadas,
		Defectos:                       resumenDefectos,
	}, nil
}

// ReporteSprint devuelve el reporte del Sprint indicado.
// Devuelve domain.ErrProyectoNoEncontrado sin envolver si el proyecto no
// existe y domain.ErrSprintNoEncontrado sin envolver si el Sprint no existe
// (V2); el resto de los errores de repositorio se envuelven con %w.
func (s *ServicioReportes) ReporteSprint(proyectoID string, sprintNumero int) (ReporteSprint, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return ReporteSprint{}, envolverErrorRepositorio("verificar proyecto", err)
	}
	if !existe {
		return ReporteSprint{}, domain.ErrProyectoNoEncontrado
	}

	listaSprints, err := s.sprints.ListarSprints(proyectoID)
	if err != nil {
		return ReporteSprint{}, err
	}
	sprint, encontrado := sprintPorNumero(listaSprints, sprintNumero)
	if !encontrado {
		return ReporteSprint{}, domain.ErrSprintNoEncontrado
	}

	planificados, err := s.sprints.StoryPointsPlanificados(proyectoID, sprintNumero)
	if err != nil {
		return ReporteSprint{}, err
	}
	completados, err := s.sprints.StoryPointsCompletados(proyectoID, sprintNumero)
	if err != nil {
		return ReporteSprint{}, err
	}
	velocidad, err := s.sprints.VelocidadSprint(proyectoID, sprintNumero)
	if err != nil {
		return ReporteSprint{}, err
	}

	asignadas, err := s.sprints.ListarItemsDelSprint(proyectoID, sprintNumero)
	if err != nil {
		return ReporteSprint{}, err
	}
	var historias []HistoriaReporte
	for _, item := range asignadas {
		metricas, err := s.esfuerzo.MetricasHistoria(proyectoID, item.ID)
		if err != nil {
			return ReporteSprint{}, err
		}
		historias = append(historias, historiaReporteDesde(item, metricas))
	}

	reporte := ReporteSprint{
		Numero:             sprint.Numero,
		Objetivo:           sprint.Objetivo,
		FechaInicio:        sprint.FechaInicio,
		FechaFin:           sprint.FechaFin,
		Estado:             sprint.Estado,
		HistoriasAsignadas: historias,
		SPPlanificados:     planificados,
		SPCompletados:      completados,
		Velocidad:          velocidad,
	}
	metricasSprint, err := s.esfuerzo.MetricasSprint(proyectoID, sprintNumero)
	switch {
	case errors.Is(err, domain.ErrMetricasSprintNoDisponibles):
		// Sprint cerrado sin snapshot de esfuerzo: no es un error (D16).
	case err != nil:
		return ReporteSprint{}, err
	default:
		reporte.MetricasDisponibles = true
		reporte.HorasEstimadas = metricasSprint.HorasEstimadas
		reporte.HorasReales = metricasSprint.HorasReales
		reporte.EstimacionesCompletas = metricasSprint.EstimacionesCompletas
		reporte.DesviacionHoras = desviacionSprint(metricasSprint)
	}

	metricasDefectos, err := s.defectos.MetricasDefectosPorSprint(proyectoID, sprintNumero)
	if err != nil {
		return ReporteSprint{}, err
	}
	listaDefectos, err := s.defectos.ListarDefectos(proyectoID)
	if err != nil {
		return ReporteSprint{}, err
	}
	resumenDefectos := ResumenDefectos{
		Detectados: metricasDefectos.Detectados,
		Resueltos:  metricasDefectos.Resueltos,
	}
	var detectados, resueltos []domain.Defecto
	for _, defecto := range listaDefectos {
		if defecto.SprintDeteccion == sprintNumero {
			detectados = append(detectados, defecto)
		}
		if defecto.SprintResolucion != nil && *defecto.SprintResolucion == sprintNumero {
			resueltos = append(resueltos, defecto)
		}
	}
	contarDefectosPorSeveridadYEstado(&resumenDefectos, detectados)
	resumenDefectos.Total = len(detectados)
	resumenDefectos.Defectos = detectados
	resumenDefectos.DefectosResueltos = resueltos
	reporte.Defectos = resumenDefectos

	return reporte, nil
}

func sprintPorNumero(sprints []domain.Sprint, numero int) (domain.Sprint, bool) {
	for _, sprint := range sprints {
		if sprint.Numero == numero {
			return sprint, true
		}
	}
	return domain.Sprint{}, false
}
