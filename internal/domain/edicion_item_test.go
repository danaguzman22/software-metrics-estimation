package domain_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// itemUS06Base arma el ítem de referencia de los escenarios de US-06.
func itemUS06Base() domain.ItemBacklog {
	return domain.ItemBacklog{
		ID:          "HU-01",
		Titulo:      "Registrar proyecto",
		Descripcion: "Permitir dar de alta proyectos en el sistema",
		Prioridad:   "Alta",
		Estado:      "Pendiente",
		StoryPoints: 5,
		Criterios: []string{
			"El nombre del proyecto no puede estar vacío",
		},
	}
}

// verificarItemUS06 compara campo por campo un ítem del backlog.
func verificarItemUS06(t *testing.T, obtenido, esperado domain.ItemBacklog) {
	t.Helper()
	if obtenido.ID != esperado.ID {
		t.Errorf("ID = %q, se esperaba %q", obtenido.ID, esperado.ID)
	}
	if obtenido.Titulo != esperado.Titulo {
		t.Errorf("Titulo = %q, se esperaba %q", obtenido.Titulo, esperado.Titulo)
	}
	if obtenido.Descripcion != esperado.Descripcion {
		t.Errorf("Descripcion = %q, se esperaba %q", obtenido.Descripcion, esperado.Descripcion)
	}
	if obtenido.Prioridad != esperado.Prioridad {
		t.Errorf("Prioridad = %q, se esperaba %q", obtenido.Prioridad, esperado.Prioridad)
	}
	if obtenido.Estado != esperado.Estado {
		t.Errorf("Estado = %q, se esperaba %q", obtenido.Estado, esperado.Estado)
	}
	if obtenido.StoryPoints != esperado.StoryPoints {
		t.Errorf("StoryPoints = %d, se esperaba %d", obtenido.StoryPoints, esperado.StoryPoints)
	}
	if len(obtenido.Criterios) != len(esperado.Criterios) {
		t.Errorf("len(Criterios) = %d, se esperaba %d", len(obtenido.Criterios), len(esperado.Criterios))
		return
	}
	for i := range esperado.Criterios {
		if obtenido.Criterios[i] != esperado.Criterios[i] {
			t.Errorf("Criterios[%d] = %q, se esperaba %q", i, obtenido.Criterios[i], esperado.Criterios[i])
		}
	}
}

// CA1 / R1: editar el título actualiza solo ese campo y recorta los espacios
// externos.
func TestUS06_Editar_Titulo(t *testing.T) {
	casos := []struct {
		descripcion string
		titulo      string
		esperado    string
	}{
		{
			descripcion: "cambio de título",
			titulo:      "Registrar y listar proyectos",
			esperado:    "Registrar y listar proyectos",
		},
		{
			descripcion: "título con espacios externos se recorta",
			titulo:      "  Registrar y listar proyectos  ",
			esperado:    "Registrar y listar proyectos",
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			original := itemUS06Base()

			editado, err := original.Editar(caso.titulo, original.Descripcion, original.Prioridad, original.Criterios, original.StoryPoints)
			if err != nil {
				t.Fatalf("Editar() devolvió un error inesperado: %v", err)
			}

			esperado := original
			esperado.Titulo = caso.esperado
			verificarItemUS06(t, editado, esperado)
		})
	}
}

// CA1 / R1: editar la descripción actualiza solo ese campo.
func TestUS06_Editar_Descripcion(t *testing.T) {
	original := itemUS06Base()

	editado, err := original.Editar(
		original.Titulo,
		"Permitir dar de alta y consultar proyectos en el sistema",
		original.Prioridad,
		original.Criterios,
		original.StoryPoints,
	)
	if err != nil {
		t.Fatalf("Editar() devolvió un error inesperado: %v", err)
	}

	esperado := original
	esperado.Descripcion = "Permitir dar de alta y consultar proyectos en el sistema"
	verificarItemUS06(t, editado, esperado)
}

// CA1 / R1: editar la prioridad actualiza solo ese campo y recorta los espacios
// externos.
func TestUS06_Editar_Prioridad(t *testing.T) {
	casos := []struct {
		descripcion string
		prioridad   string
		esperada    string
	}{
		{
			descripcion: "cambio a Media",
			prioridad:   "Media",
			esperada:    "Media",
		},
		{
			descripcion: "prioridad con espacios externos se recorta",
			prioridad:   "  Media  ",
			esperada:    "Media",
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			original := itemUS06Base()

			editado, err := original.Editar(original.Titulo, original.Descripcion, caso.prioridad, original.Criterios, original.StoryPoints)
			if err != nil {
				t.Fatalf("Editar() devolvió un error inesperado: %v", err)
			}

			esperado := original
			esperado.Prioridad = caso.esperada
			verificarItemUS06(t, editado, esperado)
		})
	}
}

// CA1 / R1: editar los criterios de aceptación los guarda tal como se reciben.
func TestUS06_Editar_Criterios(t *testing.T) {
	original := itemUS06Base()
	criterios := []string{
		"El nombre del proyecto no puede estar vacío",
		"El proyecto queda disponible para asignar integrantes",
		"El proyecto se lista en el backlog",
	}

	editado, err := original.Editar(original.Titulo, original.Descripcion, original.Prioridad, criterios, original.StoryPoints)
	if err != nil {
		t.Fatalf("Editar() devolvió un error inesperado: %v", err)
	}

	esperado := original
	esperado.Criterios = criterios
	verificarItemUS06(t, editado, esperado)
}

// CA1 / R1 / R3: editar los story points admite todos los valores de la escala.
func TestUS06_Editar_StoryPoints(t *testing.T) {
	for _, storyPoints := range []int{1, 2, 3, 5, 8, 13, 21} {
		t.Run(strconv.Itoa(storyPoints), func(t *testing.T) {
			original := itemUS06Base()

			editado, err := original.Editar(original.Titulo, original.Descripcion, original.Prioridad, original.Criterios, storyPoints)
			if err != nil {
				t.Fatalf("Editar() devolvió un error inesperado: %v", err)
			}

			esperado := original
			esperado.StoryPoints = storyPoints
			verificarItemUS06(t, editado, esperado)
		})
	}
}

// CA1 / R1 / R3: editar un ítem sin estimar con 0 conserva los story points en 0.
func TestUS06_Editar_ItemSinEstimarConservaCero(t *testing.T) {
	original := itemUS06Base()
	original.StoryPoints = 0

	editado, err := original.Editar("Registrar y listar proyectos", original.Descripcion, original.Prioridad, original.Criterios, 0)
	if err != nil {
		t.Fatalf("Editar() devolvió un error inesperado: %v", err)
	}

	esperado := original
	esperado.Titulo = "Registrar y listar proyectos"
	esperado.StoryPoints = 0
	verificarItemUS06(t, editado, esperado)
}

// CA1 / R1 / R3: editar con 0 deja el ítem sin estimar.
func TestUS06_Editar_DejarSinEstimarConCero(t *testing.T) {
	original := itemUS06Base()

	editado, err := original.Editar(original.Titulo, original.Descripcion, original.Prioridad, original.Criterios, 0)
	if err != nil {
		t.Fatalf("Editar() devolvió un error inesperado: %v", err)
	}

	esperado := original
	esperado.StoryPoints = 0
	verificarItemUS06(t, editado, esperado)
}

// CA1 / R1 / R2: editar todos los campos editables a la vez actualiza los cinco
// y conserva ID y estado.
func TestUS06_Editar_TodosLosCampos(t *testing.T) {
	original := itemUS06Base()

	editado, err := original.Editar(
		"Registrar y listar proyectos",
		"Permitir dar de alta y consultar proyectos en el sistema",
		"Media",
		[]string{
			"El nombre del proyecto no puede estar vacío",
			"El proyecto se lista en el backlog",
		},
		8,
	)
	if err != nil {
		t.Fatalf("Editar() devolvió un error inesperado: %v", err)
	}

	esperado := domain.ItemBacklog{
		ID:          original.ID,
		Titulo:      "Registrar y listar proyectos",
		Descripcion: "Permitir dar de alta y consultar proyectos en el sistema",
		Prioridad:   "Media",
		Estado:      original.Estado,
		StoryPoints: 8,
		Criterios: []string{
			"El nombre del proyecto no puede estar vacío",
			"El proyecto se lista en el backlog",
		},
	}
	verificarItemUS06(t, editado, esperado)
}

// CA1 / R2: el ID y el estado no cambian al editar.
func TestUS06_Editar_PreservaIDYEstado(t *testing.T) {
	casos := []struct {
		descripcion string
		estado      string
	}{
		{
			descripcion: "estado Pendiente",
			estado:      "Pendiente",
		},
		{
			descripcion: "estado En curso",
			estado:      "En curso",
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			original := itemUS06Base()
			original.Estado = caso.estado

			editado, err := original.Editar(
				"Registrar y listar proyectos",
				"Permitir dar de alta y consultar proyectos en el sistema",
				"Media",
				[]string{"El proyecto se lista en el backlog"},
				8,
			)
			if err != nil {
				t.Fatalf("Editar() devolvió un error inesperado: %v", err)
			}

			if editado.ID != original.ID {
				t.Errorf("ID = %q, se esperaba conservar %q", editado.ID, original.ID)
			}
			if editado.Estado != caso.estado {
				t.Errorf("Estado = %q, se esperaba conservar %q", editado.Estado, caso.estado)
			}
		})
	}
}

// CA1 / R3 / V4: un título vacío devuelve ErrTituloItemVacio.
func TestUS06_Editar_TituloVacioDevuelveError(t *testing.T) {
	original := itemUS06Base()

	editado, err := original.Editar("", original.Descripcion, original.Prioridad, original.Criterios, original.StoryPoints)
	if err == nil {
		t.Fatal("Editar() con título vacío debería devolver un error de validación")
	}
	if !errors.Is(err, domain.ErrTituloItemVacio) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrTituloItemVacio)
	}
	verificarItemBacklogVacio(t, editado)
}

// CA1 / R3 / V4: un título compuesto solo por espacios devuelve
// ErrTituloItemVacio.
func TestUS06_Editar_TituloSoloEspaciosDevuelveError(t *testing.T) {
	original := itemUS06Base()

	editado, err := original.Editar("   ", original.Descripcion, original.Prioridad, original.Criterios, original.StoryPoints)
	if err == nil {
		t.Fatal("Editar() con título solo de espacios debería devolver un error de validación")
	}
	if !errors.Is(err, domain.ErrTituloItemVacio) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrTituloItemVacio)
	}
	verificarItemBacklogVacio(t, editado)
}

// CA1 / R3 / V5: una prioridad inválida devuelve ErrPrioridadInvalida.
func TestUS06_Editar_PrioridadInvalidaDevuelveError(t *testing.T) {
	casos := []struct {
		descripcion string
		prioridad   string
	}{
		{descripcion: "prioridad no permitida", prioridad: "Urgente"},
		{descripcion: "prioridad en minúsculas", prioridad: "alta"},
		{descripcion: "prioridad en mayúsculas", prioridad: "ALTA"},
		{descripcion: "prioridad compuesta solo por espacios", prioridad: "   "},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			original := itemUS06Base()

			editado, err := original.Editar(original.Titulo, original.Descripcion, caso.prioridad, original.Criterios, original.StoryPoints)
			if err == nil {
				t.Fatal("Editar() con prioridad inválida debería devolver un error de validación")
			}
			if !errors.Is(err, domain.ErrPrioridadInvalida) {
				t.Errorf("err = %v, se esperaba %v", err, domain.ErrPrioridadInvalida)
			}
			verificarItemBacklogVacio(t, editado)
		})
	}
}

// CA1 / R3 / V6: story points negativos, 4 y 22 devuelven
// ErrStoryPointsInvalidos.
func TestUS06_Editar_StoryPointsInvalidosDevuelveError(t *testing.T) {
	for _, storyPoints := range []int{-1, 4, 22} {
		t.Run(strconv.Itoa(storyPoints), func(t *testing.T) {
			original := itemUS06Base()

			editado, err := original.Editar(original.Titulo, original.Descripcion, original.Prioridad, original.Criterios, storyPoints)
			if err == nil {
				t.Fatal("Editar() con story points inválidos debería devolver un error de validación")
			}
			if !errors.Is(err, domain.ErrStoryPointsInvalidos) {
				t.Errorf("err = %v, se esperaba %v", err, domain.ErrStoryPointsInvalidos)
			}
			verificarItemBacklogVacio(t, editado)
		})
	}
}

// CA3 / R4 / V3: un ítem completado no se puede editar.
func TestUS06_Editar_ItemCompletadoDevuelveError(t *testing.T) {
	original := itemUS06Base()
	original.Estado = domain.EstadoCompletada

	editado, err := original.Editar(
		"Registrar y listar proyectos",
		"Permitir dar de alta y consultar proyectos en el sistema",
		"Media",
		[]string{"El proyecto se lista en el backlog"},
		8,
	)
	if err == nil {
		t.Fatal("Editar() sobre un ítem completado debería devolver un error de dominio")
	}
	if !errors.Is(err, domain.ErrItemCompletadoNoEditable) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrItemCompletadoNoEditable)
	}
	verificarItemBacklogVacio(t, editado)
}

// CA3 / R4 / R12: el error de ítem completado prevalece sobre las validaciones
// de los campos editables.
func TestUS06_Editar_ItemCompletadoAntesQueValidaciones(t *testing.T) {
	original := itemUS06Base()
	original.Estado = domain.EstadoCompletada

	editado, err := original.Editar("", original.Descripcion, "Urgente", original.Criterios, 4)
	if err == nil {
		t.Fatal("Editar() sobre un ítem completado debería devolver un error de dominio")
	}
	if !errors.Is(err, domain.ErrItemCompletadoNoEditable) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrItemCompletadoNoEditable)
	}
	if errors.Is(err, domain.ErrTituloItemVacio) ||
		errors.Is(err, domain.ErrPrioridadInvalida) ||
		errors.Is(err, domain.ErrStoryPointsInvalidos) {
		t.Errorf("err = %v, no debería devolverse un error de validación de campos", err)
	}
	verificarItemBacklogVacio(t, editado)
}

// CA1 / R1 / R2: la edición devuelve una copia y no modifica el ítem original.
func TestUS06_Editar_NoModificaItemOriginal(t *testing.T) {
	original := itemUS06Base()
	esperadoOriginal := itemUS06Base()

	editado, err := original.Editar(
		"Registrar y listar proyectos",
		"Permitir dar de alta y consultar proyectos en el sistema",
		"Media",
		[]string{"El proyecto se lista en el backlog"},
		8,
	)
	if err != nil {
		t.Fatalf("Editar() devolvió un error inesperado: %v", err)
	}

	verificarItemUS06(t, original, esperadoOriginal)
	if editado.Titulo == original.Titulo {
		t.Error("el ítem editado debería tener el título nuevo")
	}
}
