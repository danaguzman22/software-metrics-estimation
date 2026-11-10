package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage"
)

// sembrarIntegrantes persiste la colección completa bajo el nombre lógico
// "integrantes" (D3).
func sembrarIntegrantes(t *testing.T, almacen *storage.JSONStore, registros []integranteRegistrado) {
	t.Helper()
	if err := almacen.Save("integrantes", registros); err != nil {
		t.Fatalf("Save(\"integrantes\"): %v", err)
	}
}

// cargarIntegrantes recarga la colección desde disco usando el nombre lógico
// "integrantes" (D3), para verificar la persistencia real.
func cargarIntegrantes(t *testing.T, almacen *storage.JSONStore) []integranteRegistrado {
	t.Helper()
	var registros []integranteRegistrado
	if err := almacen.Load("integrantes", &registros); err != nil {
		t.Fatalf("Load(\"integrantes\"): %v", err)
	}
	return registros
}

// buscarRegistroIntegrante devuelve el registro del proyecto y nombre indicados.
func buscarRegistroIntegrante(t *testing.T, registros []integranteRegistrado, proyectoID, nombre string) integranteRegistrado {
	t.Helper()
	for _, registro := range registros {
		if registro.ProyectoID == proyectoID && registro.Integrante.Nombre == nombre {
			return registro
		}
	}
	t.Fatalf("no se encontró el integrante %q del proyecto %q", nombre, proyectoID)
	return integranteRegistrado{}
}

// verificarIntegrante compara nombre y rol de un integrante.
func verificarIntegrante(t *testing.T, obtenido, esperado domain.Integrante) {
	t.Helper()
	if obtenido.Nombre != esperado.Nombre {
		t.Errorf("Nombre = %q, se esperaba %q", obtenido.Nombre, esperado.Nombre)
	}
	if obtenido.Rol != esperado.Rol {
		t.Errorf("Rol = %q, se esperaba %q", obtenido.Rol, esperado.Rol)
	}
}

// verificarIntegranteVacio comprueba que un error devuelve el valor cero de
// domain.Integrante (D5, R6).
func verificarIntegranteVacio(t *testing.T, integrante domain.Integrante) {
	t.Helper()
	if integrante != (domain.Integrante{}) {
		t.Errorf("se esperaba domain.Integrante{} ante un error, se obtuvo %+v", integrante)
	}
}

// CA1 / R1: registrar un integrante válido en un proyecto existente lo
// devuelve y lo persiste en la colección "integrantes" (D3, D5).
func TestUS03_AgregarIntegrante_Exitoso(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)},
	})
	servicio := NuevoServicioProyectos(almacen)

	agregado, err := servicio.AgregarIntegrante("p1", "Dana Guzmán", "Desarrollador")
	if err != nil {
		t.Fatalf("AgregarIntegrante() devolvió un error inesperado: %v", err)
	}
	esperado := domain.Integrante{Nombre: "Dana Guzmán", Rol: "Desarrollador"}
	verificarIntegrante(t, agregado, esperado)

	coleccion := cargarIntegrantes(t, almacen)
	if len(coleccion) != 1 {
		t.Fatalf("se esperaba 1 integrante persistido, hay %d", len(coleccion))
	}
	registro := buscarRegistroIntegrante(t, coleccion, "p1", "Dana Guzmán")
	if registro.ProyectoID != "p1" {
		t.Errorf("ProyectoID = %q, se esperaba %q", registro.ProyectoID, "p1")
	}
	verificarIntegrante(t, registro.Integrante, esperado)
}

// CA1 / R1: un segundo integrante se agrega junto al que ya estaba (D3, D5).
func TestUS03_AgregarIntegrante_SegundoIntegrante(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)},
	})
	sembrarIntegrantes(t, almacen, []integranteRegistrado{
		{ProyectoID: "p1", Integrante: domain.Integrante{Nombre: "Dana Guzmán", Rol: "Desarrollador"}},
	})
	servicio := NuevoServicioProyectos(almacen)

	agregado, err := servicio.AgregarIntegrante("p1", "Victoria Caro Boldrini", "Tester")
	if err != nil {
		t.Fatalf("AgregarIntegrante() devolvió un error inesperado: %v", err)
	}
	verificarIntegrante(t, agregado, domain.Integrante{Nombre: "Victoria Caro Boldrini", Rol: "Tester"})

	coleccion := cargarIntegrantes(t, almacen)
	if len(coleccion) != 2 {
		t.Fatalf("se esperaban 2 integrantes persistidos, hay %d", len(coleccion))
	}
	verificarIntegrante(t, buscarRegistroIntegrante(t, coleccion, "p1", "Dana Guzmán").Integrante,
		domain.Integrante{Nombre: "Dana Guzmán", Rol: "Desarrollador"})
	verificarIntegrante(t, buscarRegistroIntegrante(t, coleccion, "p1", "Victoria Caro Boldrini").Integrante,
		domain.Integrante{Nombre: "Victoria Caro Boldrini", Rol: "Tester"})
}

// CA1 / CA3 / R1 / R4: un nombre con espacios externos es válido y se guarda
// como "Ana Paz" (C6).
func TestUS03_AgregarIntegrante_NombreConEspaciosExternos(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)},
	})
	servicio := NuevoServicioProyectos(almacen)

	agregado, err := servicio.AgregarIntegrante("p1", "  Ana Paz  ", "Desarrollador Backend")
	if err != nil {
		t.Fatalf("AgregarIntegrante() devolvió un error inesperado: %v", err)
	}
	esperado := domain.Integrante{Nombre: "Ana Paz", Rol: "Desarrollador Backend"}
	verificarIntegrante(t, agregado, esperado)

	coleccion := cargarIntegrantes(t, almacen)
	registro := buscarRegistroIntegrante(t, coleccion, "p1", "Ana Paz")
	verificarIntegrante(t, registro.Integrante, esperado)
}

// CA1 / R1 / D1: sin integrantes.json la primera alta es exitosa y crea el
// archivo.
func TestUS03_AgregarIntegrante_PrimerAltaSinArchivoIntegrantes(t *testing.T) {
	almacen, dir := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)},
	})
	servicio := NuevoServicioProyectos(almacen)

	ruta := filepath.Join(dir, "integrantes.json")
	if _, err := os.Stat(ruta); !os.IsNotExist(err) {
		t.Fatalf("integrantes.json no debería existir antes del alta: err = %v", err)
	}

	agregado, err := servicio.AgregarIntegrante("p1", "Dana Guzmán", "Desarrollador")
	if err != nil {
		t.Fatalf("AgregarIntegrante() devolvió un error inesperado: %v", err)
	}
	verificarIntegrante(t, agregado, domain.Integrante{Nombre: "Dana Guzmán", Rol: "Desarrollador"})

	if _, err := os.Stat(ruta); err != nil {
		t.Fatalf("integrantes.json debería existir después del alta: %v", err)
	}
	coleccion := cargarIntegrantes(t, almacen)
	verificarIntegrante(t, buscarRegistroIntegrante(t, coleccion, "p1", "Dana Guzmán").Integrante,
		domain.Integrante{Nombre: "Dana Guzmán", Rol: "Desarrollador"})
}

// CA2 / R3: el mismo nombre en otro proyecto se permite; cada proyecto
// registra su integrante por separado (D3).
func TestUS03_AgregarIntegrante_MismoNombreEnOtroProyecto(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)},
		{ID: "p2", Proyecto: proyectoUS02Servicio("Otro proyecto", domain.EstadoPlanificado)},
	})
	servicio := NuevoServicioProyectos(almacen)

	agregadoP1, err := servicio.AgregarIntegrante("p1", "Dana Guzmán", "Desarrollador")
	if err != nil {
		t.Fatalf("AgregarIntegrante(p1) devolvió un error inesperado: %v", err)
	}
	verificarIntegrante(t, agregadoP1, domain.Integrante{Nombre: "Dana Guzmán", Rol: "Desarrollador"})

	agregadoP2, err := servicio.AgregarIntegrante("p2", "Dana Guzmán", "Tester")
	if err != nil {
		t.Fatalf("AgregarIntegrante(p2) devolvió un error inesperado: %v", err)
	}
	verificarIntegrante(t, agregadoP2, domain.Integrante{Nombre: "Dana Guzmán", Rol: "Tester"})

	coleccion := cargarIntegrantes(t, almacen)
	if len(coleccion) != 2 {
		t.Fatalf("se esperaban 2 integrantes persistidos, hay %d", len(coleccion))
	}
	verificarIntegrante(t, buscarRegistroIntegrante(t, coleccion, "p1", "Dana Guzmán").Integrante,
		domain.Integrante{Nombre: "Dana Guzmán", Rol: "Desarrollador"})
	verificarIntegrante(t, buscarRegistroIntegrante(t, coleccion, "p2", "Dana Guzmán").Integrante,
		domain.Integrante{Nombre: "Dana Guzmán", Rol: "Tester"})
}
