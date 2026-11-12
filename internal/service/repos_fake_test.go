package service

import (
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

type fakeRepositorioProyectos struct {
	proyectos       map[string]domain.Proyecto
	obtenerErr      error
	existeErr       error
	actualizarErr   error
	obtenerLlamadas int
	existeLlamadas  int
	actualizaciones int
}

func (f *fakeRepositorioProyectos) Obtener(id string) (ProyectoRegistrado, error) {
	f.obtenerLlamadas++
	if f.obtenerErr != nil {
		return ProyectoRegistrado{}, f.obtenerErr
	}
	proyecto, ok := f.proyectos[id]
	if !ok {
		return ProyectoRegistrado{}, domain.ErrProyectoNoEncontrado
	}
	return ProyectoRegistrado{ID: id, Proyecto: proyecto}, nil
}

func (f *fakeRepositorioProyectos) Existe(id string) (bool, error) {
	f.existeLlamadas++
	if f.existeErr != nil {
		return false, f.existeErr
	}
	_, ok := f.proyectos[id]
	return ok, nil
}

func (f *fakeRepositorioProyectos) Listar() ([]ProyectoRegistrado, error) {
	registros := make([]ProyectoRegistrado, 0, len(f.proyectos))
	for id, proyecto := range f.proyectos {
		registros = append(registros, ProyectoRegistrado{ID: id, Proyecto: proyecto})
	}
	return registros, nil
}

func (f *fakeRepositorioProyectos) Actualizar(id string, proyecto domain.Proyecto) error {
	f.actualizaciones++
	if f.actualizarErr != nil {
		return f.actualizarErr
	}
	if _, ok := f.proyectos[id]; !ok {
		return domain.ErrProyectoNoEncontrado
	}
	f.proyectos[id] = proyecto
	return nil
}

type fakeRepositorioIntegrantes struct {
	integrantes      map[string][]domain.Integrante
	listarErr        error
	agregarErr       error
	listarLlamadas   int
	agregados        int
	ultimoProyectoID string
}

func (f *fakeRepositorioIntegrantes) ListarPorProyecto(proyectoID string) ([]domain.Integrante, error) {
	f.listarLlamadas++
	if f.listarErr != nil {
		return nil, f.listarErr
	}
	return append([]domain.Integrante(nil), f.integrantes[proyectoID]...), nil
}

func (f *fakeRepositorioIntegrantes) Agregar(proyectoID string, integrante domain.Integrante) error {
	f.agregados++
	f.ultimoProyectoID = proyectoID
	if f.agregarErr != nil {
		return f.agregarErr
	}
	f.integrantes[proyectoID] = append(f.integrantes[proyectoID], integrante)
	return nil
}

type fakeRepositorioBacklog struct {
	items            map[string][]domain.ItemBacklog
	agregarErr       error
	obtenerErr       error
	listarErr        error
	actualizarErr    error
	agregados        int
	obtenerLlamadas  int
	actualizaciones  int
	idAsignado       string
	ultimoProyectoID string
	ultimoItemID     string
}

func (f *fakeRepositorioBacklog) AgregarItem(proyectoID string, item domain.ItemBacklog) (domain.ItemBacklog, error) {
	f.agregados++
	f.ultimoProyectoID = proyectoID
	if f.agregarErr != nil {
		return domain.ItemBacklog{}, f.agregarErr
	}

	persistido := item
	if f.idAsignado != "" {
		persistido.ID = f.idAsignado
	} else {
		ids := make([]string, 0, len(f.items[proyectoID]))
		for _, existente := range f.items[proyectoID] {
			ids = append(ids, existente.ID)
		}
		persistido.ID = domain.SiguienteIDItem(ids)
	}

	f.items[proyectoID] = append(f.items[proyectoID], persistido)
	return persistido, nil
}

func (f *fakeRepositorioBacklog) ObtenerItem(proyectoID, itemID string) (domain.ItemBacklog, error) {
	f.obtenerLlamadas++
	f.ultimoProyectoID = proyectoID
	f.ultimoItemID = itemID
	if f.obtenerErr != nil {
		return domain.ItemBacklog{}, f.obtenerErr
	}

	for _, item := range f.items[proyectoID] {
		if item.ID == itemID {
			return item, nil
		}
	}
	return domain.ItemBacklog{}, domain.ErrItemNoEncontrado
}

func (f *fakeRepositorioBacklog) ListarPorProyecto(proyectoID string) ([]domain.ItemBacklog, error) {
	if f.listarErr != nil {
		return nil, f.listarErr
	}
	return append([]domain.ItemBacklog(nil), f.items[proyectoID]...), nil
}

func (f *fakeRepositorioBacklog) ActualizarItem(proyectoID string, item domain.ItemBacklog) error {
	f.actualizaciones++
	f.ultimoProyectoID = proyectoID
	f.ultimoItemID = item.ID
	if f.actualizarErr != nil {
		return f.actualizarErr
	}

	items := f.items[proyectoID]
	for i := range items {
		if items[i].ID == item.ID {
			items[i] = item
			f.items[proyectoID] = items
			return nil
		}
	}
	return domain.ErrItemNoEncontrado
}

func nuevoFakeProyectos() *fakeRepositorioProyectos {
	return &fakeRepositorioProyectos{proyectos: map[string]domain.Proyecto{}}
}

func nuevoFakeIntegrantes() *fakeRepositorioIntegrantes {
	return &fakeRepositorioIntegrantes{integrantes: map[string][]domain.Integrante{}}
}

func nuevoFakeBacklog() *fakeRepositorioBacklog {
	return &fakeRepositorioBacklog{items: map[string][]domain.ItemBacklog{}}
}

func proyectoFakeValido() domain.Proyecto {
	return domain.Proyecto{
		Nombre:      "Sistema de Métricas",
		FechaInicio: fechaFake(2026, 1, 1),
		FechaFin:    fechaFake(2026, 12, 31),
		Estado:      domain.EstadoPlanificado,
	}
}

func fechaFake(anio, mes, dia int) time.Time {
	return time.Date(anio, time.Month(mes), dia, 0, 0, 0, 0, time.UTC)
}
