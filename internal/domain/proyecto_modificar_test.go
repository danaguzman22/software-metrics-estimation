package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// proyectoUS02 arma un proyecto existente de referencia con las fechas de los
// escenarios de US-02, para reutilizarlo en los tests de modificación.
func proyectoUS02(nombre string, estado domain.Estado) domain.Proyecto {
	return domain.Proyecto{
		Nombre:      nombre,
		FechaInicio: fechaUS01(2026, time.March, 1),
		FechaFin:    fechaUS01(2026, time.June, 30),
		Estado:      estado,
	}
}

// verificarProyectoUS02 compara nombre, fechas y estado de un proyecto.
func verificarProyectoUS02(t *testing.T, obtenido, esperado domain.Proyecto) {
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

// CA1 / R1: la modificación con datos válidos devuelve una copia con el nombre
// y las fechas nuevas, y conserva el estado del proyecto existente.
func TestUS02_ModificacionExitosa(t *testing.T) {
	casos := []struct {
		descripcion string
		nombre      string
		fechaInicio time.Time
		fechaFin    time.Time
	}{
		{
			descripcion: "datos válidos",
			nombre:      "Sistema de Métricas v2",
			fechaInicio: fechaUS01(2026, time.March, 15),
			fechaFin:    fechaUS01(2026, time.July, 31),
		},
		{
			// C7: un nombre válido conserva los espacios externos tal como se ingresó.
			descripcion: "nombre con espacios externos se conserva",
			nombre:      "  Proyecto con espacios  ",
			fechaInicio: fechaUS01(2026, time.March, 15),
			fechaFin:    fechaUS01(2026, time.July, 31),
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			original := proyectoUS02("Sistema de Métricas", domain.Estado("En curso"))

			modificado, err := domain.ModificarProyecto(original, caso.nombre, caso.fechaInicio, caso.fechaFin)
			if err != nil {
				t.Fatalf("ModificarProyecto() devolvió un error inesperado: %v", err)
			}
			if modificado.Nombre != caso.nombre {
				t.Errorf("Nombre = %q, se esperaba %q", modificado.Nombre, caso.nombre)
			}
			if !modificado.FechaInicio.Equal(caso.fechaInicio) {
				t.Errorf("FechaInicio = %v, se esperaba %v", modificado.FechaInicio, caso.fechaInicio)
			}
			if !modificado.FechaFin.Equal(caso.fechaFin) {
				t.Errorf("FechaFin = %v, se esperaba %v", modificado.FechaFin, caso.fechaFin)
			}
			if modificado.Estado != original.Estado {
				t.Errorf("Estado = %q, se esperaba conservar %q", modificado.Estado, original.Estado)
			}
		})
	}
}

// CA5 / R5 / D9: la modificación conserva el estado actual del proyecto, en
// cualquier estado, sin resetearlo a Planificado.
func TestUS02_ConservaEstado(t *testing.T) {
	casos := []struct {
		descripcion string
		estado      domain.Estado
	}{
		{
			descripcion: "estado En curso",
			estado:      domain.Estado("En curso"),
		},
		{
			descripcion: "estado Finalizado",
			estado:      domain.Estado("Finalizado"),
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			original := proyectoUS02("Sistema de Métricas", caso.estado)

			modificado, err := domain.ModificarProyecto(
				original,
				"Sistema de Métricas v2",
				fechaUS01(2026, time.March, 15),
				fechaUS01(2026, time.July, 31),
			)
			if err != nil {
				t.Fatalf("ModificarProyecto() devolvió un error inesperado: %v", err)
			}
			if modificado.Estado != caso.estado {
				t.Errorf("Estado = %q, se esperaba conservar %q", modificado.Estado, caso.estado)
			}
			if modificado.Estado == domain.EstadoPlanificado {
				t.Errorf("Estado = %q, no debería resetearse a Planificado", modificado.Estado)
			}
		})
	}
}

// CA1 / R1 / D8: la modificación devuelve una copia y no muta el proyecto original.
func TestUS02_NoMutaOriginal(t *testing.T) {
	original := proyectoUS02("Sistema de Métricas", domain.Estado("En curso"))

	modificado, err := domain.ModificarProyecto(
		original,
		"Sistema de Métricas v2",
		fechaUS01(2026, time.March, 15),
		fechaUS01(2026, time.July, 31),
	)
	if err != nil {
		t.Fatalf("ModificarProyecto() devolvió un error inesperado: %v", err)
	}

	esperadoOriginal := proyectoUS02("Sistema de Métricas", domain.Estado("En curso"))
	verificarProyectoUS02(t, original, esperadoOriginal)

	esperadoModificado := domain.Proyecto{
		Nombre:      "Sistema de Métricas v2",
		FechaInicio: fechaUS01(2026, time.March, 15),
		FechaFin:    fechaUS01(2026, time.July, 31),
		Estado:      domain.Estado("En curso"),
	}
	verificarProyectoUS02(t, modificado, esperadoModificado)
}

// CA2 / R2 / V1: un nombre vacío o compuesto solo por espacios devuelve
// ErrNombreVacio y no muta el proyecto original.
func TestUS02_NombreVacio(t *testing.T) {
	casos := []struct {
		descripcion string
		nombre      string
	}{
		{
			descripcion: "nombre vacío",
			nombre:      "",
		},
		{
			// C1: un nombre compuesto únicamente por espacios se considera vacío.
			descripcion: "nombre compuesto solo por espacios",
			nombre:      "   ",
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			original := proyectoUS02("Sistema de Métricas", domain.Estado("En curso"))

			modificado, err := domain.ModificarProyecto(
				original,
				caso.nombre,
				fechaUS01(2026, time.March, 1),
				fechaUS01(2026, time.June, 30),
			)
			if err == nil {
				t.Fatal("ModificarProyecto() con nombre vacío o solo espacios debería devolver un error de validación")
			}
			if !errors.Is(err, domain.ErrNombreVacio) {
				t.Errorf("err = %v, se esperaba %v", err, domain.ErrNombreVacio)
			}
			verificarProyectoVacio(t, modificado)
			verificarProyectoUS02(t, original, proyectoUS02("Sistema de Métricas", domain.Estado("En curso")))
		})
	}
}

// CA3 / R3 / V2: una fecha de fin anterior a la de inicio devuelve
// ErrFechaFinAnteriorAInicio y no muta el proyecto original.
func TestUS02_FechaFinAnteriorAInicio(t *testing.T) {
	casos := []struct {
		descripcion string
		nombre      string
		fechaInicio time.Time
		fechaFin    time.Time
	}{
		{
			descripcion: "fin anterior al inicio",
			nombre:      "Sistema de Métricas v2",
			fechaInicio: fechaUS01(2026, time.June, 30),
			fechaFin:    fechaUS01(2026, time.March, 1),
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			original := proyectoUS02("Sistema de Métricas", domain.Estado("En curso"))

			modificado, err := domain.ModificarProyecto(original, caso.nombre, caso.fechaInicio, caso.fechaFin)
			if err == nil {
				t.Fatal("ModificarProyecto() con fecha de fin anterior al inicio debería devolver un error")
			}
			if !errors.Is(err, domain.ErrFechaFinAnteriorAInicio) {
				t.Errorf("err = %v, se esperaba %v", err, domain.ErrFechaFinAnteriorAInicio)
			}
			verificarProyectoVacio(t, modificado)
			verificarProyectoUS02(t, original, proyectoUS02("Sistema de Métricas", domain.Estado("En curso")))
		})
	}
}

// CA1 / CA3 / R3: fecha de fin igual a la de inicio es válida al modificar.
func TestUS02_FechasIgualesValidas(t *testing.T) {
	casos := []struct {
		descripcion string
		nombre      string
		fecha       time.Time
	}{
		{
			descripcion: "inicio y fin iguales",
			nombre:      "Sistema de Métricas v2",
			fecha:       fechaUS01(2026, time.March, 1),
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			original := proyectoUS02("Sistema de Métricas", domain.Estado("En curso"))

			modificado, err := domain.ModificarProyecto(original, caso.nombre, caso.fecha, caso.fecha)
			if err != nil {
				t.Fatalf("ModificarProyecto() con fechas iguales devolvió un error inesperado: %v", err)
			}
			if !modificado.FechaInicio.Equal(caso.fecha) {
				t.Errorf("FechaInicio = %v, se esperaba %v", modificado.FechaInicio, caso.fecha)
			}
			if !modificado.FechaFin.Equal(caso.fecha) {
				t.Errorf("FechaFin = %v, se esperaba %v", modificado.FechaFin, caso.fecha)
			}
			if modificado.Estado != original.Estado {
				t.Errorf("Estado = %q, se esperaba conservar %q", modificado.Estado, original.Estado)
			}
		})
	}
}

// CA1 / R7 / D10: modificar con los mismos valores actuales es válido y
// conserva el estado.
func TestUS02_ValoresIgualesValidos(t *testing.T) {
	casos := []struct {
		descripcion string
		estado      domain.Estado
	}{
		{
			descripcion: "mismos valores en estado En curso",
			estado:      domain.Estado("En curso"),
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			original := proyectoUS02("Sistema de Métricas", caso.estado)

			modificado, err := domain.ModificarProyecto(
				original,
				original.Nombre,
				original.FechaInicio,
				original.FechaFin,
			)
			if err != nil {
				t.Fatalf("ModificarProyecto() con los mismos valores devolvió un error inesperado: %v", err)
			}
			verificarProyectoUS02(t, modificado, original)
		})
	}
}

// CA2 / CA3 / R6 / D3 / D8: ante cualquier error de validación se devuelve
// Proyecto{} y el proyecto original queda intacto.
func TestUS02_ErrorDevuelveProyectoVacioYNoMutaOriginal(t *testing.T) {
	casos := []struct {
		descripcion   string
		nombre        string
		fechaInicio   time.Time
		fechaFin      time.Time
		errorEsperado error
	}{
		{
			descripcion:   "nombre vacío",
			nombre:        "",
			fechaInicio:   fechaUS01(2026, time.March, 1),
			fechaFin:      fechaUS01(2026, time.June, 30),
			errorEsperado: domain.ErrNombreVacio,
		},
		{
			descripcion:   "nombre compuesto solo por espacios",
			nombre:        "   ",
			fechaInicio:   fechaUS01(2026, time.March, 1),
			fechaFin:      fechaUS01(2026, time.June, 30),
			errorEsperado: domain.ErrNombreVacio,
		},
		{
			descripcion:   "fecha de fin anterior al inicio",
			nombre:        "Sistema de Métricas v2",
			fechaInicio:   fechaUS01(2026, time.June, 30),
			fechaFin:      fechaUS01(2026, time.March, 1),
			errorEsperado: domain.ErrFechaFinAnteriorAInicio,
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			original := proyectoUS02("Sistema de Métricas", domain.Estado("En curso"))

			modificado, err := domain.ModificarProyecto(original, caso.nombre, caso.fechaInicio, caso.fechaFin)
			if err == nil {
				t.Fatal("ModificarProyecto() con datos inválidos debería devolver un error")
			}
			if !errors.Is(err, caso.errorEsperado) {
				t.Errorf("err = %v, se esperaba %v", err, caso.errorEsperado)
			}
			verificarProyectoVacio(t, modificado)
			verificarProyectoUS02(t, original, proyectoUS02("Sistema de Métricas", domain.Estado("En curso")))
		})
	}
}

// CA2 / CA3 / R2 / R3 / D4: con nombre y fechas inválidos a la vez prevalece
// el error de nombre.
func TestUS02_OrdenValidacionesNombreAntesQueFechas(t *testing.T) {
	original := proyectoUS02("Sistema de Métricas", domain.Estado("En curso"))

	modificado, err := domain.ModificarProyecto(
		original,
		"   ",
		fechaUS01(2026, time.June, 30),
		fechaUS01(2026, time.March, 1),
	)
	if err == nil {
		t.Fatal("ModificarProyecto() con nombre y fechas inválidos debería devolver un error")
	}
	if !errors.Is(err, domain.ErrNombreVacio) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrNombreVacio)
	}
	if errors.Is(err, domain.ErrFechaFinAnteriorAInicio) {
		t.Errorf("err = %v, no debería devolverse el error de fechas cuando el nombre también es inválido", err)
	}
	verificarProyectoVacio(t, modificado)
	verificarProyectoUS02(t, original, proyectoUS02("Sistema de Métricas", domain.Estado("En curso")))
}
