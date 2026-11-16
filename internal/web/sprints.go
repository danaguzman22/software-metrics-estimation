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

// listarSprintsAnteriores muestra el historial de Sprints cerrados del proyecto.
func (a *aplicacion) listarSprintsAnteriores(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	proyecto, err := a.deps.Proyectos.ObtenerProyecto(id)
	if err != nil {
		a.responderErrorSprints(w, id, err)
		return
	}

	sprints, err := a.deps.Sprints.ListarSprintsCerrados(id)
	if err != nil {
		a.responderErrorSprints(w, id, err)
		return
	}

	vistas := make([]sprintVista, 0, len(sprints))
	for _, sprint := range sprints {
		planificados, err := a.deps.Sprints.StoryPointsPlanificados(id, sprint.Numero)
		if err != nil {
			a.responderErrorSprints(w, id, err)
			return
		}
		completados, err := a.deps.Sprints.StoryPointsCompletados(id, sprint.Numero)
		if err != nil {
			a.responderErrorSprints(w, id, err)
			return
		}
		velocidad, err := a.deps.Sprints.VelocidadSprint(id, sprint.Numero)
		if err != nil {
			a.responderErrorSprints(w, id, err)
			return
		}
		vistas = append(vistas, sprintVistaDesde(sprint, planificados, completados, velocidad, nil))
	}

	a.renderizar(w, http.StatusOK, "sprints-historicos", datosPagina{
		Titulo:   "Histórico de Sprints",
		Nav:      navConProyecto(id),
		Proyecto: &proyectoVista{ID: id, Nombre: proyecto.Proyecto.Nombre},
		Sprints:  vistas,
		Migas: migasDeProyecto(proyecto.Proyecto.Nombre, id,
			Miga{Texto: "Sprints", URL: "/proyectos/" + id + "/sprints"},
			Miga{Texto: "Histórico", URL: "/proyectos/" + id + "/sprints/historico"},
		),
		Volver: "/proyectos/" + id + "/sprints",
	})
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

// completarItemSprint delega el marcado como completada en el servicio y
// redirige al listado con 303.
func (a *aplicacion) completarItemSprint(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	numero, err := parsearEntero(strings.TrimSpace(r.PathValue("numero")), "El número del sprint debe ser un número entero")
	if err != nil {
		a.responderErrorRuta(w, id, err.Error())
		return
	}

	itemID := r.PathValue("itemID")
	if _, ok := domain.NumeroDeItem(itemID); !ok {
		a.responderErrorRuta(w, id, "El ID de la historia debe tener el formato HU-NN")
		return
	}

	if err := a.deps.Sprints.MarcarItemCompletada(id, numero, itemID); err != nil {
		switch {
		case errors.Is(err, domain.ErrSprintCerrado),
			errors.Is(err, domain.ErrItemNoPendiente),
			errors.Is(err, domain.ErrItemNoAsignadoAlSprint):
			a.renderizarSprints(w, id, http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, domain.ErrProyectoNoEncontrado),
			errors.Is(err, domain.ErrSprintNoEncontrado),
			errors.Is(err, domain.ErrItemNoEncontrado):
			a.responderErrorSprints(w, id, err)
		default:
			a.responderErrorInfraestructura(w, "marcar ítem como completado", err)
		}
		return
	}

	w.Header().Set("Location", "/proyectos/"+id+"/sprints")
	w.WriteHeader(http.StatusSeeOther)
}

// cerrarSprint delega el cierre del Sprint en el servicio y redirige al
// listado con 303.
func (a *aplicacion) cerrarSprint(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	numero, err := parsearEntero(strings.TrimSpace(r.PathValue("numero")), "El número del sprint debe ser un número entero")
	if err != nil {
		a.responderErrorRuta(w, id, err.Error())
		return
	}

	if err := a.deps.Sprints.CerrarSprint(id, numero); err != nil {
		switch {
		case errors.Is(err, domain.ErrSprintCerrado),
			errors.Is(err, domain.ErrItemNoPendiente),
			errors.Is(err, domain.ErrItemNoAsignadoAlSprint):
			a.renderizarSprints(w, id, http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, domain.ErrProyectoNoEncontrado),
			errors.Is(err, domain.ErrSprintNoEncontrado):
			a.responderErrorSprints(w, id, err)
		default:
			a.responderErrorInfraestructura(w, "cerrar sprint", err)
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
		completados, err := a.deps.Sprints.StoryPointsCompletados(proyectoID, sprint.Numero)
		if err != nil {
			a.responderErrorSprints(w, proyectoID, err)
			return
		}
		velocidad, err := a.deps.Sprints.VelocidadSprint(proyectoID, sprint.Numero)
		if err != nil {
			a.responderErrorSprints(w, proyectoID, err)
			return
		}
		itemsSprint, err := a.deps.Sprints.ListarItemsDelSprint(proyectoID, sprint.Numero)
		if err != nil {
			a.responderErrorSprints(w, proyectoID, err)
			return
		}
		vistas = append(vistas, sprintVistaDesde(sprint, planificados, completados, velocidad, itemsSprint))
	}

	promedio, err := a.deps.Sprints.VelocidadPromedio(proyectoID)
	if err != nil {
		a.responderErrorSprints(w, proyectoID, err)
		return
	}
	porcentaje, err := a.deps.Sprints.PorcentajeHistoriasCompletadas(proyectoID)
	if err != nil {
		a.responderErrorSprints(w, proyectoID, err)
		return
	}

	items, err := a.deps.Backlog.ListarBacklog(proyectoID)
	if err != nil {
		a.responderErrorSprints(w, proyectoID, err)
		return
	}

	a.renderizar(w, estado, "sprints", datosPagina{
		Titulo:                "Sprints",
		Nav:                   navConProyecto(proyectoID),
		Proyecto:              &proyectoVista{ID: proyectoID, Nombre: proyecto.Proyecto.Nombre},
		Sprints:               vistas,
		Items:                 items,
		VelocidadPromedio:     promedio,
		PorcentajeCompletadas: porcentaje,
		Error:                 mensaje,
		Migas: migasDeProyecto(proyecto.Proyecto.Nombre, proyectoID,
			Miga{Texto: "Sprints", URL: "/proyectos/" + proyectoID + "/sprints"},
		),
		Volver: "/proyectos/" + proyectoID,
	})
}

// renderizarFormularioSprint muestra el formulario de alta de Sprint.
func (a *aplicacion) renderizarFormularioSprint(w http.ResponseWriter, estado int, proyectoID string, formulario formularioSprint, mensaje string) {
	proyecto, err := a.deps.Proyectos.ObtenerProyecto(proyectoID)
	if err != nil {
		a.responderErrorSprints(w, proyectoID, err)
		return
	}

	a.renderizar(w, estado, "formulario-sprint", datosPagina{
		Titulo:           "Nuevo Sprint",
		Nav:              navConProyecto(proyectoID),
		Proyecto:         &proyectoVista{ID: proyectoID, Nombre: proyecto.Proyecto.Nombre},
		FormularioSprint: formulario,
		Error:            mensaje,
		Migas: migasDeProyecto(proyecto.Proyecto.Nombre, proyectoID,
			Miga{Texto: "Sprints", URL: "/proyectos/" + proyectoID + "/sprints"},
			Miga{Texto: "Nuevo Sprint", URL: "/proyectos/" + proyectoID + "/sprints/nuevo"},
		),
		Volver: "/proyectos/" + proyectoID + "/sprints",
	})
}

// responderErrorRuta muestra el mensaje de formato con HTTP 422 sin consultar
// la persistencia: los parámetros de la ruta se validan antes que el proyecto
// (V11).
func (a *aplicacion) responderErrorRuta(w http.ResponseWriter, proyectoID, mensaje string) {
	a.renderizar(w, http.StatusUnprocessableEntity, "error", datosPagina{
		Titulo: "Datos inválidos",
		Nav:    navConProyecto(proyectoID),
		Error:  mensaje,
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

// sprintVistaDesde arma la vista de un Sprint con los totales y las historias
// devueltos por el servicio.
func sprintVistaDesde(sprint domain.Sprint, planificados, completados, velocidad int, items []domain.ItemBacklog) sprintVista {
	historias := make([]itemSprintVista, 0, len(items))
	for _, item := range items {
		historias = append(historias, itemSprintVista{
			ID:          item.ID,
			Titulo:      item.Titulo,
			Estado:      item.Estado,
			StoryPoints: item.StoryPoints,
			Pendiente:   item.Estado == domain.EstadoItemPendiente,
		})
	}
	return sprintVista{
		Numero:         sprint.Numero,
		Objetivo:       sprint.Objetivo,
		FechaInicio:    sprint.FechaInicio.Format(formatoFecha),
		FechaFin:       sprint.FechaFin.Format(formatoFecha),
		Estado:         sprint.Estado,
		SPPlanificados: planificados,
		SPCompletados:  completados,
		Velocidad:      float64(velocidad),
		Activo:         sprint.Estado == domain.EstadoSprintActivo,
		Items:          historias,
	}
}
