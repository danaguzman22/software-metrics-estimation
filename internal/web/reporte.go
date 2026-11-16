package web

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// reporteProyecto muestra el reporte imprimible del proyecto completo. El
// estado del proyecto se calcula con domain.EstadoPorFechas y el reloj
// inyectable; las cifras llegan del servicio y la web solo formatea.
func (a *aplicacion) reporteProyecto(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	registro, err := a.deps.Proyectos.ObtenerProyecto(id)
	if err != nil {
		a.responderErrorProyecto(w, id, err)
		return
	}

	estado, err := domain.EstadoPorFechas(registro.Proyecto.FechaInicio, registro.Proyecto.FechaFin, a.deps.Ahora())
	if err != nil {
		a.responderErrorInfraestructura(w, "calcular estado del proyecto "+id, err)
		return
	}

	datos, err := a.deps.Reportes.ReporteProyecto(id)
	if err != nil {
		if errors.Is(err, domain.ErrProyectoNoEncontrado) {
			a.renderizar(w, http.StatusNotFound, "error", datosPagina{
				Titulo: "Proyecto no encontrado",
				Nav:    navSinProyecto(false),
				Error:  err.Error(),
			})
			return
		}
		a.responderErrorInfraestructura(w, "generar el reporte del proyecto "+id, err)
		return
	}

	a.renderizar(w, http.StatusOK, "reporte", datosPagina{
		Titulo:  "Reporte de " + registro.Proyecto.Nombre,
		Nav:     navConProyecto(id),
		Reporte: reporteProyectoVistaDesde(datos, registro, estado),
		Migas: migasDeProyecto(registro.Proyecto.Nombre, id,
			Miga{Texto: "Reporte del proyecto", URL: "/proyectos/" + id + "/reporte"},
		),
		Volver: "/proyectos/" + id,
	})
}

// reporteSprint muestra el reporte imprimible de un Sprint. El número de la
// ruta se valida antes de invocar el servicio; el estado del proyecto no se
// muestra en esta página (D11).
func (a *aplicacion) reporteSprint(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	numero, err := strconv.Atoi(r.PathValue("numero"))
	if err != nil || numero < 1 {
		a.renderizar(w, http.StatusNotFound, "error", datosPagina{
			Titulo: "Sprint no encontrado",
			Nav:    navConProyecto(id),
			Error:  domain.ErrSprintNoEncontrado.Error(),
		})
		return
	}

	registro, err := a.deps.Proyectos.ObtenerProyecto(id)
	if err != nil {
		a.responderErrorProyecto(w, id, err)
		return
	}

	datos, err := a.deps.Reportes.ReporteSprint(id, numero)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrProyectoNoEncontrado),
			errors.Is(err, domain.ErrSprintNoEncontrado):
			a.renderizar(w, http.StatusNotFound, "error", datosPagina{
				Titulo: "Sprint no encontrado",
				Nav:    navConProyecto(id),
				Error:  err.Error(),
			})
		default:
			a.responderErrorInfraestructura(w, "generar el reporte del sprint "+id, err)
		}
		return
	}

	numeroTexto := strconv.Itoa(numero)
	a.renderizar(w, http.StatusOK, "reporte-sprint", datosPagina{
		Titulo:        "Reporte del Sprint " + numeroTexto,
		Nav:           navConProyecto(id),
		ReporteSprint: reporteSprintVistaDesde(datos, registro),
		Migas: migasDeProyecto(registro.Proyecto.Nombre, id,
			Miga{Texto: "Sprints", URL: "/proyectos/" + id + "/sprints"},
			Miga{Texto: "Reporte del Sprint " + numeroTexto, URL: "/proyectos/" + id + "/sprints/" + numeroTexto + "/reporte"},
		),
		Volver: "/proyectos/" + id + "/sprints#sprint-" + numeroTexto,
	})
}
