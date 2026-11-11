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
	agregarErr      error
	obtenerErr      error
	listarErr       error
	actualizarErr   error
	agregados       int
	obtenerLlamadas int
	actualizaciones int
}

func (f *fakeRepositorioBacklog) AgregarItem(string, domain.ItemBacklog) (domain.ItemBacklog, error) {
	f.agregados++
	return domain.ItemBacklog{}, f.agregarErr
}

func (f *fakeRepositorioBacklog) ObtenerItem(string, string) (domain.ItemBacklog, error) {
	f.obtenerLlamadas++
	return domain.ItemBacklog{}, f.obtenerErr
}

func (f *fakeRepositorioBacklog) ListarPorProyecto(string) ([]domain.ItemBacklog, error) {
	return nil, f.listarErr
}

func (f *fakeRepositorioBacklog) ActualizarItem(string, domain.ItemBacklog) error {
	f.actualizaciones++
	return f.actualizarErr
}

func nuevoFakeProyectos() *fakeRepositorioProyectos {
	return &fakeRepositorioProyectos{proyectos: map[string]domain.Proyecto{}}
}

func nuevoFakeIntegrantes() *fakeRepositorioIntegrantes {
	return &fakeRepositorioIntegrantes{integrantes: map[string][]domain.Integrante{}}
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
