package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// CA1 / R29 / D19: ListarSprints devuelve los Sprints del proyecto en el orden
// del repositorio, sin recalcular ni sumar Story Points en el servicio.
func TestUS22_Fase2_ServicioSprints_Listar(t *testing.T) {
	entorno := nuevoEntornoSprints()
	sembrarProyectosSprint(t, entorno, "p1")

	cerrado := sprintUS07Servicio(1, "Cerrado")
	spPlanificados := 8
	cerrado.SPPlanificados = &spPlanificados
	activo := sprintUS07Servicio(2, "Activo")
	entorno.sprints["p1"] = []domain.Sprint{cerrado, activo}
	servicio := entorno.servicio()

	sprints, err := servicio.ListarSprints("p1")
	if err != nil {
		t.Fatalf("ListarSprints() devolvió un error inesperado: %v", err)
	}
	if len(sprints) != 2 {
		t.Fatalf("se esperaban 2 Sprints, se obtuvieron %d", len(sprints))
	}

	verificarSprintServicioUS07(t, sprints[0], cerrado)
	verificarSprintServicioUS07(t, sprints[1], activo)

	if entorno.existeLlamadas != 1 {
		t.Errorf("Existe fue invocado %d veces, se esperaba 1", entorno.existeLlamadas)
	}
	if entorno.fakeRepositorioSprints.listarLlamadas != 1 {
		t.Errorf("ListarPorProyecto fue invocado %d veces, se esperaba 1", entorno.fakeRepositorioSprints.listarLlamadas)
	}
	if entorno.fakeRepositorioSprints.obtenerLlamadas != 0 {
		t.Errorf("Obtener fue invocado %d veces; ListarSprints no debe consultar Sprint por número", entorno.fakeRepositorioSprints.obtenerLlamadas)
	}
}

// CA2 / R29 / V9: un proyecto inexistente devuelve domain.ErrProyectoNoEncontrado
// sin envolver y no consulta los Sprints.
func TestUS22_Fase2_ServicioSprints_ListarProyectoNoEncontrado(t *testing.T) {
	entorno := nuevoEntornoSprints()
	sembrarProyectosSprint(t, entorno, "p1")
	servicio := entorno.servicio()

	sprints, err := servicio.ListarSprints("ausente")
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if err != domain.ErrProyectoNoEncontrado {
		t.Errorf("el error de proyecto inexistente debe devolverse sin envolver; err = %v", err)
	}
	if len(sprints) != 0 {
		t.Errorf("se esperaba una lista vacía de Sprints, se obtuvieron %d", len(sprints))
	}
	if entorno.fakeRepositorioSprints.listarLlamadas != 0 {
		t.Errorf("ListarPorProyecto fue invocado %d veces; el proyecto debía validarse primero", entorno.fakeRepositorioSprints.listarLlamadas)
	}
}

// CA2 / R29: los errores de repositorio se propagan envueltos con %w y no se
// mapean a domain.ErrProyectoNoEncontrado.
func TestUS22_Fase2_ServicioSprints_ListarErrorEnvuelto(t *testing.T) {
	causa := errors.New("fallo de repositorio")

	casos := []struct {
		descripcion string
		sembrar     bool
		preparar    func(*entornoSprints)
		contexto    string
	}{
		{
			descripcion: "Existe",
			preparar:    func(e *entornoSprints) { e.existeErr = causa },
			contexto:    "verificar proyecto",
		},
		{
			descripcion: "ListarPorProyecto",
			sembrar:     true,
			preparar:    func(e *entornoSprints) { e.fakeRepositorioSprints.listarErr = causa },
			contexto:    "listar sprints del proyecto",
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoSprints()
			if caso.sembrar {
				sembrarProyectosSprint(t, entorno, "p1")
			}
			caso.preparar(entorno)
			servicio := entorno.servicio()

			sprints, err := servicio.ListarSprints("p1")
			if err == nil {
				t.Fatal("ListarSprints() debía devolver el error de repositorio")
			}
			if !errors.Is(err, causa) {
				t.Fatalf("errors.Is(err, causa) = false; err = %v", err)
			}
			if err == causa {
				t.Fatal("el error debe estar envuelto; err == causa")
			}
			if !strings.HasPrefix(err.Error(), "service:") {
				t.Errorf("mensaje = %q; se esperaba prefijo service:", err.Error())
			}
			if !strings.Contains(err.Error(), caso.contexto) {
				t.Errorf("mensaje = %q; se esperaba el contexto %q", err.Error(), caso.contexto)
			}
			if errors.Is(err, domain.ErrProyectoNoEncontrado) {
				t.Errorf("err = %v, no debería mapearse a ErrProyectoNoEncontrado", err)
			}
			if len(sprints) != 0 {
				t.Errorf("se esperaba una lista vacía de Sprints, se obtuvieron %d", len(sprints))
			}
		})
	}
}
