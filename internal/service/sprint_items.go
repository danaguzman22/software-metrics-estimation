package service

import (
	"fmt"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// RepositorioAsignaciones persiste la relación historia–Sprint (D4, D5).
type RepositorioAsignaciones interface {
	// Agregar registra que la historia itemNumero pertenece a sprintNumero en
	// el proyecto. Devuelve domain.ErrItemYaAsignado si ya está asignada
	// (PK), domain.ErrSprintNoEncontrado o domain.ErrItemNoEncontrado si una
	// FK falla (defensa ante concurrencia, D7).
	Agregar(proyectoID string, sprintNumero, itemNumero int) error

	// ExisteItem informa si la historia ya está asignada a algún Sprint del
	// proyecto. La consulta es por proyecto (D2, R5).
	ExisteItem(proyectoID string, itemNumero int) (bool, error)

	// ListarNumerosItem devuelve los números de historia asignados al Sprint,
	// ordenados por número ascendente (CA3).
	ListarNumerosItem(proyectoID string, sprintNumero int) ([]int, error)

	// ExisteItemEnSprint informa si la historia está asignada al Sprint indicado.
	ExisteItemEnSprint(proyectoID string, sprintNumero, itemNumero int) (bool, error)
}

// AsignarItem asigna la historia itemID al Sprint sprintNumero del proyecto.
//
// Orden de evaluación: proyecto, Sprint, ítem, ítem ya asignado, ítem estimado
// y ítem pendiente. Devuelve nil si la asignación se persiste.
func (s *ServicioSprints) AsignarItem(proyectoID string, sprintNumero int, itemID string) error {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return envolverErrorRepositorio("verificar proyecto", err)
	}
	if !existe {
		return domain.ErrProyectoNoEncontrado
	}

	sprint, err := s.sprints.Obtener(proyectoID, sprintNumero)
	if err != nil {
		return envolverErrorRepositorio("verificar sprint", err)
	}
	if sprint.Estado == domain.EstadoSprintCerrado {
		return domain.ErrSprintCerrado
	}

	item, err := s.backlog.ObtenerItem(proyectoID, itemID)
	if err != nil {
		return envolverErrorRepositorio("obtener ítem del backlog", err)
	}

	numero, ok := domain.NumeroDeItem(item.ID)
	if !ok {
		return domain.ErrItemNoEncontrado
	}

	yaAsignado, err := s.asignaciones.ExisteItem(proyectoID, numero)
	if err != nil {
		return envolverErrorRepositorio("verificar asignación", err)
	}
	if yaAsignado {
		return domain.ErrItemYaAsignado
	}

	if err := domain.ValidarItemAsignable(item); err != nil {
		return err
	}

	if err := s.asignaciones.Agregar(proyectoID, sprintNumero, numero); err != nil {
		return envolverErrorRepositorio("asignar ítem al sprint", err)
	}

	return nil
}

// StoryPointsPlanificados devuelve la suma de story points de las historias
// asignadas al Sprint (CA3); 0 si no tiene historias.
func (s *ServicioSprints) StoryPointsPlanificados(proyectoID string, sprintNumero int) (int, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return 0, envolverErrorRepositorio("verificar proyecto", err)
	}
	if !existe {
		return 0, domain.ErrProyectoNoEncontrado
	}

	sprint, err := s.sprints.Obtener(proyectoID, sprintNumero)
	if err != nil {
		return 0, envolverErrorRepositorio("verificar sprint", err)
	}
	if sprint.Estado == domain.EstadoSprintCerrado {
		if sprint.SPPlanificados == nil {
			return 0, fmt.Errorf("service: el sprint cerrado %d no tiene story points planificados guardados", sprintNumero)
		}
		return *sprint.SPPlanificados, nil
	}

	numeros, err := s.asignaciones.ListarNumerosItem(proyectoID, sprintNumero)
	if err != nil {
		return 0, envolverErrorRepositorio("listar ítems del sprint", err)
	}

	items, err := s.backlog.ListarPorProyecto(proyectoID)
	if err != nil {
		return 0, envolverErrorRepositorio("listar ítems del proyecto", err)
	}

	asignados := make(map[int]struct{}, len(numeros))
	for _, numero := range numeros {
		asignados[numero] = struct{}{}
	}

	filtrados := make([]domain.ItemBacklog, 0, len(asignados))
	for _, item := range items {
		numero, ok := domain.NumeroDeItem(item.ID)
		if !ok {
			continue
		}
		if _, asignado := asignados[numero]; asignado {
			filtrados = append(filtrados, item)
		}
	}

	return domain.StoryPointsPlanificados(filtrados), nil
}

// MarcarItemCompletada cambia a Completada una historia Pendiente asignada al
// Sprint activo indicado.
func (s *ServicioSprints) MarcarItemCompletada(proyectoID string, sprintNumero int, itemID string) error {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return envolverErrorRepositorio("verificar proyecto", err)
	}
	if !existe {
		return domain.ErrProyectoNoEncontrado
	}

	sprint, err := s.sprints.Obtener(proyectoID, sprintNumero)
	if err != nil {
		return envolverErrorRepositorio("obtener sprint", err)
	}
	if sprint.Estado != domain.EstadoSprintActivo {
		return domain.ErrSprintCerrado
	}

	item, err := s.backlog.ObtenerItem(proyectoID, itemID)
	if err != nil {
		return envolverErrorRepositorio("obtener ítem del backlog", err)
	}
	numero, ok := domain.NumeroDeItem(item.ID)
	if !ok {
		return domain.ErrItemNoEncontrado
	}
	asignado, err := s.asignaciones.ExisteItemEnSprint(proyectoID, sprintNumero, numero)
	if err != nil {
		return envolverErrorRepositorio("verificar asignación al sprint", err)
	}
	if !asignado {
		return domain.ErrItemNoAsignadoAlSprint
	}

	completado, err := item.MarcarComoCompletada()
	if err != nil {
		return err
	}
	if err := s.backlog.ActualizarItem(proyectoID, completado); err != nil {
		return envolverErrorRepositorio("marcar ítem como completado", err)
	}
	return nil
}

// StoryPointsCompletados devuelve la suma de Story Points de las historias
// completadas asignadas al Sprint.
func (s *ServicioSprints) StoryPointsCompletados(proyectoID string, sprintNumero int) (int, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return 0, envolverErrorRepositorio("verificar proyecto", err)
	}
	if !existe {
		return 0, domain.ErrProyectoNoEncontrado
	}
	if _, err := s.sprints.Obtener(proyectoID, sprintNumero); err != nil {
		return 0, envolverErrorRepositorio("obtener sprint", err)
	}
	items, err := s.backlog.ListarPorProyecto(proyectoID)
	if err != nil {
		return 0, envolverErrorRepositorio("listar ítems del proyecto", err)
	}
	numeros, err := s.asignaciones.ListarNumerosItem(proyectoID, sprintNumero)
	if err != nil {
		return 0, envolverErrorRepositorio("listar ítems del sprint", err)
	}
	asignados := make(map[int]struct{}, len(numeros))
	for _, numero := range numeros {
		asignados[numero] = struct{}{}
	}
	total := 0
	for _, item := range items {
		if item.Estado != domain.EstadoCompletada {
			continue
		}
		numero, ok := domain.NumeroDeItem(item.ID)
		if !ok {
			continue
		}
		if _, ok := asignados[numero]; ok {
			total += item.StoryPoints
		}
	}
	return total, nil
}

// VelocidadSprint devuelve la suma de Story Points completados del Sprint.
func (s *ServicioSprints) VelocidadSprint(proyectoID string, sprintNumero int) (int, error) {
	return s.StoryPointsCompletados(proyectoID, sprintNumero)
}
