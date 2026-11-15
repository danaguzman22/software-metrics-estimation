package postgres_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage/postgres"
)

func TestUS15_RepositorioDefectosCrearListarYActualizar(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-015")
	sembrarSprintUS08Postgres(t, pool, "PRJ-015", 1)
	sprintCerrado := sprintUS08Postgres(2)
	sprintCerrado.Estado = domain.EstadoSprintCerrado
	if err := postgres.NuevoSprints(pool).Agregar("PRJ-015", sprintCerrado); err != nil {
		t.Fatalf("sembrar Sprint cerrado de integración: %v", err)
	}
	itemNumero := sembrarItemUS08Postgres(t, pool, "PRJ-015")
	repo := postgres.NuevoDefectos(pool)

	defecto := domain.Defecto{
		Descripcion:     "No persiste el cambio",
		Severidad:       domain.SeveridadDefectoAlta,
		Estado:          domain.EstadoDefectoAbierto,
		ItemID:          fmt.Sprintf("HU-%02d", itemNumero),
		SprintDeteccion: 1,
	}
	creado, err := repo.Crear("PRJ-015", defecto)
	if err != nil {
		t.Fatalf("Crear() devolvió error: %v", err)
	}
	if creado.ID <= 0 || creado.Descripcion != defecto.Descripcion || creado.ItemID != defecto.ItemID {
		t.Fatalf("defecto creado = %+v", creado)
	}

	lista, err := repo.ListarPorProyecto("PRJ-015")
	if err != nil || len(lista) != 1 || lista[0].ID != creado.ID {
		t.Fatalf("lista = %+v, err=%v", lista, err)
	}
	obtenido, err := repo.Obtener("PRJ-015", creado.ID)
	if err != nil || obtenido.ID != creado.ID {
		t.Fatalf("Obtener() = %+v, err=%v", obtenido, err)
	}

	sprintResolucion := 2
	obtenido.Estado = domain.EstadoDefectoResuelto
	obtenido.SprintResolucion = &sprintResolucion
	if err := repo.Actualizar("PRJ-015", obtenido); err != nil {
		t.Fatalf("Actualizar() devolvió error: %v", err)
	}
	actualizado, err := repo.Obtener("PRJ-015", creado.ID)
	if err != nil || actualizado.Estado != domain.EstadoDefectoResuelto ||
		actualizado.SprintResolucion == nil || *actualizado.SprintResolucion != 2 {
		t.Fatalf("defecto actualizado = %+v, err=%v", actualizado, err)
	}
}

func TestUS15_RepositorioDefectosValidaProyectoEIntegridad(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-015")
	repo := postgres.NuevoDefectos(pool)
	defecto := domain.Defecto{
		Descripcion:     "Error",
		Severidad:       domain.SeveridadDefectoMedia,
		Estado:          domain.EstadoDefectoAbierto,
		ItemID:          "HU-01",
		SprintDeteccion: 1,
	}

	if _, err := repo.Crear("PRJ-INEXISTENTE", defecto); err != domain.ErrProyectoNoEncontrado {
		t.Fatalf("error de proyecto = %v; se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if _, err := repo.Crear("PRJ-015", defecto); !errors.Is(err, domain.ErrItemNoEncontrado) {
		t.Fatalf("error de historia = %v; se esperaba %v", err, domain.ErrItemNoEncontrado)
	}

	sembrarSprintUS08Postgres(t, pool, "PRJ-015", 1)
	itemNumero := sembrarItemUS08Postgres(t, pool, "PRJ-015")
	defecto.ItemID = fmt.Sprintf("HU-%02d", itemNumero)
	defecto.SprintDeteccion = 99
	if _, err := repo.Crear("PRJ-015", defecto); !errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Fatalf("error de Sprint = %v; se esperaba %v", err, domain.ErrSprintNoEncontrado)
	}

	// La búsqueda por proyecto debe omitir IDs pertenecientes a otros proyectos.
	seedProject(t, pool, "PRJ-016")
	sembrarSprintUS08Postgres(t, pool, "PRJ-016", 1)
	itemOtroProyecto := sembrarItemUS08Postgres(t, pool, "PRJ-016")
	otro := defecto
	otro.ItemID = fmt.Sprintf("HU-%02d", itemOtroProyecto)
	otro.SprintDeteccion = 1
	if _, err := repo.Crear("PRJ-016", otro); err != nil {
		t.Fatalf("crear defecto del segundo proyecto: %v", err)
	}
	if _, err := repo.Obtener("PRJ-015", 1); !errors.Is(err, domain.ErrDefectoNoEncontrado) {
		t.Fatalf("Obtener() cruzando proyectos = %v; se esperaba %v", err, domain.ErrDefectoNoEncontrado)
	}

	var conteo int
	if err := pool.QueryRow(context.Background(), `SELECT count(*) FROM defectos WHERE proyecto_id = 'PRJ-016'`).Scan(&conteo); err != nil {
		t.Fatal(err)
	}
	if conteo != 1 {
		t.Fatalf("cantidad de defectos del segundo proyecto = %d, se esperaba 1", conteo)
	}
}
