package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage/postgres"
)

// sprintUS07Postgres arma un Sprint activo de referencia para la suite de
// integración de US-07.
func sprintUS07Postgres(numero int) domain.Sprint {
	return domain.Sprint{
		Numero:      numero,
		Objetivo:    "Planificar la primera iteración",
		FechaInicio: time.Date(2026, time.March, 2, 0, 0, 0, 0, time.UTC),
		FechaFin:    time.Date(2026, time.March, 20, 0, 0, 0, 0, time.UTC),
		Estado:      domain.EstadoSprintActivo,
	}
}

// CA1 / R8: la clave primaria (proyecto_id, numero) impide duplicar el número
// dentro del proyecto y el adaptador traduce el error de unicidad.
func TestUS07_RepositorioSprints_NumeroDuplicado(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	repo := postgres.NuevoSprints(pool)

	sprint := sprintUS07Postgres(1)
	if err := repo.Agregar("PRJ-001", sprint); err != nil {
		t.Fatalf("agregar primer Sprint: %v", err)
	}

	err := repo.Agregar("PRJ-001", sprint)
	if !errors.Is(err, domain.ErrSprintDuplicado) {
		t.Fatalf("error de duplicado = %v, se esperaba %v", err, domain.ErrSprintDuplicado)
	}
}

// CA3 / R8: el índice único parcial impide dos Sprints activos en el mismo
// proyecto y el adaptador traduce el error de unicidad.
func TestUS07_RepositorioSprints_IndiceUnicoActivo(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	repo := postgres.NuevoSprints(pool)

	if err := repo.Agregar("PRJ-001", sprintUS07Postgres(1)); err != nil {
		t.Fatalf("agregar primer Sprint activo: %v", err)
	}

	err := repo.Agregar("PRJ-001", sprintUS07Postgres(2))
	if !errors.Is(err, domain.ErrSprintActivoExistente) {
		t.Fatalf("error de Sprint activo = %v, se esperaba %v", err, domain.ErrSprintActivoExistente)
	}
}

// CA1 / R4: el mismo número en otro proyecto es válido porque la unicidad es
// por proyecto.
func TestUS07_RepositorioSprints_MismoNumeroOtroProyecto(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	seedProject(t, pool, "PRJ-002")
	repo := postgres.NuevoSprints(pool)

	if err := repo.Agregar("PRJ-001", sprintUS07Postgres(1)); err != nil {
		t.Fatalf("agregar Sprint en PRJ-001: %v", err)
	}
	if err := repo.Agregar("PRJ-002", sprintUS07Postgres(1)); err != nil {
		t.Fatalf("agregar Sprint en PRJ-002: %v", err)
	}

	var cantidad int
	if err := pool.QueryRow(context.Background(), `
		SELECT count(*) FROM sprints WHERE numero = 1
	`).Scan(&cantidad); err != nil {
		t.Fatalf("contar Sprints con número 1: %v", err)
	}
	if cantidad != 2 {
		t.Fatalf("se esperaban 2 Sprints con número 1 en proyectos distintos, hay %d", cantidad)
	}
}
