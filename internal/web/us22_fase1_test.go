package web_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
	"github.com/danaguzman22/software-metrics-estimation/internal/web"
)

// fechaUS22 arma una fecha del reloj inyectable de los tests de US-22.
func fechaUS22(anio int, mes time.Month, dia int) time.Time {
	return time.Date(anio, mes, dia, 0, 0, 0, 0, time.UTC)
}

func registroProyectoUS22(id, nombre string, inicio, fin time.Time) service.ProyectoRegistrado {
	return service.ProyectoRegistrado{
		ID: id,
		Proyecto: domain.Proyecto{
			Nombre:      nombre,
			FechaInicio: inicio,
			FechaFin:    fin,
			Estado:      domain.EstadoPlanificado,
		},
	}
}

func itemUS22(id, titulo, prioridad, estado string, storyPoints int) domain.ItemBacklog {
	return domain.ItemBacklog{
		ID:          id,
		Titulo:      titulo,
		Descripcion: "Descripción de " + titulo,
		Prioridad:   prioridad,
		Estado:      estado,
		StoryPoints: storyPoints,
	}
}

type creacionProyectoUS22 struct {
	nombre      string
	fechaInicio time.Time
	fechaFin    time.Time
}

// servicioProyectosFalso es el doble en memoria de web.ServicioProyectosWeb.
type servicioProyectosFalso struct {
	registros       []service.ProyectoRegistrado
	crearResultado  service.ProyectoRegistrado
	listarErr       error
	obtenerErr      error
	crearErr        error
	listarLlamadas  int
	obtenerLlamadas int
	crearLlamadas   int
	ultimaCreacion  creacionProyectoUS22
	ultimoID        string
}

func (f *servicioProyectosFalso) ListarProyectos() ([]service.ProyectoRegistrado, error) {
	f.listarLlamadas++
	if f.listarErr != nil {
		return nil, f.listarErr
	}
	return append([]service.ProyectoRegistrado(nil), f.registros...), nil
}

func (f *servicioProyectosFalso) ObtenerProyecto(id string) (service.ProyectoRegistrado, error) {
	f.obtenerLlamadas++
	f.ultimoID = id
	if f.obtenerErr != nil {
		return service.ProyectoRegistrado{}, f.obtenerErr
	}
	for _, registro := range f.registros {
		if registro.ID == id {
			return registro, nil
		}
	}
	return service.ProyectoRegistrado{}, domain.ErrProyectoNoEncontrado
}

func (f *servicioProyectosFalso) CrearProyecto(nombre string, fechaInicio, fechaFin time.Time) (service.ProyectoRegistrado, error) {
	f.crearLlamadas++
	f.ultimaCreacion = creacionProyectoUS22{nombre: nombre, fechaInicio: fechaInicio, fechaFin: fechaFin}
	if f.crearErr != nil {
		return service.ProyectoRegistrado{}, f.crearErr
	}
	resultado := f.crearResultado
	if resultado.ID == "" {
		resultado = registroProyectoUS22("PRJ-777", nombre, fechaInicio, fechaFin)
	}
	f.registros = append(f.registros, resultado)
	return resultado, nil
}

func (f *servicioProyectosFalso) ListarIntegrantes(proyectoID string) ([]domain.Integrante, error) {
	return nil, nil
}

func (f *servicioProyectosFalso) AgregarIntegrante(proyectoID, nombre, rol string) (domain.Integrante, error) {
	return domain.Integrante{}, nil
}

type llamadaCrearItemUS22 struct {
	proyectoID  string
	titulo      string
	descripcion string
	prioridad   string
	criterios   []string
}

type llamadaEditarItemUS22 struct {
	proyectoID  string
	itemID      string
	titulo      string
	descripcion string
	prioridad   string
	criterios   []string
	storyPoints int
}

type llamadaEstimarItemUS22 struct {
	proyectoID  string
	itemID      string
	storyPoints int
}

// servicioBacklogFalso es el doble en memoria de web.ServicioBacklogWeb.
type servicioBacklogFalso struct {
	items               map[string][]domain.ItemBacklog
	proyectosExistentes map[string]bool
	listarErr           error
	obtenerErr          error
	crearErr            error
	editarErr           error
	estimarErr          error
	listarLlamadas      int
	obtenerLlamadas     int
	crearLlamadas       int
	editarLlamadas      int
	estimarLlamadas     int
	ultimaCreacion      llamadaCrearItemUS22
	ultimaEdicion       llamadaEditarItemUS22
	ultimaEstimacion    llamadaEstimarItemUS22
	ultimoProyectoID    string
	ultimoItemID        string
}

func nuevoServicioBacklogFalso() *servicioBacklogFalso {
	return &servicioBacklogFalso{
		items:               map[string][]domain.ItemBacklog{},
		proyectosExistentes: map[string]bool{},
	}
}

func (f *servicioBacklogFalso) agregarItem(proyectoID string, item domain.ItemBacklog) {
	f.items[proyectoID] = append(f.items[proyectoID], item)
}

func (f *servicioBacklogFalso) ListarBacklog(proyectoID string) ([]domain.ItemBacklog, error) {
	f.listarLlamadas++
	f.ultimoProyectoID = proyectoID
	if f.listarErr != nil {
		return nil, f.listarErr
	}
	if !f.proyectosExistentes[proyectoID] {
		return nil, domain.ErrProyectoNoEncontrado
	}
	return append([]domain.ItemBacklog(nil), f.items[proyectoID]...), nil
}

func (f *servicioBacklogFalso) ObtenerItem(proyectoID, itemID string) (domain.ItemBacklog, error) {
	f.obtenerLlamadas++
	f.ultimoProyectoID = proyectoID
	f.ultimoItemID = itemID
	if f.obtenerErr != nil {
		return domain.ItemBacklog{}, f.obtenerErr
	}
	if !f.proyectosExistentes[proyectoID] {
		return domain.ItemBacklog{}, domain.ErrProyectoNoEncontrado
	}
	for _, item := range f.items[proyectoID] {
		if item.ID == itemID {
			return item, nil
		}
	}
	return domain.ItemBacklog{}, domain.ErrItemNoEncontrado
}

func (f *servicioBacklogFalso) CrearItem(proyectoID, titulo, descripcion, prioridad string, criterios []string) (domain.ItemBacklog, error) {
	f.crearLlamadas++
	f.ultimoProyectoID = proyectoID
	f.ultimaCreacion = llamadaCrearItemUS22{
		proyectoID:  proyectoID,
		titulo:      titulo,
		descripcion: descripcion,
		prioridad:   prioridad,
		criterios:   append([]string(nil), criterios...),
	}
	if f.crearErr != nil {
		return domain.ItemBacklog{}, f.crearErr
	}
	if !f.proyectosExistentes[proyectoID] {
		return domain.ItemBacklog{}, domain.ErrProyectoNoEncontrado
	}
	item := domain.ItemBacklog{
		ID:          fmt.Sprintf("HU-%02d", len(f.items[proyectoID])+1),
		Titulo:      titulo,
		Descripcion: descripcion,
		Prioridad:   prioridad,
		Estado:      "Pendiente",
		StoryPoints: 0,
		Criterios:   append([]string(nil), criterios...),
	}
	f.items[proyectoID] = append(f.items[proyectoID], item)
	return item, nil
}

func (f *servicioBacklogFalso) EditarItem(
	proyectoID, itemID, titulo, descripcion, prioridad string,
	criterios []string,
	storyPoints int,
) (domain.ItemBacklog, error) {
	f.editarLlamadas++
	f.ultimoProyectoID = proyectoID
	f.ultimoItemID = itemID
	f.ultimaEdicion = llamadaEditarItemUS22{
		proyectoID:  proyectoID,
		itemID:      itemID,
		titulo:      titulo,
		descripcion: descripcion,
		prioridad:   prioridad,
		criterios:   append([]string(nil), criterios...),
		storyPoints: storyPoints,
	}
	if f.editarErr != nil {
		return domain.ItemBacklog{}, f.editarErr
	}
	if !f.proyectosExistentes[proyectoID] {
		return domain.ItemBacklog{}, domain.ErrProyectoNoEncontrado
	}
	items := f.items[proyectoID]
	for i := range items {
		if items[i].ID == itemID {
			items[i].Titulo = titulo
			items[i].Descripcion = descripcion
			items[i].Prioridad = prioridad
			items[i].Criterios = append([]string(nil), criterios...)
			items[i].StoryPoints = storyPoints
			f.items[proyectoID] = items
			return items[i], nil
		}
	}
	return domain.ItemBacklog{}, domain.ErrItemNoEncontrado
}

func (f *servicioBacklogFalso) EstimarItem(proyectoID, itemID string, storyPoints int) (domain.ItemBacklog, error) {
	f.estimarLlamadas++
	f.ultimoProyectoID = proyectoID
	f.ultimoItemID = itemID
	f.ultimaEstimacion = llamadaEstimarItemUS22{proyectoID: proyectoID, itemID: itemID, storyPoints: storyPoints}
	if f.estimarErr != nil {
		return domain.ItemBacklog{}, f.estimarErr
	}
	if !f.proyectosExistentes[proyectoID] {
		return domain.ItemBacklog{}, domain.ErrProyectoNoEncontrado
	}
	items := f.items[proyectoID]
	for i := range items {
		if items[i].ID == itemID {
			items[i].StoryPoints = storyPoints
			f.items[proyectoID] = items
			return items[i], nil
		}
	}
	return domain.ItemBacklog{}, domain.ErrItemNoEncontrado
}

func nuevoEntornoUS22(registros ...service.ProyectoRegistrado) (*servicioProyectosFalso, *servicioBacklogFalso) {
	proyectos := &servicioProyectosFalso{}
	backlog := nuevoServicioBacklogFalso()
	for _, registro := range registros {
		proyectos.registros = append(proyectos.registros, registro)
		backlog.proyectosExistentes[registro.ID] = true
	}
	return proyectos, backlog
}

func nuevoHandlerUS22(t *testing.T, proyectos *servicioProyectosFalso, backlog *servicioBacklogFalso, sprintsOpcionales ...*servicioSprintsFalso) http.Handler {
	t.Helper()
	if proyectos == nil {
		proyectos = &servicioProyectosFalso{}
	}
	if backlog == nil {
		backlog = nuevoServicioBacklogFalso()
	}

	// La FASE 2 agregó ServicioSprintsWeb a las dependencias. El parámetro es
	// variádico para no modificar las llamadas existentes de la FASE 1.
	var sprints *servicioSprintsFalso
	if len(sprintsOpcionales) > 0 && sprintsOpcionales[0] != nil {
		sprints = sprintsOpcionales[0]
	} else {
		sprints = nuevoServicioSprintsFalso()
	}
	for _, registro := range proyectos.registros {
		sprints.marcarProyectoExistente(registro.ID)
	}

	return web.NuevoHandlerAplicacion(&pingFalso{}, web.Dependencias{
		Proyectos: proyectos,
		Backlog:   backlog,
		Sprints:   sprints,
		Dashboard: &servicioDashboardFalso{},
		Reportes:  &servicioReportesFalso{},
		Ahora: func() time.Time {
			return fechaUS22(2026, time.February, 1)
		},
	})
}

func solicitarUS22(handler http.Handler, metodo, ruta string, formulario url.Values) *httptest.ResponseRecorder {
	var cuerpo io.Reader
	if formulario != nil {
		cuerpo = strings.NewReader(formulario.Encode())
	}
	peticion := httptest.NewRequest(metodo, ruta, cuerpo)
	if formulario != nil {
		peticion.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	grabador := httptest.NewRecorder()
	handler.ServeHTTP(grabador, peticion)
	return grabador
}

func exigirEstadoUS22(t *testing.T, rec *httptest.ResponseRecorder, esperado int) {
	t.Helper()
	if rec.Code != esperado {
		t.Fatalf("código = %d, se esperaba %d; cuerpo:\n%s", rec.Code, esperado, rec.Body.String())
	}
}

func exigirCuerpoContieneUS22(t *testing.T, rec *httptest.ResponseRecorder, fragmentos ...string) {
	t.Helper()
	cuerpo := rec.Body.String()
	for _, fragmento := range fragmentos {
		if !strings.Contains(cuerpo, fragmento) {
			t.Errorf("el cuerpo no contiene %q; cuerpo:\n%s", fragmento, cuerpo)
		}
	}
}

func exigirCuerpoNoContieneUS22(t *testing.T, rec *httptest.ResponseRecorder, fragmentos ...string) {
	t.Helper()
	cuerpo := rec.Body.String()
	for _, fragmento := range fragmentos {
		if strings.Contains(cuerpo, fragmento) {
			t.Errorf("el cuerpo no debería contener %q; cuerpo:\n%s", fragmento, cuerpo)
		}
	}
}

func posicionUS22(t *testing.T, cuerpo, fragmento string) int {
	t.Helper()
	posicion := strings.Index(cuerpo, fragmento)
	if posicion < 0 {
		t.Fatalf("el cuerpo no contiene %q; cuerpo:\n%s", fragmento, cuerpo)
	}
	return posicion
}

func filaProyectoUS22(t *testing.T, cuerpo, id string) string {
	t.Helper()
	marca := "<td>" + id + "</td>"
	inicioID := strings.Index(cuerpo, marca)
	if inicioID < 0 {
		t.Fatalf("no se encontró la fila del proyecto %s; cuerpo:\n%s", id, cuerpo)
	}
	inicioFila := strings.LastIndex(cuerpo[:inicioID], "<tr>")
	finRelativo := strings.Index(cuerpo[inicioID:], "</tr>")
	if inicioFila < 0 || finRelativo < 0 {
		t.Fatalf("no se pudo delimitar la fila del proyecto %s", id)
	}
	finFila := inicioID + finRelativo + len("</tr>")
	return cuerpo[inicioFila:finFila]
}

func exigirRedireccionUS22(t *testing.T, rec *httptest.ResponseRecorder, destino string) {
	t.Helper()
	exigirEstadoUS22(t, rec, http.StatusSeeOther)
	if obtenido := rec.Header().Get("Location"); obtenido != destino {
		t.Fatalf("Location = %q, se esperaba %q", obtenido, destino)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("el cuerpo de la redirección debe estar vacío, se obtuvo %q", rec.Body.String())
	}
}

func exigirServiciosSinLlamadasUS22(t *testing.T, proyectos *servicioProyectosFalso, backlog *servicioBacklogFalso) {
	t.Helper()
	if proyectos.listarLlamadas+proyectos.obtenerLlamadas+proyectos.crearLlamadas != 0 {
		t.Errorf("se llamó al servicio de proyectos: %+v", proyectos)
	}
	if backlog.listarLlamadas+backlog.obtenerLlamadas+backlog.crearLlamadas+backlog.editarLlamadas+backlog.estimarLlamadas != 0 {
		t.Errorf("se llamó al servicio de backlog: %+v", backlog)
	}
}

// CA1 / R1: la navegación principal está presente en todas las páginas y la
// leyenda aparece en inicio y proyectos.
func TestUS22_Fase1_NavegacionEnTodasLasPaginas(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	backlog.agregarItem("PRJ-001", itemUS22("HU-01", "Iniciar sesión", "Media", "Pendiente", 3))
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	casos := []struct {
		ruta    string
		estado  int
		leyenda bool
	}{
		{ruta: "/", estado: http.StatusOK, leyenda: true},
		{ruta: "/proyectos", estado: http.StatusOK, leyenda: true},
		{ruta: "/proyectos/nuevo", estado: http.StatusOK, leyenda: false},
		{ruta: "/proyectos/PRJ-001", estado: http.StatusOK, leyenda: false},
		{ruta: "/proyectos/PRJ-001/backlog", estado: http.StatusOK, leyenda: false},
		{ruta: "/proyectos/PRJ-001/backlog/nuevo", estado: http.StatusOK, leyenda: false},
		{ruta: "/proyectos/PRJ-001/backlog/HU-01/editar", estado: http.StatusOK, leyenda: false},
		{ruta: "/proyectos/PRJ-001/sprints", estado: http.StatusOK, leyenda: false},
		{ruta: "/defectos", estado: http.StatusOK, leyenda: false},
		{ruta: "/proyectos/PRJ-999", estado: http.StatusNotFound, leyenda: false},
	}

	for _, caso := range casos {
		t.Run(caso.ruta, func(t *testing.T) {
			rec := solicitarUS22(handler, http.MethodGet, caso.ruta, nil)
			exigirEstadoUS22(t, rec, caso.estado)
			// D18 (US-19): la navegación gana el enlace Dashboard; los cuatro
			// enlaces preexistentes de US-22 (Proyectos/Backlog/Sprints/Defectos)
			// siguen presentes (R1 de US-22).
			exigirCuerpoContieneUS22(t, rec, "Proyectos", "Backlog", "Sprints", "Defectos", "Dashboard")

			cuerpo := rec.Body.String()
			tieneLeyenda := strings.Contains(cuerpo, "Elegí un proyecto para ver su backlog y sus sprints")
			if caso.leyenda && !tieneLeyenda {
				t.Errorf("la página %s debería mostrar la leyenda de contexto", caso.ruta)
			}
			if !caso.leyenda && tieneLeyenda {
				t.Errorf("la página %s no debería mostrar la leyenda de contexto", caso.ruta)
			}
		})
	}
}

// CA1 / R2: la página de inicio muestra el título y accesos a proyectos.
func TestUS22_Fase1_Inicio(t *testing.T) {
	handler := nuevoHandlerUS22(t, nil, nil)

	rec := solicitarUS22(handler, http.MethodGet, "/", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "Software Metrics & Estimation", "/proyectos", "/proyectos/nuevo")
}

// CA1 / R3: la sección de defectos informa que estará disponible próximamente
// sin ejecutar lógica de negocio.
func TestUS22_Fase1_DefectosProximamente(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22()
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	rec := solicitarUS22(handler, http.MethodGet, "/defectos", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "Disponible próximamente")
	if !strings.Contains(strings.ToLower(rec.Body.String()), "defectos") {
		t.Errorf("el cuerpo debería informar sobre la sección de defectos; cuerpo:\n%s", rec.Body.String())
	}
	exigirServiciosSinLlamadasUS22(t, proyectos, backlog)
}

// CA1 / R4: en la FASE 1 la sección de sprints informaba "Disponible
// próximamente". Ese comportamiento cambió en la FASE 2: ahora la página lista
// los Sprints del proyecto. El test conserva la verificación de que el mensaje
// de FASE 1 ya no aparece y de que un proyecto inexistente sigue devolviendo 404.
func TestUS22_Fase1_SprintsProximamente(t *testing.T) {
	t.Run("proyecto existente", func(t *testing.T) {
		proyectos, backlog := nuevoEntornoUS22(
			registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
		)
		handler := nuevoHandlerUS22(t, proyectos, backlog)

		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints", nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		exigirCuerpoNoContieneUS22(t, rec, "Disponible próximamente")
		if !strings.Contains(strings.ToLower(rec.Body.String()), "sprints") {
			t.Errorf("el cuerpo debería informar sobre la sección de sprints; cuerpo:\n%s", rec.Body.String())
		}
	})

	t.Run("proyecto inexistente", func(t *testing.T) {
		proyectos, backlog := nuevoEntornoUS22()
		handler := nuevoHandlerUS22(t, proyectos, backlog)

		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-999/sprints", nil)
		exigirEstadoUS22(t, rec, http.StatusNotFound)
		exigirCuerpoContieneUS22(t, rec, domain.ErrProyectoNoEncontrado.Error())
	})
}

// CA1 / R5: listar los proyectos con ID, nombre, fechas, estado calculado y en
// el orden devuelto por el servicio.
func TestUS22_Fase1_Proyectos_Listar(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
		registroProyectoUS22("PRJ-002", "Portal de Clientes", fechaUS22(2026, time.April, 15), fechaUS22(2026, time.August, 31)),
		registroProyectoUS22("PRJ-003", "Sistema legado", fechaUS22(2025, time.January, 1), fechaUS22(2025, time.December, 31)),
	)
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec,
		"PRJ-001", "Sistema de Métricas", "2026-03-01", "2026-06-30",
		"PRJ-002", "Portal de Clientes", "2026-04-15", "2026-08-31",
		"PRJ-003", "Sistema legado", "2025-01-01", "2025-12-31",
		"Planificado", "Finalizado",
	)

	cuerpo := rec.Body.String()
	filaPlanificado := filaProyectoUS22(t, cuerpo, "PRJ-001")
	if !strings.Contains(filaPlanificado, "Planificado") {
		t.Errorf("la fila de PRJ-001 debería mostrar el estado Planificado; fila: %s", filaPlanificado)
	}
	filaFinalizado := filaProyectoUS22(t, cuerpo, "PRJ-003")
	if !strings.Contains(filaFinalizado, "Finalizado") {
		t.Errorf("la fila de PRJ-003 debería mostrar el estado Finalizado; fila: %s", filaFinalizado)
	}
	if posicionUS22(t, cuerpo, "Sistema de Métricas") > posicionUS22(t, cuerpo, "Portal de Clientes") {
		t.Error("Sistema de Métricas debería aparecer antes que Portal de Clientes")
	}
	if posicionUS22(t, cuerpo, "Portal de Clientes") > posicionUS22(t, cuerpo, "Sistema legado") {
		t.Error("Portal de Clientes debería aparecer antes que Sistema legado")
	}
}

// CA1 / R6: una lista vacía muestra un mensaje sin error.
func TestUS22_Fase1_Proyectos_ListaVacia(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22()
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "No hay proyectos registrados")
}

// CA1 / R7: crear un proyecto válido responde 303 con la Location del detalle y
// delega los datos parseados en el servicio.
func TestUS22_Fase1_Proyectos_Crear(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22()
	inicio := fechaUS22(2026, time.March, 1)
	fin := fechaUS22(2026, time.June, 30)
	proyectos.crearResultado = registroProyectoUS22("PRJ-777", "Sistema de Métricas", inicio, fin)
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos", url.Values{
		"nombre":      {"Sistema de Métricas"},
		"fechaInicio": {"2026-03-01"},
		"fechaFin":    {"2026-06-30"},
	})
	exigirRedireccionUS22(t, rec, "/proyectos/PRJ-777")

	if proyectos.crearLlamadas != 1 {
		t.Fatalf("CrearProyecto se llamó %d veces, se esperaba 1", proyectos.crearLlamadas)
	}
	if proyectos.ultimaCreacion.nombre != "Sistema de Métricas" {
		t.Errorf("nombre = %q, se esperaba %q", proyectos.ultimaCreacion.nombre, "Sistema de Métricas")
	}
	if !proyectos.ultimaCreacion.fechaInicio.Equal(inicio) {
		t.Errorf("fechaInicio = %v, se esperaba %v", proyectos.ultimaCreacion.fechaInicio, inicio)
	}
	if !proyectos.ultimaCreacion.fechaFin.Equal(fin) {
		t.Errorf("fechaFin = %v, se esperaba %v", proyectos.ultimaCreacion.fechaFin, fin)
	}
}

// CA1 / R8: el detalle muestra nombre, fechas y estado calculado por fechas.
func TestUS22_Fase1_Proyectos_Detalle(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "Sistema de Métricas", "2026-03-01", "2026-06-30", "Planificado")
}

// CA1 / R5 / R8: el estado se calcula con domain.EstadoPorFechas y el reloj
// inyectable, no con el estado almacenado del proyecto.
func TestUS22_Fase1_Proyectos_EstadoCalculado(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-003", "Sistema legado", fechaUS22(2025, time.January, 1), fechaUS22(2025, time.December, 31)),
	)
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	t.Run("detalle", func(t *testing.T) {
		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-003", nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		exigirCuerpoContieneUS22(t, rec, "Sistema legado", "Finalizado")
		exigirCuerpoNoContieneUS22(t, rec, "Planificado")
	})

	t.Run("listado", func(t *testing.T) {
		rec := solicitarUS22(handler, http.MethodGet, "/proyectos", nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		exigirCuerpoContieneUS22(t, rec, "Sistema legado", "Finalizado")
		exigirCuerpoNoContieneUS22(t, rec, "Planificado")
	})
}

// CA1 / CA3 / R13: el backlog muestra los datos de cada ítem en el orden del
// servicio.
func TestUS22_Fase1_Backlog_Listar(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	items := []domain.ItemBacklog{
		itemUS22("HU-02", "Cerrar sesión", "Alta", "Pendiente", 5),
		itemUS22("HU-03", "Recuperar contraseña", "Media", "Pendiente", 3),
		itemUS22("HU-01", "Iniciar sesión", "Baja", "Pendiente", 0),
	}
	for _, item := range items {
		backlog.agregarItem("PRJ-001", item)
	}
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/backlog", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	for _, item := range items {
		exigirCuerpoContieneUS22(t, rec,
			item.ID, item.Titulo, item.Prioridad, item.Estado, fmt.Sprintf("%d", item.StoryPoints),
		)
	}

	cuerpo := rec.Body.String()
	if posicionUS22(t, cuerpo, "Cerrar sesión") > posicionUS22(t, cuerpo, "Recuperar contraseña") {
		t.Error("Cerrar sesión debería aparecer antes que Recuperar contraseña")
	}
	if posicionUS22(t, cuerpo, "Recuperar contraseña") > posicionUS22(t, cuerpo, "Iniciar sesión") {
		t.Error("Recuperar contraseña debería aparecer antes que Iniciar sesión")
	}
}

// CA1 / R14: un proyecto sin ítems muestra un mensaje sin error.
func TestUS22_Fase1_Backlog_ListaVacia(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/backlog", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "El backlog está vacío")
}

// CA1 / R15 / R16: el formulario de alta muestra los campos y el POST crea el
// ítem con los criterios recortados y sin líneas vacías.
func TestUS22_Fase1_Backlog_Crear(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	t.Run("formulario de alta", func(t *testing.T) {
		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/backlog/nuevo", nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		exigirCuerpoContieneUS22(t, rec,
			`name="titulo"`, `name="descripcion"`, `name="prioridad"`, `name="criterios"`,
			"Alta", "Media", "Baja",
		)
	})

	t.Run("alta válida", func(t *testing.T) {
		rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog", url.Values{
			"titulo":      {"Registrar horas"},
			"descripcion": {"Registrar las horas trabajadas"},
			"prioridad":   {"Media"},
			"criterios":   {"  La hora debe ser válida \n\n El registro queda guardado \n"},
		})
		exigirRedireccionUS22(t, rec, "/proyectos/PRJ-001/backlog")

		if backlog.crearLlamadas != 1 {
			t.Fatalf("CrearItem se llamó %d veces, se esperaba 1", backlog.crearLlamadas)
		}
		llamada := backlog.ultimaCreacion
		if llamada.proyectoID != "PRJ-001" {
			t.Errorf("proyectoID = %q, se esperaba %q", llamada.proyectoID, "PRJ-001")
		}
		if llamada.titulo != "Registrar horas" {
			t.Errorf("titulo = %q, se esperaba %q", llamada.titulo, "Registrar horas")
		}
		if llamada.descripcion != "Registrar las horas trabajadas" {
			t.Errorf("descripcion = %q, se esperaba %q", llamada.descripcion, "Registrar las horas trabajadas")
		}
		if llamada.prioridad != "Media" {
			t.Errorf("prioridad = %q, se esperaba %q", llamada.prioridad, "Media")
		}
		esperados := []string{"La hora debe ser válida", "El registro queda guardado"}
		if len(llamada.criterios) != len(esperados) {
			t.Fatalf("criterios = %#v, se esperaban %#v", llamada.criterios, esperados)
		}
		for i := range esperados {
			if llamada.criterios[i] != esperados[i] {
				t.Errorf("criterios[%d] = %q, se esperaba %q", i, llamada.criterios[i], esperados[i])
			}
		}
	})
}

// CA1 / R17 / R18: el formulario de edición precarga los datos y story points, y
// el POST actualiza el ítem y redirige al backlog.
func TestUS22_Fase1_Backlog_Editar(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	backlog.agregarItem("PRJ-001", itemUS22("HU-01", "Iniciar sesión", "Media", "Pendiente", 3))
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	t.Run("formulario precargado", func(t *testing.T) {
		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/backlog/HU-01/editar", nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		exigirCuerpoContieneUS22(t, rec, "Iniciar sesión", `name="storyPoints"`, `value="3"`)
	})

	t.Run("edición válida", func(t *testing.T) {
		antes := backlog.editarLlamadas
		rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/editar", url.Values{
			"titulo":      {"Iniciar sesión con correo"},
			"descripcion": {"Iniciar sesión con correo electrónico"},
			"prioridad":   {"Media"},
			"criterios":   {"El correo debe ser válido"},
			"storyPoints": {"3"},
		})
		exigirRedireccionUS22(t, rec, "/proyectos/PRJ-001/backlog")

		if backlog.editarLlamadas != antes+1 {
			t.Fatalf("EditarItem se llamó %d veces, se esperaba %d", backlog.editarLlamadas, antes+1)
		}
		llamada := backlog.ultimaEdicion
		if llamada.proyectoID != "PRJ-001" {
			t.Errorf("proyectoID = %q, se esperaba %q", llamada.proyectoID, "PRJ-001")
		}
		if llamada.itemID != "HU-01" {
			t.Errorf("itemID = %q, se esperaba %q", llamada.itemID, "HU-01")
		}
		if llamada.titulo != "Iniciar sesión con correo" {
			t.Errorf("titulo = %q, se esperaba %q", llamada.titulo, "Iniciar sesión con correo")
		}
		if llamada.descripcion != "Iniciar sesión con correo electrónico" {
			t.Errorf("descripcion = %q, se esperaba %q", llamada.descripcion, "Iniciar sesión con correo electrónico")
		}
		if llamada.prioridad != "Media" {
			t.Errorf("prioridad = %q, se esperaba %q", llamada.prioridad, "Media")
		}
		esperados := []string{"El correo debe ser válido"}
		if len(llamada.criterios) != len(esperados) {
			t.Fatalf("criterios = %#v, se esperaban %#v", llamada.criterios, esperados)
		}
		for i := range esperados {
			if llamada.criterios[i] != esperados[i] {
				t.Errorf("criterios[%d] = %q, se esperaba %q", i, llamada.criterios[i], esperados[i])
			}
		}
		if llamada.storyPoints != 3 {
			t.Errorf("storyPoints = %d, se esperaba 3", llamada.storyPoints)
		}
	})

	t.Run("story points vacío llega como cero", func(t *testing.T) {
		antes := backlog.editarLlamadas
		rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/editar", url.Values{
			"titulo":      {"Iniciar sesión con correo"},
			"descripcion": {"Iniciar sesión con correo electrónico"},
			"prioridad":   {"Media"},
			"criterios":   {"El correo debe ser válido"},
			"storyPoints": {""},
		})
		exigirRedireccionUS22(t, rec, "/proyectos/PRJ-001/backlog")

		if backlog.editarLlamadas != antes+1 {
			t.Fatalf("EditarItem se llamó %d veces, se esperaba %d", backlog.editarLlamadas, antes+1)
		}
		if backlog.ultimaEdicion.storyPoints != 0 {
			t.Errorf("storyPoints = %d, se esperaba 0", backlog.ultimaEdicion.storyPoints)
		}
	})
}

// CA1 / R19: estimar un ítem delega el entero parseado y redirige al backlog.
func TestUS22_Fase1_Backlog_Estimar(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	backlog.agregarItem("PRJ-001", itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 0))
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/estimar", url.Values{
		"storyPoints": {"5"},
	})
	exigirRedireccionUS22(t, rec, "/proyectos/PRJ-001/backlog")

	if backlog.estimarLlamadas != 1 {
		t.Fatalf("EstimarItem se llamó %d veces, se esperaba 1", backlog.estimarLlamadas)
	}
	if backlog.ultimaEstimacion.itemID != "HU-01" || backlog.ultimaEstimacion.storyPoints != 5 {
		t.Errorf("EstimarItem recibió %+v, se esperaba itemID HU-01 y storyPoints 5", backlog.ultimaEstimacion)
	}
}

// CA2 / R9 / R10: los errores de dominio del alta se muestran con HTTP 422 y el
// formulario conserva los datos ingresados.
func TestUS22_Fase1_Proyectos_ErroresDominio(t *testing.T) {
	casos := []struct {
		descripcion   string
		nombre        string
		fechaInicio   string
		fechaFin      string
		errorDominio  error
		valorRetenido string
	}{
		{
			descripcion:   "nombre vacío",
			nombre:        "",
			fechaInicio:   "2026-03-01",
			fechaFin:      "2026-06-30",
			errorDominio:  domain.ErrNombreVacio,
			valorRetenido: "2026-03-01",
		},
		{
			descripcion:   "fecha de fin anterior a la de inicio",
			nombre:        "Sistema de Métricas",
			fechaInicio:   "2026-06-30",
			fechaFin:      "2026-03-01",
			errorDominio:  domain.ErrFechaFinAnteriorAInicio,
			valorRetenido: "Sistema de Métricas",
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoUS22()
			proyectos.crearErr = caso.errorDominio
			handler := nuevoHandlerUS22(t, proyectos, backlog)

			rec := solicitarUS22(handler, http.MethodPost, "/proyectos", url.Values{
				"nombre":      {caso.nombre},
				"fechaInicio": {caso.fechaInicio},
				"fechaFin":    {caso.fechaFin},
			})
			exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
			exigirCuerpoContieneUS22(t, rec, caso.errorDominio.Error(), caso.valorRetenido)
			if len(proyectos.registros) != 0 {
				t.Errorf("no debería crearse ningún proyecto, hay %d", len(proyectos.registros))
			}
		})
	}
}

// CA2 / R9: el formulario de alta de proyecto muestra sus tres campos.
func TestUS22_Fase1_Proyectos_FormularioVacio(t *testing.T) {
	handler := nuevoHandlerUS22(t, nil, nil)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/nuevo", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, `name="nombre"`, `name="fechaInicio"`, `name="fechaFin"`)
}

// CA2 / R12: una fecha con formato inválido o vacía produce HTTP 422 con un
// mensaje en español y no llama al servicio.
func TestUS22_Fase1_Proyectos_FechaInvalida(t *testing.T) {
	casos := []struct {
		descripcion string
		fechaInicio string
		fechaFin    string
		mensaje     string
	}{
		{
			descripcion: "fecha inexistente",
			fechaInicio: "2026-13-01",
			fechaFin:    "2026-06-30",
			mensaje:     "La fecha debe tener el formato AAAA-MM-DD",
		},
		{
			descripcion: "texto en lugar de fecha",
			fechaInicio: "mañana",
			fechaFin:    "2026-06-30",
			mensaje:     "La fecha debe tener el formato AAAA-MM-DD",
		},
		{
			descripcion: "fecha de inicio vacía",
			fechaInicio: "",
			fechaFin:    "2026-06-30",
			mensaje:     "La fecha de inicio es obligatoria",
		},
		{
			descripcion: "fecha de fin vacía",
			fechaInicio: "2026-03-01",
			fechaFin:    "",
			mensaje:     "La fecha de fin es obligatoria",
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoUS22()
			handler := nuevoHandlerUS22(t, proyectos, backlog)

			rec := solicitarUS22(handler, http.MethodPost, "/proyectos", url.Values{
				"nombre":      {"Sistema de Métricas"},
				"fechaInicio": {caso.fechaInicio},
				"fechaFin":    {caso.fechaFin},
			})
			exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
			exigirCuerpoContieneUS22(t, rec, caso.mensaje, "Sistema de Métricas")
			if proyectos.crearLlamadas != 0 {
				t.Errorf("CrearProyecto se llamó %d veces y no debía ejecutarse", proyectos.crearLlamadas)
			}
		})
	}
}

// CA2 / CA3 / R16 / R20: los errores de dominio del alta de ítem se muestran con
// HTTP 422 y conservan lo ingresado.
func TestUS22_Fase1_Backlog_ErroresDominio(t *testing.T) {
	casos := []struct {
		descripcion   string
		titulo        string
		prioridad     string
		errorDominio  error
		valorRetenido string
	}{
		{
			descripcion:   "título vacío",
			titulo:        "",
			prioridad:     "Alta",
			errorDominio:  domain.ErrTituloItemVacio,
			valorRetenido: "Alta",
		},
		{
			descripcion:   "prioridad inválida",
			titulo:        "Registrar horas",
			prioridad:     "Urgente",
			errorDominio:  domain.ErrPrioridadInvalida,
			valorRetenido: "Registrar horas",
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoUS22(
				registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
			)
			backlog.crearErr = caso.errorDominio
			handler := nuevoHandlerUS22(t, proyectos, backlog)

			rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog", url.Values{
				"titulo":      {caso.titulo},
				"descripcion": {"Descripción"},
				"prioridad":   {caso.prioridad},
				"criterios":   {"Un criterio"},
			})
			exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
			exigirCuerpoContieneUS22(t, rec, caso.errorDominio.Error(), caso.valorRetenido)
			if len(backlog.items["PRJ-001"]) != 0 {
				t.Errorf("el backlog debería seguir vacío, tiene %d ítems", len(backlog.items["PRJ-001"]))
			}
		})
	}
}

// CA2 / R20: un error de story points del dominio se muestra con HTTP 422 y el
// ítem conserva su estimación.
func TestUS22_Fase1_Backlog_StoryPointsInvalidos(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	backlog.agregarItem("PRJ-001", itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 0))
	backlog.estimarErr = domain.ErrStoryPointsInvalidos
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/estimar", url.Values{
		"storyPoints": {"4"},
	})
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, domain.ErrStoryPointsInvalidos.Error())
	if backlog.items["PRJ-001"][0].StoryPoints != 0 {
		t.Errorf("StoryPoints = %d, se esperaba conservar 0", backlog.items["PRJ-001"][0].StoryPoints)
	}
}

// CA2 / R20: un ítem completado devuelve ErrItemCompletadoNoEditable con HTTP
// 422.
func TestUS22_Fase1_Backlog_ItemCompletado(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	backlog.agregarItem("PRJ-001", itemUS22("HU-01", "Iniciar sesión", "Alta", "Completada", 3))
	backlog.editarErr = domain.ErrItemCompletadoNoEditable
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/editar", url.Values{
		"titulo":      {"Iniciar sesión con correo"},
		"descripcion": {"Iniciar sesión con correo electrónico"},
		"prioridad":   {"Media"},
		"criterios":   {"El correo debe ser válido"},
		"storyPoints": {"3"},
	})
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, domain.ErrItemCompletadoNoEditable.Error())
}

// CA2 / R21: editar o estimar un ítem inexistente devuelve HTTP 404 con el
// mensaje del dominio.
func TestUS22_Fase1_Backlog_ItemInexistente404(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	t.Run("editar", func(t *testing.T) {
		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/backlog/HU-99/editar", nil)
		exigirEstadoUS22(t, rec, http.StatusNotFound)
		exigirCuerpoContieneUS22(t, rec, domain.ErrItemNoEncontrado.Error())
	})

	t.Run("estimar", func(t *testing.T) {
		rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-99/estimar", url.Values{
			"storyPoints": {"5"},
		})
		exigirEstadoUS22(t, rec, http.StatusNotFound)
		exigirCuerpoContieneUS22(t, rec, domain.ErrItemNoEncontrado.Error())
	})
}

// CA2 / R11: un proyecto inexistente devuelve HTTP 404 con el mensaje del
// dominio en el detalle, el backlog y los sprints.
func TestUS22_Fase1_Proyectos_Inexistente404(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22()
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	for _, ruta := range []string{
		"/proyectos/PRJ-999",
		"/proyectos/PRJ-999/backlog",
		"/proyectos/PRJ-999/sprints",
	} {
		t.Run(ruta, func(t *testing.T) {
			rec := solicitarUS22(handler, http.MethodGet, ruta, nil)
			exigirEstadoUS22(t, rec, http.StatusNotFound)
			exigirCuerpoContieneUS22(t, rec, domain.ErrProyectoNoEncontrado.Error())
		})
	}
}

// CA2 / R22: un story points no numérico produce HTTP 422 y no llama al
// servicio.
func TestUS22_Fase1_ParseoNumeroInvalido(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	backlog.agregarItem("PRJ-001", itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 0))
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/estimar", url.Values{
		"storyPoints": {"muchos"},
	})
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, "Los story points deben ser un número entero")
	if backlog.estimarLlamadas != 0 {
		t.Errorf("EstimarItem se llamó %d veces y no debía ejecutarse", backlog.estimarLlamadas)
	}
}

// CA2 / R23: un error de infraestructura se muestra con HTTP 500 y un mensaje
// genérico, sin exponer el detalle interno.
func TestUS22_Fase1_ErrorInfraestructura(t *testing.T) {
	errInfra := errors.New("detalle interno: conexión caída")

	t.Run("listar proyectos", func(t *testing.T) {
		proyectos, backlog := nuevoEntornoUS22()
		proyectos.listarErr = errInfra
		handler := nuevoHandlerUS22(t, proyectos, backlog)

		rec := solicitarUS22(handler, http.MethodGet, "/proyectos", nil)
		exigirEstadoUS22(t, rec, http.StatusInternalServerError)
		exigirCuerpoContieneUS22(t, rec, "Ocurrió un error inesperado. Intentá nuevamente.")
		exigirCuerpoNoContieneUS22(t, rec, "detalle interno")
	})

	t.Run("crear proyecto", func(t *testing.T) {
		proyectos, backlog := nuevoEntornoUS22()
		proyectos.crearErr = errInfra
		handler := nuevoHandlerUS22(t, proyectos, backlog)

		rec := solicitarUS22(handler, http.MethodPost, "/proyectos", url.Values{
			"nombre":      {"Sistema de Métricas"},
			"fechaInicio": {"2026-03-01"},
			"fechaFin":    {"2026-06-30"},
		})
		exigirEstadoUS22(t, rec, http.StatusInternalServerError)
		exigirCuerpoContieneUS22(t, rec, "Ocurrió un error inesperado. Intentá nuevamente.")
		exigirCuerpoNoContieneUS22(t, rec, "detalle interno")
	})
}

// CA2 / R24: todo POST exitoso responde 303 con Location a una ruta GET y sin
// contenido.
func TestUS22_Fase1_PostRedirige303(t *testing.T) {
	t.Run("crear proyecto", func(t *testing.T) {
		proyectos, backlog := nuevoEntornoUS22()
		handler := nuevoHandlerUS22(t, proyectos, backlog)

		rec := solicitarUS22(handler, http.MethodPost, "/proyectos", url.Values{
			"nombre":      {"Sistema de Métricas"},
			"fechaInicio": {"2026-03-01"},
			"fechaFin":    {"2026-06-30"},
		})
		exigirRedireccionUS22(t, rec, "/proyectos/PRJ-777")
	})

	t.Run("crear ítem", func(t *testing.T) {
		proyectos, backlog := nuevoEntornoUS22(
			registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
		)
		handler := nuevoHandlerUS22(t, proyectos, backlog)

		rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog", url.Values{
			"titulo":    {"Registrar horas"},
			"prioridad": {"Alta"},
		})
		exigirRedireccionUS22(t, rec, "/proyectos/PRJ-001/backlog")
	})

	t.Run("editar ítem", func(t *testing.T) {
		proyectos, backlog := nuevoEntornoUS22(
			registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
		)
		backlog.agregarItem("PRJ-001", itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 3))
		handler := nuevoHandlerUS22(t, proyectos, backlog)

		rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/editar", url.Values{
			"titulo":      {"Iniciar sesión con correo"},
			"prioridad":   {"Alta"},
			"storyPoints": {"3"},
		})
		exigirRedireccionUS22(t, rec, "/proyectos/PRJ-001/backlog")
	})

	t.Run("estimar ítem", func(t *testing.T) {
		proyectos, backlog := nuevoEntornoUS22(
			registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
		)
		backlog.agregarItem("PRJ-001", itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 0))
		handler := nuevoHandlerUS22(t, proyectos, backlog)

		rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog/HU-01/estimar", url.Values{
			"storyPoints": {"5"},
		})
		exigirRedireccionUS22(t, rec, "/proyectos/PRJ-001/backlog")
	})
}

// CA2 / R25: un método no permitido sobre una ruta existente devuelve HTTP 405
// sin ejecutar el handler del caso de uso.
func TestUS22_Fase1_MetodoNoPermitido(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	backlog.agregarItem("PRJ-001", itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 0))
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	t.Run("GET sobre estimar", func(t *testing.T) {
		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/backlog/HU-01/estimar", nil)
		exigirEstadoUS22(t, rec, http.StatusMethodNotAllowed)
		if backlog.estimarLlamadas != 0 {
			t.Errorf("EstimarItem se llamó %d veces y no debía ejecutarse", backlog.estimarLlamadas)
		}
	})

	t.Run("DELETE sobre proyectos", func(t *testing.T) {
		rec := solicitarUS22(handler, http.MethodDelete, "/proyectos", nil)
		exigirEstadoUS22(t, rec, http.StatusMethodNotAllowed)
	})
}

// CA3 / R13 / R28: la web respeta el orden del backlog entregado por el
// servicio y no lo reordena.
func TestUS22_Fase1_BacklogRespetaOrdenDelServicio(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	backlog.agregarItem("PRJ-001", itemUS22("HU-03", "Historia C", "Baja", "Pendiente", 1))
	backlog.agregarItem("PRJ-001", itemUS22("HU-01", "Historia A", "Alta", "Pendiente", 2))
	backlog.agregarItem("PRJ-001", itemUS22("HU-02", "Historia B", "Media", "Pendiente", 3))
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/backlog", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)

	cuerpo := rec.Body.String()
	if posicionUS22(t, cuerpo, "Historia C") > posicionUS22(t, cuerpo, "Historia A") {
		t.Error("Historia C debería aparecer antes que Historia A: la web no debe reordenar")
	}
	if posicionUS22(t, cuerpo, "Historia A") > posicionUS22(t, cuerpo, "Historia B") {
		t.Error("Historia A debería aparecer antes que Historia B: la web no debe reordenar")
	}
}

// CA3 / R20 / R26: los handlers delegan en el servicio y muestran el error de
// dominio devuelto, sin reemplazarlo por uno de la web.
func TestUS22_Fase1_HandlersDeleganEnServicios(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	backlog.crearErr = domain.ErrPrioridadInvalida
	handler := nuevoHandlerUS22(t, proyectos, backlog)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/backlog", url.Values{
		"titulo":    {"Registrar horas"},
		"prioridad": {"Urgente"},
	})
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, domain.ErrPrioridadInvalida.Error())
	if backlog.crearLlamadas != 1 {
		t.Errorf("CrearItem se llamó %d veces, se esperaba 1", backlog.crearLlamadas)
	}
}

// CA3 / R27: la capa web no importa internal/storage ni pgx.
func TestUS22_Fase1_WebNoImportaStorageNiPgx(t *testing.T) {
	_, archivo, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("no se pudo obtener la ruta del test")
	}
	raizModulo := filepath.Clean(filepath.Join(filepath.Dir(archivo), "..", ".."))

	comando := exec.Command("go", "list", "-deps", "./internal/web")
	comando.Dir = raizModulo
	salida, err := comando.CombinedOutput()
	if err != nil {
		t.Fatalf("go list -deps ./internal/web falló: %v\n%s", err, salida)
	}

	prohibidos := []string{
		"github.com/danaguzman22/software-metrics-estimation/internal/storage",
		"github.com/danaguzman22/software-metrics-estimation/internal/storage/postgres",
		"github.com/danaguzman22/software-metrics-estimation/internal/storage/jsonrepo",
		"github.com/jackc/pgx",
	}
	for _, linea := range strings.Split(string(salida), "\n") {
		dependencia := strings.TrimSpace(linea)
		for _, prohibido := range prohibidos {
			if dependencia == prohibido || strings.HasPrefix(dependencia, prohibido+"/") {
				t.Errorf("internal/web no debe depender de %q", dependencia)
			}
		}
	}
}
