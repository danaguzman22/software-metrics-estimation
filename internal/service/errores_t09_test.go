package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func TestT09_ErrorSimpleDeRepositorioSeEnvuelve(t *testing.T) {
	causa := errors.New("falla")

	t.Run("ModificarProyecto/Obtener", func(t *testing.T) {
		proyectos := nuevoFakeProyectos()
		proyectos.obtenerErr = causa
		servicio := NuevoServicioProyectos(proyectos, nuevoFakeIntegrantes())

		_, err := servicio.ModificarProyecto("p1", "Nombre", fechaFake(2026, 2, 1), fechaFake(2026, 3, 1))
		if !errors.Is(err, causa) {
			t.Fatalf("errors.Is(err, causa) = false; err = %v", err)
		}
		if err == causa {
			t.Fatalf("err debe estar envuelto; err == causa")
		}
		if !strings.HasPrefix(err.Error(), "service:") {
			t.Fatalf("mensaje = %q; se esperaba prefijo service:", err.Error())
		}
	})

	t.Run("ModificarProyecto/Actualizar", func(t *testing.T) {
		proyectos := nuevoFakeProyectos()
		proyectos.proyectos["p1"] = proyectoFakeValido()
		proyectos.actualizarErr = causa
		servicio := NuevoServicioProyectos(proyectos, nuevoFakeIntegrantes())

		_, err := servicio.ModificarProyecto("p1", "Nuevo nombre", fechaFake(2026, 2, 1), fechaFake(2026, 3, 1))
		if !errors.Is(err, causa) {
			t.Fatalf("errors.Is(err, causa) = false; err = %v", err)
		}
		if err == causa {
			t.Fatalf("err debe estar envuelto; err == causa")
		}
		if !strings.HasPrefix(err.Error(), "service:") {
			t.Fatalf("mensaje = %q; se esperaba prefijo service:", err.Error())
		}
	})

	t.Run("AgregarIntegrante/Existe", func(t *testing.T) {
		proyectos := nuevoFakeProyectos()
		proyectos.existeErr = causa
		servicio := NuevoServicioProyectos(proyectos, nuevoFakeIntegrantes())

		_, err := servicio.AgregarIntegrante("p1", "Ana", "Desarrollador")
		if !errors.Is(err, causa) {
			t.Fatalf("errors.Is(err, causa) = false; err = %v", err)
		}
		if err == causa {
			t.Fatalf("err debe estar envuelto; err == causa")
		}
		if !strings.HasPrefix(err.Error(), "service:") {
			t.Fatalf("mensaje = %q; se esperaba prefijo service:", err.Error())
		}
	})

	t.Run("AgregarIntegrante/ListarPorProyecto", func(t *testing.T) {
		proyectos := nuevoFakeProyectos()
		proyectos.proyectos["p1"] = proyectoFakeValido()
		integrantes := nuevoFakeIntegrantes()
		integrantes.listarErr = causa
		servicio := NuevoServicioProyectos(proyectos, integrantes)

		_, err := servicio.AgregarIntegrante("p1", "Ana", "Desarrollador")
		if !errors.Is(err, causa) {
			t.Fatalf("errors.Is(err, causa) = false; err = %v", err)
		}
		if err == causa {
			t.Fatalf("err debe estar envuelto; err == causa")
		}
		if !strings.HasPrefix(err.Error(), "service:") {
			t.Fatalf("mensaje = %q; se esperaba prefijo service:", err.Error())
		}
	})

	t.Run("AgregarIntegrante/Agregar", func(t *testing.T) {
		proyectos := nuevoFakeProyectos()
		proyectos.proyectos["p1"] = proyectoFakeValido()
		integrantes := nuevoFakeIntegrantes()
		integrantes.agregarErr = causa
		servicio := NuevoServicioProyectos(proyectos, integrantes)

		_, err := servicio.AgregarIntegrante("p1", "Ana", "Desarrollador")
		if !errors.Is(err, causa) {
			t.Fatalf("errors.Is(err, causa) = false; err = %v", err)
		}
		if err == causa {
			t.Fatalf("err debe estar envuelto; err == causa")
		}
		if !strings.HasPrefix(err.Error(), "service:") {
			t.Fatalf("mensaje = %q; se esperaba prefijo service:", err.Error())
		}
	})
}

func TestT09_ErroresDeDominioDelRepositorioSeDevuelvenTalCual(t *testing.T) {
	t.Run("ModificarProyecto/Obtener/ErrProyectoNoEncontrado", func(t *testing.T) {
		proyectos := nuevoFakeProyectos()
		proyectos.obtenerErr = domain.ErrProyectoNoEncontrado
		servicio := NuevoServicioProyectos(proyectos, nuevoFakeIntegrantes())

		_, err := servicio.ModificarProyecto("p1", "Nombre", fechaFake(2026, 2, 1), fechaFake(2026, 3, 1))
		if err != domain.ErrProyectoNoEncontrado {
			t.Fatalf("err = %v; se esperaba sentinel idéntico ErrProyectoNoEncontrado", err)
		}
	})

	t.Run("ModificarProyecto/Actualizar/ErrProyectoNoEncontrado", func(t *testing.T) {
		proyectos := nuevoFakeProyectos()
		proyectos.proyectos["p1"] = proyectoFakeValido()
		proyectos.actualizarErr = domain.ErrProyectoNoEncontrado
		servicio := NuevoServicioProyectos(proyectos, nuevoFakeIntegrantes())

		_, err := servicio.ModificarProyecto("p1", "Nuevo nombre", fechaFake(2026, 2, 1), fechaFake(2026, 3, 1))
		if err != domain.ErrProyectoNoEncontrado {
			t.Fatalf("err = %v; se esperaba sentinel idéntico ErrProyectoNoEncontrado", err)
		}
	})

	t.Run("AgregarIntegrante/Agregar/ErrIntegranteDuplicado", func(t *testing.T) {
		proyectos := nuevoFakeProyectos()
		proyectos.proyectos["p1"] = proyectoFakeValido()
		integrantes := nuevoFakeIntegrantes()
		integrantes.agregarErr = domain.ErrIntegranteDuplicado
		servicio := NuevoServicioProyectos(proyectos, integrantes)

		_, err := servicio.AgregarIntegrante("p1", "Ana", "Desarrollador")
		if err != domain.ErrIntegranteDuplicado {
			t.Fatalf("err = %v; se esperaba sentinel idéntico ErrIntegranteDuplicado", err)
		}
	})
}
