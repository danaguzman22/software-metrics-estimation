package domain_test

import (
	"errors"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// R5 / D3 (US-12): EsValorVotoPokerValido informa si el valor pertenece a la
// serie Fibonacci extendida 1, 2, 3, 5, 8, 13, 21.
func TestUS12_Dominio_EsValorVotoPokerValido(t *testing.T) {
	casos := []struct {
		valor    int
		esperado bool
	}{
		{valor: 1, esperado: true},
		{valor: 2, esperado: true},
		{valor: 3, esperado: true},
		{valor: 5, esperado: true},
		{valor: 8, esperado: true},
		{valor: 13, esperado: true},
		{valor: 21, esperado: true},
		{valor: 0, esperado: false},
		{valor: 4, esperado: false},
		{valor: 6, esperado: false},
		{valor: 22, esperado: false},
		{valor: -1, esperado: false},
		{valor: -13, esperado: false},
	}

	for _, caso := range casos {
		if obtenido := domain.EsValorVotoPokerValido(caso.valor); obtenido != caso.esperado {
			t.Errorf("EsValorVotoPokerValido(%d) = %t, se esperaba %t", caso.valor, obtenido, caso.esperado)
		}
	}
}

// CA1 / R5 / D3 (US-12): un voto válido se construye con integrante y valor.
func TestUS12_Dominio_NuevoVotoPoker(t *testing.T) {
	t.Run("caso válido", func(t *testing.T) {
		voto, err := domain.NuevoVotoPoker("Ana", 5)
		if err != nil {
			t.Fatalf("NuevoVotoPoker() devolvió un error inesperado: %v", err)
		}
		if voto.Integrante != "Ana" {
			t.Errorf("Integrante = %q, se esperaba %q", voto.Integrante, "Ana")
		}
		if voto.Valor != 5 {
			t.Errorf("Valor = %d, se esperaba %d", voto.Valor, 5)
		}
	})

	t.Run("integrante vacío", func(t *testing.T) {
		voto, err := domain.NuevoVotoPoker("", 5)
		if !errors.Is(err, domain.ErrNombreIntegranteVacio) {
			t.Fatalf("err = %v, se esperaba %v", err, domain.ErrNombreIntegranteVacio)
		}
		if voto != (domain.VotoPoker{}) {
			t.Errorf("voto = %+v, se esperaba domain.VotoPoker{}", voto)
		}
	})

	t.Run("integrante solo espacios", func(t *testing.T) {
		voto, err := domain.NuevoVotoPoker("   ", 5)
		if !errors.Is(err, domain.ErrNombreIntegranteVacio) {
			t.Fatalf("err = %v, se esperaba %v", err, domain.ErrNombreIntegranteVacio)
		}
		if voto != (domain.VotoPoker{}) {
			t.Errorf("voto = %+v, se esperaba domain.VotoPoker{}", voto)
		}
	})

	t.Run("valor inválido", func(t *testing.T) {
		valoresInvalidos := []int{0, 4, 6, 22, -1}
		for _, valor := range valoresInvalidos {
			voto, err := domain.NuevoVotoPoker("Ana", valor)
			if !errors.Is(err, domain.ErrValorVotoPokerInvalido) {
				t.Errorf("NuevoVotoPoker(\"Ana\", %d): err = %v, se esperaba %v", valor, err, domain.ErrValorVotoPokerInvalido)
			}
			if voto != (domain.VotoPoker{}) {
				t.Errorf("NuevoVotoPoker(\"Ana\", %d): voto = %+v, se esperaba domain.VotoPoker{}", valor, voto)
			}
		}
	})
}

// R9 / D7 (US-12): el número de la siguiente ronda es 1 si no hay rondas,
// max(numero)+1 si todas están Revelada, o ErrRondaAbiertaExistente si hay
// una ronda Abierta.
func TestUS12_Dominio_SiguienteNumeroRonda(t *testing.T) {
	t.Run("sin rondas devuelve 1", func(t *testing.T) {
		casos := []struct {
			descripcion string
			rondas      []domain.RondaPoker
		}{
			{descripcion: "lista nil", rondas: nil},
			{descripcion: "lista vacía", rondas: []domain.RondaPoker{}},
		}

		for _, caso := range casos {
			t.Run(caso.descripcion, func(t *testing.T) {
				numero, err := domain.SiguienteNumeroRonda(caso.rondas)
				if err != nil {
					t.Fatalf("SiguienteNumeroRonda() devolvió un error inesperado: %v", err)
				}
				if numero != 1 {
					t.Errorf("numero = %d, se esperaba %d", numero, 1)
				}
			})
		}
	})

	t.Run("con la última ronda Revelada devuelve max+1", func(t *testing.T) {
		casos := []struct {
			descripcion string
			rondas      []domain.RondaPoker
			esperado    int
		}{
			{
				descripcion: "una ronda revelada",
				rondas: []domain.RondaPoker{
					{Numero: 1, Estado: domain.EstadoRondaRevelada},
				},
				esperado: 2,
			},
			{
				descripcion: "varias rondas reveladas",
				rondas: []domain.RondaPoker{
					{Numero: 1, Estado: domain.EstadoRondaRevelada},
					{Numero: 2, Estado: domain.EstadoRondaRevelada},
				},
				esperado: 3,
			},
			{
				descripcion: "max+1 y no cantidad+1",
				rondas: []domain.RondaPoker{
					{Numero: 1, Estado: domain.EstadoRondaRevelada},
					{Numero: 3, Estado: domain.EstadoRondaRevelada},
				},
				esperado: 4,
			},
		}

		for _, caso := range casos {
			t.Run(caso.descripcion, func(t *testing.T) {
				numero, err := domain.SiguienteNumeroRonda(caso.rondas)
				if err != nil {
					t.Fatalf("SiguienteNumeroRonda() devolvió un error inesperado: %v", err)
				}
				if numero != caso.esperado {
					t.Errorf("numero = %d, se esperaba %d", numero, caso.esperado)
				}
			})
		}
	})

	t.Run("con una ronda Abierta devuelve ErrRondaAbiertaExistente", func(t *testing.T) {
		casos := []struct {
			descripcion string
			rondas      []domain.RondaPoker
		}{
			{
				descripcion: "única ronda abierta",
				rondas: []domain.RondaPoker{
					{Numero: 1, Estado: domain.EstadoRondaAbierta},
				},
			},
			{
				descripcion: "ronda abierta tras una revelada",
				rondas: []domain.RondaPoker{
					{Numero: 1, Estado: domain.EstadoRondaRevelada},
					{Numero: 2, Estado: domain.EstadoRondaAbierta},
				},
			},
		}

		for _, caso := range casos {
			t.Run(caso.descripcion, func(t *testing.T) {
				numero, err := domain.SiguienteNumeroRonda(caso.rondas)
				if !errors.Is(err, domain.ErrRondaAbiertaExistente) {
					t.Fatalf("err = %v, se esperaba %v", err, domain.ErrRondaAbiertaExistente)
				}
				if numero != 0 {
					t.Errorf("numero = %d, se esperaba %d", numero, 0)
				}
			})
		}
	})
}

// R7 / R22 (US-12): Revelar pasa una ronda Abierta a Revelada; una ronda ya
// Revelada devuelve ErrRondaRevelada (irreversible).
func TestUS12_Dominio_RondaPokerRevelar(t *testing.T) {
	t.Run("Abierta pasa a Revelada", func(t *testing.T) {
		ronda := domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaAbierta}

		revelada, err := ronda.Revelar()
		if err != nil {
			t.Fatalf("Revelar() devolvió un error inesperado: %v", err)
		}
		if revelada.Estado != domain.EstadoRondaRevelada {
			t.Errorf("Estado = %q, se esperaba %q", revelada.Estado, domain.EstadoRondaRevelada)
		}
		if revelada.Estado != "Revelada" {
			t.Errorf("Estado = %q, se esperaba %q", revelada.Estado, "Revelada")
		}
		if revelada.Numero != ronda.Numero {
			t.Errorf("Numero = %d, se esperaba %d", revelada.Numero, ronda.Numero)
		}
		if ronda.Estado != domain.EstadoRondaAbierta {
			t.Errorf("la ronda original cambió a %q, se esperaba que siguiera %q", ronda.Estado, domain.EstadoRondaAbierta)
		}
	})

	t.Run("ya Revelada devuelve ErrRondaRevelada", func(t *testing.T) {
		ronda := domain.RondaPoker{Numero: 1, Estado: domain.EstadoRondaRevelada}

		revelada, err := ronda.Revelar()
		if !errors.Is(err, domain.ErrRondaRevelada) {
			t.Fatalf("err = %v, se esperaba %v", err, domain.ErrRondaRevelada)
		}
		if revelada.Estado != "" {
			t.Errorf("revelada.Estado = %q, se esperaba %q", revelada.Estado, "")
		}
	})
}

// R8 / D2 (US-12): los estados de la ronda son exactamente "Abierta" y
// "Revelada".
func TestUS12_Dominio_EstadoRondaConstantes(t *testing.T) {
	if domain.EstadoRondaAbierta != "Abierta" {
		t.Errorf("EstadoRondaAbierta = %q, se esperaba %q", domain.EstadoRondaAbierta, "Abierta")
	}
	if domain.EstadoRondaRevelada != "Revelada" {
		t.Errorf("EstadoRondaRevelada = %q, se esperaba %q", domain.EstadoRondaRevelada, "Revelada")
	}
}
