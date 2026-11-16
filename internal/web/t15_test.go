package web_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/web"
)

// T-15 (specs/T-15.md, features/T-15.feature): navegación jerárquica con
// migas de pan y enlace "← Volver" en todas las pantallas. Los tests de este
// archivo verifican la tabla de la sección 11 de la spec (CA1-CA6).

// entornoT15 reúne los dobles en memoria con los que se ejercitan las
// pantallas de la jerarquía: proyecto PRJ-001 con la historia HU-01, el
// Sprint número 1 y defectos, dashboard, reportes y Planning Poker.
type entornoT15 struct {
	proyectos *servicioProyectosFalso
	backlog   *servicioBacklogFalso
	sprints   *servicioSprintsFalso
	defectos  *servicioDefectosFalsoUS15
	dashboard *servicioDashboardFalso
	reportes  *servicioReportesFalso
	poker     *servicioPokerFalso
}

const nombreProyectoT15 = "Sistema de Métricas"

func nuevoEntornoT15() *entornoT15 {
	inicio := fechaUS22(2026, time.January, 5)
	fin := fechaUS22(2026, time.December, 31)
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", nombreProyectoT15, inicio, fin),
		registroProyectoUS22("PRJ-002", `<script>alert(1)</script>`, inicio, fin),
		registroProyectoUS22("PRJ-003", "", inicio, fin),
		registroProyectoUS22("PRJ-004", "Sistema de métricas y estimación para equipos distribuidos en múltiples zonas horarias con reportes semanales consolidados", inicio, fin),
	)
	backlog.agregarItem("PRJ-001", itemUS22("HU-01", "Iniciar sesión", "Media", "Pendiente", 3))

	sprints := nuevoServicioSprintsFalso()
	for _, registro := range proyectos.registros {
		sprints.marcarProyectoExistente(registro.ID)
	}
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 8)
	sprints.itemsPorSprint["PRJ-001"] = map[int][]domain.ItemBacklog{
		1: {itemUS22("HU-01", "Iniciar sesión", "Media", "Pendiente", 3)},
	}

	return &entornoT15{
		proyectos: proyectos,
		backlog:   backlog,
		sprints:   sprints,
		defectos:  &servicioDefectosFalsoUS15{},
		dashboard: &servicioDashboardFalso{dashboard: dashboardConDatosUS19()},
		reportes: &servicioReportesFalso{
			reporteProyecto: reporteProyectoConDatosUS20(),
			reporteSprint:   reporteSprintConDatosUS20(),
		},
		poker: &servicioPokerFalso{},
	}
}

func (e *entornoT15) handler(t *testing.T) http.Handler {
	t.Helper()
	return web.NuevoHandlerAplicacion(&pingFalso{}, web.Dependencias{
		Proyectos: e.proyectos,
		Backlog:   e.backlog,
		Sprints:   e.sprints,
		Poker:     e.poker,
		Defectos:  e.defectos,
		Dashboard: e.dashboard,
		Reportes:  e.reportes,
		Ahora: func() time.Time {
			return fechaUS22(2026, time.February, 1)
		},
	})
}

// migaEsperadaT15 describe un nivel esperado de la migaja de pan.
type migaEsperadaT15 struct {
	texto string
	url   string
}

// pantallaT15 describe una fila de la tabla de la sección 11 de la spec:
// ruta, migas esperadas (en orden) y destino del "← Volver" (vacío cuando la
// pantalla no tiene Volver).
type pantallaT15 struct {
	ruta   string
	migas  []migaEsperadaT15
	volver string
}

func tablaPantallasT15() []pantallaT15 {
	proyectos := migaEsperadaT15{"Proyectos", "/proyectos"}
	nombre := migaEsperadaT15{nombreProyectoT15, "/proyectos/PRJ-001"}
	backlog := migaEsperadaT15{"Backlog", "/proyectos/PRJ-001/backlog"}
	sprints := migaEsperadaT15{"Sprints", "/proyectos/PRJ-001/sprints"}
	return []pantallaT15{
		{ruta: "/proyectos", migas: []migaEsperadaT15{proyectos}, volver: ""},
		{ruta: "/proyectos/nuevo", migas: []migaEsperadaT15{proyectos, {"Nuevo proyecto", "/proyectos/nuevo"}}, volver: "/proyectos"},
		{ruta: "/proyectos/PRJ-001", migas: []migaEsperadaT15{proyectos, nombre}, volver: "/proyectos"},
		{ruta: "/proyectos/PRJ-001/backlog", migas: []migaEsperadaT15{proyectos, nombre, backlog}, volver: "/proyectos/PRJ-001"},
		{ruta: "/proyectos/PRJ-001/backlog/nuevo", migas: []migaEsperadaT15{proyectos, nombre, backlog, {"Nuevo ítem", "/proyectos/PRJ-001/backlog/nuevo"}}, volver: "/proyectos/PRJ-001/backlog"},
		{ruta: "/proyectos/PRJ-001/backlog/HU-01/editar", migas: []migaEsperadaT15{proyectos, nombre, backlog, {"Editar ítem", "/proyectos/PRJ-001/backlog/HU-01/editar"}}, volver: "/proyectos/PRJ-001/backlog#item-HU-01"},
		{ruta: "/proyectos/PRJ-001/backlog/HU-01/poker", migas: []migaEsperadaT15{proyectos, nombre, backlog, {"Planning Poker", "/proyectos/PRJ-001/backlog/HU-01/poker"}}, volver: "/proyectos/PRJ-001/backlog#item-HU-01"},
		{ruta: "/proyectos/PRJ-001/sprints", migas: []migaEsperadaT15{proyectos, nombre, sprints}, volver: "/proyectos/PRJ-001"},
		{ruta: "/proyectos/PRJ-001/sprints/nuevo", migas: []migaEsperadaT15{proyectos, nombre, sprints, {"Nuevo Sprint", "/proyectos/PRJ-001/sprints/nuevo"}}, volver: "/proyectos/PRJ-001/sprints"},
		{ruta: "/proyectos/PRJ-001/sprints/historico", migas: []migaEsperadaT15{proyectos, nombre, sprints, {"Histórico", "/proyectos/PRJ-001/sprints/historico"}}, volver: "/proyectos/PRJ-001/sprints"},
		{ruta: "/proyectos/PRJ-001/sprints/1/reporte", migas: []migaEsperadaT15{proyectos, nombre, sprints, {"Reporte del Sprint 1", "/proyectos/PRJ-001/sprints/1/reporte"}}, volver: "/proyectos/PRJ-001/sprints#sprint-1"},
		{ruta: "/proyectos/PRJ-001/defectos", migas: []migaEsperadaT15{proyectos, nombre, {"Defectos", "/proyectos/PRJ-001/defectos"}}, volver: "/proyectos/PRJ-001"},
		{ruta: "/proyectos/PRJ-001/integrantes", migas: []migaEsperadaT15{proyectos, nombre, {"Integrantes", "/proyectos/PRJ-001/integrantes"}}, volver: "/proyectos/PRJ-001"},
		{ruta: "/proyectos/PRJ-001/dashboard", migas: []migaEsperadaT15{proyectos, nombre, {"Dashboard", "/proyectos/PRJ-001/dashboard"}}, volver: "/proyectos/PRJ-001"},
		{ruta: "/proyectos/PRJ-001/reporte", migas: []migaEsperadaT15{proyectos, nombre, {"Reporte del proyecto", "/proyectos/PRJ-001/reporte"}}, volver: "/proyectos/PRJ-001"},
	}
}

var rxHrefT15 = regexp.MustCompile(`href="([^"]*)"`)

// bloqueMigasT15 extrae el bloque <nav class="migas">…</nav> del cuerpo.
func bloqueMigasT15(t *testing.T, cuerpo string) string {
	t.Helper()
	const apertura = `<nav class="migas"`
	inicio := strings.Index(cuerpo, apertura)
	if inicio < 0 {
		t.Fatalf("la página no contiene el bloque %q; cuerpo:\n%s", apertura, cuerpo)
	}
	resto := cuerpo[inicio:]
	fin := strings.Index(resto, "</nav>")
	if fin < 0 {
		t.Fatalf("no se encontró el cierre </nav> del bloque de migas; cuerpo:\n%s", cuerpo)
	}
	return resto[:fin+len("</nav>")]
}

// enlacesDelBloqueT15 devuelve las etiquetas <a …>…</a> del bloque, en orden.
func enlacesDelBloqueT15(bloque string) []string {
	var enlaces []string
	resto := bloque
	for {
		ini := strings.Index(resto, "<a ")
		if ini < 0 {
			return enlaces
		}
		fin := strings.Index(resto[ini:], "</a>")
		if fin < 0 {
			return enlaces
		}
		enlaces = append(enlaces, resto[ini:ini+fin+len("</a>")])
		resto = resto[ini+fin+len("</a>"):]
	}
}

func exigirMigasT15(t *testing.T, cuerpo string, esperadas []migaEsperadaT15) {
	t.Helper()
	bloque := bloqueMigasT15(t, cuerpo)
	if n := strings.Count(cuerpo, `<nav class="migas"`); n != 1 {
		t.Errorf("el bloque de migas aparece %d veces, se esperaba 1 (R1)", n)
	}
	enlaces := enlacesDelBloqueT15(bloque)
	if len(enlaces) != len(esperadas) {
		t.Fatalf("la migaja tiene %d enlaces, se esperaban %d (%v); bloque:\n%s",
			len(enlaces), len(esperadas), esperadas, bloque)
	}
	for i, enlace := range enlaces {
		esperada := esperadas[i]
		coincidencia := rxHrefT15.FindStringSubmatch(enlace)
		if coincidencia == nil {
			t.Errorf("el enlace %d no tiene href; enlace: %s", i, enlace)
			continue
		}
		if coincidencia[1] != esperada.url {
			t.Errorf("enlace %d: href = %q, se esperaba %q", i, coincidencia[1], esperada.url)
		}
		inicioTexto := strings.Index(enlace, ">")
		texto := ""
		if inicioTexto >= 0 {
			texto = strings.TrimSuffix(enlace[inicioTexto+1:], "</a>")
		}
		if texto != esperada.texto {
			t.Errorf("enlace %d: texto = %q, se esperaba %q", i, texto, esperada.texto)
		}
		vigente := strings.Contains(enlace, `aria-current="page"`)
		if i == len(enlaces)-1 {
			if !vigente {
				t.Errorf("el nivel vigente (%q) no lleva aria-current=\"page\" (R2); enlace: %s", esperada.texto, enlace)
			}
		} else if vigente {
			t.Errorf("el enlace %q lleva aria-current=\"page\" pero no es el nivel vigente (R2); enlace: %s", esperada.texto, enlace)
		}
	}
	if len(esperadas) > 1 && !strings.Contains(bloque, "›") {
		t.Errorf("la migaja no contiene el separador › entre niveles (D11); bloque:\n%s", bloque)
	}
}

func exigirVolverT15(t *testing.T, cuerpo, destino string) {
	t.Helper()
	if destino == "" {
		if strings.Contains(cuerpo, "← Volver") {
			t.Errorf("la pantalla no debería mostrar el enlace ← Volver; cuerpo:\n%s", cuerpo)
		}
		return
	}
	quiero := `href="` + destino + `">← Volver</a>`
	if !strings.Contains(cuerpo, quiero) {
		t.Errorf("el enlace ← Volver no apunta a %q; cuerpo:\n%s", destino, cuerpo)
	}
}

// CA1 / R1 / R2 / R3 / R4: la tabla de la sección 11 de la spec — las migas
// de cada pantalla, con todos los niveles clicables y en orden.
func TestT15_MigasEnCadaPantalla(t *testing.T) {
	entorno := nuevoEntornoT15()
	handler := entorno.handler(t)

	for _, caso := range tablaPantallasT15() {
		t.Run(caso.ruta, func(t *testing.T) {
			rec := solicitarUS22(handler, http.MethodGet, caso.ruta, nil)
			exigirEstadoUS22(t, rec, http.StatusOK)
			exigirMigasT15(t, rec.Body.String(), caso.migas)
		})
	}
}

// CA2 / R4: el "← Volver" de cada pantalla apunta al nivel padre inmediato,
// incluidos los desvíos anclados al ítem y al Sprint.
func TestT15_VolverEnCadaPantalla(t *testing.T) {
	entorno := nuevoEntornoT15()
	handler := entorno.handler(t)

	for _, caso := range tablaPantallasT15() {
		t.Run(caso.ruta, func(t *testing.T) {
			rec := solicitarUS22(handler, http.MethodGet, caso.ruta, nil)
			exigirEstadoUS22(t, rec, http.StatusOK)
			exigirVolverT15(t, rec.Body.String(), caso.volver)
		})
	}
}

// CA1 / R1: las migas se renderizan una única vez por página y solo desde
// base.html.
func TestT15_MigasSeRenderizanUnaVez(t *testing.T) {
	entorno := nuevoEntornoT15()
	handler := entorno.handler(t)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/backlog", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	if n := strings.Count(rec.Body.String(), `<nav class="migas"`); n != 1 {
		t.Errorf("el bloque de migas aparece %d veces, se esperaba exactamente 1 (R1)", n)
	}
}

// CA1 / D12: las páginas fuera de la jerarquía (inicio y placeholder global
// de defectos) no muestran migas ni "← Volver".
func TestT15_FueraDeJerarquiaSinMigas(t *testing.T) {
	entorno := nuevoEntornoT15()
	handler := entorno.handler(t)

	for _, ruta := range []string{"/", "/defectos"} {
		t.Run(ruta, func(t *testing.T) {
			rec := solicitarUS22(handler, http.MethodGet, ruta, nil)
			exigirEstadoUS22(t, rec, http.StatusOK)
			exigirCuerpoNoContieneUS22(t, rec, `<nav class="migas"`, "← Volver")
		})
	}
}

// Tanda 2: anclas y regresión.

// CA2 / R5 / D14: cada fila del backlog lleva el ancla id="item-{ID}" (ID
// completo) para que Poker y Editar vuelvan al lugar exacto.
func TestT15_BacklogFilasConAncla(t *testing.T) {
	entorno := nuevoEntornoT15()
	handler := entorno.handler(t)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/backlog", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, `<tr id="item-HU-01">`)
}

// CA2 / R6 / C7: el encabezado por Sprint lleva el ancla id="sprint-{numero}".
func TestT15_SprintsEncabezadoConAncla(t *testing.T) {
	entorno := nuevoEntornoT15()
	handler := entorno.handler(t)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, `<h2 id="sprint-1">`, "Sprint 1: historias asignadas")
}

// CA2 / R4 / D14: el "← Volver" del Planning Poker apunta al backlog anclado
// al ítem.
func TestT15_PokerVolverAnclaAlItem(t *testing.T) {
	entorno := nuevoEntornoT15()
	handler := entorno.handler(t)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/backlog/HU-01/poker", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirVolverT15(t, rec.Body.String(), "/proyectos/PRJ-001/backlog#item-HU-01")
}

// CA2 / R4 / D14: el "← Volver" de Editar ítem apunta al backlog anclado al
// ítem.
func TestT15_EditarVolverAnclaAlItem(t *testing.T) {
	entorno := nuevoEntornoT15()
	handler := entorno.handler(t)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/backlog/HU-01/editar", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirVolverT15(t, rec.Body.String(), "/proyectos/PRJ-001/backlog#item-HU-01")
}

// CA2 / R4 / R6: el "← Volver" del Reporte del Sprint apunta a Sprints
// anclado al Sprint.
func TestT15_ReporteSprintVolverAncladoAlSprint(t *testing.T) {
	entorno := nuevoEntornoT15()
	handler := entorno.handler(t)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints/1/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirVolverT15(t, rec.Body.String(), "/proyectos/PRJ-001/sprints#sprint-1")
}

// CA1 / CA6 / R12: el 422 que re-renderiza Editar ítem muestra las mismas
// migas y el mismo "← Volver" que la pantalla en GET (escenario @CA1 @CA6).
func TestT15_EditarItem422MuestraMismasMigas(t *testing.T) {
	entorno := nuevoEntornoT15()
	entorno.backlog.editarErr = domain.ErrTituloItemVacio
	handler := entorno.handler(t)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/editar", url.Values{
		"titulo":      {""},
		"prioridad":   {"Media"},
		"descripcion": {"Descripción"},
	})
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirMigasT15(t, rec.Body.String(), []migaEsperadaT15{
		{"Proyectos", "/proyectos"},
		{nombreProyectoT15, "/proyectos/PRJ-001"},
		{"Backlog", "/proyectos/PRJ-001/backlog"},
		{"Editar ítem", "/proyectos/PRJ-001/backlog/HU-01/editar"},
	})
	exigirVolverT15(t, rec.Body.String(), "/proyectos/PRJ-001/backlog#item-HU-01")
}

// CA1 / CA6 / R12 / C2: la estimación es una acción inline; el 422
// re-renderiza el Backlog con las migas del Backlog y Volver a /proyectos/{id}.
func TestT15_Estimar422MuestraMigasDelBacklog(t *testing.T) {
	entorno := nuevoEntornoT15()
	handler := entorno.handler(t)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/estimar", url.Values{
		"storyPoints": {"abc"},
	})
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirMigasT15(t, rec.Body.String(), []migaEsperadaT15{
		{"Proyectos", "/proyectos"},
		{nombreProyectoT15, "/proyectos/PRJ-001"},
		{"Backlog", "/proyectos/PRJ-001/backlog"},
	})
	exigirVolverT15(t, rec.Body.String(), "/proyectos/PRJ-001")
}

// CA1 / CA6 / R12 / C8: el 422 del alta de proyecto muestra las migas
// Proyectos › Nuevo proyecto y Volver a /proyectos.
func TestT15_FormularioProyecto422MuestraMigas(t *testing.T) {
	entorno := nuevoEntornoT15()
	entorno.proyectos.crearErr = domain.ErrNombreVacio
	handler := entorno.handler(t)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos", url.Values{
		"nombre":      {"   "},
		"fechaInicio": {"2026-01-05"},
		"fechaFin":    {"2026-12-31"},
	})
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirMigasT15(t, rec.Body.String(), []migaEsperadaT15{
		{"Proyectos", "/proyectos"},
		{"Nuevo proyecto", "/proyectos/nuevo"},
	})
	exigirVolverT15(t, rec.Body.String(), "/proyectos")
}

// CA4 / R7 / D6 / E1: un proyecto inexistente muestra el error sin migas ni
// "← Volver".
func TestT15_ErrorProyectoInexistenteSinMigas(t *testing.T) {
	entorno := nuevoEntornoT15()
	handler := entorno.handler(t)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-999/backlog", nil)
	exigirEstadoUS22(t, rec, http.StatusNotFound)
	exigirCuerpoContieneUS22(t, rec, domain.ErrProyectoNoEncontrado.Error())
	exigirCuerpoNoContieneUS22(t, rec, `<nav class="migas"`, "← Volver")
}

// CA4 / R7 / D6 / E2: un ítem inexistente muestra el error sin migas ni
// "← Volver".
func TestT15_ErrorItemInexistenteSinMigas(t *testing.T) {
	entorno := nuevoEntornoT15()
	handler := entorno.handler(t)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/backlog/HU-99/editar", nil)
	exigirEstadoUS22(t, rec, http.StatusNotFound)
	exigirCuerpoContieneUS22(t, rec, domain.ErrItemNoEncontrado.Error())
	exigirCuerpoNoContieneUS22(t, rec, `<nav class="migas"`, "← Volver")
}

// CA4 / R7 / D6 / E3: un Sprint inexistente muestra el error sin migas ni
// "← Volver".
func TestT15_ErrorSprintInexistenteSinMigas(t *testing.T) {
	entorno := nuevoEntornoT15()
	entorno.reportes.sprintErr = domain.ErrSprintNoEncontrado
	handler := entorno.handler(t)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints/9/reporte", nil)
	exigirEstadoUS22(t, rec, http.StatusNotFound)
	exigirCuerpoContieneUS22(t, rec, domain.ErrSprintNoEncontrado.Error())
	exigirCuerpoNoContieneUS22(t, rec, `<nav class="migas"`, "← Volver")
}

// CA4 / R7 / D15 / E4: un número de Sprint inválido muestra el error sin
// migas ni "← Volver".
func TestT15_ErrorSprintNumeroInvalidoSinMigas(t *testing.T) {
	entorno := nuevoEntornoT15()
	handler := entorno.handler(t)

	for _, numero := range []string{"abc", "0", "-1"} {
		t.Run(numero, func(t *testing.T) {
			rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints/"+numero+"/reporte", nil)
			exigirEstadoUS22(t, rec, http.StatusNotFound)
			exigirCuerpoNoContieneUS22(t, rec, `<nav class="migas"`, "← Volver")
		})
	}
}

// CA4 / R7 / D15 / E5: un {itemID} con formato inválido en el GET de editar
// muestra el error de formato sin migas ni "← Volver", igual que poker.
func TestT15_EditarItemFormatoInvalidoSinMigas(t *testing.T) {
	entorno := nuevoEntornoT15()
	handler := entorno.handler(t)

	for _, itemID := range []string{"hu-01", "HU-", "HU-abc", "HU-01x"} {
		t.Run(itemID, func(t *testing.T) {
			rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/backlog/"+itemID+"/editar", nil)
			exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
			exigirCuerpoContieneUS22(t, rec, "El ID de la historia debe tener el formato HU-NN")
			exigirCuerpoNoContieneUS22(t, rec, `<nav class="migas"`, "← Volver")
		})
	}
}

// CA4 / R7 / D15 / E5: un {itemID} con formato inválido en el POST de editar
// muestra el error de formato sin migas ni "← Volver", igual que poker.
func TestT15_EditarItemPostFormatoInvalidoSinMigas(t *testing.T) {
	entorno := nuevoEntornoT15()
	handler := entorno.handler(t)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/hu-01/editar", url.Values{
		"titulo":    {"Título"},
		"prioridad": {"Media"},
	})
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, "El ID de la historia debe tener el formato HU-NN")
	exigirCuerpoNoContieneUS22(t, rec, `<nav class="migas"`, "← Volver")
}

// CA4 / R7 / D15 / E4: un {numero} de Sprint no numerico en la ruta de
// completar/cerrar es un error de formato de ruta: error.html 422 sin migas ni
// "← Volver" (R12 alcanza solo a los 422 de dominio que re-renderizan Sprints).
func TestT15_SprintNumeroRutaInvalidoSinMigas(t *testing.T) {
	casos := []struct {
		nombre string
		ruta   string
	}{
		{nombre: "completar", ruta: "/proyectos/PRJ-001/sprints/abc/items/HU-01/completar"},
		{nombre: "cerrar", ruta: "/proyectos/PRJ-001/sprints/abc/cerrar"},
		{nombre: "asignar", ruta: "/proyectos/PRJ-001/sprints/abc/asignar"},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			entorno := nuevoEntornoT15()
			handler := entorno.handler(t)

			rec := solicitarUS22(handler, http.MethodPost, caso.ruta, url.Values{"itemID": {"HU-01"}})
			exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
			exigirCuerpoContieneUS22(t, rec, "El número del sprint debe ser un número entero")
			exigirCuerpoNoContieneUS22(t, rec, `<nav class="migas"`, "← Volver")
		})
	}
}

// CA4 / R7 / D15 / E5: un {itemID} con formato inválido muestra el error de
// formato sin migas ni "← Volver".
func TestT15_ErrorItemFormatoInvalidoSinMigas(t *testing.T) {
	entorno := nuevoEntornoT15()
	handler := entorno.handler(t)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/backlog/hu-01/poker", nil)
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoNoContieneUS22(t, rec, `<nav class="migas"`, "← Volver")
}

// CA4 / R7 / E6: un error de infraestructura muestra 500 con el mensaje
// genérico, sin migas ni "← Volver" ni detalle interno.
func TestT15_Error500SinMigas(t *testing.T) {
	entorno := nuevoEntornoT15()
	detalle := "detalle-t15-postgres-caida"
	entorno.backlog.listarErr = errors.New(detalle)
	handler := entorno.handler(t)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/backlog", nil)
	exigirEstadoUS22(t, rec, http.StatusInternalServerError)
	exigirCuerpoContieneUS22(t, rec, "Ocurrió un error inesperado. Intentá nuevamente.")
	exigirCuerpoNoContieneUS22(t, rec, `<nav class="migas"`, "← Volver", detalle)
}

// CA5 / R9 / D13: el nombre del proyecto con markup se emite escapado en el
// enlace de la miga; nunca se interpreta como HTML (patrón de T-14).
func TestT15_MigasEscapadas(t *testing.T) {
	entorno := nuevoEntornoT15()
	handler := entorno.handler(t)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-002", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	bloque := bloqueMigasT15(t, rec.Body.String())
	if !strings.Contains(bloque, "&lt;script&gt;alert(1)&lt;/script&gt;") {
		t.Errorf("la miga del proyecto no muestra el nombre escapado; bloque:\n%s", bloque)
	}
	if strings.Contains(bloque, "<script>alert(1)</script>") {
		t.Errorf("la miga interpreta el nombre como HTML; bloque:\n%s", bloque)
	}
}

// CA5 / R9 / D13: un nombre vacío muestra el texto genérico "Proyecto".
func TestT15_NombreVacioMuestraProyecto(t *testing.T) {
	entorno := nuevoEntornoT15()
	handler := entorno.handler(t)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-003", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	bloque := bloqueMigasT15(t, rec.Body.String())
	if !strings.Contains(bloque, ">Proyecto<") {
		t.Errorf("la miga no muestra el texto genérico %q; bloque:\n%s", "Proyecto", bloque)
	}
}

// CA5 / R9 / D13: un nombre muy largo no se trunca; se muestra completo en
// la miga.
func TestT15_NombreMuyLargoSeMuestraCompleto(t *testing.T) {
	entorno := nuevoEntornoT15()
	handler := entorno.handler(t)

	nombre := "Sistema de métricas y estimación para equipos distribuidos en múltiples zonas horarias con reportes semanales consolidados"
	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-004", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	bloque := bloqueMigasT15(t, rec.Body.String())
	if !strings.Contains(bloque, nombre) {
		t.Errorf("la miga no muestra el nombre completo; bloque:\n%s", bloque)
	}
}

// CA3 / R8 / D16 (@CA3): el bloque delimitado T-15 al FINAL de estilos.css
// oculta .migas y .volver dentro de @media print acotado a .reporte
// (patrón de TestUS20_ReporteEstilosMediaPrintPresentes).
func TestT15_ImprimirReporteOcultaMigas(t *testing.T) {
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

	inicio := strings.Index(css, "/* === T-15")
	if inicio < 0 {
		t.Fatal("estilos.css no contiene el bloque delimitado T-15 (/* === T-15 … */)")
	}
	fin := strings.Index(css[inicio:], "(fin) === */")
	if fin < 0 {
		t.Fatal("el bloque delimitado T-15 no tiene su comentario de cierre (fin) === */")
	}
	bloque := css[inicio : inicio+fin]
	for _, fragmento := range []string{
		"@media print",
		"body:has(.reporte)",
		".migas",
		".volver",
		"display: none",
	} {
		if !strings.Contains(bloque, fragmento) {
			t.Errorf("el bloque T-15 de estilos.css no contiene %q; bloque:\n%s", fragmento, bloque)
		}
	}
}

// CA3 / R8 / D4 / @CA3: el PDF del reporte (generado fuera del HTML) no
// contiene el bloque de migas ni el enlace "← Volver".
func TestT15_PDFNoIncluyeMigas(t *testing.T) {
	entorno := nuevoEntornoT15()
	handler := entorno.handler(t)

	for _, ruta := range []string{"/proyectos/PRJ-001/reporte.pdf", "/proyectos/PRJ-001/sprints/1/reporte.pdf"} {
		t.Run(ruta, func(t *testing.T) {
			rec := solicitarUS22(handler, http.MethodGet, ruta, nil)
			exigirEstadoUS22(t, rec, http.StatusOK)
			pdf := rec.Body.Bytes()
			for _, prohibido := range [][]byte{[]byte("migas"), []byte("<nav"), []byte("← Volver")} {
				if bytes.Contains(pdf, prohibido) {
					t.Errorf("el PDF de %s contiene %q; el generador no debe emitir migas (D4)", ruta, prohibido)
				}
			}
		})
	}
}
