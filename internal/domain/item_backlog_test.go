package domain_test

import (
	"errors"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// verificarItemBacklogVacio comprueba que un error de validación devuelve el
// valor cero de domain.ItemBacklog.
func verificarItemBacklogVacio(t *testing.T, item domain.ItemBacklog) {
	t.Helper()
	if item.ID != "" || item.Titulo != "" || item.Descripcion != "" ||
		item.Prioridad != "" || item.Estado != "" || item.StoryPoints != 0 ||
		len(item.Criterios) != 0 {
		t.Errorf("se esperaba domain.ItemBacklog{} ante un error, se obtuvo %+v", item)
	}
}

// CA3 / R3 / V2: un título vacío o compuesto solo por espacios devuelve
// ErrTituloItemVacio.
func TestUS05_TituloVacio(t *testing.T) {
	casos := []struct {
		descripcion string
		titulo      string
	}{
		{
			descripcion: "título vacío",
			titulo:      "",
		},
		{
			// D4: un título compuesto únicamente por espacios se considera vacío.
			descripcion: "título compuesto solo por espacios",
			titulo:      "   ",
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			item, err := domain.NuevoItemBacklog(caso.titulo, "Permitir dar de alta proyectos en el sistema", "Alta", nil)
			if err == nil {
				t.Fatal("NuevoItemBacklog() con título vacío o solo espacios debería devolver un error de validación")
			}
			if !errors.Is(err, domain.ErrTituloItemVacio) {
				t.Errorf("err = %v, se esperaba %v", err, domain.ErrTituloItemVacio)
			}
			verificarItemBacklogVacio(t, item)
		})
	}
}

// CA1 / CA3 / R3 / D4: el título se guarda sin espacios al inicio y al final.
func TestUS05_NuevoItemBacklog_RecortaTitulo(t *testing.T) {
	item, err := domain.NuevoItemBacklog(
		"  Registrar proyecto  ",
		"Permitir dar de alta proyectos en el sistema",
		"Alta",
		nil,
	)
	if err != nil {
		t.Fatalf("NuevoItemBacklog() devolvió un error inesperado: %v", err)
	}
	if item.Titulo != "Registrar proyecto" {
		t.Errorf("Titulo = %q, se esperaba %q", item.Titulo, "Registrar proyecto")
	}
}

// CA3 / R4 / V3: una prioridad distinta de Alta, Media o Baja devuelve
// ErrPrioridadInvalida.
func TestUS05_PrioridadInvalida(t *testing.T) {
	casos := []struct {
		descripcion string
		prioridad   string
	}{
		{
			descripcion: "prioridad no permitida",
			prioridad:   "Urgente",
		},
		{
			// D5: la comparación distingue mayúsculas.
			descripcion: "prioridad en minúsculas",
			prioridad:   "alta",
		},
		{
			descripcion: "prioridad en mayúsculas",
			prioridad:   "ALTA",
		},
		{
			// D5: una prioridad compuesta solo por espacios no es válida.
			descripcion: "prioridad compuesta solo por espacios",
			prioridad:   " ",
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			item, err := domain.NuevoItemBacklog("Registrar proyecto", "Permitir dar de alta proyectos en el sistema", caso.prioridad, nil)
			if err == nil {
				t.Fatal("NuevoItemBacklog() con prioridad inválida debería devolver un error de validación")
			}
			if !errors.Is(err, domain.ErrPrioridadInvalida) {
				t.Errorf("err = %v, se esperaba %v", err, domain.ErrPrioridadInvalida)
			}
			verificarItemBacklogVacio(t, item)
		})
	}
}

// CA1 / CA3 / R4 / D5: la prioridad válida se guarda recortada.
func TestUS05_NuevoItemBacklog_PrioridadRecortada(t *testing.T) {
	casos := []struct {
		descripcion string
		entrada     string
		esperada    string
	}{
		{
			descripcion: "Alta con espacios externos",
			entrada:     "  Alta  ",
			esperada:    "Alta",
		},
		{
			descripcion: "Media con espacios externos",
			entrada:     "  Media  ",
			esperada:    "Media",
		},
		{
			descripcion: "Baja con espacios externos",
			entrada:     "  Baja  ",
			esperada:    "Baja",
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			item, err := domain.NuevoItemBacklog("Registrar proyecto", "Permitir dar de alta proyectos en el sistema", caso.entrada, nil)
			if err != nil {
				t.Fatalf("NuevoItemBacklog() devolvió un error inesperado: %v", err)
			}
			if item.Prioridad != caso.esperada {
				t.Errorf("Prioridad = %q, se esperaba %q", item.Prioridad, caso.esperada)
			}
		})
	}
}

// CA1 / R4 / D5: cada prioridad permitida es válida sin espacios externos.
func TestUS05_NuevoItemBacklog_PrioridadesValidas(t *testing.T) {
	casos := []struct {
		descripcion string
		prioridad   string
	}{
		{descripcion: "prioridad Alta", prioridad: "Alta"},
		{descripcion: "prioridad Media", prioridad: "Media"},
		{descripcion: "prioridad Baja", prioridad: "Baja"},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			item, err := domain.NuevoItemBacklog("Registrar proyecto", "Permitir dar de alta proyectos en el sistema", caso.prioridad, nil)
			if err != nil {
				t.Fatalf("NuevoItemBacklog() devolvió un error inesperado: %v", err)
			}
			if item.Prioridad != caso.prioridad {
				t.Errorf("Prioridad = %q, se esperaba %q", item.Prioridad, caso.prioridad)
			}
		})
	}
}

// CA4 / R6 / D6: NuevoItemBacklog fija el estado "Pendiente" sin recibirlo.
func TestUS05_NuevoItemBacklog_EstadoInicialPendiente(t *testing.T) {
	item, err := domain.NuevoItemBacklog("Registrar proyecto", "Permitir dar de alta proyectos en el sistema", "Alta", nil)
	if err != nil {
		t.Fatalf("NuevoItemBacklog() devolvió un error inesperado: %v", err)
	}
	if item.Estado != "Pendiente" {
		t.Errorf("Estado = %q, se esperaba %q", item.Estado, "Pendiente")
	}
}

// CA1 / R7 / D7: NuevoItemBacklog fija los story points en 0 sin recibirlos.
func TestUS05_NuevoItemBacklog_StoryPointsCero(t *testing.T) {
	item, err := domain.NuevoItemBacklog("Registrar proyecto", "Permitir dar de alta proyectos en el sistema", "Alta", nil)
	if err != nil {
		t.Fatalf("NuevoItemBacklog() devolvió un error inesperado: %v", err)
	}
	if item.StoryPoints != 0 {
		t.Errorf("StoryPoints = %d, se esperaba 0", item.StoryPoints)
	}
}
