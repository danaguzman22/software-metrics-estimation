package web_test

import (
	"bytes"
	"errors"
	"net/http"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// pdfValidoUS21 comprueba la estructura mínima exigida por D5/R11.
func pdfValidoUS21(t *testing.T, cuerpo []byte) {
	t.Helper()
	if len(cuerpo) == 0 {
		t.Fatal("el PDF está vacío")
	}
	if !bytes.HasPrefix(cuerpo, []byte("%PDF-")) {
		t.Fatal("el PDF no comienza con la cabecera esperada")
	}
	if !bytes.HasSuffix(bytes.TrimRight(cuerpo, "\r\n"), []byte("%%"+"EOF")) {
		t.Fatal("el PDF no termina con el marcador EOF esperado")
	}
	if !bytes.Contains(cuerpo, []byte("xref")) {
		t.Fatal("el PDF no contiene xref")
	}
	if !bytes.Contains(cuerpo, []byte("startxref")) {
		t.Fatal("el PDF no contiene startxref")
	}
	if !bytes.Contains(cuerpo, []byte("/Type /Page")) && !bytes.Contains(cuerpo, []byte("/Type/Page")) {
		t.Fatal("el PDF no contiene al menos una página")
	}
}

// R1,CA1: descarga PDF del proyecto con cabeceras, nombre seguro y PDF válido.
func TestUS21_WebReporteProyectoPDF_DescargaOK(t *testing.T) {
	reportes := &servicioReportesFalso{reporteProyecto: reporteProyectoConDatosUS20()}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/reporte.pdf", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	if ct := rec.Header().Get("Content-Type"); ct != "application/pdf" {
		t.Fatalf("Content-Type=%q, se esperaba application/pdf", ct)
	}
	cd := rec.Header().Get("Content-Disposition")
	if cd != `attachment; filename="reporte-proyecto-PRJ-001.pdf"` {
		t.Fatalf("Content-Disposition=%q", cd)
	}
	pdfValidoUS21(t, rec.Body.Bytes())
	if reportes.proyectoLlamadas != 1 || reportes.ultimoProyectoID != "PRJ-001" {
		t.Fatalf("reportes: llamadas=%d proyecto=%q", reportes.proyectoLlamadas, reportes.ultimoProyectoID)
	}
}

// R2,CA1: descarga PDF del Sprint con cabeceras, nombre seguro y PDF válido.
func TestUS21_WebReporteSprintPDF_DescargaOK(t *testing.T) {
	reportes := &servicioReportesFalso{reporteSprint: reporteSprintConDatosUS20()}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints/1/reporte.pdf", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	if ct := rec.Header().Get("Content-Type"); ct != "application/pdf" {
		t.Fatalf("Content-Type=%q, se esperaba application/pdf", ct)
	}
	cd := rec.Header().Get("Content-Disposition")
	if cd != `attachment; filename="reporte-sprint-PRJ-001-1.pdf"` {
		t.Fatalf("Content-Disposition=%q", cd)
	}
	pdfValidoUS21(t, rec.Body.Bytes())
	if reportes.sprintLlamadas != 1 || reportes.ultimoSprintNumero != 1 {
		t.Fatalf("reportes sprint: llamadas=%d numero=%d", reportes.sprintLlamadas, reportes.ultimoSprintNumero)
	}
}

// R6,CA1: proyecto inexistente responde 404 con mensaje del dominio.
func TestUS21_WebReporteProyectoPDF_ProyectoInexistente404(t *testing.T) {
	reportes := &servicioReportesFalso{proyectoErr: domain.ErrProyectoNoEncontrado}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-999/reporte.pdf", nil)
	exigirEstadoUS22(t, rec, http.StatusNotFound)
	exigirCuerpoContieneUS22(t, rec, domain.ErrProyectoNoEncontrado.Error())
}

// R6,CA1: Sprint inexistente responde 404 con mensaje del dominio.
func TestUS21_WebReporteSprintPDF_SprintInexistente404(t *testing.T) {
	reportes := &servicioReportesFalso{sprintErr: domain.ErrSprintNoEncontrado}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints/9/reporte.pdf", nil)
	exigirEstadoUS22(t, rec, http.StatusNotFound)
	exigirCuerpoContieneUS22(t, rec, domain.ErrSprintNoEncontrado.Error())
}

// R7,CA1: número de Sprint inválido responde 404 sin invocar el servicio.
func TestUS21_WebReporteSprintPDF_NumeroInvalido404SinInvocarServicio(t *testing.T) {
	reportes := &servicioReportesFalso{reporteSprint: reporteSprintConDatosUS20()}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	for _, numero := range []string{"abc", "0", "-1"} {
		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints/"+numero+"/reporte.pdf", nil)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("numero %q: código=%d, se esperaba 404", numero, rec.Code)
		}
	}
	if reportes.sprintLlamadas != 0 {
		t.Fatalf("no debía invocar el servicio; llamadas=%d", reportes.sprintLlamadas)
	}
}

// R7,CA1: método no permitido responde 405.
func TestUS21_WebReportePDF_MetodoNoPermitido405(t *testing.T) {
	reportes := &servicioReportesFalso{reporteProyecto: reporteProyectoConDatosUS20(), reporteSprint: reporteSprintConDatosUS20()}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	for _, ruta := range []string{"/proyectos/PRJ-001/reporte.pdf", "/proyectos/PRJ-001/sprints/1/reporte.pdf"} {
		rec := solicitarUS22(handler, http.MethodPost, ruta, nil)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("POST %s: código=%d, se esperaba 405", ruta, rec.Code)
		}
	}
	if reportes.proyectoLlamadas != 0 || reportes.sprintLlamadas != 0 {
		t.Fatalf("no debía invocar el servicio; llamadas proyecto=%d sprint=%d", reportes.proyectoLlamadas, reportes.sprintLlamadas)
	}
}

// R8,CA1: error de infraestructura responde 500 genérico sin exponer detalle.
func TestUS21_WebReporteProyectoPDF_ErrorInfraestructura500Generico(t *testing.T) {
	reportes := &servicioReportesFalso{proyectoErr: errors.New("detalle-secreto-bd")}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/reporte.pdf", nil)
	exigirEstadoUS22(t, rec, http.StatusInternalServerError)
	exigirCuerpoContieneUS22(t, rec, "Ocurrió un error inesperado. Intentá nuevamente.")
	exigirCuerpoNoContieneUS22(t, rec, "detalle-secreto-bd")
}

// R13,CA1: las páginas HTML de reporte conservan comportamiento y agregan enlace "Descargar PDF".
func TestUS21_Web_EnlaceDescargaPDF_EnPaginasHTML(t *testing.T) {
	reportes := &servicioReportesFalso{reporteProyecto: reporteProyectoConDatosUS20(), reporteSprint: reporteSprintConDatosUS20()}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	recP := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/reporte", nil)
	exigirEstadoUS22(t, recP, http.StatusOK)
	exigirCuerpoContieneUS22(t, recP,
		`href="/proyectos/PRJ-001/reporte.pdf"`,
		"Descargar PDF",
		`class="boton-imprimir"`,
		"Imprimir",
	)

	recS := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints/1/reporte", nil)
	exigirEstadoUS22(t, recS, http.StatusOK)
	exigirCuerpoContieneUS22(t, recS,
		`href="/proyectos/PRJ-001/sprints/1/reporte.pdf"`,
		"Descargar PDF",
		`class="boton-imprimir"`,
		"Imprimir",
	)
}
