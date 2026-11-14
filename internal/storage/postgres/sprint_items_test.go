package postgres_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

// sprintUS08Postgres arma un Sprint activo de referencia para la suite de
// integración de US-08.
func sprintUS08Postgres(numero int) domain.Sprint {
	return domain.Sprint{
		Numero:      numero,
		Objetivo:    "Planificar el Sprint Backlog",
		FechaInicio: time.Date(2026, time.March, 2, 0, 0, 0, 0, time.UTC),
		FechaFin:    time.Date(2026, time.March, 20, 0, 0, 0, 0, time.UTC),
		Estado:      domain.EstadoSprintActivo,
	}
}

func sembrarSprintUS08Postgres(t *testing.T, pool *pgxpool.Pool, proyectoID string, numero int) {
	t.Helper()
	if err := postgres.NuevoSprints(pool).Agregar(proyectoID, sprintUS08Postgres(numero)); err != nil {
		t.Fatalf("sembrar Sprint %d del proyecto %s: %v", numero, proyectoID, err)
	}
}

// sembrarItemUS08Postgres crea una historia estimada con el adaptador real de
// backlog y devuelve su número.
func sembrarItemUS08Postgres(t *testing.T, pool *pgxpool.Pool, proyectoID string) int {
	t.Helper()
	item, err := domain.NuevoItemBacklog("Historia de integración", "Historia para US-08", "Alta", nil)
	if err != nil {
		t.Fatalf("crear ítem de integración: %v", err)
	}
	item, err = item.Estimar(3)
	if err != nil {
		t.Fatalf("estimar ítem de integración: %v", err)
	}
	agregado, err := postgres.NuevoBacklog(pool).AgregarItem(proyectoID, item)
	if err != nil {
		t.Fatalf("sembrar ítem del proyecto %s: %v", proyectoID, err)
	}
	numero, ok := domain.NumeroDeItem(agregado.ID)
	if !ok {
		t.Fatalf("el ID sembrado %q no es válido", agregado.ID)
	}
	return numero
}

// nuevoRepositorioAsignacionesUS08 instancia el adaptador PostgreSQL de la
// relación historia–Sprint. Es el único punto que depende del nombre del
// constructor del adaptador.
func nuevoRepositorioAsignacionesUS08(pool *pgxpool.Pool) service.RepositorioAsignaciones {
	return postgres.NuevoAsignaciones(pool)
}

// verificarErrorDominioUS08 comprueba que el adaptador devuelva el error de
// dominio exacto, sin envolver.
func verificarErrorDominioUS08(t *testing.T, err, esperado error) {
	t.Helper()
	if !errors.Is(err, esperado) {
		t.Fatalf("err = %v, se esperaba %v", err, esperado)
	}
	if err != esperado {
		t.Errorf("el error de dominio debe devolverse sin envolver: err = %v", err)
	}
}

// CA3 / R9: asignar y consultar la existencia del ítem en el proyecto.
func TestUS08_RepositorioAsignaciones_AsignarYExiste(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	sembrarSprintUS08Postgres(t, pool, "PRJ-001", 1)
	numero := sembrarItemUS08Postgres(t, pool, "PRJ-001")
	repo := nuevoRepositorioAsignacionesUS08(pool)

	if err := repo.Agregar("PRJ-001", 1, numero); err != nil {
		t.Fatalf("Agregar() devolvió un error inesperado: %v", err)
	}

	existe, err := repo.ExisteItem("PRJ-001", numero)
	if err != nil {
		t.Fatalf("ExisteItem() devolvió un error inesperado: %v", err)
	}
	if !existe {
		t.Errorf("ExisteItem() = false, se esperaba true para el ítem asignado")
	}

	existe, err = repo.ExisteItem("PRJ-001", numero+1)
	if err != nil {
		t.Fatalf("ExisteItem() de un ítem no asignado devolvió un error: %v", err)
	}
	if existe {
		t.Errorf("ExisteItem() = true para un ítem no asignado")
	}

	numeros, err := repo.ListarNumerosItem("PRJ-001", 1)
	if err != nil {
		t.Fatalf("ListarNumerosItem() devolvió un error inesperado: %v", err)
	}
	if !reflect.DeepEqual(numeros, []int{numero}) {
		t.Errorf("ListarNumerosItem() = %v, se esperaba [%d]", numeros, numero)
	}
}

// CA2 / R4 / D7: la PK (proyecto_id, item_numero) impide asignar dos veces la
// misma historia y el adaptador traduce el error de unicidad.
func TestUS08_RepositorioAsignaciones_ItemYaAsignado(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	sembrarSprintUS08Postgres(t, pool, "PRJ-001", 1)
	numero := sembrarItemUS08Postgres(t, pool, "PRJ-001")
	repo := nuevoRepositorioAsignacionesUS08(pool)

	if err := repo.Agregar("PRJ-001", 1, numero); err != nil {
		t.Fatalf("Agregar() del primer ítem devolvió un error inesperado: %v", err)
	}

	err := repo.Agregar("PRJ-001", 1, numero)
	verificarErrorDominioUS08(t, err, domain.ErrItemYaAsignado)

	numeros, err := repo.ListarNumerosItem("PRJ-001", 1)
	if err != nil {
		t.Fatalf("ListarNumerosItem() devolvió un error inesperado: %v", err)
	}
	if !reflect.DeepEqual(numeros, []int{numero}) {
		t.Errorf("se esperaba una sola fila con el ítem %d, se obtuvo %v", numero, numeros)
	}
}

// CA2 / R5: el mismo número de historia en otro proyecto no interfiere porque
// la PK es por proyecto.
func TestUS08_RepositorioAsignaciones_MismoItemOtroProyecto(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	seedProject(t, pool, "PRJ-002")
	sembrarSprintUS08Postgres(t, pool, "PRJ-001", 1)
	sembrarSprintUS08Postgres(t, pool, "PRJ-002", 1)
	numeroP1 := sembrarItemUS08Postgres(t, pool, "PRJ-001")
	numeroP2 := sembrarItemUS08Postgres(t, pool, "PRJ-002")
	repo := nuevoRepositorioAsignacionesUS08(pool)

	if err := repo.Agregar("PRJ-001", 1, numeroP1); err != nil {
		t.Fatalf("Agregar() en PRJ-001 devolvió un error inesperado: %v", err)
	}
	if err := repo.Agregar("PRJ-002", 1, numeroP2); err != nil {
		t.Fatalf("Agregar() en PRJ-002 devolvió un error inesperado: %v", err)
	}

	numerosP1, err := repo.ListarNumerosItem("PRJ-001", 1)
	if err != nil {
		t.Fatalf("ListarNumerosItem() de PRJ-001 devolvió un error inesperado: %v", err)
	}
	if !reflect.DeepEqual(numerosP1, []int{numeroP1}) {
		t.Errorf("PRJ-001 = %v, se esperaba [%d]", numerosP1, numeroP1)
	}
	numerosP2, err := repo.ListarNumerosItem("PRJ-002", 1)
	if err != nil {
		t.Fatalf("ListarNumerosItem() de PRJ-002 devolvió un error inesperado: %v", err)
	}
	if !reflect.DeepEqual(numerosP2, []int{numeroP2}) {
		t.Errorf("PRJ-002 = %v, se esperaba [%d]", numerosP2, numeroP2)
	}
}

// CA1 / R7 / D7: la FK de Sprint se traduce a domain.ErrSprintNoEncontrado sin
// envolver y no se inserta la fila.
func TestUS08_RepositorioAsignaciones_SprintInexistente(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	numero := sembrarItemUS08Postgres(t, pool, "PRJ-001")
	repo := nuevoRepositorioAsignacionesUS08(pool)

	err := repo.Agregar("PRJ-001", 1, numero)
	verificarErrorDominioUS08(t, err, domain.ErrSprintNoEncontrado)

	numeros, err := repo.ListarNumerosItem("PRJ-001", 1)
	if err != nil {
		t.Fatalf("ListarNumerosItem() devolvió un error inesperado: %v", err)
	}
	if len(numeros) != 0 {
		t.Errorf("no debía persistirse la asignación, se obtuvo %v", numeros)
	}
}

// CA1 / R8 / D7: la FK de ítem se traduce a domain.ErrItemNoEncontrado sin
// envolver y no se inserta la fila.
func TestUS08_RepositorioAsignaciones_ItemInexistente(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	sembrarSprintUS08Postgres(t, pool, "PRJ-001", 1)
	repo := nuevoRepositorioAsignacionesUS08(pool)

	err := repo.Agregar("PRJ-001", 1, 1)
	verificarErrorDominioUS08(t, err, domain.ErrItemNoEncontrado)

	numeros, err := repo.ListarNumerosItem("PRJ-001", 1)
	if err != nil {
		t.Fatalf("ListarNumerosItem() devolvió un error inesperado: %v", err)
	}
	if len(numeros) != 0 {
		t.Errorf("no debía persistirse la asignación, se obtuvo %v", numeros)
	}
}

// CA3 / R9: ListarNumerosItem devuelve los números ordenados de forma
// ascendente y una lista vacía para un Sprint sin historias.
func TestUS08_RepositorioAsignaciones_ListarNumerosItem(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	sembrarSprintUS08Postgres(t, pool, "PRJ-001", 1)
	numeros := make([]int, 0, 3)
	for i := 0; i < 3; i++ {
		numeros = append(numeros, sembrarItemUS08Postgres(t, pool, "PRJ-001"))
	}
	repo := nuevoRepositorioAsignacionesUS08(pool)

	// Se asignan en orden inverso para comprobar que el adaptador ordena.
	for _, numero := range []int{numeros[2], numeros[0], numeros[1]} {
		if err := repo.Agregar("PRJ-001", 1, numero); err != nil {
			t.Fatalf("Agregar(%d) devolvió un error inesperado: %v", numero, err)
		}
	}

	obtenidos, err := repo.ListarNumerosItem("PRJ-001", 1)
	if err != nil {
		t.Fatalf("ListarNumerosItem() devolvió un error inesperado: %v", err)
	}
	if !reflect.DeepEqual(obtenidos, numeros) {
		t.Errorf("ListarNumerosItem() = %v, se esperaba %v", obtenidos, numeros)
	}

	seedProject(t, pool, "PRJ-002")
	sembrarSprintUS08Postgres(t, pool, "PRJ-002", 1)
	vacios, err := repo.ListarNumerosItem("PRJ-002", 1)
	if err != nil {
		t.Fatalf("ListarNumerosItem() de un Sprint vacío devolvió un error: %v", err)
	}
	if len(vacios) != 0 {
		t.Errorf("ListarNumerosItem() de un Sprint vacío = %v, se esperaba una lista vacía", vacios)
	}
}
