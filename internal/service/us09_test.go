package service

import (
	"errors"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func TestUS09_MarcarItemCompletada_Exitoso(t *testing.T) {
	entorno, servicio := entornoConHistoriaUS08(
		t,
		itemUS08Servicio("HU-01", 3, domain.EstadoItemPendiente),
	)
	entorno.asignaciones.sembrar("p1", 1, 1)

	if err := servicio.MarcarItemCompletada("p1", 1, "HU-01"); err != nil {
		t.Fatalf("MarcarItemCompletada() devolvió un error inesperado: %v", err)
	}
	item, err := entorno.backlog.ObtenerItem("p1", "HU-01")
	if err != nil {
		t.Fatalf("ObtenerItem() devolvió un error inesperado: %v", err)
	}
	if item.Estado != domain.EstadoCompletada {
		t.Errorf("Estado = %q, se esperaba %q", item.Estado, domain.EstadoCompletada)
	}
}

func TestUS09_MarcarItemCompletada_OrdenDeValidacion(t *testing.T) {
	t.Run("proyecto antes que Sprint e historia", func(t *testing.T) {
		entorno := nuevoEntornoAsignaciones()
		err := entorno.servicio().MarcarItemCompletada("p999", 9, "HU-99")
		if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
			t.Fatalf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
		}
		if entorno.sprints.obtenerLlamadas != 0 || entorno.backlog.obtenerLlamadas != 0 {
			t.Error("se consultó Sprint o historia antes de validar el proyecto")
		}
	})

	t.Run("Sprint antes que historia", func(t *testing.T) {
		entorno := nuevoEntornoAsignaciones()
		sembrarProyectosUS08(t, entorno, "p1")
		err := entorno.servicio().MarcarItemCompletada("p1", 9, "HU-99")
		if !errors.Is(err, domain.ErrSprintNoEncontrado) {
			t.Fatalf("err = %v, se esperaba %v", err, domain.ErrSprintNoEncontrado)
		}
		if entorno.backlog.obtenerLlamadas != 0 {
			t.Error("se consultó la historia antes de validar el Sprint")
		}
	})

	t.Run("Sprint cerrado antes que historia", func(t *testing.T) {
		entorno := nuevoEntornoAsignaciones()
		sembrarProyectosUS08(t, entorno, "p1")
		sembrarSprintUS08(t, entorno, "p1", 1)
		entorno.sprints.sprints["p1"][0].Estado = "Cerrado"
		err := entorno.servicio().MarcarItemCompletada("p1", 1, "HU-99")
		if !errors.Is(err, domain.ErrSprintCerrado) {
			t.Fatalf("err = %v, se esperaba %v", err, domain.ErrSprintCerrado)
		}
		if entorno.backlog.obtenerLlamadas != 0 {
			t.Error("se consultó la historia antes de validar el estado del Sprint")
		}
	})
}

func TestUS09_MarcarItemCompletada_ValidaAsignacionYEstado(t *testing.T) {
	t.Run("historia no asignada", func(t *testing.T) {
		entorno, servicio := entornoConHistoriaUS08(
			t,
			itemUS08Servicio("HU-01", 3, domain.EstadoItemPendiente),
		)
		err := servicio.MarcarItemCompletada("p1", 1, "HU-01")
		if !errors.Is(err, domain.ErrItemNoAsignadoAlSprint) {
			t.Fatalf("err = %v, se esperaba %v", err, domain.ErrItemNoAsignadoAlSprint)
		}
		if entorno.backlog.actualizaciones != 0 {
			t.Errorf("ActualizarItem fue invocado %d veces; no debía modificarse la historia", entorno.backlog.actualizaciones)
		}
	})

	t.Run("historia asignada pero no pendiente", func(t *testing.T) {
		entorno, servicio := entornoConHistoriaUS08(
			t,
			itemUS08Servicio("HU-01", 3, domain.EstadoCompletada),
		)
		entorno.asignaciones.sembrar("p1", 1, 1)
		err := servicio.MarcarItemCompletada("p1", 1, "HU-01")
		if !errors.Is(err, domain.ErrItemNoPendiente) {
			t.Fatalf("err = %v, se esperaba %v", err, domain.ErrItemNoPendiente)
		}
		if entorno.backlog.actualizaciones != 0 {
			t.Errorf("ActualizarItem fue invocado %d veces; no debía modificarse la historia", entorno.backlog.actualizaciones)
		}
	})
}

func TestUS09_CerrarSprint_CalculaPlanificadosYValidaProyectoSprint(t *testing.T) {
	t.Run("cierre exitoso", func(t *testing.T) {
		entorno := nuevoEntornoAsignaciones()
		sembrarProyectosUS08(t, entorno, "p1")
		sembrarSprintUS08(t, entorno, "p1", 1)
		sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-01", 3, domain.EstadoCompletada))
		sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-02", 5, domain.EstadoItemPendiente))
		entorno.asignaciones.sembrar("p1", 1, 1)
		entorno.asignaciones.sembrar("p1", 1, 2)

		err := entorno.servicio().CerrarSprint("p1", 1)
		if err != nil {
			t.Fatalf("CerrarSprint() devolvió un error inesperado: %v", err)
		}
		if entorno.sprints.cerrarLlamadas != 1 || entorno.sprints.ultimoSPPlanificados != 8 {
			t.Errorf("Cerrar invocado %d veces con %d SP planificados; se esperaba una vez con 8",
				entorno.sprints.cerrarLlamadas, entorno.sprints.ultimoSPPlanificados)
		}
	})

	t.Run("proyecto antes que Sprint", func(t *testing.T) {
		entorno := nuevoEntornoAsignaciones()
		err := entorno.servicio().CerrarSprint("p999", 1)
		if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
			t.Fatalf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
		}
		if entorno.sprints.obtenerLlamadas != 0 {
			t.Error("se consultó el Sprint antes de validar el proyecto")
		}
	})

	t.Run("Sprint inexistente", func(t *testing.T) {
		entorno := nuevoEntornoAsignaciones()
		sembrarProyectosUS08(t, entorno, "p1")
		err := entorno.servicio().CerrarSprint("p1", 9)
		if !errors.Is(err, domain.ErrSprintNoEncontrado) {
			t.Fatalf("err = %v, se esperaba %v", err, domain.ErrSprintNoEncontrado)
		}
		if entorno.sprints.cerrarLlamadas != 0 {
			t.Error("se intentó cerrar un Sprint inexistente")
		}
	})

	t.Run("Sprint ya cerrado", func(t *testing.T) {
		entorno := nuevoEntornoAsignaciones()
		sembrarProyectosUS08(t, entorno, "p1")
		sembrarSprintUS08(t, entorno, "p1", 1)
		entorno.sprints.sprints["p1"][0].Estado = "Cerrado"
		err := entorno.servicio().CerrarSprint("p1", 1)
		if !errors.Is(err, domain.ErrSprintCerrado) {
			t.Fatalf("err = %v, se esperaba %v", err, domain.ErrSprintCerrado)
		}
		if entorno.sprints.cerrarLlamadas != 0 {
			t.Error("se intentó cerrar un Sprint que ya estaba cerrado")
		}
	})

	t.Run("cierre vacío guarda cero", func(t *testing.T) {
		entorno := nuevoEntornoAsignaciones()
		sembrarProyectosUS08(t, entorno, "p1")
		sembrarSprintUS08(t, entorno, "p1", 1)

		if err := entorno.servicio().CerrarSprint("p1", 1); err != nil {
			t.Fatalf("CerrarSprint() devolvió un error inesperado: %v", err)
		}
		if entorno.sprints.ultimoSPPlanificados != 0 {
			t.Errorf("spPlanificados = %d, se esperaba 0", entorno.sprints.ultimoSPPlanificados)
		}
	})

	t.Run("historias completadas también se cuentan como planificadas", func(t *testing.T) {
		entorno := nuevoEntornoAsignaciones()
		sembrarProyectosUS08(t, entorno, "p1")
		sembrarSprintUS08(t, entorno, "p1", 1)
		sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-01", 3, domain.EstadoCompletada))
		sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-02", 5, domain.EstadoCompletada))
		entorno.asignaciones.sembrar("p1", 1, 1)
		entorno.asignaciones.sembrar("p1", 1, 2)

		if err := entorno.servicio().CerrarSprint("p1", 1); err != nil {
			t.Fatalf("CerrarSprint() devolvió un error inesperado: %v", err)
		}
		if entorno.sprints.ultimoSPPlanificados != 8 {
			t.Errorf("spPlanificados = %d, se esperaba 8", entorno.sprints.ultimoSPPlanificados)
		}
	})
}

func TestUS09_StoryPointsPlanificados_UsaValorGuardadoAlCerrar(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	sembrarProyectosUS08(t, entorno, "p1")
	sembrarSprintUS08(t, entorno, "p1", 1)
	valorGuardado := 8
	entorno.sprints.sprints["p1"][0].Estado = domain.EstadoSprintCerrado
	entorno.sprints.sprints["p1"][0].SPPlanificados = &valorGuardado
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-01", 3, domain.EstadoCompletada))
	entorno.asignaciones.sembrar("p1", 1, 1)

	total, err := entorno.servicio().StoryPointsPlanificados("p1", 1)
	if err != nil {
		t.Fatalf("StoryPointsPlanificados() devolvió un error inesperado: %v", err)
	}
	if total != valorGuardado {
		t.Errorf("StoryPointsPlanificados() = %d, se esperaba el valor guardado %d", total, valorGuardado)
	}
	if entorno.asignaciones.listarLlamadas != 0 || entorno.backlog.listarLlamadas != 0 {
		t.Error("se consultaron asignaciones o backlog para un Sprint cerrado")
	}
}

func TestUS09_StoryPointsPlanificados_SprintActivoSeCalculaEnVivo(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	sembrarProyectosUS08(t, entorno, "p1")
	sembrarSprintUS08(t, entorno, "p1", 1)
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-01", 3, domain.EstadoItemPendiente))
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-02", 5, domain.EstadoCompletada))
	entorno.asignaciones.sembrar("p1", 1, 1)
	entorno.asignaciones.sembrar("p1", 1, 2)

	total, err := entorno.servicio().StoryPointsPlanificados("p1", 1)
	if err != nil {
		t.Fatalf("StoryPointsPlanificados() devolvió un error inesperado: %v", err)
	}
	if total != 8 {
		t.Errorf("StoryPointsPlanificados() = %d, se esperaba 8", total)
	}
}

func TestUS09_ConsultasMetricas_ValidanProyectoAntesQueSprint(t *testing.T) {
	consultas := map[string]func(*ServicioSprints) (int, error){
		"planificados": func(s *ServicioSprints) (int, error) { return s.StoryPointsPlanificados("p1", 9) },
		"completados":  func(s *ServicioSprints) (int, error) { return s.StoryPointsCompletados("p1", 9) },
		"velocidad":    func(s *ServicioSprints) (int, error) { return s.VelocidadSprint("p1", 9) },
	}

	for nombre, consultar := range consultas {
		t.Run(nombre+"/proyecto inexistente", func(t *testing.T) {
			entorno := nuevoEntornoAsignaciones()
			_, err := consultar(entorno.servicio())
			if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
			}
			if entorno.sprints.obtenerLlamadas != 0 {
				t.Error("se consultó el Sprint antes de validar el proyecto")
			}
		})

		t.Run(nombre+"/Sprint inexistente", func(t *testing.T) {
			entorno := nuevoEntornoAsignaciones()
			sembrarProyectosUS08(t, entorno, "p1")
			_, err := consultar(entorno.servicio())
			if !errors.Is(err, domain.ErrSprintNoEncontrado) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrSprintNoEncontrado)
			}
		})
	}
}

func TestUS09_StoryPointsCompletadosYVelocidadSprint(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	sembrarProyectosUS08(t, entorno, "p1")
	sembrarSprintUS08(t, entorno, "p1", 1)
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-01", 3, domain.EstadoCompletada))
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-02", 5, domain.EstadoItemPendiente))
	entorno.asignaciones.sembrar("p1", 1, 1)
	entorno.asignaciones.sembrar("p1", 1, 2)

	servicio := entorno.servicio()
	completados, err := servicio.StoryPointsCompletados("p1", 1)
	if err != nil {
		t.Fatalf("StoryPointsCompletados() devolvió un error inesperado: %v", err)
	}
	if completados != 3 {
		t.Errorf("StoryPointsCompletados() = %d, se esperaba 3", completados)
	}
	velocidad, err := servicio.VelocidadSprint("p1", 1)
	if err != nil {
		t.Fatalf("VelocidadSprint() devolvió un error inesperado: %v", err)
	}
	if velocidad != 3 {
		t.Errorf("VelocidadSprint() = %d, se esperaba 3", velocidad)
	}
}

func TestUS09_AsignarItem_RechazaSprintCerradoAntesDeConsultarHistoria(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	sembrarProyectosUS08(t, entorno, "p1")
	sembrarSprintUS08(t, entorno, "p1", 1)
	entorno.sprints.sprints["p1"][0].Estado = "Cerrado"
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-01", 3, domain.EstadoItemPendiente))

	err := entorno.servicio().AsignarItem("p1", 1, "HU-01")
	if !errors.Is(err, domain.ErrSprintCerrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrSprintCerrado)
	}
	if entorno.backlog.obtenerLlamadas != 0 {
		t.Error("se consultó la historia antes de rechazar el Sprint cerrado")
	}
	if entorno.asignaciones.agregarLlamadas != 0 {
		t.Error("se intentó agregar una asignación a un Sprint cerrado")
	}
}
