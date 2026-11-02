package domain

import (
	"fmt"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// TestDomainPackageCompiles es un smoke test complementario de CA1 (R1).
// Comprueba que internal/domain contiene al menos un archivo Go de producción
// (no _test.go) y que el paquete se puede importar/construir como paquete Go.
// No reemplaza la verificación completa con `go build ./...`.
func TestDomainPackageCompiles(t *testing.T) {
	tests := []struct {
		name string
		dir  string
	}{
		{name: "internal/domain tiene archivos Go de producción", dir: "."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg, err := build.ImportDir(tt.dir, 0)
			if err != nil {
				if _, ok := err.(*build.NoGoError); ok {
					t.Fatalf("internal/domain no contiene archivos Go de producción (solo _test.go): %v", err)
				}
				t.Fatalf("no se pudo importar internal/domain: %v", err)
			}

			if len(pkg.GoFiles) == 0 {
				t.Fatal("internal/domain no contiene archivos Go de producción (no _test.go)")
			}
		})
	}
}

// TestDomainDoesNotImportWeb verifica CA3 (R3): internal/domain no debe importar
// internal/web ni net/http. Además falla si no encuentra archivos Go de producción,
// para que el RED sea significativo y no un falso positivo.
func TestDomainDoesNotImportWeb(t *testing.T) {
	t.Run("clasificación de imports prohibidos", func(t *testing.T) {
		tests := []struct {
			name       string
			importPath string
			forbidden  bool
		}{
			{name: "ruta interna directa internal/web", importPath: "internal/web", forbidden: true},
			{name: "ruta interna con prefijo de módulo", importPath: "github.com/danaguzman22/software-metrics-estimation/internal/web", forbidden: true},
			{name: "paquete estándar net/http", importPath: "net/http", forbidden: true},
			{name: "paquete estándar permitido fmt", importPath: "fmt", forbidden: false},
			{name: "paquete estándar permitido net/url", importPath: "net/url", forbidden: false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := isForbiddenImport(tt.importPath); got != tt.forbidden {
					t.Fatalf("isForbiddenImport(%q) = %v, se esperaba %v", tt.importPath, got, tt.forbidden)
				}
			})
		}
	})

	t.Run("imports reales de internal/domain", func(t *testing.T) {
		files, err := productionGoFiles(".")
		if err != nil {
			t.Fatalf("no se pudieron inspeccionar los archivos Go de producción de internal/domain: %v", err)
		}

		imports, err := importsFromFiles(files)
		if err != nil {
			t.Fatalf("no se pudieron parsear los imports de internal/domain: %v", err)
		}

		for _, importPath := range imports {
			if isForbiddenImport(importPath) {
				t.Errorf("internal/domain importa %q, prohibido por la separación domain/web", importPath)
			}
		}
	})
}

// productionGoFiles devuelve los archivos .go de producción del directorio dado,
// excluyendo los archivos _test.go. Falla explícitamente si no hay ninguno.
func productionGoFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		files = append(files, filepath.Join(dir, name))
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no hay archivos Go de producción (no _test.go) en %s", dir)
	}

	return files, nil
}

// importsFromFiles parsea los imports de los archivos indicados usando go/parser.
func importsFromFiles(files []string) ([]string, error) {
	fset := token.NewFileSet()

	var imports []string
	for _, file := range files {
		parsed, err := parser.ParseFile(fset, file, nil, parser.ImportsOnly)
		if err != nil {
			return nil, fmt.Errorf("error al parsear %s: %w", file, err)
		}

		for _, spec := range parsed.Imports {
			importPath, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				return nil, fmt.Errorf("import inválido en %s: %w", file, err)
			}
			imports = append(imports, importPath)
		}
	}

	return imports, nil
}

// isForbiddenImport retorna true si el import rompe la separación domain/web
// o la guarda técnica de desacople de net/http.
func isForbiddenImport(importPath string) bool {
	if importPath == "net/http" {
		return true
	}

	return importPath == "internal/web" || strings.HasSuffix(importPath, "/internal/web")
}
