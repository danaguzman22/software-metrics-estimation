package web_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/web"
)

type servicioDefectosMetricasFalsoUS18 struct {
	*servicioDefectosFalsoUS15
	porProyecto        domain.MetricasDefectos
	porSprint          map[int]domain.MetricasDefectos
	proyectoErr        error
	sprintErr          error
	sprintsConsultados []int
}

func (f *servicioDefectosMetricasFalsoUS18) MetricasDefectosPorProyecto(string) (domain.MetricasDefectos, error) {
	return f.porProyecto, f.proyectoErr
}

func (f *servicioDefectosMetricasFalsoUS18) MetricasDefectosPorSprint(_ string, sprintNumero int) (domain.MetricasDefectos, error) {
	f.sprintsConsultados = append(f.sprintsConsultados, sprintNumero)
	return f.porSprint[sprintNumero], f.sprintErr
}

func TestUS18_WebMuestraTotalesYTodosLosSprintsEnOrden(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(registroProyectoUS22(
		"PRJ-018", "Proyecto de métricas", fechaUS22(2026, 1, 1), fechaUS22(2026, 12, 31),
	))
	backlog.proyectosExistentes["PRJ-018"] = true
	sprints := nuevoServicioSprintsFalso()
	sprints.marcarProyectoExistente("PRJ-018")
	sprints.agregarSprint("PRJ-018", domain.Sprint{Numero: 1, Objetivo: "Sprint uno"}, 0)
	sprints.agregarSprint("PRJ-018", domain.Sprint{Numero: 2, Objetivo: "Sprint dos"}, 0)
	sprints.agregarSprint("PRJ-018", domain.Sprint{Numero: 3, Objetivo: "Sprint tres"}, 0)
	defectos := &servicioDefectosMetricasFalsoUS18{
		servicioDefectosFalsoUS15: &servicioDefectosFalsoUS15{},
		porProyecto:               domain.MetricasDefectos{Detectados: 4, Resueltos: 2},
		porSprint: map[int]domain.MetricasDefectos{
			1: {Detectados: 2, Resueltos: 1},
			2: {},
			3: {Detectados: 2, Resueltos: 1},
		},
	}
	handler := web.NuevoHandlerAplicacion(&pingFalso{}, web.Dependencias{
		Proyectos: proyectos,
		Backlog:   backlog,
		Sprints:   sprints,
		Defectos:  defectos,
	})

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-018/defectos", nil)
	exigirEstadoUS22(t, rec, http.StatusOK)
	for _, texto := range []string{
		"Defectos detectados: 4",
		"Defectos resueltos: 2",
		"Sprint 1",
		"Sprint 2",
		"Sprint 3",
		"Detectados: 0",
		"Resueltos: 0",
	} {
		if !strings.Contains(rec.Body.String(), texto) {
			t.Errorf("la página no contiene %q; cuerpo:\n%s", texto, rec.Body.String())
		}
	}
	if len(defectos.sprintsConsultados) != 3 ||
		defectos.sprintsConsultados[0] != 1 ||
		defectos.sprintsConsultados[1] != 2 ||
		defectos.sprintsConsultados[2] != 3 {
		t.Fatalf("métricas consultadas para Sprints = %v; se esperaban todos los Sprints del proyecto", defectos.sprintsConsultados)
	}
}

func TestUS18_WebNoMuestraMetricasParcialesSiFallaConsulta(t *testing.T) {
	proyectos, backlog := nuevoEntornoUS22(registroProyectoUS22(
		"PRJ-018", "Proyecto de métricas", fechaUS22(2026, 1, 1), fechaUS22(2026, 12, 31),
	))
	backlog.proyectosExistentes["PRJ-018"] = true
	sprints := nuevoServicioSprintsFalso()
	sprints.marcarProyectoExistente("PRJ-018")
	sprints.agregarSprint("PRJ-018", domain.Sprint{Numero: 1, Objetivo: "Sprint uno"}, 0)
	defectos := &servicioDefectosMetricasFalsoUS18{
		servicioDefectosFalsoUS15: &servicioDefectosFalsoUS15{},
		sprintErr:                 errors.New("fallo de conexión"),
	}
	handler := web.NuevoHandlerAplicacion(&pingFalso{}, web.Dependencias{
		Proyectos: proyectos,
		Backlog:   backlog,
		Sprints:   sprints,
		Defectos:  defectos,
	})

	rec := solicitarUS22(handler, http.MethodGet, "/proyectos/PRJ-018/defectos", nil)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("código = %d, se esperaba %d", rec.Code, http.StatusInternalServerError)
	}
	if strings.Contains(rec.Body.String(), "Defectos detectados: 0") {
		t.Fatal("la web no debe presentar métricas por defecto si falta la consulta de métricas")
	}
}

var _ web.ServicioDefectosWeb = (*servicioDefectosMetricasFalsoUS18)(nil)
