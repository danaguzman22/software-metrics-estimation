package service

import (
	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// ServicioDefectos ejecuta los casos de uso de gestión de defectos.
type ServicioDefectos struct {
	proyectos RepositorioProyectos
	backlog   RepositorioBacklog
	sprints   RepositorioSprints
	defectos  RepositorioDefectos
}

// NuevoServicioDefectos crea el servicio con los repositorios indicados.
func NuevoServicioDefectos(
	proyectos RepositorioProyectos,
	backlog RepositorioBacklog,
	sprints RepositorioSprints,
	defectos RepositorioDefectos,
) *ServicioDefectos {
	return &ServicioDefectos{
		proyectos: proyectos,
		backlog:   backlog,
		sprints:   sprints,
		defectos:  defectos,
	}
}

// RegistrarDefecto valida las referencias dentro del proyecto y persiste el
// defecto.
func (s *ServicioDefectos) RegistrarDefecto(proyectoID string, entrada domain.Defecto) (domain.Defecto, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return domain.Defecto{}, envolverErrorRepositorio("verificar proyecto del defecto", err)
	}
	if !existe {
		return domain.Defecto{}, domain.ErrProyectoNoEncontrado
	}

	defecto, err := domain.NuevoDefecto(
		entrada.Descripcion,
		entrada.Severidad,
		entrada.Estado,
		entrada.ItemID,
		entrada.SprintDeteccion,
		entrada.SprintResolucion,
	)
	if err != nil {
		return domain.Defecto{}, err
	}
	if _, valido := domain.NumeroDeItem(defecto.ItemID); !valido {
		return domain.Defecto{}, domain.ErrItemNoEncontrado
	}
	if _, err := s.backlog.ObtenerItem(proyectoID, defecto.ItemID); err != nil {
		return domain.Defecto{}, envolverErrorRepositorio("obtener historia del defecto", err)
	}
	if _, err := s.sprints.Obtener(proyectoID, defecto.SprintDeteccion); err != nil {
		return domain.Defecto{}, envolverErrorRepositorio("obtener Sprint de detección del defecto", err)
	}
	if defecto.SprintResolucion != nil {
		if _, err := s.sprints.Obtener(proyectoID, *defecto.SprintResolucion); err != nil {
			return domain.Defecto{}, envolverErrorRepositorio("obtener Sprint de resolución del defecto", err)
		}
	}
	creado, err := s.defectos.Crear(proyectoID, defecto)
	if err != nil {
		return domain.Defecto{}, envolverErrorRepositorio("crear defecto", err)
	}
	return creado, nil
}

// ListarDefectos devuelve los defectos del proyecto en orden de creación.
func (s *ServicioDefectos) ListarDefectos(proyectoID string) ([]domain.Defecto, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return nil, envolverErrorRepositorio("verificar proyecto de defectos", err)
	}
	if !existe {
		return nil, domain.ErrProyectoNoEncontrado
	}
	defectos, err := s.defectos.ListarPorProyecto(proyectoID)
	if err != nil {
		return nil, envolverErrorRepositorio("listar defectos del proyecto", err)
	}
	if defectos == nil {
		defectos = []domain.Defecto{}
	}
	return defectos, nil
}

// CambiarEstadoDefecto actualiza el estado y, si corresponde, el Sprint de
// resolución en una sola operación.
func (s *ServicioDefectos) CambiarEstadoDefecto(
	proyectoID string,
	defectoID int64,
	estado string,
	sprintResolucion *int,
) (domain.Defecto, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return domain.Defecto{}, envolverErrorRepositorio("verificar proyecto del defecto", err)
	}
	if !existe {
		return domain.Defecto{}, domain.ErrProyectoNoEncontrado
	}
	defecto, err := s.defectos.Obtener(proyectoID, defectoID)
	if err != nil {
		return domain.Defecto{}, envolverErrorRepositorio("obtener defecto", err)
	}
	actualizado, err := defecto.CambiarEstado(estado, sprintResolucion)
	if err != nil {
		return domain.Defecto{}, err
	}
	if sprintResolucion != nil {
		if _, err := s.sprints.Obtener(proyectoID, *sprintResolucion); err != nil {
			return domain.Defecto{}, envolverErrorRepositorio("obtener Sprint de resolución del defecto", err)
		}
	}
	if err := s.defectos.Actualizar(proyectoID, actualizado); err != nil {
		return domain.Defecto{}, envolverErrorRepositorio("actualizar defecto", err)
	}
	return actualizado, nil
}

// MetricasDefectosPorSprint devuelve los conteos del Sprint luego de validar
// primero el proyecto y luego el Sprint dentro de ese proyecto.
func (s *ServicioDefectos) MetricasDefectosPorSprint(
	proyectoID string,
	sprintNumero int,
) (domain.MetricasDefectos, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return domain.MetricasDefectos{}, envolverErrorRepositorio("verificar proyecto de métricas de defectos", err)
	}
	if !existe {
		return domain.MetricasDefectos{}, domain.ErrProyectoNoEncontrado
	}
	if _, err := s.sprints.Obtener(proyectoID, sprintNumero); err != nil {
		return domain.MetricasDefectos{}, envolverErrorRepositorio("obtener sprint de métricas de defectos", err)
	}
	metricas, err := s.defectos.MetricasPorSprint(proyectoID, sprintNumero)
	if err != nil {
		return domain.MetricasDefectos{}, envolverErrorRepositorio("consultar métricas de defectos del sprint", err)
	}
	return metricas, nil
}

// MetricasDefectosPorProyecto devuelve los conteos agregados de defectos del
// proyecto.
func (s *ServicioDefectos) MetricasDefectosPorProyecto(proyectoID string) (domain.MetricasDefectos, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return domain.MetricasDefectos{}, envolverErrorRepositorio("verificar proyecto de métricas de defectos", err)
	}
	if !existe {
		return domain.MetricasDefectos{}, domain.ErrProyectoNoEncontrado
	}
	metricas, err := s.defectos.MetricasPorProyecto(proyectoID)
	if err != nil {
		return domain.MetricasDefectos{}, envolverErrorRepositorio("consultar métricas de defectos del proyecto", err)
	}
	return metricas, nil
}
