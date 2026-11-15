package web_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func (f *servicioSprintsFalso) ListarSprintsCerrados(proyectoID string) ([]domain.Sprint, error) {
	if !f.proyectosExistentes[proyectoID] {
		return nil, domain.ErrProyectoNoEncontrado
	}
	cerrados := make([]domain.Sprint, 0)
	for _, sprint := range f.sprints[proyectoID] {
		if sprint.Estado == domain.EstadoSprintCerrado {
			cerrados = append(cerrados, sprint)
		}
	}
	return cerrados, nil
}

func TestUS10_HistorialMuestraSprintsCerradosYMetricas(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.December, 31)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Primera entrega", domain.EstadoSprintCerrado,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 8)
	sprints.agregarSprint("PRJ-001", sprintUS22(
		2, "Iteración actual", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 15), fechaUS22(2026, time.March, 28),
	), 13)
	sprints.completadosPorSprint["PRJ-001"] = map[int]int{1: 5}
	sprints.velocidadPorSprint["PRJ-001"] = map[int]int{1: 5}
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints/historico", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec,
		"Histórico de Sprints",
		"Primera entrega",
		"8",
		"5",
	)
	exigirCuerpoNoContieneUS22(t, rec, "Iteración actual")
	if strings.Contains(rec.Body.String(), "No hay Sprints cerrados") {
		t.Fatal("se mostró el estado vacío con Sprints cerrados disponibles")
	}
	if sprints.completadosLlamadas != 1 || sprints.velocidadLlamadas != 1 {
		t.Errorf("consultas de métricas completados=%d velocidad=%d; se esperaba una de cada una",
			sprints.completadosLlamadas, sprints.velocidadLlamadas)
	}
}

func TestUS10_HistorialSinSprintsCerradosMuestraListaVacia(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.December, 31)),
	)
	sprints := nuevoServicioSprintsFalso()
	sprints.agregarSprint("PRJ-001", sprintUS22(
		1, "Iteración actual", domain.EstadoSprintActivo,
		fechaUS22(2026, time.March, 1), fechaUS22(2026, time.March, 14),
	), 8)
	handler := nuevoHandlerUS22(t, proyectos, backlog, sprints)

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints/historico", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, "No hay Sprints cerrados")
	exigirCuerpoNoContieneUS22(t, rec, "Iteración actual")
}

func TestUS10_PaginaSprintsEnlazaAlHistorico(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(
		registroProyectoUS22("PRJ-001", "Sistema de Métricas", fechaUS22(2026, time.March, 1), fechaUS22(2026, time.December, 31)),
	)
	handler := nuevoHandlerUS22(t, proyectos, backlog, nuevoServicioSprintsFalso())

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-001/sprints", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	exigirCuerpoContieneUS22(t, rec, `href="/proyectos/PRJ-001/sprints/historico"`)
}
