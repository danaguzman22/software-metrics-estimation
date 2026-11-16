package domain_test

import (
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// itemUS06Orden arma un ítem mínimo para probar el ordenamiento del backlog.
func itemUS06Orden(id, prioridad string) domain.ItemBacklog {
	return domain.ItemBacklog{
		ID:        id,
		Titulo:    "Ítem " + id,
		Prioridad: prioridad,
		Estado:    "Pendiente",
	}
}

// verificarOrdenUS06 compara la secuencia de IDs devuelta por OrdenarBacklog.
func verificarOrdenUS06(t *testing.T, obtenido []domain.ItemBacklog, idsEsperados []string) {
	t.Helper()
	if len(obtenido) != len(idsEsperados) {
		t.Fatalf("len(OrdenarBacklog()) = %d, se esperaban %d", len(obtenido), len(idsEsperados))
	}
	for i, idEsperado := range idsEsperados {
		if obtenido[i].ID != idEsperado {
			t.Errorf("ID en posición %d = %q, se esperaba %q", i, obtenido[i].ID, idEsperado)
		}
	}
}

// CA2 / R5: el backlog se ordena por prioridad Alta, Media, Baja.
func TestUS06_OrdenarBacklog_PorPrioridad(t *testing.T) {
	entrada := []domain.ItemBacklog{
		itemUS06Orden("HU-03", "Baja"),
		itemUS06Orden("HU-01", "Alta"),
		itemUS06Orden("HU-02", "Media"),
	}

	obtenido := domain.OrdenarBacklog(entrada)
	verificarOrdenUS06(t, obtenido, []string{"HU-01", "HU-02", "HU-03"})
}

// CA2 / R6: a igual prioridad, el desempate es por número ascendente.
func TestUS06_OrdenarBacklog_DesempatePorNumero(t *testing.T) {
	entrada := []domain.ItemBacklog{
		itemUS06Orden("HU-03", "Media"),
		itemUS06Orden("HU-01", "Media"),
		itemUS06Orden("HU-02", "Media"),
	}

	obtenido := domain.OrdenarBacklog(entrada)
	verificarOrdenUS06(t, obtenido, []string{"HU-01", "HU-02", "HU-03"})
}

// CA2 / R6: la comparación es numérica y ubica HU-2 antes que HU-10.
func TestUS06_OrdenarBacklog_HU10VsHU2(t *testing.T) {
	entrada := []domain.ItemBacklog{
		itemUS06Orden("HU-10", "Alta"),
		itemUS06Orden("HU-2", "Alta"),
		itemUS06Orden("HU-1", "Alta"),
	}

	obtenido := domain.OrdenarBacklog(entrada)
	verificarOrdenUS06(t, obtenido, []string{"HU-1", "HU-2", "HU-10"})
}

// CA2 / R5 / R6: las prioridades desconocidas quedan después de Baja y entre sí
// se ordenan por número ascendente.
func TestUS06_OrdenarBacklog_PrioridadDesconocidaAlFinal(t *testing.T) {
	entrada := []domain.ItemBacklog{
		itemUS06Orden("HU-03", "Baja"),
		itemUS06Orden("HU-10", "Urgente"),
		itemUS06Orden("HU-05", "Alta"),
		itemUS06Orden("HU-02", "Urgente"),
		itemUS06Orden("HU-04", "Media"),
	}

	obtenido := domain.OrdenarBacklog(entrada)
	verificarOrdenUS06(t, obtenido, []string{"HU-05", "HU-04", "HU-03", "HU-02", "HU-10"})
}

// CA2 / R6: los IDs malformados quedan después de los válidos y entre sí se
// ordenan por ID como texto ascendente.
func TestUS06_OrdenarBacklog_IDMalformadoDespuesDeLosValidos(t *testing.T) {
	entrada := []domain.ItemBacklog{
		itemUS06Orden("ZZZ", "Alta"),
		itemUS06Orden("HU-10", "Alta"),
		itemUS06Orden("HU-X", "Alta"),
		itemUS06Orden("HU-01", "Alta"),
		itemUS06Orden("HU-0", "Alta"),
	}

	obtenido := domain.OrdenarBacklog(entrada)
	verificarOrdenUS06(t, obtenido, []string{"HU-01", "HU-10", "HU-0", "HU-X", "ZZZ"})
}

// CA2 / R11: una entrada vacía o nil devuelve una lista vacía.
func TestUS06_OrdenarBacklog_Vacio(t *testing.T) {
	casos := []struct {
		descripcion string
		entrada     []domain.ItemBacklog
	}{
		{descripcion: "lista vacía", entrada: []domain.ItemBacklog{}},
		{descripcion: "lista nil", entrada: nil},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			obtenido := domain.OrdenarBacklog(caso.entrada)
			if len(obtenido) != 0 {
				t.Errorf("len(OrdenarBacklog()) = %d, se esperaba 0", len(obtenido))
			}
		})
	}
}

// CA2 / R5 / R6: OrdenarBacklog devuelve una copia y no modifica la entrada.
func TestUS06_OrdenarBacklog_NoModificaEntrada(t *testing.T) {
	entrada := []domain.ItemBacklog{
		itemUS06Orden("HU-03", "Baja"),
		itemUS06Orden("HU-01", "Alta"),
		itemUS06Orden("HU-02", "Media"),
	}
	copia := append([]domain.ItemBacklog(nil), entrada...)

	obtenido := domain.OrdenarBacklog(entrada)
	verificarOrdenUS06(t, obtenido, []string{"HU-01", "HU-02", "HU-03"})
	verificarOrdenUS06(t, entrada, []string{"HU-03", "HU-01", "HU-02"})

	for i := range copia {
		if entrada[i].ID != copia[i].ID || entrada[i].Prioridad != copia[i].Prioridad {
			t.Errorf("la entrada cambió en posición %d: %+v, se esperaba %+v", i, entrada[i], copia[i])
		}
	}
}
