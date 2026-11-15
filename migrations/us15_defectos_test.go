package migrations

import (
	"io/fs"
	"regexp"
	"strings"
	"testing"
)

func TestUS15_MigracionDefectos(t *testing.T) {
	archivos, err := fs.Glob(FS, "00009_*.sql")
	if err != nil {
		t.Fatalf("buscar migración 00009: %v", err)
	}
	if len(archivos) != 1 {
		t.Fatalf("se esperaba una migración 00009, se encontraron %v", archivos)
	}
	contenido, err := FS.ReadFile(archivos[0])
	if err != nil {
		t.Fatalf("leer migración: %v", err)
	}
	sql := string(contenido)
	for _, fragmento := range []string{
		"-- +goose Up",
		"-- +goose Down",
		"CREATE TABLE defectos",
		"severidad",
		"estado",
		"item_numero",
		"sprint_deteccion_numero",
		"sprint_resolucion_numero",
		"Alta",
		"Media",
		"Baja",
		"Abierto",
		"En curso",
		"Resuelto",
	} {
		if !strings.Contains(sql, fragmento) {
			t.Errorf("la migración no contiene %q", fragmento)
		}
	}
	for _, patron := range []string{
		`(?i)CREATE TABLE\s+defectos`,
		`(?i)FOREIGN KEY\s*\(\s*proyecto_id\s*,\s*item_numero\s*\)`,
		`(?i)FOREIGN KEY\s*\(\s*proyecto_id\s*,\s*sprint_deteccion_numero\s*\)`,
		`(?i)FOREIGN KEY\s*\(\s*proyecto_id\s*,\s*sprint_resolucion_numero\s*\)`,
		`(?i)CHECK\s*\([^)]*severidad`,
		`(?i)CHECK\s*\([^)]*estado`,
	} {
		if !regexp.MustCompile(patron).MatchString(sql) {
			t.Errorf("la migración no cumple el patrón %q", patron)
		}
	}
}
