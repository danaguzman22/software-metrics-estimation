package web

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// formularioItem conserva los datos ingresados en el alta y la edición de un
// ítem del backlog.
type formularioItem struct {
	Titulo      string
	Descripcion string
	Prioridad   string
	Criterios   string
	StoryPoints string
}

// listarBacklog muestra los ítems del proyecto en el orden del servicio.
func (a *aplicacion) listarBacklog(w http.ResponseWriter, r *http.Request) {
	a.renderizarBacklog(w, r.PathValue("id"), http.StatusOK, "")
}

// formularioAltaItem muestra el formulario de alta de un ítem del backlog.
func (a *aplicacion) formularioAltaItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if _, err := a.deps.Proyectos.ObtenerProyecto(id); err != nil {
		a.responderErrorProyecto(w, id, err)
		return
	}

	a.renderizarFormularioItem(w, http.StatusOK, id, "", formularioItem{}, false, "")
}

// crearItem delega el alta en el servicio y redirige al backlog con 303.
func (a *aplicacion) crearItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		a.renderizarFormularioItem(w, http.StatusUnprocessableEntity, id, "", formularioItem{}, false, "No se pudo leer el formulario.")
		return
	}

	formulario := formularioItem{
		Titulo:      r.PostFormValue("titulo"),
		Descripcion: r.PostFormValue("descripcion"),
		Prioridad:   r.PostFormValue("prioridad"),
		Criterios:   r.PostFormValue("criterios"),
	}

	_, err := a.deps.Backlog.CrearItem(
		id,
		formulario.Titulo,
		formulario.Descripcion,
		formulario.Prioridad,
		parsearCriterios(formulario.Criterios),
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTituloItemVacio),
			errors.Is(err, domain.ErrPrioridadInvalida):
			a.renderizarFormularioItem(w, http.StatusUnprocessableEntity, id, "", formulario, false, err.Error())
		case errors.Is(err, domain.ErrProyectoNoEncontrado),
			errors.Is(err, domain.ErrItemNoEncontrado):
			a.responderErrorBacklog(w, id, err)
		default:
			a.responderErrorInfraestructura(w, "crear ítem del backlog", err)
		}
		return
	}

	w.Header().Set("Location", "/proyectos/"+id+"/backlog")
	w.WriteHeader(http.StatusSeeOther)
}

// formularioEditarItem muestra el formulario precargado con los datos del ítem.
func (a *aplicacion) formularioEditarItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	itemID := r.PathValue("itemID")
	if _, ok := domain.NumeroDeItem(itemID); !ok {
		a.responderErrorRuta(w, id, "El ID de la historia debe tener el formato HU-NN")
		return
	}

	item, err := a.deps.Backlog.ObtenerItem(id, itemID)
	if err != nil {
		a.responderErrorBacklog(w, id, err)
		return
	}

	a.renderizarFormularioItem(w, http.StatusOK, id, itemID, formularioItemDesde(item), true, "")
}

// editarItem delega la edición en el servicio y redirige al backlog con 303.
// Un campo storyPoints vacío se envía como 0.
func (a *aplicacion) editarItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	itemID := r.PathValue("itemID")
	if _, ok := domain.NumeroDeItem(itemID); !ok {
		a.responderErrorRuta(w, id, "El ID de la historia debe tener el formato HU-NN")
		return
	}

	if err := r.ParseForm(); err != nil {
		a.renderizarFormularioItem(w, http.StatusUnprocessableEntity, id, itemID, formularioItem{}, true, "No se pudo leer el formulario.")
		return
	}

	formulario := formularioItem{
		Titulo:      r.PostFormValue("titulo"),
		Descripcion: r.PostFormValue("descripcion"),
		Prioridad:   r.PostFormValue("prioridad"),
		Criterios:   r.PostFormValue("criterios"),
		StoryPoints: strings.TrimSpace(r.PostFormValue("storyPoints")),
	}

	storyPoints, err := parsearStoryPointsEdicion(formulario.StoryPoints)
	if err != nil {
		a.renderizarFormularioItem(w, http.StatusUnprocessableEntity, id, itemID, formulario, true, err.Error())
		return
	}

	_, err = a.deps.Backlog.EditarItem(
		id,
		itemID,
		formulario.Titulo,
		formulario.Descripcion,
		formulario.Prioridad,
		parsearCriterios(formulario.Criterios),
		storyPoints,
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTituloItemVacio),
			errors.Is(err, domain.ErrPrioridadInvalida),
			errors.Is(err, domain.ErrStoryPointsInvalidos),
			errors.Is(err, domain.ErrItemCompletadoNoEditable):
			a.renderizarFormularioItem(w, http.StatusUnprocessableEntity, id, itemID, formulario, true, err.Error())
		case errors.Is(err, domain.ErrProyectoNoEncontrado),
			errors.Is(err, domain.ErrItemNoEncontrado):
			a.responderErrorBacklog(w, id, err)
		default:
			a.responderErrorInfraestructura(w, "editar ítem del backlog", err)
		}
		return
	}

	w.Header().Set("Location", "/proyectos/"+id+"/backlog")
	w.WriteHeader(http.StatusSeeOther)
}

// estimarItem delega la estimación en el servicio y redirige al backlog con 303.
func (a *aplicacion) estimarItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	itemID := r.PathValue("itemID")

	if err := r.ParseForm(); err != nil {
		a.renderizarBacklog(w, id, http.StatusUnprocessableEntity, "No se pudo leer el formulario.")
		return
	}

	storyPoints, err := parsearEntero(
		strings.TrimSpace(r.PostFormValue("storyPoints")),
		"Los story points deben ser un número entero",
	)
	if err != nil {
		a.renderizarBacklog(w, id, http.StatusUnprocessableEntity, err.Error())
		return
	}

	_, err = a.deps.Backlog.EstimarItem(id, itemID, storyPoints)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrStoryPointsInvalidos):
			a.renderizarBacklog(w, id, http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, domain.ErrProyectoNoEncontrado),
			errors.Is(err, domain.ErrItemNoEncontrado):
			a.responderErrorBacklog(w, id, err)
		default:
			a.responderErrorInfraestructura(w, "estimar ítem del backlog", err)
		}
		return
	}

	w.Header().Set("Location", "/proyectos/"+id+"/backlog")
	w.WriteHeader(http.StatusSeeOther)
}

// renderizarBacklog muestra el backlog del proyecto con un mensaje opcional.
func (a *aplicacion) renderizarBacklog(w http.ResponseWriter, proyectoID string, estado int, mensaje string) {
	proyecto, err := a.deps.Proyectos.ObtenerProyecto(proyectoID)
	if err != nil {
		a.responderErrorBacklog(w, proyectoID, err)
		return
	}

	items, err := a.deps.Backlog.ListarBacklog(proyectoID)
	if err != nil {
		a.responderErrorBacklog(w, proyectoID, err)
		return
	}

	a.renderizar(w, estado, "backlog", datosPagina{
		Titulo:   "Backlog",
		Nav:      navConProyecto(proyectoID),
		Proyecto: &proyectoVista{ID: proyectoID, Nombre: proyecto.Proyecto.Nombre},
		Items:    items,
		Error:    mensaje,
		Migas:    migasDeProyecto(proyecto.Proyecto.Nombre, proyectoID, Miga{Texto: "Backlog", URL: "/proyectos/" + proyectoID + "/backlog"}),
		Volver:   "/proyectos/" + proyectoID,
	})
}

// renderizarFormularioItem muestra el formulario de alta o edición de un ítem.
func (a *aplicacion) renderizarFormularioItem(w http.ResponseWriter, estado int, proyectoID, itemID string, formulario formularioItem, editando bool, mensaje string) {
	titulo := "Nuevo ítem"
	nivel := "Nuevo ítem"
	enlace := "/proyectos/" + proyectoID + "/backlog/nuevo"
	volver := "/proyectos/" + proyectoID + "/backlog"
	if editando {
		titulo = "Editar ítem"
		nivel = "Editar ítem"
		enlace = "/proyectos/" + proyectoID + "/backlog/" + itemID + "/editar"
		volver = "/proyectos/" + proyectoID + "/backlog#item-" + itemID
	}

	proyecto, err := a.deps.Proyectos.ObtenerProyecto(proyectoID)
	if err != nil {
		a.responderErrorBacklog(w, proyectoID, err)
		return
	}

	a.renderizar(w, estado, "formulario-item", datosPagina{
		Titulo:         titulo,
		Nav:            navConProyecto(proyectoID),
		Proyecto:       &proyectoVista{ID: proyectoID, Nombre: proyecto.Proyecto.Nombre},
		ItemID:         itemID,
		FormularioItem: formulario,
		Editando:       editando,
		Error:          mensaje,
		Migas: migasDeProyecto(proyecto.Proyecto.Nombre, proyectoID,
			Miga{Texto: "Backlog", URL: "/proyectos/" + proyectoID + "/backlog"},
			Miga{Texto: nivel, URL: enlace},
		),
		Volver: volver,
	})
}

// responderErrorBacklog muestra 404 con el mensaje del dominio para un proyecto
// o un ítem inexistente y 500 genérico para el resto de los errores.
func (a *aplicacion) responderErrorBacklog(w http.ResponseWriter, proyectoID string, err error) {
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
	default:
		a.responderErrorInfraestructura(w, "operar el backlog", err)
	}
}

// formularioItemDesde arma el formulario de edición con los datos del ítem.
func formularioItemDesde(item domain.ItemBacklog) formularioItem {
	return formularioItem{
		Titulo:      item.Titulo,
		Descripcion: item.Descripcion,
		Prioridad:   item.Prioridad,
		Criterios:   strings.Join(item.Criterios, "\n"),
		StoryPoints: strconv.Itoa(item.StoryPoints),
	}
}

// parsearCriterios recorta cada línea y descarta las vacías (D14).
func parsearCriterios(valor string) []string {
	lineas := strings.Split(valor, "\n")
	criterios := make([]string, 0, len(lineas))
	for _, linea := range lineas {
		linea = strings.TrimSpace(linea)
		if linea == "" {
			continue
		}
		criterios = append(criterios, linea)
	}
	return criterios
}

// parsearStoryPointsEdicion interpreta un campo vacío como 0 (D15).
func parsearStoryPointsEdicion(valor string) (int, error) {
	if valor == "" {
		return 0, nil
	}
	return parsearEntero(valor, "Los story points deben ser un número entero")
}

// parsearEntero convierte un texto en entero o devuelve el mensaje indicado.
func parsearEntero(valor, mensaje string) (int, error) {
	numero, err := strconv.Atoi(valor)
	if err != nil {
		return 0, errors.New(mensaje)
	}
	return numero, nil
}
