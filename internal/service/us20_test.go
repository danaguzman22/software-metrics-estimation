package service

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// entornoReporte agrupa los dobles de proyectos, Sprints, backlog, asignaciones,
// integrantes, esfuerzo y defectos para los tests del servicio de reportes
// (US-20). Los servicios concretos comparten las mismas instancias de los
// dobles para actuar sobre el mismo estado.
type entornoReporte struct {
	proyectos    *fakeRepositorioProyectos
	sprints      *fakeRepositorioSprints
	backlog      *fakeRepositorioBacklog
	asignaciones *fakeRepositorioAsignaciones
	integrantes  *fakeRepositorioIntegrantes
	esfuerzo     *fakeRepositorioEsfuerzo
	defectos     *fakeRepositorioMetricasDefectosUS18
}

func nuevoEntornoReporte() *entornoReporte {
	return &entornoReporte{
		proyectos:    nuevoFakeProyectos(),
		sprints:      nuevoFakeSprints(),
		backlog:      nuevoFakeBacklog(),
		asignaciones: nuevoFakeAsignaciones(),
		integrantes:  nuevoFakeIntegrantes(),
		esfuerzo:     nuevoFakeEsfuerzo(),
		defectos:     nuevoFakeMetricasDefectosUS18(),
	}
}

func (e *entornoReporte) backlogServicio() *ServicioBacklog {
	return NuevoServicioBacklog(e.proyectos, e.backlog)
}

func (e *entornoReporte) sprintsServicio() *ServicioSprints {
	return NuevoServicioSprints(e.proyectos, e.sprints, e.backlog, e.asignaciones)
}

func (e *entornoReporte) esfuerzoServicio() *ServicioEsfuerzo {
	return NuevoServicioEsfuerzo(e.proyectos, e.backlog, e.integrantes, e.sprints, e.asignaciones, e.esfuerzo)
}

func (e *entornoReporte) dashboardServicio() *ServicioDashboard {
	return NuevoServicioDashboard(e.proyectos, e.sprintsServicio(), e.esfuerzoServicio())
}

func (e *entornoReporte) defectosServicio() *ServicioDefectos {
	return NuevoServicioDefectos(e.proyectos, e.backlog, e.sprints, e.defectos)
}

func (e *entornoReporte) servicio() *ServicioReportes {
	return NuevoServicioReportes(
		e.proyectos,
		e.backlogServicio(),
		e.sprintsServicio(),
		e.esfuerzoServicio(),
		e.dashboardServicio(),
		e.defectosServicio(),
	)
}

func sembrarProyectoReporte(t *testing.T, e *entornoReporte, proyectoID string) {
	t.Helper()
	e.proyectos.proyectos[proyectoID] = proyectoUS02Servicio(
		"Sistema de Métricas",
		domain.Estado("En curso"),
	)
}

func sembrarSprintReporte(
	t *testing.T,
	e *entornoReporte,
	proyectoID string,
	numero int,
	estado string,
	spPlanificados *int,
) {
	t.Helper()
	sprint := sprintUS07Servicio(numero, estado)
	sprint.SPPlanificados = spPlanificados
	e.sprints.sprints[proyectoID] = append(e.sprints.sprints[proyectoID], sprint)
}

func sembrarHistoriaReporte(
	t *testing.T,
	e *entornoReporte,
	proyectoID, id string,
	storyPoints int,
	estado string,
) {
	t.Helper()
	e.backlog.items[proyectoID] = append(
		e.backlog.items[proyectoID],
		itemUS08Servicio(id, storyPoints, estado),
	)
}

func estimarHorasReporte(
	t *testing.T,
	e *entornoReporte,
	proyectoID, itemID string,
	horas *float64,
) {
	t.Helper()
	for i := range e.backlog.items[proyectoID] {
		if e.backlog.items[proyectoID][i].ID == itemID {
			e.backlog.items[proyectoID][i].HorasEstimadas = horas
			return
		}
	}
	t.Fatalf("no se encontró la historia %q para estimar horas", itemID)
}

func asignarReporte(
	e *entornoReporte,
	proyectoID string,
	sprintNumero, itemNumero int,
) {
	e.asignaciones.sembrar(proyectoID, sprintNumero, itemNumero)
}

func registrarEsfuerzoReporte(
	e *entornoReporte,
	proyectoID string,
	itemNumero int,
	sprintNum *int,
	horas float64,
) {
	e.esfuerzo.registros = append(
		e.esfuerzo.registros,
		registroEsfuerzoFake{
			proyectoID: proyectoID,
			itemNumero: itemNumero,
			sprintNum:  sprintNum,
			registro:   domain.RegistroEsfuerzo{Horas: horas},
		},
	)
}

func sembrarDefectoReporte(e *entornoReporte, proyectoID string, defecto domain.Defecto) {
	e.defectos.defectos[proyectoID] = append(e.defectos.defectos[proyectoID], defecto)
}

func defectoReporte(severidad, estado string, sprintDeteccion int, sprintResolucion *int) domain.Defecto {
	return domain.Defecto{
		Descripcion:      "Defecto sembrado de prueba",
		Severidad:        severidad,
		Estado:           estado,
		ItemID:           "HU-01",
		SprintDeteccion:  sprintDeteccion,
		SprintResolucion: sprintResolucion,
	}
}

// R1 / V1 / CA1 (US-20): un proyecto inexistente devuelve
// domain.ErrProyectoNoEncontrado sin envolver; no se consultan los Sprints
// ni el resto de los dobles.
func TestUS20_Fase1_ServicioReportes_ProyectoInexistente(t *testing.T) {
	entorno := nuevoEntornoReporte()

	reporte, err := entorno.servicio().ReporteProyecto("PRJ-999")
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if err != domain.ErrProyectoNoEncontrado {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrProyectoNoEncontrado)
	}
	if !reflect.DeepEqual(reporte, ReporteProyecto{}) {
		t.Errorf("reporte = %+v, se esperaba ReporteProyecto{}", reporte)
	}
	if entorno.proyectos.existeLlamadas != 1 {
		t.Errorf("Existe fue invocado %d veces, se esperaba 1", entorno.proyectos.existeLlamadas)
	}
	if entorno.sprints.listarLlamadas != 0 {
		t.Errorf("ListarPorProyecto (Sprints) fue invocado %d veces; el proyecto debía validarse primero",
			entorno.sprints.listarLlamadas)
	}
	if entorno.defectos.listarLlamadas != 0 {
		t.Errorf("ListarPorProyecto (defectos) fue invocado %d veces; nunca debe llamarse", entorno.defectos.listarLlamadas)
	}
}

// R6 / V2 / CA2 (US-20): un Sprint inexistente (no aparece en ListarSprints)
// devuelve domain.ErrSprintNoEncontrado sin envolver.
func TestUS20_Fase1_ServicioReportes_SprintInexistente(t *testing.T) {
	entorno := nuevoEntornoReporte()
	sembrarProyectoReporte(t, entorno, "p1")
	sprint1 := 8
	sembrarSprintReporte(t, entorno, "p1", 1, domain.EstadoSprintCerrado, &sprint1)

	reporte, err := entorno.servicio().ReporteSprint("p1", 9)
	if !errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrSprintNoEncontrado)
	}
	if err != domain.ErrSprintNoEncontrado {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrSprintNoEncontrado)
	}
	if !reflect.DeepEqual(reporte, ReporteSprint{}) {
		t.Errorf("reporte = %+v, se esperaba ReporteSprint{}", reporte)
	}
}

// V1 / V2 (US-20): el proyecto se valida antes que cualquier otra cosa en
// ambos reportes; el error de proyecto prevalece sobre cualquier otro dato
// inconsistente sembrado.
func TestUS20_Fase1_ServicioReportes_PrecedenciaValidaciones(t *testing.T) {
	sembrarSinSPPlanificados := func(t *testing.T, e *entornoReporte) {
		t.Helper()
		sembrarSprintReporte(t, e, "p1", 1, domain.EstadoSprintCerrado, nil)
		sembrarHistoriaReporte(t, e, "p1", "HU-01", 8, domain.EstadoCompletada)
	}

	t.Run("proyecto inexistente prevalece en el reporte del proyecto", func(t *testing.T) {
		entorno := nuevoEntornoReporte()
		sembrarSinSPPlanificados(t, entorno)
		entorno.proyectos.proyectos = map[string]domain.Proyecto{} // proyecto inexistente

		_, err := entorno.servicio().ReporteProyecto("p1")
		if err != domain.ErrProyectoNoEncontrado {
			t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrProyectoNoEncontrado)
		}
	})

	t.Run("proyecto inexistente prevalece en el reporte del sprint", func(t *testing.T) {
		entorno := nuevoEntornoReporte()
		sembrarSinSPPlanificados(t, entorno)
		entorno.proyectos.proyectos = map[string]domain.Proyecto{} // proyecto inexistente

		_, err := entorno.servicio().ReporteSprint("p1", 1)
		if err != domain.ErrProyectoNoEncontrado {
			t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrProyectoNoEncontrado)
		}
	})
}

// R2 / R3 / R4 / R5 / CA1 (US-20): el servicio arma el reporte del proyecto
// reutilizando ObtenerDashboard y los servicios existentes: totales copiados
// del dashboard, historias del backlog separadas por estado (orden del
// backlog), horas reales acumuladas por historia, Sprints cerrados con
// velocidad y defectos con las dos fuentes (US-18 para Detectados/Resueltos
// y ListarDefectos para severidad/estado y la lista).
func TestUS20_Fase1_ServicioReportes_ProyectoCompleto(t *testing.T) {
	entorno := nuevoEntornoReporte()
	sembrarProyectoReporte(t, entorno, "p1")
	sprint1 := 8
	sembrarSprintReporte(t, entorno, "p1", 1, domain.EstadoSprintCerrado, &sprint1)
	sprint2 := 13
	sembrarSprintReporte(t, entorno, "p1", 2, domain.EstadoSprintCerrado, &sprint2)
	sembrarSprintReporte(t, entorno, "p1", 3, domain.EstadoSprintActivo, nil)
	sembrarHistoriaReporte(t, entorno, "p1", "HU-01", 5, domain.EstadoCompletada)
	sembrarHistoriaReporte(t, entorno, "p1", "HU-02", 3, domain.EstadoItemPendiente)
	sembrarHistoriaReporte(t, entorno, "p1", "HU-03", 8, domain.EstadoItemPendiente)
	asignarReporte(entorno, "p1", 1, 1)
	asignarReporte(entorno, "p1", 2, 3)
	estimaran := []float64{40, 10, 80}
	estimaran1 := estimaran[0]
	estimaran2 := estimaran[1]
	estimaran3 := estimaran[2]
	estimarHorasReporte(t, entorno, "p1", "HU-01", &estimaran1)
	estimarHorasReporte(t, entorno, "p1", "HU-02", &estimaran2)
	estimarHorasReporte(t, entorno, "p1", "HU-03", &estimaran3)
	// Esfuerzo histórico de HU-01 en todo el proyecto (registros de dos Sprints).
	registrarEsfuerzoReporte(entorno, "p1", 1, intPuntero(1), 50)
	registrarEsfuerzoReporte(entorno, "p1", 1, nil, 12)
	registrarEsfuerzoReporte(entorno, "p1", 3, intPuntero(2), 99)
	// Snapshots de esfuerzo de los Sprints cerrados (US-17).
	entorno.esfuerzo.snapshots["p1:1"] = snapshotEsfuerzoFake{
		disponible: true,
		items: []domain.ItemSprintEsfuerzo{
			{ItemNumero: 1, HorasEstimadas: horasPuntero(40)},
		},
	}
	entorno.esfuerzo.snapshots["p1:2"] = snapshotEsfuerzoFake{
		disponible: true,
		items: []domain.ItemSprintEsfuerzo{
			{ItemNumero: 3, HorasEstimadas: horasPuntero(80)},
		},
	}
	sembrarDefectoReporte(entorno, "p1", defectoReporte(domain.SeveridadDefectoAlta, domain.EstadoDefectoAbierto, 1, nil))
	sembrarDefectoReporte(entorno, "p1", defectoReporte(domain.SeveridadDefectoAlta, domain.EstadoDefectoResuelto, 1, intPuntero(1)))
	sembrarDefectoReporte(entorno, "p1", defectoReporte(domain.SeveridadDefectoBaja, domain.EstadoDefectoEnCurso, 2, nil))
	entorno.defectos.metricasProyecto = domain.MetricasDefectos{Detectados: 121, Resueltos: 117}

	reporte, err := entorno.servicio().ReporteProyecto("p1")
	if err != nil {
		t.Fatalf("ReporteProyecto() devolvió un error inesperado: %v", err)
	}

	// R2: los totales son los de ObtenerDashboard, no un recálculo propio.
	dashboard, err := entorno.dashboardServicio().ObtenerDashboard("p1")
	if err != nil {
		t.Fatalf("ObtenerDashboard() devolvió un error inesperado: %v", err)
	}
	if reporte.SPPlanificadosTotal != dashboard.SPPlanificadosTotal {
		t.Errorf("SPPlanificadosTotal = %d, se esperaba %d (el del dashboard)", reporte.SPPlanificadosTotal, dashboard.SPPlanificadosTotal)
	}
	if reporte.SPCompletadosTotal != dashboard.SPCompletadosTotal {
		t.Errorf("SPCompletadosTotal = %d, se esperaba %d (el del dashboard)", reporte.SPCompletadosTotal, dashboard.SPCompletadosTotal)
	}
	if reporte.VelocidadPromedio != dashboard.VelocidadPromedio {
		t.Errorf("VelocidadPromedio = %v, se esperaba %v (el del dashboard)", reporte.VelocidadPromedio, dashboard.VelocidadPromedio)
	}
	if reporte.PorcentajeHistoriasCompletadas != dashboard.PorcentajeHistoriasCompletadas {
		t.Errorf("PorcentajeHistoriasCompletadas = %v, se esperaba %v (el del dashboard)", reporte.PorcentajeHistoriasCompletadas, dashboard.PorcentajeHistoriasCompletadas)
	}
	if reporte.HorasEstimadasTotal != dashboard.HorasEstimadasTotal {
		t.Errorf("HorasEstimadasTotal = %v, se esperaba %v (el del dashboard)", reporte.HorasEstimadasTotal, dashboard.HorasEstimadasTotal)
	}
	if reporte.HorasRealesTotal != dashboard.HorasRealesTotal {
		t.Errorf("HorasRealesTotal = %v, se esperaba %v (el del dashboard)", reporte.HorasRealesTotal, dashboard.HorasRealesTotal)
	}
	if reporte.EstimacionesCompletas != dashboard.EstimacionesCompletas {
		t.Errorf("EstimacionesCompletas = %v, se esperaba %v (el del dashboard)", reporte.EstimacionesCompletas, dashboard.EstimacionesCompletas)
	}
	if !reflect.DeepEqual(reporte.DesviacionHorasTotal, dashboard.DesviacionHorasTotal) {
		t.Errorf("DesviacionHorasTotal = %v, se esperaba %v (la del dashboard)", reporte.DesviacionHorasTotal, dashboard.DesviacionHorasTotal)
	}
	if reporte.SprintsCerrados != dashboard.SprintsCerrados {
		t.Errorf("SprintsCerrados = %d, se esperaba %d (el del dashboard)", reporte.SprintsCerrados, dashboard.SprintsCerrados)
	}

	// Verificación de los datos sembrados (no del dashboard, que tiene sus propios tests).
	if reporte.SPPlanificadosTotal != 21 {
		t.Errorf("SPPlanificadosTotal = %d, se esperaba 21 (8 + 13 guardados de los cerrados; el Sprint activo no tiene historias asignadas)", reporte.SPPlanificadosTotal)
	}
	if reporte.SPCompletadosTotal != 5 {
		t.Errorf("SPCompletadosTotal = %d, se esperaba 5 (HU-01)", reporte.SPCompletadosTotal)
	}

	// R4: los Sprints cerrados con velocidad, en orden ascendente; el activo no figura.
	if !reflect.DeepEqual(reporte.SprintsCerradosVelocidad, dashboard.VelocidadPorSprint) {
		t.Errorf("SprintsCerradosVelocidad = %+v, se esperaba %+v (la del dashboard)", reporte.SprintsCerradosVelocidad, dashboard.VelocidadPorSprint)
	}
	if len(reporte.SprintsCerradosVelocidad) != 2 {
		t.Fatalf("len(SprintsCerradosVelocidad) = %d, se esperaban 2 (solo cerrados)", len(reporte.SprintsCerradosVelocidad))
	}
	if reporte.SprintsCerradosVelocidad[0] != (PuntoVelocidadSprint{SprintNumero: 1, Velocidad: 5}) {
		t.Errorf("SprintsCerradosVelocidad[0] = %+v, se esperaba {1 5}", reporte.SprintsCerradosVelocidad[0])
	}
	if reporte.SprintsCerradosVelocidad[1] != (PuntoVelocidadSprint{SprintNumero: 2, Velocidad: 0}) {
		t.Errorf("SprintsCerradosVelocidad[1] = %+v, se esperaba {2 0}", reporte.SprintsCerradosVelocidad[1])
	}

	// R3: las historias separadas por estado, en el orden del backlog, con SP,
	// estado, estimación y horas reales acumuladas (suma histórica).
	if len(reporte.HistoriasPlanificadas) != 2 {
		t.Fatalf("len(HistoriasPlanificadas) = %d, se esperaban 2 (HU-02 y HU-03)", len(reporte.HistoriasPlanificadas))
	}
	if reporte.HistoriasPlanificadas[0].ID != "HU-02" || reporte.HistoriasPlanificadas[1].ID != "HU-03" {
		t.Errorf("HistoriasPlanificadas = [%s %s], se esperaba [HU-02 HU-03] (orden del backlog)",
			reporte.HistoriasPlanificadas[0].ID, reporte.HistoriasPlanificadas[1].ID)
	}
	if len(reporte.HistoriasCompletadas) != 1 {
		t.Fatalf("len(HistoriasCompletadas) = %d, se esperaba 1 (HU-01)", len(reporte.HistoriasCompletadas))
	}
	historiaCompletada := reporte.HistoriasCompletadas[0]
	if historiaCompletada.ID != "HU-01" || historiaCompletada.StoryPoints != 5 || historiaCompletada.Estado != domain.EstadoCompletada {
		t.Errorf("HistoriasCompletadas[0] = %+v, se esperaba HU-01, 5 SP, Completada", historiaCompletada)
	}
	if historiaCompletada.HorasEstimadas == nil || *historiaCompletada.HorasEstimadas != 40 {
		t.Errorf("HorasEstimadas = %v, se esperaba 40", historiaCompletada.HorasEstimadas)
	}
	if historiaCompletada.HorasReales != 62 {
		t.Errorf("HorasReales = %v, se esperaba 62 (suma histórica: 50 + 12)", historiaCompletada.HorasReales)
	}
	historiaPendiente := reporte.HistoriasPlanificadas[0]
	if historiaPendiente.StoryPoints != 3 || historiaPendiente.Estado != domain.EstadoItemPendiente {
		t.Errorf("HistoriasPlanificadas[0] = %+v, se esperaba HU-02, 3 SP, Pendiente", historiaPendiente)
	}
	if reporte.HistoriasPlanificadas[1].HorasReales != 99 {
		t.Errorf("HistoriasPlanificadas[1].HorasReales = %v, se esperaba 99 (suma histórica de HU-03)", reporte.HistoriasPlanificadas[1].HorasReales)
	}

	// R5: Detectados/Resueltos reutilizados de US-18 (valores del doble,
	// distintos de los que daría el reconteo); severidad/estado desde la lista.
	if entorno.defectos.metricasProyectoLlamadas != 1 {
		t.Errorf("MetricasDefectosPorProyecto fue invocada %d veces, se esperaba 1", entorno.defectos.metricasProyectoLlamadas)
	}
	if entorno.defectos.listarLlamadas != 1 {
		t.Errorf("ListarDefectos fue invocada %d veces, se esperaba 1", entorno.defectos.listarLlamadas)
	}
	defectos := reporte.Defectos
	if defectos.Detectados != 121 || defectos.Resueltos != 117 {
		t.Errorf("Detectados/Resueltos = %d/%d, se esperaban 121/117 (los de US-18, sin recontar)",
			defectos.Detectados, defectos.Resueltos)
	}
	if defectos.Total != 121 {
		t.Errorf("Total = %d, se esperaba 121 (metricasDefectos.Detectados)", defectos.Total)
	}
	if defectos.Alta != 2 || defectos.Media != 0 || defectos.Baja != 1 {
		t.Errorf("severidad = Alta %d, Media %d, Baja %d; se esperaban 2, 0, 1 (sobre la lista)",
			defectos.Alta, defectos.Media, defectos.Baja)
	}
	if defectos.Abierto != 1 || defectos.EnCurso != 1 || defectos.Resuelto != 1 {
		t.Errorf("estado = Abierto %d, EnCurso %d, Resuelto %d; se esperaban 1, 1, 1 (sobre la lista)",
			defectos.Abierto, defectos.EnCurso, defectos.Resuelto)
	}
	if len(defectos.Defectos) != 3 {
		t.Errorf("len(Defectos) = %d, se esperaban 3 (la lista completa de ListarDefectos)", len(defectos.Defectos))
	}
	if len(defectos.DefectosResueltos) != 0 {
		t.Errorf("len(DefectosResueltos) = %d, se esperaban 0 (el reporte del proyecto no la completa)", len(defectos.DefectosResueltos))
	}
}

// R5 / CA1 (US-20): los conteos por severidad y por estado se arman sobre
// ListarDefectos y cubren todos los defectos del proyecto; Detectados y
// Resueltos se toman del doble de las métricas de US-18 sin recontar (los
// valores del doble difieren deliberadamente del reconteo).
func TestUS20_Fase1_ServicioReportes_DefectosPorSeveridadYEstado(t *testing.T) {
	entorno := nuevoEntornoReporte()
	sembrarProyectoReporte(t, entorno, "p1")
	sembrarDefectoReporte(entorno, "p1", defectoReporte(domain.SeveridadDefectoAlta, domain.EstadoDefectoAbierto, 1, nil))
	sembrarDefectoReporte(entorno, "p1", defectoReporte(domain.SeveridadDefectoAlta, domain.EstadoDefectoEnCurso, 1, nil))
	sembrarDefectoReporte(entorno, "p1", defectoReporte(domain.SeveridadDefectoMedia, domain.EstadoDefectoResuelto, 1, intPuntero(1)))
	sembrarDefectoReporte(entorno, "p1", defectoReporte(domain.SeveridadDefectoMedia, domain.EstadoDefectoAbierto, 2, nil))
	sembrarDefectoReporte(entorno, "p1", defectoReporte(domain.SeveridadDefectoBaja, domain.EstadoDefectoResuelto, 2, intPuntero(2)))
	entorno.defectos.metricasProyecto = domain.MetricasDefectos{Detectados: 42, Resueltos: 17}

	reporte, err := entorno.servicio().ReporteProyecto("p1")
	if err != nil {
		t.Fatalf("ReporteProyecto() devolvió un error inesperado: %v", err)
	}

	defectos := reporte.Defectos
	if defectos.Alta != 2 || defectos.Media != 2 || defectos.Baja != 1 {
		t.Errorf("severidad = Alta %d, Media %d, Baja %d; se esperaban 2, 2, 1",
			defectos.Alta, defectos.Media, defectos.Baja)
	}
	if defectos.Abierto != 2 || defectos.EnCurso != 1 || defectos.Resuelto != 2 {
		t.Errorf("estado = Abierto %d, EnCurso %d, Resuelto %d; se esperaban 2, 1, 2",
			defectos.Abierto, defectos.EnCurso, defectos.Resuelto)
	}
	if defectos.Detectados != 42 || defectos.Resueltos != 17 {
		t.Errorf("Detectados/Resueltos = %d/%d, se esperaban 42/17 (los de US-18, sin recontar: la lista daría 5/2)",
			defectos.Detectados, defectos.Resueltos)
	}
	if len(defectos.Defectos) != 5 {
		t.Errorf("len(Defectos) = %d, se esperaban 5 (la lista completa de ListarDefectos)", len(defectos.Defectos))
	}
}

// R5 / CA1 (US-20): un proyecto sin defectos devuelve todos los conteos en
// cero y la lista vacía, sin error (con las métricas de US-18 en cero).
func TestUS20_Fase1_ServicioReportes_SinDefectos(t *testing.T) {
	entorno := nuevoEntornoReporte()
	sembrarProyectoReporte(t, entorno, "p1")

	reporte, err := entorno.servicio().ReporteProyecto("p1")
	if err != nil {
		t.Fatalf("ReporteProyecto() devolvió un error inesperado: %v", err)
	}

	defectos := reporte.Defectos
	if !reflect.DeepEqual(defectos, ResumenDefectos{}) {
		t.Errorf("Defectos = %+v, se esperaba ResumenDefectos{} (todos en cero, listas vacías)", defectos)
	}
}

// R3 / CA1 (US-20): un proyecto sin historias devuelve las dos listas
// vacías y los totales en cero, sin errores ni división por cero.
func TestUS20_Fase1_ServicioReportes_SinHistorias(t *testing.T) {
	entorno := nuevoEntornoReporte()
	sembrarProyectoReporte(t, entorno, "p1")

	reporte, err := entorno.servicio().ReporteProyecto("p1")
	if err != nil {
		t.Fatalf("ReporteProyecto() devolvió un error inesperado: %v", err)
	}

	if len(reporte.HistoriasPlanificadas) != 0 {
		t.Errorf("len(HistoriasPlanificadas) = %d, se esperaba 0", len(reporte.HistoriasPlanificadas))
	}
	if len(reporte.HistoriasCompletadas) != 0 {
		t.Errorf("len(HistoriasCompletadas) = %d, se esperaba 0", len(reporte.HistoriasCompletadas))
	}
	if reporte.SPPlanificadosTotal != 0 || reporte.SPCompletadosTotal != 0 {
		t.Errorf("SP = %d/%d, se esperaban 0/0", reporte.SPPlanificadosTotal, reporte.SPCompletadosTotal)
	}
	if reporte.VelocidadPromedio != 0 || reporte.PorcentajeHistoriasCompletadas != 0 {
		t.Errorf("promedio/porcentaje = %v/%v, se esperaban 0/0 (sin división por cero)",
			reporte.VelocidadPromedio, reporte.PorcentajeHistoriasCompletadas)
	}
	if !reporte.EstimacionesCompletas {
		t.Errorf("EstimacionesCompletas = false, se esperaba true (sin historias no hay pendientes de estimar)")
	}
	if reporte.DesviacionHorasTotal != nil {
		t.Errorf("DesviacionHorasTotal = %v, se esperaba nil con la suma de estimaciones en cero",
			*reporte.DesviacionHorasTotal)
	}
}

// R3 / D13 (US-20): una historia sin estimar en horas llega al reporte con
// HorasEstimadas = nil ("sin estimar").
func TestUS20_Fase1_ServicioReportes_HistoriaSinEstimarEnHoras(t *testing.T) {
	entorno := nuevoEntornoReporte()
	sembrarProyectoReporte(t, entorno, "p1")
	sembrarHistoriaReporte(t, entorno, "p1", "HU-01", 5, domain.EstadoItemPendiente)

	reporte, err := entorno.servicio().ReporteProyecto("p1")
	if err != nil {
		t.Fatalf("ReporteProyecto() devolvió un error inesperado: %v", err)
	}
	if len(reporte.HistoriasPlanificadas) != 1 {
		t.Fatalf("len(HistoriasPlanificadas) = %d, se esperaba 1", len(reporte.HistoriasPlanificadas))
	}
	if reporte.HistoriasPlanificadas[0].HorasEstimadas != nil {
		t.Errorf("HorasEstimadas = %v, se esperaba nil (historia sin estimar en horas)",
			*reporte.HistoriasPlanificadas[0].HorasEstimadas)
	}
	if reporte.EstimacionesCompletas {
		t.Errorf("EstimacionesCompletas = true, se esperaba false (hay una historia sin estimar)")
	}
	if reporte.DesviacionHorasTotal != nil {
		t.Errorf("DesviacionHorasTotal = %v, se esperaba nil con estimaciones incompletas",
			*reporte.DesviacionHorasTotal)
	}
}

// R4 / CA1 (US-20): un proyecto sin Sprints cerrados devuelve la lista de
// velocidades vacía y la velocidad promedio en cero.
func TestUS20_Fase1_ServicioReportes_SinSprintsCerrados(t *testing.T) {
	entorno := nuevoEntornoReporte()
	sembrarProyectoReporte(t, entorno, "p1")
	sembrarHistoriaReporte(t, entorno, "p1", "HU-01", 8, domain.EstadoCompletada)

	reporte, err := entorno.servicio().ReporteProyecto("p1")
	if err != nil {
		t.Fatalf("ReporteProyecto() devolvió un error inesperado: %v", err)
	}
	if reporte.SprintsCerrados != 0 {
		t.Errorf("SprintsCerrados = %d, se esperaba 0", reporte.SprintsCerrados)
	}
	if len(reporte.SprintsCerradosVelocidad) != 0 {
		t.Errorf("SprintsCerradosVelocidad = %+v, se esperaba vacía (sin Sprints cerrados)", reporte.SprintsCerradosVelocidad)
	}
	if reporte.VelocidadPromedio != 0 {
		t.Errorf("VelocidadPromedio = %v, se esperaba 0 (sin Sprints cerrados)", reporte.VelocidadPromedio)
	}
}

// R4 / CA1 (US-20): con Sprints solo activos ninguno figura en la lista de
// velocidades; los totales igual suman los datos del activo y la velocidad
// promedio es cero.
func TestUS20_Fase1_ServicioReportes_SoloActivos(t *testing.T) {
	entorno := nuevoEntornoReporte()
	sembrarProyectoReporte(t, entorno, "p1")
	sembrarSprintReporte(t, entorno, "p1", 1, domain.EstadoSprintActivo, nil)
	sembrarHistoriaReporte(t, entorno, "p1", "HU-01", 8, domain.EstadoItemPendiente)
	estimaran := 4.0
	estimarHorasReporte(t, entorno, "p1", "HU-01", &estimaran)
	asignarReporte(entorno, "p1", 1, 1)
	registrarEsfuerzoReporte(entorno, "p1", 1, intPuntero(1), 6)

	reporte, err := entorno.servicio().ReporteProyecto("p1")
	if err != nil {
		t.Fatalf("ReporteProyecto() devolvió un error inesperado: %v", err)
	}
	if reporte.SprintsCerrados != 0 {
		t.Errorf("SprintsCerrados = %d, se esperaba 0 (solo activos)", reporte.SprintsCerrados)
	}
	if len(reporte.SprintsCerradosVelocidad) != 0 {
		t.Errorf("SprintsCerradosVelocidad = %+v, se esperaba vacía (solo activos)", reporte.SprintsCerradosVelocidad)
	}
	if reporte.VelocidadPromedio != 0 {
		t.Errorf("VelocidadPromedio = %v, se esperaba 0 (sin Sprints cerrados)", reporte.VelocidadPromedio)
	}
	if reporte.SPPlanificadosTotal != 8 {
		t.Errorf("SPPlanificadosTotal = %d, se esperaba 8 (el activo también suma)", reporte.SPPlanificadosTotal)
	}
}

// R7 / R8 / R9 / CA2 (US-20): el servicio arma el reporte del Sprint: objetivo,
// fechas, estado, historias asignadas (orden por número) con SP, estado y
// horas acumuladas, totales de SP, horas y desviación desde MetricasSprint,
// y defectos con las dos fuentes (US-18 para Detectados/Resueltos).
func TestUS20_Fase1_ServicioReportes_SprintCompleto(t *testing.T) {
	entorno := nuevoEntornoReporte()
	sembrarProyectoReporte(t, entorno, "p1")
	sprint1 := 13
	sembrarSprintReporte(t, entorno, "p1", 1, domain.EstadoSprintCerrado, &sprint1)
	sembrarSprintReporte(t, entorno, "p1", 2, domain.EstadoSprintActivo, nil)
	sembrarHistoriaReporte(t, entorno, "p1", "HU-01", 5, domain.EstadoCompletada)
	sembrarHistoriaReporte(t, entorno, "p1", "HU-02", 3, domain.EstadoItemPendiente)
	sembrarHistoriaReporte(t, entorno, "p1", "HU-03", 8, domain.EstadoItemPendiente)
	// Se asignan fuera de orden a propósito: el reporte las devuelve ordenadas
	// por número de historia (US-22 D24).
	asignarReporte(entorno, "p1", 1, 2)
	asignarReporte(entorno, "p1", 1, 1)
	estimaran := []float64{40, 10}
	estimaran1 := estimaran[0]
	estimaran2 := estimaran[1]
	estimarHorasReporte(t, entorno, "p1", "HU-01", &estimaran1)
	estimarHorasReporte(t, entorno, "p1", "HU-02", &estimaran2)
	// Snapshot del Sprint 1 (US-17): 50 estimadas.
	entorno.esfuerzo.snapshots["p1:1"] = snapshotEsfuerzoFake{
		disponible: true,
		items: []domain.ItemSprintEsfuerzo{
			{ItemNumero: 1, HorasEstimadas: horasPuntero(40)},
			{ItemNumero: 2, HorasEstimadas: horasPuntero(10)},
		},
	}
	// Esfuerzo real del Sprint 1 y registros históricos atribuidos al Sprint 2.
	registrarEsfuerzoReporte(entorno, "p1", 1, intPuntero(1), 35)
	registrarEsfuerzoReporte(entorno, "p1", 2, intPuntero(1), 25)
	registrarEsfuerzoReporte(entorno, "p1", 1, intPuntero(2), 5)
	// Defectos: detectados en el Sprint 1, resueltos en el Sprint 1 y en el
	// Sprint 2, y uno del Sprint 2 que no debe figurar en el reporte del Sprint 1.
	sembrarDefectoReporte(entorno, "p1", defectoReporte(domain.SeveridadDefectoAlta, domain.EstadoDefectoAbierto, 1, nil))
	sembrarDefectoReporte(entorno, "p1", defectoReporte(domain.SeveridadDefectoAlta, domain.EstadoDefectoResuelto, 1, intPuntero(1)))
	sembrarDefectoReporte(entorno, "p1", defectoReporte(domain.SeveridadDefectoMedia, domain.EstadoDefectoResuelto, 1, intPuntero(2)))
	sembrarDefectoReporte(entorno, "p1", defectoReporte(domain.SeveridadDefectoBaja, domain.EstadoDefectoAbierto, 2, nil))
	entorno.defectos.metricasSprint = domain.MetricasDefectos{Detectados: 122, Resueltos: 118}

	reporte, err := entorno.servicio().ReporteSprint("p1", 1)
	if err != nil {
		t.Fatalf("ReporteSprint() devolvió un error inesperado: %v", err)
	}

	if reporte.Numero != 1 || reporte.Objetivo != "Planificar la primera iteración" {
		t.Errorf("Numero/Objetivo = %d/%q, se esperaban 1/Planificar la primera iteración", reporte.Numero, reporte.Objetivo)
	}
	if reporte.FechaInicio.Format("2006-01-02") != "2026-03-02" || reporte.FechaFin.Format("2006-01-02") != "2026-03-20" {
		t.Errorf("fechas = %s..%s, se esperaban 2026-03-02..2026-03-20",
			reporte.FechaInicio.Format("2006-01-02"), reporte.FechaFin.Format("2006-01-02"))
	}
	if reporte.Estado != domain.EstadoSprintCerrado {
		t.Errorf("Estado = %q, se esperaba %q", reporte.Estado, domain.EstadoSprintCerrado)
	}

	// R7: las historias asignadas, ordenadas por número de historia.
	if len(reporte.HistoriasAsignadas) != 2 {
		t.Fatalf("len(HistoriasAsignadas) = %d, se esperaban 2", len(reporte.HistoriasAsignadas))
	}
	if reporte.HistoriasAsignadas[0].ID != "HU-01" || reporte.HistoriasAsignadas[1].ID != "HU-02" {
		t.Errorf("HistoriasAsignadas = [%s %s], se esperaba [HU-01 HU-02] (orden por número de historia)",
			reporte.HistoriasAsignadas[0].ID, reporte.HistoriasAsignadas[1].ID)
	}
	historia1 := reporte.HistoriasAsignadas[0]
	if historia1.StoryPoints != 5 || historia1.Estado != domain.EstadoCompletada {
		t.Errorf("HistoriasAsignadas[0] = %+v, se esperaba HU-01, 5 SP, Completada", historia1)
	}
	if historia1.HorasEstimadas == nil || *historia1.HorasEstimadas != 40 {
		t.Errorf("HorasEstimadas = %v, se esperaba 40", historia1.HorasEstimadas)
	}
	if historia1.HorasReales != 40 {
		t.Errorf("HorasReales = %v, se esperaba 40 (histórico acumulado: 35 + 5)", historia1.HorasReales)
	}
	historia2 := reporte.HistoriasAsignadas[1]
	if historia2.StoryPoints != 3 || historia2.Estado != domain.EstadoItemPendiente {
		t.Errorf("HistoriasAsignadas[1] = %+v, se esperaba HU-02, 3 SP, Pendiente", historia2)
	}
	if historia2.HorasReales != 25 {
		t.Errorf("HorasReales = %v, se esperaba 25 (histórico acumulado de HU-02)", historia2.HorasReales)
	}

	// R7: SP planificados, completados y velocidad tal como los devuelve el servicio.
	if reporte.SPPlanificados != 13 {
		t.Errorf("SPPlanificados = %d, se esperaba 13 (guardados al cerrar)", reporte.SPPlanificados)
	}
	if reporte.SPCompletados != 5 || reporte.Velocidad != 5 {
		t.Errorf("SPCompletados/Velocidad = %d/%d, se esperaban 5/5", reporte.SPCompletados, reporte.Velocidad)
	}

	// R8: horas estimadas vs reales del Sprint copiadas de MetricasSprint.
	if !reporte.MetricasDisponibles {
		t.Errorf("MetricasDisponibles = false, se esperaba true (Sprint cerrado con snapshot)")
	}
	if reporte.HorasEstimadas != 50 || reporte.HorasReales != 60 {
		t.Errorf("horas = %v/%v, se esperaban 50/60", reporte.HorasEstimadas, reporte.HorasReales)
	}
	if !reporte.EstimacionesCompletas {
		t.Errorf("EstimacionesCompletas = false, se esperaba true")
	}
	if reporte.DesviacionHoras == nil {
		t.Fatal("DesviacionHoras = nil, se esperaba 20 (60 reales sobre 50 estimadas)")
	} else if *reporte.DesviacionHoras != 20 {
		t.Errorf("DesviacionHoras = %v, se esperaba 20", *reporte.DesviacionHoras)
	}

	// R9: Detectados/Resueltos reutilizados de US-18; severidad/estado y
	// listas sobre los detectados en el Sprint; el de otro Sprint no figura.
	if entorno.defectos.metricasSprintLlamadas != 1 {
		t.Errorf("MetricasDefectosPorSprint fue invocada %d veces, se esperaba 1", entorno.defectos.metricasSprintLlamadas)
	}
	if entorno.defectos.listarLlamadas != 1 {
		t.Errorf("ListarDefectos fue invocada %d veces, se esperaba 1", entorno.defectos.listarLlamadas)
	}
	defectos := reporte.Defectos
	if defectos.Detectados != 122 || defectos.Resueltos != 118 {
		t.Errorf("Detectados/Resueltos = %d/%d, se esperaban 122/118 (los de US-18, sin recontar)",
			defectos.Detectados, defectos.Resueltos)
	}
	if defectos.Alta != 2 || defectos.Media != 1 || defectos.Baja != 0 {
		t.Errorf("severidad = Alta %d, Media %d, Baja %d; se esperaban 2, 1, 0 (detectados en el Sprint)",
			defectos.Alta, defectos.Media, defectos.Baja)
	}
	if defectos.Abierto != 1 || defectos.EnCurso != 0 || defectos.Resuelto != 2 {
		t.Errorf("estado = Abierto %d, EnCurso %d, Resuelto %d; se esperaban 1, 0, 2 (detectados en el Sprint)",
			defectos.Abierto, defectos.EnCurso, defectos.Resuelto)
	}
	if defectos.Total != 3 {
		t.Errorf("Total = %d, se esperaba 3 (detectados en el Sprint)", defectos.Total)
	}
	if len(defectos.Defectos) != 3 {
		t.Fatalf("len(Defectos) = %d, se esperaban 3 (los detectados en el Sprint)", len(defectos.Defectos))
	}
	for _, defecto := range defectos.Defectos {
		if defecto.SprintDeteccion != 1 {
			t.Errorf("Defectos incluye un defecto detectado en el Sprint %d; solo deben figurar los del Sprint 1",
				defecto.SprintDeteccion)
		}
	}
	if len(defectos.DefectosResueltos) != 1 {
		t.Fatalf("len(DefectosResueltos) = %d, se esperaba 1 (resueltos en el Sprint 1)", len(defectos.DefectosResueltos))
	}
	if defectos.DefectosResueltos[0].SprintResolucion == nil || *defectos.DefectosResueltos[0].SprintResolucion != 1 {
		t.Errorf("DefectosResueltos[0].SprintResolucion = %v, se esperaba 1", defectos.DefectosResueltos[0].SprintResolucion)
	}
}

// R7 / CA2 (US-20): un Sprint sin historias devuelve la lista vacía y los
// totales en cero, sin error.
func TestUS20_Fase1_ServicioReportes_SprintSinHistorias(t *testing.T) {
	entorno := nuevoEntornoReporte()
	sembrarProyectoReporte(t, entorno, "p1")
	sembrarSprintReporte(t, entorno, "p1", 1, domain.EstadoSprintActivo, nil)

	reporte, err := entorno.servicio().ReporteSprint("p1", 1)
	if err != nil {
		t.Fatalf("ReporteSprint() devolvió un error inesperado: %v", err)
	}
	if len(reporte.HistoriasAsignadas) != 0 {
		t.Errorf("len(HistoriasAsignadas) = %d, se esperaba 0", len(reporte.HistoriasAsignadas))
	}
	if reporte.SPPlanificados != 0 || reporte.SPCompletados != 0 || reporte.Velocidad != 0 {
		t.Errorf("SP planificados/completados/velocidad = %d/%d/%d, se esperaban 0/0/0",
			reporte.SPPlanificados, reporte.SPCompletados, reporte.Velocidad)
	}
	if reporte.Numero != 1 || reporte.Estado != domain.EstadoSprintActivo {
		t.Errorf("Numero/Estado = %d/%q, se esperaban 1/%q", reporte.Numero, reporte.Estado, domain.EstadoSprintActivo)
	}
}

// R9 / CA2 (US-20): un Sprint sin defectos devuelve todos los conteos en
// cero y las listas vacías, sin error.
func TestUS20_Fase1_ServicioReportes_SprintSinDefectos(t *testing.T) {
	entorno := nuevoEntornoReporte()
	sembrarProyectoReporte(t, entorno, "p1")
	sprint1 := 8
	sembrarSprintReporte(t, entorno, "p1", 1, domain.EstadoSprintCerrado, &sprint1)

	reporte, err := entorno.servicio().ReporteSprint("p1", 1)
	if err != nil {
		t.Fatalf("ReporteSprint() devolvió un error inesperado: %v", err)
	}
	if !reflect.DeepEqual(reporte.Defectos, ResumenDefectos{}) {
		t.Errorf("Defectos = %+v, se esperaba ResumenDefectos{} (todos en cero, listas vacías)", reporte.Defectos)
	}
}

// R8 / D16 / CA2 (US-20): un Sprint cerrado sin snapshot de esfuerzo
// (domain.ErrMetricasSprintNoDisponibles) no es un error: el reporte se
// arma con MetricasDisponibles = false, horas en cero y desviación nil;
// el resto del reporte funciona con normalidad.
func TestUS20_Fase1_ServicioReportes_SprintCerradoSinSnapshot(t *testing.T) {
	entorno := nuevoEntornoReporte()
	sembrarProyectoReporte(t, entorno, "p1")
	sprint1 := 8
	sembrarSprintReporte(t, entorno, "p1", 1, domain.EstadoSprintCerrado, &sprint1)
	sembrarHistoriaReporte(t, entorno, "p1", "HU-01", 5, domain.EstadoCompletada)
	estimaran := 10.0
	estimarHorasReporte(t, entorno, "p1", "HU-01", &estimaran)
	asignarReporte(entorno, "p1", 1, 1)

	reporte, err := entorno.servicio().ReporteSprint("p1", 1)
	if err != nil {
		t.Fatalf("ReporteSprint() devolvió un error inesperado: %v (el Sprint sin snapshot no debe fallar)", err)
	}
	if reporte.MetricasDisponibles {
		t.Errorf("MetricasDisponibles = true, se esperaba false (cerrado sin snapshot)")
	}
	if reporte.HorasEstimadas != 0 || reporte.HorasReales != 0 {
		t.Errorf("horas = %v/%v, se esperaban 0/0 sin snapshot", reporte.HorasEstimadas, reporte.HorasReales)
	}
	if reporte.DesviacionHoras != nil {
		t.Errorf("DesviacionHoras = %v, se esperaba nil sin snapshot", *reporte.DesviacionHoras)
	}
	// El resto del reporte funciona con normalidad.
	if len(reporte.HistoriasAsignadas) != 1 {
		t.Errorf("len(HistoriasAsignadas) = %d, se esperaba 1", len(reporte.HistoriasAsignadas))
	}
	if reporte.SPPlanificados != 8 || reporte.SPCompletados != 5 {
		t.Errorf("SP planificados/completados = %d/%d, se esperaban 8/5", reporte.SPPlanificados, reporte.SPCompletados)
	}
	if reporte.Velocidad != 5 {
		t.Errorf("Velocidad = %d, se esperaba 5", reporte.Velocidad)
	}
}

// R8 (US-20): con estimaciones incompletas en el snapshot la desviación no
// está disponible; las horas estimadas suman solo las disponibles.
func TestUS20_Fase1_ServicioReportes_SprintEstimacionesIncompletas(t *testing.T) {
	entorno := nuevoEntornoReporte()
	sembrarProyectoReporte(t, entorno, "p1")
	sprint1 := 8
	sembrarSprintReporte(t, entorno, "p1", 1, domain.EstadoSprintCerrado, &sprint1)
	sembrarHistoriaReporte(t, entorno, "p1", "HU-01", 3, domain.EstadoItemPendiente)
	sembrarHistoriaReporte(t, entorno, "p1", "HU-02", 5, domain.EstadoItemPendiente)
	asignarReporte(entorno, "p1", 1, 1)
	asignarReporte(entorno, "p1", 1, 2)
	estimaran1 := 40.0
	estimarHorasReporte(t, entorno, "p1", "HU-01", &estimaran1)
	// HU-02 queda sin estimación ni en el backlog ni en el snapshot.
	entorno.esfuerzo.snapshots["p1:1"] = snapshotEsfuerzoFake{
		disponible: true,
		items: []domain.ItemSprintEsfuerzo{
			{ItemNumero: 1, HorasEstimadas: horasPuntero(40)},
			{ItemNumero: 2, HorasEstimadas: nil},
		},
	}
	registrarEsfuerzoReporte(entorno, "p1", 1, intPuntero(1), 30)

	reporte, err := entorno.servicio().ReporteSprint("p1", 1)
	if err != nil {
		t.Fatalf("ReporteSprint() devolvió un error inesperado: %v", err)
	}
	if reporte.EstimacionesCompletas {
		t.Errorf("EstimacionesCompletas = true, se esperaba false (HU-02 sin estimar)")
	}
	if reporte.HorasEstimadas != 40 {
		t.Errorf("HorasEstimadas = %v, se esperaba 40 (solo las disponibles)", reporte.HorasEstimadas)
	}
	if reporte.HorasReales != 30 {
		t.Errorf("HorasReales = %v, se esperaba 30", reporte.HorasReales)
	}
	if reporte.DesviacionHoras != nil {
		t.Errorf("DesviacionHoras = %v, se esperaba nil con estimaciones incompletas", *reporte.DesviacionHoras)
	}
}

// R8 (US-20): con la suma de estimaciones del Sprint en cero la desviación
// no está disponible, aunque las estimaciones estén completas.
func TestUS20_Fase1_ServicioReportes_SprintEstimacionesSumaCero(t *testing.T) {
	entorno := nuevoEntornoReporte()
	sembrarProyectoReporte(t, entorno, "p1")
	sprint1 := 8
	sembrarSprintReporte(t, entorno, "p1", 1, domain.EstadoSprintCerrado, &sprint1)
	sembrarHistoriaReporte(t, entorno, "p1", "HU-01", 3, domain.EstadoItemPendiente)
	asignarReporte(entorno, "p1", 1, 1)
	estimaran := 0.0
	estimarHorasReporte(t, entorno, "p1", "HU-01", &estimaran)
	entorno.esfuerzo.snapshots["p1:1"] = snapshotEsfuerzoFake{
		disponible: true,
		items: []domain.ItemSprintEsfuerzo{
			{ItemNumero: 1, HorasEstimadas: horasPuntero(0)},
		},
	}
	registrarEsfuerzoReporte(entorno, "p1", 1, intPuntero(1), 5)

	reporte, err := entorno.servicio().ReporteSprint("p1", 1)
	if err != nil {
		t.Fatalf("ReporteSprint() devolvió un error inesperado: %v", err)
	}
	if !reporte.EstimacionesCompletas {
		t.Errorf("EstimacionesCompletas = false, se esperaba true (todas con estimación explícita)")
	}
	if reporte.HorasEstimadas != 0 {
		t.Errorf("HorasEstimadas = %v, se esperaba 0", reporte.HorasEstimadas)
	}
	if reporte.DesviacionHoras != nil {
		t.Errorf("DesviacionHoras = %v, se esperaba nil con la suma de estimaciones en cero", *reporte.DesviacionHoras)
	}
}

// R8 / E2 (US-20): un Sprint cerrado sin SPPlanificados guardado (dato
// inconsistente de US-09/US-16) produce un error no de dominio que se
// propaga; no se mapea a un error de dominio.
func TestUS20_Fase1_ServicioReportes_SprintSPPlanificadosFaltante(t *testing.T) {
	entorno := nuevoEntornoReporte()
	sembrarProyectoReporte(t, entorno, "p1")
	sembrarSprintReporte(t, entorno, "p1", 1, domain.EstadoSprintCerrado, nil)

	reporte, err := entorno.servicio().ReporteSprint("p1", 1)
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
	if !reflect.DeepEqual(reporte, ReporteSprint{}) {
		t.Errorf("reporte = %+v, se esperaba ReporteSprint{} (nunca un éxito parcial)", reporte)
	}
}

// R2 / E2 (US-20): la misma inconsistencia en un Sprint cerrado hace fallar
// el reporte del proyecto (el dashboard propaga el error); no se mapea a
// un error de dominio.
func TestUS20_Fase1_ServicioReportes_ProyectoSPPlanificadosFaltante(t *testing.T) {
	entorno := nuevoEntornoReporte()
	sembrarProyectoReporte(t, entorno, "p1")
	sembrarSprintReporte(t, entorno, "p1", 1, domain.EstadoSprintCerrado, nil)

	reporte, err := entorno.servicio().ReporteProyecto("p1")
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
	if !reflect.DeepEqual(reporte, ReporteProyecto{}) {
		t.Errorf("reporte = %+v, se esperaba ReporteProyecto{} (nunca un éxito parcial)", reporte)
	}
}

// R1 / R2 / R5 / V4 / CA1 (US-20): los errores de infraestructura del reporte
// del proyecto se propagan envueltos con %w (prefijo "service:"), conservando
// la causa con errors.Is; el reporte no los mapea a errores de dominio ni
// devuelve éxitos parciales.
func TestUS20_Fase1_ServicioReportes_ErrorDeRepositorioEnvuelto(t *testing.T) {
	causa := errors.New("fallo simulado del repositorio")
	cerrarSprint := func(t *testing.T, e *entornoReporte, numero int) {
		t.Helper()
		sp := 8
		sembrarSprintReporte(t, e, "p1", numero, domain.EstadoSprintCerrado, &sp)
	}
	casos := []struct {
		descripcion string
		preparar    func(t *testing.T, e *entornoReporte)
		contexto    string
	}{
		{
			descripcion: "error al verificar el proyecto",
			preparar:    func(t *testing.T, e *entornoReporte) { e.proyectos.existeErr = causa },
			contexto:    "verificar proyecto",
		},
		{
			descripcion: "error al listar los Sprints del proyecto",
			preparar:    func(t *testing.T, e *entornoReporte) { e.sprints.listarErr = causa },
			contexto:    "listar sprints del proyecto",
		},
		{
			descripcion: "error al listar las historias del backlog",
			preparar:    func(t *testing.T, e *entornoReporte) { cerrarSprint(t, e, 1); e.backlog.listarErr = causa },
			contexto:    "listar ítems del backlog",
		},
		{
			descripcion: "error al sumar horas reales del proyecto",
			preparar:    func(t *testing.T, e *entornoReporte) { cerrarSprint(t, e, 1); e.esfuerzo.horasProyectoErr = causa },
			contexto:    "sumar horas reales del proyecto",
		},
		{
			descripcion: "error al sumar horas reales de una historia",
			preparar: func(t *testing.T, e *entornoReporte) {
				cerrarSprint(t, e, 1)
				sembrarHistoriaReporte(t, e, "p1", "HU-01", 5, domain.EstadoCompletada)
				e.esfuerzo.horasItemErr = causa
			},
			contexto: "sumar horas reales de la historia",
		},
		{
			descripcion: "error al consultar las métricas de defectos del proyecto",
			preparar:    func(t *testing.T, e *entornoReporte) { cerrarSprint(t, e, 1); e.defectos.metricasProyectoErr = causa },
			contexto:    "consultar métricas de defectos del proyecto",
		},
		{
			descripcion: "error al listar los defectos del proyecto",
			preparar:    func(t *testing.T, e *entornoReporte) { cerrarSprint(t, e, 1); e.defectos.listarErr = causa },
			contexto:    "listar defectos del proyecto",
		},
	}
	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoReporte()
			sembrarProyectoReporte(t, entorno, "p1")
			caso.preparar(t, entorno)

			reporte, err := entorno.servicio().ReporteProyecto("p1")
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
			if !reflect.DeepEqual(reporte, ReporteProyecto{}) {
				t.Errorf("reporte = %+v, se esperaba ReporteProyecto{} (nunca un éxito parcial)", reporte)
			}
		})
	}
}

// R6 / R7 / R8 / R9 / V4 / CA2 (US-20): los errores de infraestructura del
// reporte del Sprint se propagan envueltos con %w, conservando la causa con
// errors.Is; el reporte no los mapea a errores de dominio ni devuelve
// éxitos parciales.
func TestUS20_Fase1_ServicioReportes_ErrorDeRepositorioEnvueltoEnSprint(t *testing.T) {
	causa := errors.New("fallo simulado del repositorio")
	cerrarSprint := func(t *testing.T, e *entornoReporte, numero int) {
		t.Helper()
		sp := 8
		sembrarSprintReporte(t, e, "p1", numero, domain.EstadoSprintCerrado, &sp)
	}
	entornoConSnapshot := func(t *testing.T, e *entornoReporte, numero int) {
		t.Helper()
		e.esfuerzo.snapshots["p1:1"] = snapshotEsfuerzoFake{disponible: true}
	}
	casos := []struct {
		descripcion string
		preparar    func(t *testing.T, e *entornoReporte)
		contexto    string
	}{
		{
			descripcion: "error al listar los Sprints del proyecto",
			preparar:    func(t *testing.T, e *entornoReporte) { e.sprints.listarErr = causa },
			contexto:    "listar sprints del proyecto",
		},
		{
			descripcion: "error al listar las historias asignadas al Sprint",
			preparar: func(t *testing.T, e *entornoReporte) {
				cerrarSprint(t, e, 1)
				e.asignaciones.listarErr = causa
			},
			contexto: "listar ítems del sprint",
		},
		{
			descripcion: "error al obtener el Sprint al sumar los planificados",
			preparar: func(t *testing.T, e *entornoReporte) {
				cerrarSprint(t, e, 1)
				e.sprints.obtenerErr = causa
			},
			contexto: "verificar sprint",
		},
		{
			descripcion: "error al obtener el snapshot del Sprint",
			preparar: func(t *testing.T, e *entornoReporte) {
				cerrarSprint(t, e, 1)
				e.esfuerzo.snapshots["p1:1"] = snapshotEsfuerzoFake{disponible: true, err: causa}
			},
			contexto: "obtener snapshot del sprint",
		},
		{
			descripcion: "error al sumar horas reales del Sprint",
			preparar: func(t *testing.T, e *entornoReporte) {
				cerrarSprint(t, e, 1)
				entornoConSnapshot(t, e, 1)
				e.esfuerzo.horasSprintErr = causa
			},
			contexto: "sumar horas reales del sprint",
		},
		{
			descripcion: "error al sumar horas reales de una historia asignada",
			preparar: func(t *testing.T, e *entornoReporte) {
				cerrarSprint(t, e, 1)
				entornoConSnapshot(t, e, 1)
				sembrarHistoriaReporte(t, e, "p1", "HU-01", 5, domain.EstadoCompletada)
				asignarReporte(e, "p1", 1, 1)
				e.esfuerzo.horasItemErr = causa
			},
			contexto: "sumar horas reales de la historia",
		},
		{
			descripcion: "error al consultar las métricas de defectos del Sprint",
			preparar: func(t *testing.T, e *entornoReporte) {
				cerrarSprint(t, e, 1)
				entornoConSnapshot(t, e, 1)
				e.defectos.metricasSprintErr = causa
			},
			contexto: "consultar métricas de defectos del sprint",
		},
		{
			descripcion: "error al listar los defectos del proyecto",
			preparar: func(t *testing.T, e *entornoReporte) {
				cerrarSprint(t, e, 1)
				entornoConSnapshot(t, e, 1)
				e.defectos.listarErr = causa
			},
			contexto: "listar defectos del proyecto",
		},
	}
	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoReporte()
			sembrarProyectoReporte(t, entorno, "p1")
			caso.preparar(t, entorno)

			reporte, err := entorno.servicio().ReporteSprint("p1", 1)
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
			if errors.Is(err, domain.ErrProyectoNoEncontrado) ||
				errors.Is(err, domain.ErrSprintNoEncontrado) {
				t.Errorf("err = %v, un error de infraestructura no debe mapearse a un error de dominio", err)
			}
			if !reflect.DeepEqual(reporte, ReporteSprint{}) {
				t.Errorf("reporte = %+v, se esperaba ReporteSprint{} (nunca un éxito parcial)", reporte)
			}
		})
	}
}
