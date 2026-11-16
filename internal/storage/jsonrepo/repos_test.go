package jsonrepo

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage"
)

// Aserciones de compilación (CA1, R1, R2)
var _ service.RepositorioProyectos = (*Proyectos)(nil)
var _ service.RepositorioIntegrantes = (*Integrantes)(nil)
var _ service.RepositorioBacklog = (*Backlog)(nil)

type dobleAlmacen struct {
	load func(name string, v any) error
	save func(name string, v any) error
}

func (d *dobleAlmacen) Load(name string, v any) error {
	if d.load != nil {
		return d.load(name, v)
	}
	return storage.ErrNotFound
}

func (d *dobleAlmacen) Save(name string, v any) error {
	if d.save != nil {
		return d.save(name, v)
	}
	return nil
}

func nuevoProyecto(nombre string) domain.Proyecto {
	return domain.Proyecto{
		Nombre:      nombre,
		FechaInicio: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		FechaFin:    time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		Estado:      domain.EstadoPlanificado,
	}
}

func nuevoItem(titulo, prioridad string) domain.ItemBacklog {
	it, _ := domain.NuevoItemBacklog(titulo, "desc", prioridad, []string{"c1"})
	return it
}

func TestT07_Proyectos_LeeFormatoActual(t *testing.T) {
	dir := t.TempDir()
	store, err := storage.NewJSONStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	// formato actual: [{"id":"...","proyecto":{...}}]
	type proyectoActual struct {
		Nombre      string    `json:"Nombre"`
		FechaInicio time.Time `json:"FechaInicio"`
		FechaFin    time.Time `json:"FechaFin"`
		Estado      string    `json:"Estado"`
	}
	type regActual struct {
		ID       string         `json:"id"`
		Proyecto proyectoActual `json:"proyecto"`
	}
	esperado := []regActual{{
		ID: "p1",
		Proyecto: proyectoActual{
			Nombre: "Alpha", FechaInicio: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			FechaFin: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), Estado: "Planificado",
		},
	}}
	if err := store.Save("proyectos", esperado); err != nil {
		t.Fatal(err)
	}
	repo := NuevoProyectos(store)
	got, err := repo.Listar()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "p1" || got[0].Proyecto.Nombre != "Alpha" || got[0].Proyecto.Estado != domain.EstadoPlanificado {
		t.Fatalf("lectura inesperada: %+v", got)
	}
}

func TestT07_Integrantes_LeeFormatoActual(t *testing.T) {
	dir := t.TempDir()
	store, err := storage.NewJSONStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	type integActual struct {
		Nombre string `json:"Nombre"`
		Rol    string `json:"Rol"`
	}
	type regActual struct {
		ProyectoID string      `json:"proyecto_id"`
		Integrante integActual `json:"integrante"`
	}
	if err := store.Save("integrantes", []regActual{
		{ProyectoID: "p1", Integrante: integActual{Nombre: "Ana", Rol: "Dev"}},
	}); err != nil {
		t.Fatal(err)
	}
	repo := NuevoIntegrantes(store)
	got, err := repo.ListarPorProyecto("p1")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Nombre != "Ana" || got[0].Rol != "Dev" {
		t.Fatalf("lectura inesperada: %+v", got)
	}
}

func TestT07_Backlog_LeeFormatoActual(t *testing.T) {
	dir := t.TempDir()
	store, err := storage.NewJSONStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	type itemActual struct {
		ID          string   `json:"ID"`
		Titulo      string   `json:"Titulo"`
		Descripcion string   `json:"Descripcion"`
		Prioridad   string   `json:"Prioridad"`
		Estado      string   `json:"Estado"`
		StoryPoints int      `json:"StoryPoints"`
		Criterios   []string `json:"Criterios"`
	}
	type regActual struct {
		ProyectoID string     `json:"proyecto_id"`
		Item       itemActual `json:"item"`
	}
	if err := store.Save("backlog", []regActual{{
		ProyectoID: "p1",
		Item:       itemActual{ID: "HU-01", Titulo: "Login", Descripcion: "d", Prioridad: "Alta", Estado: "Pendiente", Criterios: []string{"c"}},
	}}); err != nil {
		t.Fatal(err)
	}
	repo := NuevoBacklog(store)
	got, err := repo.ListarPorProyecto("p1")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "HU-01" || got[0].Titulo != "Login" {
		t.Fatalf("lectura inesperada: %+v", got)
	}
}

func TestT07_Proyectos_ActualizaYPersiste(t *testing.T) {
	dir := t.TempDir()
	store, _ := storage.NewJSONStore(dir)
	repo := NuevoProyectos(store)
	// sembrar con formato actual via store directo
	type proyectoActual struct {
		Nombre      string    `json:"Nombre"`
		FechaInicio time.Time `json:"FechaInicio"`
		FechaFin    time.Time `json:"FechaFin"`
		Estado      string    `json:"Estado"`
	}
	type regActual struct {
		ID       string         `json:"id"`
		Proyecto proyectoActual `json:"proyecto"`
	}
	p := nuevoProyecto("Viejo")
	_ = store.Save("proyectos", []regActual{{ID: "p1", Proyecto: proyectoActual{Nombre: p.Nombre, FechaInicio: p.FechaInicio, FechaFin: p.FechaFin, Estado: string(p.Estado)}}})
	act := nuevoProyecto("Nuevo")
	if err := repo.Actualizar("p1", act); err != nil {
		t.Fatal(err)
	}
	store2, _ := storage.NewJSONStore(dir)
	repo2 := NuevoProyectos(store2)
	got, err := repo2.Obtener("p1")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "p1" || got.Proyecto.Nombre != "Nuevo" {
		t.Fatalf("no persistió/conservó ID: %+v", got)
	}
}

func TestT07_Integrantes_AgregaYPersisteEnOrden(t *testing.T) {
	dir := t.TempDir()
	store, _ := storage.NewJSONStore(dir)
	repo := NuevoIntegrantes(store)
	if err := repo.Agregar("p1", domain.Integrante{Nombre: "Ana", Rol: "Dev"}); err != nil {
		t.Fatal(err)
	}
	if err := repo.Agregar("p2", domain.Integrante{Nombre: "Beto", Rol: "QA"}); err != nil {
		t.Fatal(err)
	}
	if err := repo.Agregar("p1", domain.Integrante{Nombre: "Caro", Rol: "PM"}); err != nil {
		t.Fatal(err)
	}
	store2, _ := storage.NewJSONStore(dir)
	repo2 := NuevoIntegrantes(store2)
	p1, err := repo2.ListarPorProyecto("p1")
	if err != nil {
		t.Fatal(err)
	}
	if len(p1) != 2 || p1[0].Nombre != "Ana" || p1[1].Nombre != "Caro" {
		t.Fatalf("orden/filtro p1: %+v", p1)
	}
	p2, err := repo2.ListarPorProyecto("p2")
	if err != nil || len(p2) != 1 || p2[0].Nombre != "Beto" {
		t.Fatalf("filtro p2: %+v err=%v", p2, err)
	}
}

func TestT07_Proyectos_ColeccionAusenteEsVacia(t *testing.T) {
	d := &dobleAlmacen{load: func(name string, v any) error { return storage.ErrNotFound }}
	repo := NuevoProyectos(d)
	lista, err := repo.Listar()
	if err != nil || len(lista) != 0 {
		t.Fatalf("Listar: %+v err=%v", lista, err)
	}
	existe, err := repo.Existe("x")
	if err != nil || existe {
		t.Fatalf("Existe: %v err=%v", existe, err)
	}
	_, err = repo.Obtener("x")
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("Obtener esperaba ErrProyectoNoEncontrado, got %v", err)
	}
}

func TestT07_Integrantes_ColeccionAusenteEsVacia(t *testing.T) {
	d := &dobleAlmacen{load: func(name string, v any) error { return storage.ErrNotFound }}
	repo := NuevoIntegrantes(d)
	lista, err := repo.ListarPorProyecto("p1")
	if err != nil || len(lista) != 0 {
		t.Fatalf("ListarPorProyecto: %+v err=%v", lista, err)
	}
}

func TestT07_Backlog_ColeccionAusenteEsVacia(t *testing.T) {
	d := &dobleAlmacen{load: func(name string, v any) error { return storage.ErrNotFound }}
	repo := NuevoBacklog(d)
	lista, err := repo.ListarPorProyecto("p1")
	if err != nil || len(lista) != 0 {
		t.Fatalf("ListarPorProyecto: %+v err=%v", lista, err)
	}
	it, err := repo.AgregarItem("p1", nuevoItem("Primero", "Alta"))
	if err != nil || it.ID != "HU-01" {
		t.Fatalf("AgregarItem sobre ausente: %+v err=%v", it, err)
	}
}

func TestT07_ErrCorruptSePropagaConW(t *testing.T) {
	corrupt := fmt.Errorf("boom: %w", storage.ErrCorrupt)
	d := &dobleAlmacen{load: func(name string, v any) error { return corrupt }}

	t.Run("proyectos/Obtener", func(t *testing.T) {
		_, err := NuevoProyectos(d).Obtener("p1")
		if !errors.Is(err, storage.ErrCorrupt) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("proyectos/Existe", func(t *testing.T) {
		_, err := NuevoProyectos(d).Existe("p1")
		if !errors.Is(err, storage.ErrCorrupt) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("proyectos/Listar", func(t *testing.T) {
		_, err := NuevoProyectos(d).Listar()
		if !errors.Is(err, storage.ErrCorrupt) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("proyectos/Actualizar", func(t *testing.T) {
		err := NuevoProyectos(d).Actualizar("p1", nuevoProyecto("N"))
		if !errors.Is(err, storage.ErrCorrupt) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("integrantes/ListarPorProyecto", func(t *testing.T) {
		_, err := NuevoIntegrantes(d).ListarPorProyecto("p1")
		if !errors.Is(err, storage.ErrCorrupt) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("integrantes/Agregar", func(t *testing.T) {
		err := NuevoIntegrantes(d).Agregar("p1", domain.Integrante{Nombre: "A", Rol: "R"})
		if !errors.Is(err, storage.ErrCorrupt) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("backlog/AgregarItem", func(t *testing.T) {
		_, err := NuevoBacklog(d).AgregarItem("p1", nuevoItem("T", "Alta"))
		if !errors.Is(err, storage.ErrCorrupt) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("backlog/ObtenerItem", func(t *testing.T) {
		_, err := NuevoBacklog(d).ObtenerItem("p1", "HU-01")
		if !errors.Is(err, storage.ErrCorrupt) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("backlog/ListarPorProyecto", func(t *testing.T) {
		_, err := NuevoBacklog(d).ListarPorProyecto("p1")
		if !errors.Is(err, storage.ErrCorrupt) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("backlog/ActualizarItem", func(t *testing.T) {
		it := nuevoItem("T", "Alta")
		it.ID = "HU-01"
		err := NuevoBacklog(d).ActualizarItem("p1", it)
		if !errors.Is(err, storage.ErrCorrupt) {
			t.Fatalf("err=%v", err)
		}
	})
}

func TestT07_ErrorDeGuardadoSePropaga(t *testing.T) {
	saveErr := fmt.Errorf("disk: %w", errors.New("fail"))
	dir := t.TempDir()
	store, err := storage.NewJSONStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Save("proyectos", []map[string]any{{
		"id": "p1",
		"proyecto": map[string]any{
			"Nombre": "X", "FechaInicio": time.Now(), "FechaFin": time.Now(), "Estado": "Planificado",
		},
	}}); err != nil {
		t.Fatal(err)
	}
	if err := store.Save("backlog", []map[string]any{{
		"proyecto_id": "p1",
		"item": map[string]any{
			"ID": "HU-01", "Titulo": "T", "Descripcion": "d", "Prioridad": "Alta",
			"Estado": "Pendiente", "StoryPoints": 0, "Criterios": []string{"c"},
		},
	}}); err != nil {
		t.Fatal(err)
	}

	failSave := &dobleAlmacen{
		load: store.Load,
		save: func(name string, v any) error { return saveErr },
	}

	t.Run("proyectos/Actualizar", func(t *testing.T) {
		if err := NuevoProyectos(failSave).Actualizar("p1", nuevoProyecto("Y")); !errors.Is(err, saveErr) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("integrantes/Agregar", func(t *testing.T) {
		if err := NuevoIntegrantes(failSave).Agregar("p1", domain.Integrante{Nombre: "A", Rol: "R"}); !errors.Is(err, saveErr) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("backlog/AgregarItem", func(t *testing.T) {
		if _, err := NuevoBacklog(failSave).AgregarItem("p1", nuevoItem("T", "Alta")); !errors.Is(err, saveErr) {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("backlog/ActualizarItem", func(t *testing.T) {
		it := nuevoItem("T2", "Media")
		it.ID = "HU-01"
		if err := NuevoBacklog(failSave).ActualizarItem("p1", it); !errors.Is(err, saveErr) {
			t.Fatalf("err=%v", err)
		}
	})
}

func TestT07_Backlog_AgregarItemNumeraPorProyecto(t *testing.T) {
	dir := t.TempDir()
	store, _ := storage.NewJSONStore(dir)
	repo := NuevoBacklog(store)
	i1, err := repo.AgregarItem("p1", nuevoItem("A", "Alta"))
	if err != nil || i1.ID != "HU-01" {
		t.Fatalf("p1-1: %+v %v", i1, err)
	}
	i2, err := repo.AgregarItem("p1", nuevoItem("B", "Alta"))
	if err != nil || i2.ID != "HU-02" {
		t.Fatalf("p1-2: %+v %v", i2, err)
	}
	j1, err := repo.AgregarItem("p2", nuevoItem("C", "Baja"))
	if err != nil || j1.ID != "HU-01" {
		t.Fatalf("p2-1: %+v %v", j1, err)
	}
	store2, _ := storage.NewJSONStore(dir)
	repo2 := NuevoBacklog(store2)
	lp1, _ := repo2.ListarPorProyecto("p1")
	lp2, _ := repo2.ListarPorProyecto("p2")
	if len(lp1) != 2 || lp1[0].ID != "HU-01" || lp1[1].ID != "HU-02" {
		t.Fatalf("p1: %+v", lp1)
	}
	if len(lp2) != 1 || lp2[0].ID != "HU-01" {
		t.Fatalf("p2: %+v", lp2)
	}
}

func TestT07_Backlog_ListarPorProyectoOrdenaPorNumero(t *testing.T) {
	dir := t.TempDir()
	store, _ := storage.NewJSONStore(dir)
	// sembrar fuera de orden
	_ = store.Save("backlog", []map[string]any{
		{"proyecto_id": "p1", "item": map[string]any{"ID": "HU-10", "Titulo": "diez", "Descripcion": "", "Prioridad": "Alta", "Estado": "Pendiente", "StoryPoints": 0, "Criterios": []string{}}},
		{"proyecto_id": "p1", "item": map[string]any{"ID": "HU-02", "Titulo": "dos", "Descripcion": "", "Prioridad": "Alta", "Estado": "Pendiente", "StoryPoints": 0, "Criterios": []string{}}},
		{"proyecto_id": "p1", "item": map[string]any{"ID": "HU-01", "Titulo": "uno", "Descripcion": "", "Prioridad": "Alta", "Estado": "Pendiente", "StoryPoints": 0, "Criterios": []string{}}},
	})
	repo := NuevoBacklog(store)
	got, err := repo.ListarPorProyecto("p1")
	if err != nil {
		t.Fatal(err)
	}
	ids := []string{got[0].ID, got[1].ID, got[2].ID}
	want := []string{"HU-01", "HU-02", "HU-10"}
	if !reflect.DeepEqual(ids, want) {
		t.Fatalf("orden: got %v want %v", ids, want)
	}
}

func TestT07_Backlog_ObtenerYActualizarItem(t *testing.T) {
	dir := t.TempDir()
	store, _ := storage.NewJSONStore(dir)
	repo := NuevoBacklog(store)
	it, _ := repo.AgregarItem("p1", nuevoItem("Orig", "Alta"))
	it.Titulo = "Mod"
	it.Prioridad = "Media"
	it.Criterios = []string{"a", "b"}
	if err := repo.ActualizarItem("p1", it); err != nil {
		t.Fatal(err)
	}
	store2, _ := storage.NewJSONStore(dir)
	repo2 := NuevoBacklog(store2)
	got, err := repo2.ObtenerItem("p1", it.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Titulo != "Mod" || got.Prioridad != "Media" || !reflect.DeepEqual(got.Criterios, []string{"a", "b"}) {
		t.Fatalf("round-trip: %+v", got)
	}
}

func TestT07_Backlog_AgregarItemConcurrente(t *testing.T) {
	dir := t.TempDir()
	store, _ := storage.NewJSONStore(dir)
	repo := NuevoBacklog(store)
	const n = 20
	var wg sync.WaitGroup
	wg.Add(n)
	errs := make(chan error, n)
	ids := make(chan string, n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			it, err := repo.AgregarItem("p1", nuevoItem(fmt.Sprintf("t%d", i), "Alta"))
			if err != nil {
				errs <- err
				return
			}
			ids <- it.ID
		}(i)
	}
	wg.Wait()
	close(errs)
	close(ids)
	for err := range errs {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for id := range ids {
		if seen[id] {
			t.Fatalf("id duplicado: %s", id)
		}
		seen[id] = true
	}
	if len(seen) != n {
		t.Fatalf("esperaba %d ids, got %d", n, len(seen))
	}
	for i := 1; i <= n; i++ {
		want := fmt.Sprintf("HU-%02d", i)
		if !seen[want] {
			t.Fatalf("falta id %s en %v", want, seen)
		}
	}
	store2, _ := storage.NewJSONStore(dir)
	repo2 := NuevoBacklog(store2)
	lista, err := repo2.ListarPorProyecto("p1")
	if err != nil || len(lista) != n {
		t.Fatalf("persistidos: %d err=%v", len(lista), err)
	}
}


