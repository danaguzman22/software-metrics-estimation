package service

import (
	"errors"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// CA1 / R7 / D13: crear un proyecto válido lo persiste y devuelve el registro
// con el ID generado por el repositorio.
func TestUS22_Fase1_ServicioProyectos_Crear(t *testing.T) {
	t.Run("crea y persiste el proyecto", func(t *testing.T) {
		servicio, proyectos, _ := nuevoServicioProyectosConDatos()
		fechaInicio := fechaUS02(2026, time.March, 1)
		fechaFin := fechaUS02(2026, time.June, 30)

		creado, err := servicio.CrearProyecto("Sistema de Métricas", fechaInicio, fechaFin)
		if err != nil {
			t.Fatalf("CrearProyecto() devolvió un error inesperado: %v", err)
		}
		if creado.ID == "" {
			t.Fatal("CrearProyecto() debería devolver un ID no vacío")
		}
		if proyectos.crearLlamadas != 1 {
			t.Errorf("Crear fue invocado %d veces, se esperaba 1", proyectos.crearLlamadas)
		}

		esperado := domain.Proyecto{
			Nombre:      "Sistema de Métricas",
			FechaInicio: fechaInicio,
			FechaFin:    fechaFin,
			Estado:      domain.EstadoPlanificado,
		}
		verificarProyecto(t, creado.Proyecto, esperado)
		verificarProyecto(t, buscarProyectoEnFake(t, proyectos, creado.ID), esperado)
	})

	casos := []struct {
		descripcion   string
		nombre        string
		fechaInicio   time.Time
		fechaFin      time.Time
		errorEsperado error
		sinEnvolver   bool
	}{
		{
			descripcion:   "nombre vacío",
			nombre:        "",
			fechaInicio:   fechaUS02(2026, time.March, 1),
			fechaFin:      fechaUS02(2026, time.June, 30),
			errorEsperado: domain.ErrNombreVacio,
			sinEnvolver:   true,
		},
		{
			descripcion:   "fecha de fin anterior a la de inicio",
			nombre:        "Sistema de Métricas",
			fechaInicio:   fechaUS02(2026, time.June, 30),
			fechaFin:      fechaUS02(2026, time.March, 1),
			errorEsperado: domain.ErrFechaFinAnteriorAInicio,
			sinEnvolver:   true,
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			servicio, proyectos, _ := nuevoServicioProyectosConDatos()

			creado, err := servicio.CrearProyecto(caso.nombre, caso.fechaInicio, caso.fechaFin)
			if err == nil {
				t.Fatal("CrearProyecto() con datos inválidos debería devolver un error")
			}
			if !errors.Is(err, caso.errorEsperado) {
				t.Errorf("err = %v, se esperaba %v", err, caso.errorEsperado)
			}
			if caso.sinEnvolver && err != caso.errorEsperado {
				t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, caso.errorEsperado)
			}
			if creado.ID != "" || creado.Proyecto != (domain.Proyecto{}) {
				t.Errorf("se esperaba ProyectoRegistrado{} ante un error, se obtuvo %+v", creado)
			}
			if proyectos.crearLlamadas != 0 {
				t.Errorf("Crear fue invocado %d veces ante un error de validación; no debía escribirse nada", proyectos.crearLlamadas)
			}
		})
	}
}

// CA1 / R5 / R8: listar devuelve todos los proyectos registrados.
func TestUS22_Fase1_ServicioProyectos_Listar(t *testing.T) {
	servicio, _, _ := nuevoServicioProyectosConDatos(
		ProyectoRegistrado{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)},
		ProyectoRegistrado{ID: "p2", Proyecto: proyectoUS02Servicio("App Móvil", domain.Estado("En curso"))},
	)

	obtenidos, err := servicio.ListarProyectos()
	if err != nil {
		t.Fatalf("ListarProyectos() devolvió un error inesperado: %v", err)
	}
	if len(obtenidos) != 2 {
		t.Fatalf("len(ListarProyectos()) = %d, se esperaban 2", len(obtenidos))
	}

	porID := make(map[string]ProyectoRegistrado, len(obtenidos))
	for _, registro := range obtenidos {
		porID[registro.ID] = registro
	}
	p1, ok := porID["p1"]
	if !ok {
		t.Fatalf("ListarProyectos() no devolvió el proyecto p1: %+v", obtenidos)
	}
	verificarProyecto(t, p1.Proyecto, proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado))
	p2, ok := porID["p2"]
	if !ok {
		t.Fatalf("ListarProyectos() no devolvió el proyecto p2: %+v", obtenidos)
	}
	verificarProyecto(t, p2.Proyecto, proyectoUS02Servicio("App Móvil", domain.Estado("En curso")))
}

// CA1 / R5 / R8: obtener devuelve el proyecto identificado por id y propaga
// domain.ErrProyectoNoEncontrado cuando no existe.
func TestUS22_Fase1_ServicioProyectos_Obtener(t *testing.T) {
	t.Run("devuelve el proyecto", func(t *testing.T) {
		original := proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))
		servicio, _, _ := nuevoServicioProyectosConDatos(ProyectoRegistrado{ID: "p1", Proyecto: original})

		obtenido, err := servicio.ObtenerProyecto("p1")
		if err != nil {
			t.Fatalf("ObtenerProyecto() devolvió un error inesperado: %v", err)
		}
		if obtenido.ID != "p1" {
			t.Errorf("ID = %q, se esperaba %q", obtenido.ID, "p1")
		}
		verificarProyecto(t, obtenido.Proyecto, original)
	})

	t.Run("proyecto inexistente", func(t *testing.T) {
		servicio, _, _ := nuevoServicioProyectosConDatos(
			ProyectoRegistrado{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)},
		)

		obtenido, err := servicio.ObtenerProyecto("ausente")
		if err == nil {
			t.Fatal("ObtenerProyecto() con un ID inexistente debería devolver un error")
		}
		if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
			t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
		}
		if err != domain.ErrProyectoNoEncontrado {
			t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrProyectoNoEncontrado)
		}
		if obtenido.ID != "" || obtenido.Proyecto != (domain.Proyecto{}) {
			t.Errorf("se esperaba ProyectoRegistrado{} ante un error, se obtuvo %+v", obtenido)
		}
	})
}

// CA1 / CA2 / R11 / R21: los errores de repositorio se envuelven con el error
// original y no se mapean a errores de dominio.
func TestUS22_Fase1_ServicioProyectos_ErrorDeRepositorioEnvuelto(t *testing.T) {
	t.Run("crear proyecto", func(t *testing.T) {
		errSimulado := errors.New("fallo simulado del repositorio")
		servicio, proyectos, _ := nuevoServicioProyectosConDatos()
		proyectos.crearErr = errSimulado

		creado, err := servicio.CrearProyecto("Sistema de Métricas", fechaUS02(2026, time.March, 1), fechaUS02(2026, time.June, 30))
		if err == nil {
			t.Fatal("CrearProyecto() con un error de repositorio debería devolver un error")
		}
		if !errors.Is(err, errSimulado) {
			t.Errorf("err = %v, se esperaba un error que envuelva %v", err, errSimulado)
		}
		if errors.Is(err, domain.ErrProyectoNoEncontrado) {
			t.Errorf("err = %v, no debería mapearse a ErrProyectoNoEncontrado", err)
		}
		if creado.ID != "" || creado.Proyecto != (domain.Proyecto{}) {
			t.Errorf("se esperaba ProyectoRegistrado{} ante un error, se obtuvo %+v", creado)
		}
	})

	t.Run("listar proyectos", func(t *testing.T) {
		errSimulado := errors.New("fallo simulado del repositorio")
		servicio, proyectos, _ := nuevoServicioProyectosConDatos()
		proyectos.listarErr = errSimulado

		obtenidos, err := servicio.ListarProyectos()
		if err == nil {
			t.Fatal("ListarProyectos() con un error de repositorio debería devolver un error")
		}
		if !errors.Is(err, errSimulado) {
			t.Errorf("err = %v, se esperaba un error que envuelva %v", err, errSimulado)
		}
		if errors.Is(err, domain.ErrProyectoNoEncontrado) {
			t.Errorf("err = %v, no debería mapearse a ErrProyectoNoEncontrado", err)
		}
		if len(obtenidos) != 0 {
			t.Errorf("len(ListarProyectos()) = %d, se esperaba 0", len(obtenidos))
		}
	})

	t.Run("obtener proyecto", func(t *testing.T) {
		errSimulado := errors.New("fallo simulado del repositorio")
		servicio, proyectos, _ := nuevoServicioProyectosConDatos()
		proyectos.obtenerErr = errSimulado

		obtenido, err := servicio.ObtenerProyecto("p1")
		if err == nil {
			t.Fatal("ObtenerProyecto() con un error de repositorio debería devolver un error")
		}
		if !errors.Is(err, errSimulado) {
			t.Errorf("err = %v, se esperaba un error que envuelva %v", err, errSimulado)
		}
		if errors.Is(err, domain.ErrProyectoNoEncontrado) {
			t.Errorf("err = %v, no debería mapearse a ErrProyectoNoEncontrado", err)
		}
		if obtenido.ID != "" || obtenido.Proyecto != (domain.Proyecto{}) {
			t.Errorf("se esperaba ProyectoRegistrado{} ante un error, se obtuvo %+v", obtenido)
		}
	})
}

// CA1 / R17 / D12: obtener un ítem valida el proyecto y devuelve el ítem.
func TestUS22_Fase1_ServicioBacklog_ObtenerItem(t *testing.T) {
	original := itemUS06ServicioBase()
	almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p1", Item: original})
	servicio := almacen.servicio()

	obtenido, err := servicio.ObtenerItem("p1", "HU-01")
	if err != nil {
		t.Fatalf("ObtenerItem() devolvió un error inesperado: %v", err)
	}
	verificarItemBacklog(t, obtenido, original)
	if almacen.fakeRepositorioProyectos.existeLlamadas != 1 {
		t.Errorf("Existe fue invocado %d veces, se esperaba 1", almacen.fakeRepositorioProyectos.existeLlamadas)
	}
	if almacen.fakeRepositorioBacklog.obtenerLlamadas != 1 {
		t.Errorf("ObtenerItem del repositorio fue invocado %d veces, se esperaba 1", almacen.fakeRepositorioBacklog.obtenerLlamadas)
	}
}

// CA2 / R11 / R21: un proyecto inexistente devuelve
// domain.ErrProyectoNoEncontrado antes de consultar el backlog.
func TestUS22_Fase1_ServicioBacklog_ObtenerItemProyectoNoEncontrado(t *testing.T) {
	almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p1", Item: itemUS06ServicioBase()})
	servicio := almacen.servicio()

	obtenido, err := servicio.ObtenerItem("ausente", "HU-01")
	if err == nil {
		t.Fatal("ObtenerItem() con un proyecto inexistente debería devolver un error")
	}
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if err != domain.ErrProyectoNoEncontrado {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrProyectoNoEncontrado)
	}
	verificarItemBacklogVacio(t, obtenido)
	if almacen.fakeRepositorioBacklog.obtenerLlamadas != 0 {
		t.Errorf("ObtenerItem del backlog fue invocado %d veces; el proyecto debía validarse primero", almacen.fakeRepositorioBacklog.obtenerLlamadas)
	}
}

// CA1 / CA2 / R21: los errores de repositorio de ObtenerItem se envuelven con
// el error original y no se mapean a errores de dominio.
func TestUS22_Fase1_ServicioBacklog_ObtenerItemErrorEnvuelto(t *testing.T) {
	casos := []struct {
		descripcion string
		preparar    func(almacen *entornoBacklog, err error)
	}{
		{
			descripcion: "error al verificar el proyecto",
			preparar: func(almacen *entornoBacklog, err error) {
				almacen.fakeRepositorioProyectos.existeErr = err
			},
		},
		{
			descripcion: "error al obtener el ítem",
			preparar: func(almacen *entornoBacklog, err error) {
				almacen.fakeRepositorioBacklog.obtenerErr = err
			},
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			errSimulado := errors.New("fallo simulado del repositorio")
			almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p1", Item: itemUS06ServicioBase()})
			caso.preparar(almacen, errSimulado)
			servicio := almacen.servicio()

			obtenido, err := servicio.ObtenerItem("p1", "HU-01")
			if err == nil {
				t.Fatal("ObtenerItem() con un error de repositorio debería devolver un error")
			}
			if !errors.Is(err, errSimulado) {
				t.Errorf("err = %v, se esperaba un error que envuelva %v", err, errSimulado)
			}
			if errors.Is(err, domain.ErrProyectoNoEncontrado) {
				t.Errorf("err = %v, no debería mapearse a ErrProyectoNoEncontrado", err)
			}
			if errors.Is(err, domain.ErrItemNoEncontrado) {
				t.Errorf("err = %v, no debería mapearse a ErrItemNoEncontrado", err)
			}
			verificarItemBacklogVacio(t, obtenido)
		})
	}
}
