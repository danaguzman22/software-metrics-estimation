package storage

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func newTestStore(t *testing.T) (*JSONStore, string) {
	t.Helper()
	dir := t.TempDir()
	s, err := NewJSONStore(dir)
	if err != nil {
		t.Fatalf("NewJSONStore: %v", err)
	}
	return s, dir
}

func TestSaveLoad_RoundTrip(t *testing.T) {
	s, _ := newTestStore(t)
	want := item{ID: 1, Name: "Demo"}
	if err := s.Save("proyectos", want); err != nil {
		t.Fatalf("Save: %v", err)
	}
	var got item
	if err := s.Load("proyectos", &got); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got != want {
		t.Errorf("esperaba %+v, obtuve %+v", want, got)
	}
}

func TestLoad_NotFound(t *testing.T) {
	s, _ := newTestStore(t)
	var got item
	if err := s.Load("inexistente", &got); !errors.Is(err, ErrNotFound) {
		t.Errorf("esperaba ErrNotFound, obtuve %v", err)
	}
}

func TestLoad_CorruptJSON(t *testing.T) {
	s, dir := newTestStore(t)
	if err := os.WriteFile(filepath.Join(dir, "roto.json"), []byte("{no es json"), 0o644); err != nil {
		t.Fatal(err)
	}
	var got item
	if err := s.Load("roto", &got); !errors.Is(err, ErrCorrupt) {
		t.Errorf("esperaba ErrCorrupt, obtuve %v", err)
	}
}

func TestLoad_NonPointerTarget_NoEsCorrupt(t *testing.T) {
	s, _ := newTestStore(t)
	if err := s.Save("x", item{ID: 1}); err != nil {
		t.Fatal(err)
	}
	err := s.Load("x", item{})
	if err == nil || errors.Is(err, ErrCorrupt) {
		t.Errorf("esperaba error distinto de ErrCorrupt, obtuve %v", err)
	}
}

func TestSave_InvalidName(t *testing.T) {
	s, dir := newTestStore(t)
	for _, name := range []string{"", ".", "..", "../x", "a/b", `a\b`} {
		t.Run(name, func(t *testing.T) {
			if err := s.Save(name, item{ID: 1}); !errors.Is(err, ErrInvalidName) {
				t.Errorf("esperaba ErrInvalidName, obtuve %v", err)
			}
			var got item
			if err := s.Load(name, &got); !errors.Is(err, ErrInvalidName) {
				t.Errorf("Load: esperaba ErrInvalidName, obtuve %v", err)
			}
		})
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Errorf("no debía crearse ningún archivo, hay %d", len(entries))
	}
}

func TestNewJSONStore_CreatesNestedDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "a", "b", "data")
	if _, err := NewJSONStore(dir); err != nil {
		t.Fatalf("NewJSONStore: %v", err)
	}
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		t.Errorf("el directorio %s debía existir", dir)
	}
}

func TestNewJSONStore_EmptyDir(t *testing.T) {
	if _, err := NewJSONStore("  "); err == nil {
		t.Error("esperaba error con directorio vacío")
	}
}

func TestSave_OverwritesAndLeavesNoTempFiles(t *testing.T) {
	s, dir := newTestStore(t)
	if err := s.Save("p", item{ID: 1, Name: "Viejo"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Save("p", item{ID: 1, Name: "Nuevo"}); err != nil {
		t.Fatal(err)
	}
	var got item
	if err := s.Load("p", &got); err != nil || got.Name != "Nuevo" {
		t.Errorf("esperaba Nuevo, obtuve %+v (err=%v)", got, err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 || entries[0].Name() != "p.json" {
		t.Errorf("debía quedar solo p.json, hay %v", entries)
	}
}

func TestSave_UnmarshalableValue(t *testing.T) {
	s, dir := newTestStore(t)
	if err := s.Save("canal", make(chan int)); err == nil {
		t.Error("esperaba error con valor no serializable")
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Errorf("no debía quedar ningún archivo, hay %d", len(entries))
	}
}
