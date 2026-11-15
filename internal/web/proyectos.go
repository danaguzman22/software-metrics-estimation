package web

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
)

// formatoFecha es el formato de fecha usado en formularios y vistas.
const formatoFecha = "2006-01-02"

// inicio renderiza la página de inicio.
func (a *aplicacion) inicio(w http.ResponseWriter, r *http.Request) {
	a.renderizar(w, http.StatusOK, "inicio", datosPagina{
		Titulo: "Software Metrics & Estimation",
		Nav:    navSinProyecto(true),
	})
}

// listarProyectos muestra todos los proyectos con su estado calculado.
func (a *aplicacion) listarProyectos(w http.ResponseWriter, r *http.Request) {
	registros, err := a.deps.Proyectos.ListarProyectos()
	if err != nil {
		a.responderErrorInfraestructura(w, "listar proyectos", err)
		return
	}

	hoy := a.deps.Ahora()
	proyectos := make([]proyectoVista, 0, len(registros))
	for _, registro := range registros {
		vista, err := proyectoVistaDesde(registro, hoy)
		if err != nil {
			a.responderErrorInfraestructura(w, "calcular estado del proyecto", err)
			return
		}
		proyectos = append(proyectos, vista)
	}

	a.renderizar(w, http.StatusOK, "proyectos", datosPagina{
		Titulo:    "Proyectos",
		Nav:       navSinProyecto(true),
		Proyectos: proyectos,
	})
}

// formularioAltaProyecto muestra el formulario de alta de proyecto.
func (a *aplicacion) formularioAltaProyecto(w http.ResponseWriter, r *http.Request) {
	a.renderizarFormularioProyecto(w, http.StatusOK, formularioProyecto{}, "")
}

// crearProyecto valida el formato, delega el alta en el servicio y redirige al
// detalle con 303.
func (a *aplicacion) crearProyecto(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		a.renderizarFormularioProyecto(w, http.StatusUnprocessableEntity, formularioProyecto{}, "No se pudo leer el formulario.")
		return
	}

	formulario := formularioProyecto{
		Nombre:      r.PostFormValue("nombre"),
		FechaInicio: strings.TrimSpace(r.PostFormValue("fechaInicio")),
		FechaFin:    strings.TrimSpace(r.PostFormValue("fechaFin")),
	}

	fechaInicio, err := parsearFechaObligatoria(formulario.FechaInicio, "inicio")
	if err != nil {
		a.renderizarFormularioProyecto(w, http.StatusUnprocessableEntity, formulario, err.Error())
		return
	}

	fechaFin, err := parsearFechaObligatoria(formulario.FechaFin, "fin")
	if err != nil {
		a.renderizarFormularioProyecto(w, http.StatusUnprocessableEntity, formulario, err.Error())
		return
	}

	creado, err := a.deps.Proyectos.CrearProyecto(formulario.Nombre, fechaInicio, fechaFin)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNombreVacio), errors.Is(err, domain.ErrFechaFinAnteriorAInicio):
			a.renderizarFormularioProyecto(w, http.StatusUnprocessableEntity, formulario, err.Error())
		case errors.Is(err, domain.ErrProyectoNoEncontrado):
			a.renderizar(w, http.StatusNotFound, "error", datosPagina{
				Titulo: "Proyecto no encontrado",
				Nav:    navSinProyecto(false),
				Error:  err.Error(),
			})
		default:
			a.responderErrorInfraestructura(w, "crear proyecto", err)
		}
		return
	}

	w.Header().Set("Location", "/proyectos/"+creado.ID)
	w.WriteHeader(http.StatusSeeOther)
}

// detalleProyecto muestra un proyecto con su estado calculado.
func (a *aplicacion) detalleProyecto(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	registro, err := a.deps.Proyectos.ObtenerProyecto(id)
	if err != nil {
		a.responderErrorProyecto(w, id, err)
		return
	}

	vista, err := proyectoVistaDesde(registro, a.deps.Ahora())
	if err != nil {
		a.responderErrorInfraestructura(w, "calcular estado del proyecto", err)
		return
	}

	a.renderizar(w, http.StatusOK, "proyecto", datosPagina{
		Titulo:   vista.Nombre,
		Nav:      navConProyecto(id),
		Proyecto: &vista,
	})
}

// renderizarFormularioProyecto vuelve a mostrar el formulario conservando los
// datos ingresados y el mensaje de error.
func (a *aplicacion) renderizarFormularioProyecto(w http.ResponseWriter, estado int, formulario formularioProyecto, mensaje string) {
	a.renderizar(w, estado, "formulario-proyecto", datosPagina{
		Titulo:     "Nuevo proyecto",
		Nav:        navSinProyecto(false),
		Formulario: formulario,
		Error:      mensaje,
	})
}

// proyectoVistaDesde calcula el estado con domain.EstadoPorFechas y el reloj
// inyectable, sin usar el estado almacenado del proyecto.
func proyectoVistaDesde(registro service.ProyectoRegistrado, hoy time.Time) (proyectoVista, error) {
	estado, err := domain.EstadoPorFechas(registro.Proyecto.FechaInicio, registro.Proyecto.FechaFin, hoy)
	if err != nil {
		return proyectoVista{}, err
	}

	return proyectoVista{
		ID:          registro.ID,
		Nombre:      registro.Proyecto.Nombre,
		FechaInicio: registro.Proyecto.FechaInicio.Format(formatoFecha),
		FechaFin:    registro.Proyecto.FechaFin.Format(formatoFecha),
		Estado:      estado,
	}, nil
}

// parsearFechaObligatoria parsea una fecha AAAA-MM-DD con mensajes en español.
func parsearFechaObligatoria(valor, campo string) (time.Time, error) {
	if valor == "" {
		return time.Time{}, errors.New("La fecha de " + campo + " es obligatoria")
	}

	fecha, err := time.Parse(formatoFecha, valor)
	if err != nil {
		return time.Time{}, errors.New("La fecha debe tener el formato AAAA-MM-DD")
	}

	return fecha, nil
}
