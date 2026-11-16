package web

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/reportepdf"
)

// reporteProyectoPDF descarga el reporte del proyecto como PDF. Obtiene el
// nombre desde el registro del proyecto y los datos del servicio tal cual; la
// fecha de generación la provee el reloj inyectable (D15).
func (a *aplicacion) reporteProyectoPDF(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	registro, err := a.deps.Proyectos.ObtenerProyecto(id)
	if err != nil {
		a.responderErrorProyecto(w, id, err)
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

	pdf, err := reportepdf.GenerarReporteProyecto(datos, registro.Proyecto.Nombre, a.deps.Ahora())
	if err != nil {
		a.responderErrorInfraestructura(w, "generar el PDF del proyecto "+id, err)
		return
	}

	escribirDescargaPDF(w, "reporte-proyecto-"+sanitizarNombreArchivo(id)+".pdf", pdf)
}

// reporteSprintPDF descarga el reporte de un Sprint como PDF. El número de la
// ruta se valida antes de invocar el servicio (V3, R7) y la fecha de
// generación la provee el reloj inyectable (D15).
func (a *aplicacion) reporteSprintPDF(w http.ResponseWriter, r *http.Request) {
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

	pdf, err := reportepdf.GenerarReporteSprint(datos, registro.Proyecto.Nombre, a.deps.Ahora())
	if err != nil {
		a.responderErrorInfraestructura(w, "generar el PDF del sprint "+id, err)
		return
	}

	nombreArchivo := fmt.Sprintf("reporte-sprint-%s-%d.pdf", sanitizarNombreArchivo(id), numero)
	escribirDescargaPDF(w, nombreArchivo, pdf)
}

// escribirDescargaPDF escribe el PDF con las cabeceras de descarga seguras.
func escribirDescargaPDF(w http.ResponseWriter, nombreArchivo string, pdf []byte) {
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="`+nombreArchivo+`"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pdf)
}

// sanitizarNombreArchivo deja solo caracteres seguros para un nombre de
// archivo ([A-Za-z0-9._-]) y reemplaza el resto por "_" (D14).
func sanitizarNombreArchivo(valor string) string {
	var b strings.Builder
	for _, r := range valor {
		switch {
		case r >= 'A' && r <= 'Z',
			r >= 'a' && r <= 'z',
			r >= '0' && r <= '9',
			r == '.', r == '_', r == '-':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	if b.Len() == 0 {
		return "_"
	}
	return b.String()
}
