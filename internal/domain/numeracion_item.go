package domain

import (
	"fmt"
	"strconv"
)

const (
	prefijoIDItem = "HU-"
	formatoIDItem = "HU-%02d"
)

// NumeroDeItem extrae el número de un ID con el patrón "HU-NN".
// El sufijo debe ser solo dígitos ASCII 0-9, con valor mayor o igual a 1.
// Devuelve false si el ID no tiene ese patrón, incluidos signos, valor 0
// y overflow del sufijo.
func NumeroDeItem(id string) (int, bool) {
	if len(id) <= len(prefijoIDItem) || id[:len(prefijoIDItem)] != prefijoIDItem {
		return 0, false
	}

	sufijo := id[len(prefijoIDItem):]
	for i := 0; i < len(sufijo); i++ {
		if sufijo[i] < '0' || sufijo[i] > '9' {
			return 0, false
		}
	}

	numero, err := strconv.Atoi(sufijo)
	if err != nil || numero < 1 {
		return 0, false
	}
	return numero, true
}

// SiguienteIDItem devuelve el próximo ID de ítem del proyecto:
// "HU-01" si no hay ítems previos (o ninguno con patrón válido) y, si no,
// el mayor número existente más uno, con el formato actual "HU-%02d".
func SiguienteIDItem(idsExistentes []string) string {
	maximo := 0
	for _, id := range idsExistentes {
		if numero, ok := NumeroDeItem(id); ok && numero > maximo {
			maximo = numero
		}
	}
	return fmt.Sprintf(formatoIDItem, maximo+1)
}
