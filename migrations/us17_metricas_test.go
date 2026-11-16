package migrations

import (
	"io/fs"
	"regexp"
	"strings"
	"testing"
)

func TestUS17_MigracionMetricasEsfuerzo(t *testing.T) {
	archivos, err := fs.Glob(FS, "00008_*.sql")
	if err != nil {
		t.Fatalf("buscar migración 00008: %v", err)
	}
	if len(archivos) != 1 {
		t.Fatalf("se esperaba una migración 00008, se encontraron %v", archivos)
	}
	contenido, err := FS.ReadFile(archivos[0])
	if err != nil {
		t.Fatalf("leer migración: %v", err)
	}
	sql := string(contenido)
	for _, fragmento := range []string{
		"-- +goose Up",
		"-- +goose Down",
		"horas_estimadas",
		"sprint_numero",
		"sprint_esfuerzo_snapshots",
		"sprint_esfuerzo_items",
	} {
		if !strings.Contains(sql, fragmento) {
			t.Errorf("la migración no contiene %q", fragmento)
		}
	}
	for _, patron := range []string{
		`(?i)ALTER TABLE\s+items_backlog`,
		`(?i)ALTER TABLE\s+esfuerzos`,
		`(?i)CREATE TABLE\s+sprint_esfuerzo_snapshots`,
		`(?i)CREATE TABLE\s+sprint_esfuerzo_items`,
		`(?i)FOREIGN KEY\s*\(\s*proyecto_id\s*,\s*sprint_numero\s*\)`,
	} {
		if !regexp.MustCompile(patron).MatchString(sql) {
			t.Errorf("la migración no cumple el patrón %q", patron)
		}
	}
}
