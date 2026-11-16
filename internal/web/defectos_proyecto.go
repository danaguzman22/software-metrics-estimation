package web

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func (a *aplicacion) listarDefectos(w http.ResponseWriter, r *http.Request) {
	a.renderizarDefectos(w, r.PathValue("id"), http.StatusOK, "")
}

func (a *aplicacion) renderizarDefectos(w http.ResponseWriter, proyectoID string, codigo int, mensaje string) {
	proyecto, err := a.deps.Proyectos.ObtenerProyecto(proyectoID)
	if err != nil {
		a.responderErrorProyecto(w, proyectoID, err)
		return
	}
	defectos, err := a.deps.Defectos.ListarDefectos(proyectoID)
	if err != nil {
		a.responderErrorDefectos(w, proyectoID, err)
		return
	}
	items, err := a.deps.Backlog.ListarBacklog(proyectoID)
	if err != nil {
		a.responderErrorDefectos(w, proyectoID, err)
		return
	}
	sprints, err := a.deps.Sprints.ListarSprints(proyectoID)
	if err != nil {
		a.responderErrorDefectos(w, proyectoID, err)
		return
	}
	metricasProyecto, err := a.deps.Defectos.MetricasDefectosPorProyecto(proyectoID)
	if err != nil {
		a.responderErrorDefectos(w, proyectoID, err)
		return
	}
	metricasSprints := make([]sprintMetricasDefectosVista, 0, len(sprints))
	for _, sprint := range sprints {
		metricas, err := a.deps.Defectos.MetricasDefectosPorSprint(proyectoID, sprint.Numero)
		if err != nil {
			a.responderErrorDefectos(w, proyectoID, err)
			return
		}
		metricasSprints = append(metricasSprints, sprintMetricasDefectosVista{
			Numero:     sprint.Numero,
			Objetivo:   sprint.Objetivo,
			Detectados: metricas.Detectados,
			Resueltos:  metricas.Resueltos,
		})
	}
	a.renderizar(w, codigo, "defectos", datosPagina{
		Titulo:            "Defectos",
		Nav:               navConProyecto(proyectoID),
		Proyecto:          &proyectoVista{ID: proyectoID, Nombre: proyecto.Proyecto.Nombre},
		Defectos:          defectos,
		Items:             items,
		SprintsOpciones:   sprints,
		MetricasDefectos:  metricasProyecto,
		MetricasPorSprint: metricasSprints,
		Error:             mensaje,
	})
}

func (a *aplicacion) registrarDefecto(w http.ResponseWriter, r *http.Request) {
	proyectoID := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		a.renderizarDefectos(w, proyectoID, http.StatusUnprocessableEntity, "No se pudo leer el formulario.")
		return
	}
	descripcion := r.PostFormValue("descripcion")
	severidad := r.PostFormValue("severidad")
	estado := r.PostFormValue("estado")
	itemID := r.PostFormValue("item_id")
	sprintDeteccion, err := strconv.Atoi(strings.TrimSpace(r.PostFormValue("sprint_deteccion")))
	if err != nil {
		a.renderizarDefectos(w, proyectoID, http.StatusUnprocessableEntity, "El Sprint de detección debe ser un número entero.")
		return
	}
	sprintResolucion, err := enteroOpcionalDefecto(strings.TrimSpace(r.PostFormValue("sprint_resolucion")))
	if err != nil {
		a.renderizarDefectos(w, proyectoID, http.StatusUnprocessableEntity, "El Sprint de resolución debe ser un número entero.")
		return
	}
	defecto := domain.Defecto{
		Descripcion:      descripcion,
		Severidad:        severidad,
		Estado:           estado,
		ItemID:           itemID,
		SprintDeteccion:  sprintDeteccion,
		SprintResolucion: sprintResolucion,
	}
	if _, err := a.deps.Defectos.RegistrarDefecto(proyectoID, defecto); err != nil {
		a.responderErrorDefectos(w, proyectoID, err)
		return
	}
	w.Header().Set("Location", "/proyectos/"+proyectoID+"/defectos")
	w.WriteHeader(http.StatusSeeOther)
}

func (a *aplicacion) cambiarEstadoDefecto(w http.ResponseWriter, r *http.Request) {
	proyectoID := r.PathValue("id")
	defectoID, err := strconv.ParseInt(strings.TrimSpace(r.PathValue("defectoID")), 10, 64)
	if err != nil || defectoID < 1 {
		a.responderErrorRuta(w, proyectoID, "El identificador del defecto debe ser un número entero positivo.")
		return
	}
	if err := r.ParseForm(); err != nil {
		a.renderizarDefectos(w, proyectoID, http.StatusUnprocessableEntity, "No se pudo leer el formulario.")
		return
	}
	estado := r.PostFormValue("estado")
	sprintResolucion, err := enteroOpcionalDefecto(strings.TrimSpace(r.PostFormValue("sprint_resolucion")))
	if err != nil {
		a.renderizarDefectos(w, proyectoID, http.StatusUnprocessableEntity, "El Sprint de resolución debe ser un número entero.")
		return
	}
	if _, err := a.deps.Defectos.CambiarEstadoDefecto(proyectoID, defectoID, estado, sprintResolucion); err != nil {
		a.responderErrorDefectos(w, proyectoID, err)
		return
	}
	w.Header().Set("Location", "/proyectos/"+proyectoID+"/defectos")
	w.WriteHeader(http.StatusSeeOther)
}

func enteroOpcionalDefecto(valor string) (*int, error) {
	if valor == "" {
		return nil, nil
	}
	numero, err := strconv.Atoi(valor)
	if err != nil {
		return nil, err
	}
	return &numero, nil
}

func (a *aplicacion) responderErrorDefectos(w http.ResponseWriter, proyectoID string, err error) {
	switch {
	case errors.Is(err, domain.ErrProyectoNoEncontrado):
		a.renderizar(w, http.StatusNotFound, "error", datosPagina{
			Titulo: "Proyecto no encontrado",
			Nav:    navSinProyecto(false),
			Error:  err.Error(),
		})
	case errors.Is(err, domain.ErrItemNoEncontrado),
		errors.Is(err, domain.ErrSprintNoEncontrado),
		errors.Is(err, domain.ErrDefectoNoEncontrado):
		a.renderizar(w, http.StatusNotFound, "error", datosPagina{
			Titulo: "Referencia no encontrada",
			Nav:    navConProyecto(proyectoID),
			Error:  err.Error(),
		})
	case errors.Is(err, domain.ErrDescripcionDefectoVacia),
		errors.Is(err, domain.ErrSeveridadDefectoInvalida),
		errors.Is(err, domain.ErrEstadoDefectoInvalido),
		errors.Is(err, domain.ErrSprintDeteccionDefectoInvalido),
		errors.Is(err, domain.ErrSprintResolucionDefectoObligatorio),
		errors.Is(err, domain.ErrSprintResolucionDefectoInvalido),
		errors.Is(err, domain.ErrSprintResolucionAnteriorADeteccion):
		a.renderizarDefectos(w, proyectoID, http.StatusUnprocessableEntity, err.Error())
	default:
		a.responderErrorInfraestructura(w, "operar defectos del proyecto", err)
	}
}
