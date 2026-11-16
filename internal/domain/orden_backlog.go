package domain

import "sort"

// OrdenarBacklog devuelve una copia ordenada de los ítems: primero por
// prioridad Alta, Media, Baja; las prioridades desconocidas quedan después de
// Baja. A igual prioridad, por número de ítem ascendente usando NumeroDeItem;
// los IDs malformados quedan después de los válidos y entre sí por ID como
// texto ascendente. No modifica la slice de entrada.
func OrdenarBacklog(items []ItemBacklog) []ItemBacklog {
	copia := make([]ItemBacklog, len(items))
	copy(copia, items)

	sort.SliceStable(copia, func(a, b int) bool {
		return itemAntesQue(copia[a], copia[b])
	})

	return copia
}

// itemAntesQue informa si el ítem a debe ubicarse antes que el ítem b.
func itemAntesQue(a, b ItemBacklog) bool {
	rangoA := rangoPrioridad(a.Prioridad)
	rangoB := rangoPrioridad(b.Prioridad)
	if rangoA != rangoB {
		return rangoA < rangoB
	}

	numeroA, validoA := NumeroDeItem(a.ID)
	numeroB, validoB := NumeroDeItem(b.ID)
	if validoA && validoB {
		return numeroA < numeroB
	}
	if validoA != validoB {
		return validoA
	}
	return a.ID < b.ID
}

// rangoPrioridad devuelve la posición de la prioridad en el orden del backlog.
// Las prioridades desconocidas comparten el último rango.
func rangoPrioridad(prioridad string) int {
	switch prioridad {
	case "Alta":
		return 0
	case "Media":
		return 1
	case "Baja":
		return 2
	default:
		return 3
	}
}
