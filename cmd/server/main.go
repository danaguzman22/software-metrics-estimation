package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/danaguzman22/software-metrics-estimation/internal/storage/postgres"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func requiredDatabaseURL() (string, error) {
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		return "", errors.New("la variable DATABASE_URL es obligatoria")
	}
	return databaseURL, nil
}

func run(ctx context.Context) error {
	databaseURL, err := requiredDatabaseURL()
	if err != nil {
		return err
	}

	pool, err := postgres.ConnectAndMigrate(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("inicializar PostgreSQL: %w", err)
	}
	defer pool.Close()

	return nil
}
