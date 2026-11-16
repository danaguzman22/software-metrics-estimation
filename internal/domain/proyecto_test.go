package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func fechaUS01(anio int, mes time.Month, dia int) time.Time {
	return time.Date(anio, mes, dia, 0, 0, 0, 0, time.UTC)
}

// verificarProyectoVacio comprueba que un error de validación devuelve el
// valor cero de domain.Proyecto.
func verificarProyectoVacio(t *testing.T, proyecto domain.Proyecto) {
	t.Helper()
	if proyecto != (domain.Proyecto{}) {
		t.Errorf("se esperaba domain.Proyecto{} ante un error, se obtuvo %+v", proyecto)
	}
}

// CA1 / R1: el proyecto se crea con los tres datos y los conserva.
func TestUS01_CreacionExitosa(t *testing.T) {
	casos := []struct {
		descripcion string
		nombre      string
		fechaInicio time.Time
		fechaFin    time.Time
	}{
		{
			descripcion: "datos válidos",
			nombre:      "Sistema de Métricas",
			fechaInicio: fechaUS01(2026, time.March, 1),
			fechaFin:    fechaUS01(2026, time.June, 30),
		},
		{
			descripcion: "otro proyecto",
			nombre:      "Otro proyecto",
			fechaInicio: fechaUS01(2026, time.January, 15),
			fechaFin:    fechaUS01(2026, time.December, 31),
		},
		{
			// D2: un nombre válido conserva los espacios externos tal como se ingresó.
			descripcion: "nombre con espacios externos se conserva",
			nombre:      "  Proyecto con espacios  ",
			fechaInicio: fechaUS01(2026, time.March, 1),
			fechaFin:    fechaUS01(2026, time.June, 30),
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			proyecto, err := domain.NuevoProyecto(caso.nombre, caso.fechaInicio, caso.fechaFin)
			if err != nil {
				t.Fatalf("NuevoProyecto() devolvió un error inesperado: %v", err)
			}
			if proyecto.Nombre != caso.nombre {
				t.Errorf("Nombre = %q, se esperaba %q", proyecto.Nombre, caso.nombre)
			}
			if !proyecto.FechaInicio.Equal(caso.fechaInicio) {
				t.Errorf("FechaInicio = %v, se esperaba %v", proyecto.FechaInicio, caso.fechaInicio)
			}
			if !proyecto.FechaFin.Equal(caso.fechaFin) {
				t.Errorf("FechaFin = %v, se esperaba %v", proyecto.FechaFin, caso.fechaFin)
			}
		})
	}
}

// CA2 / R2 / V1: un nombre vacío o compuesto solo por espacios devuelve ErrNombreVacio.
func TestUS01_NombreVacio(t *testing.T) {
	casos := []struct {
		descripcion string
		nombre      string
	}{
		{
			descripcion: "nombre vacío",
			nombre:      "",
		},
		{
			// D2: un nombre compuesto únicamente por espacios se considera vacío.
			descripcion: "nombre compuesto solo por espacios",
			nombre:      "   ",
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			proyecto, err := domain.NuevoProyecto(caso.nombre, fechaUS01(2026, time.March, 1), fechaUS01(2026, time.June, 30))
			if err == nil {
				t.Fatal("NuevoProyecto() con nombre vacío o solo espacios debería devolver un error de validación")
			}
			if !errors.Is(err, domain.ErrNombreVacio) {
				t.Errorf("err = %v, se esperaba %v", err, domain.ErrNombreVacio)
			}
			verificarProyectoVacio(t, proyecto)
		})
	}
}

// CA3 / R3 / V2: una fecha de fin anterior a la de inicio devuelve ErrFechaFinAnteriorAInicio.
func TestUS01_FechaFinAnteriorAInicio(t *testing.T) {
	casos := []struct {
		descripcion string
		nombre      string
		fechaInicio time.Time
		fechaFin    time.Time
	}{
		{
			descripcion: "fin anterior al inicio",
			nombre:      "Sistema de Métricas",
			fechaInicio: fechaUS01(2026, time.June, 30),
			fechaFin:    fechaUS01(2026, time.March, 1),
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			proyecto, err := domain.NuevoProyecto(caso.nombre, caso.fechaInicio, caso.fechaFin)
			if err == nil {
				t.Fatal("NuevoProyecto() con fecha de fin anterior al inicio debería devolver un error")
			}
			if !errors.Is(err, domain.ErrFechaFinAnteriorAInicio) {
				t.Errorf("err = %v, se esperaba %v", err, domain.ErrFechaFinAnteriorAInicio)
			}
			verificarProyectoVacio(t, proyecto)
		})
	}
}

// CA4 / R4: todo proyecto creado válidamente queda en estado Planificado.
func TestUS01_EstadoInicialPlanificado(t *testing.T) {
	casos := []struct {
		descripcion string
		nombre      string
		fechaInicio time.Time
		fechaFin    time.Time
	}{
		{
			descripcion: "creación con datos válidos",
			nombre:      "Sistema de Métricas",
			fechaInicio: fechaUS01(2026, time.March, 1),
			fechaFin:    fechaUS01(2026, time.June, 30),
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			proyecto, err := domain.NuevoProyecto(caso.nombre, caso.fechaInicio, caso.fechaFin)
			if err != nil {
				t.Fatalf("NuevoProyecto() devolvió un error inesperado: %v", err)
			}
			if proyecto.Estado != domain.EstadoPlanificado {
				t.Errorf("Estado = %q, se esperaba %q", proyecto.Estado, domain.EstadoPlanificado)
			}
		})
	}
}

// CA1 / CA3: fecha de fin igual a la de inicio es válida (D1 confirmada por la persona).
func TestUS01_FechasIgualesValidas(t *testing.T) {
	casos := []struct {
		descripcion string
		nombre      string
		fecha       time.Time
	}{
		{
			descripcion: "inicio y fin iguales",
			nombre:      "Sistema de Métricas",
			fecha:       fechaUS01(2026, time.March, 1),
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			proyecto, err := domain.NuevoProyecto(caso.nombre, caso.fecha, caso.fecha)
			if err != nil {
				t.Fatalf("NuevoProyecto() con fechas iguales devolvió un error inesperado: %v", err)
			}
			if !proyecto.FechaInicio.Equal(caso.fecha) {
				t.Errorf("FechaInicio = %v, se esperaba %v", proyecto.FechaInicio, caso.fecha)
			}
			if !proyecto.FechaFin.Equal(caso.fecha) {
				t.Errorf("FechaFin = %v, se esperaba %v", proyecto.FechaFin, caso.fecha)
			}
		})
	}
}
