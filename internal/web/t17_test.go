package web_test

import (
	"net/http"
	"regexp"
	"strings"
	"testing"
)

// T-17 (specs/T-17.md): barra de navegación superior contextual. Fuera de un
// proyecto solo muestra "Proyectos"; dentro de un proyecto enlaza a sus
// secciones y "Defectos" apunta a los defectos del proyecto.

var reBarraT17 = regexp.MustCompile(`(?s)<nav class="navegacion".*?</nav>`)

func barraT17(t *testing.T, cuerpo string) string {
	t.Helper()
	barra := reBarraT17.FindString(cuerpo)
	if barra == "" {
		t.Fatalf("no se encontró la barra de navegación principal en:\n%s", cuerpo)
	}
	return barra
}

// CA1: fuera de un proyecto la barra muestra solo "Proyectos".
func TestT17_BarraFueraDeProyectoSoloMuestraProyectos(t *testing.T) {
	handler := nuevoEntornoT15().handler(t)
	for _, ruta := range []string{"/", "/proyectos", "/proyectos/nuevo", "/defectos"} {
		t.Run(ruta, func(t *testing.T) {
			rec := solicitarUS22(handler, http.MethodGet, ruta, nil)
			exigirEstadoUS22(t, rec, http.StatusOK)
			barra := barraT17(t, rec.Body.String())
			if !strings.Contains(barra, `href="/proyectos"`) || !strings.Contains(barra, "Proyectos") {
				t.Errorf("la barra debe enlazar a Proyectos:\n%s", barra)
			}
			for _, sobrante := range []string{"Backlog", "Sprints", "Dashboard", "Defectos", "Reporte", "Integrantes"} {
				if strings.Contains(barra, sobrante) {
					t.Errorf("fuera de un proyecto la barra no debe mostrar %q:\n%s", sobrante, barra)
				}
			}
		})
	}
}

// CA2 y CA3: dentro de un proyecto la barra enlaza a sus secciones y Defectos
// apunta a los defectos del proyecto, no al marcador global.
func TestT17_BarraDentroDeProyectoEnlazaSusSecciones(t *testing.T) {
	handler := nuevoEntornoT15().handler(t)
	rutas := []string{
		"/proyectos/PRJ-001",
		"/proyectos/PRJ-001/backlog",
		"/proyectos/PRJ-001/sprints",
		"/proyectos/PRJ-001/dashboard",
		"/proyectos/PRJ-001/defectos",
		"/proyectos/PRJ-001/reporte",
		"/proyectos/PRJ-001/integrantes",
	}
	esperados := []string{
		`href="/proyectos"`,
		`href="/proyectos/PRJ-001/backlog"`,
		`href="/proyectos/PRJ-001/sprints"`,
		`href="/proyectos/PRJ-001/dashboard"`,
		`href="/proyectos/PRJ-001/defectos"`,
		`href="/proyectos/PRJ-001/reporte"`,
		`href="/proyectos/PRJ-001/integrantes"`,
	}
	for _, ruta := range rutas {
		t.Run(ruta, func(t *testing.T) {
			rec := solicitarUS22(handler, http.MethodGet, ruta, nil)
			exigirEstadoUS22(t, rec, http.StatusOK)
			barra := barraT17(t, rec.Body.String())
			for _, fragmento := range esperados {
				if !strings.Contains(barra, fragmento) {
					t.Errorf("la barra debe contener %s:\n%s", fragmento, barra)
				}
			}
			if strings.Contains(barra, `href="/defectos"`) {
				t.Errorf("la barra no debe apuntar al marcador global /defectos:\n%s", barra)
			}
		})
	}
}

// CA4: el proyecto inexistente (404) no muestra enlaces de sección.
func TestT17_BarraEnProyectoInexistenteSoloMuestraProyectos(t *testing.T) {
	handler := nuevoEntornoT15().handler(t)
	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-999", nil)
	exigirEstadoUS22(t, rec, http.StatusNotFound)
	barra := barraT17(t, rec.Body.String())
	if strings.Contains(barra, "Backlog") || strings.Contains(barra, "Defectos") {
		t.Errorf("un proyecto inexistente no debe mostrar secciones:\n%s", barra)
	}
}

// CA5: la ruta global /defectos sigue respondiendo 200 (compatibilidad).
func TestT17_RutaGlobalDefectosSigueRespondiendo(t *testing.T) {
	handler := nuevoEntornoT15().handler(t)
	rec := solicitarUS22(handler, http.MethodGet, "/defectos", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
}

// CA6: el dashboard organiza las métricas como tarjetas y los gráficos en un
// contenedor propio, sin perder los datos ni las tablas accesibles.
func TestT17_DashboardMuestraTarjetasYContenedorDeGraficos(t *testing.T) {
	handler := nuevoEntornoT15().handler(t)
	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/dashboard", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec,
		`class="tarjetas-kpi"`, `class="tarjeta-kpi"`, `class="paneles-graficos"`,
		"Story Points planificados", "Velocidad promedio", "Horas estimadas",
		"Velocidad por Sprint", "Horas estimadas vs reales por Sprint", "<table",
	)
}
