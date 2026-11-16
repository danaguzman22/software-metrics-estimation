package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func TestT09_ModificarProyecto_PropagaErrorRepositorioConWrap(t *testing.T) {
	causa := errors.New("fallo al obtener proyecto")
	proyectos := nuevoFakeProyectos()
	proyectos.obtenerErr = fmt.Errorf("repositorio: %w", causa)
	servicio := NuevoServicioProyectos(proyectos, nuevoFakeIntegrantes())

	proyecto, err := servicio.ModificarProyecto(
		"p1",
		"Nuevo nombre",
		fechaFake(2026, 2, 1),
		fechaFake(2026, 3, 1),
	)
	if err == nil || !errors.Is(err, causa) {
		t.Fatalf("error = %v; se esperaba una causa envuelta", err)
	}
	if proyecto != (domain.Proyecto{}) {
		t.Fatalf("proyecto = %+v; se esperaba el valor cero", proyecto)
	}
}

func TestT09_ModificarProyecto_PropagaErrorDeActualizacionConWrap(t *testing.T) {
	causa := errors.New("fallo al actualizar proyecto")
	proyectos := nuevoFakeProyectos()
	proyectos.proyectos["p1"] = proyectoFakeValido()
	proyectos.actualizarErr = fmt.Errorf("repositorio: %w", causa)
	servicio := NuevoServicioProyectos(proyectos, nuevoFakeIntegrantes())

	proyecto, err := servicio.ModificarProyecto(
		"p1",
		"Nuevo nombre",
		fechaFake(2026, 2, 1),
		fechaFake(2026, 3, 1),
	)
	if err == nil || !errors.Is(err, causa) {
		t.Fatalf("error = %v; se esperaba una causa envuelta", err)
	}
	if proyecto != (domain.Proyecto{}) {
		t.Fatalf("proyecto = %+v; se esperaba el valor cero", proyecto)
	}
}

func TestT09_ModificarProyecto_NoActualizaSiValidacionFalla(t *testing.T) {
	proyectos := nuevoFakeProyectos()
	proyectos.proyectos["p1"] = proyectoFakeValido()
	servicio := NuevoServicioProyectos(proyectos, nuevoFakeIntegrantes())

	_, err := servicio.ModificarProyecto(
		"p1",
		"",
		fechaFake(2026, 2, 1),
		fechaFake(2026, 3, 1),
	)
	if !errors.Is(err, domain.ErrNombreVacio) {
		t.Fatalf("error = %v; se esperaba ErrNombreVacio", err)
	}
	if proyectos.actualizaciones != 0 {
		t.Fatalf("Actualizar fue invocado %d veces", proyectos.actualizaciones)
	}
}

func TestT09_ModificarProyecto_ErrorDominioSinEnvolver(t *testing.T) {
	proyectos := nuevoFakeProyectos()
	proyectos.proyectos["p1"] = proyectoFakeValido()
	servicio := NuevoServicioProyectos(proyectos, nuevoFakeIntegrantes())

	_, err := servicio.ModificarProyecto(
		"p1",
		"",
		fechaFake(2026, 2, 1),
		fechaFake(2026, 3, 1),
	)
	if !errors.Is(err, domain.ErrNombreVacio) {
		t.Fatalf("error = %v; se esperaba ErrNombreVacio", err)
	}
	if errors.Unwrap(err) != nil {
		t.Fatalf("error de dominio envuelto: %v", err)
	}
}
