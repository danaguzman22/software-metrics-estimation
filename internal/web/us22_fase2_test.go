package web_test

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

type llamadaCrearSprintUS22 struct {
	proyectoID  string
	numero      int
	objetivo    string
	fechaInicio time.Time
	fechaFin    time.Time
}

// servicioSprintsFalso es el doble en memoria de web.ServicioSprintsWeb.
type servicioSprintsFalso struct {
	sprints               map[string][]domain.Sprint
	spPlanificados        map[string]map[int]int
	proyectosExistentes   map[string]bool
	listarErr             error
	crearErr              error
	asignarErr            error
	storyPointsErr        error
	listarLlamadas        int
	crearLlamadas         int
	asignarLlamadas       int
	storyPointsLlamadas   int
	ultimoProyectoID      string
	ultimoSprintNumero    int
	ultimoItemID          string
	ultimaCreacion        llamadaCrearSprintUS22
	itemsPorSprint        map[string]map[int][]domain.ItemBacklog
	completadosPorSprint  map[string]map[int]int
	velocidadPorSprint    map[string]map[int]int
	promedioPorID         map[string]float64
	porcentajePorID       map[string]float64
	itemsDelSprintErr     error
	completarErr          error
	cerrarErr             error
	completadosErr        error
	velocidadErr          error
	promedioErr           error
	porcentajeErr         error
	listarItemsLlamadas   int
	completarLlamadas     int
	cerrarLlamadas        int
	completadosLlamadas   int
	velocidadLlamadas     int
	promedioLlamadas      int
	porcentajeLlamadas    int
}

func nuevoServicioSprintsFalso() *servicioSprintsFalso {
	return &servicioSprintsFalso{
		sprints:              map[string][]domain.Sprint{},
		spPlanificados:       map[string]map[int]int{},
		proyectosExistentes:  map[string]bool{},
		itemsPorSprint:       map[string]map[int][]domain.ItemBacklog{},
		completadosPorSprint: map[string]map[int]int{},
		velocidadPorSprint:   map[string]map[int]int{},
		promedioPorID:        map[string]float64{},
		porcentajePorID:      map[string]float64{},
	}
}

func (f *servicioSprintsFalso) marcarProyectoExistente(proyectoID string) {
	f.proyectosExistentes[proyectoID] = true
}

func (f *servicioSprintsFalso) agregarSprint(proyectoID string, sprint domain.Sprint, planificados int) {
	f.sprints[proyectoID] = append(f.sprints[proyectoID], sprint)
	if f.spPlanificados[proyectoID] == nil {
		f.spPlanificados[proyectoID] = map[int]int{}
	}
	f.spPlanificados[proyectoID][sprint.Numero] = planificados
}

func (f *servicioSprintsFalso) ListarSprints(proyectoID string) ([]domain.Sprint, error) {
	f.listarLlamadas++
	f.ultimoProyectoID = proyectoID
	if f.listarErr != nil {
		return nil, f.listarErr
	}
	if !f.proyectosExistentes[proyectoID] {
		return nil, domain.ErrProyectoNoEncontrado
	}
	return append([]domain.Sprint(nil), f.sprints[proyectoID]...), nil
}

func (f *servicioSprintsFalso) StoryPointsPlanificados(proyectoID string, numero int) (int, error) {
	f.storyPointsLlamadas++
	f.ultimoProyectoID = proyectoID
	f.ultimoSprintNumero = numero
	if f.storyPointsErr != nil {
		return 0, f.storyPointsErr
	}
	if !f.proyectosExistentes[proyectoID] {
		return 0, domain.ErrProyectoNoEncontrado
	}
	return f.spPlanificados[proyectoID][numero], nil
}

func (f *servicioSprintsFalso) CrearSprint(
	proyectoID string,
	numero int,
	objetivo string,
	fechaInicio, fechaFin time.Time,
) (domain.Sprint, error) {
	f.crearLlamadas++
	f.ultimoProyectoID = proyectoID
	f.ultimaCreacion = llamadaCrearSprintUS22{
		proyectoID:  proyectoID,
		numero:      numero,
		objetivo:    objetivo,
		fechaInicio: fechaInicio,
		fechaFin:    fechaFin,
	}
	if f.crearErr != nil {
		return domain.Sprint{}, f.crearErr
	}

	sprint := domain.Sprint{
		Numero:      numero,
		Objetivo:    objetivo,
		FechaInicio: fechaInicio,
		FechaFin:    fechaFin,
		Estado:      domain.EstadoSprintActivo,
	}
	f.agregarSprint(proyectoID, sprint, 0)
	return sprint, nil
}

func (f *servicioSprintsFalso) AsignarItem(proyectoID string, numero int, itemID string) error {
	f.asignarLlamadas++
	f.ultimoProyectoID = proyectoID
	f.ultimoSprintNumero = numero
	f.ultimoItemID = itemID
	return f.asignarErr
}

// sprintUS22 arma un Sprint de referencia para sembrar el doble.
func sprintUS22(numero int, objetivo, estado string, inicio, fin time.Time) domain.Sprint {
	return domain.Sprint{
		Numero:      numero,
		Objetivo:    objetivo,
		FechaInicio: inicio,
		FechaFin:    fin,
		Estado:      estado,
	}
}

// CA1 / R29: el listado muestra número, objetivo, fechas, estado y Story Points
// planificados de cada Sprint, en el orden devuelto por el servicio.
func TestUS22_Fase2_Sprints_Listar(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	backlog.agregarItem("PRJ-001", itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 3))

	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 8)
	sprints.agregarSprint("PRJ-001", sprintUS22(
		2, "Segunda entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 15), fechaUS22(2026, time.March, 28),
	), 0)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec,
		"Primera entrega", "2026-03-01", "2026-03-14",
		"Segunda entrega", "2026-03-15", "2026-03-28",
		"Activo", "8", "0",
	)
	exigirCuerpoNoContieneUS22(t, rec, "Disponible próximamente")

	cuerpo := rec.Body.String()
	if posicionUS22(t, cuerpo, "Primera entrega") > posicionUS22(t, cuerpo, "Segunda entrega") {
		t.Error("Primera entrega debería aparecer antes que Segunda entrega: la web no debe reordenar")
	}

	if sprints.listarLlamadas != 1 {
		t.Errorf("ListarSprints se llamó %d veces, se esperaba 1", sprints.listarLlamadas)
	}
	if sprints.ultimoProyectoID != "PRJ-001" {
		t.Errorf("ListarSprints recibió el proyecto %q, se esperaba PRJ-001", sprints.ultimoProyectoID)
	}
	if sprints.storyPointsLlamadas != 2 {
		t.Errorf("StoryPointsPlanificados se llamó %d veces, se esperaban 2", sprints.storyPointsLlamadas)
	}
}

// CA3 / R28 / R29 / D19: los Story Points mostrados son el total del servicio,
// no una suma recalculada por la web.
func TestUS22_Fase2_Sprints_StoryPointsPlanificados(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	backlog.agregarItem("PRJ-001", itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 1))
	backlog.agregarItem("PRJ-001", itemUS22("HU-02", "Cerrar sesión", "Media", "Pendiente", 2))

	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		3, "Iteración de marzo", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 34)
	sprints.agregarSprint("PRJ-001", sprintUS22(
		4, "Iteración de abril", domain.EstadoSprintActivo,
		fechaUS22(2026, time.April, 1), fechaUS22(2026, time.April, 14),
	), 55)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "Iteración de marzo", "Iteración de abril", "34", "55")

	if sprints.storyPointsLlamadas != 2 {
		t.Errorf("StoryPointsPlanificados se llamó %d veces, se esperaban 2", sprints.storyPointsLlamadas)
	}
}

// CA1 / R29 / D19: un Sprint cerrado muestra el total guardado al cerrarlo.
func TestUS22_Fase2_Sprints_ListarSprintCerrado(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintCerrado,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 8)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "Primera entrega", "Cerrado", "8")
}

// CA1 / R30 / R24: el formulario de alta muestra sus campos y un alta válida
// delega en el servicio y responde 303 hacia el listado.
func TestUS22_Fase2_Sprints_Crear(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	sprints := nuevoServicioSprintsFalso()
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	t.Run("formulario de alta", func(t *testing.T) {
		rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints/nuevo", nil)
		exigirEstadoUS22(t, rec, http.StatusOK)
		exigirCuerpoContieneUS22(t, rec,
			`name="numero"`, `name="objetivo"`, `name="fechaInicio"`, `name="fechaFin"`,
		)
	})

	t.Run("alta válida", func(t *testing.T) {
		rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/sprints", url.Values{
			"numero":      {"1"},
			"objetivo":    {"Primera entrega"},
			"fechaInicio": {"2026-03-01"},
			"fechaFin":    {"2026-03-14"},
		})
		exigirRedireccionUS22(t, rec, "/proyectos/PRJ-001/sprints")

		if sprints.crearLlamadas != 1 {
			t.Fatalf("CrearSprint se llamó %d veces, se esperaba 1", sprints.crearLlamadas)
		}
		llamada := sprints.ultimaCreacion
		if llamada.proyectoID != "PRJ-001" {
			t.Errorf("proyectoID = %q, se esperaba %q", llamada.proyectoID, "PRJ-001")
		}
		if llamada.numero != 1 {
			t.Errorf("numero = %d, se esperaba 1", llamada.numero)
		}
		if llamada.objetivo != "Primera entrega" {
			t.Errorf("objetivo = %q, se esperaba %q", llamada.objetivo, "Primera entrega")
		}
		if !llamada.fechaInicio.Equal(fechaUS22(2026, time.March, 1)) {
			t.Errorf("fechaInicio = %v, se esperaba 2026-03-01", llamada.fechaInicio)
		}
		if !llamada.fechaFin.Equal(fechaUS22(2026, time.March, 14)) {
			t.Errorf("fechaFin = %v, se esperaba 2026-03-14", llamada.fechaFin)
		}
	})
}

// CA2 / R30 / E12: los errores de dominio del alta se muestran con HTTP 422 y el
// formulario conserva los datos ingresados.
func TestUS22_Fase2_Sprints_CrearErrores(t *testing.T) {
	casos := []struct {
		descripcion   string
		numero        string
		objetivo      string
		fechaInicio   string
		fechaFin      string
		errorDominio  error
		mensaje       string
		llamaServicio bool
		valorRetenido string
	}{
		{
			descripcion:   "número inválido",
			numero:        "0",
			objetivo:      "Primera entrega",
			fechaInicio:   "2026-03-01",
			fechaFin:      "2026-03-14",
			errorDominio:  domain.ErrNumeroSprintInvalido,
			mensaje:       domain.ErrNumeroSprintInvalido.Error(),
			llamaServicio: true,
			valorRetenido: "Primera entrega",
		},
		{
			descripcion:   "objetivo vacío",
			numero:        "1",
			objetivo:      "",
			fechaInicio:   "2026-03-01",
			fechaFin:      "2026-03-14",
			errorDominio:  domain.ErrObjetivoSprintVacio,
			mensaje:       domain.ErrObjetivoSprintVacio.Error(),
			llamaServicio: true,
			valorRetenido: "2026-03-01",
		},
		{
			descripcion:   "fechas inválidas",
			numero:        "1",
			objetivo:      "Primera entrega",
			fechaInicio:   "2026-03-14",
			fechaFin:      "2026-03-01",
			errorDominio:  domain.ErrFechasSprintInvalidas,
			mensaje:       domain.ErrFechasSprintInvalidas.Error(),
			llamaServicio: true,
			valorRetenido: "Primera entrega",
		},
		{
			descripcion:   "sprint activo existente",
			numero:        "2",
			objetivo:      "Segunda entrega",
			fechaInicio:   "2026-03-15",
			fechaFin:      "2026-03-28",
			errorDominio:  domain.ErrSprintActivoExistente,
			mensaje:       domain.ErrSprintActivoExistente.Error(),
			llamaServicio: true,
			valorRetenido: "Segunda entrega",
		},
		{
			descripcion:   "número duplicado",
			numero:        "1",
			objetivo:      "Primera entrega",
			fechaInicio:   "2026-03-01",
			fechaFin:      "2026-03-14",
			errorDominio:  domain.ErrSprintDuplicado,
			mensaje:       domain.ErrSprintDuplicado.Error(),
			llamaServicio: true,
			valorRetenido: "Primera entrega",
		},
		{
			descripcion:   "número no numérico",
			numero:        "abc",
			objetivo:      "Primera entrega",
			fechaInicio:   "2026-03-01",
			fechaFin:      "2026-03-14",
			mensaje:       "El número del sprint debe ser un número entero",
			llamaServicio: false,
			valorRetenido: "abc",
		},
		{
			descripcion:   "fecha con formato inválido",
			numero:        "1",
			objetivo:      "Primera entrega",
			fechaInicio:   "01/03/2026",
			fechaFin:      "2026-03-14",
			mensaje:       "La fecha debe tener el formato AAAA-MM-DD",
			llamaServicio: false,
			valorRetenido: "Primera entrega",
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoUS22(
				registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
			)
			sprints := nuevoServicioSprintsFalso()
			if caso.errorDominio != nil {
				sprints.crearErr = caso.errorDominio
			}
			handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

			rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/sprints", url.Values{
				"numero":      {caso.numero},
				"objetivo":    {caso.objetivo},
				"fechaInicio": {caso.fechaInicio},
				"fechaFin":    {caso.fechaFin},
			})
			exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
			exigirCuerpoContieneUS22(t, rec, caso.mensaje, caso.valorRetenido)

			if caso.llamaServicio {
				if sprints.crearLlamadas != 1 {
					t.Errorf("CrearSprint se llamó %d veces, se esperaba 1", sprints.crearLlamadas)
				}
			} else if sprints.crearLlamadas != 0 {
				t.Errorf("CrearSprint se llamó %d veces y no debía ejecutarse", sprints.crearLlamadas)
			}
		})
	}
}

// CA1 / R31 / R24: asignar una historia válida delega en el servicio y responde
// 303 hacia el listado.
func TestUS22_Fase2_Sprints_Asignar(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	backlog.agregarItem("PRJ-001", itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 3))
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 0)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/sprints/1/asignar", url.Values{
		"itemID": {"HU-01"},
	})
	exigirRedireccionUS22(t, rec, "/proyectos/PRJ-001/sprints")

	if sprints.asignarLlamadas != 1 {
		t.Fatalf("AsignarItem se llamó %d veces, se esperaba 1", sprints.asignarLlamadas)
	}
	if sprints.ultimoProyectoID != "PRJ-001" {
		t.Errorf("proyectoID = %q, se esperaba PRJ-001", sprints.ultimoProyectoID)
	}
	if sprints.ultimoSprintNumero != 1 {
		t.Errorf("sprintNumero = %d, se esperaba 1", sprints.ultimoSprintNumero)
	}
	if sprints.ultimoItemID != "HU-01" {
		t.Errorf("itemID = %q, se esperaba HU-01", sprints.ultimoItemID)
	}
}

// CA2 / R31 / E13: las historias no estimadas, no pendientes o ya asignadas se
// muestran con HTTP 422 y el mensaje del dominio.
func TestUS22_Fase2_Sprints_AsignarErrores(t *testing.T) {
	casos := []struct {
		descripcion  string
		errorDominio error
	}{
		{descripcion: "historia no estimada", errorDominio: domain.ErrItemNoEstimado},
		{descripcion: "historia no pendiente", errorDominio: domain.ErrItemNoPendiente},
		{descripcion: "historia ya asignada", errorDominio: domain.ErrItemYaAsignado},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoUS22(
				registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
			)
			backlog.agregarItem("PRJ-001", itemUS22("HU-04", "Ajustar la configuración", "Media", "Pendiente", 0))
			sprints := nuevoServicioSprintsFalso()
			sprints.agregarSprint("PRJ-001", sprintUS22(
				1, "Primera entrega", domain.EstadoSprintActivo,
				fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
			), 0)
			sprints.asignarErr = caso.errorDominio
			handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

			rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/sprints/1/asignar", url.Values{
				"itemID": {"HU-04"},
			})
			exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
			exigirCuerpoContieneUS22(t, rec, caso.errorDominio.Error())

			if sprints.asignarLlamadas != 1 {
				t.Errorf("AsignarItem se llamó %d veces, se esperaba 1", sprints.asignarLlamadas)
			}
		})
	}
}

// CA2 / R31 / E3: asignar a un Sprint inexistente muestra HTTP 404 con el
// mensaje del dominio.
func TestUS22_Fase2_Sprints_Inexistente404(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	backlog.agregarItem("PRJ-001", itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 3))
	sprints := nuevoServicioSprintsFalso()
	sprints.asignarErr = domain.ErrSprintNoEncontrado
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/sprints/9/asignar", url.Values{
		"itemID": {"HU-01"},
	})
	exigirEstadoUS22(t, rec, http.StatusNotFound)
	exigirCuerpoContieneUS22(t, rec, domain.ErrSprintNoEncontrado.Error())
}

// CA2 / R31 / E18: asignar a un Sprint cerrado muestra HTTP 422 con el mensaje
// del dominio.
func TestUS22_Fase2_Sprints_AsignarSprintCerrado(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	backlog.agregarItem("PRJ-001", itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 3))
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintCerrado,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 8)
	sprints.asignarErr = domain.ErrSprintCerrado
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/sprints/1/asignar", url.Values{
		"itemID": {"HU-01"},
	})
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, domain.ErrSprintCerrado.Error())
}

// CA3 / R32: el formulario de asignación ofrece las historias del backlog del
// proyecto y no las de otros proyectos; la web no filtra por estado ni
// estimación.
func TestUS22_Fase2_Sprints_FormularioAsignacion(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
		registroProyectoUS22("PRJ-002", "Portal de Clientes", fechaUS22(2026, time.April, 15), fechaUS22(2026, time.August, 31)),
	)
	backlog.agregarItem("PRJ-001", itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 3))
	backlog.agregarItem("PRJ-001", itemUS22("HU-02", "Cerrar sesión", "Media", "Completada", 0))
	backlog.agregarItem("PRJ-002", itemUS22("HU-01", "Iniciar sesión del portal", "Alta", "Pendiente", 5))

	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 0)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec,
		`name="itemID"`,
		"/proyectos/PRJ-001/sprints/1/asignar",
		"HU-01", "Iniciar sesión",
		"HU-02", "Cerrar sesión",
	)
	exigirCuerpoNoContieneUS22(t, rec, "Iniciar sesión del portal")

	if backlog.listarLlamadas < 1 {
		t.Errorf("ListarBacklog se llamó %d veces; el formulario debe ofrecer las historias del proyecto", backlog.listarLlamadas)
	}
}

// CA2 / R22 / E11: un número de Sprint no numérico en la ruta produce HTTP 422
// y no se llama al servicio de asignación.
func TestUS22_Fase2_Sprints_AsignarNumeroNoNumerico(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	backlog.agregarItem("PRJ-001", itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 3))
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 0)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/sprints/abc/asignar", url.Values{
		"itemID": {"HU-01"},
	})
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, "El número del sprint debe ser un número entero")

	if sprints.asignarLlamadas != 0 {
		t.Errorf("AsignarItem se llamó %d veces y no debía ejecutarse", sprints.asignarLlamadas)
	}
}

// CA2 / R31 / E2: si AsignarItem devuelve domain.ErrItemNoEncontrado, se muestra
// HTTP 404 con el mensaje del dominio.
func TestUS22_Fase2_Sprints_AsignarItemNoEncontrado404(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	backlog.agregarItem("PRJ-001", itemUS22("HU-99", "Historia inexistente", "Media", "Pendiente", 3))
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 0)
	sprints.asignarErr = domain.ErrItemNoEncontrado
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/sprints/1/asignar", url.Values{
		"itemID": {"HU-99"},
	})
	exigirEstadoUS22(t, rec, http.StatusNotFound)
	exigirCuerpoContieneUS22(t, rec, domain.ErrItemNoEncontrado.Error())
}

// CA1 / R29: la página de Sprints muestra el NOMBRE del proyecto (obtenido con
// Dependencias.Proyectos.ObtenerProyecto) y no su ID.
func TestUS22_Fase2_Sprints_ListarMuestraNombreProyecto(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 0)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "Sistema de Métricas")
	exigirCuerpoNoContieneUS22(t, rec, "Proyecto: PRJ-001")

	if proyectos.obtenerLlamadas != 1 {
		t.Errorf("ObtenerProyecto se llamó %d veces, se esperaba 1", proyectos.obtenerLlamadas)
	}
}
