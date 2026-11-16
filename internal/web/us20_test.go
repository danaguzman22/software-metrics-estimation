package web_test

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
	"github.com/danaguzman22/software-metrics-estimation/internal/web"
)

// servicioReportesFalso es el doble en memoria de web.ServicioReportesWeb.
// No deriva ningún valor: devuelve siempre lo sembrado, para detectar
// recálculos en la capa web (R18, R23).
type servicioReportesFalso struct {
	reporteProyecto    service.ReporteProyecto
	reporteSprint      service.ReporteSprint
	proyectoErr        error
	sprintErr          error
	proyectoLlamadas   int
	sprintLlamadas     int
	ultimoProyectoID   string
	ultimoSprintNumero int
}

func (f *servicioReportesFalso) ReporteProyecto(proyectoID string) (service.ReporteProyecto, error) {
	f.proyectoLlamadas++
	f.ultimoProyectoID = proyectoID
	if f.proyectoErr != nil {
		return service.ReporteProyecto{}, f.proyectoErr
	}
	return f.reporteProyecto, nil
}

func (f *servicioReportesFalso) ReporteSprint(proyectoID string, sprintNumero int) (service.ReporteSprint, error) {
	f.sprintLlamadas++
	f.ultimoProyectoID = proyectoID
	f.ultimoSprintNumero = sprintNumero
	if f.sprintErr != nil {
		return service.ReporteSprint{}, f.sprintErr
	}
	return f.reporteSprint, nil
}

func nuevoHandlerUS20(t *testing.T, reportes *servicioReportesFalso, registros ...service.ProyectoRegistrado) http.Handler {
	t.Helper()
	proyectos, backlog := nuevoEntornoUS22(registros...)
	sprints := nuevoServicioSprintsFalso()
	for _, registro := range proyectos.registros {
		sprints.marcarProyectoExistente(registro.ID)
	}
	return web.NuevoHandlerAplicacion(&pingFalso{}, web.Dependencias{
		Proyectos: proyectos,
		Backlog:   backlog,
		Sprints:   sprints,
		Dashboard: &servicioDashboardFalso{},
		Reportes:  reportes,
		Ahora: func() time.Time {
			return fechaUS22(2026, time.February, 1)
		},
	})
}

func proyectoEnCursoUS20() service.ProyectoRegistrado {
	return registroProyectoUS22("PRJ-001", "Sistema de Métricas",
		fechaUS22(2026, time.January, 5), fechaUS22(2026, time.December, 31))
}

// reporteProyectoConDatosUS20 siembra un reporte con valores deliberadamente
// distintos de los que la web podría derivar de las listas, para verificar
// que la web muestra los valores del servicio sin recalcularlos (R18, R23).
func reporteProyectoConDatosUS20() service.ReporteProyecto {
	desv := 25.0
	horasHU02 := 10.0
	horasHU01 := 40.0
	return service.ReporteProyecto{
		SPPlanificadosTotal:            88,
		SPCompletadosTotal:             34,
		VelocidadPromedio:              77.25,
		PorcentajeHistoriasCompletadas: 93.75, // ninguna historia completada suma 34 SP: no derivable
		HorasEstimadasTotal:            111.11,
		HorasRealesTotal:               222.22,
		EstimacionesCompletas:          false, // HU-03 sin estimar
		DesviacionHorasTotal:           &desv,
		SprintsCerrados:                4,
		SprintsCerradosVelocidad: []service.PuntoVelocidadSprint{
			{SprintNumero: 1, Velocidad: 8},
			{SprintNumero: 2, Velocidad: 13},
		},
		HistoriasPlanificadas: []service.HistoriaReporte{
			{
				ID:             "HU-02",
				Titulo:         "Registrar horas",
				Estado:         "Pendiente",
				StoryPoints:    3,
				HorasEstimadas: &horasHU02,
				HorasReales:    9.5,
			},
			{
				ID:             "HU-03",
				Titulo:         "Cerrar sesión",
				Estado:         "Pendiente",
				StoryPoints:    8,
				HorasEstimadas: nil, // sin estimar
				HorasReales:    0,
			},
		},
		HistoriasCompletadas: []service.HistoriaReporte{
			{
				ID:             "HU-01",
				Titulo:         "Iniciar sesión",
				Estado:         "Completada",
				StoryPoints:    5,
				HorasEstimadas: &horasHU01,
				HorasReales:    35.5,
			},
		},
		Defectos: service.ResumenDefectos{
			Alta:       91,
			Media:      82,
			Baja:       73,
			Abierto:    51,
			EnCurso:    42,
			Resuelto:   33,
			Total:      121,
			Detectados: 121,
			Resueltos:  117,
			Defectos: []domain.Defecto{
				{ID: 11, Descripcion: "Falla de impresión del reporte", Severidad: domain.SeveridadDefectoAlta, Estado: domain.EstadoDefectoAbierto, ItemID: "HU-01", SprintDeteccion: 1},
				{ID: 12, Descripcion: "Totales desalineados en la tabla", Severidad: domain.SeveridadDefectoMedia, Estado: domain.EstadoDefectoEnCurso, ItemID: "HU-02", SprintDeteccion: 1},
			},
		},
	}
}

// reporteProyectoEnCeroUS20 siembra un proyecto sin datos: ceros y listas nil
// (la web usa len(), nunca nil checks).
func reporteProyectoEnCeroUS20() service.ReporteProyecto {
	return service.ReporteProyecto{EstimacionesCompletas: false}
}

// reporteSprintConDatosUS20 siembra el reporte de un Sprint con valores
// distintos de los que la web podría derivar de las listas (R18, R23).
func reporteSprintConDatosUS20() service.ReporteSprint {
	desv := -12.5
	horasHU01 := 40.0
	return service.ReporteSprint{
		Numero:                1,
		Objetivo:              "Entregar el núcleo del sistema",
		FechaInicio:           fechaUS22(2026, time.February, 2),
		FechaFin:              fechaUS22(2026, time.February, 27),
		Estado:                domain.EstadoSprintCerrado,
		HistoriasAsignadas: []service.HistoriaReporte{
			{
				ID:             "HU-01",
				Titulo:         "Iniciar sesión",
				Estado:         "Completada",
				StoryPoints:    5,
				HorasEstimadas: &horasHU01,
				HorasReales:    35.5,
			},
			{
				ID:             "HU-02",
				Titulo:         "Registrar horas",
				Estado:         "Pendiente",
				StoryPoints:    3,
				HorasEstimadas: nil, // sin estimar
				HorasReales:    0,
			},
		},
		SPPlanificados:        88, // ≠ 5 + 3: no derivable
		SPCompletados:         34,
		Velocidad:             13,
		HorasEstimadas:        111.11,
		HorasReales:           222.22,
		EstimacionesCompletas: false,
		MetricasDisponibles:   true,
		DesviacionHoras:       &desv,
		Defectos: service.ResumenDefectos{
			Alta:       9,
			Media:      8,
			Baja:       7,
			Abierto:    5,
			EnCurso:    4,
			Resuelto:   3,
			Total:      24,
			Detectados: 61,
			Resueltos:  57,
			Defectos: []domain.Defecto{
				{ID: 21, Descripcion: "Falla de impresión del reporte", Severidad: domain.SeveridadDefectoAlta, Estado: domain.EstadoDefectoAbierto, ItemID: "HU-01", SprintDeteccion: 1},
				{ID: 22, Descripcion: "Encabezado cortado entre páginas", Severidad: domain.SeveridadDefectoMedia, Estado: domain.EstadoDefectoResuelto, ItemID: "HU-01", SprintDeteccion: 1, SprintResolucion: intPunteroWeb(1)},
			},
			DefectosResueltos: []domain.Defecto{
				{ID: 22, Descripcion: "Encabezado cortado entre páginas", Severidad: domain.SeveridadDefectoMedia, Estado: domain.EstadoDefectoResuelto, ItemID: "HU-01", SprintDeteccion: 1, SprintResolucion: intPunteroWeb(1)},
			},
		},
	}
}

// reporteSprintEnCeroUS20 siembra un Sprint sin datos: ceros y listas nil.
func reporteSprintEnCeroUS20() service.ReporteSprint {
	return service.ReporteSprint{Numero: 1, Estado: domain.EstadoSprintActivo}
}

func intPunteroWeb(valor int) *int {
	return &valor
}

// R11 / CA1 (US-20): la página del proyecto responde HTTP 200 con el estado
// calculado por fechas, el nombre del proyecto, las secciones de historias
// con SP, los totales, la velocidad promedio, el porcentaje, las horas con
// desviación, la tabla de Sprints cerrados, el bloque de defectos y el
// botón Imprimir.
func TestUS20_ReporteProyecto_RenderOK(t *testing.T) {
	reportes := &servicioReportesFalso{reporteProyecto: reporteProyectoConDatosUS20()}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec,
		"En curso",             // estado calculado por fechas (el registro guarda Planificado)
		"Sistema de Métricas",  // nombre del proyecto
		"Iniciar sesión",       // HU-01
		"Registrar horas",      // HU-02
		"Cerrar sesión",        // HU-03
		"Horas reales (acumuladas)",
		"77.25",                // velocidad promedio con dos decimales
		"93.75%",               // porcentaje con sufijo %
		"111.11 h",             // horas estimadas con sufijo h
		"222.22 h",             // horas reales con sufijo h
		"+25.00%",              // desviación con signo explícito
		"Sprints cerrados",
		"Historias planificadas",
		"Historias completadas",
		"Falla de impresión del reporte", // lista de defectos
	)
	exigirCuerpoNoContieneUS22(t, rec, "NaN", "Inf")
	if reportes.proyectoLlamadas != 1 || reportes.ultimoProyectoID != "PRJ-001" {
		t.Fatalf("reportes: llamadas=%d proyecto=%q", reportes.proyectoLlamadas, reportes.ultimoProyectoID)
	}
}

// R12 / CA2 (US-20): la página del Sprint responde HTTP 200 con el objetivo,
// las fechas, el estado, las historias asignadas con SP, estado y horas,
// los totales, las horas estimadas vs reales, los defectos detectados y
// resueltos y el botón Imprimir.
func TestUS20_ReporteSprint_RenderOK(t *testing.T) {
	reportes := &servicioReportesFalso{reporteSprint: reporteSprintConDatosUS20()}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints/1/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)

	exigirCuerpoContieneUS22(t, rec,
		"Entregar el núcleo del sistema", // Sprint Goal
		"2026-02-02",                      // fecha de inicio
		"2026-02-27",                      // fecha de fin
		"Cerrado",                         // estado del Sprint
		"Iniciar sesión",
		"Registrar horas",
		"Horas reales (acumuladas)",
		"111.11 h",                        // horas estimadas del Sprint
		"222.22 h",                        // horas reales del Sprint
		"-12.50%",                         // desviación con signo explícito
		"Falla de impresión del reporte",  // defectos detectados en el Sprint
	)
	exigirCuerpoNoContieneUS22(t, rec, "NaN", "Inf")
	if reportes.sprintLlamadas != 1 || reportes.ultimoProyectoID != "PRJ-001" || reportes.ultimoSprintNumero != 1 {
		t.Fatalf("reportes: llamadas=%d proyecto=%q numero=%d", reportes.sprintLlamadas, reportes.ultimoProyectoID, reportes.ultimoSprintNumero)
	}
}

// R11 / R18 (US-20): la desviación total se muestra con signo explícito,
// dos decimales y sufijo %.
func TestUS20_ReporteProyecto_DesviacionConSigno(t *testing.T) {
	t.Run("positiva", func(t *testing.T) {
		datos := reporteProyectoConDatosUS20()
		d := 25.0
		datos.DesviacionHorasTotal = &d
		reportes := &servicioReportesFalso{reporteProyecto: datos}
		handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/reporte", nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		exigirCuerpoContieneUS22(t, rec, "+25.00%")
	})

	t.Run("negativa", func(t *testing.T) {
		datos := reporteProyectoConDatosUS20()
		d := -25.0
		datos.DesviacionHorasTotal = &d
		reportes := &servicioReportesFalso{reporteProyecto: datos}
		handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/reporte", nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		exigirCuerpoContieneUS22(t, rec, "-25.00%")
	})
}

// R18 (US-20): cuando el servicio informa la desviación nil se muestra "—".
func TestUS20_ReporteProyecto_DesviacionNilMuestraGuion(t *testing.T) {
	datos := reporteProyectoConDatosUS20()
	datos.DesviacionHorasTotal = nil
	reportes := &servicioReportesFalso{reporteProyecto: datos}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "—")
	exigirCuerpoNoContieneUS22(t, rec, "+25.00%")
}

// R18 (US-20): con cero historias completadas el porcentaje se muestra como
// 0.00% y las horas como 0.00 h, sin errores ni división por cero.
func TestUS20_ReporteProyecto_PorcentajeCeroMuestraCero(t *testing.T) {
	datos := reporteProyectoEnCeroUS20()
	datos.EstimacionesCompletas = true
	reportes := &servicioReportesFalso{reporteProyecto: datos}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "0.00%", "0.00 h", "0.00")
	exigirCuerpoNoContieneUS22(t, rec, "NaN", "Inf")
}

// R13 / E1 / CA1 (US-20): un proyecto inexistente muestra la página de error
// con el mensaje del dominio y HTTP 404; el servicio de reportes no se invoca.
func TestUS20_ReporteProyecto_ProyectoInexistente404(t *testing.T) {
	reportes := &servicioReportesFalso{reporteProyecto: reporteProyectoConDatosUS20()}
	handler := nuevoHandlerUS20(t, reportes)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-999/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusNotFound)
	exigirCuerpoContieneUS22(t, rec, domain.ErrProyectoNoEncontrado.Error())
	if reportes.proyectoLlamadas != 0 {
		t.Fatalf("el servicio de reportes no debía invocarse; llamadas=%d", reportes.proyectoLlamadas)
	}
}

// R13 / E1 / E5 / CA2 (US-20): un Sprint inexistente muestra la página de
// error con el mensaje del dominio y HTTP 404.
func TestUS20_ReporteSprint_SprintInexistente404(t *testing.T) {
	reportes := &servicioReportesFalso{sprintErr: domain.ErrSprintNoEncontrado}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints/9/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusNotFound)
	exigirCuerpoContieneUS22(t, rec, domain.ErrSprintNoEncontrado.Error())
}

// R13 / E1 / CA2 (US-20): un proyecto inexistente en la ruta del reporte del
// Sprint responde 404 con el mensaje del dominio; el servicio de reportes no
// se invoca.
func TestUS20_ReporteSprint_ProyectoInexistente404(t *testing.T) {
	reportes := &servicioReportesFalso{}
	handler := nuevoHandlerUS20(t, reportes)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-999/sprints/1/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusNotFound)
	exigirCuerpoContieneUS22(t, rec, domain.ErrProyectoNoEncontrado.Error())
	if reportes.sprintLlamadas != 0 {
		t.Fatalf("el servicio de reportes no debía invocarse; llamadas=%d", reportes.sprintLlamadas)
	}
}

// R14 / V3 / CA2 (US-20): un número de Sprint inválido en la ruta responde
// HTTP 404 y no se llama al servicio.
func TestUS20_ReporteSprint_NumeroInvalido404(t *testing.T) {
	for _, numero := range []string{"abc", "0", "-1"} {
		t.Run(numero, func(t *testing.T) {
			reportes := &servicioReportesFalso{reporteSprint: reporteSprintConDatosUS20()}
			handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

			rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints/"+numero+"/reporte", nil)
			exigirEstadoUS22(t, rec, http.StatusNotFound)
			if reportes.sprintLlamadas != 0 {
				t.Fatalf("el servicio de reportes no debía invocarse; llamadas=%d", reportes.sprintLlamadas)
			}
		})
	}
}

// R13 / E2 / CA1 (US-20): un error del servicio de reportes del proyecto se
// muestra con HTTP 500 y un mensaje genérico; el detalle solo va al log y no
// aparece en la respuesta.
func TestUS20_ReporteProyecto_ErrorInfraestructura500(t *testing.T) {
	reportes := &servicioReportesFalso{proyectoErr: errors.New("falla interna: connection reset by peer")}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusInternalServerError)
	exigirCuerpoContieneUS22(t, rec, "Ocurrió un error inesperado. Intentá nuevamente.")
	exigirCuerpoNoContieneUS22(t, rec, "connection reset by peer")
}

// R13 / E2 / CA2 (US-20): un error del servicio de reportes del Sprint se
// muestra con HTTP 500 y un mensaje genérico; el detalle solo va al log y no
// aparece en la respuesta.
func TestUS20_ReporteSprint_ErrorInfraestructura500(t *testing.T) {
	reportes := &servicioReportesFalso{sprintErr: errors.New("falla interna: timeout al leer la base")}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints/1/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusInternalServerError)
	exigirCuerpoContieneUS22(t, rec, "Ocurrió un error inesperado. Intentá nuevamente.")
	exigirCuerpoNoContieneUS22(t, rec, "timeout al leer la base")
}

// R13 / V4 / E8 (US-20): un error de infraestructura al obtener el proyecto se
// muestra con HTTP 500 genérico; el detalle solo va al log.
func TestUS20_ReporteSprint_ErrorProyectoInfraestructura500(t *testing.T) {
	proyectos := &servicioProyectosFalso{obtenerErr: errors.New("falla interna: timeout al leer la base")}
	reportes := &servicioReportesFalso{reporteSprint: reporteSprintConDatosUS20()}
	sprints := nuevoServicioSprintsFalso()
	handler := web.NuevoHandlerAplicacion(&pingFalso{}, web.Dependencias{
		Proyectos: proyectos,
		Backlog:   nuevoServicioBacklogFalso(),
		Sprints:   sprints,
		Dashboard: &servicioDashboardFalso{},
		Reportes:  reportes,
		Ahora: func() time.Time {
			return fechaUS22(2026, time.February, 1)
		},
	})

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints/1/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusInternalServerError)
	exigirCuerpoContieneUS22(t, rec, "Ocurrió un error inesperado. Intentá nuevamente.")
	exigirCuerpoNoContieneUS22(t, rec, "timeout al leer la base")
	if reportes.sprintLlamadas != 0 {
		t.Fatalf("el servicio de reportes no debía invocarse; llamadas=%d", reportes.sprintLlamadas)
	}
}

// R11 / V5 (US-20): fechas inconsistentes al calcular el estado del proyecto se
// muestran con HTTP 500 genérico; el detalle solo va al log.
func TestUS20_ReporteProyecto_EstadoConError500(t *testing.T) {
	registro := service.ProyectoRegistrado{
		ID: "PRJ-001",
		Proyecto: domain.Proyecto{
			Nombre:      "Fechas cruzadas",
			FechaInicio: fechaUS22(2026, time.June, 1),
			FechaFin:    fechaUS22(2026, time.January, 1),
			Estado:      domain.EstadoPlanificado,
		},
	}
	reportes := &servicioReportesFalso{reporteProyecto: reporteProyectoConDatosUS20()}
	handler := nuevoHandlerUS20(t, reportes, registro)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusInternalServerError)
	exigirCuerpoContieneUS22(t, rec, "Ocurrió un error inesperado. Intentá nuevamente.")
	exigirCuerpoNoContieneUS22(t, rec, domain.ErrFechaFinAnteriorAInicio.Error())
	if reportes.proyectoLlamadas != 0 {
		t.Fatalf("el servicio de reportes no debía invocarse; llamadas=%d", reportes.proyectoLlamadas)
	}
}

// R15 / CA1 (US-20): un proyecto sin historias muestra "No hay historias" en
// las dos listas; responde HTTP 200 sin errores ni división por cero.
func TestUS20_ReporteProyecto_SinHistorias200(t *testing.T) {
	reportes := &servicioReportesFalso{reporteProyecto: reporteProyectoEnCeroUS20()}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	cuerpo := rec.Body.String()

	if n := strings.Count(cuerpo, "No hay historias"); n < 2 {
		t.Errorf("la página debería mostrar \"No hay historias\" en las dos listas; se encontraron %d", n)
	}
	exigirCuerpoNoContieneUS22(t, rec, "NaN", "Inf")
}

// R15 / CA1 (US-20): un proyecto sin defectos muestra "No hay defectos" con
// los conteos en cero; responde HTTP 200.
func TestUS20_ReporteProyecto_SinDefectos200(t *testing.T) {
	datos := reporteProyectoEnCeroUS20()
	datos.EstimacionesCompletas = true
	reportes := &servicioReportesFalso{reporteProyecto: datos}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "No hay defectos", "0.00%")
	exigirCuerpoNoContieneUS22(t, rec, "NaN", "Inf")
}

// R15 / CA1 (US-20): un proyecto sin Sprints cerrados muestra "Todavía no hay
// Sprints cerrados"; responde HTTP 200.
func TestUS20_ReporteProyecto_SinSprintsCerrados200(t *testing.T) {
	datos := reporteProyectoEnCeroUS20()
	datos.EstimacionesCompletas = true
	reportes := &servicioReportesFalso{reporteProyecto: datos}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "Todavía no hay Sprints cerrados")
	exigirCuerpoNoContieneUS22(t, rec, "NaN", "Inf")
}

// R15 / CA2 (US-20): un Sprint sin historias asignadas muestra "No hay
// historias asignadas" con los SP en cero; responde HTTP 200.
func TestUS20_ReporteSprint_SinHistorias200(t *testing.T) {
	reportes := &servicioReportesFalso{reporteSprint: reporteSprintEnCeroUS20()}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints/1/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "No hay historias asignadas")
	exigirCuerpoNoContieneUS22(t, rec, "NaN", "Inf")
}

// R15 / CA2 (US-20): un Sprint sin defectos muestra "No hay defectos" y
// responde HTTP 200.
func TestUS20_ReporteSprint_SinDefectos200(t *testing.T) {
	datos := reporteSprintEnCeroUS20()
	datos.MetricasDisponibles = true
	reportes := &servicioReportesFalso{reporteSprint: datos}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints/1/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "No hay defectos")
	exigirCuerpoNoContieneUS22(t, rec, "NaN", "Inf")
}

// R22 / E3 / CA2 (US-20): un Sprint cerrado sin snapshot de esfuerzo muestra
// "métricas no disponibles" en el bloque de horas y responde HTTP 200; el
// resto del reporte se muestra con normalidad.
func TestUS20_ReporteSprint_SinSnapshot200(t *testing.T) {
	datos := reporteSprintConDatosUS20()
	datos.MetricasDisponibles = false
	datos.HorasEstimadas = 0
	datos.HorasReales = 0
	datos.DesviacionHoras = nil
	reportes := &servicioReportesFalso{reporteSprint: datos}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints/1/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	cuerpo := rec.Body.String()

	exigirCuerpoContieneUS22(t, rec, "métricas no disponibles")
	// El resto del reporte funciona con normalidad.
	exigirCuerpoContieneUS22(t, rec, "Entregar el núcleo del sistema", "Iniciar sesión", "Registrar horas")
	exigirCuerpoNoContieneUS22(t, rec, "NaN", "Inf")
	if n := strings.Count(cuerpo, "métricas no disponibles"); n != 1 {
		t.Errorf("\"métricas no disponibles\" debería aparecer una vez; se encontraron %d", n)
	}
}

// R16 / D10 (US-20): las dos páginas usan la clase reporte, muestran el
// botón Imprimir con window.print() solo en pantalla y ocultan la navegación
// en impresión mediante @media print acotado a .reporte.
func TestUS20_ReporteProyecto_BotonImprimirYMediaPrint(t *testing.T) {
	t.Run("reporte del proyecto", func(t *testing.T) {
		reportes := &servicioReportesFalso{reporteProyecto: reporteProyectoConDatosUS20()}
		handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/reporte", nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		exigirCuerpoContieneUS22(t, rec, "class=\"reporte\"", "window.print()", "Imprimir", "boton-imprimir")
	})

	t.Run("reporte del Sprint", func(t *testing.T) {
		reportes := &servicioReportesFalso{reporteSprint: reporteSprintConDatosUS20()}
		handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints/1/reporte", nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		exigirCuerpoContieneUS22(t, rec, "class=\"reporte\"", "window.print()", "Imprimir", "boton-imprimir")
	})
}

// R16 / D19 (US-20): estilos.css entrega los estilos de impresión
// acotados a .reporte: oculta la navegación, el encabezado, las acciones y
// el botón Imprimir; las tablas evitan cortarse entre páginas y repiten el
// encabezado en cada página.
func TestUS20_ReporteEstilosMediaPrintPresentes(t *testing.T) {
	_, archivo, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("no se pudo obtener la ruta del test")
	}
	raizModulo := filepath.Clean(filepath.Join(filepath.Dir(archivo), "..", ".."))
	contenido, err := os.ReadFile(filepath.Join(raizModulo, "internal", "web", "estaticos", "estilos.css"))
	if err != nil {
		t.Fatalf("no se pudo leer estilos.css: %v", err)
	}
	css := string(contenido)

	indexPrint := strings.Index(css, "@media print")
	if indexPrint < 0 {
		t.Fatal("estilos.css no contiene un bloque @media print")
	}
	bloquePrint := css[indexPrint:]
	for _, fragmento := range []string{
		".reporte",
		".navegacion",
		".encabezado",
		".acciones",
		".boton-imprimir",
		"display: none",
		"page-break-inside: avoid",
		"table-header-group",
	} {
		if !strings.Contains(bloquePrint, fragmento) {
			t.Errorf("el bloque @media print no contiene %q; bloque:\n%s", fragmento, bloquePrint)
		}
	}
}

// R19 / D10 (US-20): las dos páginas no contienen recursos externos ni un
// <script> obligatorio; la única excepción tolerada es el onclick del botón
// Imprimir. El contenido dinámico se emite escapado por html/template.
func TestUS20_ReporteProyecto_SinScriptObligatorio(t *testing.T) {
	for _, caso := range []struct {
		ruta   string
		valores *servicioReportesFalso
	}{
		{ruta: "/proyectos/PRJ-001/reporte", valores: &servicioReportesFalso{reporteProyecto: reporteProyectoConDatosUS20()}},
		{ruta: "/proyectos/PRJ-001/sprints/1/reporte", valores: &servicioReportesFalso{reporteSprint: reporteSprintConDatosUS20()}},
	} {
		t.Run(caso.ruta, func(t *testing.T) {
			handler := nuevoHandlerUS20(t, caso.valores, proyectoEnCursoUS20())
			rec := solicitarUS22(handler, http.MethodGet, caso.ruta, nil)
			exigirEstadoUS22(t, rec, http.StatusOK)
			cuerpo := rec.Body.String()

			exigirCuerpoNoContieneUS22(t, rec, "http://", "https://", "src=\"http")
			if strings.Contains(cuerpo, "<script") {
				t.Errorf("la página no debería contener etiquetas <script>; cuerpo:\n%s", cuerpo)
			}
		})
	}
}

// R19 / D16 (US-20): el nombre del proyecto se renderiza escapado por
// html/template; nunca se emite crudo en la página.
func TestUS20_ReporteProyecto_NombreEscapado(t *testing.T) {
	registro := registroProyectoUS22(
		"PRJ-001",
		"<script>alert(1)</script>",
		fechaUS22(2026, time.January, 5),
		fechaUS22(2026, time.December, 31),
	)
	reportes := &servicioReportesFalso{reporteProyecto: reporteProyectoConDatosUS20()}
	handler := nuevoHandlerUS20(t, reportes, registro)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	cuerpo := rec.Body.String()

	if !strings.Contains(cuerpo, "&lt;script&gt;alert(1)&lt;/script&gt;") {
		t.Errorf("el nombre del proyecto debería renderizarse escapado (&lt;script&gt;...); cuerpo:\n%s", cuerpo)
	}
	if strings.Contains(cuerpo, "<script>alert(1)</script>") {
		t.Errorf("el nombre del proyecto no debería renderizarse crudo; cuerpo:\n%s", cuerpo)
	}
}

// R23 / CA1 / CA2 (US-20): la web no recalcula ni deriva ninguna cifra —
// porcentajes, desviaciones, conteos de defectos ni estado de las historias:
// muestra exactamente lo que devuelve el servicio (valores que no coinciden
// con ninguna combinación derivable de las listas) y solo formatea.
func TestUS20_WebReporteNoDerivaCifras(t *testing.T) {
	t.Run("proyecto", func(t *testing.T) {
		reportes := &servicioReportesFalso{reporteProyecto: reporteProyectoConDatosUS20()}
		handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)

	exigirCuerpoContieneUS22(t, rec,
			"88",      // SP planificados totales del servicio
			"34",      // SP completados totales del servicio
			"77.25",   // velocidad promedio del servicio (la media de 8 y 13 sería 10.50)
			"93.75%",  // porcentaje del servicio (1 de 3 completadas sería 33.33%)
			"111.11 h",
			"222.22 h",
			"+25.00%",
			"91",      // Alta (la lista trae 1)
			"82",      // Media (la lista trae 1)
			"73",      // Baja
			"121",     // detectados y Total (la lista trae 2)
			"117",     // resueltos
		)
		exigirCuerpoNoContieneUS22(t, rec, "10.50", "33.33%")
	})

	t.Run("sprint", func(t *testing.T) {
		reportes := &servicioReportesFalso{reporteSprint: reporteSprintConDatosUS20()}
		handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints/1/reporte", nil)
		exigirEstadoUS22(t, rec, http.StatusOK)

		exigirCuerpoContieneUS22(t, rec,
			"88",       // SP planificados del servicio (≠ 5 + 3)
			"34",       // SP completados del servicio
			"13",       // velocidad del servicio
			"111.11 h",
			"222.22 h",
			"-12.50%",  // desviación tal como la devuelve el servicio
			"61",       // detectados del servicio (la lista trae 2)
			"57",       // resueltos del servicio
			"9",        // Alta (la lista trae 1)
			"8",        // Media (la lista trae 1)
		)
	})
}

// R20 (US-20): un método no permitido sobre las rutas registradas devuelve
// HTTP 405; los handlers no se ejecutan.
func TestUS20_ReporteProyecto_MetodoNoPermitido(t *testing.T) {
	reportes := &servicioReportesFalso{reporteProyecto: reporteProyectoConDatosUS20()}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusMethodNotAllowed)
	if reportes.proyectoLlamadas != 0 {
		t.Fatalf("el handler no debía ejecutarse; llamadas=%d", reportes.proyectoLlamadas)
	}
}

// R20 (US-20): un método no permitido sobre la ruta del reporte del Sprint
// devuelve HTTP 405; el handler no se ejecuta.
func TestUS20_ReporteSprint_MetodoNoPermitido(t *testing.T) {
	reportes := &servicioReportesFalso{reporteSprint: reporteSprintConDatosUS20()}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/sprints/1/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusMethodNotAllowed)
	if reportes.sprintLlamadas != 0 {
		t.Fatalf("el handler no debía ejecutarse; llamadas=%d", reportes.sprintLlamadas)
	}
}

// R20 / E7 (US-20): una ruta inexistente devuelve HTTP 404.
func TestUS20_ReporteProyecto_RutaInexistente(t *testing.T) {
	reportes := &servicioReportesFalso{}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/inexistente", nil)
	exigirEstadoUS22(t, rec, http.StatusNotFound)
}

// R15 / R18 / CA1 (US-20): con listas nil el reporte responde HTTP 200 y no
// entra en pánico ni imprime nil/NaN; la web usa len(), nunca != nil.
func TestUS20_ReporteProyecto_ListasNil200(t *testing.T) {
	datos := reporteProyectoEnCeroUS20() // listas ya son nil
	datos.EstimacionesCompletas = true
	reportes := &servicioReportesFalso{reporteProyecto: datos}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoNoContieneUS22(t, rec, "<nil>", "NaN", "Inf")
}

// R15 / R18 / CA2 (US-20): con listas nil el reporte del Sprint responde
// HTTP 200 sin pánico ni valores inválidos.
func TestUS20_ReporteSprint_ListasNil200(t *testing.T) {
	datos := reporteSprintEnCeroUS20() // listas ya son nil
	datos.MetricasDisponibles = true
	reportes := &servicioReportesFalso{reporteSprint: datos}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints/1/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoNoContieneUS22(t, rec, "<nil>", "NaN", "Inf")
}

// R18 / D15 (US-20): Story Points enteros tal como vienen del servicio,
// historias sin estimar con "—" en las horas y nota "sin estimar".
func TestUS20_ReporteProyecto_HistoriaSinEstimarMuestraGuion(t *testing.T) {
	datos := reporteProyectoConDatosUS20()
	reportes := &servicioReportesFalso{reporteProyecto: datos}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	cuerpo := rec.Body.String()

	exigirCuerpoContieneUS22(t, rec, "—", "sin estimar", "5", "3", "8")
	if !strings.Contains(cuerpo, "35.50 h") {
		t.Errorf("las horas reales acumuladas de la historia completada deberían mostrarse con dos decimales y sufijo h; cuerpo:\n%s", cuerpo)
	}
}

// R17 / D20 / CA1 (US-20): el detalle del proyecto ofrece un enlace "Ver
// reporte" a /proyectos/{id}/reporte; los enlaces de US-22, US-15 y US-19
// siguen presentes.
func TestUS20_ReporteProyecto_EnlaceEnDetalle(t *testing.T) {
	reportes := &servicioReportesFalso{}
	handler := nuevoHandlerUS20(t, reportes, proyectoEnCursoUS20())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)

	exigirCuerpoContieneUS22(t, rec,
		"Ver reporte",
		`href="/proyectos/PRJ-001/reporte"`,
		// Enlaces preexistentes que no deben romperse (US-22, US-15, US-19).
		"Ver backlog",
		"Ver sprints",
		"Ver defectos",
		"Ver dashboard",
		`href="/proyectos/PRJ-001/dashboard"`,
	)
	if reportes.proyectoLlamadas != 0 || reportes.sprintLlamadas != 0 {
		t.Errorf("el detalle no debía invocar al servicio de reportes; llamadas=%d/%d",
			reportes.proyectoLlamadas, reportes.sprintLlamadas)
	}
}

// R17 / D20 / CA2 (US-20): el listado de Sprints ofrece un enlace "Ver reporte
// del Sprint" a /proyectos/{id}/sprints/{numero}/reporte por cada Sprint
// listado; los enlaces de US-22/US-15/US-19 siguen presentes.
func TestUS20_ReporteSprint_EnlaceEnListado(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.January, 5), fechaUS22(2026, time.December, 31)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.marcarProyectoExistente("PRJ-001")
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintCerrado,
		fechaUS22(2026, time.February, 2), fechaUS22(2026, time.February, 27),
	), 8)
	sprints.agregarSprint("PRJ-001", sprintUS22(
		2, "Segunda entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 2), fechaUS22(2026, time.March, 27),
	), 0)
	reportes := &servicioReportesFalso{}
	handler := web.NuevoHandlerAplicacion(&pingFalso{}, web.Dependencias{
		Proyectos: proyectos,
		Backlog:   backlog,
		Sprints:   sprints,
		Dashboard: &servicioDashboardFalso{},
		Reportes:  reportes,
		Ahora: func() time.Time {
			return fechaUS22(2026, time.February, 1)
		},
	})

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)

	exigirCuerpoContieneUS22(t, rec,
		"Ver reporte del Sprint",
		`href="/proyectos/PRJ-001/sprints/1/reporte"`,
		`href="/proyectos/PRJ-001/sprints/2/reporte"`,
		// Enlaces preexistentes que no deben romperse (US-22).
		"Histórico de Sprints",
		`href="/proyectos/PRJ-001/sprints/historico"`,
		"Nuevo Sprint",
	)
	if reportes.proyectoLlamadas != 0 || reportes.sprintLlamadas != 0 {
		t.Errorf("el listado no debía invocar al servicio de reportes; llamadas=%d/%d",
			reportes.proyectoLlamadas, reportes.sprintLlamadas)
	}
}


