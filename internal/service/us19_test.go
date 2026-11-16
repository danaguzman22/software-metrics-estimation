package service

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// entornoDashboard agrupa los dobles de proyectos, Sprints, backlog,
// asignaciones, integrantes y esfuerzo para los tests del servicio de
// dashboard (US-19). Los servicios concretos comparten las mismas instancias
// de los dobles para actuar sobre el mismo estado.
type entornoDashboard struct {
	proyectos    *fakeRepositorioProyectos
	sprints      *fakeRepositorioSprints
	backlog      *fakeRepositorioBacklog
	asignaciones *fakeRepositorioAsignaciones
	integrantes  *fakeRepositorioIntegrantes
	esfuerzo     *fakeRepositorioEsfuerzo
}

func nuevoEntornoDashboard() *entornoDashboard {
	return &entornoDashboard{
		proyectos:    nuevoFakeProyectos(),
		sprints:      nuevoFakeSprints(),
		backlog:      nuevoFakeBacklog(),
		asignaciones: nuevoFakeAsignaciones(),
		integrantes:  nuevoFakeIntegrantes(),
		esfuerzo:     nuevoFakeEsfuerzo(),
	}
}

func (e *entornoDashboard) servicio() *ServicioDashboard {
	return NuevoServicioDashboard(
		e.proyectos,
		NuevoServicioSprints(e.proyectos, e.sprints, e.backlog, e.asignaciones),
		NuevoServicioEsfuerzo(
			e.proyectos,
			e.backlog,
			e.integrantes,
			e.sprints,
			e.asignaciones,
			e.esfuerzo,
		),
	)
}

// sembrarProyectoDashboard carga el proyecto "Sistema de Métricas".
func sembrarProyectoDashboard(t *testing.T, entorno *entornoDashboard, proyectoID string) {
	t.Helper()
	entorno.proyectos.proyectos[proyectoID] = proyectoUS02Servicio(
		"Sistema de Métricas",
		domain.Estado("En curso"),
	)
}

// sembrarSprintDashboard carga un Sprint activo o cerrado. Los cerrados
// necesitan el valor de SPPlanificados guardado al cerrar (US-09).
func sembrarSprintDashboard(
	t *testing.T,
	entorno *entornoDashboard,
	proyectoID string,
	numero int,
	estado string,
	spPlanificados *int,
) {
	t.Helper()
	sprint := sprintUS07Servicio(numero, estado)
	sprint.SPPlanificados = spPlanificados
	entorno.sprints.sprints[proyectoID] = append(entorno.sprints.sprints[proyectoID], sprint)
}

// sembrarHistoriaDashboard carga una historia en el backlog del proyecto.
func sembrarHistoriaDashboard(
	t *testing.T,
	entorno *entornoDashboard,
	proyectoID, id string,
	storyPoints int,
	estado string,
) {
	t.Helper()
	entorno.backlog.items[proyectoID] = append(
		entorno.backlog.items[proyectoID],
		itemUS08Servicio(id, storyPoints, estado),
	)
}

// asignarDashboard asigna el número de historia al Sprint sin contar la
// llamada a Agregar ni pasar por sus validaciones (patrón del doble de US-08).
func asignarDashboard(
	t *testing.T,
	entorno *entornoDashboard,
	proyectoID string,
	sprintNumero, itemNumero int,
) {
	t.Helper()
	entorno.asignaciones.sembrar(proyectoID, sprintNumero, itemNumero)
}

// estimarHorasHistoriaDashboard fija las horas estimadas de una historia.
func estimarHorasHistoriaDashboard(
	t *testing.T,
	entorno *entornoDashboard,
	proyectoID, itemID string,
	horas *float64,
) {
	t.Helper()
	for i := range entorno.backlog.items[proyectoID] {
		if entorno.backlog.items[proyectoID][i].ID == itemID {
			entorno.backlog.items[proyectoID][i].HorasEstimadas = horas
			return
		}
	}
	t.Fatalf("no se encontró la historia %q para estimar horas", itemID)
}

// registrarEsfuerzoDashboard agrega un registro de esfuerzo del proyecto,
// opcionalmente atribuido a un Sprint.
func registrarEsfuerzoDashboard(
	entorno *entornoDashboard,
	proyectoID string,
	itemNumero int,
	sprintNum *int,
	horas float64,
) {
	entorno.esfuerzo.registros = append(
		entorno.esfuerzo.registros,
		registroEsfuerzoFake{
			proyectoID: proyectoID,
			itemNumero: itemNumero,
			sprintNum:  sprintNum,
			registro:   domain.RegistroEsfuerzo{Horas: horas},
		},
	)
}

// R2 / R3 / R4 / CA1 (US-19): el servicio arma los totales sobre TODOS los
// Sprints del proyecto (D10): SP planificados y completados, velocidad
// promedio, porcentaje de historias completadas, horas totales con su
// desviación y cantidad de Sprints cerrados.
func TestUS19_Fase1_ServicioDashboard_ArmaMetricas(t *testing.T) {
	entorno := nuevoEntornoDashboard()
	sembrarProyectoDashboard(t, entorno, "p1")
	// Sprint 1 cerrado con 8 SP completados; planificados guardados al cerrar.
	sprint1 := 13
	sembrarSprintDashboard(t, entorno, "p1", 1, domain.EstadoSprintCerrado, &sprint1)
	sembrarHistoriaDashboard(t, entorno, "p1", "HU-01", 3, domain.EstadoCompletada)
	sembrarHistoriaDashboard(t, entorno, "p1", "HU-02", 5, domain.EstadoCompletada)
	asignarDashboard(t, entorno, "p1", 1, 1)
	asignarDashboard(t, entorno, "p1", 1, 2)
	// Sprint 2 cerrado con 21 SP completados; planificados guardados al cerrar
	// (30 > 21 porque el cierre retira las asignaciones no completadas).
	sprint2 := 30
	sembrarSprintDashboard(t, entorno, "p1", 2, domain.EstadoSprintCerrado, &sprint2)
	sembrarHistoriaDashboard(t, entorno, "p1", "HU-03", 8, domain.EstadoCompletada)
	sembrarHistoriaDashboard(t, entorno, "p1", "HU-04", 13, domain.EstadoCompletada)
	asignarDashboard(t, entorno, "p1", 2, 3)
	asignarDashboard(t, entorno, "p1", 2, 4)
	// Sprint 3 activo con 8 SP pendientes: los totales también lo suman (D10).
	sembrarSprintDashboard(t, entorno, "p1", 3, domain.EstadoSprintActivo, nil)
	sembrarHistoriaDashboard(t, entorno, "p1", "HU-05", 8, domain.EstadoItemPendiente)
	asignarDashboard(t, entorno, "p1", 3, 5)
	// Horas: estimaciones completas, 10 estimadas y 12 reales.
	estimaran := []float64{1, 1, 3, 3, 2}
	estimarHorasHistoriaDashboard(t, entorno, "p1", "HU-01", &estimaran[0])
	estimarHorasHistoriaDashboard(t, entorno, "p1", "HU-02", &estimaran[1])
	estimarHorasHistoriaDashboard(t, entorno, "p1", "HU-03", &estimaran[2])
	estimarHorasHistoriaDashboard(t, entorno, "p1", "HU-04", &estimaran[3])
	estimarHorasHistoriaDashboard(t, entorno, "p1", "HU-05", &estimaran[4])
	registrarEsfuerzoDashboard(entorno, "p1", 1, intPuntero(1), 4)
	registrarEsfuerzoDashboard(entorno, "p1", 3, nil, 8)

	dashboard, err := entorno.servicio().ObtenerDashboard("p1")
	if err != nil {
		t.Fatalf("ObtenerDashboard() devolvió un error inesperado: %v", err)
	}
	if dashboard.SPPlanificadosTotal != 51 {
		t.Errorf("SPPlanificadosTotal = %d, se esperaba 51 (13+30+8; en los cerrados manda el guardado del cierre)",
			dashboard.SPPlanificadosTotal)
	}
	if dashboard.SPCompletadosTotal != 29 {
		t.Errorf("SPCompletadosTotal = %d, se esperaba 29 (8+21+0)", dashboard.SPCompletadosTotal)
	}
	if dashboard.VelocidadPromedio != 14.5 {
		t.Errorf("VelocidadPromedio = %v, se esperaba 14.5 (promedio de 8 y 21)", dashboard.VelocidadPromedio)
	}
	if dashboard.PorcentajeHistoriasCompletadas != 80 {
		t.Errorf("PorcentajeHistoriasCompletadas = %v, se esperaba 80 (4 de 5)",
			dashboard.PorcentajeHistoriasCompletadas)
	}
	if dashboard.HorasEstimadasTotal != 10 {
		t.Errorf("HorasEstimadasTotal = %v, se esperaba 10", dashboard.HorasEstimadasTotal)
	}
	if dashboard.HorasRealesTotal != 12 {
		t.Errorf("HorasRealesTotal = %v, se esperaba 12", dashboard.HorasRealesTotal)
	}
	if !dashboard.EstimacionesCompletas {
		t.Errorf("EstimacionesCompletas = false, se esperaba true")
	}
	if dashboard.DesviacionHorasTotal == nil {
		t.Fatal("DesviacionHorasTotal = nil, se esperaba 20")
	} else if *dashboard.DesviacionHorasTotal != 20 {
		t.Errorf("DesviacionHorasTotal = %v, se esperaba 20", *dashboard.DesviacionHorasTotal)
	}
	if dashboard.SprintsCerrados != 2 {
		t.Errorf("SprintsCerrados = %d, se esperaba 2", dashboard.SprintsCerrados)
	}
	if entorno.asignaciones.listarLlamadas < 3 {
		t.Errorf("ListarNumerosItem fue invocado %d veces, se esperaban al menos 3 (una por Sprint en StoryPointsCompletados)",
			entorno.asignaciones.listarLlamadas)
	}
}

// R16 / CA3 (US-19): un proyecto sin Sprints ni historias devuelve el
// dashboard en cero, sin errores ni división por cero; las listas vienen
// vacías y la desviación no está disponible.
func TestUS19_Fase1_ServicioDashboard_SinSprintsNiHistorias(t *testing.T) {
	entorno := nuevoEntornoDashboard()
	sembrarProyectoDashboard(t, entorno, "p1")

	dashboard, err := entorno.servicio().ObtenerDashboard("p1")
	if err != nil {
		t.Fatalf("ObtenerDashboard() devolvió un error inesperado: %v", err)
	}
	if dashboard.SPPlanificadosTotal != 0 {
		t.Errorf("SPPlanificadosTotal = %d, se esperaba 0", dashboard.SPPlanificadosTotal)
	}
	if dashboard.SPCompletadosTotal != 0 {
		t.Errorf("SPCompletadosTotal = %d, se esperaba 0", dashboard.SPCompletadosTotal)
	}
	if dashboard.VelocidadPromedio != 0 {
		t.Errorf("VelocidadPromedio = %v, se esperaba 0 (sin división por cero)", dashboard.VelocidadPromedio)
	}
	if dashboard.PorcentajeHistoriasCompletadas != 0 {
		t.Errorf("PorcentajeHistoriasCompletadas = %v, se esperaba 0", dashboard.PorcentajeHistoriasCompletadas)
	}
	if dashboard.HorasEstimadasTotal != 0 {
		t.Errorf("HorasEstimadasTotal = %v, se esperaba 0", dashboard.HorasEstimadasTotal)
	}
	if dashboard.HorasRealesTotal != 0 {
		t.Errorf("HorasRealesTotal = %v, se esperaba 0", dashboard.HorasRealesTotal)
	}
	if !dashboard.EstimacionesCompletas {
		t.Errorf("EstimacionesCompletas = false, se esperaba true (sin historias no hay pendientes de estimar)")
	}
	if dashboard.DesviacionHorasTotal != nil {
		t.Errorf("DesviacionHorasTotal = %v, se esperaba nil con estimaciones que suman cero",
			*dashboard.DesviacionHorasTotal)
	}
	if dashboard.SprintsCerrados != 0 {
		t.Errorf("SprintsCerrados = %d, se esperaba 0", dashboard.SprintsCerrados)
	}
	if len(dashboard.VelocidadPorSprint) != 0 {
		t.Errorf("VelocidadPorSprint = %+v, se esperaba vacía", dashboard.VelocidadPorSprint)
	}
	if len(dashboard.HorasPorSprint) != 0 {
		t.Errorf("HorasPorSprint = %+v, se esperaba vacía", dashboard.HorasPorSprint)
	}
}

// D10 / CA3 (US-19): con Sprints solo activos ningún Sprint cerrado figura en
// las listas; los totales sí suman los datos del activo y la velocidad
// promedio es cero.
func TestUS19_Fase1_ServicioDashboard_SoloSprintsActivos(t *testing.T) {
	entorno := nuevoEntornoDashboard()
	sembrarProyectoDashboard(t, entorno, "p1")
	sembrarSprintDashboard(t, entorno, "p1", 1, domain.EstadoSprintActivo, nil)
	sembrarHistoriaDashboard(t, entorno, "p1", "HU-01", 5, domain.EstadoItemPendiente)
	estimaran := 4.0
	estimarHorasHistoriaDashboard(t, entorno, "p1", "HU-01", &estimaran)
	asignarDashboard(t, entorno, "p1", 1, 1)
	registrarEsfuerzoDashboard(entorno, "p1", 1, intPuntero(1), 6)

	dashboard, err := entorno.servicio().ObtenerDashboard("p1")
	if err != nil {
		t.Fatalf("ObtenerDashboard() devolvió un error inesperado: %v", err)
	}
	if dashboard.SPPlanificadosTotal != 5 {
		t.Errorf("SPPlanificadosTotal = %d, se esperaba 5 (los activos también suman, D10)",
			dashboard.SPPlanificadosTotal)
	}
	if dashboard.SPCompletadosTotal != 0 {
		t.Errorf("SPCompletadosTotal = %d, se esperaba 0", dashboard.SPCompletadosTotal)
	}
	if dashboard.VelocidadPromedio != 0 {
		t.Errorf("VelocidadPromedio = %v, se esperaba 0 (sin Sprints cerrados)", dashboard.VelocidadPromedio)
	}
	if dashboard.PorcentajeHistoriasCompletadas != 0 {
		t.Errorf("PorcentajeHistoriasCompletadas = %v, se esperaba 0 (0 de 1)",
			dashboard.PorcentajeHistoriasCompletadas)
	}
	if dashboard.HorasEstimadasTotal != 4 {
		t.Errorf("HorasEstimadasTotal = %v, se esperaba 4", dashboard.HorasEstimadasTotal)
	}
	if dashboard.HorasRealesTotal != 6 {
		t.Errorf("HorasRealesTotal = %v, se esperaba 6", dashboard.HorasRealesTotal)
	}
	if !dashboard.EstimacionesCompletas {
		t.Errorf("EstimacionesCompletas = false, se esperaba true")
	}
	if dashboard.DesviacionHorasTotal == nil {
		t.Fatal("DesviacionHorasTotal = nil, se esperaba 50")
	} else if *dashboard.DesviacionHorasTotal != 50 {
		t.Errorf("DesviacionHorasTotal = %v, se esperaba 50", *dashboard.DesviacionHorasTotal)
	}
	if dashboard.SprintsCerrados != 0 {
		t.Errorf("SprintsCerrados = %d, se esperaba 0", dashboard.SprintsCerrados)
	}
	if len(dashboard.VelocidadPorSprint) != 0 {
		t.Errorf("VelocidadPorSprint = %+v, se esperaba vacía (solo cerrados)", dashboard.VelocidadPorSprint)
	}
	if len(dashboard.HorasPorSprint) != 0 {
		t.Errorf("HorasPorSprint = %+v, se esperaba vacía (solo cerrados)", dashboard.HorasPorSprint)
	}
}

// R3 / CA1 (US-19): VelocidadPromedio y PorcentajeHistoriasCompletadas se
// exponen tal como las devuelve US-16 (el dominio redondea el porcentaje a
// dos decimales); el servicio no las recalcula ni las redondea.
func TestUS19_Fase1_ServicioDashboard_VelocidadPromedioYPorcentajeTalComoUS16(t *testing.T) {
	casos := []struct {
		descripcion        string
		preparar           func(t *testing.T, entorno *entornoDashboard)
		promedioEsperado   float64
		porcentajeEsperado float64
	}{
		{
			descripcion: "sin Sprints cerrados devuelve cero sin error",
			preparar: func(t *testing.T, entorno *entornoDashboard) {
				sembrarSprintDashboard(t, entorno, "p1", 1, domain.EstadoSprintActivo, nil)
				sembrarHistoriaDashboard(t, entorno, "p1", "HU-01", 3, domain.EstadoItemPendiente)
				asignarDashboard(t, entorno, "p1", 1, 1)
			},
			promedioEsperado:   0,
			porcentajeEsperado: 0,
		},
		{
			descripcion: "promedio fraccionario con cerrados y porcentaje con dos decimales",
			preparar: func(t *testing.T, entorno *entornoDashboard) {
				sprint1 := 8
				sembrarSprintDashboard(t, entorno, "p1", 1, domain.EstadoSprintCerrado, &sprint1)
				sprint2 := 13
				sembrarSprintDashboard(t, entorno, "p1", 2, domain.EstadoSprintCerrado, &sprint2)
				sembrarHistoriaDashboard(t, entorno, "p1", "HU-01", 8, domain.EstadoCompletada)
				sembrarHistoriaDashboard(t, entorno, "p1", "HU-02", 13, domain.EstadoCompletada)
				sembrarHistoriaDashboard(t, entorno, "p1", "HU-03", 5, domain.EstadoItemPendiente)
				asignarDashboard(t, entorno, "p1", 1, 1)
				asignarDashboard(t, entorno, "p1", 2, 2)
			},
			promedioEsperado:   10.5,
			porcentajeEsperado: 66.67,
		},
		{
			descripcion: "sin historias el porcentaje es cero y el promedio es cero",
			preparar: func(t *testing.T, entorno *entornoDashboard) {
				sprint1 := 8
				sembrarSprintDashboard(t, entorno, "p1", 1, domain.EstadoSprintCerrado, &sprint1)
			},
			promedioEsperado:   0,
			porcentajeEsperado: 0,
		},
	}
	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoDashboard()
			sembrarProyectoDashboard(t, entorno, "p1")
			caso.preparar(t, entorno)

			dashboard, err := entorno.servicio().ObtenerDashboard("p1")
			if err != nil {
				t.Fatalf("ObtenerDashboard() devolvió un error inesperado: %v", err)
			}
			if dashboard.VelocidadPromedio != caso.promedioEsperado {
				t.Errorf("VelocidadPromedio = %v, se esperaba %v", dashboard.VelocidadPromedio, caso.promedioEsperado)
			}
			if dashboard.PorcentajeHistoriasCompletadas != caso.porcentajeEsperado {
				t.Errorf("PorcentajeHistoriasCompletadas = %v, se esperaba %v",
					dashboard.PorcentajeHistoriasCompletadas, caso.porcentajeEsperado)
			}
		})
	}
}

// R4 / CA1 (US-19): con estimaciones incompletas las horas estimadas suman
// solo las disponibles, EstimacionesCompletas es false y la desviación no
// está disponible (nil); lo demás se devuelve con normalidad.
func TestUS19_Fase1_ServicioDashboard_EstimacionesIncompletas(t *testing.T) {
	entorno := nuevoEntornoDashboard()
	sembrarProyectoDashboard(t, entorno, "p1")
	sembrarSprintDashboard(t, entorno, "p1", 1, domain.EstadoSprintActivo, nil)
	sembrarHistoriaDashboard(t, entorno, "p1", "HU-01", 3, domain.EstadoItemPendiente)
	sembrarHistoriaDashboard(t, entorno, "p1", "HU-02", 5, domain.EstadoItemPendiente)
	estimaran := 4.0
	estimarHorasHistoriaDashboard(t, entorno, "p1", "HU-01", &estimaran)
	sembrarHistoriaDashboard(t, entorno, "p1", "HU-03", 8, domain.EstadoItemPendiente)
	estimaran3 := 0.0
	estimarHorasHistoriaDashboard(t, entorno, "p1", "HU-03", &estimaran3)
	asignarDashboard(t, entorno, "p1", 1, 1)
	asignarDashboard(t, entorno, "p1", 1, 2)
	registrarEsfuerzoDashboard(entorno, "p1", 1, intPuntero(1), 3)

	dashboard, err := entorno.servicio().ObtenerDashboard("p1")
	if err != nil {
		t.Fatalf("ObtenerDashboard() devolvió un error inesperado: %v", err)
	}
	if dashboard.EstimacionesCompletas {
		t.Errorf("EstimacionesCompletas = true, se esperaba false (HU-02 sin estimar)")
	}
	if dashboard.HorasEstimadasTotal != 4 {
		t.Errorf("HorasEstimadasTotal = %v, se esperaba 4 (solo las estimaciones disponibles)",
			dashboard.HorasEstimadasTotal)
	}
	if dashboard.HorasRealesTotal != 3 {
		t.Errorf("HorasRealesTotal = %v, se esperaba 3", dashboard.HorasRealesTotal)
	}
	if dashboard.DesviacionHorasTotal != nil {
		t.Errorf("DesviacionHorasTotal = %v, se esperaba nil con estimaciones incompletas",
			*dashboard.DesviacionHorasTotal)
	}
}

// R4 / CA1 (US-19): con la suma de estimaciones en cero la desviación no
// está disponible (nil), aunque las estimaciones estén completas (US-17 D4).
func TestUS19_Fase1_ServicioDashboard_DesviacionNilConSumaCero(t *testing.T) {
	entorno := nuevoEntornoDashboard()
	sembrarProyectoDashboard(t, entorno, "p1")
	sembrarSprintDashboard(t, entorno, "p1", 1, domain.EstadoSprintActivo, nil)
	sembrarHistoriaDashboard(t, entorno, "p1", "HU-01", 3, domain.EstadoItemPendiente)
	sembrarHistoriaDashboard(t, entorno, "p1", "HU-02", 5, domain.EstadoItemPendiente)
	cerrero := 0.0
	estimarHorasHistoriaDashboard(t, entorno, "p1", "HU-01", &cerrero)
	estimarHorasHistoriaDashboard(t, entorno, "p1", "HU-02", &cerrero)
	registrarEsfuerzoDashboard(entorno, "p1", 1, nil, 5)

	dashboard, err := entorno.servicio().ObtenerDashboard("p1")
	if err != nil {
		t.Fatalf("ObtenerDashboard() devolvió un error inesperado: %v", err)
	}
	if !dashboard.EstimacionesCompletas {
		t.Errorf("EstimacionesCompletas = false, se esperaba true (todas con estimación explícita)")
	}
	if dashboard.HorasEstimadasTotal != 0 {
		t.Errorf("HorasEstimadasTotal = %v, se esperaba 0", dashboard.HorasEstimadasTotal)
	}
	if dashboard.HorasRealesTotal != 5 {
		t.Errorf("HorasRealesTotal = %v, se esperaba 5", dashboard.HorasRealesTotal)
	}
	if dashboard.DesviacionHorasTotal != nil {
		t.Errorf("DesviacionHorasTotal = %v, se esperaba nil con la suma de estimaciones en cero",
			*dashboard.DesviacionHorasTotal)
	}
}

// R4 / CA1 (US-19): sin registros de esfuerzo las horas reales son cero; con
// estimaciones completas y mayores que cero la desviación es -100 (comportamiento
// de domain.CalcularDesviacion).
func TestUS19_Fase1_ServicioDashboard_SinRegistrosDeEsfuerzo(t *testing.T) {
	entorno := nuevoEntornoDashboard()
	sembrarProyectoDashboard(t, entorno, "p1")
	sembrarSprintDashboard(t, entorno, "p1", 1, domain.EstadoSprintActivo, nil)
	sembrarHistoriaDashboard(t, entorno, "p1", "HU-01", 3, domain.EstadoItemPendiente)
	estimaran := 10.0
	estimarHorasHistoriaDashboard(t, entorno, "p1", "HU-01", &estimaran)

	dashboard, err := entorno.servicio().ObtenerDashboard("p1")
	if err != nil {
		t.Fatalf("ObtenerDashboard() devolvió un error inesperado: %v", err)
	}
	if dashboard.HorasEstimadasTotal != 10 {
		t.Errorf("HorasEstimadasTotal = %v, se esperaba 10", dashboard.HorasEstimadasTotal)
	}
	if dashboard.HorasRealesTotal != 0 {
		t.Errorf("HorasRealesTotal = %v, se esperaba 0 (sin registros de esfuerzo)",
			dashboard.HorasRealesTotal)
	}
	if !dashboard.EstimacionesCompletas {
		t.Errorf("EstimacionesCompletas = false, se esperaba true")
	}
	if dashboard.DesviacionHorasTotal == nil {
		t.Fatal("DesviacionHorasTotal = nil, se esperaba -100 (reales 0 sobre estimadas 10)")
	} else if *dashboard.DesviacionHorasTotal != -100 {
		t.Errorf("DesviacionHorasTotal = %v, se esperaba -100", *dashboard.DesviacionHorasTotal)
	}
}

// R1 / V1 / CA1 (US-19): un proyecto inexistente devuelve
// domain.ErrProyectoNoEncontrado sin envolver; no se consultan los Sprints ni
// el resto de los dobles.
func TestUS19_Fase1_ServicioDashboard_ProyectoInexistente(t *testing.T) {
	entorno := nuevoEntornoDashboard()

	dashboard, err := entorno.servicio().ObtenerDashboard("PRJ-999")
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if err != domain.ErrProyectoNoEncontrado {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrProyectoNoEncontrado)
	}
	if !reflect.DeepEqual(dashboard, DashboardProyecto{}) {
		t.Errorf("dashboard = %+v, se esperaba DashboardProyecto{}", dashboard)
	}
	if entorno.proyectos.existeLlamadas != 1 {
		t.Errorf("Existe fue invocado %d veces, se esperaba 1", entorno.proyectos.existeLlamadas)
	}
	if entorno.sprints.listarLlamadas != 0 {
		t.Errorf("ListarPorProyecto (Sprints) fue invocado %d veces; el proyecto debía validarse primero",
			entorno.sprints.listarLlamadas)
	}
	if entorno.asignaciones.listarLlamadas != 0 {
		t.Errorf("ListarNumerosItem fue invocado %d veces; nunca debe llamarse sin un proyecto inexistente",
			entorno.asignaciones.listarLlamadas)
	}
	if entorno.esfuerzo.agregarLlamadas != 0 {
		t.Errorf("Agregar (esfuerzo) fue invocado %d veces; nunca debe llamarse",
			entorno.esfuerzo.agregarLlamadas)
	}
}

// R5 / CA1 / CA2 (US-19): VelocidadPorSprint y HorasPorSprint contienen solo
// los Sprints cerrados, en orden ascendente por número; el activo no figura.
func TestUS19_Fase1_ServicioDashboard_CerradosEnOrdenAscendente(t *testing.T) {
	entorno := nuevoEntornoDashboard()
	sembrarProyectoDashboard(t, entorno, "p1")
	// Sprint 1 cerrado: 8 SP completados.
	sprint1 := 8
	sembrarSprintDashboard(t, entorno, "p1", 1, domain.EstadoSprintCerrado, &sprint1)
	sembrarHistoriaDashboard(t, entorno, "p1", "HU-01", 8, domain.EstadoCompletada)
	estimarHorasHistoriaDashboard(t, entorno, "p1", "HU-01", horasPuntero(10))
	asignarDashboard(t, entorno, "p1", 1, 1)
	// Sprint 2 activo con 3 SP pendientes: no debe figurar en las series.
	sembrarSprintDashboard(t, entorno, "p1", 2, domain.EstadoSprintActivo, nil)
	sembrarHistoriaDashboard(t, entorno, "p1", "HU-02", 3, domain.EstadoItemPendiente)
	estimarHorasHistoriaDashboard(t, entorno, "p1", "HU-02", horasPuntero(100))
	asignarDashboard(t, entorno, "p1", 2, 2)
	// Sprint 3 cerrado: 5 SP completados.
	sprint3 := 5
	sembrarSprintDashboard(t, entorno, "p1", 3, domain.EstadoSprintCerrado, &sprint3)
	sembrarHistoriaDashboard(t, entorno, "p1", "HU-03", 5, domain.EstadoCompletada)
	estimarHorasHistoriaDashboard(t, entorno, "p1", "HU-03", horasPuntero(20))
	asignarDashboard(t, entorno, "p1", 3, 3)
	// Snapshots de esfuerzo de los cerrados (US-17).
	entorno.esfuerzo.snapshots["p1:1"] = snapshotEsfuerzoFake{
		disponible: true,
		items: []domain.ItemSprintEsfuerzo{
			{ItemNumero: 1, HorasEstimadas: horasPuntero(10)},
		},
	}
	entorno.esfuerzo.snapshots["p1:3"] = snapshotEsfuerzoFake{
		disponible: true,
		items: []domain.ItemSprintEsfuerzo{
			{ItemNumero: 3, HorasEstimadas: horasPuntero(20)},
		},
	}
	registrarEsfuerzoDashboard(entorno, "p1", 1, intPuntero(1), 12)
	registrarEsfuerzoDashboard(entorno, "p1", 3, intPuntero(3), 25)

	dashboard, err := entorno.servicio().ObtenerDashboard("p1")
	if err != nil {
		t.Fatalf("ObtenerDashboard() devolvió un error inesperado: %v", err)
	}
	if dashboard.SprintsCerrados != 2 {
		t.Errorf("SprintsCerrados = %d, se esperaba 2 (solo cerrados)", dashboard.SprintsCerrados)
	}
	esperadoVelocidad := []PuntoVelocidadSprint{
		{SprintNumero: 1, Velocidad: 8},
		{SprintNumero: 3, Velocidad: 5},
	}
	if !reflect.DeepEqual(dashboard.VelocidadPorSprint, esperadoVelocidad) {
		t.Errorf("VelocidadPorSprint = %+v, se esperaba %+v (solo cerrados, orden ascendente)",
			dashboard.VelocidadPorSprint, esperadoVelocidad)
	}
	esperadoHoras := []PuntoHorasSprint{
		{SprintNumero: 1, HorasEstimadas: 10, HorasReales: 12, EstimacionesCompletas: true, MetricasDisponibles: true},
		{SprintNumero: 3, HorasEstimadas: 20, HorasReales: 25, EstimacionesCompletas: true, MetricasDisponibles: true},
	}
	if !reflect.DeepEqual(dashboard.HorasPorSprint, esperadoHoras) {
		t.Errorf("HorasPorSprint = %+v, se esperaba %+v (solo cerrados, orden ascendente)",
			dashboard.HorasPorSprint, esperadoHoras)
	}
	for _, punto := range dashboard.VelocidadPorSprint {
		if punto.SprintNumero == 2 {
			t.Errorf("VelocidadPorSprint incluye el Sprint 2 activo; solo debe incluir cerrados")
		}
	}
	for _, punto := range dashboard.HorasPorSprint {
		if punto.SprintNumero == 2 {
			t.Errorf("HorasPorSprint incluye el Sprint 2 activo; solo debe incluir cerrados")
		}
	}
}

// R6 / D12 / CA2 / CA3 (US-19): un Sprint cerrado sin snapshot de esfuerzo
// (domain.ErrMetricasSprintNoDisponibles) no produce error: la entrada de
// HorasPorSprint queda con MetricasDisponibles=false y horas en cero; el
// resto del dashboard se arma con normalidad.
func TestUS19_Fase1_ServicioDashboard_SprintCerradoSinSnapshot(t *testing.T) {
	entorno := nuevoEntornoDashboard()
	sembrarProyectoDashboard(t, entorno, "p1")
	// Sprint 1 cerrado sin snapshot de esfuerzo.
	sprint1 := 8
	sembrarSprintDashboard(t, entorno, "p1", 1, domain.EstadoSprintCerrado, &sprint1)
	sembrarHistoriaDashboard(t, entorno, "p1", "HU-01", 8, domain.EstadoCompletada)
	estimarHorasHistoriaDashboard(t, entorno, "p1", "HU-01", horasPuntero(10))
	asignarDashboard(t, entorno, "p1", 1, 1)
	registrarEsfuerzoDashboard(entorno, "p1", 1, intPuntero(1), 12)
	// Sprint 2 cerrado con snapshot de esfuerzo.
	sprint2 := 13
	sembrarSprintDashboard(t, entorno, "p1", 2, domain.EstadoSprintCerrado, &sprint2)
	sembrarHistoriaDashboard(t, entorno, "p1", "HU-02", 13, domain.EstadoCompletada)
	estimarHorasHistoriaDashboard(t, entorno, "p1", "HU-02", horasPuntero(15))
	asignarDashboard(t, entorno, "p1", 2, 2)
	entorno.esfuerzo.snapshots["p1:2"] = snapshotEsfuerzoFake{
		disponible: true,
		items: []domain.ItemSprintEsfuerzo{
			{ItemNumero: 2, HorasEstimadas: horasPuntero(15)},
		},
	}
	registrarEsfuerzoDashboard(entorno, "p1", 2, intPuntero(2), 20)

	dashboard, err := entorno.servicio().ObtenerDashboard("p1")
	if err != nil {
		t.Fatalf("ObtenerDashboard() devolvió un error inesperado: %v (el Sprint sin snapshot no debe fallar)", err)
	}
	if dashboard.SprintsCerrados != 2 {
		t.Errorf("SprintsCerrados = %d, se esperaba 2", dashboard.SprintsCerrados)
	}
	if len(dashboard.HorasPorSprint) != 2 {
		t.Fatalf("len(HorasPorSprint) = %d, se esperaban 2 entradas", len(dashboard.HorasPorSprint))
	}
	sinSnapshot := dashboard.HorasPorSprint[0]
	if sinSnapshot.SprintNumero != 1 {
		t.Errorf("HorasPorSprint[0].SprintNumero = %d, se esperaba 1", sinSnapshot.SprintNumero)
	}
	if sinSnapshot.MetricasDisponibles {
		t.Errorf("HorasPorSprint[0].MetricasDisponibles = true, se esperaba false (cerrado sin snapshot)")
	}
	if sinSnapshot.HorasEstimadas != 0 || sinSnapshot.HorasReales != 0 {
		t.Errorf("HorasPorSprint[0] = %+v, se esperaban horas en 0 sin snapshot", sinSnapshot)
	}
	if sinSnapshot.EstimacionesCompletas {
		t.Errorf("HorasPorSprint[0].EstimacionesCompletas = true, se esperaba false sin snapshot")
	}
	conSnapshot := dashboard.HorasPorSprint[1]
	if conSnapshot.SprintNumero != 2 {
		t.Errorf("HorasPorSprint[1].SprintNumero = %d, se esperaba 2", conSnapshot.SprintNumero)
	}
	if !conSnapshot.MetricasDisponibles {
		t.Errorf("HorasPorSprint[1].MetricasDisponibles = false, se esperaba true")
	}
	if conSnapshot.HorasEstimadas != 15 || conSnapshot.HorasReales != 20 {
		t.Errorf("HorasPorSprint[1] = %+v, se esperaban 15 estimadas y 20 reales", conSnapshot)
	}
	// La velocidad no depende del snapshot: el Sprint 1 aparece igualmente.
	if len(dashboard.VelocidadPorSprint) != 2 {
		t.Fatalf("len(VelocidadPorSprint) = %d, se esperaban 2 entradas", len(dashboard.VelocidadPorSprint))
	}
	if dashboard.VelocidadPorSprint[0] != (PuntoVelocidadSprint{SprintNumero: 1, Velocidad: 8}) {
		t.Errorf("VelocidadPorSprint[0] = %+v, se esperaba {1 8}", dashboard.VelocidadPorSprint[0])
	}
	if dashboard.VelocidadPorSprint[1] != (PuntoVelocidadSprint{SprintNumero: 2, Velocidad: 13}) {
		t.Errorf("VelocidadPorSprint[1] = %+v, se esperaba {2 13}", dashboard.VelocidadPorSprint[1])
	}
}

// R8 / E2 / CA1 (US-19): los errores de infraestructura o de dato
// inconsistente no reconocidos se propagan envueltos con %w (prefijo
// "service:"), conservando la causa con errors.Is y el contexto en el mensaje.
func TestUS19_Fase1_ServicioDashboard_ErrorDeRepositorioEnvuelto(t *testing.T) {
	causa := errors.New("fallo simulado del repositorio")
	cerrarSprint := func(t *testing.T, e *entornoDashboard, numero int) {
		sp := 8
		sembrarSprintDashboard(t, e, "p1", numero, domain.EstadoSprintCerrado, &sp)
	}
	casos := []struct {
		descripcion string
		preparar    func(t *testing.T, e *entornoDashboard)
		contexto    string
	}{
		{
			descripcion: "error al verificar el proyecto",
			preparar:    func(t *testing.T, e *entornoDashboard) { e.proyectos.existeErr = causa },
			contexto:    "verificar proyecto",
		},
		{
			descripcion: "error al listar los Sprints del proyecto",
			preparar:    func(t *testing.T, e *entornoDashboard) { e.sprints.listarErr = causa },
			contexto:    "listar sprints del proyecto",
		},
		{
			descripcion: "error al obtener el Sprint al sumar los planificados",
			preparar: func(t *testing.T, e *entornoDashboard) {
				sembrarSprintDashboard(t, e, "p1", 1, domain.EstadoSprintActivo, nil)
				e.sprints.obtenerErr = causa
			},
			contexto: "verificar sprint",
		},
		{
			descripcion: "error al listar las asignaciones del Sprint al sumar los planificados",
			preparar: func(t *testing.T, e *entornoDashboard) {
				sembrarSprintDashboard(t, e, "p1", 1, domain.EstadoSprintActivo, nil)
				e.asignaciones.listarErr = causa
			},
			contexto: "listar ítems del sprint",
		},
		{
			descripcion: "error al listar los ítems del proyecto al sumar los completados",
			preparar: func(t *testing.T, e *entornoDashboard) {
				cerrarSprint(t, e, 1)
				e.backlog.listarErr = causa
			},
			contexto: "listar ítems del proyecto",
		},
		{
			descripcion: "error al obtener el snapshot del Sprint",
			preparar: func(t *testing.T, e *entornoDashboard) {
				cerrarSprint(t, e, 1)
				e.esfuerzo.snapshots["p1:1"] = snapshotEsfuerzoFake{disponible: true, err: causa}
			},
			contexto: "obtener snapshot del sprint",
		},
		{
			descripcion: "error al sumar horas reales del Sprint",
			preparar: func(t *testing.T, e *entornoDashboard) {
				cerrarSprint(t, e, 1)
				e.esfuerzo.snapshots["p1:1"] = snapshotEsfuerzoFake{disponible: true}
				e.esfuerzo.horasSprintErr = causa
			},
			contexto: "sumar horas reales del sprint",
		},
		{
			descripcion: "error al sumar horas reales del proyecto",
			preparar: func(t *testing.T, e *entornoDashboard) {
				cerrarSprint(t, e, 1)
				e.esfuerzo.snapshots["p1:1"] = snapshotEsfuerzoFake{disponible: true}
				e.esfuerzo.horasProyectoErr = causa
			},
			contexto: "sumar horas reales del proyecto",
		},
	}
	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoDashboard()
			sembrarProyectoDashboard(t, entorno, "p1")
			caso.preparar(t, entorno)

			dashboard, err := entorno.servicio().ObtenerDashboard("p1")
			if err == nil {
				t.Fatal("se esperaba un error de repositorio y se obtuvo nil")
			}
			if !errors.Is(err, causa) {
				t.Fatalf("errors.Is(err, causa) = false; err = %v", err)
			}
			if err == causa {
				t.Errorf("err = %v, se esperaba un error envuelto, no el error del repositorio tal cual", err)
			}
			if !strings.HasPrefix(err.Error(), "service:") {
				t.Errorf("mensaje = %q, se esperaba el prefijo service:", err.Error())
			}
			if !strings.Contains(err.Error(), caso.contexto) {
				t.Errorf("mensaje = %q, se esperaba el contexto %q", err.Error(), caso.contexto)
			}
			if errors.Is(err, domain.ErrProyectoNoEncontrado) {
				t.Errorf("err = %v, un error de infraestructura no debe mapearse a ErrProyectoNoEncontrado", err)
			}
			if !reflect.DeepEqual(dashboard, DashboardProyecto{}) {
				t.Errorf("dashboard = %+v, se esperaba DashboardProyecto{} (nunca un éxito parcial)", dashboard)
			}
		})
	}
}

// E2 (US-19): un Sprint cerrado sin SPPlanificados guardado devuelve un error
// de dato no reconocido que se propaga; no se mapea a un error de dominio.
func TestUS19_Fase1_ServicioDashboard_SprintCerradoSinSPPlanificadosPropagaError(t *testing.T) {
	entorno := nuevoEntornoDashboard()
	sembrarProyectoDashboard(t, entorno, "p1")
	sembrarSprintDashboard(t, entorno, "p1", 1, domain.EstadoSprintCerrado, nil)

	dashboard, err := entorno.servicio().ObtenerDashboard("p1")
	if err == nil {
		t.Fatal("se esperaba un error por el Sprint cerrado sin SPPlanificados guardado y se obtuvo nil")
	}
	for _, errorDominio := range []error{
		domain.ErrProyectoNoEncontrado,
		domain.ErrSprintNoEncontrado,
		domain.ErrSprintCerrado,
		domain.ErrMetricasSprintNoDisponibles,
	} {
		if errors.Is(err, errorDominio) {
			t.Errorf("err = %v, un error de dato no reconocido no debe mapearse a %v", err, errorDominio)
		}
	}
	if !strings.Contains(err.Error(), "story points planificados guardados") {
		t.Errorf("err = %q, se esperaba el detalle del SPPlanificados faltante", err.Error())
	}
	if !reflect.DeepEqual(dashboard, DashboardProyecto{}) {
		t.Errorf("dashboard = %+v, se esperaba DashboardProyecto{} (nunca un éxito parcial)", dashboard)
	}
}

// V1 / V2 / CA1 (US-19): orden de evaluación — proyecto (V1) antes de
// cualquier otra llamada; el error de más arriba prevalece sobre el de más
// abajo.
func TestUS19_Fase1_ServicioDashboard_PrecedenciaValidaciones(t *testing.T) {
	causa := errors.New("fallo simulado del repositorio")
	casos := []struct {
		descripcion string
		preparar    func(t *testing.T, e *entornoDashboard)
		exacto      error // error de dominio esperado sin envolver; nil si se espera envuelto
		envuelto    error // causa esperada con %w; nil si se espera exacto
	}{
		{
			descripcion: "proyecto inexistente prevalece sobre Sprint cerrado sin guardado sembrado",
			preparar: func(t *testing.T, e *entornoDashboard) {
				sembrarSprintDashboard(t, e, "p1", 1, domain.EstadoSprintCerrado, nil)
				e.proyectos.proyectos = map[string]domain.Proyecto{} // proyecto inexistente
			},
			exacto: domain.ErrProyectoNoEncontrado,
		},
		{
			descripcion: "error al verificar el proyecto prevalece sobre Sprint cerrado sin guardado",
			preparar: func(t *testing.T, e *entornoDashboard) {
				sembrarSprintDashboard(t, e, "p1", 1, domain.EstadoSprintCerrado, nil)
				e.proyectos.existeErr = causa
			},
			envuelto: causa,
		},
		{
			descripcion: "error al listar los Sprints prevalece sobre Sprint cerrado sin guardado",
			preparar: func(t *testing.T, e *entornoDashboard) {
				sembrarSprintDashboard(t, e, "p1", 1, domain.EstadoSprintCerrado, nil)
				e.sprints.listarErr = causa
			},
			envuelto: causa,
		},
	}
	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoDashboard()
			sembrarProyectoDashboard(t, entorno, "p1")
			caso.preparar(t, entorno)

			dashboard, err := entorno.servicio().ObtenerDashboard("p1")
			if err == nil {
				t.Fatal("se esperaba un error y se obtuvo nil")
			}
			if caso.envuelto != nil {
				if !errors.Is(err, caso.envuelto) {
					t.Fatalf("errors.Is(err, causa) = false; err = %v", err)
				}
				if err == caso.envuelto {
					t.Errorf("err = %v, se esperaba un error envuelto, no el error tal cual", err)
				}
				if !strings.HasPrefix(err.Error(), "service:") {
					t.Errorf("mensaje = %q, se esperaba el prefijo service:", err.Error())
				}
			} else {
				if err != caso.exacto {
					t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, caso.exacto)
				}
				if entorno.sprints.listarLlamadas != 0 && caso.exacto == domain.ErrProyectoNoEncontrado {
					t.Errorf("ListarPorProyecto (Sprints) fue invocado %d veces; el proyecto debía validarse primero",
						entorno.sprints.listarLlamadas)
				}
			}
			if !reflect.DeepEqual(dashboard, DashboardProyecto{}) {
				t.Errorf("dashboard = %+v, se esperaba DashboardProyecto{} (nunca un éxito parcial)", dashboard)
			}
		})
	}
}
