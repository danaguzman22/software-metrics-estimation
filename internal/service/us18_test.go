package service

import (
	"errors"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

type fakeRepositorioMetricasDefectosUS18 struct {
	*fakeRepositorioDefectosUS15
	metricasSprint           domain.MetricasDefectos
	metricasProyecto         domain.MetricasDefectos
	metricasSprintErr        error
	metricasProyectoErr      error
	metricasSprintLlamadas   int
	metricasProyectoLlamadas int
	ultimoProyectoID         string
	ultimoSprintNumero       int
}

func nuevoFakeMetricasDefectosUS18() *fakeRepositorioMetricasDefectosUS18 {
	return &fakeRepositorioMetricasDefectosUS18{
		fakeRepositorioDefectosUS15: nuevoFakeDefectosUS15(),
	}
}

func (f *fakeRepositorioMetricasDefectosUS18) MetricasPorSprint(proyectoID string, sprintNumero int) (domain.MetricasDefectos, error) {
	f.metricasSprintLlamadas++
	f.ultimoProyectoID = proyectoID
	f.ultimoSprintNumero = sprintNumero
	if f.metricasSprintErr != nil {
		return domain.MetricasDefectos{}, f.metricasSprintErr
	}
	return f.metricasSprint, nil
}

func (f *fakeRepositorioMetricasDefectosUS18) MetricasPorProyecto(proyectoID string) (domain.MetricasDefectos, error) {
	f.metricasProyectoLlamadas++
	f.ultimoProyectoID = proyectoID
	if f.metricasProyectoErr != nil {
		return domain.MetricasDefectos{}, f.metricasProyectoErr
	}
	return f.metricasProyecto, nil
}

func TestUS18_MetricasDefectosPorSprintValidaProyectoYSprint(t *testing.T) {
	proyectos, backlog, sprints, _, _ := prepararServicioDefectosUS15(t)
	repo := nuevoFakeMetricasDefectosUS18()
	repo.metricasSprint = domain.MetricasDefectos{Detectados: 2, Resueltos: 1}
	servicio := NuevoServicioDefectos(proyectos, backlog, sprints, repo)

	metricas, err := servicio.MetricasDefectosPorSprint("p1", 2)
	if err != nil {
		t.Fatalf("MetricasDefectosPorSprint() devolvió error: %v", err)
	}
	if metricas.Detectados != 2 || metricas.Resueltos != 1 {
		t.Fatalf("métricas = %+v", metricas)
	}
	if repo.metricasSprintLlamadas != 1 || repo.ultimoProyectoID != "p1" || repo.ultimoSprintNumero != 2 {
		t.Fatalf("consulta al repositorio: llamadas=%d proyecto=%q sprint=%d",
			repo.metricasSprintLlamadas, repo.ultimoProyectoID, repo.ultimoSprintNumero)
	}

	_, err = servicio.MetricasDefectosPorSprint("inexistente", 2)
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("error de proyecto = %v; se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if repo.metricasSprintLlamadas != 1 {
		t.Fatal("se consultó el repositorio para un proyecto inexistente")
	}

	_, err = servicio.MetricasDefectosPorSprint("p1", 99)
	if !errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Fatalf("error de Sprint = %v; se esperaba %v", err, domain.ErrSprintNoEncontrado)
	}
	if repo.metricasSprintLlamadas != 1 {
		t.Fatal("se consultó el repositorio para un Sprint inexistente")
	}
}

func TestUS18_MetricasDefectosPorProyectoDevuelveCerosYSuma(t *testing.T) {
	proyectos, backlog, sprints, _, _ := prepararServicioDefectosUS15(t)
	repo := nuevoFakeMetricasDefectosUS18()
	servicio := NuevoServicioDefectos(proyectos, backlog, sprints, repo)

	metricas, err := servicio.MetricasDefectosPorProyecto("p1")
	if err != nil {
		t.Fatalf("proyecto sin defectos: %v", err)
	}
	if metricas != (domain.MetricasDefectos{}) {
		t.Fatalf("métricas vacías = %+v; se esperaban ambos contadores en 0", metricas)
	}

	repo.metricasProyecto = domain.MetricasDefectos{Detectados: 3, Resueltos: 2}
	metricas, err = servicio.MetricasDefectosPorProyecto("p1")
	if err != nil || metricas.Detectados != 3 || metricas.Resueltos != 2 {
		t.Fatalf("métricas = %+v, err=%v", metricas, err)
	}
	if repo.metricasProyectoLlamadas != 2 || repo.ultimoProyectoID != "p1" {
		t.Fatalf("consulta al repositorio: llamadas=%d proyecto=%q", repo.metricasProyectoLlamadas, repo.ultimoProyectoID)
	}
}

func TestUS18_MetricasDefectosEnvuelveErroresDeRepositorio(t *testing.T) {
	proyectos, backlog, sprints, _, _ := prepararServicioDefectosUS15(t)
	repo := nuevoFakeMetricasDefectosUS18()
	servicio := NuevoServicioDefectos(proyectos, backlog, sprints, repo)
	errorInfraestructura := errors.New("fallo de conexión")

	repo.metricasSprintErr = errorInfraestructura
	_, err := servicio.MetricasDefectosPorSprint("p1", 1)
	if err == nil || !errors.Is(err, errorInfraestructura) {
		t.Fatalf("error Sprint = %v; debía envolver la causa", err)
	}

	repo.metricasProyectoErr = errorInfraestructura
	_, err = servicio.MetricasDefectosPorProyecto("p1")
	if err == nil || !errors.Is(err, errorInfraestructura) {
		t.Fatalf("error proyecto = %v; debía envolver la causa", err)
	}
}
