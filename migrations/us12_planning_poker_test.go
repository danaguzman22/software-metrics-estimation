package migrations

import (
	"strings"
	"testing"
)

const nombreMigracion00006 = "00006_planning_poker.sql"

func migr00006PlanningPoker(t *testing.T) string {
	t.Helper()
	migration, err := FS.ReadFile(nombreMigracion00006)
	if err != nil {
		t.Fatalf("leer migración %s: %v", nombreMigracion00006, err)
	}
	return string(migration)
}

func TestUS12_MigracionPlanningPokerUnicaVersion00006(t *testing.T) {
	entries, err := FS.ReadDir(".")
	if err != nil {
		t.Fatalf("listar migraciones embebidas: %v", err)
	}
	var nombres []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		nombre := entry.Name()
		if strings.HasPrefix(nombre, "00006_") && strings.HasSuffix(nombre, ".sql") {
			nombres = append(nombres, nombre)
		}
	}
	if len(nombres) != 1 {
		t.Fatalf("se esperaba exactamente una migración 00006_*.sql, se encontraron %d: %v", len(nombres), nombres)
	}
	if nombres[0] != nombreMigracion00006 {
		t.Errorf("la migración 00006 se llama %q, se esperaba %q", nombres[0], nombreMigracion00006)
	}
}

func TestUS12_MigracionPlanningPokerUp(t *testing.T) {
	sql := migr00006PlanningPoker(t)

	indiceDown := strings.Index(sql, "-- +goose Down")
	if indiceDown < 0 {
		t.Fatal("la migración no contiene el marcador -- +goose Down")
	}
	up := sql[:indiceDown]

	requiredFragments := []string{
		"-- +goose Up",
		"CREATE TABLE rondas_poker",
		"CREATE TABLE votos_poker",
		"PRIMARY KEY (proyecto_id, item_numero, numero)",
		"REFERENCES items_backlog (proyecto_id, numero)",
		"REFERENCES rondas_poker (proyecto_id, item_numero, numero)",
		"CREATE UNIQUE INDEX rondas_poker_historia_abierta_idx",
		"ON rondas_poker (proyecto_id, item_numero)",
		"WHERE estado = 'Abierta'",
		"IN (1, 2, 3, 5, 8, 13, 21)",
		"('Abierta', 'Revelada')",
	}
	for _, fragment := range requiredFragments {
		if !strings.Contains(sql, fragment) {
			t.Errorf("la migración no contiene %q", fragment)
		}
	}

	for _, estado := range []string{"Abierta", "Revelada"} {
		if !strings.Contains(up, "'"+estado+"'") {
			t.Errorf("la sección Up no admite el estado %q", estado)
		}
	}
}

func TestUS12_MigracionPlanningPokerDownEliminaVotosAntesQueRondas(t *testing.T) {
	sql := migr00006PlanningPoker(t)

	indiceDown := strings.Index(sql, "-- +goose Down")
	if indiceDown < 0 {
		t.Fatal("la migración no contiene el marcador -- +goose Down")
	}
	down := sql[indiceDown:]

	indiceVotos := strings.Index(down, "DROP TABLE votos_poker")
	if indiceVotos < 0 {
		t.Fatal("el Down no elimina votos_poker")
	}
	indiceRondas := strings.Index(down, "DROP TABLE rondas_poker")
	if indiceRondas < 0 {
		t.Fatal("el Down no elimina rondas_poker")
	}
	if indiceVotos > indiceRondas {
		t.Errorf("el Down elimina rondas_poker antes que votos_poker; votos_poker (FK) debe eliminarse primero")
	}
}
