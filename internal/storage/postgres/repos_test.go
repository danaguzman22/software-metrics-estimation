package postgres_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
	"github.com/danaguzman22/software-metrics-estimation/internal/service/repotest"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresRepo_Contrato(t *testing.T) {
	repotest.Ejecutar(t, nuevaEntornoPostgres)
}

func TestPostgresRepo_IntegridadReferencialYChecks(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")

	if _, err := pool.Exec(context.Background(), `
		INSERT INTO integrantes (proyecto_id, nombre, rol)
		VALUES ('inexistente', 'Ana', 'Dev')
	`); err == nil {
		t.Fatal("se permitió insertar integrante con proyecto inexistente")
	}
	if _, err := pool.Exec(context.Background(), `
		INSERT INTO items_backlog
			(proyecto_id, numero, titulo, descripcion, prioridad, estado, story_points)
		VALUES ('inexistente', 1, 'HU-01', '', 'Alta', 'Pendiente', 0)
	`); err == nil {
		t.Fatal("se permitió insertar ítem con proyecto inexistente")
	}

	for _, invalid := range []struct {
		name string
		sql  string
	}{
		{
			name: "prioridad",
			sql: `INSERT INTO items_backlog
				(proyecto_id, numero, titulo, descripcion, prioridad, estado, story_points)
				VALUES ('PRJ-001', 1, 'HU-01', '', 'Urgente', 'Pendiente', 0)`,
		},
		{
			name: "estado",
			sql: `INSERT INTO items_backlog
				(proyecto_id, numero, titulo, descripcion, prioridad, estado, story_points)
				VALUES ('PRJ-001', 1, 'HU-01', '', 'Alta', 'En curso', 0)`,
		},
		{
			name: "story_points",
			sql: `INSERT INTO items_backlog
				(proyecto_id, numero, titulo, descripcion, prioridad, estado, story_points)
				VALUES ('PRJ-001', 1, 'HU-01', '', 'Alta', 'Pendiente', 4)`,
		},
	} {
		t.Run(invalid.name, func(t *testing.T) {
			if _, err := pool.Exec(context.Background(), invalid.sql); err == nil {
				t.Fatalf("se permitió un %s fuera del CHECK", invalid.name)
			}
		})
	}
}

func TestPostgresRepo_DuplicadoIntegranteTraduceError(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	repo := postgres.NuevoIntegrantes(pool)
	integrante := domain.Integrante{Nombre: "Ana", Rol: "Dev"}

	if err := repo.Agregar("PRJ-001", integrante); err != nil {
		t.Fatalf("agregar primer integrante: %v", err)
	}
	err := repo.Agregar("PRJ-001", domain.Integrante{Nombre: "ANA", Rol: "QA"})
	if err != domain.ErrIntegranteDuplicado {
		t.Fatalf("error de duplicado = %v, se esperaba exactamente %v", err, domain.ErrIntegranteDuplicado)
	}
}

func TestPostgresRepo_UnicidadIntegrantesYClaveBacklog(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")

	if _, err := pool.Exec(context.Background(), `
		INSERT INTO integrantes (proyecto_id, nombre, rol)
		VALUES ('PRJ-001', 'Ana', 'Dev')
	`); err != nil {
		t.Fatalf("insertar integrante: %v", err)
	}
	if _, err := pool.Exec(context.Background(), `
		INSERT INTO integrantes (proyecto_id, nombre, rol)
		VALUES ('PRJ-001', 'ANA', 'QA')
	`); err == nil {
		t.Fatal("se permitió integrante duplicado ignorando mayúsculas")
	}

	insertItem := `INSERT INTO items_backlog
		(proyecto_id, numero, titulo, descripcion, prioridad, estado, story_points)
		VALUES ('PRJ-001', 1, 'HU-01', '', 'Alta', 'Pendiente', 0)`
	if _, err := pool.Exec(context.Background(), insertItem); err != nil {
		t.Fatalf("insertar primer ítem: %v", err)
	}
	if _, err := pool.Exec(context.Background(), insertItem); err == nil {
		t.Fatal("se permitió duplicar la PK compuesta de backlog")
	}
}

func TestPostgresRepo_AgregarItemConcurrente(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	repo := postgres.NuevoBacklog(pool)

	const cantidad = 20
	ids := make(chan string, cantidad)
	errs := make(chan error, cantidad)
	var wg sync.WaitGroup
	for i := 0; i < cantidad; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			item, err := domain.NuevoItemBacklog(fmt.Sprintf("Item %d", i), "", "Media", nil)
			if err != nil {
				errs <- err
				return
			}
			agregado, err := repo.AgregarItem("PRJ-001", item)
			if err != nil {
				errs <- err
				return
			}
			ids <- agregado.ID
		}(i)
	}
	wg.Wait()
	close(ids)
	close(errs)

	for err := range errs {
		t.Errorf("AgregarItem concurrente: %v", err)
	}
	gotIDs := make(map[string]bool, cantidad)
	for id := range ids {
		if gotIDs[id] {
			t.Errorf("ID duplicado: %s", id)
		}
		gotIDs[id] = true
	}
	if len(gotIDs) != cantidad {
		t.Fatalf("IDs únicos: esperado %d, obtenido %d", cantidad, len(gotIDs))
	}
	for i := 1; i <= cantidad; i++ {
		id := fmt.Sprintf("HU-%02d", i)
		if !gotIDs[id] {
			t.Errorf("falta ID consecutivo %s", id)
		}
	}

	items, err := repo.ListarPorProyecto("PRJ-001")
	if err != nil {
		t.Fatalf("listar ítems concurrentes: %v", err)
	}
	if len(items) != cantidad {
		t.Fatalf("ítems persistidos: esperado %d, obtenido %d", cantidad, len(items))
	}
}

func nuevaEntornoPostgres(t *testing.T) repotest.Entorno {
	t.Helper()
	pool := nuevoPoolPostgres(t)
	return repotest.Entorno{
		Proyectos:   postgres.NuevoProyectos(pool),
		Integrantes: postgres.NuevoIntegrantes(pool),
		Backlog:     postgres.NuevoBacklog(pool),
		SembrarProyecto: func(t *testing.T, proyecto service.ProyectoRegistrado) {
			t.Helper()
			_, err := pool.Exec(context.Background(), `
				INSERT INTO proyectos (id, nombre, fecha_inicio, fecha_fin, estado)
				VALUES ($1, $2, $3::date, $4::date, $5)
			`, proyecto.ID, proyecto.Proyecto.Nombre,
				proyecto.Proyecto.FechaInicio.Format("2006-01-02"),
				proyecto.Proyecto.FechaFin.Format("2006-01-02"),
				proyecto.Proyecto.Estado)
			if err != nil {
				t.Fatalf("sembrar proyecto %s: %v", proyecto.ID, err)
			}
		},
	}
}

func nuevoPoolPostgres(t *testing.T) *pgxpool.Pool {
	t.Helper()
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL no está configurada")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	t.Cleanup(cancel)
	adminPool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("conectar a PostgreSQL: %v", err)
	}
	if err := adminPool.Ping(ctx); err != nil {
		adminPool.Close()
		t.Fatalf("verificar conexión a PostgreSQL: %v", err)
	}

	schema := fmt.Sprintf("t13_%d", time.Now().UnixNano())
	if _, err := adminPool.Exec(ctx, "CREATE SCHEMA "+pgx.Identifier{schema}.Sanitize()); err != nil {
		adminPool.Close()
		t.Fatalf("crear esquema aislado: %v", err)
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		_, _ = adminPool.Exec(ctx, "DROP SCHEMA "+pgx.Identifier{schema}.Sanitize()+" CASCADE")
		adminPool.Close()
		t.Fatalf("analizar configuración PostgreSQL: %v", err)
	}
	config.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		_, _ = adminPool.Exec(ctx, "DROP SCHEMA "+pgx.Identifier{schema}.Sanitize()+" CASCADE")
		adminPool.Close()
		t.Fatalf("crear pool de esquema aislado: %v", err)
	}
	if err := postgres.Migrate(ctx, pool); err != nil {
		pool.Close()
		_, _ = adminPool.Exec(ctx, "DROP SCHEMA "+pgx.Identifier{schema}.Sanitize()+" CASCADE")
		adminPool.Close()
		t.Fatalf("aplicar migraciones al esquema aislado: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
		if _, err := adminPool.Exec(context.Background(), "DROP SCHEMA "+pgx.Identifier{schema}.Sanitize()+" CASCADE"); err != nil {
			t.Errorf("eliminar esquema aislado: %v", err)
		}
		adminPool.Close()
	})
	return pool
}

func seedProject(t *testing.T, pool *pgxpool.Pool, id string) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `
		INSERT INTO proyectos (id, nombre, fecha_inicio, fecha_fin, estado)
		VALUES ($1, 'Proyecto de integración', DATE '2026-01-01', DATE '2026-12-31', 'Planificado')
	`, id)
	if err != nil {
		t.Fatalf("sembrar proyecto %s: %v", id, err)
	}
}
