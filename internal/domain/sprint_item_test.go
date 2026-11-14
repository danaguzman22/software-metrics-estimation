package domain_test

import (
	"errors"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// itemUS08Dominio arma una historia del backlog para los tests de dominio de
// US-08.
func itemUS08Dominio(id string, storyPoints int, estado string) domain.ItemBacklog {
	return domain.ItemBacklog{
		ID:          id,
		Titulo:      "Historia " + id,
		Descripcion: "Historia de prueba de US-08",
		Prioridad:   "Alta",
		Estado:      estado,
		StoryPoints: storyPoints,
	}
}

// CA1 / R2: una historia con story points 0 o fuera de la serie no está
// estimada. La estimación se evalúa antes que el estado.
func TestUS08_ValidarItemAsignable_NoEstimado(t *testing.T) {
	casos := []struct {
		nombre string
		item   domain.ItemBacklog
	}{
		{
			nombre: "story points 0",
			item:   itemUS08Dominio("HU-04", 0, domain.EstadoItemPendiente),
		},
		{
			nombre: "story points 4 fuera de la serie",
			item:   itemUS08Dominio("HU-04", 4, domain.EstadoItemPendiente),
		},
		{
			nombre: "sin estimar y no pendiente devuelve primero la estimación",
			item:   itemUS08Dominio("HU-05", 0, "Completada"),
		},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			err := domain.ValidarItemAsignable(caso.item)
			if !errors.Is(err, domain.ErrItemNoEstimado) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrItemNoEstimado)
			}
			if errors.Is(err, domain.ErrItemNoPendiente) {
				t.Errorf("err = %v, no debería devolverse ErrItemNoPendiente antes que ErrItemNoEstimado", err)
			}
		})
	}
}

// CA1 / R3 / R13: solo el estado Pendiente es asignable.
func TestUS08_ValidarItemAsignable_NoPendiente(t *testing.T) {
	casos := []struct {
		nombre string
		estado string
	}{
		{nombre: "Completada", estado: "Completada"},
		{nombre: "En curso", estado: "En curso"},
		{nombre: "estado vacío", estado: ""},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			err := domain.ValidarItemAsignable(itemUS08Dominio("HU-05", 5, caso.estado))
			if !errors.Is(err, domain.ErrItemNoPendiente) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrItemNoPendiente)
			}
			if errors.Is(err, domain.ErrItemNoEstimado) {
				t.Errorf("err = %v, no debería devolverse ErrItemNoEstimado con story points válidos", err)
			}
		})
	}
}

// CA3 / R9: la suma incluye historias con IDs sin relleno de ceros ("HU-1",
// "HU-10") para no asumir el formato del ID.
func TestUS08_StoryPointsPlanificados_Suma(t *testing.T) {
	items := []domain.ItemBacklog{
		itemUS08Dominio("HU-1", 3, domain.EstadoItemPendiente),
		itemUS08Dominio("HU-10", 5, domain.EstadoItemPendiente),
		itemUS08Dominio("HU-2", 2, domain.EstadoItemPendiente),
	}

	total := domain.StoryPointsPlanificados(items)
	if total != 10 {
		t.Errorf("StoryPointsPlanificados() = %d, se esperaba 10", total)
	}
}

// CA3 / R9: un Sprint sin historias tiene 0 Story Points planificados.
func TestUS08_StoryPointsPlanificados_SinItems(t *testing.T) {
	casos := []struct {
		nombre string
		items  []domain.ItemBacklog
	}{
		{nombre: "lista nil", items: nil},
		{nombre: "lista vacía", items: []domain.ItemBacklog{}},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			if total := domain.StoryPointsPlanificados(caso.items); total != 0 {
				t.Errorf("StoryPointsPlanificados() = %d, se esperaba 0", total)
			}
		})
	}
}
