package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func nuevoEntornoBacklogConProyecto() *entornoBacklog {
	entorno := nuevoEntornoBacklog()
	entorno.proyectos["p1"] = proyectoFakeValido()
	return entorno
}

func TestT09_CrearItem_PropagaErrorRepositorioConWrap(t *testing.T) {
	causa := errors.New("fallo de repositorio")

	t.Run("Existe", func(t *testing.T) {
		entorno := nuevoEntornoBacklog()
		entorno.existeErr = causa
		servicio := entorno.servicio()

		item, err := servicio.CrearItem("p1", "Registrar proyecto", "", "Alta", nil)
		if !errors.Is(err, causa) {
			t.Fatalf("errors.Is(err, causa) = false; err = %v", err)
		}
		if err == causa {
			t.Fatalf("el error debe estar envuelto; err == causa")
		}
		if !strings.HasPrefix(err.Error(), "service:") {
			t.Fatalf("mensaje = %q; se esperaba prefijo service:", err.Error())
		}
		verificarItemBacklogVacio(t, item)
		if entorno.agregados != 0 {
			t.Fatalf("AgregarItem fue invocado %d veces", entorno.agregados)
		}
	})

	t.Run("AgregarItem", func(t *testing.T) {
		entorno := nuevoEntornoBacklogConProyecto()
		entorno.agregarErr = causa
		servicio := entorno.servicio()

		item, err := servicio.CrearItem("p1", "Registrar proyecto", "", "Alta", nil)
		if !errors.Is(err, causa) {
			t.Fatalf("errors.Is(err, causa) = false; err = %v", err)
		}
		if err == causa {
			t.Fatalf("el error debe estar envuelto; err == causa")
		}
		if !strings.HasPrefix(err.Error(), "service:") {
			t.Fatalf("mensaje = %q; se esperaba prefijo service:", err.Error())
		}
		verificarItemBacklogVacio(t, item)
	})
}

func TestT09_CrearItem_ConservaIDDevueltoPorRepositorio(t *testing.T) {
	entorno := nuevoEntornoBacklogConProyecto()
	entorno.idAsignado = "HU-77-DEVUELTO-POR-REPOSITORIO"
	servicio := entorno.servicio()

	creado, err := servicio.CrearItem("p1", "Registrar proyecto", "Descripción", "Alta", nil)
	if err != nil {
		t.Fatalf("CrearItem() devolvió un error inesperado: %v", err)
	}
	if creado.ID != entorno.idAsignado {
		t.Fatalf("ID = %q; se esperaba el ID devuelto por el repositorio %q", creado.ID, entorno.idAsignado)
	}

	persistidos := entorno.items["p1"]
	if len(persistidos) != 1 {
		t.Fatalf("se esperaba 1 ítem persistido, hay %d", len(persistidos))
	}
	if persistidos[0].ID != entorno.idAsignado {
		t.Fatalf("ID persistido = %q; se esperaba %q", persistidos[0].ID, entorno.idAsignado)
	}
}

func TestT09_CrearItem_NoEscribeSiProyectoNoExiste(t *testing.T) {
	entorno := nuevoEntornoBacklog()
	servicio := entorno.servicio()

	creado, err := servicio.CrearItem("ausente", "Registrar proyecto", "", "Alta", nil)
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("error = %v; se esperaba ErrProyectoNoEncontrado", err)
	}
	verificarItemBacklogVacio(t, creado)
	if entorno.agregados != 0 {
		t.Fatalf("AgregarItem fue invocado %d veces", entorno.agregados)
	}
}

func TestT09_CrearItem_NoEscribeSiValidacionFalla(t *testing.T) {
	casos := []struct {
		descripcion   string
		titulo        string
		prioridad     string
		errorEsperado error
	}{
		{descripcion: "título vacío", titulo: "", prioridad: "Alta", errorEsperado: domain.ErrTituloItemVacio},
		{descripcion: "prioridad inválida", titulo: "Registrar proyecto", prioridad: "Urgente", errorEsperado: domain.ErrPrioridadInvalida},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoBacklogConProyecto()
			servicio := entorno.servicio()

			creado, err := servicio.CrearItem("p1", caso.titulo, "", caso.prioridad, nil)
			if !errors.Is(err, caso.errorEsperado) {
				t.Fatalf("error = %v; se esperaba %v", err, caso.errorEsperado)
			}
			verificarItemBacklogVacio(t, creado)
			if entorno.agregados != 0 {
				t.Fatalf("AgregarItem fue invocado %d veces", entorno.agregados)
			}
		})
	}
}

func TestT09_CrearItem_ErrorDominioSinEnvolver(t *testing.T) {
	casos := []struct {
		descripcion   string
		proyectoID    string
		titulo        string
		prioridad     string
		errorEsperado error
	}{
		{descripcion: "proyecto inexistente", proyectoID: "ausente", titulo: "Registrar proyecto", prioridad: "Alta", errorEsperado: domain.ErrProyectoNoEncontrado},
		{descripcion: "título vacío", proyectoID: "p1", titulo: "", prioridad: "Alta", errorEsperado: domain.ErrTituloItemVacio},
		{descripcion: "prioridad inválida", proyectoID: "p1", titulo: "Registrar proyecto", prioridad: "Urgente", errorEsperado: domain.ErrPrioridadInvalida},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoBacklogConProyecto()
			servicio := entorno.servicio()

			_, err := servicio.CrearItem(caso.proyectoID, caso.titulo, "", caso.prioridad, nil)
			if !errors.Is(err, caso.errorEsperado) {
				t.Fatalf("error = %v; se esperaba %v", err, caso.errorEsperado)
			}
			if errors.Unwrap(err) != nil {
				t.Fatalf("error de dominio envuelto: %v", err)
			}
		})
	}
}
