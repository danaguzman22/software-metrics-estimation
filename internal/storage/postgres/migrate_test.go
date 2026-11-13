package postgres_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/storage/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestMigrateOnEmptyDatabase(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL no está configurada")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	adminPool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("conectar a PostgreSQL: %v", err)
	}
	defer adminPool.Close()
	if err := adminPool.Ping(ctx); err != nil {
		t.Fatalf("verificar conexión a PostgreSQL: %v", err)
	}
	var serverVersion int
	if err := adminPool.QueryRow(ctx, "SELECT current_setting('server_version_num')::integer").Scan(&serverVersion); err != nil {
		t.Fatalf("leer versión de PostgreSQL: %v", err)
	}
	if serverVersion < 130000 {
		t.Fatalf("PostgreSQL %d no está soportado; se requiere 13 o superior", serverVersion)
	}
	t.Logf("PostgreSQL server_version_num=%d", serverVersion)

	schema := fmt.Sprintf("t12_%d", time.Now().UnixNano())
	if _, err := adminPool.Exec(ctx, "CREATE SCHEMA "+pgx.Identifier{schema}.Sanitize()); err != nil {
		t.Fatalf("crear esquema aislado: %v", err)
	}
	defer func() {
		if _, err := adminPool.Exec(context.Background(), "DROP SCHEMA "+pgx.Identifier{schema}.Sanitize()+" CASCADE"); err != nil {
			t.Errorf("eliminar esquema aislado: %v", err)
		}
	}()

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		t.Fatalf("analizar configuración PostgreSQL: %v", err)
	}
	config.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatalf("crear pool de esquema aislado: %v", err)
	}
	defer pool.Close()

	if err := postgres.Migrate(ctx, pool); err != nil {
		t.Fatalf("aplicar migraciones: %v", err)
	}
	if err := postgres.Migrate(ctx, pool); err != nil {
		t.Fatalf("reaplicar migraciones: %v", err)
	}

	assertColumnType(t, ctx, pool, schema, "proyectos", "fecha_inicio", "date")
	assertColumnType(t, ctx, pool, schema, "proyectos", "fecha_fin", "date")
	assertColumnType(t, ctx, pool, schema, "items_backlog", "criterios_aceptacion", "ARRAY")

	var projectID string
	err = pool.QueryRow(ctx, `
		INSERT INTO proyectos (nombre, fecha_inicio, fecha_fin, estado)
		VALUES ('Migración T-12', DATE '2026-01-01', DATE '2026-01-31', 'Planificado')
		RETURNING id
	`).Scan(&projectID)
	if err != nil {
		t.Fatalf("insertar proyecto con valores predeterminados: %v", err)
	}
	if projectID == "" {
		t.Fatal("el ID predeterminado del proyecto está vacío")
	}

	var generatedIDIsUUID bool
	if err := pool.QueryRow(ctx, "SELECT id::uuid IS NOT NULL FROM proyectos WHERE id = $1", projectID).Scan(&generatedIDIsUUID); err != nil {
		t.Fatalf("validar ID UUID generado: %v", err)
	}
	if !generatedIDIsUUID {
		t.Fatal("el ID predeterminado del proyecto no es UUID")
	}

	for _, status := range []string{"Planificado", "En curso", "Finalizado"} {
		if _, err := pool.Exec(ctx, `
			INSERT INTO proyectos (id, nombre, fecha_inicio, fecha_fin, estado)
			VALUES ($1, 'Estado válido', DATE '2026-01-01', DATE '2026-01-31', $2)
		`, "valid-"+status, status); err != nil {
			t.Errorf("no se permitió el estado de dominio %q: %v", status, err)
		}
	}

	if _, err := pool.Exec(ctx, `
		INSERT INTO integrantes (proyecto_id, nombre, rol)
		VALUES ($1, 'Victoria', 'Product Builder')
	`, projectID); err != nil {
		t.Fatalf("insertar integrante: %v", err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO integrantes (proyecto_id, nombre, rol)
		VALUES ($1, 'VICTORIA', 'Product Builder')
	`, projectID); err == nil {
		t.Fatal("se permitió un integrante duplicado sin distinguir mayúsculas")
	}

	if _, err := pool.Exec(ctx, `
		INSERT INTO items_backlog
			(proyecto_id, numero, titulo, descripcion, prioridad, estado, story_points)
		VALUES ($1, 1, 'HU-01', '', 'Alta', 'Pendiente', 0)
	`, projectID); err != nil {
		t.Fatalf("insertar ítem con criterios predeterminados: %v", err)
	}
	var emptyCriteria bool
	if err := pool.QueryRow(ctx, `
		SELECT criterios_aceptacion = '{}'::text[]
		FROM items_backlog WHERE proyecto_id = $1 AND numero = 1
	`, projectID).Scan(&emptyCriteria); err != nil {
		t.Fatalf("leer criterios predeterminados: %v", err)
	}
	if !emptyCriteria {
		t.Fatal("criterios_aceptacion no tiene una lista vacía por defecto")
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO items_backlog
			(proyecto_id, numero, titulo, descripcion, prioridad, estado, story_points)
		VALUES ($1, 1, 'HU-01 duplicado', '', 'Alta', 'Pendiente', 0)
	`, projectID); err == nil {
		t.Fatal("se permitió duplicar la PK compuesta de items_backlog")
	}

	for _, invalidStatus := range []string{"Pendiente", "En progreso", "Cerrado"} {
		if _, err := pool.Exec(ctx, `
			INSERT INTO proyectos (id, nombre, fecha_inicio, fecha_fin, estado)
			VALUES ($1, 'Estado inválido', DATE '2026-01-01', DATE '2026-01-31', $2)
		`, "invalid-"+invalidStatus, invalidStatus); err == nil {
			t.Errorf("se permitió estado de proyecto %q", invalidStatus)
		}
	}
}

func assertColumnType(t *testing.T, ctx context.Context, pool *pgxpool.Pool, schema, table, column, want string) {
	t.Helper()
	var got string
	err := pool.QueryRow(ctx, `
		SELECT data_type
		FROM information_schema.columns
		WHERE table_schema = $1 AND table_name = $2 AND column_name = $3
	`, schema, table, column).Scan(&got)
	if err != nil {
		t.Fatalf("leer tipo de columna %s.%s: %v", table, column, err)
	}
	if got != want {
		t.Errorf("tipo de %s.%s = %q, se esperaba %q", table, column, got, want)
	}
}
