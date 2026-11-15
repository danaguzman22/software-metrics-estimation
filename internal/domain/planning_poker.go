package domain

import (
	"errors"
	"sort"
	"strings"
)

// Estados de la ronda (D2). Una ronda nace Abierta y solo Revelar la lleva a
// Revelada (R22).
const (
	EstadoRondaAbierta  = "Abierta"
	EstadoRondaRevelada = "Revelada"
)

// Errores de dominio de Planning Poker (D15).
var (
	// ErrIntegranteNoDelProyecto indica que el integrante no pertenece al
	// proyecto.
	ErrIntegranteNoDelProyecto = errors.New("el integrante no pertenece al proyecto")
	// ErrValorVotoPokerInvalido indica que el valor del voto no pertenece a
	// la serie 1, 2, 3, 5, 8, 13, 21.
	ErrValorVotoPokerInvalido = errors.New("el voto no es válido: debe ser 1, 2, 3, 5, 8, 13 o 21")
	// ErrNumeroRondaInvalido indica que el número de la ronda no es mayor
	// que cero.
	ErrNumeroRondaInvalido = errors.New("el número de la ronda debe ser mayor que cero")
	// ErrRondaNoEncontrada indica que la ronda no existe para la historia.
	ErrRondaNoEncontrada = errors.New("la ronda no existe para la historia")
	// ErrRondaAbiertaExistente indica que la historia ya tiene una ronda
	// Abierta.
	ErrRondaAbiertaExistente = errors.New("la historia ya tiene una ronda abierta")
	// ErrRondaRevelada indica que la ronda ya fue revelada.
	ErrRondaRevelada = errors.New("la ronda ya fue revelada")
	// ErrRondaNoRevelada indica que la ronda debe estar Revelada. Texto
	// neutro: aplica a acordar con ronda no Revelada (E6) y a
	// DiferenciasRonda con ronda Abierta (E9).
	ErrRondaNoRevelada = errors.New("la ronda aún no fue revelada")
	// ErrRondaSinVotos indica que la ronda Revelada no tiene votos.
	// Únicamente en DiferenciasRonda (D25, opción A).
	ErrRondaSinVotos = errors.New("la ronda no tiene votos")
)

// RondaPoker representa una ronda de votación sobre una historia.
type RondaPoker struct {
	Numero int
	Estado string
}

// VotoPoker representa el voto de un integrante en una ronda.
type VotoPoker struct {
	Integrante string
	Valor      int
}

// DiferenciaPoker describe las diferencias de una ronda revelada (CA1).
type DiferenciaPoker struct {
	HayDiferencia bool
	Minimo        int
	Maximo        int
	QuienesMinimo []string
	QuienesMaximo []string
}

// RondaHistorialPoker es una ronda del historial con sus votos; Votos queda
// nil en rondas Abierta (D4, R15).
type RondaHistorialPoker struct {
	Ronda RondaPoker
	Votos []VotoPoker
}

// ConsultaRondaPoker es el resultado de ConsultarRonda. Con la ronda Abierta
// solo trae YaVotaron (cantidad = len(YaVotaron)); con la ronda Revelada solo
// trae Votos (D4, R17).
type ConsultaRondaPoker struct {
	Ronda     RondaPoker
	YaVotaron []string
	Votos     []VotoPoker
}

// EsValorVotoPokerValido informa si el valor de voto es válido (D3).
func EsValorVotoPokerValido(v int) bool {
	switch v {
	case 1, 2, 3, 5, 8, 13, 21:
		return true
	default:
		return false
	}
}

// NuevoVotoPoker valida un voto: integrante no vacío (reutiliza
// ErrNombreIntegranteVacio) y valor válido. Ante cualquier error devuelve
// VotoPoker{}.
func NuevoVotoPoker(integrante string, valor int) (VotoPoker, error) {
	integrante = strings.TrimSpace(integrante)
	if integrante == "" {
		return VotoPoker{}, ErrNombreIntegranteVacio
	}
	if !EsValorVotoPokerValido(valor) {
		return VotoPoker{}, ErrValorVotoPokerInvalido
	}
	return VotoPoker{Integrante: integrante, Valor: valor}, nil
}

// SiguienteNumeroRonda calcula el número de la siguiente ronda de una historia:
// 1 si no hay rondas; max(numero)+1 si todas están Revelada;
// ErrRondaAbiertaExistente si hay una Abierta (D7, R9).
func SiguienteNumeroRonda(rondas []RondaPoker) (int, error) {
	mayor := 0
	for _, ronda := range rondas {
		if ronda.Estado == EstadoRondaAbierta {
			return 0, ErrRondaAbiertaExistente
		}
		if ronda.Numero > mayor {
			mayor = ronda.Numero
		}
	}
	return mayor + 1, nil
}

// DiferenciaDeVotos calcula si hay diferencia entre ≥1 votos. Cero votos
// devuelve ErrRondaSinVotos (R3). HayDiferencia es true si y solo si existe
// al menos un par de votos con valores distintos. QuienesMinimo y
// QuienesMaximo van ordenadas alfabéticamente sin distinguir mayúsculas,
// con desempate por el nombre exacto (R2, D16).
func DiferenciaDeVotos(votos []VotoPoker) (DiferenciaPoker, error) {
	if len(votos) == 0 {
		return DiferenciaPoker{}, ErrRondaSinVotos
	}
	minimo := votos[0].Valor
	maximo := votos[0].Valor
	quienesMinimo := []string{}
	quienesMaximo := []string{}
	for _, voto := range votos {
		if voto.Valor < minimo {
			minimo = voto.Valor
		}
		if voto.Valor > maximo {
			maximo = voto.Valor
		}
	}
	for _, voto := range votos {
		if voto.Valor == minimo {
			quienesMinimo = append(quienesMinimo, voto.Integrante)
		}
		if voto.Valor == maximo {
			quienesMaximo = append(quienesMaximo, voto.Integrante)
		}
	}
	ordenarNombresPoker(quienesMinimo)
	ordenarNombresPoker(quienesMaximo)
	return DiferenciaPoker{
		HayDiferencia: minimo != maximo,
		Minimo:        minimo,
		Maximo:        maximo,
		QuienesMinimo: quienesMinimo,
		QuienesMaximo: quienesMaximo,
	}, nil
}

// ordenarNombresPoker ordena los nombres sin distinguir mayúsculas y
// desempatando por el nombre exacto (D20).
func ordenarNombresPoker(nombres []string) {
	sort.Slice(nombres, func(i, j int) bool {
		ni, nj := strings.ToLower(nombres[i]), strings.ToLower(nombres[j])
		if ni == nj {
			return nombres[i] < nombres[j]
		}
		return ni < nj
	})
}

// Revelar marca una ronda como Revelada; si ya lo estaba, ErrRondaRevelada
// (irreversible, D6). La atomicidad real la garantiza el adaptador (R11).
func (r RondaPoker) Revelar() (RondaPoker, error) {
	if r.Estado == EstadoRondaRevelada {
		return RondaPoker{}, ErrRondaRevelada
	}
	return RondaPoker{Numero: r.Numero, Estado: EstadoRondaRevelada}, nil
}
