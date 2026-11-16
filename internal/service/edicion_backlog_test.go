package service

import (
	"errors"
	"strconv"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// itemUS06ServicioBase arma el ítem de referencia de los escenarios de US-06.
func itemUS06ServicioBase() domain.ItemBacklog {
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

// itemUS06ServicioLista arma un ítem mínimo para las pruebas de listado.
func itemUS06ServicioLista(id, prioridad string) domain.ItemBacklog {
	return domain.ItemBacklog{
		ID:          id,
		Titulo:      "Ítem " + id,
		Descripcion: "Descripción " + id,
		Prioridad:   prioridad,
		Estado:      "Pendiente",
		StoryPoints: 5,
	}
}

// sembrarEntornoUS06 crea un entorno con los proyectos p1 y p2 y siembra los
// ítems indicados.
func sembrarEntornoUS06(t *testing.T, items ...itemConProyecto) *entornoBacklog {
	t.Helper()
	almacen := nuevoEntornoBacklog()
	sembrarProyectos(t, almacen, []ProyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
		{ID: "p2", Proyecto: proyectoUS02Servicio("App Móvil", domain.EstadoPlanificado)},
	})
	sembrarBacklog(t, almacen, items)
	return almacen
}

// verificarIDsUS06 compara la secuencia de IDs de una lista de ítems.
func verificarIDsUS06(t *testing.T, items []domain.ItemBacklog, idsEsperados []string) {
	t.Helper()
	if len(items) != len(idsEsperados) {
		t.Fatalf("len(items) = %d, se esperaban %d", len(items), len(idsEsperados))
	}
	for i, idEsperado := range idsEsperados {
		if items[i].ID != idEsperado {
			t.Errorf("ID en posición %d = %q, se esperaba %q", i, items[i].ID, idEsperado)
		}
	}
}

// CA1 / R1: editar el título actualiza el ítem devuelto y el persistido.
func TestUS06_EditarItem_Titulo(t *testing.T) {
	original := itemUS06ServicioBase()
	almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p1", Item: original})
	servicio := almacen.servicio()

	editado, err := servicio.EditarItem(
		"p1",
		"HU-01",
		"Registrar y listar proyectos",
		original.Descripcion,
		original.Prioridad,
		original.Criterios,
		original.StoryPoints,
	)
	if err != nil {
		t.Fatalf("EditarItem() devolvió un error inesperado: %v", err)
	}

	esperado := original
	esperado.Titulo = "Registrar y listar proyectos"
	verificarItemBacklog(t, editado, esperado)
	verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacen), "p1", "HU-01").Item, esperado)
}

// CA1 / R1: editar la descripción actualiza el ítem devuelto y el persistido.
func TestUS06_EditarItem_Descripcion(t *testing.T) {
	original := itemUS06ServicioBase()
	almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p1", Item: original})
	servicio := almacen.servicio()

	editado, err := servicio.EditarItem(
		"p1",
		"HU-01",
		original.Titulo,
		"Permitir dar de alta y consultar proyectos en el sistema",
		original.Prioridad,
		original.Criterios,
		original.StoryPoints,
	)
	if err != nil {
		t.Fatalf("EditarItem() devolvió un error inesperado: %v", err)
	}

	esperado := original
	esperado.Descripcion = "Permitir dar de alta y consultar proyectos en el sistema"
	verificarItemBacklog(t, editado, esperado)
	verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacen), "p1", "HU-01").Item, esperado)
}

// CA1 / R1: editar la prioridad actualiza el ítem devuelto y el persistido.
func TestUS06_EditarItem_Prioridad(t *testing.T) {
	original := itemUS06ServicioBase()
	almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p1", Item: original})
	servicio := almacen.servicio()

	editado, err := servicio.EditarItem(
		"p1",
		"HU-01",
		original.Titulo,
		original.Descripcion,
		"Media",
		original.Criterios,
		original.StoryPoints,
	)
	if err != nil {
		t.Fatalf("EditarItem() devolvió un error inesperado: %v", err)
	}

	esperado := original
	esperado.Prioridad = "Media"
	verificarItemBacklog(t, editado, esperado)
	verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacen), "p1", "HU-01").Item, esperado)
}

// CA1 / R1: editar los criterios de aceptación actualiza el ítem devuelto y el
// persistido.
func TestUS06_EditarItem_Criterios(t *testing.T) {
	original := itemUS06ServicioBase()
	almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p1", Item: original})
	servicio := almacen.servicio()

	criterios := []string{
		"El nombre del proyecto no puede estar vacío",
		"El proyecto queda disponible para asignar integrantes",
		"El proyecto se lista en el backlog",
	}
	editado, err := servicio.EditarItem(
		"p1",
		"HU-01",
		original.Titulo,
		original.Descripcion,
		original.Prioridad,
		criterios,
		original.StoryPoints,
	)
	if err != nil {
		t.Fatalf("EditarItem() devolvió un error inesperado: %v", err)
	}

	esperado := original
	esperado.Criterios = criterios
	verificarItemBacklog(t, editado, esperado)
	verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacen), "p1", "HU-01").Item, esperado)
}

// CA1 / R1 / R3: editar los story points admite los extremos de la escala.
func TestUS06_EditarItem_StoryPoints(t *testing.T) {
	for _, storyPoints := range []int{1, 8, 21} {
		t.Run(strconv.Itoa(storyPoints), func(t *testing.T) {
			original := itemUS06ServicioBase()
			almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p1", Item: original})
			servicio := almacen.servicio()

			editado, err := servicio.EditarItem(
				"p1",
				"HU-01",
				original.Titulo,
				original.Descripcion,
				original.Prioridad,
				original.Criterios,
				storyPoints,
			)
			if err != nil {
				t.Fatalf("EditarItem() devolvió un error inesperado: %v", err)
			}

			esperado := original
			esperado.StoryPoints = storyPoints
			verificarItemBacklog(t, editado, esperado)
			verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacen), "p1", "HU-01").Item, esperado)
		})
	}
}

// CA1 / R1 / R3: editar un ítem sin estimar con 0 conserva los story points en 0.
func TestUS06_EditarItem_ItemSinEstimar(t *testing.T) {
	original := itemUS06ServicioBase()
	original.StoryPoints = 0
	almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p1", Item: original})
	servicio := almacen.servicio()

	editado, err := servicio.EditarItem(
		"p1",
		"HU-01",
		"Registrar y listar proyectos",
		original.Descripcion,
		original.Prioridad,
		original.Criterios,
		0,
	)
	if err != nil {
		t.Fatalf("EditarItem() devolvió un error inesperado: %v", err)
	}

	esperado := original
	esperado.Titulo = "Registrar y listar proyectos"
	esperado.StoryPoints = 0
	verificarItemBacklog(t, editado, esperado)
	verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacen), "p1", "HU-01").Item, esperado)
}

// CA1 / R1 / R3: editar con 0 deja el ítem sin estimar.
func TestUS06_EditarItem_DejarSinEstimarConCero(t *testing.T) {
	original := itemUS06ServicioBase()
	almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p1", Item: original})
	servicio := almacen.servicio()

	editado, err := servicio.EditarItem(
		"p1",
		"HU-01",
		original.Titulo,
		original.Descripcion,
		original.Prioridad,
		original.Criterios,
		0,
	)
	if err != nil {
		t.Fatalf("EditarItem() devolvió un error inesperado: %v", err)
	}

	esperado := original
	esperado.StoryPoints = 0
	verificarItemBacklog(t, editado, esperado)
	verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacen), "p1", "HU-01").Item, esperado)
}

// CA1 / R1 / R2: editar todos los campos editables a la vez actualiza el ítem
// devuelto y el persistido, conservando ID y estado.
func TestUS06_EditarItem_Exitoso(t *testing.T) {
	original := itemUS06ServicioBase()
	almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p1", Item: original})
	servicio := almacen.servicio()

	editado, err := servicio.EditarItem(
		"p1",
		"HU-01",
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
		t.Fatalf("EditarItem() devolvió un error inesperado: %v", err)
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
	verificarItemBacklog(t, editado, esperado)
	verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacen), "p1", "HU-01").Item, esperado)
}

// CA1 / R2: el ID y el estado no cambian al editar y el ítem no se duplica.
func TestUS06_EditarItem_ConservaIDYEstado(t *testing.T) {
	original := itemUS06ServicioBase()
	almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p1", Item: original})
	servicio := almacen.servicio()

	editado, err := servicio.EditarItem(
		"p1",
		"HU-01",
		"Registrar y listar proyectos",
		"Permitir dar de alta y consultar proyectos en el sistema",
		"Media",
		[]string{"El proyecto se lista en el backlog"},
		8,
	)
	if err != nil {
		t.Fatalf("EditarItem() devolvió un error inesperado: %v", err)
	}

	if editado.ID != "HU-01" {
		t.Errorf("ID = %q, se esperaba %q", editado.ID, "HU-01")
	}
	if editado.Estado != "Pendiente" {
		t.Errorf("Estado = %q, se esperaba %q", editado.Estado, "Pendiente")
	}
	if len(almacen.fakeRepositorioBacklog.items["p1"]) != 1 {
		t.Errorf("el proyecto debería tener un solo ítem, tiene %d", len(almacen.fakeRepositorioBacklog.items["p1"]))
	}
	persistido := buscarRegistroItem(t, cargarBacklog(t, almacen), "p1", "HU-01").Item
	if persistido.ID != "HU-01" || persistido.Estado != "Pendiente" {
		t.Errorf("el ítem persistido debería conservar ID HU-01 y estado Pendiente, se obtuvo %+v", persistido)
	}
}

// CA1 / R3 / V4 / R12: un título vacío o solo de espacios devuelve
// ErrTituloItemVacio, no persiste y prevalece sobre prioridad y story points
// inválidos.
func TestUS06_EditarItem_TituloVacioNoPersiste(t *testing.T) {
	casos := []struct {
		descripcion string
		titulo      string
		prioridad   string
		storyPoints int
	}{
		{descripcion: "título vacío", titulo: "", prioridad: "Alta", storyPoints: 5},
		{descripcion: "título solo de espacios", titulo: "   ", prioridad: "Alta", storyPoints: 5},
		{descripcion: "título vacío prevalece sobre prioridad y story points inválidos", titulo: "", prioridad: "Urgente", storyPoints: 4},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			original := itemUS06ServicioBase()
			almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p1", Item: original})
			servicio := almacen.servicio()

			editado, err := servicio.EditarItem(
				"p1",
				"HU-01",
				caso.titulo,
				original.Descripcion,
				caso.prioridad,
				original.Criterios,
				caso.storyPoints,
			)
			if err == nil {
				t.Fatal("EditarItem() con título vacío debería devolver un error")
			}
			if !errors.Is(err, domain.ErrTituloItemVacio) {
				t.Errorf("err = %v, se esperaba %v", err, domain.ErrTituloItemVacio)
			}
			if errors.Is(err, domain.ErrPrioridadInvalida) ||
				errors.Is(err, domain.ErrStoryPointsInvalidos) {
				t.Errorf("err = %v, no debería devolverse un error de prioridad o story points", err)
			}
			verificarItemBacklogVacio(t, editado)
			if almacen.fakeRepositorioBacklog.actualizaciones != 0 {
				t.Errorf("ActualizarItem fue invocado %d veces; no debía escribirse nada", almacen.fakeRepositorioBacklog.actualizaciones)
			}
			verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacen), "p1", "HU-01").Item, original)
		})
	}
}

// CA1 / R3 / V5 / R12: una prioridad inválida devuelve ErrPrioridadInvalida, no
// persiste y prevalece sobre story points inválidos.
func TestUS06_EditarItem_PrioridadInvalidaNoPersiste(t *testing.T) {
	casos := []struct {
		descripcion string
		prioridad   string
		storyPoints int
	}{
		{descripcion: "prioridad no permitida", prioridad: "Urgente", storyPoints: 5},
		{descripcion: "prioridad en minúsculas", prioridad: "alta", storyPoints: 5},
		{descripcion: "prioridad en mayúsculas", prioridad: "ALTA", storyPoints: 5},
		{descripcion: "prioridad solo de espacios", prioridad: "   ", storyPoints: 5},
		{descripcion: "prioridad inválida prevalece sobre story points inválidos", prioridad: "Urgente", storyPoints: 4},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			original := itemUS06ServicioBase()
			almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p1", Item: original})
			servicio := almacen.servicio()

			editado, err := servicio.EditarItem(
				"p1",
				"HU-01",
				original.Titulo,
				original.Descripcion,
				caso.prioridad,
				original.Criterios,
				caso.storyPoints,
			)
			if err == nil {
				t.Fatal("EditarItem() con prioridad inválida debería devolver un error")
			}
			if !errors.Is(err, domain.ErrPrioridadInvalida) {
				t.Errorf("err = %v, se esperaba %v", err, domain.ErrPrioridadInvalida)
			}
			if errors.Is(err, domain.ErrStoryPointsInvalidos) {
				t.Errorf("err = %v, no debería devolverse un error de story points", err)
			}
			verificarItemBacklogVacio(t, editado)
			if almacen.fakeRepositorioBacklog.actualizaciones != 0 {
				t.Errorf("ActualizarItem fue invocado %d veces; no debía escribirse nada", almacen.fakeRepositorioBacklog.actualizaciones)
			}
			verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacen), "p1", "HU-01").Item, original)
		})
	}
}

// CA1 / R3 / V6: story points negativos, 4 y 22 devuelven
// ErrStoryPointsInvalidos y no persisten.
func TestUS06_EditarItem_StoryPointsInvalidosNoPersiste(t *testing.T) {
	for _, storyPoints := range []int{-1, 4, 22} {
		t.Run(strconv.Itoa(storyPoints), func(t *testing.T) {
			original := itemUS06ServicioBase()
			almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p1", Item: original})
			servicio := almacen.servicio()

			editado, err := servicio.EditarItem(
				"p1",
				"HU-01",
				original.Titulo,
				original.Descripcion,
				original.Prioridad,
				original.Criterios,
				storyPoints,
			)
			if err == nil {
				t.Fatal("EditarItem() con story points inválidos debería devolver un error")
			}
			if !errors.Is(err, domain.ErrStoryPointsInvalidos) {
				t.Errorf("err = %v, se esperaba %v", err, domain.ErrStoryPointsInvalidos)
			}
			verificarItemBacklogVacio(t, editado)
			if almacen.fakeRepositorioBacklog.actualizaciones != 0 {
				t.Errorf("ActualizarItem fue invocado %d veces; no debía escribirse nada", almacen.fakeRepositorioBacklog.actualizaciones)
			}
			verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacen), "p1", "HU-01").Item, original)
		})
	}
}

// CA1 / R7 / R12: un proyecto inexistente devuelve ErrProyectoNoEncontrado
// antes que cualquier validación de campos y no escribe nada.
func TestUS06_EditarItem_ProyectoNoEncontrado(t *testing.T) {
	original := itemUS06ServicioBase()
	almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p1", Item: original})
	servicio := almacen.servicio()

	editado, err := servicio.EditarItem(
		"ausente",
		"HU-01",
		"",
		original.Descripcion,
		"Urgente",
		original.Criterios,
		4,
	)
	if err == nil {
		t.Fatal("EditarItem() con un proyecto inexistente debería devolver un error")
	}
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if errors.Is(err, domain.ErrItemNoEncontrado) ||
		errors.Is(err, domain.ErrTituloItemVacio) ||
		errors.Is(err, domain.ErrPrioridadInvalida) ||
		errors.Is(err, domain.ErrStoryPointsInvalidos) {
		t.Errorf("err = %v, no debería devolverse otro error de dominio", err)
	}
	verificarItemBacklogVacio(t, editado)
	if almacen.fakeRepositorioBacklog.actualizaciones != 0 {
		t.Errorf("ActualizarItem fue invocado %d veces; no debía escribirse nada", almacen.fakeRepositorioBacklog.actualizaciones)
	}
	verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacen), "p1", "HU-01").Item, original)
}

// CA1 / R8 / R12: un ítem inexistente devuelve ErrItemNoEncontrado antes que
// las validaciones de campos y no escribe nada. El mismo ID en otro proyecto no
// cuenta como existente.
func TestUS06_EditarItem_ItemNoEncontrado(t *testing.T) {
	t.Run("ítem ausente en el proyecto", func(t *testing.T) {
		original := itemUS06ServicioBase()
		almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p1", Item: original})
		servicio := almacen.servicio()

		editado, err := servicio.EditarItem(
			"p1",
			"HU-02",
			"",
			original.Descripcion,
			"Urgente",
			original.Criterios,
			4,
		)
		if err == nil {
			t.Fatal("EditarItem() con un ítem inexistente debería devolver un error")
		}
		if !errors.Is(err, domain.ErrItemNoEncontrado) {
			t.Errorf("err = %v, se esperaba %v", err, domain.ErrItemNoEncontrado)
		}
		if errors.Is(err, domain.ErrProyectoNoEncontrado) ||
			errors.Is(err, domain.ErrTituloItemVacio) ||
			errors.Is(err, domain.ErrPrioridadInvalida) ||
			errors.Is(err, domain.ErrStoryPointsInvalidos) {
			t.Errorf("err = %v, no debería devolverse otro error de dominio", err)
		}
		verificarItemBacklogVacio(t, editado)
		if almacen.fakeRepositorioBacklog.actualizaciones != 0 {
			t.Errorf("ActualizarItem fue invocado %d veces; no debía escribirse nada", almacen.fakeRepositorioBacklog.actualizaciones)
		}
		verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacen), "p1", "HU-01").Item, original)
	})

	t.Run("mismo ID en otro proyecto no cuenta", func(t *testing.T) {
		itemOtroProyecto := itemUS06ServicioBase()
		almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p2", Item: itemOtroProyecto})
		servicio := almacen.servicio()

		editado, err := servicio.EditarItem(
			"p1",
			"HU-01",
			"Registrar y listar proyectos",
			itemOtroProyecto.Descripcion,
			itemOtroProyecto.Prioridad,
			itemOtroProyecto.Criterios,
			itemOtroProyecto.StoryPoints,
		)
		if err == nil {
			t.Fatal("EditarItem() con un ítem de otro proyecto debería devolver un error")
		}
		if !errors.Is(err, domain.ErrItemNoEncontrado) {
			t.Errorf("err = %v, se esperaba %v", err, domain.ErrItemNoEncontrado)
		}
		verificarItemBacklogVacio(t, editado)
		if almacen.fakeRepositorioBacklog.actualizaciones != 0 {
			t.Errorf("ActualizarItem fue invocado %d veces; no debía escribirse nada", almacen.fakeRepositorioBacklog.actualizaciones)
		}
		verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacen), "p2", "HU-01").Item, itemOtroProyecto)
	})
}

// CA3 / R4: un ítem completado no se puede editar y no se persiste.
func TestUS06_EditarItem_ItemCompletadoNoPersiste(t *testing.T) {
	original := itemUS06ServicioBase()
	original.Estado = domain.EstadoCompletada
	almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p1", Item: original})
	servicio := almacen.servicio()

	editado, err := servicio.EditarItem(
		"p1",
		"HU-01",
		"Registrar y listar proyectos",
		"Permitir dar de alta y consultar proyectos en el sistema",
		"Media",
		[]string{"El proyecto se lista en el backlog"},
		8,
	)
	if err == nil {
		t.Fatal("EditarItem() sobre un ítem completado debería devolver un error")
	}
	if !errors.Is(err, domain.ErrItemCompletadoNoEditable) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrItemCompletadoNoEditable)
	}
	verificarItemBacklogVacio(t, editado)
	if almacen.fakeRepositorioBacklog.actualizaciones != 0 {
		t.Errorf("ActualizarItem fue invocado %d veces; no debía escribirse nada", almacen.fakeRepositorioBacklog.actualizaciones)
	}
	verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacen), "p1", "HU-01").Item, original)
}

// CA3 / R4 / R12: el error de ítem completado prevalece sobre las validaciones
// de los campos editables.
func TestUS06_EditarItem_ItemCompletadoAntesQueValidaciones(t *testing.T) {
	original := itemUS06ServicioBase()
	original.Estado = domain.EstadoCompletada
	almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p1", Item: original})
	servicio := almacen.servicio()

	editado, err := servicio.EditarItem(
		"p1",
		"HU-01",
		"",
		original.Descripcion,
		"Urgente",
		original.Criterios,
		4,
	)
	if err == nil {
		t.Fatal("EditarItem() sobre un ítem completado debería devolver un error")
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
	if almacen.fakeRepositorioBacklog.actualizaciones != 0 {
		t.Errorf("ActualizarItem fue invocado %d veces; no debía escribirse nada", almacen.fakeRepositorioBacklog.actualizaciones)
	}
	verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacen), "p1", "HU-01").Item, original)
}

// CA2 / R5: el backlog se lista ordenado por prioridad Alta, Media, Baja.
func TestUS06_ListarBacklog_OrdenadoPorPrioridad(t *testing.T) {
	almacen := sembrarEntornoUS06(t,
		itemConProyecto{ProyectoID: "p1", Item: itemUS06ServicioLista("HU-03", "Baja")},
		itemConProyecto{ProyectoID: "p1", Item: itemUS06ServicioLista("HU-01", "Alta")},
		itemConProyecto{ProyectoID: "p1", Item: itemUS06ServicioLista("HU-02", "Media")},
	)
	servicio := almacen.servicio()

	obtenido, err := servicio.ListarBacklog("p1")
	if err != nil {
		t.Fatalf("ListarBacklog() devolvió un error inesperado: %v", err)
	}
	verificarIDsUS06(t, obtenido, []string{"HU-01", "HU-02", "HU-03"})

	prioridadesEsperadas := []string{"Alta", "Media", "Baja"}
	for i, prioridadEsperada := range prioridadesEsperadas {
		if obtenido[i].Prioridad != prioridadEsperada {
			t.Errorf("Prioridad en posición %d = %q, se esperaba %q", i, obtenido[i].Prioridad, prioridadEsperada)
		}
	}
}

// CA2 / R6: a igual prioridad, el desempate es por número ascendente.
func TestUS06_ListarBacklog_DesempatePorNumero(t *testing.T) {
	almacen := sembrarEntornoUS06(t,
		itemConProyecto{ProyectoID: "p1", Item: itemUS06ServicioLista("HU-03", "Media")},
		itemConProyecto{ProyectoID: "p1", Item: itemUS06ServicioLista("HU-01", "Media")},
		itemConProyecto{ProyectoID: "p1", Item: itemUS06ServicioLista("HU-02", "Media")},
	)
	servicio := almacen.servicio()

	obtenido, err := servicio.ListarBacklog("p1")
	if err != nil {
		t.Fatalf("ListarBacklog() devolvió un error inesperado: %v", err)
	}
	verificarIDsUS06(t, obtenido, []string{"HU-01", "HU-02", "HU-03"})
}

// CA2 / R6: el desempate numérico ubica HU-2 antes que HU-10.
func TestUS06_ListarBacklog_HU10VsHU2(t *testing.T) {
	almacen := sembrarEntornoUS06(t,
		itemConProyecto{ProyectoID: "p1", Item: itemUS06ServicioLista("HU-10", "Alta")},
		itemConProyecto{ProyectoID: "p1", Item: itemUS06ServicioLista("HU-2", "Alta")},
		itemConProyecto{ProyectoID: "p1", Item: itemUS06ServicioLista("HU-1", "Alta")},
	)
	servicio := almacen.servicio()

	obtenido, err := servicio.ListarBacklog("p1")
	if err != nil {
		t.Fatalf("ListarBacklog() devolvió un error inesperado: %v", err)
	}
	verificarIDsUS06(t, obtenido, []string{"HU-1", "HU-2", "HU-10"})
}

// CA2 / R11: un proyecto sin ítems devuelve una lista vacía sin error.
func TestUS06_ListarBacklog_BacklogVacio(t *testing.T) {
	almacen := sembrarEntornoUS06(t)
	servicio := almacen.servicio()

	obtenido, err := servicio.ListarBacklog("p1")
	if err != nil {
		t.Fatalf("ListarBacklog() devolvió un error inesperado: %v", err)
	}
	if len(obtenido) != 0 {
		t.Errorf("len(ListarBacklog()) = %d, se esperaba 0", len(obtenido))
	}
}

// CA2 / R7: un proyecto inexistente al listar devuelve
// ErrProyectoNoEncontrado.
func TestUS06_ListarBacklog_ProyectoNoEncontrado(t *testing.T) {
	almacen := sembrarEntornoUS06(t)
	servicio := almacen.servicio()

	obtenido, err := servicio.ListarBacklog("ausente")
	if err == nil {
		t.Fatal("ListarBacklog() con un proyecto inexistente debería devolver un error")
	}
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if len(obtenido) != 0 {
		t.Errorf("len(ListarBacklog()) = %d, se esperaba 0", len(obtenido))
	}
}

// R10 / E7: los errores de repositorio de la edición se envuelven con el error
// original y no se mapean a errores de no encontrado.
func TestUS06_EditarItem_ErrorDeRepositorioEnvuelto(t *testing.T) {
	casos := []struct {
		descripcion string
		preparar    func(almacen *entornoBacklog, err error)
	}{
		{
			descripcion: "error al verificar el proyecto",
			preparar: func(almacen *entornoBacklog, err error) {
				almacen.fakeRepositorioProyectos.existeErr = err
			},
		},
		{
			descripcion: "error al obtener el ítem",
			preparar: func(almacen *entornoBacklog, err error) {
				almacen.fakeRepositorioBacklog.obtenerErr = err
			},
		},
		{
			descripcion: "error al actualizar el ítem",
			preparar: func(almacen *entornoBacklog, err error) {
				almacen.fakeRepositorioBacklog.actualizarErr = err
			},
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			errSimulado := errors.New("fallo simulado del repositorio")
			original := itemUS06ServicioBase()
			almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p1", Item: original})
			caso.preparar(almacen, errSimulado)
			servicio := almacen.servicio()

			editado, err := servicio.EditarItem(
				"p1",
				"HU-01",
				"Registrar y listar proyectos",
				original.Descripcion,
				original.Prioridad,
				original.Criterios,
				original.StoryPoints,
			)
			if err == nil {
				t.Fatal("EditarItem() con un error de repositorio debería devolver un error")
			}
			if !errors.Is(err, errSimulado) {
				t.Errorf("err = %v, se esperaba un error que envuelva %v", err, errSimulado)
			}
			if errors.Is(err, domain.ErrProyectoNoEncontrado) {
				t.Errorf("err = %v, no debería mapearse a ErrProyectoNoEncontrado", err)
			}
			if errors.Is(err, domain.ErrItemNoEncontrado) {
				t.Errorf("err = %v, no debería mapearse a ErrItemNoEncontrado", err)
			}
			verificarItemBacklogVacio(t, editado)
		})
	}
}

// R10 / E7: los errores de repositorio del listado se envuelven con el error
// original y no se mapean a errores de no encontrado.
func TestUS06_ListarBacklog_ErrorDeRepositorioEnvuelto(t *testing.T) {
	casos := []struct {
		descripcion string
		preparar    func(almacen *entornoBacklog, err error)
	}{
		{
			descripcion: "error al verificar el proyecto",
			preparar: func(almacen *entornoBacklog, err error) {
				almacen.fakeRepositorioProyectos.existeErr = err
			},
		},
		{
			descripcion: "error al listar el backlog",
			preparar: func(almacen *entornoBacklog, err error) {
				almacen.fakeRepositorioBacklog.listarErr = err
			},
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			errSimulado := errors.New("fallo simulado del repositorio")
			almacen := sembrarEntornoUS06(t,
				itemConProyecto{ProyectoID: "p1", Item: itemUS06ServicioLista("HU-01", "Alta")},
			)
			caso.preparar(almacen, errSimulado)
			servicio := almacen.servicio()

			obtenido, err := servicio.ListarBacklog("p1")
			if err == nil {
				t.Fatal("ListarBacklog() con un error de repositorio debería devolver un error")
			}
			if !errors.Is(err, errSimulado) {
				t.Errorf("err = %v, se esperaba un error que envuelva %v", err, errSimulado)
			}
			if errors.Is(err, domain.ErrProyectoNoEncontrado) {
				t.Errorf("err = %v, no debería mapearse a ErrProyectoNoEncontrado", err)
			}
			if errors.Is(err, domain.ErrItemNoEncontrado) {
				t.Errorf("err = %v, no debería mapearse a ErrItemNoEncontrado", err)
			}
			if len(obtenido) != 0 {
				t.Errorf("len(ListarBacklog()) = %d, se esperaba 0", len(obtenido))
			}
		})
	}
}

// R9: editar un ítem entre varios conserva los demás ítems y su orden.
func TestUS06_EditarItem_ConservaOtrosItemsYOrden(t *testing.T) {
	hu01 := itemUS06ServicioLista("HU-01", "Alta")
	hu02 := itemUS06ServicioLista("HU-02", "Media")
	hu03 := itemUS06ServicioLista("HU-03", "Baja")
	almacen := sembrarEntornoUS06(t,
		itemConProyecto{ProyectoID: "p1", Item: hu01},
		itemConProyecto{ProyectoID: "p1", Item: hu02},
		itemConProyecto{ProyectoID: "p1", Item: hu03},
	)
	servicio := almacen.servicio()

	editado, err := servicio.EditarItem(
		"p1",
		"HU-02",
		"Modificar proyecto",
		"Permitir modificar proyectos existentes",
		"Alta",
		[]string{"El proyecto modificado conserva su ID"},
		8,
	)
	if err != nil {
		t.Fatalf("EditarItem() devolvió un error inesperado: %v", err)
	}

	esperadoHu02 := hu02
	esperadoHu02.Titulo = "Modificar proyecto"
	esperadoHu02.Descripcion = "Permitir modificar proyectos existentes"
	esperadoHu02.Prioridad = "Alta"
	esperadoHu02.Criterios = []string{"El proyecto modificado conserva su ID"}
	esperadoHu02.StoryPoints = 8
	verificarItemBacklog(t, editado, esperadoHu02)

	if almacen.fakeRepositorioBacklog.ultimoItemID != "HU-02" {
		t.Errorf("ActualizarItem recibió el ítem %q, se esperaba %q", almacen.fakeRepositorioBacklog.ultimoItemID, "HU-02")
	}

	coleccion := cargarBacklog(t, almacen)
	if len(coleccion) != 3 {
		t.Fatalf("se esperaban 3 ítems persistidos, hay %d", len(coleccion))
	}
	verificarIDsUS06(t, []domain.ItemBacklog{coleccion[0].Item, coleccion[1].Item, coleccion[2].Item}, []string{"HU-01", "HU-02", "HU-03"})
	verificarItemBacklog(t, coleccion[0].Item, hu01)
	verificarItemBacklog(t, coleccion[1].Item, esperadoHu02)
	verificarItemBacklog(t, coleccion[2].Item, hu03)
}

// R7 / R8 / R9 / R10 / R12: ningún camino de error escribe ni altera lo
// persistido.
func TestUS06_EditarItem_ErrorNoEscribe(t *testing.T) {
	itemCompletado := itemUS06ServicioBase()
	itemCompletado.Estado = domain.EstadoCompletada

	casos := []struct {
		descripcion   string
		proyectoID    string
		itemID        string
		titulo        string
		prioridad     string
		storyPoints   int
		itemInicial   domain.ItemBacklog
		errorEsperado error
	}{
		{
			descripcion:   "proyecto inexistente",
			proyectoID:    "ausente",
			itemID:        "HU-01",
			titulo:        "Registrar y listar proyectos",
			prioridad:     "Media",
			storyPoints:   8,
			itemInicial:   itemUS06ServicioBase(),
			errorEsperado: domain.ErrProyectoNoEncontrado,
		},
		{
			descripcion:   "ítem inexistente",
			proyectoID:    "p1",
			itemID:        "HU-02",
			titulo:        "Registrar y listar proyectos",
			prioridad:     "Media",
			storyPoints:   8,
			itemInicial:   itemUS06ServicioBase(),
			errorEsperado: domain.ErrItemNoEncontrado,
		},
		{
			descripcion:   "título vacío",
			proyectoID:    "p1",
			itemID:        "HU-01",
			titulo:        "",
			prioridad:     "Alta",
			storyPoints:   5,
			itemInicial:   itemUS06ServicioBase(),
			errorEsperado: domain.ErrTituloItemVacio,
		},
		{
			descripcion:   "prioridad inválida",
			proyectoID:    "p1",
			itemID:        "HU-01",
			titulo:        "Registrar proyecto",
			prioridad:     "Urgente",
			storyPoints:   5,
			itemInicial:   itemUS06ServicioBase(),
			errorEsperado: domain.ErrPrioridadInvalida,
		},
		{
			descripcion:   "story points inválidos",
			proyectoID:    "p1",
			itemID:        "HU-01",
			titulo:        "Registrar proyecto",
			prioridad:     "Alta",
			storyPoints:   4,
			itemInicial:   itemUS06ServicioBase(),
			errorEsperado: domain.ErrStoryPointsInvalidos,
		},
		{
			descripcion:   "ítem completado",
			proyectoID:    "p1",
			itemID:        "HU-01",
			titulo:        "Registrar y listar proyectos",
			prioridad:     "Media",
			storyPoints:   8,
			itemInicial:   itemCompletado,
			errorEsperado: domain.ErrItemCompletadoNoEditable,
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			almacen := sembrarEntornoUS06(t, itemConProyecto{ProyectoID: "p1", Item: caso.itemInicial})
			servicio := almacen.servicio()

			editado, err := servicio.EditarItem(
				caso.proyectoID,
				caso.itemID,
				caso.titulo,
				caso.itemInicial.Descripcion,
				caso.prioridad,
				caso.itemInicial.Criterios,
				caso.storyPoints,
			)
			if err == nil {
				t.Fatal("EditarItem() con datos inválidos debería devolver un error")
			}
			if !errors.Is(err, caso.errorEsperado) {
				t.Errorf("err = %v, se esperaba %v", err, caso.errorEsperado)
			}
			verificarItemBacklogVacio(t, editado)
			if almacen.fakeRepositorioBacklog.actualizaciones != 0 {
				t.Errorf("ActualizarItem fue invocado %d veces; no debía escribirse nada", almacen.fakeRepositorioBacklog.actualizaciones)
			}
			verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacen), "p1", "HU-01").Item, caso.itemInicial)
		})
	}
}

// R9: el listado de un proyecto no mezcla ítems de otros proyectos.
func TestUS06_ListarBacklog_NoMezclaProyectos(t *testing.T) {
	almacen := sembrarEntornoUS06(t,
		itemConProyecto{ProyectoID: "p1", Item: itemUS06ServicioLista("HU-01", "Alta")},
		itemConProyecto{ProyectoID: "p1", Item: itemUS06ServicioLista("HU-02", "Media")},
		itemConProyecto{ProyectoID: "p2", Item: itemUS06ServicioLista("HU-03", "Baja")},
		itemConProyecto{ProyectoID: "p2", Item: itemUS06ServicioLista("HU-04", "Alta")},
	)
	servicio := almacen.servicio()

	obtenido, err := servicio.ListarBacklog("p1")
	if err != nil {
		t.Fatalf("ListarBacklog() devolvió un error inesperado: %v", err)
	}
	verificarIDsUS06(t, obtenido, []string{"HU-01", "HU-02"})
	for _, item := range obtenido {
		if item.ID == "HU-03" || item.ID == "HU-04" {
			t.Errorf("el listado del proyecto p1 no debería incluir el ítem %q", item.ID)
		}
	}
}
