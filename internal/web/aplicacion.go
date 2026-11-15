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
}

// ServicioBacklogWeb es el contrato de consumidor del backlog de la web.
type ServicioBacklogWeb interface {
	ListarBacklog(proyectoID string) ([]domain.ItemBacklog, error)
	CrearItem(proyectoID, titulo, descripcion, prioridad string, criterios []string) (domain.ItemBacklog, error)
	ObtenerItem(proyectoID, itemID string) (domain.ItemBacklog, error)
	EditarItem(proyectoID, itemID, titulo, descripcion, prioridad string, criterios []string, storyPoints int) (domain.ItemBacklog, error)
	EstimarItem(proyectoID, itemID string, storyPoints int) (domain.ItemBacklog, error)
}

// ServicioSprintsWeb es el contrato de consumidor de Sprints de la web.
type ServicioSprintsWeb interface {
	CrearSprint(proyectoID string, numero int, objetivo string, fechaInicio, fechaFin time.Time) (domain.Sprint, error)
	AsignarItem(proyectoID string, sprintNumero int, itemID string) error
	StoryPointsPlanificados(proyectoID string, sprintNumero int) (int, error)
	ListarSprints(proyectoID string) ([]domain.Sprint, error)
	// Nuevos en FASE 3
	ListarItemsDelSprint(proyectoID string, sprintNumero int) ([]domain.ItemBacklog, error)
	MarcarItemCompletada(proyectoID string, sprintNumero int, itemID string) error
	CerrarSprint(proyectoID string, sprintNumero int) error
	StoryPointsCompletados(proyectoID string, sprintNumero int) (int, error)
	VelocidadSprint(proyectoID string, sprintNumero int) (int, error)
	VelocidadPromedio(proyectoID string) (float64, error)
	PorcentajeHistoriasCompletadas(proyectoID string) (float64, error)
}

// Dependencias reúne los servicios y el reloj que usa la aplicación web.
type Dependencias struct {
	Proyectos ServicioProyectosWeb
	Backlog   ServicioBacklogWeb
	Sprints   ServicioSprintsWeb
	// Ahora es el reloj inyectable usado para calcular el estado del proyecto
	// con domain.EstadoPorFechas. Si es nil se usa time.Now.
	Ahora func() time.Time
}

// aplicacion agrupa las dependencias y las plantillas ya parseadas.
type aplicacion struct {
	deps   Dependencias
	vistas *vistas
}

// NuevoHandlerAplicacion arma la aplicación web completa y conserva /healthz.
// Usa plantillas embebidas; un fallo de parseo es un error de programación y
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
		panic("web: no se pudieron preparar los archivos estáticos: " + err.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", manejadorSalud(ping))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(estaticos))))
	mux.HandleFunc("GET /{$}", app.inicio)
	mux.HandleFunc("GET /proyectos", app.listarProyectos)
	mux.HandleFunc("GET /proyectos/nuevo", app.formularioAltaProyecto)
	mux.HandleFunc("POST /proyectos", app.crearProyecto)
	mux.HandleFunc("GET /proyectos/{id}", app.detalleProyecto)
	mux.HandleFunc("GET /proyectos/{id}/backlog", app.listarBacklog)
	mux.HandleFunc("GET /proyectos/{id}/backlog/nuevo", app.formularioAltaItem)
	mux.HandleFunc("POST /proyectos/{id}/backlog", app.crearItem)
	mux.HandleFunc("GET /proyectos/{id}/backlog/{itemID}/editar", app.formularioEditarItem)
	mux.HandleFunc("POST /proyectos/{id}/backlog/{itemID}/editar", app.editarItem)
	mux.HandleFunc("POST /proyectos/{id}/backlog/{itemID}/estimar", app.estimarItem)
	mux.HandleFunc("GET /proyectos/{id}/sprints", app.listarSprints)
	mux.HandleFunc("GET /proyectos/{id}/sprints/nuevo", app.formularioAltaSprint)
	mux.HandleFunc("POST /proyectos/{id}/sprints", app.crearSprint)
	mux.HandleFunc("POST /proyectos/{id}/sprints/{numero}/asignar", app.asignarItemSprint)
	mux.HandleFunc("POST /proyectos/{id}/sprints/{numero}/items/{itemID}/completar", app.completarItemSprint)
	mux.HandleFunc("POST /proyectos/{id}/sprints/{numero}/cerrar", app.cerrarSprint)
	mux.HandleFunc("GET /defectos", app.defectos)

	return mux
}

// renderizar ejecuta la plantilla indicada y escribe la respuesta.
func (a *aplicacion) renderizar(w http.ResponseWriter, estado int, nombre string, datos datosPagina) {
	a.vistas.render(w, estado, nombre, datos)
}

// responderErrorInfraestructura registra el detalle y muestra un mensaje
// genérico con HTTP 500.
func (a *aplicacion) responderErrorInfraestructura(w http.ResponseWriter, contexto string, err error) {
	log.Printf("web: %s: %v", contexto, err)
	a.renderizar(w, http.StatusInternalServerError, "error", datosPagina{
		Titulo: "Error",
		Nav:    navSinProyecto(false),
		Error:  mensajeErrorInfraestructura,
	})
}

// responderErrorProyecto muestra el mensaje del dominio con HTTP 404 cuando el
// proyecto no existe y delega el resto de los errores en el manejador genérico.
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
