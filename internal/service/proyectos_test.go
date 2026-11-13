package service

import (
	"errors"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func fechaUS02(anio int, mes time.Month, dia int) time.Time {
	return time.Date(anio, mes, dia, 0, 0, 0, 0, time.UTC)
}

func proyectoUS02Servicio(nombre string, estado domain.Estado) domain.Proyecto {
	return domain.Proyecto{
		Nombre:      nombre,
		FechaInicio: fechaUS02(2026, time.March, 1),
		FechaFin:    fechaUS02(2026, time.June, 30),
		Estado:      estado,
	}
}

func nuevoServicioProyectosConDatos(registros ...ProyectoRegistrado) (*ServicioProyectos, *fakeRepositorioProyectos, *fakeRepositorioIntegrantes) {
	proyectos := nuevoFakeProyectos()
	for _, r := range registros {
		proyectos.proyectos[r.ID] = r.Proyecto
	}
	integrantes := nuevoFakeIntegrantes()
	return NuevoServicioProyectos(proyectos, integrantes), proyectos, integrantes
}

func buscarProyectoEnFake(t *testing.T, proyectos *fakeRepositorioProyectos, id string) domain.Proyecto {
	t.Helper()
	p, ok := proyectos.proyectos[id]
	if !ok {
		t.Fatalf("no se encontró el registro con ID %q", id)
	}
	return p
}

func verificarProyecto(t *testing.T, obtenido, esperado domain.Proyecto) {
	t.Helper()
	if obtenido.Nombre != esperado.Nombre {
		t.Errorf("Nombre = %q, se esperaba %q", obtenido.Nombre, esperado.Nombre)
	}
	if !obtenido.FechaInicio.Equal(esperado.FechaInicio) {
		t.Errorf("FechaInicio = %v, se esperaba %v", obtenido.FechaInicio, esperado.FechaInicio)
	}
	if !obtenido.FechaFin.Equal(esperado.FechaFin) {
		t.Errorf("FechaFin = %v, se esperaba %v", obtenido.FechaFin, esperado.FechaFin)
	}
	if obtenido.Estado != esperado.Estado {
		t.Errorf("Estado = %q, se esperaba %q", obtenido.Estado, esperado.Estado)
	}
}

func verificarProyectoVacio(t *testing.T, proyecto domain.Proyecto) {
	t.Helper()
	if proyecto != (domain.Proyecto{}) {
		t.Errorf("se esperaba domain.Proyecto{} ante un error, se obtuvo %+v", proyecto)
	}
}

func TestUS02_Modificar_PersisteCambios(t *testing.T) {
	casos := []struct {
		descripcion string
		nombre      string
	}{
		{descripcion: "datos válidos", nombre: "Sistema de Métricas v2"},
		{descripcion: "nombre con espacios externos", nombre: "  Proyecto con espacios  "},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			original := proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))
			servicio, proyectos, _ := nuevoServicioProyectosConDatos(ProyectoRegistrado{ID: "p1", Proyecto: original})

			fechaInicio := fechaUS02(2026, time.March, 15)
			fechaFin := fechaUS02(2026, time.July, 31)
			modificado, err := servicio.ModificarProyecto("p1", caso.nombre, fechaInicio, fechaFin)
			if err != nil {
				t.Fatalf("ModificarProyecto() devolvió un error inesperado: %v", err)
			}

			esperado := domain.Proyecto{Nombre: caso.nombre, FechaInicio: fechaInicio, FechaFin: fechaFin, Estado: domain.Estado("En curso")}
			verificarProyecto(t, modificado, esperado)
			verificarProyecto(t, buscarProyectoEnFake(t, proyectos, "p1"), esperado)
		})
	}
}

func TestUS02_Modificar_ConservaEstado(t *testing.T) {
	casos := []struct {
		descripcion string
		estado      domain.Estado
	}{
		{descripcion: "En curso", estado: domain.Estado("En curso")},
		{descripcion: "Finalizado", estado: domain.Estado("Finalizado")},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			original := proyectoUS02Servicio("Sistema de Métricas", caso.estado)
			servicio, proyectos, _ := nuevoServicioProyectosConDatos(ProyectoRegistrado{ID: "p1", Proyecto: original})

			modificado, err := servicio.ModificarProyecto("p1", "Sistema de Métricas v2", fechaUS02(2026, time.March, 15), fechaUS02(2026, time.July, 31))
			if err != nil {
				t.Fatalf("ModificarProyecto() devolvió un error inesperado: %v", err)
			}
			if modificado.Estado != caso.estado {
				t.Errorf("Estado = %q, se esperaba conservar %q", modificado.Estado, caso.estado)
			}
			if modificado.Estado == domain.EstadoPlanificado {
				t.Errorf("Estado = %q, no debería resetearse a Planificado", modificado.Estado)
			}
			if buscarProyectoEnFake(t, proyectos, "p1").Estado != caso.estado {
				t.Errorf("Estado persistido = %q, se esperaba conservar %q", buscarProyectoEnFake(t, proyectos, "p1").Estado, caso.estado)
			}
		})
	}
}

func TestUS02_Modificar_NoModificaOtrosProyectos(t *testing.T) {
	p1 := proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))
	p2 := domain.Proyecto{Nombre: "Otro proyecto", FechaInicio: fechaUS02(2026, time.January, 15), FechaFin: fechaUS02(2026, time.December, 31), Estado: domain.EstadoPlanificado}
	p3 := domain.Proyecto{Nombre: "Tercer proyecto", FechaInicio: fechaUS02(2026, time.February, 10), FechaFin: fechaUS02(2026, time.May, 20), Estado: domain.Estado("Finalizado")}
	servicio, proyectos, _ := nuevoServicioProyectosConDatos(
		ProyectoRegistrado{ID: "p1", Proyecto: p1},
		ProyectoRegistrado{ID: "p2", Proyecto: p2},
		ProyectoRegistrado{ID: "p3", Proyecto: p3},
	)

	_, err := servicio.ModificarProyecto("p1", "Sistema de Métricas v2", fechaUS02(2026, time.March, 15), fechaUS02(2026, time.July, 31))
	if err != nil {
		t.Fatalf("ModificarProyecto() devolvió un error inesperado: %v", err)
	}

	if len(proyectos.proyectos) != 3 {
		t.Fatalf("la cantidad de registros cambió: se esperaban 3, hay %d", len(proyectos.proyectos))
	}
	esperadoP1 := domain.Proyecto{Nombre: "Sistema de Métricas v2", FechaInicio: fechaUS02(2026, time.March, 15), FechaFin: fechaUS02(2026, time.July, 31), Estado: domain.Estado("En curso")}
	verificarProyecto(t, buscarProyectoEnFake(t, proyectos, "p1"), esperadoP1)
	verificarProyecto(t, buscarProyectoEnFake(t, proyectos, "p2"), p2)
	verificarProyecto(t, buscarProyectoEnFake(t, proyectos, "p3"), p3)
}

func TestUS02_Modificar_ValoresIguales(t *testing.T) {
	original := proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))
	servicio, proyectos, _ := nuevoServicioProyectosConDatos(ProyectoRegistrado{ID: "p1", Proyecto: original})

	modificado, err := servicio.ModificarProyecto("p1", original.Nombre, original.FechaInicio, original.FechaFin)
	if err != nil {
		t.Fatalf("ModificarProyecto() con los mismos valores devolvió un error inesperado: %v", err)
	}
	verificarProyecto(t, modificado, original)
	verificarProyecto(t, buscarProyectoEnFake(t, proyectos, "p1"), original)
}

func TestUS02_Modificar_ValidacionFallidaNoPersiste(t *testing.T) {
	casos := []struct {
		descripcion   string
		nombre        string
		fechaInicio   time.Time
		fechaFin      time.Time
		errorEsperado error
	}{
		{descripcion: "nombre vacío", nombre: "", fechaInicio: fechaUS02(2026, time.March, 1), fechaFin: fechaUS02(2026, time.June, 30), errorEsperado: domain.ErrNombreVacio},
		{descripcion: "nombre compuesto solo por espacios", nombre: "   ", fechaInicio: fechaUS02(2026, time.March, 1), fechaFin: fechaUS02(2026, time.June, 30), errorEsperado: domain.ErrNombreVacio},
		{descripcion: "fecha de fin anterior al inicio", nombre: "Sistema de Métricas v2", fechaInicio: fechaUS02(2026, time.June, 30), fechaFin: fechaUS02(2026, time.March, 1), errorEsperado: domain.ErrFechaFinAnteriorAInicio},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			original := proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))
			servicio, proyectos, _ := nuevoServicioProyectosConDatos(ProyectoRegistrado{ID: "p1", Proyecto: original})

			modificado, err := servicio.ModificarProyecto("p1", caso.nombre, caso.fechaInicio, caso.fechaFin)
			if err == nil {
				t.Fatal("ModificarProyecto() con datos inválidos debería devolver un error")
			}
			if !errors.Is(err, caso.errorEsperado) {
				t.Errorf("err = %v, se esperaba %v", err, caso.errorEsperado)
			}
			verificarProyectoVacio(t, modificado)
			if proyectos.actualizaciones != 0 {
				t.Errorf("Actualizar fue invocado %d veces ante un error de validación; no debía escribirse nada", proyectos.actualizaciones)
			}
			verificarProyecto(t, buscarProyectoEnFake(t, proyectos, "p1"), original)
		})
	}
}

func TestUS02_Modificar_ProyectoNoEncontrado(t *testing.T) {
	servicio, _, _ := nuevoServicioProyectosConDatos(
		ProyectoRegistrado{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
		ProyectoRegistrado{ID: "p2", Proyecto: proyectoUS02Servicio("Otro proyecto", domain.EstadoPlanificado)},
	)

	modificado, err := servicio.ModificarProyecto("ausente", "Proyecto actualizado", fechaUS02(2026, time.March, 1), fechaUS02(2026, time.June, 30))
	if err == nil {
		t.Fatal("ModificarProyecto() con un ID inexistente debería devolver un error")
	}
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	verificarProyectoVacio(t, modificado)
}

func TestUS02_Modificar_ProyectoNoEncontradoPrecedeValidaciones(t *testing.T) {
	servicio, _, _ := nuevoServicioProyectosConDatos(
		ProyectoRegistrado{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
	)

	modificado, err := servicio.ModificarProyecto("ausente", "", fechaUS02(2026, time.June, 30), fechaUS02(2026, time.March, 1))
	if err == nil {
		t.Fatal("ModificarProyecto() con un ID inexistente debería devolver un error")
	}
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if errors.Is(err, domain.ErrNombreVacio) {
		t.Errorf("err = %v, no debería devolverse ErrNombreVacio si el proyecto no existe", err)
	}
	if errors.Is(err, domain.ErrFechaFinAnteriorAInicio) {
		t.Errorf("err = %v, no debería devolverse ErrFechaFinAnteriorAInicio si el proyecto no existe", err)
	}
	verificarProyectoVacio(t, modificado)
}

func TestUS02_Modificar_ArchivoInexistente(t *testing.T) {
	servicio, proyectos, _ := nuevoServicioProyectosConDatos()

	modificado, err := servicio.ModificarProyecto("p1", "Proyecto actualizado", fechaUS02(2026, time.March, 1), fechaUS02(2026, time.June, 30))
	if err == nil {
		t.Fatal("ModificarProyecto() sin archivo de colección debería devolver un error")
	}
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	verificarProyectoVacio(t, modificado)
	if proyectos.actualizaciones != 0 {
		t.Errorf("en ningún camino de error se debe escribir; actualizaciones = %d", proyectos.actualizaciones)
	}
}

var errInfraCorruptoUS02 = errors.New("contenido inválido")

func TestUS02_Modificar_JSONCorruptoNoEsNoEncontrado(t *testing.T) {
	proyectos := nuevoFakeProyectos()
	proyectos.obtenerErr = errInfraCorruptoUS02
	servicio := NuevoServicioProyectos(proyectos, nuevoFakeIntegrantes())

	modificado, err := servicio.ModificarProyecto("p1", "Proyecto actualizado", fechaUS02(2026, time.March, 1), fechaUS02(2026, time.June, 30))
	if err == nil {
		t.Fatal("ModificarProyecto() con JSON corrupto debería devolver un error")
	}
	if !errors.Is(err, errInfraCorruptoUS02) {
		t.Errorf("err = %v, se esperaba un error que envuelva %v", err, errInfraCorruptoUS02)
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, no debería mapearse a ErrProyectoNoEncontrado", err)
	}
	verificarProyectoVacio(t, modificado)
}

var errFalloEscrituraUS02 = errors.New("fallo simulado de escritura")

func TestUS02_Modificar_FalloDeEscrituraNoInformaExito(t *testing.T) {
	original := proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))
	proyectos := nuevoFakeProyectos()
	proyectos.proyectos["p1"] = original
	proyectos.actualizarErr = errFalloEscrituraUS02
	servicio := NuevoServicioProyectos(proyectos, nuevoFakeIntegrantes())

	modificado, err := servicio.ModificarProyecto("p1", "Sistema de Métricas v2", fechaUS02(2026, time.March, 15), fechaUS02(2026, time.July, 31))
	if err == nil {
		t.Fatal("ModificarProyecto() con fallo de escritura debería devolver un error")
	}
	if !errors.Is(err, errFalloEscrituraUS02) {
		t.Errorf("err = %v, se esperaba un error que envuelva %v", err, errFalloEscrituraUS02)
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, no debería ser ErrProyectoNoEncontrado", err)
	}
	if errors.Is(err, domain.ErrNombreVacio) || errors.Is(err, domain.ErrFechaFinAnteriorAInicio) {
		t.Errorf("err = %v, no debería ser un error de validación de dominio", err)
	}
	verificarProyectoVacio(t, modificado)
}
