package service

import (
	"errors"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// errReporteNoImplementado indica que el caso de uso todavía no fue
// implementado (stub de la fase RED).
var errReporteNoImplementado = errors.New("no implementado")

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

// ReporteProyecto devuelve el reporte completo del proyecto.
// Devuelve domain.ErrProyectoNoEncontrado sin envolver si el proyecto no
// existe; el resto de los errores de repositorio se envuelven con %w.
func (s *ServicioReportes) ReporteProyecto(proyectoID string) (ReporteProyecto, error) {
	return ReporteProyecto{}, errReporteNoImplementado
}

// ReporteSprint devuelve el reporte del Sprint indicado.
// Devuelve domain.ErrProyectoNoEncontrado sin envolver si el proyecto no
// existe y domain.ErrSprintNoEncontrado sin envolver si el Sprint no existe
// (V2); el resto de los errores de repositorio se envuelven con %w.
func (s *ServicioReportes) ReporteSprint(proyectoID string, sprintNumero int) (ReporteSprint, error) {
	return ReporteSprint{}, errReporteNoImplementado
}
