package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

// repositorioEsfuerzoUS14 es el contrato mínimo que el adaptador PostgreSQL de
// esfuerzo debe cumplir para la suite de integración de US-14.
type repositorioEsfuerzoUS14 interface {
	Agregar(proyectoID string, itemNumero int, sprintNumero *int, registro domain.RegistroEsfuerzo) error
}

// nuevoRepositorioEsfuerzoUS14 instancia el adaptador PostgreSQL de esfuerzo.
// Es el único punto que depende del nombre del constructor del adaptador.
func nuevoRepositorioEsfuerzoUS14(pool *pgxpool.Pool) repositorioEsfuerzoUS14 {
	return postgres.NuevoEsfuerzo(pool)
}

// CA1 / R1 / D17: el adaptador persiste integrante, fecha, actividad y horas
// asociados al proyecto y a la historia.
func TestUS14_RepositorioEsfuerzo_Agregar(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	numero := sembrarItemUS08Postgres(t, pool, "PRJ-001")
	repo := nuevoRepositorioEsfuerzoUS14(pool)

	fecha := time.Date(2026, time.November, 15, 0, 0, 0, 0, time.UTC)
	registro := domain.RegistroEsfuerzo{
		Integrante: "Ana Pérez",
		Fecha:      fecha,
		Actividad:  "Análisis",
		Horas:      1.5,
	}

	if err := repo.Agregar("PRJ-001", numero, nil, registro); err != nil {
		t.Fatalf("Agregar() devolvió un error inesperado: %v", err)
	}

	var (
		integrante string
		fechaDB    time.Time
		actividad  string
		horas      float64
	)
	err := pool.QueryRow(context.Background(), `
		SELECT integrante, fecha, actividad, horas
		FROM esfuerzos
		WHERE proyecto_id = $1 AND item_numero = $2
	`, "PRJ-001", numero).Scan(&integrante, &fechaDB, &actividad, &horas)
	if err != nil {
		t.Fatalf("leer el esfuerzo persistido: %v", err)
	}

	if integrante != registro.Integrante {
		t.Errorf("integrante = %q, se esperaba %q", integrante, registro.Integrante)
	}
	if fechaDB.Format("2006-01-02") != "2026-11-15" {
		t.Errorf("fecha = %v, se esperaba 2026-11-15", fechaDB)
	}
	if actividad != registro.Actividad {
		t.Errorf("actividad = %q, se esperaba %q", actividad, registro.Actividad)
	}
	if horas != registro.Horas {
		t.Errorf("horas = %v, se esperaba %v", horas, registro.Horas)
	}
}

// CA3 / R3 / D17: la FK compuesta impide registrar esfuerzo en una historia
// inexistente y no deja filas.
func TestUS14_RepositorioEsfuerzo_ItemInexistente(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	numero := sembrarItemUS08Postgres(t, pool, "PRJ-001")
	repo := nuevoRepositorioEsfuerzoUS14(pool)

	registro := domain.RegistroEsfuerzo{
		Integrante: "Ana Pérez",
		Fecha:      time.Date(2026, time.November, 15, 0, 0, 0, 0, time.UTC),
		Actividad:  "Análisis",
		Horas:      1.5,
	}

	if err := repo.Agregar("PRJ-001", numero+1, nil, registro); err == nil {
		t.Fatal("Agregar() con una historia inexistente debía devolver un error")
	}

	var filas int
	if err := pool.QueryRow(context.Background(), `
		SELECT count(*) FROM esfuerzos WHERE proyecto_id = $1 AND item_numero = $2
	`, "PRJ-001", numero+1).Scan(&filas); err != nil {
		t.Fatalf("contar esfuerzos persistidos: %v", err)
	}
	if filas != 0 {
		t.Errorf("no debía persistirse ningún registro, hay %d", filas)
	}
}

// D20: SprintDeItem devuelve el Sprint al que pertenece la historia.
func TestUS14_RepositorioAsignaciones_SprintDeItem(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	sembrarSprintUS08Postgres(t, pool, "PRJ-001", 1)
	numero := sembrarItemUS08Postgres(t, pool, "PRJ-001")
	repo := postgres.NuevoAsignaciones(pool)

	if err := repo.Agregar("PRJ-001", 1, numero); err != nil {
		t.Fatalf("Agregar() devolvió un error inesperado: %v", err)
	}

	sprintNumero, asignado, err := repo.SprintDeItem("PRJ-001", numero)
	if err != nil {
		t.Fatalf("SprintDeItem() devolvió un error inesperado: %v", err)
	}
	if !asignado {
		t.Fatal("SprintDeItem() = false, se esperaba true para una historia asignada")
	}
	if sprintNumero != 1 {
		t.Errorf("SprintDeItem() = %d, se esperaba 1", sprintNumero)
	}
}

// D20 / D19: una historia sin Sprint asignado devuelve (0, false, nil).
func TestUS14_RepositorioAsignaciones_SprintDeItemSinAsignar(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	numero := sembrarItemUS08Postgres(t, pool, "PRJ-001")
	repo := postgres.NuevoAsignaciones(pool)

	sprintNumero, asignado, err := repo.SprintDeItem("PRJ-001", numero)
	if err != nil {
		t.Fatalf("SprintDeItem() devolvió un error inesperado: %v", err)
	}
	if asignado {
		t.Fatal("SprintDeItem() = true, se esperaba false para una historia sin Sprint")
	}
	if sprintNumero != 0 {
		t.Errorf("SprintDeItem() = %d, se esperaba 0", sprintNumero)
	}
}
