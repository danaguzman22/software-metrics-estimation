package service

import (
	"errors"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// R36 / D30 (US-12 FASE 2): sin rondas, RondaActual devuelve el cero de la
// ronda, false y nil.
func TestUS12_Fase2_ServicioPlanningPoker_RondaActualSinRondas(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	servicio := entorno.servicio()

	ronda, existe, err := servicio.RondaActual("p1", "HU-01")
	if err != nil {
		t.Fatalf("RondaActual() devolvió un error inesperado: %v", err)
	}
	if existe {
		t.Errorf("existe = true, se esperaba false para una historia sin rondas")
	}
	if ronda != (domain.RondaPoker{}) {
		t.Errorf("ronda = %+v, se esperaba domain.RondaPoker{}", ronda)
	}
}

// R36 / D30 (US-12 FASE 2): con rondas, RondaActual devuelve la de mayor
// número (la última de la lista ordenada por número ASC), true y nil.
func TestUS12_Fase2_ServicioPlanningPoker_RondaActualDevuelveLaDeMayorNumero(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 3, Estado: domain.EstadoRondaRevelada})
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada})
	sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 2, Estado: domain.EstadoRondaAbierta})
	servicio := entorno.servicio()

	ronda, existe, err := servicio.RondaActual("p1", "HU-01")
	if err != nil {
		t.Fatalf("RondaActual() devolvió un error inesperado: %v", err)
	}
	if !existe {
		t.Errorf("existe = false, se esperaba true")
	}
	if ronda.Numero != 3 {
		t.Errorf("ronda.Numero = %d, se esperaba 3 (la de mayor número)", ronda.Numero)
	}
	if ronda.Estado != domain.EstadoRondaRevelada {
		t.Errorf("ronda.Estado = %q, se esperaba %q", ronda.Estado, domain.EstadoRondaRevelada)
	}
}

// D30 (US-12 FASE 2): proyecto inexistente -> domain.ErrProyectoNoEncontrado,
// sin consultar historia ni rondas.
func TestUS12_Fase2_ServicioPlanningPoker_RondaActualProyectoNoEncontrado(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	servicio := entorno.servicio()

	ronda, existe, err := servicio.RondaActual("ausente", "HU-01")
	if err == nil {
		t.Fatal("RondaActual() con un proyecto inexistente debería devolver un error")
	}
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if existe {
		t.Errorf("existe = true, se esperaba false")
	}
	if ronda != (domain.RondaPoker{}) {
		t.Errorf("ronda = %+v, se esperaba domain.RondaPoker{}", ronda)
	}
	if entorno.fakeRepositorioBacklog.obtenerLlamadas != 0 {
		t.Errorf("ObtenerItem fue invocado %d veces; no debía consultarse la historia", entorno.fakeRepositorioBacklog.obtenerLlamadas)
	}
	if entorno.fakeRepositorioRondas.listarLlamadas != 0 {
		t.Errorf("ListarPorHistoria fue invocado %d veces; no debían listarse las rondas", entorno.fakeRepositorioRondas.listarLlamadas)
	}
}

// D30 (US-12 FASE 2): historia inexistente -> domain.ErrItemNoEncontrado
// (pasa sin envolver), sin listar las rondas.
func TestUS12_Fase2_ServicioPlanningPoker_RondaActualHistoriaNoEncontrada(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	servicio := entorno.servicio()

	ronda, existe, err := servicio.RondaActual("p1", "HU-99")
	if err == nil {
		t.Fatal("RondaActual() con una historia inexistente debería devolver un error")
	}
	if !errors.Is(err, domain.ErrItemNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrItemNoEncontrado)
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, no debería mapearse a ErrProyectoNoEncontrado", err)
	}
	if existe {
		t.Errorf("existe = true, se esperaba false")
	}
	if ronda != (domain.RondaPoker{}) {
		t.Errorf("ronda = %+v, se esperaba domain.RondaPoker{}", ronda)
	}
	if entorno.fakeRepositorioRondas.listarLlamadas != 0 {
		t.Errorf("ListarPorHistoria fue invocado %d veces; no debían listarse las rondas", entorno.fakeRepositorioRondas.listarLlamadas)
	}
}

// D30 (US-12 FASE 2): los errores de repositorio se envuelven con su
// contexto ("verificar proyecto", "obtener ítem del backlog",
// "listar rondas de la historia") y no se mapean a errores de dominio.
func TestUS12_Fase2_ServicioPlanningPoker_RondaActualErrorEnvuelto(t *testing.T) {
	errSimulado := errors.New("fallo simulado del repositorio")
	casos := []struct {
		descripcion string
		preparar    func(e *entornoPoker)
	}{
		{"error al verificar el proyecto", func(e *entornoPoker) { e.existeErr = errSimulado }},
		{"error al obtener la historia", func(e *entornoPoker) { e.fakeRepositorioBacklog.obtenerErr = errSimulado }},
		{"error al listar las rondas", func(e *entornoPoker) { e.fakeRepositorioRondas.listarErr = errSimulado }},
	}
	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoPokerSembrado(t)
			sembrarRondaPoker(t, entorno, "p1", "HU-01", domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta})
			caso.preparar(entorno)
			servicio := entorno.servicio()

			ronda, existe, err := servicio.RondaActual("p1", "HU-01")
			if err == nil {
				t.Fatal("RondaActual() con un error de repositorio debería devolver un error")
			}
			if !errors.Is(err, errSimulado) {
				t.Errorf("err = %v, se esperaba un error que envuelva %v", err, errSimulado)
			}
			if errors.Is(err, domain.ErrProyectoNoEncontrado) || errors.Is(err, domain.ErrItemNoEncontrado) {
				t.Errorf("err = %v, no debería mapearse a un error de dominio", err)
			}
			if existe {
				t.Errorf("existe = true, se esperaba false")
			}
			if ronda != (domain.RondaPoker{}) {
				t.Errorf("ronda = %+v, se esperaba domain.RondaPoker{}", ronda)
			}
		})
	}
}

// D31 (US-12 FASE 2): ListarIntegrantes devuelve los integrantes del proyecto
// en orden de inserción.
func TestUS12_Fase2_ServicioProyectos_ListarIntegrantes(t *testing.T) {
	servicio, _, integrantes := nuevoServicioProyectosConDatos(
		ProyectoRegistrado{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)},
	)
	integrantes.integrantes["p1"] = []domain.Integrante{
		{Nombre: "Victoria Caro Boldrini", Rol: "Tester"},
		{Nombre: "Dana Guzmán", Rol: "Desarrollador"},
		{Nombre: "Ana Paz", Rol: "Diseñadora"},
	}

	obtenidos, err := servicio.ListarIntegrantes("p1")
	if err != nil {
		t.Fatalf("ListarIntegrantes() devolvió un error inesperado: %v", err)
	}
	esperados := []domain.Integrante{
		{Nombre: "Victoria Caro Boldrini", Rol: "Tester"},
		{Nombre: "Dana Guzmán", Rol: "Desarrollador"},
		{Nombre: "Ana Paz", Rol: "Diseñadora"},
	}
	if len(obtenidos) != len(esperados) {
		t.Fatalf("se esperaban %d integrantes, se obtuvieron %d", len(esperados), len(obtenidos))
	}
	for i, esperado := range esperados {
		if obtenidos[i].Nombre != esperado.Nombre || obtenidos[i].Rol != esperado.Rol {
			t.Errorf("integrante[%d] = %+v, se esperaba %+v (orden de inserción)", i, obtenidos[i], esperado)
		}
	}
}

// D31 (US-12 FASE 2): proyecto inexistente -> domain.ErrProyectoNoEncontrado,
// sin consultar el repositorio de integrantes.
func TestUS12_Fase2_ServicioProyectos_ListarIntegrantesProyectoNoEncontrado(t *testing.T) {
	servicio, _, integrantes := nuevoServicioProyectosConDatos(
		ProyectoRegistrado{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)},
	)

	obtenidos, err := servicio.ListarIntegrantes("ausente")
	if err == nil {
		t.Fatal("ListarIntegrantes() con un proyecto inexistente debería devolver un error")
	}
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if obtenidos != nil {
		t.Errorf("integrantes = %+v, se esperaba nil", obtenidos)
	}
	if integrantes.listarLlamadas != 0 {
		t.Errorf("ListarPorProyecto fue invocado %d veces; no debía consultarse", integrantes.listarLlamadas)
	}
}

// D31 (US-12 FASE 2): los errores de repositorio se envuelven con su contexto
// ("verificar proyecto", "listar integrantes") y no se mapean a errores de
// dominio.
func TestUS12_Fase2_ServicioProyectos_ListarIntegrantesErrorEnvuelto(t *testing.T) {
	errSimulado := errors.New("fallo simulado del repositorio")
	casos := []struct {
		descripcion string
		preparar    func(proyectos *fakeRepositorioProyectos, integrantes *fakeRepositorioIntegrantes)
	}{
		{"error al verificar el proyecto", func(proyectos *fakeRepositorioProyectos, _ *fakeRepositorioIntegrantes) {
			proyectos.existeErr = errSimulado
		}},
		{"error al listar los integrantes", func(_ *fakeRepositorioProyectos, integrantes *fakeRepositorioIntegrantes) {
			integrantes.listarErr = errSimulado
		}},
	}
	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			proyectos := nuevoFakeProyectos()
			proyectos.proyectos["p1"] = proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)
			integrantes := nuevoFakeIntegrantes()
			servicio := NuevoServicioProyectos(proyectos, integrantes)
			caso.preparar(proyectos, integrantes)

			obtenidos, err := servicio.ListarIntegrantes("p1")
			if err == nil {
				t.Fatal("ListarIntegrantes() con un error de repositorio debería devolver un error")
			}
			if !errors.Is(err, errSimulado) {
				t.Errorf("err = %v, se esperaba un error que envuelva %v", err, errSimulado)
			}
			if errors.Is(err, domain.ErrProyectoNoEncontrado) {
				t.Errorf("err = %v, no debería mapearse a ErrProyectoNoEncontrado", err)
			}
			if obtenidos != nil {
				t.Errorf("integrantes = %+v, se esperaba nil", obtenidos)
			}
		})
	}
}
