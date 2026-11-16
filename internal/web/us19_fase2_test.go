package web_test

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
	"github.com/danaguzman22/software-metrics-estimation/internal/web"
)

// servicioDashboardFalso es el doble en memoria de web.ServicioDashboardWeb.
// No deriva ningún valor: devuelve siempre lo sembrado, para detectar
// recálculos en la capa web (R19).
type servicioDashboardFalso struct {
	dashboard      service.DashboardProyecto
	err            error
	llamadas       int
	ultimoProyecto string
}

func (f *servicioDashboardFalso) ObtenerDashboard(proyectoID string) (service.DashboardProyecto, error) {
	f.llamadas++
	f.ultimoProyecto = proyectoID
	if f.err != nil {
		return service.DashboardProyecto{}, f.err
	}
	return f.dashboard, nil
}

// hoyUS19 es el reloj inyectable de los tests de US-19.
func hoyUS19() func() time.Time {
	return func() time.Time {
		return fechaUS22(2026, time.February, 1)
	}
}

// proyectoEnCursoUS19 es un proyecto con el estado almacenado en Planificado
// pero vigente a la fecha del reloj inyectable (demuestra que el estado se
// calcula por fechas y no se lee del registro, D14).
func proyectoEnCursoUS19() service.ProyectoRegistrado {
	return registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.January, 5), fechaUS22(2026, time.December, 31))
}

// dashboardConDatosUS19 siembra un dashboard con valores deliberadamente
// distintos de los que la web podría derivar de las listas, para verificar
// que la web muestra los valores del servicio sin recalcularlos (R19).
func dashboardConDatosUS19() service.DashboardProyecto {
	desv := 25.0
	desvSprint1 := -99.0  // recalcular (40→35) daría -12.50%
	desvSprint2 := 42.5   // recalcular (60→55.75) daría -7.08%
	return service.DashboardProyecto{
		SPPlanificadosTotal:            34,
		SPCompletadosTotal:             26,
		VelocidadPromedio:               77.25,
		PorcentajeHistoriasCompletadas: 12.34,
		HorasEstimadasTotal:            111.11,
		HorasRealesTotal:               222.22,
		EstimacionesCompletas:          true,
		DesviacionHorasTotal:           &desv,
		SprintsCerrados:                4,
		VelocidadPorSprint: []service.PuntoVelocidadSprint{
			{SprintNumero: 1, Velocidad: 8},
			{SprintNumero: 2, Velocidad: 13},
		},
		HorasPorSprint: []service.PuntoHorasSprint{
			{SprintNumero: 1, HorasEstimadas: 40, HorasReales: 35, EstimacionesCompletas: true, MetricasDisponibles: true, Desviacion: &desvSprint1},
			{SprintNumero: 2, HorasEstimadas: 60, HorasReales: 55.75, EstimacionesCompletas: true, MetricasDisponibles: true, Desviacion: &desvSprint2},
		},
	}
}

// dashboardEnCeroUS19 siembra un proyecto sin datos: ceros y listas nil
// (la web usa len(), nunca nil checks).
func dashboardEnCeroUS19() service.DashboardProyecto {
	return service.DashboardProyecto{
		EstimacionesCompletas: false,
	}
}

func nuevoHandlerUS19(t *testing.T, dashboard *servicioDashboardFalso, registros ...service.ProyectoRegistrado) http.Handler {
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
		Dashboard: dashboard,
		Reportes:  &servicioReportesFalso{},
		Ahora:     hoyUS19(),
	})
}

// exigirContieneUS19 verifica que el cuerpo contenga todos los fragmentos.
func exigirContieneUS19(t *testing.T, cuerpo string, fragmentos ...string) {
	t.Helper()
	for _, fragmento := range fragmentos {
		if !strings.Contains(cuerpo, fragmento) {
			t.Errorf("el fragmento %q no se encontró; cuerpo:\n%s", fragmento, cuerpo)
		}
	}
}

// svgPorTituloUS19 devuelve el bloque <svg>...</svg> cuyo contenido incluye el
// título indicado (los gráficos llevan su título en el <title> del SVG).
func svgPorTituloUS19(t *testing.T, cuerpo, titulo string) string {
	t.Helper()
	restante := cuerpo
	for {
		inicio := strings.Index(restante, "<svg")
		if inicio < 0 {
			t.Fatalf("no se encontró un SVG con %q; cuerpo:\n%s", titulo, cuerpo)
		}
		fin := strings.Index(restante[inicio:], "</svg>")
		if fin < 0 {
			t.Fatalf("el SVG que contiene %q no está cerrado con </svg>", titulo)
		}
		bloque := restante[inicio : inicio+fin]
		if strings.Contains(bloque, titulo) {
			return bloque
		}
		restante = restante[inicio+fin:]
	}
}

// R13 / CA2: el gráfico (a) "Velocidad por Sprint" es un SVG con una barra por
// Sprint cerrado, role="img", <title>, ejes con etiquetas y valores sobre las
// barras; debajo lleva una tabla de datos con <caption>.
func TestUS19_Fase2_Dashboard_GraficoVelocidadSVG(t *testing.T) {
	dashboard := &servicioDashboardFalso{dashboard: dashboardConDatosUS19()}
	handler := nuevoHandlerUS19(t, dashboard, proyectoEnCursoUS19())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	cuerpo := rec.Body.String()

	svg := svgPorTituloUS19(t, cuerpo, "Velocidad por Sprint")
	exigirContieneUS19(t, svg,
		`role="img"`,
		"<title>Velocidad por Sprint</title>",
		"<line",
		">8<",  // valor sobre la barra del Sprint 1
		">13<",  // valor sobre la barra del Sprint 2
	)
	if n := strings.Count(svg, "<rect"); n != 2 {
		t.Errorf("el gráfico de velocidad debería tener 2 barras (una por Sprint cerrado); se encontraron %d", n)
	}
	if n := strings.Count(svg, "<text"); n < 4 {
		t.Errorf("el gráfico de velocidad debería tener etiquetas numéricas en el eje; se encontraron %d <text>", n)
	}

	exigirContieneUS19(t, cuerpo, "<caption>", "Sprint", "Story Points completados")
}

// R13 / CA2: la tabla del gráfico de velocidad lista los Sprints cerrados con
// sus Story Points completados (que coinciden con la velocidad).
func TestUS19_Fase2_Dashboard_TablaVelocidad(t *testing.T) {
	dashboard := &servicioDashboardFalso{dashboard: dashboardConDatosUS19()}
	handler := nuevoHandlerUS19(t, dashboard, proyectoEnCursoUS19())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	cuerpo := rec.Body.String()

	exigirContieneUS19(t, cuerpo,
		"<caption>",
		"Story Points completados",
		">8<",
		">13<",
	)
	if strings.Count(cuerpo, "<caption>") < 2 {
		t.Errorf("cada gráfico debería llevar una tabla de datos con <caption>; se encontraron %d", strings.Count(cuerpo, "<caption>"))
	}
}

// R14 / CA2: el gráfico (b) "Horas estimadas vs reales por Sprint" es un SVG
// con barras agrupadas (estimadas #1f4e79, reales #b45309), leyenda, ejes con
// etiquetas y valores sobre las barras.
func TestUS19_Fase2_Dashboard_GraficoHorasSVG(t *testing.T) {
	dashboard := &servicioDashboardFalso{dashboard: dashboardConDatosUS19()}
	handler := nuevoHandlerUS19(t, dashboard, proyectoEnCursoUS19())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	cuerpo := rec.Body.String()

	svg := svgPorTituloUS19(t, cuerpo, "Horas estimadas vs reales por Sprint")
	exigirContieneUS19(t, svg,
		`role="img"`,
		"<title>Horas estimadas vs reales por Sprint</title>",
		"Horas estimadas",
		"Horas reales",
		"#1f4e79",
		"#b45309",
		"40.00",
		"35.00",
		"60.00",
		"55.75",
	)
	if n := strings.Count(svg, "<rect"); n != 4 {
		t.Errorf("el gráfico de horas debería tener 4 barras agrupadas (2 Sprints x 2 series); se encontraron %d", n)
	}
	if n := strings.Count(svg, "<text"); n < 4 {
		t.Errorf("el gráfico de horas debería tener etiquetas numéricas en el eje; se encontraron %d <text>", n)
	}
}

// R14 / CA2: la tabla del gráfico de horas lista Sprint, horas estimadas,
// horas reales y desviación (o "—").
func TestUS19_Fase2_Dashboard_TablaHoras(t *testing.T) {
	dashboard := &servicioDashboardFalso{dashboard: dashboardConDatosUS19()}
	handler := nuevoHandlerUS19(t, dashboard, proyectoEnCursoUS19())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	cuerpo := rec.Body.String()

	exigirContieneUS19(t, cuerpo,
		"<caption>",
		"Horas estimadas",
		"Horas reales",
		"Desviación",
		"40.00",
		"35.00",
	)
}

// R14 / R6 / E3: un Sprint cerrado sin snapshot de esfuerzo no dibuja barras y
// la tabla muestra "métricas no disponibles" para ese Sprint; el resto del
// dashboard funciona.
func TestUS19_Fase2_Dashboard_SprintSinMetricasDisponibles(t *testing.T) {
	datos := dashboardConDatosUS19()
	datos.HorasPorSprint = []service.PuntoHorasSprint{
		{SprintNumero: 1, MetricasDisponibles: false},
		{SprintNumero: 2, HorasEstimadas: 60, HorasReales: 55.75, EstimacionesCompletas: true, MetricasDisponibles: true},
	}
	dashboard := &servicioDashboardFalso{dashboard: datos}
	handler := nuevoHandlerUS19(t, dashboard, proyectoEnCursoUS19())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	cuerpo := rec.Body.String()

	exigirContieneUS19(t, cuerpo, "métricas no disponibles", "60.00", "55.75")
	exigirCuerpoNoContieneUS22(t, rec, "40.00")

	svg := svgPorTituloUS19(t, cuerpo, "Horas estimadas vs reales por Sprint")
	if n := strings.Count(svg, "<rect"); n != 2 {
		t.Errorf("el gráfico de horas debería dibujar barras solo para el Sprint con métricas (2 barras); se encontraron %d", n)
	}
}

// R16 / CA3: un proyecto sin datos responde 200 con los valores en cero, la
// desviación en "—", el mensaje "Todavía no hay Sprints cerrados" en cada
// gráfico y sin NaN ni división por cero. Las listas llegan en nil.
func TestUS19_Fase2_Dashboard_SinDatos200(t *testing.T) {
	assertos := func(t *testing.T, dashboard service.DashboardProyecto) {
		d := &servicioDashboardFalso{dashboard: dashboard}
		handler := nuevoHandlerUS19(t, d, proyectoEnCursoUS19())

		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		cuerpo := rec.Body.String()

		if n := strings.Count(cuerpo, "Todavía no hay Sprints cerrados"); n != 2 {
			t.Errorf("cada gráfico debería mostrar el mensaje \"Todavía no hay Sprints cerrados\"; se encontraron %d", n)
		}
		exigirContieneUS19(t, cuerpo,
			"0.00",    // velocidad promedio y totales
			"0.00%",   // porcentaje de historias completadas
			"0.00 h",  // horas estimadas y reales
			"—",       // desviación no disponible
			"Sprints cerrados",
		)
		exigirCuerpoNoContieneUS22(t, rec, "NaN", "Inf")
	}

	t.Run("listas nil", func(t *testing.T) {
		assertos(t, dashboardEnCeroUS19())
	})

	t.Run("listas vacías", func(t *testing.T) {
		vacias := dashboardEnCeroUS19()
		vacias.VelocidadPorSprint = []service.PuntoVelocidadSprint{}
		vacias.HorasPorSprint = []service.PuntoHorasSprint{}
		assertos(t, vacias)
	})
}

// R16 / CA3 / D10: un proyecto con Sprints solo activos no muestra cerrados en
// los gráficos; los totales sí reflejan el proyecto entero.
func TestUS19_Fase2_Dashboard_SoloActivos200(t *testing.T) {
	datos := dashboardEnCeroUS19()
	datos.SPPlanificadosTotal = 8 // SP planificados del Sprint activo
	datos.SPCompletadosTotal = 0
	dashboard := &servicioDashboardFalso{dashboard: datos}
	handler := nuevoHandlerUS19(t, dashboard, proyectoEnCursoUS19())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	cuerpo := rec.Body.String()

	exigirContieneUS19(t, cuerpo, "8", "0.00", "0.00%")
	if n := strings.Count(cuerpo, "Todavía no hay Sprints cerrados"); n != 2 {
		t.Errorf("cada gráfico debería mostrar el mensaje \"Todavía no hay Sprints cerrados\"; se encontraron %d", n)
	}
}

// R15 / CA2: la página es SVG inline generado en el servidor, sin <script> ni
// recursos externos; exactamente dos gráficos.
func TestUS19_Fase2_Dashboard_SVGSinScriptNiCDN(t *testing.T) {
	dashboard := &servicioDashboardFalso{dashboard: dashboardConDatosUS19()}
	handler := nuevoHandlerUS19(t, dashboard, proyectoEnCursoUS19())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)

	exigirCuerpoNoContieneUS22(t, rec, "<script", "http://", "https://", "src=\"http")
	if n := strings.Count(rec.Body.String(), "<svg"); n != 2 {
		t.Errorf("la página debería contener exactamente 2 gráficos SVG inline; se encontraron %d", n)
	}
}

// R20 / D19: la página no contiene secciones ni placeholders de defectos.
func TestUS19_Fase2_Dashboard_SinSeccionDefectos(t *testing.T) {
	dashboard := &servicioDashboardFalso{dashboard: dashboardConDatosUS19()}
	handler := nuevoHandlerUS19(t, dashboard, proyectoEnCursoUS19())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)

	exigirCuerpoNoContieneUS22(t, rec, "Ver defectos", "Registrar defecto", "Severidad")
}

// R19: la web no recalcula ninguna métrica; muestra los valores del servicio
// tal como llegan (el doble trae valores que no coinciden con los derivados).
func TestUS19_Fase2_Dashboard_WebNoRecalculaMetricas(t *testing.T) {
	dashboard := &servicioDashboardFalso{dashboard: dashboardConDatosUS19()}
	handler := nuevoHandlerUS19(t, dashboard, proyectoEnCursoUS19())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	cuerpo := rec.Body.String()

	exigirContieneUS19(t, cuerpo,
		"77.25",   // velocidad promedio del servicio (la media de las listas sería 10.50)
		"12.34%",  // porcentaje del servicio
		"111.11 h", // horas estimadas del servicio
		"222.22 h", // horas reales del servicio
		"34",       // SP planificados totales del servicio
		"26",       // SP completados totales del servicio
	)
	exigirCuerpoNoContieneUS22(t, rec, "10.50")
}

// R18: la desviación total se muestra con signo explícito, dos decimales y
// sufijo %.
func TestUS19_Fase2_Dashboard_DesviacionConSigno(t *testing.T) {
	t.Run("positiva", func(t *testing.T) {
		datos := dashboardConDatosUS19()
		d := 25.0
		datos.DesviacionHorasTotal = &d
		dashboard := &servicioDashboardFalso{dashboard: datos}
		handler := nuevoHandlerUS19(t, dashboard, proyectoEnCursoUS19())

		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		exigirCuerpoContieneUS22(t, rec, "+25.00%")
	})

	t.Run("negativa", func(t *testing.T) {
		datos := dashboardConDatosUS19()
		d := -25.0
		datos.DesviacionHorasTotal = &d
		dashboard := &servicioDashboardFalso{dashboard: datos}
		handler := nuevoHandlerUS19(t, dashboard, proyectoEnCursoUS19())

		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		exigirCuerpoContieneUS22(t, rec, "-25.00%")
	})
}

// R18: desviación nil → se muestra "—"; con estimaciones incompletas, además,
// una nota que indica que hay historias sin estimar en horas.
func TestUS19_Fase2_Dashboard_DesviacionNilMuestraGuion(t *testing.T) {
	t.Run("estimaciones completas con suma cero", func(t *testing.T) {
		datos := dashboardConDatosUS19()
		datos.EstimacionesCompletas = true
		datos.DesviacionHorasTotal = nil
		dashboard := &servicioDashboardFalso{dashboard: datos}
		handler := nuevoHandlerUS19(t, dashboard, proyectoEnCursoUS19())

		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		exigirContieneUS19(t, rec.Body.String(), "—")
		exigirCuerpoNoContieneUS22(t, rec, "historias sin estimar en horas")
	})

	t.Run("estimaciones incompletas", func(t *testing.T) {
		datos := dashboardConDatosUS19()
		datos.EstimacionesCompletas = false
		datos.DesviacionHorasTotal = nil
		dashboard := &servicioDashboardFalso{dashboard: datos}
		handler := nuevoHandlerUS19(t, dashboard, proyectoEnCursoUS19())

		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		exigirContieneUS19(t, rec.Body.String(), "—", "historias sin estimar en horas")
	})
}

// R17 / D1 / D18: la navegación principal agrega el enlace Dashboard apuntando
// al dashboard del proyecto en contexto; el detalle agrega "Ver dashboard".
func TestUS19_Fase2_Dashboard_EnlacesNavegacionYDetalle(t *testing.T) {
	dashboard := &servicioDashboardFalso{dashboard: dashboardConDatosUS19()}
	handler := nuevoHandlerUS19(t, dashboard, proyectoEnCursoUS19())

	t.Run("detalle del proyecto", func(t *testing.T) {
		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001", nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		exigirCuerpoContieneUS22(t, rec, "Ver dashboard", `href="/proyectos/PRJ-001/dashboard"`)
	})

	t.Run("navegación en una página del proyecto", func(t *testing.T) {
		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints", nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		exigirCuerpoContieneUS22(t, rec, "Dashboard", `href="/proyectos/PRJ-001/dashboard"`)
	})

	t.Run("los cuatro enlaces preexistentes siguen presentes", func(t *testing.T) {
		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001", nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		exigirCuerpoContieneUS22(t, rec, "Proyectos", "Backlog", "Sprints", "Defectos")
	})
}

// R17 / D18: sin proyecto en contexto (inicio y /proyectos), el enlace
// Dashboard apunta a /proyectos.
func TestUS19_Fase2_Dashboard_EnlaceDashboardSinProyecto(t *testing.T) {
	dashboard := &servicioDashboardFalso{dashboard: dashboardConDatosUS19()}
	handler := nuevoHandlerUS19(t, dashboard)

	for _, ruta := range []string{"/", "/proyectos"} {
		t.Run(ruta, func(t *testing.T) {
			rec := solicitarUS22(handler, http.MethodGet, ruta, nil)
			exigirEstadoUS22(t, rec, http.StatusOK)
			exigirCuerpoContieneUS22(t, rec, `href="/proyectos">Dashboard</a>`)
		})
	}
}

// R11 / CA1: el dashboard responde 200 con el estado calculado por fechas y
// todas las métricas del servicio tal como llegan.
func TestUS19_Fase2_Dashboard_RenderOK(t *testing.T) {
	dashboard := &servicioDashboardFalso{dashboard: dashboardConDatosUS19()}
	handler := nuevoHandlerUS19(t, dashboard, proyectoEnCursoUS19())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)

	exigirCuerpoContieneUS22(t, rec,
		"En curso",   // estado calculado por fechas (el registro guarda "Planificado")
		"34",          // Story Points planificados totales
		"26",          // Story Points completados totales
		"77.25",       // velocidad promedio con dos decimales (%.2f)
		"12.34%",      // porcentaje de historias completadas con sufijo %
		"111.11 h",    // horas estimadas totales con sufijo h
		"222.22 h",    // horas reales totales con sufijo h
		"+25.00%",     // desviación con signo explícito
		"Sprints cerrados",
	)
	if dashboard.llamadas != 1 || dashboard.ultimoProyecto != "PRJ-001" {
		t.Fatalf("dashboard: llamadas=%d proyecto=%q", dashboard.llamadas, dashboard.ultimoProyecto)
	}
}

// R12 / E1: proyecto inexistente → 404 con el mensaje del dominio; el servicio
// de dashboard no se invoca.
func TestUS19_Fase2_Dashboard_ProyectoInexistente404(t *testing.T) {
	dashboard := &servicioDashboardFalso{dashboard: dashboardConDatosUS19()}
	handler := nuevoHandlerUS19(t, dashboard)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-999/dashboard", nil)
	exigirEstadoUS22(t, rec, http.StatusNotFound)
	exigirCuerpoContieneUS22(t, rec, domain.ErrProyectoNoEncontrado.Error())
	if dashboard.llamadas != 0 {
		t.Fatalf("el servicio de dashboard no debía invocarse; llamadas=%d", dashboard.llamadas)
	}
}

// R12 / E1: un id que no corresponde a ningún proyecto (la aplicación no
// define un formato de id de proyecto; el mux acepta cualquier segmento) →
// 404 con el mensaje del dominio; el servicio de dashboard no se invoca.
func TestUS19_Fase2_Dashboard_IdInvalido404(t *testing.T) {
	dashboard := &servicioDashboardFalso{dashboard: dashboardConDatosUS19()}
	handler := nuevoHandlerUS19(t, dashboard, proyectoEnCursoUS19())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/proyecto-que-no-existe/dashboard", nil)
	exigirEstadoUS22(t, rec, http.StatusNotFound)
	exigirCuerpoContieneUS22(t, rec, domain.ErrProyectoNoEncontrado.Error())
	if dashboard.llamadas != 0 {
		t.Fatalf("el servicio de dashboard no debía invocarse; llamadas=%d", dashboard.llamadas)
	}
}

// R21 / E4: método no permitido sobre una ruta registrada → 405; el handler
// no se ejecuta.
func TestUS19_Fase2_Dashboard_MetodoNoPermitido(t *testing.T) {
	dashboard := &servicioDashboardFalso{dashboard: dashboardConDatosUS19()}
	handler := nuevoHandlerUS19(t, dashboard, proyectoEnCursoUS19())

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/dashboard", nil)
	exigirEstadoUS22(t, rec, http.StatusMethodNotAllowed)
	if dashboard.llamadas != 0 {
		t.Fatalf("el handler no debía ejecutarse; llamadas=%d", dashboard.llamadas)
	}
}

// R21 / E5: una ruta inexistente sobre un proyecto existente → 404.
func TestUS19_Fase2_Dashboard_RutaInexistente(t *testing.T) {
	dashboard := &servicioDashboardFalso{dashboard: dashboardConDatosUS19()}
	handler := nuevoHandlerUS19(t, dashboard, proyectoEnCursoUS19())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/inexistente", nil)
	exigirEstadoUS22(t, rec, http.StatusNotFound)
}

// R12 / E2: un error del servicio de dashboard → 500 con el mensaje genérico;
// el detalle solo va al log y no aparece en la respuesta.
func TestUS19_Fase2_Dashboard_ErrorDeInfraestructura500(t *testing.T) {
	dashboard := &servicioDashboardFalso{err: errors.New("falla interna: connection reset by peer")}
	handler := nuevoHandlerUS19(t, dashboard, proyectoEnCursoUS19())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
	exigirEstadoUS22(t, rec, http.StatusInternalServerError)
	exigirCuerpoContieneUS22(t, rec, "Ocurrió un error inesperado. Intentá nuevamente.")
	exigirCuerpoNoContieneUS22(t, rec, "connection reset by peer")
}

// R12 / E2: un error de infraestructura al obtener el proyecto → 500 con el
// mensaje genérico; el detalle solo va al log.
func TestUS19_Fase2_Dashboard_ErrorDeInfraestructuraAlObtenerProyecto500(t *testing.T) {
	proyectos := &servicioProyectosFalso{obtenerErr: errors.New("falla interna: timeout al leer la base")}
	dashboard := &servicioDashboardFalso{dashboard: dashboardConDatosUS19()}
	sprints := nuevoServicioSprintsFalso()
	handler := web.NuevoHandlerAplicacion(&pingFalso{}, web.Dependencias{
		Proyectos: proyectos,
		Backlog:   nuevoServicioBacklogFalso(),
		Sprints:   sprints,
		Dashboard: dashboard,
		Reportes:  &servicioReportesFalso{},
		Ahora:     hoyUS19(),
	})

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
	exigirEstadoUS22(t, rec, http.StatusInternalServerError)
	exigirCuerpoContieneUS22(t, rec, "Ocurrió un error inesperado. Intentá nuevamente.")
	exigirCuerpoNoContieneUS22(t, rec, "timeout al leer la base")
	if dashboard.llamadas != 0 {
		t.Fatalf("el servicio de dashboard no debía invocarse; llamadas=%d", dashboard.llamadas)
	}
}

// R11 / V3 / E6: fechas inconsistentes al calcular el estado del proyecto →
// 500 con el mensaje genérico; el detalle solo va al log.
func TestUS19_Fase2_Dashboard_EstadoConError500(t *testing.T) {
	registro := service.ProyectoRegistrado{
		ID: "PRJ-001",
		Proyecto: domain.Proyecto{
			Nombre:      "Fechas cruzadas",
			FechaInicio: fechaUS22(2026, time.June, 1),
			FechaFin:    fechaUS22(2026, time.January, 1),
			Estado:      domain.EstadoPlanificado,
		},
	}
	dashboard := &servicioDashboardFalso{dashboard: dashboardConDatosUS19()}
	handler := nuevoHandlerUS19(t, dashboard, registro)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
	exigirEstadoUS22(t, rec, http.StatusInternalServerError)
	exigirCuerpoContieneUS22(t, rec, "Ocurrió un error inesperado. Intentá nuevamente.")
	exigirCuerpoNoContieneUS22(t, rec, domain.ErrFechaFinAnteriorAInicio.Error())
}

// filaTablaHorasUS19 devuelve la fila <tr>...</tr> del Sprint indicado en la
// tabla del gráfico de horas (se identifica por Sprint y horas estimadas).
func filaTablaHorasUS19(t *testing.T, cuerpo string, sprint int, horasEstimadas string) string {
	t.Helper()
	inicio := strings.Index(cuerpo, fmt.Sprintf("<tr><td>%d</td><td>%s</td>", sprint, horasEstimadas))
	if inicio < 0 {
		t.Fatalf("no se encontró la fila del Sprint %d en la tabla de horas; cuerpo:\n%s", sprint, cuerpo)
	}
	fin := strings.Index(cuerpo[inicio:], "</tr>")
	if fin < 0 {
		t.Fatalf("la fila del Sprint %d no está cerrada con </tr>; cuerpo:\n%s", sprint, cuerpo)
	}
	return cuerpo[inicio : inicio+fin]
}

// R23 / CA2 (US-19, FASE 2b): la desviación por Sprint se muestra tal como la
// devuelve el servicio (dos decimales y sufijo %; el signo "+" lo escribe la
// plantilla cuando corresponde). La web no la recalcula: los valores sembrados
// no coinciden con los derivados de las horas (40→35 da -12.50% y 60→55.75 da
// -7.08%). Cuando el servicio informa nil, la tabla muestra "—".
func TestUS19_Fase2b_Dashboard_DesviacionPorSprint(t *testing.T) {
	t.Run("muestra la desviación del servicio y no la derivada", func(t *testing.T) {
		dashboard := &servicioDashboardFalso{dashboard: dashboardConDatosUS19()}
		handler := nuevoHandlerUS19(t, dashboard, proyectoEnCursoUS19())

		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		cuerpo := rec.Body.String()

		exigirContieneUS19(t, cuerpo, "-99.00%", "+42.50%")
		exigirCuerpoNoContieneUS22(t, rec, "-12.50%", "-7.08%")
	})

	t.Run("nil muestra — en la fila del Sprint", func(t *testing.T) {
		datos := dashboardConDatosUS19()
		datos.HorasPorSprint[0].Desviacion = nil
		dashboard := &servicioDashboardFalso{dashboard: datos}
		handler := nuevoHandlerUS19(t, dashboard, proyectoEnCursoUS19())

		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		cuerpo := rec.Body.String()

		fila := filaTablaHorasUS19(t, cuerpo, 1, "40.00")
		if !strings.Contains(fila, "—") {
			t.Errorf("la fila del Sprint 1 debería mostrar — para la desviación; fila: %s", fila)
		}
		if strings.Contains(fila, "%") {
			t.Errorf("la fila del Sprint 1 no debería mostrar ningún porcentaje; fila: %s", fila)
		}
		fila2 := filaTablaHorasUS19(t, cuerpo, 2, "60.00")
		if !strings.Contains(fila2, "+42.50%") {
			t.Errorf("la fila del Sprint 2 debería mostrar +42.50%%; fila: %s", fila2)
		}
	})
}

// R23 / D16 (US-19, FASE 2b): el nombre del proyecto se renderiza escapado por
// html/template; nunca se emite crudo en la página.
func TestUS19_Fase2b_Dashboard_NombreProyectoEscapado(t *testing.T) {
	registro := registroProyectoUS22(
		"PRJ-001",
		"<script>alert(1)</script>",
		fechaUS22(2026, time.January, 5),
		fechaUS22(2026, time.December, 31),
	)
	dashboard := &servicioDashboardFalso{dashboard: dashboardConDatosUS19()}
	handler := nuevoHandlerUS19(t, dashboard, registro)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	cuerpo := rec.Body.String()

	if !strings.Contains(cuerpo, "&lt;script&gt;alert(1)&lt;/script&gt;") {
		t.Errorf("el nombre del proyecto debería renderizarse escapado (&lt;script&gt;...); cuerpo:\n%s", cuerpo)
	}
	if strings.Contains(cuerpo, "<script>alert(1)</script>") {
		t.Errorf("el nombre del proyecto no debería renderizarse crudo; cuerpo:\n%s", cuerpo)
	}
}
