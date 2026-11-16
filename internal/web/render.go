package web

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

//go:embed plantillas/*.html
var archivosPlantillas embed.FS

// mensajeErrorInfraestructura es el texto genérico mostrado ante un error 500.
const mensajeErrorInfraestructura = "Ocurrió un error inesperado. Intentá nuevamente."

// navContexto indica a dónde apuntan los enlaces de Backlog y Sprints según si
// la página pertenece a un proyecto.
type navContexto struct {
	Backlog        string
	Sprints        string
	Dashboard      string
	MostrarLeyenda bool
}

// proyectoVista es la representación de un proyecto lista para renderizar.
type proyectoVista struct {
	ID          string
	Nombre      string
	FechaInicio string
	FechaFin    string
	Estado      domain.Estado
}

// formularioProyecto conserva los datos ingresados en el alta de proyecto.
type formularioProyecto struct {
	Nombre      string
	FechaInicio string
	FechaFin    string
}

// formularioSprint conserva los datos ingresados en el alta de Sprint.
type formularioSprint struct {
	Numero      string
	Objetivo    string
	FechaInicio string
	FechaFin    string
}

// itemSprintVista es la representación de una historia asignada a un Sprint
// lista para renderizar.
type itemSprintVista struct {
	ID          string
	Titulo      string
	Estado      string
	StoryPoints int
	Pendiente   bool
}

// sprintVista es la representación de un Sprint lista para renderizar, con sus
// totales devueltos por el servicio.
type sprintVista struct {
	Numero         int
	Objetivo       string
	FechaInicio    string
	FechaFin       string
	Estado         string
	SPPlanificados int
	SPCompletados  int
	Velocidad      float64
	Activo         bool
	Items          []itemSprintVista
}

// dashboardTick es una división del eje Y de un gráfico: posición vertical y
// etiqueta numérica.
type dashboardTick struct {
	Y        int
	Etiqueta string
}

// dashboardPuntoVelocidad es una barra del gráfico de velocidad por Sprint con
// la geometría ya calculada (D15).
type dashboardPuntoVelocidad struct {
	Sprint    int
	Velocidad int
	X         int
	Y         int
	Ancho     int
	Alto      int
	TextoX    int
	TextoY    int
	Etiqueta  string
}

// dashboardPuntoHoras son las dos barras agrupadas de un Sprint en el gráfico
// de horas estimadas vs reales (D15). MetricasDisponibles=false no dibuja
// barras (D12). La desviación se copia del servicio (R23): DesvDisponible es
// false cuando el servicio informa nil y la plantilla muestra "—"; Desv es el
// valor tal como lo devuelve el servicio (con signo), válido si DesvDisponible.
type dashboardPuntoHoras struct {
	Sprint              int
	HorasEstimadas      float64
	HorasReales         float64
	MetricasDisponibles bool
	DesvDisponible      bool
	Desv                float64
	EstX                int
	EstY                int
	EstAlto             int
	RealX               int
	RealY               int
	RealAlto            int
	Ancho               int
	TextoEstX           int
	TextoEstY           int
	TextoRealX          int
	TextoRealY          int
	GrupoX              int
	EstTexto            string
	RealTexto           string
}

// dashboardVista es la representación del dashboard lista para renderizar:
// copia los valores del servicio sin recalcularlos (D5, R19).
type dashboardVista struct {
	ProyectoID                     string
	ProyectoNombre                 string
	Estado                         domain.Estado
	SPPlanificadosTotal            int
	SPCompletadosTotal             int
	VelocidadPromedio              float64
	PorcentajeHistoriasCompletadas float64
	HorasEstimadasTotal            float64
	HorasRealesTotal               float64
	EstimacionesCompletas          bool
	DesviacionTotalDisponible      bool    // false cuando el servicio informa nil; la plantilla muestra "—"
	DesviacionHorasTotal           float64 // valor tal como lo devuelve el servicio (con signo); válida si DesviacionTotalDisponible
	SprintsCerrados                int
	Velocidad                      []dashboardPuntoVelocidad
	Horas                          []dashboardPuntoHoras
	EscalaVelocidadMax             int
	EscalaHorasMax                 float64
	TicksVelocidad                 []dashboardTick
	TicksHoras                     []dashboardTick
	HayMetricasNoDisponibles       bool
}

type sprintMetricasDefectosVista struct {
	Numero     int
	Objetivo   string
	Detectados int
	Resueltos  int
}

// datosPagina reúne los datos que consumen las plantillas.
type datosPagina struct {
	Titulo                string
	Nav                   navContexto
	Proyectos             []proyectoVista
	Proyecto              *proyectoVista
	Formulario            formularioProyecto
	Items                 []domain.ItemBacklog
	ItemID                string
	FormularioItem        formularioItem
	Editando              bool
	Sprints               []sprintVista
	SprintsOpciones       []domain.Sprint
	Defectos              []domain.Defecto
	MetricasDefectos      domain.MetricasDefectos
	MetricasPorSprint     []sprintMetricasDefectosVista
	FormularioSprint      formularioSprint
	VelocidadPromedio     float64
	PorcentajeCompletadas float64
	Mensaje               string
	Error                 string
	Poker                 *pokerVista
	Dashboard             *dashboardVista
}

// vistas contiene las plantillas embebidas ya parseadas por página.
type vistas struct {
	plantillas map[string]*template.Template
}

var paginas = []string{
	"inicio",
	"proyectos",
	"proyecto",
	"formulario-proyecto",
	"backlog",
	"formulario-item",
	"defectos",
	"sprints",
	"formulario-sprint",
	"sprints-historicos",
	"poker",
	"dashboard",
	"error",
	"reporte",
	"reporte-sprint",
}

// nuevasVistas parsea las plantillas embebidas. Un fallo de parseo es un error
// de programación y se detecta al arrancar.
func nuevasVistas() *vistas {
	v := &vistas{plantillas: make(map[string]*template.Template, len(paginas))}
	for _, pagina := range paginas {
		v.plantillas[pagina] = template.Must(template.ParseFS(
			archivosPlantillas,
			"plantillas/base.html",
			"plantillas/"+pagina+".html",
		))
	}
	return v
}

// navSinProyecto arma la navegación cuando no hay proyecto en contexto.
func navSinProyecto(conLeyenda bool) navContexto {
	return navContexto{
		Backlog:        "/proyectos",
		Sprints:        "/proyectos",
		Dashboard:      "/proyectos",
		MostrarLeyenda: conLeyenda,
	}
}

// navConProyecto arma la navegación apuntando al backlog y a los sprints del
// proyecto indicado.
func navConProyecto(id string) navContexto {
	return navContexto{
		Backlog:   "/proyectos/" + id + "/backlog",
		Sprints:   "/proyectos/" + id + "/sprints",
		Dashboard: "/proyectos/" + id + "/dashboard",
	}
}

// render ejecuta el layout base con la página indicada y escribe la respuesta.
// Renderiza primero en memoria para no dejar una respuesta parcial si falla.
func (v *vistas) render(w http.ResponseWriter, estado int, nombre string, datos datosPagina) {
	plantilla, ok := v.plantillas[nombre]
	if !ok {
		log.Printf("web: plantilla desconocida %q", nombre)
		http.Error(w, mensajeErrorInfraestructura, http.StatusInternalServerError)
		return
	}

	var buffer bytes.Buffer
	if err := plantilla.ExecuteTemplate(&buffer, "base.html", datos); err != nil {
		log.Printf("web: error al renderizar la plantilla %q: %v", nombre, err)
		http.Error(w, mensajeErrorInfraestructura, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(estado)
	_, _ = w.Write(buffer.Bytes())
}
