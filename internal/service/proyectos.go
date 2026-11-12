package service

import (
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// ServicioProyectos orquesta los casos de uso de proyectos sobre el dominio.
type ServicioProyectos struct {
	proyectos   RepositorioProyectos
	integrantes RepositorioIntegrantes
}

// NuevoServicioProyectos crea el servicio con los repositorios indicados.
func NuevoServicioProyectos(proyectos RepositorioProyectos, integrantes RepositorioIntegrantes) *ServicioProyectos {
	return &ServicioProyectos{proyectos: proyectos, integrantes: integrantes}
}

// ModificarProyecto modifica el nombre, la fecha de inicio y la fecha de fin
// del proyecto identificado por id, conservando su estado actual.
func (s *ServicioProyectos) ModificarProyecto(
	id, nombre string,
	fechaInicio, fechaFin time.Time,
) (domain.Proyecto, error) {
	registro, err := s.proyectos.Obtener(id)
	if err != nil {
		return domain.Proyecto{}, envolverErrorRepositorio("obtener proyecto", err)
	}

	modificado, err := domain.ModificarProyecto(registro.Proyecto, nombre, fechaInicio, fechaFin)
	if err != nil {
		return domain.Proyecto{}, err
	}

	if err := s.proyectos.Actualizar(id, modificado); err != nil {
		return domain.Proyecto{}, envolverErrorRepositorio("actualizar proyecto", err)
	}

	return modificado, nil
}
