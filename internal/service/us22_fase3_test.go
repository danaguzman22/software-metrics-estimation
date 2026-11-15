package service

import (
	"errors"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// CA3 / R33 / R42 / D24: ListarItemsDelSprint devuelve solo las historias
// asignadas al Sprint, con sus datos intactos.
func TestUS22_Fase3_ServicioSprints_ListarItemsDelSprint(t *testing.T) {
	t.Run("devuelve solo las asignadas", func(t *testing.T) {
		entorno := nuevoEntornoAsignaciones()
		sembrarProyectosUS08(t, entorno, "p1")
		sembrarSprintUS08(t, entorno, "p1", 1)
		sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-01", 3, domain.EstadoItemPendiente))
		sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-02", 5, domain.EstadoCompletada))
		sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-03", 8, domain.EstadoItemPendiente))
		entorno.asignaciones.sembrar("p1", 1, 1)
		entorno.asignaciones.sembrar("p1", 1, 3)
		servicio := entorno.servicio()

		items, err := servicio.ListarItemsDelSprint("p1", 1)
		if err != nil {
			t.Fatalf("ListarItemsDelSprint() devolvió un error inesperado: %v", err)
		}
		if len(items) != 2 {
			t.Fatalf("se esperaban 2 historias asignadas, se obtuvieron %d: %+v", len(items), items)
		}
		porID := make(map[string]domain.ItemBacklog, len(items))
		for _, item := range items {
			porID[item.ID] = item
		}
		if _, ok := porID["HU-02"]; ok {
			t.Errorf("se devolvió HU-02, que no está asignada al Sprint 1")
		}
		if item, ok := porID["HU-01"]; !ok {
			t.Errorf("HU-01 está asignada al Sprint 1 y no se devolvió")
		} else if item.StoryPoints != 3 || item.Estado != domain.EstadoItemPendiente {
			t.Errorf("HU-01 = %+v, se esperaban 3 story points y estado Pendiente", item)
		}
		if item, ok := porID["HU-03"]; !ok {
			t.Errorf("HU-03 está asignada al Sprint 1 y no se devolvió")
		} else if item.StoryPoints != 8 || item.Estado != domain.EstadoItemPendiente {
			t.Errorf("HU-03 = %+v, se esperaban 8 story points y estado Pendiente", item)
		}
		if entorno.asignaciones.listarLlamadas != 1 {
			t.Errorf("ListarNumerosItem fue invocado %d veces, se esperaba 1", entorno.asignaciones.listarLlamadas)
		}
		if entorno.backlog.listarLlamadas != 1 {
			t.Errorf("ListarPorProyecto fue invocado %d veces, se esperaba 1", entorno.backlog.listarLlamadas)
		}
	})

	t.Run("Sprint sin historias", func(t *testing.T) {
		entorno := nuevoEntornoAsignaciones()
		sembrarProyectosUS08(t, entorno, "p1")
		sembrarSprintUS08(t, entorno, "p1", 1)
		servicio := entorno.servicio()

		items, err := servicio.ListarItemsDelSprint("p1", 1)
		if err != nil {
			t.Fatalf("ListarItemsDelSprint() devolvió un error inesperado: %v", err)
		}
		if len(items) != 0 {
			t.Errorf("se esperaba una lista vacía de historias, se obtuvieron %d: %+v", len(items), items)
		}
	})
}

// CA3 / R33 / D24: las historias devueltas respetan el orden ascendente por
// número de historia, aunque el backlog esté cargado desordenado.
func TestUS22_Fase3_ServicioSprints_ListarItemsDelSprintOrdenAscendente(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	sembrarProyectosUS08(t, entorno, "p1")
	sembrarSprintUS08(t, entorno, "p1", 1)
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-09", 13, domain.EstadoItemPendiente))
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-03", 3, domain.EstadoItemPendiente))
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-07", 8, domain.EstadoItemPendiente))
	entorno.asignaciones.sembrar("p1", 1, 9)
	entorno.asignaciones.sembrar("p1", 1, 3)
	entorno.asignaciones.sembrar("p1", 1, 7)
	servicio := entorno.servicio()

	items, err := servicio.ListarItemsDelSprint("p1", 1)
	if err != nil {
		t.Fatalf("ListarItemsDelSprint() devolvió un error inesperado: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("se esperaban 3 historias, se obtuvieron %d: %+v", len(items), items)
	}
	esperado := []string{"HU-03", "HU-07", "HU-09"}
	for i, id := range esperado {
		if items[i].ID != id {
			t.Errorf("items[%d].ID = %q, se esperaba %q; orden devuelto: %+v", i, items[i].ID, id, items)
		}
	}
}

// CA2 / R41 / V9 / D24: un proyecto inexistente devuelve
// domain.ErrProyectoNoEncontrado sin envolver y no consulta los Sprints ni las
// asignaciones.
func TestUS22_Fase3_ServicioSprints_ListarItemsDelSprintProyectoNoEncontrado(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	sembrarProyectosUS08(t, entorno, "p1")
	servicio := entorno.servicio()

	items, err := servicio.ListarItemsDelSprint("ausente", 1)
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if err != domain.ErrProyectoNoEncontrado {
		t.Errorf("el error de proyecto inexistente debe devolverse sin envolver; err = %v", err)
	}
	if len(items) != 0 {
		t.Errorf("se esperaba una lista vacía de historias, se obtuvieron %d: %+v", len(items), items)
	}
	if entorno.sprints.obtenerLlamadas != 0 {
		t.Errorf("Obtener fue invocado %d veces; el proyecto debía validarse primero", entorno.sprints.obtenerLlamadas)
	}
	if entorno.asignaciones.listarLlamadas != 0 {
		t.Errorf("ListarNumerosItem fue invocado %d veces; el proyecto debía validarse primero", entorno.asignaciones.listarLlamadas)
	}
	if entorno.backlog.listarLlamadas != 0 {
		t.Errorf("ListarPorProyecto fue invocado %d veces; el proyecto debía validarse primero", entorno.backlog.listarLlamadas)
	}
}

// CA2 / R41 / V9 / D24: un Sprint inexistente devuelve
// domain.ErrSprintNoEncontrado sin envolver y no lista las asignaciones.
func TestUS22_Fase3_ServicioSprints_ListarItemsDelSprintSprintNoEncontrado(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	sembrarProyectosUS08(t, entorno, "p1")
	servicio := entorno.servicio()

	items, err := servicio.ListarItemsDelSprint("p1", 1)
	if !errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrSprintNoEncontrado)
	}
	if err != domain.ErrSprintNoEncontrado {
		t.Errorf("el error de Sprint inexistente debe devolverse sin envolver; err = %v", err)
	}
	if len(items) != 0 {
		t.Errorf("se esperaba una lista vacía de historias, se obtuvieron %d: %+v", len(items), items)
	}
	if entorno.asignaciones.listarLlamadas != 0 {
		t.Errorf("ListarNumerosItem fue invocado %d veces; el Sprint debía existir antes de listar", entorno.asignaciones.listarLlamadas)
	}
	if entorno.backlog.listarLlamadas != 0 {
		t.Errorf("ListarPorProyecto fue invocado %d veces; el Sprint debía existir antes de listar", entorno.backlog.listarLlamadas)
	}
}

// CA2 / R42 / D24: los errores de repositorio se propagan envueltos y no se
// mapean a errores de dominio de no encontrado.
func TestUS22_Fase3_ServicioSprints_ListarItemsDelSprintErrorEnvuelto(t *testing.T) {
	causa := errors.New("fallo de repositorio")

	casos := []struct {
		descripcion string
		preparar    func(t *testing.T, entorno *entornoAsignaciones, causa error)
	}{
		{
			descripcion: "Existe",
			preparar: func(t *testing.T, entorno *entornoAsignaciones, causa error) {
				entorno.proyectos.existeErr = causa
			},
		},
		{
			descripcion: "Obtener",
			preparar: func(t *testing.T, entorno *entornoAsignaciones, causa error) {
				sembrarProyectosUS08(t, entorno, "p1")
				entorno.sprints.obtenerErr = causa
			},
		},
		{
			descripcion: "ListarNumerosItem",
			preparar: func(t *testing.T, entorno *entornoAsignaciones, causa error) {
				sembrarProyectosUS08(t, entorno, "p1")
				sembrarSprintUS08(t, entorno, "p1", 1)
				entorno.asignaciones.listarErr = causa
			},
		},
		{
			descripcion: "ListarPorProyecto",
			preparar: func(t *testing.T, entorno *entornoAsignaciones, causa error) {
				sembrarProyectosUS08(t, entorno, "p1")
				sembrarSprintUS08(t, entorno, "p1", 1)
				entorno.backlog.listarErr = causa
			},
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoAsignaciones()
			caso.preparar(t, entorno, causa)
			servicio := entorno.servicio()

			items, err := servicio.ListarItemsDelSprint("p1", 1)
			if len(items) != 0 {
				t.Errorf("se esperaba una lista vacía de historias, se obtuvieron %d: %+v", len(items), items)
			}
			verificarErrorRepositorioUS08(t, err, causa)
		})
	}
}
