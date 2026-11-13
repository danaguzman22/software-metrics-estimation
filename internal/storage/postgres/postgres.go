package postgres

import (
	"context"
	"fmt"

	"github.com/danaguzman22/software-metrics-estimation/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// ConnectAndMigrate abre y valida el pool de PostgreSQL, y aplica las
// migraciones antes de devolverlo al llamador.
func ConnectAndMigrate(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("crear pool PostgreSQL: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("conectar a PostgreSQL: %w", err)
	}

	if err := Migrate(ctx, pool); err != nil {
		pool.Close()
		return nil, fmt.Errorf("aplicar migraciones PostgreSQL: %w", err)
	}

	return pool, nil
}

// Migrate aplica las migraciones embebidas al pool indicado.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	provider, err := goose.NewProvider(goose.DialectPostgres, db, migrations.FS)
	if err != nil {
		return fmt.Errorf("crear proveedor de migraciones: %w", err)
	}
	if _, err := provider.Up(ctx); err != nil {
		return fmt.Errorf("ejecutar migraciones: %w", err)
	}

	return nil
}
