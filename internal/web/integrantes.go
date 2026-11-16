package web

import (
	"errors"
	"net/http"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func (a *aplicacion) paginaIntegrantes(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a.renderizarIntegrantes(w, id, http.StatusOK, "", formularioIntegrante{})
}

func (a *aplicacion) crearIntegrante(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		a.renderizarIntegrantes(w, id, http.StatusUnprocessableEntity, "No se pudo leer el formulario.", formularioIntegrante{})
		return
	}

	formulario := formularioIntegrante{
		Nombre: r.PostFormValue("nombre"),
		Rol:    r.PostFormValue("rol"),
	}
	_, err := a.deps.Proyectos.AgregarIntegrante(id, formulario.Nombre, formulario.Rol)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrProyectoNoEncontrado):
			a.renderizar(w, http.StatusNotFound, "error", datosPagina{
				Titulo: "Proyecto no encontrado",
				Nav:    navSinProyecto(false),
				Error:  err.Error(),
			})
		case errors.Is(err, domain.ErrNombreIntegranteVacio),
			errors.Is(err, domain.ErrRolIntegranteVacio),
			errors.Is(err, domain.ErrIntegranteDuplicado):
			a.renderizarIntegrantes(w, id, http.StatusUnprocessableEntity, err.Error(), formulario)
		default:
			a.responderErrorInfraestructura(w, "agregar integrante", err)
		}
		return
	}

	http.Redirect(w, r, "/proyectos/"+id+"/integrantes", http.StatusSeeOther)
}

func (a *aplicacion) renderizarIntegrantes(w http.ResponseWriter, proyectoID string, estado int, mensaje string, formulario formularioIntegrante) {
	registro, err := a.deps.Proyectos.ObtenerProyecto(proyectoID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrProyectoNoEncontrado):
			a.renderizar(w, http.StatusNotFound, "error", datosPagina{
				Titulo: "Proyecto no encontrado",
				Nav:    navSinProyecto(false),
				Error:  err.Error(),
			})
		default:
			a.responderErrorInfraestructura(w, "obtener proyecto "+proyectoID, err)
		}
		return
	}

	integrantes, err := a.deps.Proyectos.ListarIntegrantes(proyectoID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrProyectoNoEncontrado):
			a.renderizar(w, http.StatusNotFound, "error", datosPagina{
				Titulo: "Proyecto no encontrado",
				Nav:    navSinProyecto(false),
				Error:  err.Error(),
			})
		default:
			a.responderErrorInfraestructura(w, "listar integrantes", err)
		}
		return
	}

	a.renderizar(w, estado, "integrantes", datosPagina{
		Titulo:               "Integrantes",
		Nav:                  navConProyecto(proyectoID),
		Proyecto:             &proyectoVista{ID: registro.ID, Nombre: registro.Proyecto.Nombre},
		Integrantes:          integrantes,
		FormularioIntegrante: formulario,
		Error:                mensaje,
	})
}
