package service

import (
	"errors"
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

// CA2 / R2 / C2: un duplicado exacto devuelve ErrIntegranteDuplicado y no se
// registra (D6).
func TestUS03_AgregarIntegrante_DuplicadoExacto(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	existente := domain.Integrante{Nombre: "Dana Guzmán", Rol: "Desarrollador"}
	sembrarProyectos(t, almacen, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)},
	})
	sembrarIntegrantes(t, almacen, []integranteRegistrado{
		{ProyectoID: "p1", Integrante: existente},
	})
	servicio := NuevoServicioProyectos(almacen)

	agregado, err := servicio.AgregarIntegrante("p1", "Dana Guzmán", "Tester")
	if err == nil {
		t.Fatal("AgregarIntegrante() con un duplicado exacto debería devolver un error")
	}
	if !errors.Is(err, domain.ErrIntegranteDuplicado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrIntegranteDuplicado)
	}
	verificarIntegranteVacio(t, agregado)

	coleccion := cargarIntegrantes(t, almacen)
	if len(coleccion) != 1 {
		t.Fatalf("se esperaba 1 integrante persistido, hay %d", len(coleccion))
	}
	verificarIntegrante(t, buscarRegistroIntegrante(t, coleccion, "p1", "Dana Guzmán").Integrante, existente)
}

// CA4 / R7 / R6: un rol vacío o compuesto solo por espacios devuelve
// ErrRolIntegranteVacio, no agrega y conserva lo persistido (D6).
func TestUS03_AgregarIntegrante_RolVacioNoAgregaNiPersiste(t *testing.T) {
	casos := []struct {
		descripcion string
		rol         string
	}{
		{descripcion: "rol vacío", rol: ""},
		{descripcion: "rol compuesto solo por espacios", rol: "   "},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			almacenReal, _ := nuevoAlmacenServicio(t)
			existente := domain.Integrante{Nombre: "Dana Guzmán", Rol: "Desarrollador"}
			sembrarProyectos(t, almacenReal, []proyectoRegistrado{
				{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)},
			})
			sembrarIntegrantes(t, almacenReal, []integranteRegistrado{
				{ProyectoID: "p1", Integrante: existente},
			})
			almacen := &almacenCuentaSaves{real: almacenReal}
			servicio := NuevoServicioProyectos(almacen)

			agregado, err := servicio.AgregarIntegrante("p1", "Victoria Caro Boldrini", caso.rol)
			if err == nil {
				t.Fatal("AgregarIntegrante() con rol vacío debería devolver un error")
			}
			if !errors.Is(err, domain.ErrRolIntegranteVacio) {
				t.Errorf("err = %v, se esperaba %v", err, domain.ErrRolIntegranteVacio)
			}
			verificarIntegranteVacio(t, agregado)

			if almacen.llamadasASave != 0 {
				t.Errorf("Save fue invocado %d veces; no debía escribirse nada", almacen.llamadasASave)
			}

			coleccion := cargarIntegrantes(t, almacenReal)
			if len(coleccion) != 1 {
				t.Fatalf("se esperaba 1 integrante persistido, hay %d", len(coleccion))
			}
			verificarIntegrante(t, buscarRegistroIntegrante(t, coleccion, "p1", "Dana Guzmán").Integrante, existente)
		})
	}
}

// CA2 / CA3 / CA4 / R6: ningún camino de error escribe la colección
// "integrantes" ni altera lo ya persistido.
func TestUS03_AgregarIntegrante_ErrorNoAgregaNiPersiste(t *testing.T) {
	casos := []struct {
		descripcion   string
		nombre        string
		rol           string
		errorEsperado error
	}{
		{
			descripcion:   "nombre vacío",
			nombre:        "",
			rol:           "Desarrollador",
			errorEsperado: domain.ErrNombreIntegranteVacio,
		},
		{
			descripcion:   "rol vacío",
			nombre:        "Victoria Caro Boldrini",
			rol:           "",
			errorEsperado: domain.ErrRolIntegranteVacio,
		},
		{
			descripcion:   "duplicado exacto",
			nombre:        "Dana Guzmán",
			rol:           "Desarrollador",
			errorEsperado: domain.ErrIntegranteDuplicado,
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			almacenReal, _ := nuevoAlmacenServicio(t)
			existente := domain.Integrante{Nombre: "Dana Guzmán", Rol: "Desarrollador"}
			sembrarProyectos(t, almacenReal, []proyectoRegistrado{
				{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)},
			})
			sembrarIntegrantes(t, almacenReal, []integranteRegistrado{
				{ProyectoID: "p1", Integrante: existente},
			})
			almacen := &almacenCuentaSaves{real: almacenReal}
			servicio := NuevoServicioProyectos(almacen)

			agregado, err := servicio.AgregarIntegrante("p1", caso.nombre, caso.rol)
			if err == nil {
				t.Fatal("AgregarIntegrante() con datos inválidos debería devolver un error")
			}
			if !errors.Is(err, caso.errorEsperado) {
				t.Errorf("err = %v, se esperaba %v", err, caso.errorEsperado)
			}
			verificarIntegranteVacio(t, agregado)

			// En ningún camino de error se escribe "integrantes".
			if almacen.llamadasASave != 0 {
				t.Errorf("Save fue invocado %d veces ante un error; no debía escribirse nada", almacen.llamadasASave)
			}

			coleccion := cargarIntegrantes(t, almacenReal)
			if len(coleccion) != 1 {
				t.Fatalf("se esperaba 1 integrante persistido, hay %d", len(coleccion))
			}
			verificarIntegrante(t, buscarRegistroIntegrante(t, coleccion, "p1", "Dana Guzmán").Integrante, existente)
		})
	}
}

// CA1 / R5 / C1: un ID ausente en una colección con otros proyectos devuelve
// ErrProyectoNoEncontrado y no escribe nada.
func TestUS03_AgregarIntegrante_ProyectoNoEncontrado(t *testing.T) {
	almacenReal, dir := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacenReal, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)},
		{ID: "p2", Proyecto: proyectoUS02Servicio("Otro proyecto", domain.EstadoPlanificado)},
	})
	almacen := &almacenCuentaSaves{real: almacenReal}
	servicio := NuevoServicioProyectos(almacen)

	agregado, err := servicio.AgregarIntegrante("ausente", "Dana Guzmán", "Desarrollador")
	if err == nil {
		t.Fatal("AgregarIntegrante() con un ID inexistente debería devolver un error")
	}
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	verificarIntegranteVacio(t, agregado)

	if almacen.llamadasASave != 0 {
		t.Errorf("Save fue invocado %d veces; no debía escribirse nada", almacen.llamadasASave)
	}
	if _, err := os.Stat(filepath.Join(dir, "integrantes.json")); !os.IsNotExist(err) {
		t.Fatalf("integrantes.json no debería existir ante un proyecto inexistente: err = %v", err)
	}
}

// CA1 / R5 / R4 / C7: la existencia del proyecto se verifica antes que el
// nombre vacío.
func TestUS03_AgregarIntegrante_ProyectoInexistenteAntesQueNombreVacio(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)},
	})
	servicio := NuevoServicioProyectos(almacen)

	agregado, err := servicio.AgregarIntegrante("ausente", "", "Desarrollador")
	if err == nil {
		t.Fatal("AgregarIntegrante() con un ID inexistente debería devolver un error")
	}
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if errors.Is(err, domain.ErrNombreIntegranteVacio) {
		t.Errorf("err = %v, no debería devolverse ErrNombreIntegranteVacio si el proyecto no existe", err)
	}
	verificarIntegranteVacio(t, agregado)
}

// CA1 / R5 / D1: sin archivo de la colección "proyectos" se interpreta como
// colección vacía, el proyecto no se encuentra y no se crea nada.
func TestUS03_AgregarIntegrante_ArchivoProyectosInexistente(t *testing.T) {
	almacen, dir := nuevoAlmacenServicio(t)
	servicio := NuevoServicioProyectos(almacen)

	agregado, err := servicio.AgregarIntegrante("p1", "Dana Guzmán", "Desarrollador")
	if err == nil {
		t.Fatal("AgregarIntegrante() sin archivo de proyectos debería devolver un error")
	}
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	verificarIntegranteVacio(t, agregado)

	entradas, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir(%q): %v", dir, err)
	}
	if len(entradas) != 0 {
		t.Errorf("en ningún camino de error se crea el archivo, hay %d entradas", len(entradas))
	}
}

// E5 / D1: un JSON corrupto propaga storage.ErrCorrupt y no se mapea a
// ErrProyectoNoEncontrado.
func TestUS03_AgregarIntegrante_JSONCorruptoNoEsNoEncontrado(t *testing.T) {
	almacen, dir := nuevoAlmacenServicio(t)
	if err := os.WriteFile(filepath.Join(dir, "proyectos.json"), []byte("{no es json"), 0o644); err != nil {
		t.Fatal(err)
	}
	servicio := NuevoServicioProyectos(almacen)

	agregado, err := servicio.AgregarIntegrante("p1", "Dana Guzmán", "Desarrollador")
	if err == nil {
		t.Fatal("AgregarIntegrante() con JSON corrupto debería devolver un error")
	}
	if !errors.Is(err, storage.ErrCorrupt) {
		t.Errorf("err = %v, se esperaba un error que envuelva %v", err, storage.ErrCorrupt)
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, no debería mapearse a ErrProyectoNoEncontrado", err)
	}
	verificarIntegranteVacio(t, agregado)
}

// E5 / D1: un fallo de escritura se propaga envuelto y no se informa un alta
// exitosa.
func TestUS03_AgregarIntegrante_FalloDeEscrituraNoInformaExito(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.EstadoPlanificado)},
	})

	errSimulado := errors.New("fallo simulado de escritura")
	servicio := NuevoServicioProyectos(almacenFalloEscritura{real: almacen, err: errSimulado})

	agregado, err := servicio.AgregarIntegrante("p1", "Dana Guzmán", "Desarrollador")
	if err == nil {
		t.Fatal("AgregarIntegrante() con fallo de escritura debería devolver un error")
	}
	if !errors.Is(err, errSimulado) {
		t.Errorf("err = %v, se esperaba un error que envuelva %v", err, errSimulado)
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, no debería ser ErrProyectoNoEncontrado", err)
	}
	if errors.Is(err, domain.ErrNombreIntegranteVacio) || errors.Is(err, domain.ErrRolIntegranteVacio) ||
		errors.Is(err, domain.ErrIntegranteDuplicado) {
		t.Errorf("err = %v, no debería ser un error de validación de dominio", err)
	}
	verificarIntegranteVacio(t, agregado)
}
