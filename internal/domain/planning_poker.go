package domain

import (
	"errors"
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

// Revelar marca una ronda como Revelada; si ya lo estaba, ErrRondaRevelada
// (irreversible, D6). La atomicidad real la garantiza el adaptador (R11).
func (r RondaPoker) Revelar() (RondaPoker, error) {
	if r.Estado == EstadoRondaRevelada {
		return RondaPoker{}, ErrRondaRevelada
	}
	return RondaPoker{Numero: r.Numero, Estado: EstadoRondaRevelada}, nil
}
