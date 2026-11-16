package web

import (
	"embed"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
)

//go:embed estaticos/*.css
var archivosEstaticos embed.FS

// ServicioProyectosWeb es el contrato de consumidor de proyectos de la web.
type ServicioProyectosWeb interface {
	CrearProyecto(nombre string, fechaInicio, fechaFin time.Time) (service.ProyectoRegistrado, error)
	ListarProyectos() ([]service.ProyectoRegistrado, error)
	ObtenerProyecto(id string) (service.ProyectoRegistrado, error)
	ListarIntegrantes(proyectoID string) ([]domain.Integrante, error)
}

// ServicioBacklogWeb es el contrato de consumidor del backlog de la web.
type ServicioBacklogWeb interface {
	ListarBacklog(proyectoID string) ([]domain.ItemBacklog, error)
	CrearItem(proyectoID, titulo, descripcion, prioridad string, criterios []string) (domain.ItemBacklog, error)
	ObtenerItem(proyectoID, itemID string) (domain.ItemBacklog, error)
	EditarItem(proyectoID, itemID, titulo, descripcion, prioridad string, criterios []string, storyPoints int) (domain.ItemBacklog, error)
	EstimarItem(proyectoID, itemID string, storyPoints int) (domain.ItemBacklog, error)
}

// ServicioPokerWeb es el contrato de consumidor de Planning Poker de la web.
// Errores exactos: domain.Err*; la web los muestra con HTTP 422 o 404 segÃºn la
// secciÃ³n 8 de specs/US-12.md; el resto se registra en el log y responde HTTP
// 500 genÃ©rico.
type ServicioPokerWeb interface {
	IniciarRonda(proyectoID, itemID string) (domain.RondaPoker, error)
	Votar(proyectoID, itemID string, rondaNumero int, integranteNombre string, valor int) error
	RevelarRonda(proyectoID, itemID string, rondaNumero int, integranteNombre string) error
	ConsultarRonda(proyectoID, itemID string, rondaNumero int) (domain.ConsultaRondaPoker, error)
	RondaActual(proyectoID, itemID string) (domain.RondaPoker, bool, error)
	DiferenciasRonda(proyectoID, itemID string, rondaNumero int) (domain.DiferenciaPoker, error)
	HistorialRondas(proyectoID, itemID string) ([]domain.RondaHistorialPoker, error)
	AcordarEstimacion(proyectoID, itemID string, rondaNumero int, integranteNombre string, valor int) (domain.ItemBacklog, error)
	HistoriaActual(proyectoID, itemID string) (domain.ItemBacklog, error)
}

// ServicioSprintsWeb es el contrato de consumidor de Sprints de la web.
type ServicioSprintsWeb interface {
	CrearSprint(proyectoID string, numero int, objetivo string, fechaInicio, fechaFin time.Time) (domain.Sprint, error)
	AsignarItem(proyectoID string, sprintNumero int, itemID string) error
	StoryPointsPlanificados(proyectoID string, sprintNumero int) (int, error)
	ListarSprints(proyectoID string) ([]domain.Sprint, error)
	ListarSprintsCerrados(proyectoID string) ([]domain.Sprint, error)
	// Nuevos en FASE 3
	ListarItemsDelSprint(proyectoID string, sprintNumero int) ([]domain.ItemBacklog, error)
	MarcarItemCompletada(proyectoID string, sprintNumero int, itemID string) error
	CerrarSprint(proyectoID string, sprintNumero int) error
	StoryPointsCompletados(proyectoID string, sprintNumero int) (int, error)
	VelocidadSprint(proyectoID string, sprintNumero int) (int, error)
	VelocidadPromedio(proyectoID string) (float64, error)
	PorcentajeHistoriasCompletadas(proyectoID string) (float64, error)
}

// ServicioDefectosWeb es el contrato de consumidor de defectos de la web.
type ServicioDefectosWeb interface {
	ListarDefectos(proyectoID string) ([]domain.Defecto, error)
	RegistrarDefecto(proyectoID string, defecto domain.Defecto) (domain.Defecto, error)
	CambiarEstadoDefecto(proyectoID string, defectoID int64, estado string, sprintResolucion *int) (domain.Defecto, error)
	MetricasDefectosPorSprint(proyectoID string, sprintNumero int) (domain.MetricasDefectos, error)
	MetricasDefectosPorProyecto(proyectoID string) (domain.MetricasDefectos, error)
}

// ServicioDashboardWeb es el contrato de consumidor del dashboard de la web.
type ServicioDashboardWeb interface {
	ObtenerDashboard(proyectoID string) (service.DashboardProyecto, error)
}

// ServicioReportesWeb es el contrato de consumidor de reportes de la web.
type ServicioReportesWeb interface {
	ReporteProyecto(proyectoID string) (service.ReporteProyecto, error)
	ReporteSprint(proyectoID string, sprintNumero int) (service.ReporteSprint, error)
}

// Dependencias reÃºne los servicios y el reloj que usa la aplicaciÃ³n web.
type Dependencias struct {
	Proyectos ServicioProyectosWeb
	Backlog   ServicioBacklogWeb
	Sprints   ServicioSprintsWeb
	Poker     ServicioPokerWeb
	Defectos  ServicioDefectosWeb
	Dashboard ServicioDashboardWeb
	Reportes  ServicioReportesWeb
	// Ahora es el reloj inyectable usado para calcular el estado del proyecto
	// con domain.EstadoPorFechas. Si es nil se usa time.Now.
	Ahora func() time.Time
}

// aplicacion agrupa las dependencias y las plantillas ya parseadas.
type aplicacion struct {
	deps   Dependencias
	vistas *vistas
}

// NuevoHandlerAplicacion arma la aplicaciÃ³n web completa y conserva /healthz.
// Usa plantillas embebidas; un fallo de parseo es un error de programaciÃ³n y
// se detecta al arrancar.
func NuevoHandlerAplicacion(ping Pinger, deps Dependencias) http.Handler {
	if deps.Ahora == nil {
		deps.Ahora = time.Now
	}

	app := &aplicacion{
		deps:   deps,
		vistas: nuevasVistas(),
	}

	estaticos, err := fs.Sub(archivosEstaticos, "estaticos")
	if err != nil {
		panic("web: no se pudieron preparar los archivos estÃ¡ticos: " + err.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", manejadorSalud(ping))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(estaticos))))
	mux.HandleFunc("GET /{$}", app.inicio)
	mux.HandleFunc("GET /proyectos", app.listarProyectos)
	mux.HandleFunc("GET /proyectos/nuevo", app.formularioAltaProyecto)
	mux.HandleFunc("POST /proyectos", app.crearProyecto)
	mux.HandleFunc("GET /proyectos/{id}", app.detalleProyecto)
	mux.HandleFunc("GET /proyectos/{id}/dashboard", app.dashboardProyecto)
	mux.HandleFunc("GET /proyectos/{id}/reporte", app.reporteProyecto)
	mux.HandleFunc("GET /proyectos/{id}/reporte.pdf", app.reporteProyectoPDF)
	mux.HandleFunc("GET /proyectos/{id}/sprints/{numero}/reporte", app.reporteSprint)
	mux.HandleFunc("GET /proyectos/{id}/sprints/{numero}/reporte.pdf", app.reporteSprintPDF)
	mux.HandleFunc("GET /proyectos/{id}/backlog", app.listarBacklog)
	mux.HandleFunc("GET /proyectos/{id}/backlog/nuevo", app.formularioAltaItem)
	mux.HandleFunc("POST /proyectos/{id}/backlog", app.crearItem)
	mux.HandleFunc("GET /proyectos/{id}/backlog/{itemID}/editar", app.formularioEditarItem)
	mux.HandleFunc("POST /proyectos/{id}/backlog/{itemID}/editar", app.editarItem)
	mux.HandleFunc("POST /proyectos/{id}/backlog/{itemID}/estimar", app.estimarItem)
	mux.HandleFunc("GET /proyectos/{id}/sprints", app.listarSprints)
	mux.HandleFunc("GET /proyectos/{id}/sprints/historico", app.listarSprintsAnteriores)
	mux.HandleFunc("GET /proyectos/{id}/sprints/nuevo", app.formularioAltaSprint)
	mux.HandleFunc("POST /proyectos/{id}/sprints", app.crearSprint)
	mux.HandleFunc("POST /proyectos/{id}/sprints/{numero}/asignar", app.asignarItemSprint)
	mux.HandleFunc("POST /proyectos/{id}/sprints/{numero}/items/{itemID}/completar", app.completarItemSprint)
	mux.HandleFunc("POST /proyectos/{id}/sprints/{numero}/cerrar", app.cerrarSprint)
	mux.HandleFunc("GET /proyectos/{id}/defectos", app.listarDefectos)
	mux.HandleFunc("POST /proyectos/{id}/defectos", app.registrarDefecto)
	mux.HandleFunc("POST /proyectos/{id}/defectos/{defectoID}/estado", app.cambiarEstadoDefecto)
	mux.HandleFunc("GET /proyectos/{id}/backlog/{itemID}/poker", app.paginaPoker)
	mux.HandleFunc("POST /proyectos/{id}/backlog/{itemID}/poker/rondas", app.iniciarRondaPoker)
	mux.HandleFunc("POST /proyectos/{id}/backlog/{itemID}/poker/rondas/{numero}/votos", app.votarPoker)
	mux.HandleFunc("POST /proyectos/{id}/backlog/{itemID}/poker/rondas/{numero}/revelar", app.revelarPoker)
	mux.HandleFunc("POST /proyectos/{id}/backlog/{itemID}/poker/rondas/{numero}/acordar", app.acordarPoker)
	mux.HandleFunc("GET /defectos", app.defectos)

	return mux
}

// renderizar ejecuta la plantilla indicada y escribe la respuesta.
func (a *aplicacion) renderizar(w http.ResponseWriter, estado int, nombre string, datos datosPagina) {
	a.vistas.render(w, estado, nombre, datos)
}

// responderErrorInfraestructura registra el detalle y muestra un mensaje
// genÃ©rico con HTTP 500.
func (a *aplicacion) responderErrorInfraestructura(w http.ResponseWriter, contexto string, err error) {
	log.Printf("web: %s: %v", contexto, err)
	a.renderizar(w, http.StatusInternalServerError, "error", datosPagina{
		Titulo: "Error",
		Nav:    navSinProyecto(false),
		Error:  mensajeErrorInfraestructura,
	})
}

// responderErrorProyecto muestra el mensaje del dominio con HTTP 404 cuando el
// proyecto no existe y delega el resto de los errores en el manejador genÃ©rico.
func (a *aplicacion) responderErrorProyecto(w http.ResponseWriter, id string, err error) {
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		a.renderizar(w, http.StatusNotFound, "error", datosPagina{
			Titulo: "Proyecto no encontrado",
			Nav:    navSinProyecto(false),
			Error:  err.Error(),
		})
		return
	}
	a.responderErrorInfraestructura(w, "obtener proyecto "+id, err)
}
