package postgres_test

import (
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage/postgres"
)

// CA1 / R7 / D13: PostgreSQL genera el ID del proyecto con
// DEFAULT gen_random_uuid()::text y el registro creado se puede recuperar.
func TestUS22_Fase1_Postgres_ProyectoIDGenerado(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	repo := postgres.NuevoProyectos(pool)

	proyecto, err := domain.NuevoProyecto(
		"Proyecto con ID generado",
		time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, time.June, 30, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("NuevoProyecto() devolvió un error inesperado: %v", err)
	}

	creado, err := repo.Crear(proyecto)
	if err != nil {
		t.Fatalf("Crear() devolvió un error inesperado: %v", err)
	}
	if creado.ID == "" {
		t.Fatal("Crear() debería devolver un ID generado no vacío")
	}

	obtenido, err := repo.Obtener(creado.ID)
	if err != nil {
		t.Fatalf("Obtener() tras Crear devolvió un error inesperado: %v", err)
	}
	if obtenido.ID != creado.ID {
		t.Fatalf("ID = %q, se esperaba %q", obtenido.ID, creado.ID)
	}
	if obtenido.Proyecto.Nombre != proyecto.Nombre {
		t.Errorf("Nombre = %q, se esperaba %q", obtenido.Proyecto.Nombre, proyecto.Nombre)
	}
	if obtenido.Proyecto.Estado != proyecto.Estado {
		t.Errorf("Estado = %q, se esperaba %q", obtenido.Proyecto.Estado, proyecto.Estado)
	}
	if !obtenido.Proyecto.FechaInicio.Equal(proyecto.FechaInicio) {
		t.Errorf("FechaInicio = %v, se esperaba %v", obtenido.Proyecto.FechaInicio, proyecto.FechaInicio)
	}
	if !obtenido.Proyecto.FechaFin.Equal(proyecto.FechaFin) {
		t.Errorf("FechaFin = %v, se esperaba %v", obtenido.Proyecto.FechaFin, proyecto.FechaFin)
	}
}
