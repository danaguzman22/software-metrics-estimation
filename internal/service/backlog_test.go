package service

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage"
)

// sembrarBacklog persiste la colecciÃ³n completa bajo el nombre lÃ³gico
// "backlog" (D2).
func sembrarBacklog(t *testing.T, almacen *storage.JSONStore, registros []itemRegistrado) {
	t.Helper()
	if err := almacen.Save("backlog", registros); err != nil {
		t.Fatalf("Save(\"backlog\"): %v", err)
	}
}

// cargarBacklog recarga la colecciÃ³n desde disco usando el nombre lÃ³gico
// "backlog" (D2), para verificar la persistencia real.
func cargarBacklog(t *testing.T, almacen *storage.JSONStore) []itemRegistrado {
	t.Helper()
	var registros []itemRegistrado
	if err := almacen.Load("backlog", &registros); err != nil {
		t.Fatalf("Load(\"backlog\"): %v", err)
	}
	return registros
}

// buscarRegistroItem devuelve el registro del proyecto e ID indicados.
func buscarRegistroItem(t *testing.T, registros []itemRegistrado, proyectoID, id string) itemRegistrado {
	t.Helper()
	for _, registro := range registros {
		if registro.ProyectoID == proyectoID && registro.Item.ID == id {
			return registro
		}
	}
	t.Fatalf("no se encontrÃ³ el Ã­tem %q del proyecto %q", id, proyectoID)
	return itemRegistrado{}
}

// verificarItemBacklog compara campo por campo un Ã­tem del backlog. Los
// criterios se comparan por contenido y nil se considera equivalente a una
// lista vacÃ­a.
func verificarItemBacklog(t *testing.T, obtenido, esperado domain.ItemBacklog) {
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

// verificarItemBacklogVacio comprueba que un error devuelve el valor cero de
// domain.ItemBacklog (D11).
func verificarItemBacklogVacio(t *testing.T, item domain.ItemBacklog) {
	t.Helper()
	if item.ID != "" || item.Titulo != "" || item.Descripcion != "" ||
		item.Prioridad != "" || item.Estado != "" || item.StoryPoints != 0 ||
		len(item.Criterios) != 0 {
		t.Errorf("se esperaba domain.ItemBacklog{} ante un error, se obtuvo %+v", item)
	}
}

// CA1 / R1: crear un Ã­tem vÃ¡lido en un proyecto existente lo devuelve y lo
// persiste en la colecciÃ³n "backlog" (D2, D11).
func TestUS05_CrearItem_Exitoso(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []ProyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de MÃ©tricas", domain.EstadoPlanificado)},
	})
	servicio := NuevoServicioBacklog(almacen)

	creado, err := servicio.CrearItem(
		"p1",
		"Registrar proyecto",
		"Permitir dar de alta proyectos en el sistema",
		"Alta",
		[]string{
			"El nombre del proyecto no puede estar vacÃ­o",
			"El proyecto queda disponible para asignar integrantes",
		},
	)
	if err != nil {
		t.Fatalf("CrearItem() devolviÃ³ un error inesperado: %v", err)
	}
	esperado := domain.ItemBacklog{
		ID:          "HU-01",
		Titulo:      "Registrar proyecto",
		Descripcion: "Permitir dar de alta proyectos en el sistema",
		Prioridad:   "Alta",
		Estado:      "Pendiente",
		StoryPoints: 0,
		Criterios: []string{
			"El nombre del proyecto no puede estar vacÃ­o",
			"El proyecto queda disponible para asignar integrantes",
		},
	}
	verificarItemBacklog(t, creado, esperado)

	coleccion := cargarBacklog(t, almacen)
	if len(coleccion) != 1 {
		t.Fatalf("se esperaba 1 Ã­tem persistido, hay %d", len(coleccion))
	}
	registro := buscarRegistroItem(t, coleccion, "p1", "HU-01")
	if registro.ProyectoID != "p1" {
		t.Errorf("ProyectoID = %q, se esperaba %q", registro.ProyectoID, "p1")
	}
	verificarItemBacklog(t, registro.Item, esperado)
}

// CA1 / R1: un segundo Ã­tem se agrega junto al que ya estaba (D2).
func TestUS05_CrearItem_SegundoItemEnMismoProyecto(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []ProyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de MÃ©tricas", domain.EstadoPlanificado)},
	})
	sembrarBacklog(t, almacen, []itemRegistrado{
		{
			ProyectoID: "p1",
			Item: domain.ItemBacklog{
				ID:          "HU-01",
				Titulo:      "Registrar proyecto",
				Descripcion: "Permitir dar de alta proyectos en el sistema",
				Prioridad:   "Alta",
				Estado:      "Pendiente",
				StoryPoints: 0,
			},
		},
	})
	servicio := NuevoServicioBacklog(almacen)

	creado, err := servicio.CrearItem("p1", "Modificar proyecto", "Permitir cambiar los datos de un proyecto", "Media", nil)
	if err != nil {
		t.Fatalf("CrearItem() devolviÃ³ un error inesperado: %v", err)
	}
	esperado := domain.ItemBacklog{
		ID:          "HU-02",
		Titulo:      "Modificar proyecto",
		Descripcion: "Permitir cambiar los datos de un proyecto",
		Prioridad:   "Media",
		Estado:      "Pendiente",
		StoryPoints: 0,
	}
	verificarItemBacklog(t, creado, esperado)

	coleccion := cargarBacklog(t, almacen)
	if len(coleccion) != 2 {
		t.Fatalf("se esperaban 2 Ã­tems persistidos, hay %d", len(coleccion))
	}
	if registro := buscarRegistroItem(t, coleccion, "p1", "HU-01"); registro.Item.Titulo != "Registrar proyecto" {
		t.Errorf("Titulo de HU-01 = %q, se esperaba %q", registro.Item.Titulo, "Registrar proyecto")
	}
	verificarItemBacklog(t, buscarRegistroItem(t, coleccion, "p1", "HU-02").Item, esperado)
}

// CA1 / R1 / R10: sin backlog.json la primera alta es exitosa y crea el
// archivo.
func TestUS05_CrearItem_PrimerAltaSinArchivoBacklog(t *testing.T) {
	almacen, dir := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []ProyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de MÃ©tricas", domain.EstadoPlanificado)},
	})
	servicio := NuevoServicioBacklog(almacen)

	ruta := filepath.Join(dir, "backlog.json")
	if _, err := os.Stat(ruta); !os.IsNotExist(err) {
		t.Fatalf("backlog.json no deberÃ­a existir antes del alta: err = %v", err)
	}

	creado, err := servicio.CrearItem("p1", "Registrar proyecto", "Permitir dar de alta proyectos en el sistema", "Alta", nil)
	if err != nil {
		t.Fatalf("CrearItem() devolviÃ³ un error inesperado: %v", err)
	}
	if creado.ID != "HU-01" {
		t.Errorf("ID = %q, se esperaba %q", creado.ID, "HU-01")
	}

	if _, err := os.Stat(ruta); err != nil {
		t.Fatalf("backlog.json deberÃ­a existir despuÃ©s del alta: %v", err)
	}
	coleccion := cargarBacklog(t, almacen)
	if len(coleccion) != 1 {
		t.Fatalf("se esperaba 1 Ã­tem persistido, hay %d", len(coleccion))
	}
}

// CA2 / R2: el primer Ã­tem de un proyecto recibe el ID "HU-01".
func TestUS05_CrearItem_PrimerItemIDHU01(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []ProyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de MÃ©tricas", domain.EstadoPlanificado)},
	})
	servicio := NuevoServicioBacklog(almacen)

	creado, err := servicio.CrearItem("p1", "Registrar proyecto", "Permitir dar de alta proyectos en el sistema", "Alta", nil)
	if err != nil {
		t.Fatalf("CrearItem() devolviÃ³ un error inesperado: %v", err)
	}
	if creado.ID != "HU-01" {
		t.Errorf("ID = %q, se esperaba %q", creado.ID, "HU-01")
	}
}

// CA2 / R2: el ID es el siguiente al mayor existente, sin rellenar huecos.
func TestUS05_CrearItem_IDConsecutivoAlMayor(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []ProyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de MÃ©tricas", domain.EstadoPlanificado)},
	})
	sembrarBacklog(t, almacen, []itemRegistrado{
		{ProyectoID: "p1", Item: domain.ItemBacklog{ID: "HU-01", Titulo: "Registrar proyecto", Prioridad: "Alta", Estado: "Pendiente"}},
		{ProyectoID: "p1", Item: domain.ItemBacklog{ID: "HU-03", Titulo: "Consultar proyecto", Prioridad: "Media", Estado: "Pendiente"}},
	})
	servicio := NuevoServicioBacklog(almacen)

	creado, err := servicio.CrearItem("p1", "Consultar mÃ©tricas", "Consultar las mÃ©tricas del proyecto", "Baja", nil)
	if err != nil {
		t.Fatalf("CrearItem() devolviÃ³ un error inesperado: %v", err)
	}
	if creado.ID != "HU-04" {
		t.Errorf("ID = %q, se esperaba %q", creado.ID, "HU-04")
	}

	coleccion := cargarBacklog(t, almacen)
	if len(coleccion) != 3 {
		t.Fatalf("se esperaban 3 Ã­tems persistidos, hay %d", len(coleccion))
	}
	if registro := buscarRegistroItem(t, coleccion, "p1", "HU-04"); registro.Item.Titulo != "Consultar mÃ©tricas" {
		t.Errorf("Titulo de HU-04 = %q, se esperaba %q", registro.Item.Titulo, "Consultar mÃ©tricas")
	}
}

// CA2 / R2: dos altas en el mismo proyecto generan IDs distintos.
func TestUS05_CrearItem_IDUnicoEnProyecto(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []ProyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de MÃ©tricas", domain.EstadoPlanificado)},
	})
	servicio := NuevoServicioBacklog(almacen)

	primero, err := servicio.CrearItem("p1", "Registrar proyecto", "Permitir dar de alta proyectos en el sistema", "Alta", nil)
	if err != nil {
		t.Fatalf("CrearItem() del primer Ã­tem devolviÃ³ un error inesperado: %v", err)
	}
	segundo, err := servicio.CrearItem("p1", "Modificar proyecto", "Permitir cambiar los datos de un proyecto", "Media", nil)
	if err != nil {
		t.Fatalf("CrearItem() del segundo Ã­tem devolviÃ³ un error inesperado: %v", err)
	}
	if primero.ID == segundo.ID {
		t.Errorf("los IDs no deberÃ­an repetirse: ambos son %q", primero.ID)
	}
	if primero.ID != "HU-01" || segundo.ID != "HU-02" {
		t.Errorf("IDs = %q y %q, se esperaban %q y %q", primero.ID, segundo.ID, "HU-01", "HU-02")
	}

	coleccion := cargarBacklog(t, almacen)
	if len(coleccion) != 2 {
		t.Fatalf("se esperaban 2 Ã­tems persistidos, hay %d", len(coleccion))
	}
}

// CA2 / R2: la numeraciÃ³n es independiente por proyecto; el mismo ID puede
// repetirse en proyectos distintos.
func TestUS05_CrearItem_MismoIDEnOtroProyecto(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []ProyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de MÃ©tricas", domain.EstadoPlanificado)},
		{ID: "p2", Proyecto: proyectoUS02Servicio("App MÃ³vil", domain.EstadoPlanificado)},
	})
	sembrarBacklog(t, almacen, []itemRegistrado{
		{ProyectoID: "p1", Item: domain.ItemBacklog{ID: "HU-01", Titulo: "Registrar proyecto", Prioridad: "Alta", Estado: "Pendiente"}},
	})
	servicio := NuevoServicioBacklog(almacen)

	creado, err := servicio.CrearItem("p2", "Registrar proyecto", "Permitir dar de alta proyectos en el sistema", "Alta", nil)
	if err != nil {
		t.Fatalf("CrearItem() devolviÃ³ un error inesperado: %v", err)
	}
	if creado.ID != "HU-01" {
		t.Errorf("ID = %q, se esperaba %q", creado.ID, "HU-01")
	}

	coleccion := cargarBacklog(t, almacen)
	if len(coleccion) != 2 {
		t.Fatalf("se esperaban 2 Ã­tems persistidos, hay %d", len(coleccion))
	}
	if registro := buscarRegistroItem(t, coleccion, "p1", "HU-01"); registro.Item.Titulo != "Registrar proyecto" {
		t.Errorf("el proyecto p1 deberÃ­a conservar su Ã­tem HU-01, Titulo = %q", registro.Item.Titulo)
	}
	verificarItemBacklog(t, buscarRegistroItem(t, coleccion, "p2", "HU-01").Item, creado)
}

// CA3 / R3 / R9: un tÃ­tulo vacÃ­o o solo espacios devuelve ErrTituloItemVacio,
// no agrega y conserva lo persistido.
func TestUS05_CrearItem_TituloVacioNoAgregaNiPersiste(t *testing.T) {
	casos := []struct {
		descripcion string
		titulo      string
	}{
		{descripcion: "tÃ­tulo vacÃ­o", titulo: ""},
		{descripcion: "tÃ­tulo compuesto solo por espacios", titulo: "   "},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			almacenReal, _ := nuevoAlmacenServicio(t)
			existente := domain.ItemBacklog{ID: "HU-01", Titulo: "Registrar proyecto", Prioridad: "Alta", Estado: "Pendiente"}
			sembrarProyectos(t, almacenReal, []ProyectoRegistrado{
				{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de MÃ©tricas", domain.EstadoPlanificado)},
			})
			sembrarBacklog(t, almacenReal, []itemRegistrado{
				{ProyectoID: "p1", Item: existente},
			})
			almacen := &almacenCuentaSaves{real: almacenReal}
			servicio := NuevoServicioBacklog(almacen)

			creado, err := servicio.CrearItem("p1", caso.titulo, "Permitir dar de alta proyectos en el sistema", "Alta", nil)
			if err == nil {
				t.Fatal("CrearItem() con tÃ­tulo vacÃ­o deberÃ­a devolver un error")
			}
			if !errors.Is(err, domain.ErrTituloItemVacio) {
				t.Errorf("err = %v, se esperaba %v", err, domain.ErrTituloItemVacio)
			}
			verificarItemBacklogVacio(t, creado)

			if almacen.llamadasASave != 0 {
				t.Errorf("Save fue invocado %d veces; no debÃ­a escribirse nada", almacen.llamadasASave)
			}

			coleccion := cargarBacklog(t, almacenReal)
			if len(coleccion) != 1 {
				t.Fatalf("se esperaba 1 Ã­tem persistido, hay %d", len(coleccion))
			}
			verificarItemBacklog(t, buscarRegistroItem(t, coleccion, "p1", "HU-01").Item, existente)
		})
	}
}

// CA3 / R4 / R9: una prioridad invÃ¡lida devuelve ErrPrioridadInvalida, no
// agrega y conserva lo persistido.
func TestUS05_CrearItem_PrioridadInvalidaNoAgregaNiPersiste(t *testing.T) {
	casos := []struct {
		descripcion string
		prioridad   string
	}{
		{descripcion: "prioridad no permitida", prioridad: "Urgente"},
		{descripcion: "prioridad en minÃºsculas", prioridad: "alta"},
		{descripcion: "prioridad en mayÃºsculas", prioridad: "ALTA"},
		{descripcion: "prioridad compuesta solo por espacios", prioridad: "   "},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			almacenReal, _ := nuevoAlmacenServicio(t)
			existente := domain.ItemBacklog{ID: "HU-01", Titulo: "Registrar proyecto", Prioridad: "Alta", Estado: "Pendiente"}
			sembrarProyectos(t, almacenReal, []ProyectoRegistrado{
				{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de MÃ©tricas", domain.EstadoPlanificado)},
			})
			sembrarBacklog(t, almacenReal, []itemRegistrado{
				{ProyectoID: "p1", Item: existente},
			})
			almacen := &almacenCuentaSaves{real: almacenReal}
			servicio := NuevoServicioBacklog(almacen)

			creado, err := servicio.CrearItem("p1", "Registrar proyecto", "Permitir dar de alta proyectos en el sistema", caso.prioridad, nil)
			if err == nil {
				t.Fatal("CrearItem() con prioridad invÃ¡lida deberÃ­a devolver un error")
			}
			if !errors.Is(err, domain.ErrPrioridadInvalida) {
				t.Errorf("err = %v, se esperaba %v", err, domain.ErrPrioridadInvalida)
			}
			verificarItemBacklogVacio(t, creado)

			if almacen.llamadasASave != 0 {
				t.Errorf("Save fue invocado %d veces; no debÃ­a escribirse nada", almacen.llamadasASave)
			}

			coleccion := cargarBacklog(t, almacenReal)
			if len(coleccion) != 1 {
				t.Fatalf("se esperaba 1 Ã­tem persistido, hay %d", len(coleccion))
			}
			verificarItemBacklog(t, buscarRegistroItem(t, coleccion, "p1", "HU-01").Item, existente)
		})
	}
}

// CA4 / R6: el estado inicial del Ã­tem devuelto y persistido es "Pendiente".
func TestUS05_CrearItem_EstadoInicialPendiente(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []ProyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de MÃ©tricas", domain.EstadoPlanificado)},
	})
	servicio := NuevoServicioBacklog(almacen)

	creado, err := servicio.CrearItem("p1", "Registrar proyecto", "Permitir dar de alta proyectos en el sistema", "Alta", nil)
	if err != nil {
		t.Fatalf("CrearItem() devolviÃ³ un error inesperado: %v", err)
	}
	if creado.Estado != "Pendiente" {
		t.Errorf("Estado devuelto = %q, se esperaba %q", creado.Estado, "Pendiente")
	}

	coleccion := cargarBacklog(t, almacen)
	persistido := buscarRegistroItem(t, coleccion, "p1", creado.ID)
	if persistido.Item.Estado != "Pendiente" {
		t.Errorf("Estado persistido = %q, se esperaba %q", persistido.Item.Estado, "Pendiente")
	}
}

// CA1 / R7: los story points del Ã­tem devuelto y persistido quedan en 0.
func TestUS05_CrearItem_StoryPointsCero(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []ProyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de MÃ©tricas", domain.EstadoPlanificado)},
	})
	servicio := NuevoServicioBacklog(almacen)

	creado, err := servicio.CrearItem("p1", "Registrar proyecto", "Permitir dar de alta proyectos en el sistema", "Alta", nil)
	if err != nil {
		t.Fatalf("CrearItem() devolviÃ³ un error inesperado: %v", err)
	}
	if creado.StoryPoints != 0 {
		t.Errorf("StoryPoints devueltos = %d, se esperaban 0", creado.StoryPoints)
	}

	coleccion := cargarBacklog(t, almacen)
	persistido := buscarRegistroItem(t, coleccion, "p1", creado.ID)
	if persistido.Item.StoryPoints != 0 {
		t.Errorf("StoryPoints persistidos = %d, se esperaban 0", persistido.Item.StoryPoints)
	}
}

// CA1 / R8: descripciÃ³n vacÃ­a y criterios ausentes son vÃ¡lidos.
func TestUS05_CrearItem_DescripcionYCriteriosOpcionales(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []ProyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de MÃ©tricas", domain.EstadoPlanificado)},
	})
	servicio := NuevoServicioBacklog(almacen)

	creado, err := servicio.CrearItem("p1", "Revisar mÃ©tricas", "", "Baja", nil)
	if err != nil {
		t.Fatalf("CrearItem() devolviÃ³ un error inesperado: %v", err)
	}
	if creado.Descripcion != "" {
		t.Errorf("Descripcion = %q, se esperaba vacÃ­a", creado.Descripcion)
	}
	if len(creado.Criterios) != 0 {
		t.Errorf("len(Criterios) = %d, se esperaba 0", len(creado.Criterios))
	}

	coleccion := cargarBacklog(t, almacen)
	persistido := buscarRegistroItem(t, coleccion, "p1", creado.ID)
	if persistido.Item.Descripcion != "" {
		t.Errorf("Descripcion persistida = %q, se esperaba vacÃ­a", persistido.Item.Descripcion)
	}
	if len(persistido.Item.Criterios) != 0 {
		t.Errorf("len(Criterios) persistidos = %d, se esperaba 0", len(persistido.Item.Criterios))
	}
}

// CA1 / R8: descripciÃ³n y criterios admiten texto libre y se guardan sin
// cambios.
func TestUS05_CrearItem_DescripcionYCriteriosSinValidar(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []ProyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de MÃ©tricas", domain.EstadoPlanificado)},
	})
	servicio := NuevoServicioBacklog(almacen)

	descripcion := "Â¿QuÃ© mÃ©tricas se van a calcular? AÃºn no estÃ¡ definido."
	criterios := []string{"Se listan las mÃ©tricas", "Cada mÃ©trica tiene una descripciÃ³n"}
	creado, err := servicio.CrearItem("p1", "Analizar mÃ©tricas", descripcion, "Media", criterios)
	if err != nil {
		t.Fatalf("CrearItem() devolviÃ³ un error inesperado: %v", err)
	}
	esperado := domain.ItemBacklog{
		ID:          "HU-01",
		Titulo:      "Analizar mÃ©tricas",
		Descripcion: descripcion,
		Prioridad:   "Media",
		Estado:      "Pendiente",
		StoryPoints: 0,
		Criterios:   criterios,
	}
	verificarItemBacklog(t, creado, esperado)

	coleccion := cargarBacklog(t, almacen)
	verificarItemBacklog(t, buscarRegistroItem(t, coleccion, "p1", "HU-01").Item, esperado)
}

// CA1 / R5: un ID ausente en una colecciÃ³n con otros proyectos devuelve
// ErrProyectoNoEncontrado y no escribe nada.
func TestUS05_CrearItem_ProyectoNoEncontrado(t *testing.T) {
	almacenReal, dir := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacenReal, []ProyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de MÃ©tricas", domain.EstadoPlanificado)},
		{ID: "p2", Proyecto: proyectoUS02Servicio("Otro proyecto", domain.EstadoPlanificado)},
	})
	almacen := &almacenCuentaSaves{real: almacenReal}
	servicio := NuevoServicioBacklog(almacen)

	creado, err := servicio.CrearItem("ausente", "Registrar proyecto", "Permitir dar de alta proyectos en el sistema", "Alta", nil)
	if err == nil {
		t.Fatal("CrearItem() con un ID inexistente deberÃ­a devolver un error")
	}
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	verificarItemBacklogVacio(t, creado)

	if almacen.llamadasASave != 0 {
		t.Errorf("Save fue invocado %d veces; no debÃ­a escribirse nada", almacen.llamadasASave)
	}
	if _, err := os.Stat(filepath.Join(dir, "backlog.json")); !os.IsNotExist(err) {
		t.Fatalf("backlog.json no deberÃ­a existir ante un proyecto inexistente: err = %v", err)
	}
}

// CA1 / R5 / D2: sin archivo de la colecciÃ³n "proyectos" se interpreta como
// colecciÃ³n vacÃ­a, el proyecto no se encuentra y no se crea nada.
func TestUS05_CrearItem_ArchivoProyectosInexistente(t *testing.T) {
	almacen, dir := nuevoAlmacenServicio(t)
	servicio := NuevoServicioBacklog(almacen)

	creado, err := servicio.CrearItem("p1", "Registrar proyecto", "Permitir dar de alta proyectos en el sistema", "Alta", nil)
	if err == nil {
		t.Fatal("CrearItem() sin archivo de proyectos deberÃ­a devolver un error")
	}
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	verificarItemBacklogVacio(t, creado)

	entradas, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir(%q): %v", dir, err)
	}
	if len(entradas) != 0 {
		t.Errorf("en ningÃºn camino de error se crea el archivo, hay %d entradas", len(entradas))
	}
}

// CA1 / CA3 / R5 / R3 / D9: la existencia del proyecto se verifica antes que
// el tÃ­tulo vacÃ­o.
func TestUS05_CrearItem_ProyectoInexistenteAntesQueTituloVacio(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []ProyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de MÃ©tricas", domain.EstadoPlanificado)},
	})
	servicio := NuevoServicioBacklog(almacen)

	creado, err := servicio.CrearItem("ausente", "", "Permitir dar de alta proyectos en el sistema", "Alta", nil)
	if err == nil {
		t.Fatal("CrearItem() con un ID inexistente deberÃ­a devolver un error")
	}
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if errors.Is(err, domain.ErrTituloItemVacio) {
		t.Errorf("err = %v, no deberÃ­a devolverse ErrTituloItemVacio si el proyecto no existe", err)
	}
	verificarItemBacklogVacio(t, creado)
}

// CA3 / R3 / R4 / D9: el tÃ­tulo vacÃ­o se verifica antes que la prioridad
// invÃ¡lida.
func TestUS05_CrearItem_TituloVacioAntesQuePrioridadInvalida(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []ProyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de MÃ©tricas", domain.EstadoPlanificado)},
	})
	servicio := NuevoServicioBacklog(almacen)

	creado, err := servicio.CrearItem("p1", "", "Permitir dar de alta proyectos en el sistema", "Urgente", nil)
	if err == nil {
		t.Fatal("CrearItem() con tÃ­tulo vacÃ­o deberÃ­a devolver un error")
	}
	if !errors.Is(err, domain.ErrTituloItemVacio) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrTituloItemVacio)
	}
	if errors.Is(err, domain.ErrPrioridadInvalida) {
		t.Errorf("err = %v, no deberÃ­a devolverse ErrPrioridadInvalida si el tÃ­tulo estÃ¡ vacÃ­o", err)
	}
	verificarItemBacklogVacio(t, creado)
}

// CA1 / CA2 / CA3 / CA4 / R9: ningÃºn camino de error escribe la colecciÃ³n
// "backlog" ni altera lo ya persistido.
func TestUS05_CrearItem_ErrorNoAgregaNiPersiste(t *testing.T) {
	casos := []struct {
		descripcion   string
		proyectoID    string
		titulo        string
		prioridad     string
		errorEsperado error
	}{
		{
			descripcion:   "proyecto inexistente",
			proyectoID:    "ausente",
			titulo:        "Registrar proyecto",
			prioridad:     "Alta",
			errorEsperado: domain.ErrProyectoNoEncontrado,
		},
		{
			descripcion:   "tÃ­tulo vacÃ­o",
			proyectoID:    "p1",
			titulo:        "",
			prioridad:     "Alta",
			errorEsperado: domain.ErrTituloItemVacio,
		},
		{
			descripcion:   "prioridad invÃ¡lida",
			proyectoID:    "p1",
			titulo:        "Registrar proyecto",
			prioridad:     "Urgente",
			errorEsperado: domain.ErrPrioridadInvalida,
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			almacenReal, _ := nuevoAlmacenServicio(t)
			existente := domain.ItemBacklog{ID: "HU-01", Titulo: "Registrar proyecto", Prioridad: "Alta", Estado: "Pendiente"}
			sembrarProyectos(t, almacenReal, []ProyectoRegistrado{
				{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de MÃ©tricas", domain.EstadoPlanificado)},
			})
			sembrarBacklog(t, almacenReal, []itemRegistrado{
				{ProyectoID: "p1", Item: existente},
			})
			almacen := &almacenCuentaSaves{real: almacenReal}
			servicio := NuevoServicioBacklog(almacen)

			creado, err := servicio.CrearItem(caso.proyectoID, caso.titulo, "Permitir dar de alta proyectos en el sistema", caso.prioridad, nil)
			if err == nil {
				t.Fatal("CrearItem() con datos invÃ¡lidos deberÃ­a devolver un error")
			}
			if !errors.Is(err, caso.errorEsperado) {
				t.Errorf("err = %v, se esperaba %v", err, caso.errorEsperado)
			}
			verificarItemBacklogVacio(t, creado)

			if almacen.llamadasASave != 0 {
				t.Errorf("Save fue invocado %d veces ante un error; no debÃ­a escribirse nada", almacen.llamadasASave)
			}

			coleccion := cargarBacklog(t, almacenReal)
			if len(coleccion) != 1 {
				t.Fatalf("se esperaba 1 Ã­tem persistido, hay %d", len(coleccion))
			}
			verificarItemBacklog(t, buscarRegistroItem(t, coleccion, "p1", "HU-01").Item, existente)
		})
	}
}

// E5 / D2: un JSON corrupto propaga storage.ErrCorrupt y no se mapea a
// ErrProyectoNoEncontrado.
func TestUS05_CrearItem_JSONCorruptoNoEsNoEncontrado(t *testing.T) {
	almacen, dir := nuevoAlmacenServicio(t)
	if err := os.WriteFile(filepath.Join(dir, "proyectos.json"), []byte("{no es json"), 0o644); err != nil {
		t.Fatal(err)
	}
	servicio := NuevoServicioBacklog(almacen)

	creado, err := servicio.CrearItem("p1", "Registrar proyecto", "Permitir dar de alta proyectos en el sistema", "Alta", nil)
	if err == nil {
		t.Fatal("CrearItem() con JSON corrupto deberÃ­a devolver un error")
	}
	if !errors.Is(err, storage.ErrCorrupt) {
		t.Errorf("err = %v, se esperaba un error que envuelva %v", err, storage.ErrCorrupt)
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, no deberÃ­a mapearse a ErrProyectoNoEncontrado", err)
	}
	verificarItemBacklogVacio(t, creado)
}

// E5 / D2: un backlog.json corrupto propaga storage.ErrCorrupt y no se mapea a
// ErrProyectoNoEncontrado.
func TestUS05_CrearItem_BacklogCorruptoNoEsNoEncontrado(t *testing.T) {
	almacenReal, dir := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacenReal, []ProyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de MÃ©tricas", domain.EstadoPlanificado)},
	})
	if err := os.WriteFile(filepath.Join(dir, "backlog.json"), []byte("{no es json"), 0o644); err != nil {
		t.Fatal(err)
	}
	almacen := &almacenCuentaSaves{real: almacenReal}
	servicio := NuevoServicioBacklog(almacen)

	creado, err := servicio.CrearItem("p1", "Registrar proyecto", "Permitir dar de alta proyectos en el sistema", "Alta", nil)
	if err == nil {
		t.Fatal("CrearItem() con backlog corrupto deberÃ­a devolver un error")
	}
	if !errors.Is(err, storage.ErrCorrupt) {
		t.Errorf("err = %v, se esperaba un error que envuelva %v", err, storage.ErrCorrupt)
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, no deberÃ­a mapearse a ErrProyectoNoEncontrado", err)
	}
	verificarItemBacklogVacio(t, creado)

	if almacen.llamadasASave != 0 {
		t.Errorf("Save fue invocado %d veces; no debÃ­a escribirse nada", almacen.llamadasASave)
	}
}

// E5 / D2: un fallo de escritura se propaga envuelto y no se informa un alta
// exitosa.
func TestUS05_CrearItem_FalloDeEscrituraNoInformaExito(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []ProyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de MÃ©tricas", domain.EstadoPlanificado)},
	})

	errSimulado := errors.New("fallo simulado de escritura")
	servicio := NuevoServicioBacklog(almacenFalloEscritura{real: almacen, err: errSimulado})

	creado, err := servicio.CrearItem("p1", "Registrar proyecto", "Permitir dar de alta proyectos en el sistema", "Alta", nil)
	if err == nil {
		t.Fatal("CrearItem() con fallo de escritura deberÃ­a devolver un error")
	}
	if !errors.Is(err, errSimulado) {
		t.Errorf("err = %v, se esperaba un error que envuelva %v", err, errSimulado)
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, no deberÃ­a ser ErrProyectoNoEncontrado", err)
	}
	if errors.Is(err, domain.ErrTituloItemVacio) || errors.Is(err, domain.ErrPrioridadInvalida) {
		t.Errorf("err = %v, no deberÃ­a ser un error de validaciÃ³n de dominio", err)
	}
	verificarItemBacklogVacio(t, creado)
}
