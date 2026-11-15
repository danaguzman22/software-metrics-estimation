package postgres_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage/postgres"
	"github.com/danaguzman22/software-metrics-estimation/migrations"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// itemIDUS12 siembra una historia en el proyecto indicado y devuelve su ID.
func itemIDUS12(t *testing.T, pool *pgxpool.Pool, proyectoID string) string {
	t.Helper()
	numero := sembrarItemUS08Postgres(t, pool, proyectoID)
	return fmt.Sprintf("HU-%02d", numero)
}

// nuevoPoolVacioPostgres arma un pool sobre un esquema aislado sin aplicar
// ninguna migración; lo usa el test de Up/Down de 00006.
func nuevoPoolVacioPostgres(t *testing.T) (*pgxpool.Pool, string) {
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

	schema := fmt.Sprintf("t14_%d", time.Now().UnixNano())
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

	t.Cleanup(func() {
		pool.Close()
		if _, err := adminPool.Exec(context.Background(), "DROP SCHEMA "+pgx.Identifier{schema}.Sanitize()+" CASCADE"); err != nil {
			t.Errorf("eliminar esquema aislado: %v", err)
		}
		adminPool.Close()
	})
	return pool, schema
}

func existeTablaUS12(t *testing.T, pool *pgxpool.Pool, schema, tabla string) bool {
	t.Helper()
	var existe bool
	if err := pool.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables
			WHERE table_schema = $1 AND table_name = $2
		)
	`, schema, tabla).Scan(&existe); err != nil {
		t.Fatalf("consultar existencia de la tabla %s.%s: %v", schema, tabla, err)
	}
	return existe
}

// CA1 / R6 / R9: round-trip de Iniciar, Obtener y ListarPorHistoria, con las
// rondas ordenadas por número.
func TestUS12_Postgres_IniciarObtenerYListarPorHistoria(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	itemID := itemIDUS12(t, pool, "PRJ-001")
	repo := postgres.NuevoRondasPoker(pool)

	if err := repo.Iniciar("PRJ-001", itemID, 1); err != nil {
		t.Fatalf("Iniciar() devolvió un error inesperado: %v", err)
	}

	ronda, err := repo.Obtener("PRJ-001", itemID, 1)
	if err != nil {
		t.Fatalf("Obtener() devolvió un error inesperado: %v", err)
	}
	if ronda.Numero != 1 || ronda.Estado != domain.EstadoRondaAbierta {
		t.Errorf("Obtener() = %+v, se esperaba ronda 1 en estado Abierta", ronda)
	}

	if err := repo.Revelar("PRJ-001", itemID, 1); err != nil {
		t.Fatalf("Revelar() devolvió un error inesperado: %v", err)
	}
	if err := repo.Iniciar("PRJ-001", itemID, 2); err != nil {
		t.Fatalf("Iniciar() de la segunda ronda devolvió un error inesperado: %v", err)
	}

	rondas, err := repo.ListarPorHistoria("PRJ-001", itemID)
	if err != nil {
		t.Fatalf("ListarPorHistoria() devolvió un error inesperado: %v", err)
	}
	if len(rondas) != 2 || rondas[0].Numero != 1 || rondas[1].Numero != 2 {
		t.Fatalf("ListarPorHistoria() = %+v, se esperaban las rondas 1 y 2 ordenadas", rondas)
	}
	if rondas[0].Estado != domain.EstadoRondaRevelada {
		t.Errorf("ronda 1 = %q, se esperaba Revelada", rondas[0].Estado)
	}
	if rondas[1].Estado != domain.EstadoRondaAbierta {
		t.Errorf("ronda 2 = %q, se esperaba Abierta", rondas[1].Estado)
	}

	ronda, err = repo.Obtener("PRJ-001", itemID, 1)
	if err != nil {
		t.Fatalf("Obtener() tras revelar devolvió un error inesperado: %v", err)
	}
	if ronda.Estado != domain.EstadoRondaRevelada {
		t.Errorf("Obtener() tras revelar = %q, se esperaba Revelada", ronda.Estado)
	}
}

// R2 / R3 / R12: el upsert reemplaza el valor del voto del integrante y deja
// una única fila.
func TestUS12_Postgres_RegistrarVotoUpsertReemplazaValor(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	itemID := itemIDUS12(t, pool, "PRJ-001")
	repo := postgres.NuevoRondasPoker(pool)

	if err := repo.Iniciar("PRJ-001", itemID, 1); err != nil {
		t.Fatalf("Iniciar() devolvió un error inesperado: %v", err)
	}
	for _, valor := range []int{5, 8} {
		if err := repo.RegistrarVoto("PRJ-001", itemID, 1, domain.VotoPoker{Integrante: "Ana", Valor: valor}); err != nil {
			t.Fatalf("RegistrarVoto(%d) devolvió un error inesperado: %v", valor, err)
		}
	}

	if err := repo.Revelar("PRJ-001", itemID, 1); err != nil {
		t.Fatalf("Revelar() devolvió un error inesperado: %v", err)
	}
	votos, err := repo.Votos("PRJ-001", itemID, 1)
	if err != nil {
		t.Fatalf("Votos() devolvió un error inesperado: %v", err)
	}
	if len(votos) != 1 {
		t.Fatalf("Votos() = %+v, se esperaba una única fila", votos)
	}
	if votos[0].Integrante != "Ana" || votos[0].Valor != 8 {
		t.Errorf("voto = %+v, se esperaba Ana con 8 (el último valor escrito gana)", votos[0])
	}
}

// R6 / D20: YaVotaron y Votos ordenan alfabéticamente sin distinguir
// mayúsculas; YaVotaron nunca devuelve valores.
func TestUS12_Postgres_YaVotaronYVotosOrdenAlfabeticoSinDistinguirMayusculas(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	itemID := itemIDUS12(t, pool, "PRJ-001")
	repo := postgres.NuevoRondasPoker(pool)

	if err := repo.Iniciar("PRJ-001", itemID, 1); err != nil {
		t.Fatalf("Iniciar() devolvió un error inesperado: %v", err)
	}
	votos := []domain.VotoPoker{
		{Integrante: "Zoe", Valor: 2},
		{Integrante: "ana", Valor: 5},
		{Integrante: "Beto", Valor: 8},
		{Integrante: "beto", Valor: 13},
		{Integrante: "CARLA", Valor: 21},
	}
	for _, voto := range votos {
		if err := repo.RegistrarVoto("PRJ-001", itemID, 1, voto); err != nil {
			t.Fatalf("RegistrarVoto(%q) devolvió un error inesperado: %v", voto.Integrante, err)
		}
	}

	nombres, err := repo.YaVotaron("PRJ-001", itemID, 1)
	if err != nil {
		t.Fatalf("YaVotaron() devolvió un error inesperado: %v", err)
	}
	esperados := []string{"ana", "Beto", "beto", "CARLA", "Zoe"}
	if len(nombres) != len(esperados) {
		t.Fatalf("YaVotaron() = %v, se esperaban %v", nombres, esperados)
	}
	for i, nombre := range esperados {
		if nombres[i] != nombre {
			t.Errorf("YaVotaron()[%d] = %q, se esperaba %q (orden = %v)", i, nombres[i], nombre, nombres)
		}
	}

	if err := repo.Revelar("PRJ-001", itemID, 1); err != nil {
		t.Fatalf("Revelar() devolvió un error inesperado: %v", err)
	}
	votosLeidos, err := repo.Votos("PRJ-001", itemID, 1)
	if err != nil {
		t.Fatalf("Votos() devolvió un error inesperado: %v", err)
	}
	esperadosVotos := []domain.VotoPoker{
		{Integrante: "ana", Valor: 5},
		{Integrante: "Beto", Valor: 8},
		{Integrante: "beto", Valor: 13},
		{Integrante: "CARLA", Valor: 21},
		{Integrante: "Zoe", Valor: 2},
	}
	if len(votosLeidos) != len(esperadosVotos) {
		t.Fatalf("Votos() = %+v, se esperaban %+v", votosLeidos, esperadosVotos)
	}
	for i, esperado := range esperadosVotos {
		if votosLeidos[i] != esperado {
			t.Errorf("Votos()[%d] = %+v, se esperaba %+v (orden = %+v)", i, votosLeidos[i], esperado, votosLeidos)
		}
	}
}

// R4 / R11: con la ronda Revelada no se puede votar ni cambiar el voto; el
// voto previo permanece intacto.
func TestUS12_Postgres_RegistrarVotoTrasRevelar(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	itemID := itemIDUS12(t, pool, "PRJ-001")
	repo := postgres.NuevoRondasPoker(pool)

	if err := repo.Iniciar("PRJ-001", itemID, 1); err != nil {
		t.Fatalf("Iniciar() devolvió un error inesperado: %v", err)
	}
	if err := repo.RegistrarVoto("PRJ-001", itemID, 1, domain.VotoPoker{Integrante: "Ana", Valor: 5}); err != nil {
		t.Fatalf("RegistrarVoto inicial devolvió un error inesperado: %v", err)
	}
	if err := repo.Revelar("PRJ-001", itemID, 1); err != nil {
		t.Fatalf("Revelar() devolvió un error inesperado: %v", err)
	}

	err := repo.RegistrarVoto("PRJ-001", itemID, 1, domain.VotoPoker{Integrante: "Ana", Valor: 8})
	if !errors.Is(err, domain.ErrRondaRevelada) {
		t.Fatalf("RegistrarVoto tras revelar = %v, se esperaba %v", err, domain.ErrRondaRevelada)
	}

	votos, err := repo.Votos("PRJ-001", itemID, 1)
	if err != nil {
		t.Fatalf("Votos() devolvió un error inesperado: %v", err)
	}
	if len(votos) != 1 || votos[0].Valor != 5 {
		t.Errorf("el voto previo fue modificado: %+v, se esperaba Ana con 5", votos)
	}
}

// R7 / V9 / R11: revelar dos veces devuelve domain.ErrRondaRevelada; el
// primer revelado es definitivo.
func TestUS12_Postgres_RevelarDosVeces(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	itemID := itemIDUS12(t, pool, "PRJ-001")
	repo := postgres.NuevoRondasPoker(pool)

	if err := repo.Iniciar("PRJ-001", itemID, 1); err != nil {
		t.Fatalf("Iniciar() devolvió un error inesperado: %v", err)
	}
	if err := repo.Revelar("PRJ-001", itemID, 1); err != nil {
		t.Fatalf("primer Revelar() devolvió un error inesperado: %v", err)
	}

	err := repo.Revelar("PRJ-001", itemID, 1)
	if !errors.Is(err, domain.ErrRondaRevelada) {
		t.Fatalf("segundo Revelar() = %v, se esperaba %v", err, domain.ErrRondaRevelada)
	}
}

// V4 / E5: operar sobre una ronda inexistente devuelve domain.ErrRondaNoEncontrada.
func TestUS12_Postgres_RondaInexistente(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	itemID := itemIDUS12(t, pool, "PRJ-001")
	repo := postgres.NuevoRondasPoker(pool)

	if _, err := repo.Obtener("PRJ-001", itemID, 99); !errors.Is(err, domain.ErrRondaNoEncontrada) {
		t.Errorf("Obtener() = %v, se esperaba %v", err, domain.ErrRondaNoEncontrada)
	}
	err := repo.RegistrarVoto("PRJ-001", itemID, 99, domain.VotoPoker{Integrante: "Ana", Valor: 5})
	if !errors.Is(err, domain.ErrRondaNoEncontrada) {
		t.Errorf("RegistrarVoto() = %v, se esperaba %v", err, domain.ErrRondaNoEncontrada)
	}
	if err := repo.Revelar("PRJ-001", itemID, 99); !errors.Is(err, domain.ErrRondaNoEncontrada) {
		t.Errorf("Revelar() = %v, se esperaba %v", err, domain.ErrRondaNoEncontrada)
	}
}

// R9 / R10: el índice único parcial impide una segunda ronda Abierta por
// historia; tras revelar, una nueva ronda Abierta es válida.
func TestUS12_Postgres_IndiceParcialUnaAbiertaPorHistoria(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	itemID := itemIDUS12(t, pool, "PRJ-001")
	repo := postgres.NuevoRondasPoker(pool)

	if err := repo.Iniciar("PRJ-001", itemID, 1); err != nil {
		t.Fatalf("Iniciar() de la primera ronda devolvió un error inesperado: %v", err)
	}
	err := repo.Iniciar("PRJ-001", itemID, 2)
	if !errors.Is(err, domain.ErrRondaAbiertaExistente) {
		t.Fatalf("Iniciar() con una ronda Abierta pendiente = %v, se esperaba %v", err, domain.ErrRondaAbiertaExistente)
	}

	if err := repo.Revelar("PRJ-001", itemID, 1); err != nil {
		t.Fatalf("Revelar() devolvió un error inesperado: %v", err)
	}
	if err := repo.Iniciar("PRJ-001", itemID, 2); err != nil {
		t.Fatalf("Iniciar() tras revelar devolvió un error inesperado: %v", err)
	}
}

// V2 / E2: iniciar una ronda sobre una historia inexistente devuelve
// domain.ErrItemNoEncontrado.
func TestUS12_Postgres_IniciarHistoriaInexistente(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	repo := postgres.NuevoRondasPoker(pool)

	err := repo.Iniciar("PRJ-001", "HU-99", 1)
	if !errors.Is(err, domain.ErrItemNoEncontrado) {
		t.Fatalf("Iniciar() sobre historia inexistente = %v, se esperaba %v", err, domain.ErrItemNoEncontrado)
	}
}

// R19 / 13.2: la FK rondas_poker_item_fk requiere que items_backlog(proyecto_id,
// numero) sea PK o UNIQUE referenciable; verificado en el catálogo del motor y
// rechazando la inserción de una ronda para una historia inexistente.
func TestUS12_Postgres_FKHistoriaInexistente(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")

	rows, err := pool.Query(context.Background(), `
		SELECT tc.constraint_name,
		       string_agg(kcu.column_name, ',' ORDER BY kcu.ordinal_position) AS columnas
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu
		  ON kcu.constraint_name = tc.constraint_name
		 AND kcu.table_schema   = tc.table_schema
		 AND kcu.table_name     = tc.table_name
		WHERE tc.table_schema = current_schema()
		  AND tc.table_name = 'items_backlog'
		  AND tc.constraint_type IN ('PRIMARY KEY', 'UNIQUE')
		GROUP BY tc.constraint_name
	`)
	if err != nil {
		t.Fatalf("consultar claves candidatas de items_backlog: %v", err)
	}
	defer rows.Close()
	refenciable := false
	for rows.Next() {
		var constraintName, columnas string
		if err := rows.Scan(&constraintName, &columnas); err != nil {
			t.Fatalf("leer clave candidata de items_backlog: %v", err)
		}
		if columnas == "proyecto_id,numero" {
			refenciable = true
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterar claves candidatas de items_backlog: %v", err)
	}
	if !refenciable {
		t.Fatal("items_backlog(proyecto_id, numero) no es PK ni UNIQUE: la FK rondas_poker_item_fk no es referenciable")
	}

	_, err = pool.Exec(context.Background(), `
		INSERT INTO rondas_poker (proyecto_id, item_numero, numero, estado)
		VALUES ('PRJ-001', 99, 1, 'Abierta')
	`)
	if err == nil {
		t.Fatal("se permitió insertar una ronda para una historia inexistente")
	}
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		t.Fatalf("error de la FK = %T (%v), se esperaba un *pgconn.PgError con SQLSTATE 23503", err, err)
	}
	if pgErr.Code != "23503" {
		t.Errorf("SQLSTATE = %s, se esperaba 23503 (violación de clave foránea)", pgErr.Code)
	}
}

// R12 / D21: N goroutines votando el mismo integrante dejan una única fila.
func TestUS12_Postgres_RegistrarVotoConcurrenteMismoIntegrante(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	itemID := itemIDUS12(t, pool, "PRJ-001")
	repo := postgres.NuevoRondasPoker(pool)

	if err := repo.Iniciar("PRJ-001", itemID, 1); err != nil {
		t.Fatalf("Iniciar() devolvió un error inesperado: %v", err)
	}

	valores := []int{1, 2, 3, 5, 8, 13, 21, 1, 2, 3}
	var wg sync.WaitGroup
	errs := make(chan error, len(valores))
	for _, valor := range valores {
		wg.Add(1)
		go func(valor int) {
			defer wg.Done()
			errs <- repo.RegistrarVoto("PRJ-001", itemID, 1, domain.VotoPoker{Integrante: "Ana", Valor: valor})
		}(valor)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Errorf("RegistrarVoto concurrente del mismo integrante: %v", err)
	}

	if err := repo.Revelar("PRJ-001", itemID, 1); err != nil {
		t.Fatalf("Revelar() devolvió un error inesperado: %v", err)
	}
	votos, err := repo.Votos("PRJ-001", itemID, 1)
	if err != nil {
		t.Fatalf("Votos() devolvió un error inesperado: %v", err)
	}
	if len(votos) != 1 {
		t.Fatalf("votos persistidos = %+v, se esperaba una única fila para el integrante", votos)
	}
	enSerie := map[int]bool{1: true, 2: true, 3: true, 5: true, 8: true, 13: true, 21: true}
	if !enSerie[votos[0].Valor] {
		t.Errorf("valor final = %d, se esperaba un valor de la serie 1, 2, 3, 5, 8, 13, 21", votos[0].Valor)
	}
}

// R12 / D21: N goroutines con integrantes distintos persisten N filas.
func TestUS12_Postgres_RegistrarVotoConcurrenteIntegrantesDistintos(t *testing.T) {
	pool := nuevoPoolPostgres(t)
	seedProject(t, pool, "PRJ-001")
	itemID := itemIDUS12(t, pool, "PRJ-001")
	repo := postgres.NuevoRondasPoker(pool)

	if err := repo.Iniciar("PRJ-001", itemID, 1); err != nil {
		t.Fatalf("Iniciar() devolvió un error inesperado: %v", err)
	}

	const cantidad = 12
	valores := []int{1, 2, 3, 5, 8, 13, 21}
	var wg sync.WaitGroup
	errs := make(chan error, cantidad)
	for i := 0; i < cantidad; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			voto := domain.VotoPoker{
				Integrante: fmt.Sprintf("Integrante %02d", i+1),
				Valor:      valores[i%len(valores)],
			}
			errs <- repo.RegistrarVoto("PRJ-001", itemID, 1, voto)
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Errorf("RegistrarVoto concurrente de integrantes distintos: %v", err)
	}

	nombres, err := repo.YaVotaron("PRJ-001", itemID, 1)
	if err != nil {
		t.Fatalf("YaVotaron() devolvió un error inesperado: %v", err)
	}
	if len(nombres) != cantidad {
		t.Errorf("YaVotaron() = %v, se esperaban %d nombres", nombres, cantidad)
	}

	if err := repo.Revelar("PRJ-001", itemID, 1); err != nil {
		t.Fatalf("Revelar() devolvió un error inesperado: %v", err)
	}
	votos, err := repo.Votos("PRJ-001", itemID, 1)
	if err != nil {
		t.Fatalf("Votos() devolvió un error inesperado: %v", err)
	}
	if len(votos) != cantidad {
		t.Fatalf("Votos() = %+v, se esperaban %d votos", votos, cantidad)
	}
}

// R19: Up y Down de 00006 sobre una base vacía.
func TestUS12_Postgres_UpDownMigracion00006SobreBaseVacia(t *testing.T) {
	if _, err := migrations.FS.ReadFile("00006_planning_poker.sql"); err != nil {
		t.Fatalf("la migración 00006_planning_poker.sql no existe: %v", err)
	}

	pool, schema := nuevoPoolVacioPostgres(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()
	provider, err := goose.NewProvider(goose.DialectPostgres, db, migrations.FS)
	if err != nil {
		t.Fatalf("crear proveedor de migraciones: %v", err)
	}

	if _, err := provider.UpTo(ctx, 5); err != nil {
		t.Fatalf("aplicar migraciones 00001-00005: %v", err)
	}
	if existeTablaUS12(t, pool, schema, "rondas_poker") || existeTablaUS12(t, pool, schema, "votos_poker") {
		t.Fatal("rondas_poker/votos_poker existen antes de aplicar 00006")
	}

	if _, err := provider.UpTo(ctx, 6); err != nil {
		t.Fatalf("aplicar 00006_planning_poker: %v", err)
	}
	if !existeTablaUS12(t, pool, schema, "rondas_poker") {
		t.Error("rondas_poker no existe tras aplicar 00006")
	}
	if !existeTablaUS12(t, pool, schema, "votos_poker") {
		t.Error("votos_poker no existe tras aplicar 00006")
	}

	if _, err := provider.DownTo(ctx, 5); err != nil {
		t.Fatalf("revertir 00006_planning_poker: %v", err)
	}
	if existeTablaUS12(t, pool, schema, "rondas_poker") {
		t.Error("rondas_poker existe tras el Down de 00006")
	}
	if existeTablaUS12(t, pool, schema, "votos_poker") {
		t.Error("votos_poker existe tras el Down de 00006")
	}

	if _, err := provider.UpTo(ctx, 6); err != nil {
		t.Fatalf("reaplicar 00006_planning_poker: %v", err)
	}
	if !existeTablaUS12(t, pool, schema, "rondas_poker") || !existeTablaUS12(t, pool, schema, "votos_poker") {
		t.Error("las tablas no existen tras reaplicar 00006")
	}
}
