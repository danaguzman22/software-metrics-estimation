package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func TestT09_AgregarIntegrante_PropagaErrorRepositorioConWrap(t *testing.T) {
	causa := errors.New("fallo al listar integrantes")
	proyectos := nuevoFakeProyectos()
	proyectos.proyectos["p1"] = proyectoFakeValido()
	integrantes := nuevoFakeIntegrantes()
	integrantes.listarErr = fmt.Errorf("repositorio: %w", causa)
	servicio := NuevoServicioProyectos(proyectos, integrantes)

	integrante, err := servicio.AgregarIntegrante("p1", "Dana", "Developer")
	if err == nil || !errors.Is(err, causa) {
		t.Fatalf("error = %v; se esperaba una causa envuelta", err)
	}
	if integrante != (domain.Integrante{}) {
		t.Fatalf("integrante = %+v; se esperaba el valor cero", integrante)
	}
	if integrantes.agregados != 0 {
		t.Fatalf("Agregar fue invocado %d veces", integrantes.agregados)
	}
}

func TestT09_AgregarIntegrante_PropagaErrorDeExisteConWrap(t *testing.T) {
	causa := errors.New("fallo al verificar proyecto")
	proyectos := nuevoFakeProyectos()
	proyectos.existeErr = fmt.Errorf("repositorio: %w", causa)
	servicio := NuevoServicioProyectos(proyectos, nuevoFakeIntegrantes())

	integrante, err := servicio.AgregarIntegrante("p1", "Dana", "Developer")
	if err == nil || !errors.Is(err, causa) {
		t.Fatalf("error = %v; se esperaba una causa envuelta", err)
	}
	if integrante != (domain.Integrante{}) {
		t.Fatalf("integrante = %+v; se esperaba el valor cero", integrante)
	}
}

func TestT09_AgregarIntegrante_PropagaErrorDeAgregarConWrap(t *testing.T) {
	causa := errors.New("fallo al agregar integrante")
	proyectos := nuevoFakeProyectos()
	proyectos.proyectos["p1"] = proyectoFakeValido()
	integrantes := nuevoFakeIntegrantes()
	integrantes.agregarErr = fmt.Errorf("repositorio: %w", causa)
	servicio := NuevoServicioProyectos(proyectos, integrantes)

	integrante, err := servicio.AgregarIntegrante("p1", "Dana", "Developer")
	if err == nil || !errors.Is(err, causa) {
		t.Fatalf("error = %v; se esperaba una causa envuelta", err)
	}
	if integrante != (domain.Integrante{}) {
		t.Fatalf("integrante = %+v; se esperaba el valor cero", integrante)
	}
}

func TestT09_AgregarIntegrante_NoEscribeSiProyectoNoExiste(t *testing.T) {
	servicio := NuevoServicioProyectos(nuevoFakeProyectos(), nuevoFakeIntegrantes())

	_, err := servicio.AgregarIntegrante("inexistente", "Dana", "Developer")
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("error = %v; se esperaba ErrProyectoNoEncontrado", err)
	}
}

func TestT09_AgregarIntegrante_NoEscribeSiDuplicado(t *testing.T) {
	proyectos := nuevoFakeProyectos()
	proyectos.proyectos["p1"] = proyectoFakeValido()
	integrantes := nuevoFakeIntegrantes()
	integrantes.integrantes["p1"] = []domain.Integrante{{Nombre: "Dana", Rol: "Tester"}}
	servicio := NuevoServicioProyectos(proyectos, integrantes)

	_, err := servicio.AgregarIntegrante("p1", " dana ", "Developer")
	if !errors.Is(err, domain.ErrIntegranteDuplicado) {
		t.Fatalf("error = %v; se esperaba ErrIntegranteDuplicado", err)
	}
	if integrantes.agregados != 0 {
		t.Fatalf("Agregar fue invocado %d veces", integrantes.agregados)
	}
}

func TestT09_AgregarIntegrante_NoEscribeSiValidacionFalla(t *testing.T) {
	proyectos := nuevoFakeProyectos()
	proyectos.proyectos["p1"] = proyectoFakeValido()
	integrantes := nuevoFakeIntegrantes()
	servicio := NuevoServicioProyectos(proyectos, integrantes)

	_, err := servicio.AgregarIntegrante("p1", "Dana", "")
	if !errors.Is(err, domain.ErrRolIntegranteVacio) {
		t.Fatalf("error = %v; se esperaba ErrRolIntegranteVacio", err)
	}
	if integrantes.listarLlamadas != 0 || integrantes.agregados != 0 {
		t.Fatalf("se accedió al repositorio después de una validación fallida")
	}
}

func TestT09_AgregarIntegrante_ErrorDominioSinEnvolver(t *testing.T) {
	proyectos := nuevoFakeProyectos()
	proyectos.proyectos["p1"] = proyectoFakeValido()
	servicio := NuevoServicioProyectos(proyectos, nuevoFakeIntegrantes())

	_, err := servicio.AgregarIntegrante("p1", "", "Developer")
	if !errors.Is(err, domain.ErrNombreIntegranteVacio) {
		t.Fatalf("error = %v; se esperaba ErrNombreIntegranteVacio", err)
	}
	if errors.Unwrap(err) != nil {
		t.Fatalf("error de dominio envuelto: %v", err)
	}
}
