package jsonrepo

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage"
)

type Almacen interface {
	Load(name string, v any) error
	Save(name string, v any) error
}

var _ service.RepositorioProyectos = (*Proyectos)(nil)
var _ service.RepositorioIntegrantes = (*Integrantes)(nil)
var _ service.RepositorioBacklog = (*Backlog)(nil)

type registroProyecto struct {
	ID       string          `json:"id"`
	Proyecto domain.Proyecto `json:"proyecto"`
}

type registroIntegrante struct {
	ProyectoID string            `json:"proyecto_id"`
	Integrante domain.Integrante `json:"integrante"`
}

type registroItem struct {
	ProyectoID string             `json:"proyecto_id"`
	Item       domain.ItemBacklog `json:"item"`
}

type Proyectos struct {
	almacen Almacen
	mu      sync.Mutex
}

func NuevoProyectos(a Almacen) *Proyectos {
	return &Proyectos{almacen: a}
}

func (r *Proyectos) cargar() ([]registroProyecto, error) {
	var registros []registroProyecto
	if err := r.almacen.Load("proyectos", &registros); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return []registroProyecto{}, nil
		}
		return nil, fmt.Errorf("cargar proyectos: %w", err)
	}
	if registros == nil {
		registros = []registroProyecto{}
	}
	return registros, nil
}

func (r *Proyectos) guardar(registros []registroProyecto) error {
	if err := r.almacen.Save("proyectos", registros); err != nil {
		return fmt.Errorf("guardar proyectos: %w", err)
	}
	return nil
}

func (r *Proyectos) Obtener(id string) (service.ProyectoRegistrado, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	registros, err := r.cargar()
	if err != nil {
		return service.ProyectoRegistrado{}, err
	}
	for _, reg := range registros {
		if reg.ID == id {
			return service.ProyectoRegistrado{ID: reg.ID, Proyecto: reg.Proyecto}, nil
		}
	}
	return service.ProyectoRegistrado{}, domain.ErrProyectoNoEncontrado
}

func (r *Proyectos) Existe(id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	registros, err := r.cargar()
	if err != nil {
		return false, err
	}
	for _, reg := range registros {
		if reg.ID == id {
			return true, nil
		}
	}
	return false, nil
}

func (r *Proyectos) Listar() ([]service.ProyectoRegistrado, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	registros, err := r.cargar()
	if err != nil {
		return nil, err
	}
	resultado := make([]service.ProyectoRegistrado, 0, len(registros))
	for _, reg := range registros {
		resultado = append(resultado, service.ProyectoRegistrado{ID: reg.ID, Proyecto: reg.Proyecto})
	}
	return resultado, nil
}

func (r *Proyectos) Actualizar(id string, proyecto domain.Proyecto) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	registros, err := r.cargar()
	if err != nil {
		return err
	}
	for i, reg := range registros {
		if reg.ID == id {
			registros[i].Proyecto = proyecto
			return r.guardar(registros)
		}
	}
	return domain.ErrProyectoNoEncontrado
}

type Integrantes struct {
	almacen Almacen
	mu      sync.Mutex
}

func NuevoIntegrantes(a Almacen) *Integrantes {
	return &Integrantes{almacen: a}
}

func (r *Integrantes) cargar() ([]registroIntegrante, error) {
	var registros []registroIntegrante
	if err := r.almacen.Load("integrantes", &registros); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return []registroIntegrante{}, nil
		}
		return nil, fmt.Errorf("cargar integrantes: %w", err)
	}
	if registros == nil {
		registros = []registroIntegrante{}
	}
	return registros, nil
}

func (r *Integrantes) guardar(registros []registroIntegrante) error {
	if err := r.almacen.Save("integrantes", registros); err != nil {
		return fmt.Errorf("guardar integrantes: %w", err)
	}
	return nil
}

func (r *Integrantes) ListarPorProyecto(proyectoID string) ([]domain.Integrante, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	registros, err := r.cargar()
	if err != nil {
		return nil, err
	}
	var resultado []domain.Integrante
	for _, reg := range registros {
		if reg.ProyectoID == proyectoID {
			resultado = append(resultado, reg.Integrante)
		}
	}
	if resultado == nil {
		resultado = []domain.Integrante{}
	}
	return resultado, nil
}

func (r *Integrantes) Agregar(proyectoID string, integrante domain.Integrante) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	registros, err := r.cargar()
	if err != nil {
		return err
	}
	registros = append(registros, registroIntegrante{ProyectoID: proyectoID, Integrante: integrante})
	return r.guardar(registros)
}

type Backlog struct {
	almacen Almacen
	mu      sync.Mutex
}

func NuevoBacklog(a Almacen) *Backlog {
	return &Backlog{almacen: a}
}

func (r *Backlog) cargar() ([]registroItem, error) {
	var registros []registroItem
	if err := r.almacen.Load("backlog", &registros); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return []registroItem{}, nil
		}
		return nil, fmt.Errorf("cargar backlog: %w", err)
	}
	if registros == nil {
		registros = []registroItem{}
	}
	return registros, nil
}

func (r *Backlog) guardar(registros []registroItem) error {
	if err := r.almacen.Save("backlog", registros); err != nil {
		return fmt.Errorf("guardar backlog: %w", err)
	}
	return nil
}

func (r *Backlog) AgregarItem(proyectoID string, item domain.ItemBacklog) (domain.ItemBacklog, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	registros, err := r.cargar()
	if err != nil {
		return domain.ItemBacklog{}, err
	}
	var ids []string
	for _, reg := range registros {
		if reg.ProyectoID == proyectoID {
			ids = append(ids, reg.Item.ID)
		}
	}
	item.ID = domain.SiguienteIDItem(ids)
	registros = append(registros, registroItem{ProyectoID: proyectoID, Item: item})
	if err := r.guardar(registros); err != nil {
		return domain.ItemBacklog{}, err
	}
	return item, nil
}

func (r *Backlog) ObtenerItem(proyectoID, itemID string) (domain.ItemBacklog, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	registros, err := r.cargar()
	if err != nil {
		return domain.ItemBacklog{}, err
	}
	for _, reg := range registros {
		if reg.ProyectoID == proyectoID && reg.Item.ID == itemID {
			return reg.Item, nil
		}
	}
	return domain.ItemBacklog{}, domain.ErrItemNoEncontrado
}

func (r *Backlog) ListarPorProyecto(proyectoID string) ([]domain.ItemBacklog, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	registros, err := r.cargar()
	if err != nil {
		return nil, err
	}
	var resultado []domain.ItemBacklog
	for _, reg := range registros {
		if reg.ProyectoID == proyectoID {
			resultado = append(resultado, reg.Item)
		}
	}
	if resultado == nil {
		resultado = []domain.ItemBacklog{}
	}
	sort.Slice(resultado, func(i, j int) bool {
		ni, _ := domain.NumeroDeItem(resultado[i].ID)
		nj, _ := domain.NumeroDeItem(resultado[j].ID)
		return ni < nj
	})
	return resultado, nil
}

func (r *Backlog) ActualizarItem(proyectoID string, item domain.ItemBacklog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	registros, err := r.cargar()
	if err != nil {
		return err
	}
	for i, reg := range registros {
		if reg.ProyectoID == proyectoID && reg.Item.ID == item.ID {
			registros[i].Item = item
			return r.guardar(registros)
		}
	}
	return domain.ErrItemNoEncontrado
}
