package domain

import (
	"errors"
	"testing"
)

func TestUS09_MarcarComoCompletada_Exitoso(t *testing.T) {
	original := ItemBacklog{
		ID:          "HU-01",
		Titulo:      "Iniciar sesión",
		Descripcion: "Descripción",
		Prioridad:   "Alta",
		Estado:      EstadoItemPendiente,
		StoryPoints: 3,
		Criterios:   []string{"criterio"},
	}

	completado, err := original.MarcarComoCompletada()
	if err != nil {
		t.Fatalf("MarcarComoCompletada() devolvió un error inesperado: %v", err)
	}
	if completado.Estado != EstadoCompletada {
		t.Errorf("Estado = %q, se esperaba %q", completado.Estado, EstadoCompletada)
	}
	if original.Estado != EstadoItemPendiente {
		t.Errorf("el ítem original cambió a %q", original.Estado)
	}
	completado.Criterios[0] = "modificado"
	if original.Criterios[0] != "criterio" {
		t.Error("la copia comparte el slice Criterios con el ítem original")
	}
}

func TestUS09_MarcarComoCompletada_ItemNoPendiente(t *testing.T) {
	item := ItemBacklog{ID: "HU-01", Estado: EstadoCompletada}

	completado, err := item.MarcarComoCompletada()
	if !errors.Is(err, ErrItemNoPendiente) {
		t.Fatalf("err = %v, se esperaba %v", err, ErrItemNoPendiente)
	}
	if completado.ID != "" || completado.Titulo != "" || completado.Estado != "" ||
		completado.Descripcion != "" || completado.Prioridad != "" ||
		completado.StoryPoints != 0 || completado.Criterios != nil {
		t.Errorf("resultado ante error = %+v, se esperaba ItemBacklog{}", completado)
	}
}
