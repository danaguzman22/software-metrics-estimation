package domain_test

import (
	"errors"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func TestUS15_NuevoDefectoValido(t *testing.T) {
	sprintResolucion := 2
	defecto, err := domain.NuevoDefecto(
		"  No guarda los cambios  ",
		domain.SeveridadDefectoAlta,
		domain.EstadoDefectoAbierto,
		"HU-01",
		1,
		&sprintResolucion,
	)
	if err != nil {
		t.Fatalf("NuevoDefecto() devolvió un error: %v", err)
	}
	if defecto.Descripcion != "No guarda los cambios" ||
		defecto.Severidad != domain.SeveridadDefectoAlta ||
		defecto.Estado != domain.EstadoDefectoAbierto ||
		defecto.ItemID != "HU-01" ||
		defecto.SprintDeteccion != 1 ||
		defecto.SprintResolucion == nil ||
		*defecto.SprintResolucion != 2 {
		t.Fatalf("defecto = %+v, sus campos no coinciden con la entrada normalizada", defecto)
	}
}

func TestUS15_NuevoDefectoValidaCamposYResolucion(t *testing.T) {
	sprint2 := 2
	sprint1 := 1
	casos := []struct {
		nombre      string
		descripcion string
		severidad   string
		estado      string
		sprintDet   int
		sprintRes   *int
		esperado    error
	}{
		{"descripción vacía", "   ", domain.SeveridadDefectoAlta, domain.EstadoDefectoAbierto, 1, nil, domain.ErrDescripcionDefectoVacia},
		{"severidad inválida", "Error", "Crítica", domain.EstadoDefectoAbierto, 1, nil, domain.ErrSeveridadDefectoInvalida},
		{"estado inválido", "Error", domain.SeveridadDefectoAlta, "Cerrado", 1, nil, domain.ErrEstadoDefectoInvalido},
		{"Sprint detección inválido", "Error", domain.SeveridadDefectoAlta, domain.EstadoDefectoAbierto, 0, nil, domain.ErrSprintDeteccionDefectoInvalido},
		{"resuelto sin Sprint", "Error", domain.SeveridadDefectoAlta, domain.EstadoDefectoResuelto, 1, nil, domain.ErrSprintResolucionDefectoObligatorio},
		{"Sprint resolución anterior", "Error", domain.SeveridadDefectoAlta, domain.EstadoDefectoAbierto, 2, &sprint1, domain.ErrSprintResolucionAnteriorADeteccion},
		{"Sprint resolución inválido", "Error", domain.SeveridadDefectoAlta, domain.EstadoDefectoAbierto, 1, intPunteroDefecto(0), domain.ErrSprintResolucionDefectoInvalido},
		{"resuelto con Sprint posterior", "Error", domain.SeveridadDefectoAlta, domain.EstadoDefectoResuelto, 1, &sprint2, nil},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			_, err := domain.NuevoDefecto(caso.descripcion, caso.severidad, caso.estado, "HU-01", caso.sprintDet, caso.sprintRes)
			if caso.esperado == nil {
				if err != nil {
					t.Fatalf("NuevoDefecto() devolvió un error inesperado: %v", err)
				}
				return
			}
			if !errors.Is(err, caso.esperado) {
				t.Fatalf("err = %v, se esperaba %v", err, caso.esperado)
			}
		})
	}
}

func TestUS15_CambiarEstadoDefectoValidaEstadoYResolucion(t *testing.T) {
	defecto, err := domain.NuevoDefecto("Error", domain.SeveridadDefectoMedia, domain.EstadoDefectoAbierto, "HU-01", 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	original := defecto

	if _, err := defecto.CambiarEstado("Resuelto", nil); !errors.Is(err, domain.ErrSprintResolucionDefectoObligatorio) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrSprintResolucionDefectoObligatorio)
	}
	if defecto.Estado != original.Estado || defecto.SprintResolucion != nil {
		t.Fatalf("el defecto original cambió tras una transición inválida: %+v", defecto)
	}

	sprintResolucion := 1
	actualizado, err := defecto.CambiarEstado(domain.EstadoDefectoResuelto, &sprintResolucion)
	if err != nil {
		t.Fatalf("CambiarEstado() devolvió un error inesperado: %v", err)
	}
	if actualizado.Estado != domain.EstadoDefectoResuelto ||
		actualizado.SprintResolucion == nil ||
		*actualizado.SprintResolucion != 1 {
		t.Fatalf("defecto actualizado = %+v", actualizado)
	}
}

func TestUS15_ReabrirDefectoLimpiaSprintResolucion(t *testing.T) {
	sprintResolucion := 2
	defecto, err := domain.NuevoDefecto(
		"Error corregido pero vuelve a ocurrir",
		domain.SeveridadDefectoAlta,
		domain.EstadoDefectoResuelto,
		"HU-01",
		1,
		&sprintResolucion,
	)
	if err != nil {
		t.Fatal(err)
	}

	reabierto, err := defecto.CambiarEstado(domain.EstadoDefectoEnCurso, nil)
	if err != nil {
		t.Fatalf("CambiarEstado() al reabrir devolvió error: %v", err)
	}
	if reabierto.Estado != domain.EstadoDefectoEnCurso || reabierto.SprintResolucion != nil {
		t.Fatalf("defecto reabierto = %+v; se esperaba estado En curso y sin Sprint de resolución", reabierto)
	}
	if defecto.Estado != domain.EstadoDefectoResuelto ||
		defecto.SprintResolucion == nil ||
		*defecto.SprintResolucion != 2 {
		t.Fatalf("la reapertura modificó el defecto original: %+v", defecto)
	}
}

func intPunteroDefecto(valor int) *int {
	return &valor
}
