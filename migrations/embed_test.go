package migrations

import (
	"strings"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func TestEmbeddedInitialMigrationMatchesSchemaDecision(t *testing.T) {
	migration, err := FS.ReadFile("00001_esquema_inicial.sql")
	if err != nil {
		t.Fatalf("leer migración embebida: %v", err)
	}
	sql := string(migration)

	requiredFragments := []string{
		"-- +goose Up",
		"-- +goose Down",
		"CREATE TABLE proyectos",
		"id           text PRIMARY KEY DEFAULT gen_random_uuid()::text",
		"fecha_inicio date NOT NULL",
		"fecha_fin    date NOT NULL",
		"CREATE TABLE integrantes",
		"CREATE UNIQUE INDEX integrantes_proyecto_nombre_idx",
		"ON integrantes (proyecto_id, lower(nombre))",
		"CREATE TABLE items_backlog",
		"PRIMARY KEY (proyecto_id, numero)",
		"criterios_aceptacion text[] NOT NULL DEFAULT '{}'",
	}
	for _, fragment := range requiredFragments {
		if !strings.Contains(sql, fragment) {
			t.Errorf("la migración no contiene %q", fragment)
		}
	}

	for _, status := range []domain.Estado{
		domain.EstadoPlanificado,
		domain.EstadoEnCurso,
		domain.EstadoFinalizado,
	} {
		if !strings.Contains(sql, "'"+string(status)+"'") {
			t.Errorf("el CHECK de proyectos no admite el estado de dominio %q", status)
		}
	}
}
