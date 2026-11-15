package migrations

import (
	"io/fs"
	"regexp"
	"strings"
	"testing"
)

// D11 / D17 / D26 / D27: la migración 00007 crea la tabla de esfuerzo con
// horas numeric(4,2), CHECK de rango, FK compuesta a items_backlog, integrante
// como texto y secciones Up y Down.
func TestUS14_MigracionEsfuerzo(t *testing.T) {
	archivos, err := fs.Glob(FS, "00007_*.sql")
	if err != nil {
		t.Fatalf("buscar la migración 00007: %v", err)
	}
	if len(archivos) != 1 {
		t.Fatalf("se esperaba exactamente una migración 00007_*.sql, se encontraron %v", archivos)
	}

	contenido, err := FS.ReadFile(archivos[0])
	if err != nil {
		t.Fatalf("leer la migración %s: %v", archivos[0], err)
	}
	sql := string(contenido)

	fragmentos := []string{
		"-- +goose Up",
		"-- +goose Down",
		"CHECK (horas > 0 AND horas <= 24)",
	}
	for _, fragmento := range fragmentos {
		if !strings.Contains(sql, fragmento) {
			t.Errorf("la migración no contiene %q", fragmento)
		}
	}

	patrones := []struct {
		descripcion string
		patron      string
	}{
		{
			descripcion: "tabla de esfuerzo en Up",
			patron:      `(?i)CREATE TABLE\s+\w*esfuerzo\w*`,
		},
		{
			descripcion: "tabla de esfuerzo en Down",
			patron:      `(?i)DROP TABLE\s+\w*esfuerzo\w*`,
		},
		{
			descripcion: "columna horas numeric(4,2)",
			patron:      `(?i)horas\s+numeric\(4,\s*2\)`,
		},
		{
			descripcion: "columna integrante de texto",
			patron:      `(?i)integrante\s+text`,
		},
		{
			descripcion: "FK compuesta de proyecto e historia",
			patron:      `(?i)FOREIGN KEY\s*\(\s*proyecto_id\s*,\s*item_numero\s*\)`,
		},
		{
			descripcion: "FK a items_backlog por (proyecto_id, numero)",
			patron:      `(?i)REFERENCES\s+items_backlog\s*\(\s*proyecto_id\s*,\s*numero\s*\)`,
		},
		{
			descripcion: "FK a proyectos",
			patron:      `(?i)REFERENCES\s+proyectos\s*\(\s*id\s*\)`,
		},
	}
	for _, patron := range patrones {
		if !regexp.MustCompile(patron.patron).MatchString(sql) {
			t.Errorf("la migración no contiene %s (patrón %q)", patron.descripcion, patron.patron)
		}
	}
}
