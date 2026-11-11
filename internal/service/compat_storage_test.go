package service

import (
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/storage"
)

func nuevoAlmacenServicio(t *testing.T) (*storage.JSONStore, string) {
	t.Helper()
	dir := t.TempDir()
	almacen, err := storage.NewJSONStore(dir)
	if err != nil {
		t.Fatalf("NewJSONStore: %v", err)
	}
	return almacen, dir
}

func sembrarProyectos(t *testing.T, a interface {
	Save(name string, v any) error
}, registros []ProyectoRegistrado) {
	t.Helper()
	if err := a.Save("proyectos", registros); err != nil {
		t.Fatalf("Save(\"proyectos\"): %v", err)
	}
}

type almacenCuentaSaves struct {
	real          *storage.JSONStore
	llamadasASave int
}

func (a *almacenCuentaSaves) Save(name string, v any) error {
	a.llamadasASave++
	return a.real.Save(name, v)
}

func (a *almacenCuentaSaves) Load(name string, v any) error {
	return a.real.Load(name, v)
}

type almacenFalloEscritura struct {
	real *storage.JSONStore
	err  error
}

func (a almacenFalloEscritura) Save(name string, v any) error {
	return a.err
}

func (a almacenFalloEscritura) Load(name string, v any) error {
	return a.real.Load(name, v)
}
