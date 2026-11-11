package jsonrepo

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
	"github.com/danaguzman22/software-metrics-estimation/internal/service/repotest"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage"
)

func TestJSONRepo_Contrato(t *testing.T) {
	repotest.Ejecutar(t, nuevaEntornoJSON)
}

func nuevaEntornoJSON(t *testing.T) repotest.Entorno {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "data")
	store, err := storage.NewJSONStore(dir)
	if err != nil {
		t.Fatalf("NewJSONStore: %v", err)
	}
	return repotest.Entorno{
		Proyectos:   NuevoProyectos(store),
		Integrantes: NuevoIntegrantes(store),
		Backlog:     NuevoBacklog(store),
		SembrarProyecto: func(t *testing.T, p service.ProyectoRegistrado) {
			t.Helper()
			proyectos, err := cargarProyectosParaSemilla(store)
			if err != nil {
				t.Fatalf("cargar proyectos para semilla: %v", err)
			}
			proyectos = append(proyectos, registroProyecto{ID: p.ID, Proyecto: p.Proyecto})
			if err := store.Save("proyectos", proyectos); err != nil {
				t.Fatalf("guardar semilla de proyecto: %v", err)
			}
		},
		InyectarFallo: func(t *testing.T) (error, func() error) {
			t.Helper()
			causa := fmt.Errorf("causa infraestructura: %w", storage.ErrCorrupt)
			d := &dobleAlmacen{load: func(name string, v any) error { return causa }}
			repo := NuevoProyectos(d)
			return causa, func() error {
				_, err := repo.Listar()
				return err
			}
		},
	}
}

func cargarProyectosParaSemilla(store *storage.JSONStore) ([]registroProyecto, error) {
	var proyectos []registroProyecto
	if err := store.Load("proyectos", &proyectos); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return []registroProyecto{}, nil
		}
		return nil, err
	}
	if proyectos == nil {
		proyectos = []registroProyecto{}
	}
	return proyectos, nil
}

var _ = domain.NuevoProyecto
