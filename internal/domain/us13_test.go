package domain_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func verificarDiferenciaPoker(t *testing.T, obtenida domain.DiferenciaPoker, esperada domain.DiferenciaPoker) {
	t.Helper()
	if obtenida.HayDiferencia != esperada.HayDiferencia {
		t.Errorf("HayDiferencia = %t, se esperaba %t", obtenida.HayDiferencia, esperada.HayDiferencia)
	}
	if obtenida.Minimo != esperada.Minimo {
		t.Errorf("Minimo = %d, se esperaba %d", obtenida.Minimo, esperada.Minimo)
	}
	if obtenida.Maximo != esperada.Maximo {
		t.Errorf("Maximo = %d, se esperaba %d", obtenida.Maximo, esperada.Maximo)
	}
	if !reflect.DeepEqual(obtenida.QuienesMinimo, esperada.QuienesMinimo) {
		t.Errorf("QuienesMinimo = %v, se esperaban %v", obtenida.QuienesMinimo, esperada.QuienesMinimo)
	}
	if !reflect.DeepEqual(obtenida.QuienesMaximo, esperada.QuienesMaximo) {
		t.Errorf("QuienesMaximo = %v, se esperaban %v", obtenida.QuienesMaximo, esperada.QuienesMaximo)
	}
}

// R2 / R3 / V8 / CA1 (US-13): una ronda revelada sin votos no informa nada;
// DiferenciaDeVotos devuelve domain.ErrRondaSinVotos y la diferencia en cero.
func TestUS13_Dominio_DiferenciaDeVotos_SinVotosDevuelveErrRondaSinVotos(t *testing.T) {
	diferencia, err := domain.DiferenciaDeVotos(nil)
	if !errors.Is(err, domain.ErrRondaSinVotos) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrRondaSinVotos)
	}
	if err != domain.ErrRondaSinVotos {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrRondaSinVotos)
	}
	verificarDiferenciaPoker(t, diferencia, domain.DiferenciaPoker{})

	diferencia, err = domain.DiferenciaDeVotos([]domain.VotoPoker{})
	if !errors.Is(err, domain.ErrRondaSinVotos) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrRondaSinVotos)
	}
	verificarDiferenciaPoker(t, diferencia, domain.DiferenciaPoker{})
}

// R2 / CA1 (US-13): un solo voto no registra diferencia; mínimo y máximo son
// ese valor y el votante figura en ambos extremos.
func TestUS13_Dominio_DiferenciaDeVotos_UnSoloVoto(t *testing.T) {
	votos := []domain.VotoPoker{{Integrante: "Dana Guzmán", Valor: 8}}

	diferencia, err := domain.DiferenciaDeVotos(votos)
	if err != nil {
		t.Fatalf("DiferenciaDeVotos() devolvió un error inesperado: %v", err)
	}
	verificarDiferenciaPoker(t, diferencia, domain.DiferenciaPoker{
		HayDiferencia: false,
		Minimo:        8,
		Maximo:        8,
		QuienesMinimo: []string{"Dana Guzmán"},
		QuienesMaximo: []string{"Dana Guzmán"},
	})
}

// R2 / CA1 (US-13): con todos los votos iguales no hay diferencia; mínimo y
// máximo son ese valor y todos los votantes figuran en ambos extremos.
func TestUS13_Dominio_DiferenciaDeVotos_TodosIguales(t *testing.T) {
	votos := []domain.VotoPoker{
		{Integrante: "Beto", Valor: 5},
		{Integrante: "ana", Valor: 5},
		{Integrante: "Carlos", Valor: 5},
	}

	diferencia, err := domain.DiferenciaDeVotos(votos)
	if err != nil {
		t.Fatalf("DiferenciaDeVotos() devolvió un error inesperado: %v", err)
	}
	verificarDiferenciaPoker(t, diferencia, domain.DiferenciaPoker{
		HayDiferencia: false,
		Minimo:        5,
		Maximo:        5,
		QuienesMinimo: []string{"ana", "Beto", "Carlos"},
		QuienesMaximo: []string{"ana", "Beto", "Carlos"},
	})
}

// R2 / CA1 (US-13): con votos distintos hay diferencia; mínimo y máximo son
// los extremos numéricos y se informa quién votó cada uno.
func TestUS13_Dominio_DiferenciaDeVotos_VotosDistintos(t *testing.T) {
	votos := []domain.VotoPoker{
		{Integrante: "Carlos", Valor: 13},
		{Integrante: "Dana Guzmán", Valor: 2},
		{Integrante: "Ana", Valor: 8},
	}

	diferencia, err := domain.DiferenciaDeVotos(votos)
	if err != nil {
		t.Fatalf("DiferenciaDeVotos() devolvió un error inesperado: %v", err)
	}
	verificarDiferenciaPoker(t, diferencia, domain.DiferenciaPoker{
		HayDiferencia: true,
		Minimo:        2,
		Maximo:        13,
		QuienesMinimo: []string{"Dana Guzmán"},
		QuienesMaximo: []string{"Carlos"},
	})
}

// R2 / CA1 (US-13): los empates en los extremos conservan a todos los votantes
// de cada extremo.
func TestUS13_Dominio_DiferenciaDeVotos_EmpatesEnAmbosExtremos(t *testing.T) {
	votos := []domain.VotoPoker{
		{Integrante: "Zoe", Valor: 13},
		{Integrante: "Ana", Valor: 5},
		{Integrante: "Beto", Valor: 5},
		{Integrante: "Carlos", Valor: 13},
	}

	diferencia, err := domain.DiferenciaDeVotos(votos)
	if err != nil {
		t.Fatalf("DiferenciaDeVotos() devolvió un error inesperado: %v", err)
	}
	verificarDiferenciaPoker(t, diferencia, domain.DiferenciaPoker{
		HayDiferencia: true,
		Minimo:        5,
		Maximo:        13,
		QuienesMinimo: []string{"Ana", "Beto"},
		QuienesMaximo: []string{"Carlos", "Zoe"},
	})
}

// R2 / CA1 (US-13): HayDiferencia es true si y solo si existe al menos un par
// de votos con valores distintos.
func TestUS13_Dominio_DiferenciaDeVotos_HayDiferenciaSiYSoloSiHayValoresDistintos(t *testing.T) {
	casos := []struct {
		descripcion   string
		votos         []domain.VotoPoker
		hayDiferencia bool
		minimo        int
		maximo        int
	}{
		{
			descripcion: "dos valores distintos en el mínimo y el resto iguales",
			votos: []domain.VotoPoker{
				{Integrante: "Ana", Valor: 5},
				{Integrante: "Beto", Valor: 13},
				{Integrante: "Carlos", Valor: 5},
			},
			hayDiferencia: true,
			minimo:        5,
			maximo:        13,
		},
		{
			descripcion: "tres valores distintos",
			votos: []domain.VotoPoker{
				{Integrante: "Ana", Valor: 2},
				{Integrante: "Beto", Valor: 8},
				{Integrante: "Carlos", Valor: 21},
			},
			hayDiferencia: true,
			minimo:        2,
			maximo:        21,
		},
		{
			descripcion: "todos los votos iguales de varios votantes",
			votos: []domain.VotoPoker{
				{Integrante: "Ana", Valor: 8},
				{Integrante: "Beto", Valor: 8},
				{Integrante: "Carlos", Valor: 8},
				{Integrante: "Zoe", Valor: 8},
			},
			hayDiferencia: false,
			minimo:        8,
			maximo:        8,
		},
	}
	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			diferencia, err := domain.DiferenciaDeVotos(caso.votos)
			if err != nil {
				t.Fatalf("DiferenciaDeVotos() devolvió un error inesperado: %v", err)
			}
			if diferencia.HayDiferencia != caso.hayDiferencia {
				t.Errorf("HayDiferencia = %t, se esperaba %t", diferencia.HayDiferencia, caso.hayDiferencia)
			}
			if diferencia.Minimo != caso.minimo {
				t.Errorf("Minimo = %d, se esperaba %d", diferencia.Minimo, caso.minimo)
			}
			if diferencia.Maximo != caso.maximo {
				t.Errorf("Maximo = %d, se esperaba %d", diferencia.Maximo, caso.maximo)
			}
		})
	}
}

// R2 / D16 (US-13): QuienesMinimo y QuienesMaximo van ordenadas alfabéticamente
// sin distinguir mayúsculas y desempate por el nombre exacto (D20 de US-12).
func TestUS13_Dominio_DiferenciaDeVotos_OrdenAlfabeticoSinMayusculasYDesempateExacto(t *testing.T) {
	votos := []domain.VotoPoker{
		{Integrante: "Carlos", Valor: 13},
		{Integrante: "beto", Valor: 5},
		{Integrante: "Ana", Valor: 5},
		{Integrante: "ANA", Valor: 5},
		{Integrante: "Zoe", Valor: 13},
	}

	diferencia, err := domain.DiferenciaDeVotos(votos)
	if err != nil {
		t.Fatalf("DiferenciaDeVotos() devolvió un error inesperado: %v", err)
	}
	verificarDiferenciaPoker(t, diferencia, domain.DiferenciaPoker{
		HayDiferencia: true,
		Minimo:        5,
		Maximo:        13,
		QuienesMinimo: []string{"ANA", "Ana", "beto"},
		QuienesMaximo: []string{"Carlos", "Zoe"},
	})
}

// R2 / D20 (US-13): con los votos ingresados en cualquier orden, las listas de
// votantes de los extremos quedan ordenadas.
func TestUS13_Dominio_DiferenciaDeVotos_VotosEnCualquierOrden(t *testing.T) {
	votos := []domain.VotoPoker{
		{Integrante: "Zoe", Valor: 13},
		{Integrante: "Ana", Valor: 5},
		{Integrante: "Carlos", Valor: 13},
		{Integrante: "Beto", Valor: 5},
	}

	diferencia, err := domain.DiferenciaDeVotos(votos)
	if err != nil {
		t.Fatalf("DiferenciaDeVotos() devolvió un error inesperado: %v", err)
	}
	if !reflect.DeepEqual(diferencia.QuienesMinimo, []string{"Ana", "Beto"}) {
		t.Errorf("QuienesMinimo = %v, se esperaban [Ana Beto]", diferencia.QuienesMinimo)
	}
	if !reflect.DeepEqual(diferencia.QuienesMaximo, []string{"Carlos", "Zoe"}) {
		t.Errorf("QuienesMaximo = %v, se esperaban [Carlos Zoe]", diferencia.QuienesMaximo)
	}
}

// D6 / D7 (US-13): los errores de dominio nuevos existen y son distintos entre sí.
func TestUS13_Dominio_ErroresNuevos(t *testing.T) {
	if domain.ErrRondaNoRevelada == nil {
		t.Fatal("domain.ErrRondaNoRevelada es nil")
	}
	if domain.ErrRondaSinVotos == nil {
		t.Fatal("domain.ErrRondaSinVotos es nil")
	}
	if domain.ErrRondaNoRevelada == domain.ErrRondaSinVotos {
		t.Error("ErrRondaNoRevelada y ErrRondaSinVotos deben ser errores distintos")
	}
	if errors.Is(domain.ErrRondaNoRevelada, domain.ErrRondaSinVotos) {
		t.Error("errors.Is(ErrRondaNoRevelada, ErrRondaSinVotos) = true, se esperaba false")
	}
}
