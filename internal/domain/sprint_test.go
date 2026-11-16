package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// fechaUS07 arma una fecha sin hora para los escenarios de US-07.
func fechaUS07(anio int, mes time.Month, dia int) time.Time {
	return time.Date(anio, mes, dia, 0, 0, 0, 0, time.UTC)
}

// verificarSprintUS07Vacio comprueba que un error de validación devuelve el
// valor cero de domain.Sprint (D7, sección 8 de la spec).
func verificarSprintUS07Vacio(t *testing.T, sprint domain.Sprint) {
	t.Helper()
	if sprint.Numero != 0 || sprint.Objetivo != "" ||
		!sprint.FechaInicio.IsZero() || !sprint.FechaFin.IsZero() ||
		sprint.Estado != "" {
		t.Errorf("se esperaba domain.Sprint{} ante un error, se obtuvo %+v", sprint)
	}
}

// verificarSprintUS07 compara campo por campo un Sprint con su esperado.
func verificarSprintUS07(t *testing.T, obtenido, esperado domain.Sprint) {
	t.Helper()
	if obtenido.Numero != esperado.Numero {
		t.Errorf("Numero = %d, se esperaba %d", obtenido.Numero, esperado.Numero)
	}
	if obtenido.Objetivo != esperado.Objetivo {
		t.Errorf("Objetivo = %q, se esperaba %q", obtenido.Objetivo, esperado.Objetivo)
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

// CA1 / R1: un Sprint válido se crea con número, objetivo y fechas.
func TestUS07_NuevoSprint_Exitoso(t *testing.T) {
	inicio := fechaUS07(2026, time.March, 2)
	fin := fechaUS07(2026, time.March, 20)

	sprint, err := domain.NuevoSprint(1, "Planificar la primera iteración", inicio, fin)
	if err != nil {
		t.Fatalf("NuevoSprint() devolvió un error inesperado: %v", err)
	}
	verificarSprintUS07(t, sprint, domain.Sprint{
		Numero:      1,
		Objetivo:    "Planificar la primera iteración",
		FechaInicio: inicio,
		FechaFin:    fin,
		Estado:      domain.EstadoSprintActivo,
	})
}

// CA1 / R1 / D2: todo Sprint creado nace en estado "Activo".
func TestUS07_NuevoSprint_EstadoActivo(t *testing.T) {
	sprint, err := domain.NuevoSprint(
		1,
		"Planificar la primera iteración",
		fechaUS07(2026, time.March, 2),
		fechaUS07(2026, time.March, 20),
	)
	if err != nil {
		t.Fatalf("NuevoSprint() devolvió un error inesperado: %v", err)
	}
	if sprint.Estado != "Activo" {
		t.Errorf("Estado = %q, se esperaba %q", sprint.Estado, "Activo")
	}
	if sprint.Estado != domain.EstadoSprintActivo {
		t.Errorf("Estado = %q, se esperaba domain.EstadoSprintActivo (%q)", sprint.Estado, domain.EstadoSprintActivo)
	}
}

// CA1 / R2 / D3: el objetivo se guarda sin espacios al inicio ni al final.
func TestUS07_NuevoSprint_RecortaObjetivo(t *testing.T) {
	sprint, err := domain.NuevoSprint(
		1,
		"  Planificar la primera iteración  ",
		fechaUS07(2026, time.March, 2),
		fechaUS07(2026, time.March, 20),
	)
	if err != nil {
		t.Fatalf("NuevoSprint() devolvió un error inesperado: %v", err)
	}
	if sprint.Objetivo != "Planificar la primera iteración" {
		t.Errorf("Objetivo = %q, se esperaba %q", sprint.Objetivo, "Planificar la primera iteración")
	}
}

// CA1 / R3 / D4: un número cero devuelve ErrNumeroSprintInvalido.
func TestUS07_NuevoSprint_NumeroCero(t *testing.T) {
	sprint, err := domain.NuevoSprint(
		0,
		"Planificar la primera iteración",
		fechaUS07(2026, time.March, 2),
		fechaUS07(2026, time.March, 20),
	)
	if !errors.Is(err, domain.ErrNumeroSprintInvalido) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrNumeroSprintInvalido)
	}
	verificarSprintUS07Vacio(t, sprint)
}

// CA1 / R3 / D4: un número negativo devuelve ErrNumeroSprintInvalido.
func TestUS07_NuevoSprint_NumeroNegativo(t *testing.T) {
	sprint, err := domain.NuevoSprint(
		-1,
		"Planificar la primera iteración",
		fechaUS07(2026, time.March, 2),
		fechaUS07(2026, time.March, 20),
	)
	if !errors.Is(err, domain.ErrNumeroSprintInvalido) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrNumeroSprintInvalido)
	}
	verificarSprintUS07Vacio(t, sprint)
}

// CA1 / R5 / D5: una fecha de fin anterior a la de inicio es inválida.
func TestUS07_NuevoSprint_FinAnteriorAInicio(t *testing.T) {
	sprint, err := domain.NuevoSprint(
		1,
		"Planificar la primera iteración",
		fechaUS07(2026, time.March, 20),
		fechaUS07(2026, time.March, 2),
	)
	if !errors.Is(err, domain.ErrFechasSprintInvalidas) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrFechasSprintInvalidas)
	}
	verificarSprintUS07Vacio(t, sprint)
}

// CA1 / R5 / D5: fecha de inicio igual a la de fin es válida.
func TestUS07_NuevoSprint_FechasIgualesValidas(t *testing.T) {
	dia := fechaUS07(2026, time.March, 10)
	sprint, err := domain.NuevoSprint(2, "Realizar una iteración de un día", dia, dia)
	if err != nil {
		t.Fatalf("NuevoSprint() devolvió un error inesperado: %v", err)
	}
	if !sprint.FechaInicio.Equal(dia) || !sprint.FechaFin.Equal(dia) {
		t.Errorf("fechas = %v-%v, se esperaban ambas %v", sprint.FechaInicio, sprint.FechaFin, dia)
	}
}

// CA1 / R5 / D5: dos fechas del mismo día con distinta hora son válidas porque
// la comparación es por año, mes y día.
func TestUS07_NuevoSprint_MismoDiaDistintaHoraValido(t *testing.T) {
	inicio := time.Date(2026, time.March, 10, 18, 0, 0, 0, time.UTC)
	fin := time.Date(2026, time.March, 10, 9, 0, 0, 0, time.UTC)

	sprint, err := domain.NuevoSprint(2, "Realizar una iteración de un día", inicio, fin)
	if err != nil {
		t.Fatalf("NuevoSprint() devolvió un error inesperado: %v", err)
	}
	if !sprint.FechaInicio.Equal(inicio) {
		t.Errorf("FechaInicio = %v, se esperaba %v", sprint.FechaInicio, inicio)
	}
	if !sprint.FechaFin.Equal(fin) {
		t.Errorf("FechaFin = %v, se esperaba %v", sprint.FechaFin, fin)
	}
}

// CA2 / R2 / D3: un objetivo vacío devuelve ErrObjetivoSprintVacio.
func TestUS07_NuevoSprint_ObjetivoVacio(t *testing.T) {
	sprint, err := domain.NuevoSprint(
		1,
		"",
		fechaUS07(2026, time.March, 2),
		fechaUS07(2026, time.March, 20),
	)
	if !errors.Is(err, domain.ErrObjetivoSprintVacio) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrObjetivoSprintVacio)
	}
	verificarSprintUS07Vacio(t, sprint)
}

// CA2 / R2 / D3: un objetivo compuesto solo por espacios devuelve
// ErrObjetivoSprintVacio.
func TestUS07_NuevoSprint_ObjetivoSoloEspacios(t *testing.T) {
	sprint, err := domain.NuevoSprint(
		1,
		"   ",
		fechaUS07(2026, time.March, 2),
		fechaUS07(2026, time.March, 20),
	)
	if !errors.Is(err, domain.ErrObjetivoSprintVacio) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrObjetivoSprintVacio)
	}
	verificarSprintUS07Vacio(t, sprint)
}

// CA1 / R3 / R10 / D7: el número inválido se valida antes que el objetivo
// vacío.
func TestUS07_NuevoSprint_OrdenNumeroAntesQueObjetivo(t *testing.T) {
	sprint, err := domain.NuevoSprint(
		0,
		"",
		fechaUS07(2026, time.March, 2),
		fechaUS07(2026, time.March, 20),
	)
	if !errors.Is(err, domain.ErrNumeroSprintInvalido) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrNumeroSprintInvalido)
	}
	if errors.Is(err, domain.ErrObjetivoSprintVacio) {
		t.Errorf("err = %v, no debería devolverse ErrObjetivoSprintVacio si el número es inválido", err)
	}
	verificarSprintUS07Vacio(t, sprint)
}

// CA2 / R2 / R10 / D7: el objetivo vacío se valida antes que las fechas
// inválidas.
func TestUS07_NuevoSprint_OrdenObjetivoAntesQueFechas(t *testing.T) {
	sprint, err := domain.NuevoSprint(
		1,
		"   ",
		fechaUS07(2026, time.March, 20),
		fechaUS07(2026, time.March, 2),
	)
	if !errors.Is(err, domain.ErrObjetivoSprintVacio) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrObjetivoSprintVacio)
	}
	if errors.Is(err, domain.ErrFechasSprintInvalidas) {
		t.Errorf("err = %v, no debería devolverse ErrFechasSprintInvalidas si el objetivo está vacío", err)
	}
	verificarSprintUS07Vacio(t, sprint)
}

// CA1 / R1: NuevoSprint no modifica las entradas recibidas.
func TestUS07_NuevoSprint_NoModificaEntradas(t *testing.T) {
	numero := 3
	objetivo := "  Planificar la iteración de julio  "
	inicio := fechaUS07(2026, time.July, 1)
	fin := fechaUS07(2026, time.July, 15)

	numeroOriginal, objetivoOriginal := numero, objetivo
	inicioOriginal, finOriginal := inicio, fin

	sprint, err := domain.NuevoSprint(numero, objetivo, inicio, fin)
	if err != nil {
		t.Fatalf("NuevoSprint() devolvió un error inesperado: %v", err)
	}

	if numero != numeroOriginal {
		t.Errorf("numero = %d, se esperaba que quedara en %d", numero, numeroOriginal)
	}
	if objetivo != objetivoOriginal {
		t.Errorf("objetivo = %q, se esperaba que quedara en %q", objetivo, objetivoOriginal)
	}
	if !inicio.Equal(inicioOriginal) {
		t.Errorf("inicio = %v, se esperaba que quedara en %v", inicio, inicioOriginal)
	}
	if !fin.Equal(finOriginal) {
		t.Errorf("fin = %v, se esperaba que quedara en %v", fin, finOriginal)
	}
	verificarSprintUS07(t, sprint, domain.Sprint{
		Numero:      numeroOriginal,
		Objetivo:    "Planificar la iteración de julio",
		FechaInicio: inicioOriginal,
		FechaFin:    finOriginal,
		Estado:      domain.EstadoSprintActivo,
	})
}
