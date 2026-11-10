package service

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage"
)

func itemUS11Backlog(storyPoints int) domain.ItemBacklog {
	return domain.ItemBacklog{
		ID:          "HU-01",
		Titulo:      "Registrar proyecto",
		Descripcion: "Permitir dar de alta proyectos en el sistema",
		Prioridad:   "Alta",
		Estado:      "Pendiente",
		StoryPoints: storyPoints,
		Criterios: []string{
			"El nombre del proyecto no puede estar vacío",
			"El proyecto queda disponible para asignar integrantes",
		},
	}
}

// CA1 / R1: estimar un ítem existente con 5 actualiza el valor devuelto y
// persistido.
func TestUS11_EstimarItem_Exitoso(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
	})
	original := itemUS11Backlog(0)
	sembrarBacklog(t, almacen, []itemRegistrado{{ProyectoID: "p1", Item: original}})
	servicio := NuevoServicioBacklog(almacen)

	estimado, err := servicio.EstimarItem("p1", "HU-01", 5)
	if err != nil {
		t.Fatalf("EstimarItem() devolvió un error inesperado: %v", err)
	}
	esperado := original
	esperado.StoryPoints = 5
	verificarItemBacklog(t, estimado, esperado)
	verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacen), "p1", "HU-01").Item, esperado)
}

// CA1 / R1: la estimación conserva ID, título, descripción, prioridad, estado y
// criterios del ítem.
func TestUS11_EstimarItem_ConservaRestoDeCampos(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
	})
	original := domain.ItemBacklog{
		ID:          "HU-07",
		Titulo:      "Consultar métricas",
		Descripcion: "Consultar las métricas del proyecto",
		Prioridad:   "Media",
		Estado:      "En curso",
		StoryPoints: 0,
		Criterios: []string{
			"Se listan las métricas",
			"Cada métrica tiene una descripción",
		},
	}
	sembrarBacklog(t, almacen, []itemRegistrado{{ProyectoID: "p1", Item: original}})
	servicio := NuevoServicioBacklog(almacen)

	estimado, err := servicio.EstimarItem("p1", "HU-07", 8)
	if err != nil {
		t.Fatalf("EstimarItem() devolvió un error inesperado: %v", err)
	}
	esperado := original
	esperado.StoryPoints = 8
	verificarItemBacklog(t, estimado, esperado)
	verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacen), "p1", "HU-07").Item, esperado)
}

// CA1 / R4: re-estimar de 3 a 8 reemplaza el valor y no duplica el registro.
func TestUS11_EstimarItem_Reestimacion(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
	})
	original := itemUS11Backlog(3)
	sembrarBacklog(t, almacen, []itemRegistrado{{ProyectoID: "p1", Item: original}})
	servicio := NuevoServicioBacklog(almacen)

	estimado, err := servicio.EstimarItem("p1", "HU-01", 8)
	if err != nil {
		t.Fatalf("EstimarItem() devolvió un error inesperado: %v", err)
	}
	esperado := original
	esperado.StoryPoints = 8
	verificarItemBacklog(t, estimado, esperado)

	coleccion := cargarBacklog(t, almacen)
	if len(coleccion) != 1 {
		t.Fatalf("se esperaba 1 registro, hay %d", len(coleccion))
	}
	verificarItemBacklog(t, buscarRegistroItem(t, coleccion, "p1", "HU-01").Item, esperado)
}

// CA1 / R4: un ítem sin estimar (0) pasa a estar estimado con 5.
func TestUS11_EstimarItem_ItemSinEstimar(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
	})
	original := itemUS11Backlog(0)
	sembrarBacklog(t, almacen, []itemRegistrado{{ProyectoID: "p1", Item: original}})
	servicio := NuevoServicioBacklog(almacen)

	estimado, err := servicio.EstimarItem("p1", "HU-01", 5)
	if err != nil {
		t.Fatalf("EstimarItem() devolvió un error inesperado: %v", err)
	}
	esperado := original
	esperado.StoryPoints = 5
	verificarItemBacklog(t, estimado, esperado)
	verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacen), "p1", "HU-01").Item, esperado)
}

// CA1 / CA2 / R5: un proyecto ausente devuelve ErrProyectoNoEncontrado y no
// escribe nada.
func TestUS11_EstimarItem_ProyectoNoEncontrado(t *testing.T) {
	almacenReal, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacenReal, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
		{ID: "p2", Proyecto: proyectoUS02Servicio("App Móvil", domain.EstadoPlanificado)},
	})
	original := itemUS11Backlog(2)
	sembrarBacklog(t, almacenReal, []itemRegistrado{{ProyectoID: "p1", Item: original}})
	almacen := &almacenCuentaSaves{real: almacenReal}
	servicio := NuevoServicioBacklog(almacen)

	estimado, err := servicio.EstimarItem("ausente", "HU-01", 5)
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	verificarItemBacklogVacio(t, estimado)
	if almacen.llamadasASave != 0 {
		t.Errorf("Save fue invocado %d veces; no debía escribirse nada", almacen.llamadasASave)
	}
	verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacenReal), "p1", "HU-01").Item, original)
}

// CA1 / CA2 / R5 / D4: sin archivo de proyectos se interpreta como colección
// vacía, devuelve ErrProyectoNoEncontrado y no crea archivos.
func TestUS11_EstimarItem_ArchivoProyectosInexistente(t *testing.T) {
	almacen, dir := nuevoAlmacenServicio(t)
	servicio := NuevoServicioBacklog(almacen)

	estimado, err := servicio.EstimarItem("p1", "HU-01", 5)
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	verificarItemBacklogVacio(t, estimado)

	entradas, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir(%q): %v", dir, err)
	}
	if len(entradas) != 0 {
		t.Errorf("en ningún camino de error se crea el archivo, hay %d entradas", len(entradas))
	}
}

// CA1 / CA2 / R6: un ítem ausente en un proyecto con otros ítems devuelve
// ErrItemNoEncontrado y no escribe nada.
func TestUS11_EstimarItem_ItemNoEncontrado(t *testing.T) {
	almacenReal, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacenReal, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
	})
	original := itemUS11Backlog(2)
	sembrarBacklog(t, almacenReal, []itemRegistrado{{ProyectoID: "p1", Item: original}})
	almacen := &almacenCuentaSaves{real: almacenReal}
	servicio := NuevoServicioBacklog(almacen)

	estimado, err := servicio.EstimarItem("p1", "HU-02", 5)
	if !errors.Is(err, domain.ErrItemNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrItemNoEncontrado)
	}
	verificarItemBacklogVacio(t, estimado)
	if almacen.llamadasASave != 0 {
		t.Errorf("Save fue invocado %d veces; no debía escribirse nada", almacen.llamadasASave)
	}
	verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacenReal), "p1", "HU-01").Item, original)
}

// CA1 / CA2 / R6 / D4: sin backlog.json se interpreta como backlog vacío,
// devuelve ErrItemNoEncontrado y no crea el archivo.
func TestUS11_EstimarItem_BacklogSinArchivo(t *testing.T) {
	almacenReal, dir := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacenReal, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
	})
	almacen := &almacenCuentaSaves{real: almacenReal}
	servicio := NuevoServicioBacklog(almacen)

	estimado, err := servicio.EstimarItem("p1", "HU-01", 5)
	if !errors.Is(err, domain.ErrItemNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrItemNoEncontrado)
	}
	verificarItemBacklogVacio(t, estimado)
	if almacen.llamadasASave != 0 {
		t.Errorf("Save fue invocado %d veces; no debía escribirse nada", almacen.llamadasASave)
	}
	if _, err := os.Stat(filepath.Join(dir, "backlog.json")); !os.IsNotExist(err) {
		t.Errorf("backlog.json no debería crearse ante un ítem inexistente: err = %v", err)
	}
}

// CA1 / CA2 / R6 / R9: el mismo ID en otro proyecto no se encuentra y el ítem
// del otro proyecto no se modifica.
func TestUS11_EstimarItem_MismoIDEnOtroProyecto(t *testing.T) {
	almacenReal, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacenReal, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
		{ID: "p2", Proyecto: proyectoUS02Servicio("App Móvil", domain.EstadoPlanificado)},
	})
	original := itemUS11Backlog(2)
	sembrarBacklog(t, almacenReal, []itemRegistrado{{ProyectoID: "p2", Item: original}})
	almacen := &almacenCuentaSaves{real: almacenReal}
	servicio := NuevoServicioBacklog(almacen)

	estimado, err := servicio.EstimarItem("p1", "HU-01", 8)
	if !errors.Is(err, domain.ErrItemNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrItemNoEncontrado)
	}
	verificarItemBacklogVacio(t, estimado)
	if almacen.llamadasASave != 0 {
		t.Errorf("Save fue invocado %d veces; no debía escribirse nada", almacen.llamadasASave)
	}
	verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacenReal), "p2", "HU-01").Item, original)
}

// CA2 / R2 / R7: los valores 0, -1, 4 y 22 devuelven ErrStoryPointsInvalidos y
// conservan el contenido persistido.
func TestUS11_EstimarItem_StoryPointsInvalidosNoPersiste(t *testing.T) {
	for _, storyPoints := range []int{0, -1, 4, 22} {
		t.Run(strconv.Itoa(storyPoints), func(t *testing.T) {
			almacenReal, _ := nuevoAlmacenServicio(t)
			sembrarProyectos(t, almacenReal, []proyectoRegistrado{
				{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
			})
			original := itemUS11Backlog(3)
			sembrarBacklog(t, almacenReal, []itemRegistrado{{ProyectoID: "p1", Item: original}})
			almacen := &almacenCuentaSaves{real: almacenReal}
			servicio := NuevoServicioBacklog(almacen)

			estimado, err := servicio.EstimarItem("p1", "HU-01", storyPoints)
			if !errors.Is(err, domain.ErrStoryPointsInvalidos) {
				t.Errorf("err = %v, se esperaba %v", err, domain.ErrStoryPointsInvalidos)
			}
			verificarItemBacklogVacio(t, estimado)
			if almacen.llamadasASave != 0 {
				t.Errorf("Save fue invocado %d veces; no debía escribirse nada", almacen.llamadasASave)
			}
			verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacenReal), "p1", "HU-01").Item, original)
		})
	}
}

// CA2 / R7: un valor inválido sobre un ítem ya estimado conserva el valor
// anterior.
func TestUS11_EstimarItem_ValorInvalidoConservaValorPrevio(t *testing.T) {
	almacenReal, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacenReal, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
	})
	original := itemUS11Backlog(3)
	sembrarBacklog(t, almacenReal, []itemRegistrado{{ProyectoID: "p1", Item: original}})
	almacen := &almacenCuentaSaves{real: almacenReal}
	servicio := NuevoServicioBacklog(almacen)

	estimado, err := servicio.EstimarItem("p1", "HU-01", 4)
	if !errors.Is(err, domain.ErrStoryPointsInvalidos) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrStoryPointsInvalidos)
	}
	verificarItemBacklogVacio(t, estimado)
	if almacen.llamadasASave != 0 {
		t.Errorf("Save fue invocado %d veces; no debía escribirse nada", almacen.llamadasASave)
	}
	verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacenReal), "p1", "HU-01").Item, original)
}

// CA1 / CA2 / R5 / R6 / R7: prevalece el error de proyecto sobre el valor
// inválido.
func TestUS11_EstimarItem_ProyectoInexistenteAntesQueValorInvalido(t *testing.T) {
	almacenReal, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacenReal, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
	})
	almacen := &almacenCuentaSaves{real: almacenReal}
	servicio := NuevoServicioBacklog(almacen)

	estimado, err := servicio.EstimarItem("ausente", "HU-01", 4)
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if errors.Is(err, domain.ErrStoryPointsInvalidos) {
		t.Errorf("err = %v, no debería devolverse ErrStoryPointsInvalidos si el proyecto no existe", err)
	}
	verificarItemBacklogVacio(t, estimado)
	if almacen.llamadasASave != 0 {
		t.Errorf("Save fue invocado %d veces; no debía escribirse nada", almacen.llamadasASave)
	}
}

// CA1 / CA2 / R6 / R7: prevalece el error de ítem sobre el valor inválido.
func TestUS11_EstimarItem_ItemInexistenteAntesQueValorInvalido(t *testing.T) {
	almacenReal, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacenReal, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
	})
	original := itemUS11Backlog(2)
	sembrarBacklog(t, almacenReal, []itemRegistrado{{ProyectoID: "p1", Item: original}})
	almacen := &almacenCuentaSaves{real: almacenReal}
	servicio := NuevoServicioBacklog(almacen)

	estimado, err := servicio.EstimarItem("p1", "HU-02", 4)
	if !errors.Is(err, domain.ErrItemNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrItemNoEncontrado)
	}
	if errors.Is(err, domain.ErrStoryPointsInvalidos) {
		t.Errorf("err = %v, no debería devolverse ErrStoryPointsInvalidos si el ítem no existe", err)
	}
	verificarItemBacklogVacio(t, estimado)
	if almacen.llamadasASave != 0 {
		t.Errorf("Save fue invocado %d veces; no debía escribirse nada", almacen.llamadasASave)
	}
	verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacenReal), "p1", "HU-01").Item, original)
}

// CA1 / R8 / R9: solo cambia el ítem objetivo y se conservan los demás ítems
// con su orden.
func TestUS11_EstimarItem_ConservaOtrosItemsYOrden(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
	})
	hu01 := itemUS11Backlog(2)
	hu01.ID = "HU-01"
	hu01.Titulo = "Registrar proyecto"
	hu02 := itemUS11Backlog(3)
	hu02.ID = "HU-02"
	hu02.Titulo = "Modificar proyecto"
	hu03 := itemUS11Backlog(5)
	hu03.ID = "HU-03"
	hu03.Titulo = "Consultar métricas"
	sembrarBacklog(t, almacen, []itemRegistrado{
		{ProyectoID: "p1", Item: hu01},
		{ProyectoID: "p1", Item: hu02},
		{ProyectoID: "p1", Item: hu03},
	})
	servicio := NuevoServicioBacklog(almacen)

	estimado, err := servicio.EstimarItem("p1", "HU-02", 13)
	if err != nil {
		t.Fatalf("EstimarItem() devolvió un error inesperado: %v", err)
	}
	esperadoHu02 := hu02
	esperadoHu02.StoryPoints = 13
	verificarItemBacklog(t, estimado, esperadoHu02)

	coleccion := cargarBacklog(t, almacen)
	if len(coleccion) != 3 {
		t.Fatalf("se esperaban 3 registros, hay %d", len(coleccion))
	}
	idsEsperados := []string{"HU-01", "HU-02", "HU-03"}
	for i, idEsperado := range idsEsperados {
		if coleccion[i].Item.ID != idEsperado {
			t.Errorf("ID en posición %d = %q, se esperaba %q", i, coleccion[i].Item.ID, idEsperado)
		}
	}
	verificarItemBacklog(t, coleccion[0].Item, hu01)
	verificarItemBacklog(t, coleccion[1].Item, esperadoHu02)
	verificarItemBacklog(t, coleccion[2].Item, hu03)
}

// CA1 / R9: estimar en un proyecto no afecta al mismo ID de otro proyecto.
func TestUS11_EstimarItem_NoAfectaItemsDeOtrosProyectos(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
		{ID: "p2", Proyecto: proyectoUS02Servicio("App Móvil", domain.EstadoPlanificado)},
	})
	p1Item := itemUS11Backlog(0)
	p2Item := itemUS11Backlog(2)
	sembrarBacklog(t, almacen, []itemRegistrado{
		{ProyectoID: "p1", Item: p1Item},
		{ProyectoID: "p2", Item: p2Item},
	})
	servicio := NuevoServicioBacklog(almacen)

	estimado, err := servicio.EstimarItem("p1", "HU-01", 8)
	if err != nil {
		t.Fatalf("EstimarItem() devolvió un error inesperado: %v", err)
	}
	esperadoP1 := p1Item
	esperadoP1.StoryPoints = 8
	verificarItemBacklog(t, estimado, esperadoP1)

	coleccion := cargarBacklog(t, almacen)
	verificarItemBacklog(t, buscarRegistroItem(t, coleccion, "p1", "HU-01").Item, esperadoP1)
	verificarItemBacklog(t, buscarRegistroItem(t, coleccion, "p2", "HU-01").Item, p2Item)
}

// CA1 / CA2 / R7: ningún camino de error invoca Save ni altera lo persistido.
func TestUS11_EstimarItem_ErrorNoEscribe(t *testing.T) {
	casos := []struct {
		descripcion   string
		proyectoID    string
		itemID        string
		storyPoints   int
		errorEsperado error
	}{
		{
			descripcion:   "proyecto inexistente",
			proyectoID:    "ausente",
			itemID:        "HU-01",
			storyPoints:   5,
			errorEsperado: domain.ErrProyectoNoEncontrado,
		},
		{
			descripcion:   "ítem inexistente",
			proyectoID:    "p1",
			itemID:        "HU-02",
			storyPoints:   5,
			errorEsperado: domain.ErrItemNoEncontrado,
		},
		{
			descripcion:   "story points inválidos",
			proyectoID:    "p1",
			itemID:        "HU-01",
			storyPoints:   4,
			errorEsperado: domain.ErrStoryPointsInvalidos,
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			almacenReal, _ := nuevoAlmacenServicio(t)
			sembrarProyectos(t, almacenReal, []proyectoRegistrado{
				{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
			})
			original := itemUS11Backlog(3)
			sembrarBacklog(t, almacenReal, []itemRegistrado{{ProyectoID: "p1", Item: original}})
			almacen := &almacenCuentaSaves{real: almacenReal}
			servicio := NuevoServicioBacklog(almacen)

			estimado, err := servicio.EstimarItem(caso.proyectoID, caso.itemID, caso.storyPoints)
			if !errors.Is(err, caso.errorEsperado) {
				t.Errorf("err = %v, se esperaba %v", err, caso.errorEsperado)
			}
			verificarItemBacklogVacio(t, estimado)
			if almacen.llamadasASave != 0 {
				t.Errorf("Save fue invocado %d veces ante un error; no debía escribirse nada", almacen.llamadasASave)
			}
			verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacenReal), "p1", "HU-01").Item, original)
		})
	}
}

// E5 / D4: un backlog.json corrupto propaga storage.ErrCorrupt y no se mapea a
// un error de no encontrado.
func TestUS11_EstimarItem_JSONCorruptoPropagaError(t *testing.T) {
	almacenReal, dir := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacenReal, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
	})
	if err := os.WriteFile(filepath.Join(dir, "backlog.json"), []byte("{no es json"), 0o644); err != nil {
		t.Fatal(err)
	}
	almacen := &almacenCuentaSaves{real: almacenReal}
	servicio := NuevoServicioBacklog(almacen)

	estimado, err := servicio.EstimarItem("p1", "HU-01", 5)
	if !errors.Is(err, storage.ErrCorrupt) {
		t.Errorf("err = %v, se esperaba un error que envuelva %v", err, storage.ErrCorrupt)
	}
	if errors.Is(err, domain.ErrItemNoEncontrado) {
		t.Errorf("err = %v, no debería mapearse a ErrItemNoEncontrado", err)
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, no debería mapearse a ErrProyectoNoEncontrado", err)
	}
	verificarItemBacklogVacio(t, estimado)
	if almacen.llamadasASave != 0 {
		t.Errorf("Save fue invocado %d veces; no debía escribirse nada", almacen.llamadasASave)
	}
}

// E5 / D4: un fallo de Save se propaga envuelto y no se informa una estimación
// exitosa.
func TestUS11_EstimarItem_FalloDeEscrituraNoInformaExito(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
	})
	original := itemUS11Backlog(0)
	sembrarBacklog(t, almacen, []itemRegistrado{{ProyectoID: "p1", Item: original}})

	errSimulado := errors.New("fallo simulado de escritura")
	servicio := NuevoServicioBacklog(almacenFalloEscritura{real: almacen, err: errSimulado})

	estimado, err := servicio.EstimarItem("p1", "HU-01", 5)
	if !errors.Is(err, errSimulado) {
		t.Errorf("err = %v, se esperaba un error que envuelva %v", err, errSimulado)
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) ||
		errors.Is(err, domain.ErrItemNoEncontrado) ||
		errors.Is(err, domain.ErrStoryPointsInvalidos) {
		t.Errorf("err = %v, no debería ser un error de dominio de no encontrado o de valor inválido", err)
	}
	verificarItemBacklogVacio(t, estimado)
	verificarItemBacklog(t, buscarRegistroItem(t, cargarBacklog(t, almacen), "p1", "HU-01").Item, original)
}
