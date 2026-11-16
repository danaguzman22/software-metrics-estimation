package postgres_test

import (
	"context"
	"errors"
	"fmt"
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

func TestUS09_RepositorioSprints_ObtenerYListarIncluyenPlanificados(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	repo := postgres.NuevoSprints(pool)

	if err := repo.Agregar("PRJ-001", sprintUS07Postgres(1)); err != nil {
		t.Fatalf("agregar Sprint 1: %v", err)
	}
	if _, err := pool.Exec(context.Background(), `
		UPDATE sprints SET estado = 'Cerrado', sp_planificados = 8
		WHERE proyecto_id = 'PRJ-001' AND numero = 1
	`); err != nil {
		t.Fatalf("cerrar Sprint sembrado: %v", err)
	}
	if err := repo.Agregar("PRJ-001", sprintUS07Postgres(2)); err != nil {
		t.Fatalf("agregar Sprint 2: %v", err)
	}

	obtenido, err := repo.Obtener("PRJ-001", 1)
	if err != nil {
		t.Fatalf("Obtener() devolvió un error inesperado: %v", err)
	}
	if obtenido.SPPlanificados == nil || *obtenido.SPPlanificados != 8 {
		t.Errorf("SPPlanificados = %v, se esperaba 8", obtenido.SPPlanificados)
	}

	listados, err := repo.ListarPorProyecto("PRJ-001")
	if err != nil {
		t.Fatalf("ListarPorProyecto() devolvió un error inesperado: %v", err)
	}
	if len(listados) != 2 || listados[0].Numero != 1 || listados[1].Numero != 2 {
		t.Fatalf("ListarPorProyecto() = %+v, se esperaban los Sprints 1 y 2 ordenados", listados)
	}
	if listados[0].SPPlanificados == nil || *listados[0].SPPlanificados != 8 {
		t.Errorf("Sprint listado SPPlanificados = %v, se esperaba 8", listados[0].SPPlanificados)
	}
	if listados[1].SPPlanificados != nil {
		t.Errorf("Sprint activo SPPlanificados = %v, se esperaba nil", listados[1].SPPlanificados)
	}
}

func TestUS09_RepositorioSprints_CerrarAtomicoYLiberaIndiceActivo(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	sprints := postgres.NuevoSprints(pool)
	asignaciones := postgres.NuevoAsignaciones(pool)
	backlog := postgres.NuevoBacklog(pool)

	if err := sprints.Agregar("PRJ-001", sprintUS07Postgres(1)); err != nil {
		t.Fatalf("agregar Sprint 1: %v", err)
	}
	completadoNumero := sembrarItemUS08Postgres(t, pool, "PRJ-001")
	pendienteNumero := sembrarItemUS08Postgres(t, pool, "PRJ-001")
	completado, err := backlog.ObtenerItem("PRJ-001", fmt.Sprintf("HU-%02d", completadoNumero))
	if err != nil {
		t.Fatalf("obtener historia completada: %v", err)
	}
	completado.Estado = domain.EstadoCompletada
	if err := backlog.ActualizarItem("PRJ-001", completado); err != nil {
		t.Fatalf("marcar historia completada: %v", err)
	}
	for _, numero := range []int{completadoNumero, pendienteNumero} {
		if err := asignaciones.Agregar("PRJ-001", 1, numero); err != nil {
			t.Fatalf("asignar historia %d: %v", numero, err)
		}
	}

	if err := sprints.Cerrar("PRJ-001", 1, 6); err != nil {
		t.Fatalf("Cerrar() devolvió un error inesperado: %v", err)
	}
	obtenido, err := sprints.Obtener("PRJ-001", 1)
	if err != nil {
		t.Fatalf("Obtener() tras cerrar: %v", err)
	}
	if obtenido.Estado != "Cerrado" || obtenido.SPPlanificados == nil || *obtenido.SPPlanificados != 6 {
		t.Errorf("Sprint tras cerrar = %+v, se esperaba Cerrado con 6 SP planificados", obtenido)
	}
	numeros, err := asignaciones.ListarNumerosItem("PRJ-001", 1)
	if err != nil {
		t.Fatalf("ListarNumerosItem() tras cerrar: %v", err)
	}
	if len(numeros) != 1 || numeros[0] != completadoNumero {
		t.Errorf("asignaciones tras cerrar = %v, se esperaba solo la historia completada %d", numeros, completadoNumero)
	}
	if err := sprints.Agregar("PRJ-001", sprintUS07Postgres(2)); err != nil {
		t.Fatalf("crear siguiente Sprint activo tras el cierre: %v", err)
	}
}

func TestUS09_RepositorioSprints_CerrarSprintNoActivoDevuelveError(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	repo := postgres.NuevoSprints(pool)
	if err := repo.Agregar("PRJ-001", sprintUS07Postgres(1)); err != nil {
		t.Fatalf("agregar Sprint: %v", err)
	}

	if err := repo.Cerrar("PRJ-001", 1, 3); err != nil {
		t.Fatalf("primer Cerrar() devolvió un error inesperado: %v", err)
	}
	if err := repo.Cerrar("PRJ-001", 1, 9); err != domain.ErrSprintCerrado {
		t.Fatalf("segundo Cerrar() = %v, se esperaba %v", err, domain.ErrSprintCerrado)
	}
	obtenido, err := repo.Obtener("PRJ-001", 1)
	if err != nil {
		t.Fatalf("Obtener() tras el segundo cierre: %v", err)
	}
	if obtenido.SPPlanificados == nil || *obtenido.SPPlanificados != 3 {
		t.Errorf("el segundo cierre modificó sp_planificados: %v", obtenido.SPPlanificados)
	}
}
