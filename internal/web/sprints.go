package web

import (
	"errors"
	"net/http"
	"strings"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// listarSprints muestra los Sprints del proyecto con su estado y sus Story
// Points planificados devueltos por el servicio.
func (a *aplicacion) listarSprints(w http.ResponseWriter, r *http.Request) {
	a.renderizarSprints(w, r.PathValue("id"), http.StatusOK, "")
}

// formularioAltaSprint muestra el formulario de alta de un Sprint.
func (a *aplicacion) formularioAltaSprint(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if _, err := a.deps.Proyectos.ObtenerProyecto(id); err != nil {
		a.responderErrorProyecto(w, id, err)
		return
	}

	a.renderizarFormularioSprint(w, http.StatusOK, id, formularioSprint{}, "")
}

// crearSprint delega el alta en el servicio y redirige al listado con 303.
func (a *aplicacion) crearSprint(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		a.renderizarFormularioSprint(w, http.StatusUnprocessableEntity, id, formularioSprint{}, "No se pudo leer el formulario.")
		return
	}

	formulario := formularioSprint{
		Numero:      strings.TrimSpace(r.PostFormValue("numero")),
		Objetivo:    r.PostFormValue("objetivo"),
		FechaInicio: strings.TrimSpace(r.PostFormValue("fechaInicio")),
		FechaFin:    strings.TrimSpace(r.PostFormValue("fechaFin")),
	}

	numero, err := parsearEntero(formulario.Numero, "El número del sprint debe ser un número entero")
	if err != nil {
		a.renderizarFormularioSprint(w, http.StatusUnprocessableEntity, id, formulario, err.Error())
		return
	}

	fechaInicio, err := parsearFechaObligatoria(formulario.FechaInicio, "inicio")
	if err != nil {
		a.renderizarFormularioSprint(w, http.StatusUnprocessableEntity, id, formulario, err.Error())
		return
	}

	fechaFin, err := parsearFechaObligatoria(formulario.FechaFin, "fin")
	if err != nil {
		a.renderizarFormularioSprint(w, http.StatusUnprocessableEntity, id, formulario, err.Error())
		return
	}

	_, err = a.deps.Sprints.CrearSprint(id, numero, formulario.Objetivo, fechaInicio, fechaFin)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNumeroSprintInvalido),
			errors.Is(err, domain.ErrObjetivoSprintVacio),
			errors.Is(err, domain.ErrFechasSprintInvalidas),
			errors.Is(err, domain.ErrSprintActivoExistente),
			errors.Is(err, domain.ErrSprintDuplicado):
			a.renderizarFormularioSprint(w, http.StatusUnprocessableEntity, id, formulario, err.Error())
		case errors.Is(err, domain.ErrProyectoNoEncontrado):
			a.responderErrorSprints(w, id, err)
		default:
			a.responderErrorInfraestructura(w, "crear sprint", err)
		}
		return
	}

	w.Header().Set("Location", "/proyectos/"+id+"/sprints")
	w.WriteHeader(http.StatusSeeOther)
}

// asignarItemSprint delega la asignación en el servicio y redirige al listado
// con 303.
func (a *aplicacion) asignarItemSprint(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	numero, err := parsearEntero(strings.TrimSpace(r.PathValue("numero")), "El número del sprint debe ser un número entero")
	if err != nil {
		a.renderizarSprints(w, id, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err := r.ParseForm(); err != nil {
		a.renderizarSprints(w, id, http.StatusUnprocessableEntity, "No se pudo leer el formulario.")
		return
	}

	itemID := r.PostFormValue("itemID")
	if err := a.deps.Sprints.AsignarItem(id, numero, itemID); err != nil {
		switch {
		case errors.Is(err, domain.ErrItemNoEstimado),
			errors.Is(err, domain.ErrItemNoPendiente),
			errors.Is(err, domain.ErrItemYaAsignado),
			errors.Is(err, domain.ErrSprintCerrado):
			a.renderizarSprints(w, id, http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, domain.ErrProyectoNoEncontrado),
			errors.Is(err, domain.ErrSprintNoEncontrado),
			errors.Is(err, domain.ErrItemNoEncontrado):
			a.responderErrorSprints(w, id, err)
		default:
			a.responderErrorInfraestructura(w, "asignar ítem al sprint", err)
		}
		return
	}

	w.Header().Set("Location", "/proyectos/"+id+"/sprints")
	w.WriteHeader(http.StatusSeeOther)
}

// renderizarSprints muestra el listado de Sprints con un mensaje opcional.
func (a *aplicacion) renderizarSprints(w http.ResponseWriter, proyectoID string, estado int, mensaje string) {
	proyecto, err := a.deps.Proyectos.ObtenerProyecto(proyectoID)
	if err != nil {
		a.responderErrorSprints(w, proyectoID, err)
		return
	}

	sprints, err := a.deps.Sprints.ListarSprints(proyectoID)
	if err != nil {
		a.responderErrorSprints(w, proyectoID, err)
		return
	}

	vistas := make([]sprintVista, 0, len(sprints))
	for _, sprint := range sprints {
		planificados, err := a.deps.Sprints.StoryPointsPlanificados(proyectoID, sprint.Numero)
		if err != nil {
			a.responderErrorSprints(w, proyectoID, err)
			return
		}
		vistas = append(vistas, sprintVistaDesde(sprint, planificados))
	}

	items, err := a.deps.Backlog.ListarBacklog(proyectoID)
	if err != nil {
		a.responderErrorSprints(w, proyectoID, err)
		return
	}

	a.renderizar(w, estado, "sprints", datosPagina{
		Titulo:   "Sprints",
		Nav:      navConProyecto(proyectoID),
		Proyecto: &proyectoVista{ID: proyectoID, Nombre: proyecto.Proyecto.Nombre},
		Sprints:  vistas,
		Items:    items,
		Error:    mensaje,
	})
}

// renderizarFormularioSprint muestra el formulario de alta de Sprint.
func (a *aplicacion) renderizarFormularioSprint(w http.ResponseWriter, estado int, proyectoID string, formulario formularioSprint, mensaje string) {
	a.renderizar(w, estado, "formulario-sprint", datosPagina{
		Titulo:           "Nuevo Sprint",
		Nav:              navConProyecto(proyectoID),
		Proyecto:         &proyectoVista{ID: proyectoID},
		FormularioSprint: formulario,
		Error:            mensaje,
	})
}

// responderErrorSprints muestra 404 con el mensaje del dominio para un proyecto,
// un ítem o un Sprint inexistente y 500 genérico para el resto de los errores.
func (a *aplicacion) responderErrorSprints(w http.ResponseWriter, proyectoID string, err error) {
	switch {
	case errors.Is(err, domain.ErrProyectoNoEncontrado):
		a.renderizar(w, http.StatusNotFound, "error", datosPagina{
			Titulo: "Proyecto no encontrado",
			Nav:    navSinProyecto(false),
			Error:  err.Error(),
		})
	case errors.Is(err, domain.ErrItemNoEncontrado):
		a.renderizar(w, http.StatusNotFound, "error", datosPagina{
			Titulo: "Ítem no encontrado",
			Nav:    navConProyecto(proyectoID),
			Error:  err.Error(),
		})
	case errors.Is(err, domain.ErrSprintNoEncontrado):
		a.renderizar(w, http.StatusNotFound, "error", datosPagina{
			Titulo: "Sprint no encontrado",
			Nav:    navConProyecto(proyectoID),
			Error:  err.Error(),
		})
	default:
		a.responderErrorInfraestructura(w, "operar los sprints", err)
	}
}

// sprintVistaDesde arma la vista de un Sprint con el total planificado ya
// calculado por el servicio.
func sprintVistaDesde(sprint domain.Sprint, planificados int) sprintVista {
	return sprintVista{
		Numero:         sprint.Numero,
		Objetivo:       sprint.Objetivo,
		FechaInicio:    sprint.FechaInicio.Format(formatoFecha),
		FechaFin:       sprint.FechaFin.Format(formatoFecha),
		Estado:         sprint.Estado,
		SPPlanificados: planificados,
	}
}
