package postgres_test

import (
	"fmt"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage/postgres"
)

func TestUS18_RepositorioDefectosMetricasPorSprintYProyecto(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-018")
	sprintCerradoUno := sprintUS08Postgres(1)
	sprintCerradoUno.Estado = domain.EstadoSprintCerrado
	if err := postgres.NuevoSprints(pool).Agregar("PRJ-018", sprintCerradoUno); err != nil {
		t.Fatalf("sembrar Sprint 1: %v", err)
	}
	sprintCerrado := sprintUS08Postgres(2)
	sprintCerrado.Estado = domain.EstadoSprintCerrado
	if err := postgres.NuevoSprints(pool).Agregar("PRJ-018", sprintCerrado); err != nil {
		t.Fatalf("sembrar Sprint 2: %v", err)
	}
	sembrarSprintUS08Postgres(t, pool, "PRJ-018", 3)
	itemNumero := sembrarItemUS08Postgres(t, pool, "PRJ-018")
	repo := postgres.NuevoDefectos(pool)

	entradas := []domain.Defecto{
		{
			Descripcion:     "Abierto, detectado en 1",
			Severidad:       domain.SeveridadDefectoAlta,
			Estado:          domain.EstadoDefectoAbierto,
			ItemID:          fmt.Sprintf("HU-%02d", itemNumero),
			SprintDeteccion: 1,
		},
		{
			Descripcion:      "Resuelto entre Sprints",
			Severidad:        domain.SeveridadDefectoMedia,
			Estado:           domain.EstadoDefectoResuelto,
			ItemID:           fmt.Sprintf("HU-%02d", itemNumero),
			SprintDeteccion:  1,
			SprintResolucion: intPunteroPostgresUS18(2),
		},
		{
			Descripcion:      "En curso con resolución informada",
			Severidad:        domain.SeveridadDefectoBaja,
			Estado:           domain.EstadoDefectoEnCurso,
			ItemID:           fmt.Sprintf("HU-%02d", itemNumero),
			SprintDeteccion:  3,
			SprintResolucion: intPunteroPostgresUS18(3),
		},
	}
	for _, defecto := range entradas {
		if _, err := repo.Crear("PRJ-018", defecto); err != nil {
			t.Fatalf("crear defecto %q: %v", defecto.Descripcion, err)
		}
	}

	casosSprint := []struct {
		numero   int
		esperado domain.MetricasDefectos
	}{
		{1, domain.MetricasDefectos{Detectados: 2}},
		{2, domain.MetricasDefectos{Resueltos: 1}},
		{3, domain.MetricasDefectos{Detectados: 1}},
		{99, domain.MetricasDefectos{}},
	}
	for _, caso := range casosSprint {
		metricas, err := repo.MetricasPorSprint("PRJ-018", caso.numero)
		if err != nil {
			t.Fatalf("MetricasPorSprint(%d): %v", caso.numero, err)
		}
		if metricas != caso.esperado {
			t.Errorf("MetricasPorSprint(%d) = %+v, se esperaban %+v", caso.numero, metricas, caso.esperado)
		}
	}

	metricas, err := repo.MetricasPorProyecto("PRJ-018")
	if err != nil {
		t.Fatalf("MetricasPorProyecto() devolvió error: %v", err)
	}
	esperadas := domain.MetricasDefectos{Detectados: 3, Resueltos: 1}
	if metricas != esperadas {
		t.Fatalf("MetricasPorProyecto() = %+v, se esperaban %+v", metricas, esperadas)
	}
}

func TestUS18_RepositorioDefectosMetricasProyectoVacioDevuelveCeros(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-018")
	sembrarSprintUS08Postgres(t, pool, "PRJ-018", 1)
	repo := postgres.NuevoDefectos(pool)

	porSprint, err := repo.MetricasPorSprint("PRJ-018", 1)
	if err != nil || porSprint != (domain.MetricasDefectos{}) {
		t.Fatalf("MetricasPorSprint() = %+v, err=%v; se esperaban ceros", porSprint, err)
	}
	porProyecto, err := repo.MetricasPorProyecto("PRJ-018")
	if err != nil || porProyecto != (domain.MetricasDefectos{}) {
		t.Fatalf("MetricasPorProyecto() = %+v, err=%v; se esperaban ceros", porProyecto, err)
	}
}

func intPunteroPostgresUS18(valor int) *int {
	return &valor
}
