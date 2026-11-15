package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage/postgres"
	"github.com/jackc/pgx/v5"
)

func TestUS17_RepositorioEsfuerzoSumaHorasPorHistoriaSprintYProyecto(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-017")
	sembrarSprintUS08Postgres(t, pool, "PRJ-017", 1)
	itemNumero := sembrarItemUS08Postgres(t, pool, "PRJ-017")
	sprintNumero := 1
	repo := postgres.NuevoEsfuerzo(pool)
	fecha := time.Date(2026, time.November, 15, 0, 0, 0, 0, time.UTC)

	for _, registro := range []domain.RegistroEsfuerzo{
		{Integrante: "Ana", Fecha: fecha, Actividad: "A", Horas: 1.25},
		{Integrante: "Ana", Fecha: fecha, Actividad: "B", Horas: 2.5},
	} {
		if err := repo.Agregar("PRJ-017", itemNumero, &sprintNumero, registro); err != nil {
			t.Fatalf("Agregar() devolvió un error: %v", err)
		}
	}
	if err := repo.Agregar("PRJ-017", itemNumero, nil, domain.RegistroEsfuerzo{
		Integrante: "Ana", Fecha: fecha, Actividad: "sin Sprint", Horas: 10,
	}); err != nil {
		t.Fatalf("Agregar() sin Sprint devolvió un error: %v", err)
	}

	porItem, err := repo.HorasRealesPorItem("PRJ-017", itemNumero)
	if err != nil || porItem != 13.75 {
		t.Fatalf("HorasRealesPorItem() = %v, %v; se esperaba 13.75, nil", porItem, err)
	}
	porSprint, err := repo.HorasRealesPorSprint("PRJ-017", sprintNumero)
	if err != nil || porSprint != 3.75 {
		t.Fatalf("HorasRealesPorSprint() = %v, %v; se esperaba 3.75, nil", porSprint, err)
	}
	porProyecto, err := repo.HorasRealesPorProyecto("PRJ-017")
	if err != nil || porProyecto != 13.75 {
		t.Fatalf("HorasRealesPorProyecto() = %v, %v; se esperaba 13.75, nil", porProyecto, err)
	}
}

func TestUS17_RepositorioBacklogPersisteHorasEstimadas(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-017")
	itemNumero := sembrarItemUS08Postgres(t, pool, "PRJ-017")
	repo := postgres.NuevoBacklog(pool)

	item, err := repo.ObtenerItem("PRJ-017", "HU-01")
	if err != nil {
		t.Fatalf("ObtenerItem() devolvió un error: %v", err)
	}
	if item.HorasEstimadas != nil {
		t.Fatalf("HorasEstimadas inicial = %v, se esperaba nil", item.HorasEstimadas)
	}
	estimadas := 4.25
	item.HorasEstimadas = &estimadas
	if err := repo.ActualizarItem("PRJ-017", item); err != nil {
		t.Fatalf("ActualizarItem() devolvió un error: %v", err)
	}
	obtenido, err := repo.ObtenerItem("PRJ-017", "HU-01")
	if err != nil {
		t.Fatalf("volver a obtener ítem: %v", err)
	}
	if obtenido.HorasEstimadas == nil || *obtenido.HorasEstimadas != 4.25 {
		t.Fatalf("HorasEstimadas = %v, se esperaba 4.25 (numero=%d)", obtenido.HorasEstimadas, itemNumero)
	}
}

func TestUS17_CerrarSprintGuardaSnapshotAntesDeLiberarPendientes(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-017")
	sembrarSprintUS08Postgres(t, pool, "PRJ-017", 1)
	completadoNumero := sembrarItemUS08Postgres(t, pool, "PRJ-017")
	pendienteNumero := sembrarItemUS08Postgres(t, pool, "PRJ-017")
	backlog := postgres.NuevoBacklog(pool)
	completado, err := backlog.ObtenerItem("PRJ-017", "HU-01")
	if err != nil {
		t.Fatal(err)
	}
	completado.Estado = domain.EstadoCompletada
	estimadas := 3.5
	completado.HorasEstimadas = &estimadas
	if err := backlog.ActualizarItem("PRJ-017", completado); err != nil {
		t.Fatal(err)
	}
	pendiente, err := backlog.ObtenerItem("PRJ-017", "HU-02")
	if err != nil {
		t.Fatal(err)
	}
	if err := postgres.NuevoAsignaciones(pool).Agregar("PRJ-017", 1, completadoNumero); err != nil {
		t.Fatal(err)
	}
	if err := postgres.NuevoAsignaciones(pool).Agregar("PRJ-017", 1, pendienteNumero); err != nil {
		t.Fatal(err)
	}

	if err := postgres.NuevoSprints(pool).Cerrar("PRJ-017", 1, 3); err != nil {
		t.Fatalf("Cerrar() devolvió un error: %v", err)
	}
	items, disponible, err := postgres.NuevoEsfuerzo(pool).SnapshotSprint("PRJ-017", 1)
	if err != nil {
		t.Fatalf("SnapshotSprint() devolvió un error: %v", err)
	}
	if !disponible || len(items) != 2 {
		t.Fatalf("snapshot disponible=%v items=%v; se esperaban dos ítems", disponible, items)
	}
	if items[0].ItemNumero != 1 || items[0].HorasEstimadas == nil || *items[0].HorasEstimadas != 3.5 {
		t.Errorf("snapshot del ítem completado = %+v, se esperaban 3.5 horas estimadas", items[0])
	}
	if items[1].ItemNumero != 2 || items[1].HorasEstimadas != nil {
		t.Errorf("snapshot del ítem pendiente = %+v, se esperaba estimación nil", items[1])
	}
	asignaciones := postgres.NuevoAsignaciones(pool)
	sprintNumero, asignado, err := asignaciones.SprintDeItem("PRJ-017", pendienteNumero)
	if err != nil || asignado || sprintNumero != 0 {
		t.Errorf("SprintDeItem para pendiente devuelto al backlog = (%d, %v, %v), se esperaba (0, false, nil)", sprintNumero, asignado, err)
	}
	if pendiente.Estado != domain.EstadoItemPendiente {
		t.Fatal("el cierre no debía modificar el estado de la historia pendiente")
	}
}

func TestUS17_SnapshotSprintAusenteNoSeReconstruye(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-017")
	sembrarSprintUS08Postgres(t, pool, "PRJ-017", 1)
	if _, err := pool.Exec(context.Background(), `
		UPDATE sprints SET estado = $3 WHERE proyecto_id = $1 AND numero = $2
	`, "PRJ-017", 1, domain.EstadoSprintCerrado); err != nil {
		t.Fatal(err)
	}
	_, disponible, err := postgres.NuevoEsfuerzo(pool).SnapshotSprint("PRJ-017", 1)
	if err != nil {
		t.Fatalf("SnapshotSprint() devolvió un error: %v", err)
	}
	if disponible {
		t.Fatal("un Sprint cerrado anterior al snapshot no debe reportarse disponible")
	}
}

func TestUS17_RepositorioEsfuerzoErroresNoEncontradosSePropagan(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	repo := postgres.NuevoEsfuerzo(pool)
	_, err := repo.HorasRealesPorProyecto("no-existe")
	if err != nil {
		t.Fatalf("una consulta válida sin registros debe sumar cero, error: %v", err)
	}

	_, disponible, err := repo.SnapshotSprint("no-existe", 99)
	if err != nil || disponible {
		t.Fatalf("SnapshotSprint de dato ausente = (%v, %v), se esperaba (false, nil)", disponible, err)
	}
	if errors.Is(err, pgx.ErrNoRows) {
		t.Fatal("no se esperaba pgx.ErrNoRows")
	}
}
