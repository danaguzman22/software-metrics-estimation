package web_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// Métodos FASE 3 del doble de ServicioSprintsWeb. El struct vive en
// us22_fase2_test.go y se amplió solo con los campos nuevos (el helper
// compartido queda como la Fase 2 lo dejó).

func (f *servicioSprintsFalso) ListarItemsDelSprint(proyectoID string, sprintNumero int) ([]domain.ItemBacklog, error) {
	f.listarItemsLlamadas++
	f.ultimoProyectoID = proyectoID
	f.ultimoSprintNumero = sprintNumero
	if f.itemsDelSprintErr != nil {
		return nil, f.itemsDelSprintErr
	}
	if !f.proyectosExistentes[proyectoID] {
		return nil, domain.ErrProyectoNoEncontrado
	}
	return append([]domain.ItemBacklog(nil), f.itemsPorSprint[proyectoID][sprintNumero]...), nil
}

func (f *servicioSprintsFalso) StoryPointsCompletados(proyectoID string, sprintNumero int) (int, error) {
	f.completadosLlamadas++
	f.ultimoProyectoID = proyectoID
	f.ultimoSprintNumero = sprintNumero
	if f.completadosErr != nil {
		return 0, f.completadosErr
	}
	if !f.proyectosExistentes[proyectoID] {
		return 0, domain.ErrProyectoNoEncontrado
	}
	return f.completadosPorSprint[proyectoID][sprintNumero], nil
}

func (f *servicioSprintsFalso) VelocidadSprint(proyectoID string, sprintNumero int) (int, error) {
	f.velocidadLlamadas++
	f.ultimoProyectoID = proyectoID
	f.ultimoSprintNumero = sprintNumero
	if f.velocidadErr != nil {
		return 0, f.velocidadErr
	}
	if !f.proyectosExistentes[proyectoID] {
		return 0, domain.ErrProyectoNoEncontrado
	}
	return f.velocidadPorSprint[proyectoID][sprintNumero], nil
}

func (f *servicioSprintsFalso) VelocidadPromedio(proyectoID string) (float64, error) {
	f.promedioLlamadas++
	f.ultimoProyectoID = proyectoID
	if f.promedioErr != nil {
		return 0, f.promedioErr
	}
	if !f.proyectosExistentes[proyectoID] {
		return 0, domain.ErrProyectoNoEncontrado
	}
	return f.promedioPorID[proyectoID], nil
}

func (f *servicioSprintsFalso) PorcentajeHistoriasCompletadas(proyectoID string) (float64, error) {
	f.porcentajeLlamadas++
	f.ultimoProyectoID = proyectoID
	if f.porcentajeErr != nil {
		return 0, f.porcentajeErr
	}
	if !f.proyectosExistentes[proyectoID] {
		return 0, domain.ErrProyectoNoEncontrado
	}
	return f.porcentajePorID[proyectoID], nil
}

func (f *servicioSprintsFalso) MarcarItemCompletada(proyectoID string, sprintNumero int, itemID string) error {
	f.completarLlamadas++
	f.ultimoProyectoID = proyectoID
	f.ultimoSprintNumero = sprintNumero
	f.ultimoItemID = itemID
	if f.completarErr != nil {
		return f.completarErr
	}
	sprint, ok := f.buscarSprintFase3(proyectoID, sprintNumero)
	if !ok {
		return domain.ErrSprintNoEncontrado
	}
	if sprint.Estado != domain.EstadoSprintActivo {
		return domain.ErrSprintCerrado
	}
	items := f.itemsPorSprint[proyectoID][sprintNumero]
	for i := range items {
		if items[i].ID != itemID {
			continue
		}
		if items[i].Estado != domain.EstadoItemPendiente {
			return domain.ErrItemNoPendiente
		}
		items[i].Estado = domain.EstadoCompletada
		f.itemsPorSprint[proyectoID][sprintNumero] = items
		return nil
	}
	return domain.ErrItemNoAsignadoAlSprint
}

func (f *servicioSprintsFalso) CerrarSprint(proyectoID string, sprintNumero int) error {
	f.cerrarLlamadas++
	f.ultimoProyectoID = proyectoID
	f.ultimoSprintNumero = sprintNumero
	if f.cerrarErr != nil {
		return f.cerrarErr
	}
	sprints := f.sprints[proyectoID]
	for i := range sprints {
		if sprints[i].Numero != sprintNumero {
			continue
		}
		if sprints[i].Estado != domain.EstadoSprintActivo {
			return domain.ErrSprintCerrado
		}
		sprints[i].Estado = domain.EstadoSprintCerrado
		f.sprints[proyectoID] = sprints
		// US-09: las historias no completadas vuelven al backlog.
		items := f.itemsPorSprint[proyectoID][sprintNumero]
		restantes := make([]domain.ItemBacklog, 0, len(items))
		for _, item := range items {
			if item.Estado == domain.EstadoCompletada {
				restantes = append(restantes, item)
			}
		}
		f.itemsPorSprint[proyectoID][sprintNumero] = restantes
		return nil
	}
	return domain.ErrSprintNoEncontrado
}

func (f *servicioSprintsFalso) buscarSprintFase3(proyectoID string, sprintNumero int) (domain.Sprint, bool) {
	for _, sprint := range f.sprints[proyectoID] {
		if sprint.Numero == sprintNumero {
			return sprint, true
		}
	}
	return domain.Sprint{}, false
}

// sembrarItemsSprint siembra las historias asignadas a un Sprint del doble.
func sembrarItemsSprintFase3(f *servicioSprintsFalso, proyectoID string, sprintNumero int, items ...domain.ItemBacklog) {
	if f.itemsPorSprint[proyectoID] == nil {
		f.itemsPorSprint[proyectoID] = map[int][]domain.ItemBacklog{}
	}
	f.itemsPorSprint[proyectoID][sprintNumero] = append([]domain.ItemBacklog(nil), items...)
}

// configurarMetricasSprint siembra los totales del doble para un Sprint:
// completados y velocidad (alias), distintos de la suma de las historias.
func configurarMetricasSprintFase3(f *servicioSprintsFalso, proyectoID string, sprintNumero, completados, velocidad int) {
	if f.completadosPorSprint[proyectoID] == nil {
		f.completadosPorSprint[proyectoID] = map[int]int{}
	}
	if f.velocidadPorSprint[proyectoID] == nil {
		f.velocidadPorSprint[proyectoID] = map[int]int{}
	}
	f.completadosPorSprint[proyectoID][sprintNumero] = completados
	f.velocidadPorSprint[proyectoID][sprintNumero] = velocidad
}

// configurarMetricasProyecto siembra la velocidad promedio y el porcentaje del
// doble para el proyecto.
func configurarMetricasProyectoFase3(f *servicioSprintsFalso, proyectoID string, promedio, porcentaje float64) {
	f.promedioPorID[proyectoID] = promedio
	f.porcentajePorID[proyectoID] = porcentaje
}

// CA1 / R33: el listado muestra las historias asignadas de cada Sprint con id,
// título, estado y story points, en el orden devuelto por el servicio.
func TestUS22_Fase3_Sprints_ListarHistoriasAsignadas(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 8)
	sembrarItemsSprintFase3(sprints, "PRJ-001", 1,
		itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 3),
		itemUS22("HU-02", "Cerrar sesión", "Media", "Completada", 5),
	)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec,
		"HU-01", "Iniciar sesión", "Pendiente", "3",
		"HU-02", "Cerrar sesión", "Completada", "5",
	)

	cuerpo := rec.Body.String()
	if posicionUS22(t, cuerpo, "HU-01") > posicionUS22(t, cuerpo, "HU-02") {
		t.Error("HU-01 debería aparecer antes que HU-02: la web debe respetar el orden recibido del servicio")
	}
	if sprints.listarItemsLlamadas != 1 {
		t.Errorf("ListarItemsDelSprint se llamó %d veces, se esperaba 1", sprints.listarItemsLlamadas)
	}
	if sprints.ultimoSprintNumero != 1 {
		t.Errorf("sprintNumero = %d, se esperaba 1", sprints.ultimoSprintNumero)
	}
}

// CA1 / R34 / D25: "Marcar completada" solo en historias Pendientes de un
// Sprint Activo; "Cerrar Sprint" solo en el Sprint Activo del proyecto.
func TestUS22_Fase3_Sprints_BotonesSoloSprintActivo(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 8)
	sprints.agregarSprint("PRJ-001", sprintUS22(
		2, "Segunda entrega", domain.EstadoSprintCerrado,
		fechaUS22(2026, time.March, 15), fechaUS22(2026, time.March, 28),
	), 13)
	sembrarItemsSprintFase3(sprints, "PRJ-001", 1,
		itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 3),
		itemUS22("HU-02", "Cambiar contraseña", "Media", "Completada", 5),
	)
	sembrarItemsSprintFase3(sprints, "PRJ-001", 2,
		itemUS22("HU-03", "Recuperar contraseña", "Media", "Pendiente", 8),
	)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	cuerpo := rec.Body.String()

	if !strings.Contains(cuerpo, `/proyectos/PRJ-001/sprints/1/items/HU-01/completar`) {
		t.Errorf("falta el botón \"Marcar completada\" de HU-01 (pendiente, Sprint Activo); cuerpo:\n%s", cuerpo)
	}
	if strings.Contains(cuerpo, `/proyectos/PRJ-001/sprints/1/items/HU-02/completar`) {
		t.Errorf("HU-02 está Completada y no debería ofrecer \"Marcar completada\"; cuerpo:\n%s", cuerpo)
	}
	if strings.Contains(cuerpo, `/proyectos/PRJ-001/sprints/2/items/HU-03/completar`) {
		t.Errorf("el Sprint 2 está Cerrado y no debería ofrecer \"Marcar completada\"; cuerpo:\n%s", cuerpo)
	}
	if !strings.Contains(cuerpo, `/proyectos/PRJ-001/sprints/1/cerrar`) {
		t.Errorf("falta el botón \"Cerrar Sprint\" del Sprint Activo; cuerpo:\n%s", cuerpo)
	}
	if strings.Contains(cuerpo, `/proyectos/PRJ-001/sprints/2/cerrar`) {
		t.Errorf("el Sprint 2 está Cerrado y no debería ofrecer \"Cerrar Sprint\"; cuerpo:\n%s", cuerpo)
	}
}

// CA1 / R35 / R24: completar una historia pendiente delega en el servicio y
// responde 303 hacia el listado con el cuerpo vacío.
func TestUS22_Fase3_Sprints_MarcarCompletada(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 0)
	sembrarItemsSprintFase3(sprints, "PRJ-001", 1,
		itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 3),
	)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/sprints/1/items/HU-01/completar", nil)
	exigirRedireccionUS22(t, rec, "/proyectos/PRJ-001/sprints")

	if sprints.completarLlamadas != 1 {
		t.Fatalf("MarcarItemCompletada se llamó %d veces, se esperaba 1", sprints.completarLlamadas)
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

// CA1 / R36: cerrar el Sprint delega en el servicio, responde 303 y la leyenda
// informa que las historias no completadas vuelven al backlog; tras el cierre,
// la historia no completada ya no aparece asignada al Sprint.
func TestUS22_Fase3_Sprints_CerrarSprint(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 8)
	sembrarItemsSprintFase3(sprints, "PRJ-001", 1,
		itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 3),
		itemUS22("HU-02", "Cambiar contraseña", "Media", "Completada", 5),
	)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/sprints/1/cerrar", nil)
	exigirRedireccionUS22(t, rec, "/proyectos/PRJ-001/sprints")

	if sprints.cerrarLlamadas != 1 {
		t.Fatalf("CerrarSprint se llamó %d veces, se esperaba 1", sprints.cerrarLlamadas)
	}
	if sprints.ultimoProyectoID != "PRJ-001" {
		t.Errorf("proyectoID = %q, se esperaba PRJ-001", sprints.ultimoProyectoID)
	}
	if sprints.ultimoSprintNumero != 1 {
		t.Errorf("sprintNumero = %d, se esperaba 1", sprints.ultimoSprintNumero)
	}

	rec = solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "Al cerrar el Sprint, las historias no completadas vuelven al backlog")
	exigirCuerpoNoContieneUS22(t, rec, "HU-01", `/proyectos/PRJ-001/sprints/1/items/HU-01/completar`)
	exigirCuerpoContieneUS22(t, rec, "HU-02")
}

// CA1 / R37: los Story Points mostrados por Sprint son los del servicio, no la
// suma de las historias exhibidas.
func TestUS22_Fase3_Sprints_PlanificadosCompletadosVelocidad(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 34)
	sembrarItemsSprintFase3(sprints, "PRJ-001", 1,
		itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 3),
		itemUS22("HU-02", "Cerrar sesión", "Media", "Pendiente", 5),
	)
	configurarMetricasSprintFase3(sprints, "PRJ-001", 1, 21, 21)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "34", "21.00")

	if sprints.completadosLlamadas != 1 {
		t.Errorf("StoryPointsCompletados se llamó %d veces, se esperaba 1", sprints.completadosLlamadas)
	}
	if sprints.velocidadLlamadas != 1 {
		t.Errorf("VelocidadSprint se llamó %d veces, se esperaba 1", sprints.velocidadLlamadas)
	}
}

// CA1 / R38: la velocidad promedio se muestra con dos decimales.
func TestUS22_Fase3_Sprints_VelocidadPromedio(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintCerrado,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 9)
	configurarMetricasProyectoFase3(sprints, "PRJ-001", 4.5, 0)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "4.50")
}

// CA1 / R38: sin Sprints cerrados la velocidad promedio es 0.00.
func TestUS22_Fase3_Sprints_VelocidadPromedioSinCerrados(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 0)
	configurarMetricasProyectoFase3(sprints, "PRJ-001", 0, 0)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "0.00")
}

// CA1 / R39: el porcentaje de historias completadas se muestra con el valor
// devuelto por el servicio, con dos decimales y sufijo %.
func TestUS22_Fase3_Sprints_PorcentajeHistoriasCompletadas(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintCerrado,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 9)
	configurarMetricasProyectoFase3(sprints, "PRJ-001", 0, 62.5)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "62.50%")
}

// CA1 / R39: con cero historias el porcentaje es 0.00%.
func TestUS22_Fase3_Sprints_PorcentajeConCeroHistorias(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 0)
	configurarMetricasProyectoFase3(sprints, "PRJ-001", 0, 0)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "0.00%")
}

// CA2 / R40 / E19: completar una historia de un Sprint cerrado muestra el
// mensaje del dominio con HTTP 422.
func TestUS22_Fase3_Sprints_MarcarCompletadaSprintCerrado(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintCerrado,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 8)
	sembrarItemsSprintFase3(sprints, "PRJ-001", 1,
		itemUS22("HU-02", "Cambiar contraseña", "Media", "Completada", 5),
	)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/sprints/1/items/HU-02/completar", nil)
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, domain.ErrSprintCerrado.Error())

	if sprints.completarLlamadas != 1 {
		t.Errorf("MarcarItemCompletada se llamó %d veces, se esperaba 1", sprints.completarLlamadas)
	}
}

// CA2 / R40 / E20: completar una historia que no está Pendiente muestra el
// mensaje del dominio con HTTP 422.
func TestUS22_Fase3_Sprints_MarcarCompletadaNoPendiente(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 5)
	sembrarItemsSprintFase3(sprints, "PRJ-001", 1,
		itemUS22("HU-02", "Cambiar contraseña", "Media", "Completada", 5),
	)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/sprints/1/items/HU-02/completar", nil)
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, domain.ErrItemNoPendiente.Error())

	if sprints.completarLlamadas != 1 {
		t.Errorf("MarcarItemCompletada se llamó %d veces, se esperaba 1", sprints.completarLlamadas)
	}
}

// CA2 / R40 / E21: completar una historia no asignada al Sprint muestra el
// mensaje del dominio con HTTP 422.
func TestUS22_Fase3_Sprints_MarcarCompletadaNoAsignadaAlSprint(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 0)
	sembrarItemsSprintFase3(sprints, "PRJ-001", 1,
		itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 3),
	)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/sprints/1/items/HU-09/completar", nil)
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, domain.ErrItemNoAsignadoAlSprint.Error())

	if sprints.completarLlamadas != 1 {
		t.Errorf("MarcarItemCompletada se llamó %d veces, se esperaba 1", sprints.completarLlamadas)
	}
}

// CA2 / R40 / E22: cerrar un Sprint que no está activo muestra el mensaje del
// dominio con HTTP 422.
func TestUS22_Fase3_Sprints_CerrarSprintCerrado(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintCerrado,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 8)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/sprints/1/cerrar", nil)
	exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
	exigirCuerpoContieneUS22(t, rec, domain.ErrSprintCerrado.Error())

	if sprints.cerrarLlamadas != 1 {
		t.Errorf("CerrarSprint se llamó %d veces, se esperaba 1", sprints.cerrarLlamadas)
	}
}

// CA2 / R41 / E23: completar con un Sprint inexistente muestra HTTP 404 con el
// mensaje del dominio.
func TestUS22_Fase3_Sprints_CompletarSprintInexistente404(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 0)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/sprints/9/items/HU-01/completar", nil)
	exigirEstadoUS22(t, rec, http.StatusNotFound)
	exigirCuerpoContieneUS22(t, rec, domain.ErrSprintNoEncontrado.Error())

	if sprints.completarLlamadas != 1 {
		t.Errorf("MarcarItemCompletada se llamó %d veces, se esperaba 1", sprints.completarLlamadas)
	}
}

// CA2 / R41 / E23: cerrar un Sprint inexistente muestra HTTP 404 con el
// mensaje del dominio.
func TestUS22_Fase3_Sprints_CerrarSprintInexistente404(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 0)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/sprints/9/cerrar", nil)
	exigirEstadoUS22(t, rec, http.StatusNotFound)
	exigirCuerpoContieneUS22(t, rec, domain.ErrSprintNoEncontrado.Error())

	if sprints.cerrarLlamadas != 1 {
		t.Errorf("CerrarSprint se llamó %d veces, se esperaba 1", sprints.cerrarLlamadas)
	}
}

// CA2 / R41 / E24: completar con un ítem inexistente muestra HTTP 404 con el
// mensaje del dominio.
func TestUS22_Fase3_Sprints_MarcarCompletadaItemInexistente404(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 0)
	sprints.completarErr = domain.ErrItemNoEncontrado
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodPost, "/proyectos/PRJ-001/sprints/1/items/HU-99/completar", nil)
	exigirEstadoUS22(t, rec, http.StatusNotFound)
	exigirCuerpoContieneUS22(t, rec, domain.ErrItemNoEncontrado.Error())

	if sprints.completarLlamadas != 1 {
		t.Errorf("MarcarItemCompletada se llamó %d veces, se esperaba 1", sprints.completarLlamadas)
	}
}

// CA2 / V11: un número de Sprint no numérico en la ruta produce HTTP 422,
// el mensaje de formato y no se llama al servicio.
func TestUS22_Fase3_Sprints_NumeroRutaInvalido(t *testing.T) {
	casos := []struct {
		nombre string
		ruta   string
	}{
		{nombre: "completar", ruta: "/proyectos/PRJ-001/sprints/abc/items/HU-01/completar"},
		{nombre: "cerrar", ruta: "/proyectos/PRJ-001/sprints/abc/cerrar"},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoUS22(
				registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
			)
			sprints := nuevoServicioSprintsFalso()
			sprints.agregarSprint("PRJ-001", sprintUS22(
				1, "Primera entrega", domain.EstadoSprintActivo,
				fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
			), 0)
			handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

			rec := solicitarUS22(handler, http.MethodPost, caso.ruta, nil)
			exigirEstadoUS22(t, rec, http.StatusUnprocessableEntity)
			exigirCuerpoContieneUS22(t, rec, "El número del sprint debe ser un número entero")

			if sprints.completarLlamadas != 0 {
				t.Errorf("MarcarItemCompletada se llamó %d veces y no debía ejecutarse", sprints.completarLlamadas)
			}
			if sprints.cerrarLlamadas != 0 {
				t.Errorf("CerrarSprint se llamó %d veces y no debía ejecutarse", sprints.cerrarLlamadas)
			}
		})
	}
}

// CA2 / R24 / R35 / R36: completar y cerrar responden 303 al listado con el
// cuerpo vacío; la respuesta no renderiza contenido.
func TestUS22_Fase3_Sprints_PostRedirige303(t *testing.T) {
	casos := []struct {
		nombre string
		ruta   string
	}{
		{nombre: "completar", ruta: "/proyectos/PRJ-001/sprints/1/items/HU-01/completar"},
		{nombre: "cerrar", ruta: "/proyectos/PRJ-001/sprints/1/cerrar"},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoUS22(
				registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
			)
			sprints := nuevoServicioSprintsFalso()
			sprints.agregarSprint("PRJ-001", sprintUS22(
				1, "Primera entrega", domain.EstadoSprintActivo,
				fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
			), 0)
			sembrarItemsSprintFase3(sprints, "PRJ-001", 1,
				itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 3),
			)
			handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

			rec := solicitarUS22(handler, http.MethodPost, caso.ruta, nil)
			exigirRedireccionUS22(t, rec, "/proyectos/PRJ-001/sprints")
		})
	}
}

// CA2 / R25: un método no permitido en las rutas de completar y cerrar
// devuelve HTTP 405 y no ejecuta el handler.
func TestUS22_Fase3_Sprints_MetodoNoPermitido(t *testing.T) {
	casos := []struct {
		nombre string
		ruta   string
	}{
		{nombre: "completar con GET", ruta: "/proyectos/PRJ-001/sprints/1/items/HU-01/completar"},
		{nombre: "cerrar con GET", ruta: "/proyectos/PRJ-001/sprints/1/cerrar"},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoUS22(
				registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
			)
			sprints := nuevoServicioSprintsFalso()
			sprints.agregarSprint("PRJ-001", sprintUS22(
				1, "Primera entrega", domain.EstadoSprintActivo,
				fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
			), 0)
			handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

			rec := solicitarUS22(handler, http.MethodGet, caso.ruta, nil)
			exigirEstadoUS22(t, rec, http.StatusMethodNotAllowed)

			if sprints.completarLlamadas != 0 {
				t.Errorf("MarcarItemCompletada se llamó %d veces y no debía ejecutarse", sprints.completarLlamadas)
			}
			if sprints.cerrarLlamadas != 0 {
				t.Errorf("CerrarSprint se llamó %d veces y no debía ejecutarse", sprints.cerrarLlamadas)
			}
		})
	}
}

// CA2 / R23 / R42: un error de infraestructura del servicio se registra en el
// log, se responde con HTTP 500 y el detalle no llega al cuerpo.
func TestUS22_Fase3_Sprints_ErrorInfraestructura(t *testing.T) {
	causa := errors.New("fallo de conexión a postgres: timeout al obtener el sprint")
	casos := []struct {
		nombre string
		ruta   string
		metodo string
		fondo  func(*servicioSprintsFalso)
	}{
		{
			nombre: "listado",
			ruta:   "/proyectos/PRJ-001/sprints",
			metodo: http.MethodGet,
			fondo:  func(f *servicioSprintsFalso) { f.itemsDelSprintErr = causa },
		},
		{
			nombre: "completar",
			ruta:   "/proyectos/PRJ-001/sprints/1/items/HU-01/completar",
			metodo: http.MethodPost,
			fondo:  func(f *servicioSprintsFalso) { f.completarErr = causa },
		},
		{
			nombre: "cerrar",
			ruta:   "/proyectos/PRJ-001/sprints/1/cerrar",
			metodo: http.MethodPost,
			fondo:  func(f *servicioSprintsFalso) { f.cerrarErr = causa },
		},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			proyectos, backlog := nuevoEntornoUS22(
				registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
			)
			sprints := nuevoServicioSprintsFalso()
			sprints.agregarSprint("PRJ-001", sprintUS22(
				1, "Primera entrega", domain.EstadoSprintActivo,
				fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
			), 0)
			sembrarItemsSprintFase3(sprints, "PRJ-001", 1,
				itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 3),
			)
			caso.fondo(sprints)
			handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

			rec := solicitarUS22(handler, caso.metodo, caso.ruta, nil)
			exigirEstadoUS22(t, rec, http.StatusInternalServerError)
			exigirCuerpoContieneUS22(t, rec, "Ocurrió un error inesperado. Intentá nuevamente.")
			exigirCuerpoNoContieneUS22(t, rec, causa.Error())
		})
	}
}

// CA3 / R37 / R38 / R39 / R42: la web muestra los totales del servicio sin
// recalcularlos, aunque no coincidan con la suma de las historias exhibidas.
func TestUS22_Fase3_Sprints_WebNoRecalculaMetricas(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.June, 30)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 34)
	sembrarItemsSprintFase3(sprints, "PRJ-001", 1,
		itemUS22("HU-01", "Iniciar sesión", "Alta", "Pendiente", 3),
		itemUS22("HU-02", "Cerrar sesión", "Media", "Pendiente", 5),
	)
	configurarMetricasSprintFase3(sprints, "PRJ-001", 1, 21, 21)
	configurarMetricasProyectoFase3(sprints, "PRJ-001", 6.5, 75)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec,
		"34",     // planificados del doble, no la suma (3 + 5)
		"21",     // completados y velocidad del doble, no la suma
		"6.50",   // velocidad promedio del doble
		"75.00%", // porcentaje del doble, no recalculado desde las historias
	)
}
