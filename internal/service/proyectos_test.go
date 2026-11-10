package service

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage"
)

// fechaUS02 arma una fecha UTC de los escenarios de US-02.
func fechaUS02(anio int, mes time.Month, dia int) time.Time {
	return time.Date(anio, mes, dia, 0, 0, 0, 0, time.UTC)
}

// proyectoUS02Servicio arma el proyecto de referencia de los escenarios de
// US-02 (inicio 2026-03-01 y fin 2026-06-30).
func proyectoUS02Servicio(nombre string, estado domain.Estado) domain.Proyecto {
	return domain.Proyecto{
		Nombre:      nombre,
		FechaInicio: fechaUS02(2026, time.March, 1),
		FechaFin:    fechaUS02(2026, time.June, 30),
		Estado:      estado,
	}
}

// nuevoAlmacenServicio crea un JSONStore real sobre un directorio temporal y
// devuelve también el directorio para inspeccionarlo desde disco.
func nuevoAlmacenServicio(t *testing.T) (*storage.JSONStore, string) {
	t.Helper()
	dir := t.TempDir()
	almacen, err := storage.NewJSONStore(dir)
	if err != nil {
		t.Fatalf("NewJSONStore: %v", err)
	}
	return almacen, dir
}

// sembrarProyectos persiste la colección completa bajo el nombre lógico
// "proyectos" (D13), tal como lo define el contrato de persistencia de US-02.
func sembrarProyectos(t *testing.T, almacen *storage.JSONStore, registros []proyectoRegistrado) {
	t.Helper()
	if err := almacen.Save("proyectos", registros); err != nil {
		t.Fatalf("Save(\"proyectos\"): %v", err)
	}
}

// cargarProyectos recarga la colección desde disco usando el nombre lógico
// "proyectos" (D13), para verificar la persistencia real.
func cargarProyectos(t *testing.T, almacen *storage.JSONStore) []proyectoRegistrado {
	t.Helper()
	var registros []proyectoRegistrado
	if err := almacen.Load("proyectos", &registros); err != nil {
		t.Fatalf("Load(\"proyectos\"): %v", err)
	}
	return registros
}

// buscarRegistro devuelve el registro de la colección con el ID indicado.
func buscarRegistro(t *testing.T, registros []proyectoRegistrado, id string) proyectoRegistrado {
	t.Helper()
	for _, registro := range registros {
		if registro.ID == id {
			return registro
		}
	}
	t.Fatalf("no se encontró el registro con ID %q", id)
	return proyectoRegistrado{}
}

// verificarProyecto compara nombre, fechas y estado de un proyecto.
func verificarProyecto(t *testing.T, obtenido, esperado domain.Proyecto) {
	t.Helper()
	if obtenido.Nombre != esperado.Nombre {
		t.Errorf("Nombre = %q, se esperaba %q", obtenido.Nombre, esperado.Nombre)
	}
	if !obtenido.FechaInicio.Equal(esperado.FechaInicio) {
		t.Errorf("FechaInicio = %v, se esperaba %v", obtenido.FechaInicio, esperado.FechaInicio)
	}
	if !obtenido.FechaFin.Equal(esperado.FechaFin) {
		t.Errorf("FechaFin = %v, se esperaba %v", obtenido.FechaFin, esperado.FechaFin)
	}
	if obtenido.Estado != esperado.Estado {
		t.Errorf("Estado = %q, se esperaba %q", obtenido.Estado, esperado.Estado)
	}
}

// verificarProyectoVacio comprueba que un error devuelve el valor cero de
// domain.Proyecto (D3, R6).
func verificarProyectoVacio(t *testing.T, proyecto domain.Proyecto) {
	t.Helper()
	if proyecto != (domain.Proyecto{}) {
		t.Errorf("se esperaba domain.Proyecto{} ante un error, se obtuvo %+v", proyecto)
	}
}

// CA1 / R1 / D17: una modificación válida persiste el nombre y las fechas
// nuevas en la colección "proyectos" (D13).
func TestUS02_Modificar_PersisteCambios(t *testing.T) {
	casos := []struct {
		descripcion string
		nombre      string
	}{
		{
			descripcion: "datos válidos",
			nombre:      "Sistema de Métricas v2",
		},
		{
			// C7: un nombre válido conserva los espacios externos tal como se ingresó.
			descripcion: "nombre con espacios externos",
			nombre:      "  Proyecto con espacios  ",
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			almacen, _ := nuevoAlmacenServicio(t)
			original := proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))
			sembrarProyectos(t, almacen, []proyectoRegistrado{
				{ID: "p1", Proyecto: original},
			})
			servicio := NuevoServicioProyectos(almacen)

			fechaInicio := fechaUS02(2026, time.March, 15)
			fechaFin := fechaUS02(2026, time.July, 31)
			modificado, err := servicio.ModificarProyecto("p1", caso.nombre, fechaInicio, fechaFin)
			if err != nil {
				t.Fatalf("ModificarProyecto() devolvió un error inesperado: %v", err)
			}

			esperado := domain.Proyecto{
				Nombre:      caso.nombre,
				FechaInicio: fechaInicio,
				FechaFin:    fechaFin,
				Estado:      domain.Estado("En curso"),
			}
			verificarProyecto(t, modificado, esperado)

			// D13: la colección se recarga desde disco con el nombre lógico "proyectos".
			coleccion := cargarProyectos(t, almacen)
			registro := buscarRegistro(t, coleccion, "p1")
			verificarProyecto(t, registro.Proyecto, esperado)
		})
	}
}

// CA5 / R5 / D9: la modificación conserva el estado en el proyecto devuelto y
// en el registro persistido, sin resetearlo a Planificado.
func TestUS02_Modificar_ConservaEstado(t *testing.T) {
	casos := []struct {
		descripcion string
		estado      domain.Estado
	}{
		{
			descripcion: "En curso",
			estado:      domain.Estado("En curso"),
		},
		{
			descripcion: "Finalizado",
			estado:      domain.Estado("Finalizado"),
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			almacen, _ := nuevoAlmacenServicio(t)
			original := proyectoUS02Servicio("Sistema de Métricas", caso.estado)
			sembrarProyectos(t, almacen, []proyectoRegistrado{
				{ID: "p1", Proyecto: original},
			})
			servicio := NuevoServicioProyectos(almacen)

			modificado, err := servicio.ModificarProyecto(
				"p1",
				"Sistema de Métricas v2",
				fechaUS02(2026, time.March, 15),
				fechaUS02(2026, time.July, 31),
			)
			if err != nil {
				t.Fatalf("ModificarProyecto() devolvió un error inesperado: %v", err)
			}
			if modificado.Estado != caso.estado {
				t.Errorf("Estado = %q, se esperaba conservar %q", modificado.Estado, caso.estado)
			}
			if modificado.Estado == domain.EstadoPlanificado {
				t.Errorf("Estado = %q, no debería resetearse a Planificado", modificado.Estado)
			}

			coleccion := cargarProyectos(t, almacen)
			registro := buscarRegistro(t, coleccion, "p1")
			if registro.Proyecto.Estado != caso.estado {
				t.Errorf("Estado persistido = %q, se esperaba conservar %q", registro.Proyecto.Estado, caso.estado)
			}
		})
	}
}

// D13: modificar un proyecto no altera los demás registros ni la cantidad de
// registros de la colección.
func TestUS02_Modificar_NoModificaOtrosProyectos(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	p1 := proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))
	p2 := domain.Proyecto{
		Nombre:      "Otro proyecto",
		FechaInicio: fechaUS02(2026, time.January, 15),
		FechaFin:    fechaUS02(2026, time.December, 31),
		Estado:      domain.EstadoPlanificado,
	}
	p3 := domain.Proyecto{
		Nombre:      "Tercer proyecto",
		FechaInicio: fechaUS02(2026, time.February, 10),
		FechaFin:    fechaUS02(2026, time.May, 20),
		Estado:      domain.Estado("Finalizado"),
	}
	sembrarProyectos(t, almacen, []proyectoRegistrado{
		{ID: "p1", Proyecto: p1},
		{ID: "p2", Proyecto: p2},
		{ID: "p3", Proyecto: p3},
	})
	servicio := NuevoServicioProyectos(almacen)

	_, err := servicio.ModificarProyecto(
		"p1",
		"Sistema de Métricas v2",
		fechaUS02(2026, time.March, 15),
		fechaUS02(2026, time.July, 31),
	)
	if err != nil {
		t.Fatalf("ModificarProyecto() devolvió un error inesperado: %v", err)
	}

	coleccion := cargarProyectos(t, almacen)
	if len(coleccion) != 3 {
		t.Fatalf("la cantidad de registros cambió: se esperaban 3, hay %d", len(coleccion))
	}

	esperadoP1 := domain.Proyecto{
		Nombre:      "Sistema de Métricas v2",
		FechaInicio: fechaUS02(2026, time.March, 15),
		FechaFin:    fechaUS02(2026, time.July, 31),
		Estado:      domain.Estado("En curso"),
	}
	verificarProyecto(t, buscarRegistro(t, coleccion, "p1").Proyecto, esperadoP1)
	verificarProyecto(t, buscarRegistro(t, coleccion, "p2").Proyecto, p2)
	verificarProyecto(t, buscarRegistro(t, coleccion, "p3").Proyecto, p3)
}

// CA1 / R7 / D10: modificar con los mismos valores actuales es válido y el
// registro persistido queda igual al original.
func TestUS02_Modificar_ValoresIguales(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	original := proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))
	sembrarProyectos(t, almacen, []proyectoRegistrado{
		{ID: "p1", Proyecto: original},
	})
	servicio := NuevoServicioProyectos(almacen)

	modificado, err := servicio.ModificarProyecto("p1", original.Nombre, original.FechaInicio, original.FechaFin)
	if err != nil {
		t.Fatalf("ModificarProyecto() con los mismos valores devolvió un error inesperado: %v", err)
	}
	verificarProyecto(t, modificado, original)

	coleccion := cargarProyectos(t, almacen)
	registro := buscarRegistro(t, coleccion, "p1")
	verificarProyecto(t, registro.Proyecto, original)
}

// almacenCuentaSaves envuelve un JSONStore real y cuenta las llamadas a Save,
// para verificar que el servicio no escribe cuando no corresponde.
type almacenCuentaSaves struct {
	real          *storage.JSONStore
	llamadasASave int
}

func (a *almacenCuentaSaves) Save(name string, v any) error {
	a.llamadasASave++
	return a.real.Save(name, v)
}

func (a *almacenCuentaSaves) Load(name string, v any) error {
	return a.real.Load(name, v)
}

// CA2 / CA3 / R6 / D3 / D14: ante un error de validación se devuelve Proyecto{},
// el error de dominio correspondiente, Save no se invoca (no hay escritura) y
// el archivo conserva los valores originales.
func TestUS02_Modificar_ValidacionFallidaNoPersiste(t *testing.T) {
	casos := []struct {
		descripcion   string
		nombre        string
		fechaInicio   time.Time
		fechaFin      time.Time
		errorEsperado error
	}{
		{
			descripcion:   "nombre vacío",
			nombre:        "",
			fechaInicio:   fechaUS02(2026, time.March, 1),
			fechaFin:      fechaUS02(2026, time.June, 30),
			errorEsperado: domain.ErrNombreVacio,
		},
		{
			// C1: un nombre compuesto únicamente por espacios se considera vacío.
			descripcion:   "nombre compuesto solo por espacios",
			nombre:        "   ",
			fechaInicio:   fechaUS02(2026, time.March, 1),
			fechaFin:      fechaUS02(2026, time.June, 30),
			errorEsperado: domain.ErrNombreVacio,
		},
		{
			descripcion:   "fecha de fin anterior al inicio",
			nombre:        "Sistema de Métricas v2",
			fechaInicio:   fechaUS02(2026, time.June, 30),
			fechaFin:      fechaUS02(2026, time.March, 1),
			errorEsperado: domain.ErrFechaFinAnteriorAInicio,
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			almacenReal, _ := nuevoAlmacenServicio(t)
			almacen := &almacenCuentaSaves{real: almacenReal}
			original := proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))
			// La siembra usa el almacenamiento real: el contador mide solo
			// las escrituras que hace el servicio.
			sembrarProyectos(t, almacenReal, []proyectoRegistrado{
				{ID: "p1", Proyecto: original},
			})
			servicio := NuevoServicioProyectos(almacen)

			modificado, err := servicio.ModificarProyecto("p1", caso.nombre, caso.fechaInicio, caso.fechaFin)
			if err == nil {
				t.Fatal("ModificarProyecto() con datos inválidos debería devolver un error")
			}
			if !errors.Is(err, caso.errorEsperado) {
				t.Errorf("err = %v, se esperaba %v", err, caso.errorEsperado)
			}
			verificarProyectoVacio(t, modificado)

			// D14: ante un error de validación el servicio no debe escribir.
			if almacen.llamadasASave != 0 {
				t.Errorf("Save fue invocado %d veces ante un error de validación; no debía escribirse nada", almacen.llamadasASave)
			}

			coleccion := cargarProyectos(t, almacenReal)
			registro := buscarRegistro(t, coleccion, "p1")
			verificarProyecto(t, registro.Proyecto, original)
		})
	}
}

// CA4 / R4 / V3 / D12: un ID ausente en una colección con otros proyectos
// devuelve ErrProyectoNoEncontrado y Proyecto{}.
func TestUS02_Modificar_ProyectoNoEncontrado(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
		{ID: "p2", Proyecto: proyectoUS02Servicio("Otro proyecto", domain.EstadoPlanificado)},
	})
	servicio := NuevoServicioProyectos(almacen)

	modificado, err := servicio.ModificarProyecto(
		"ausente",
		"Proyecto actualizado",
		fechaUS02(2026, time.March, 1),
		fechaUS02(2026, time.June, 30),
	)
	if err == nil {
		t.Fatal("ModificarProyecto() con un ID inexistente debería devolver un error")
	}
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	verificarProyectoVacio(t, modificado)
}

// CA4 / D4 / D7: la verificación de existencia precede a las validaciones de
// nombre y fechas.
func TestUS02_Modificar_ProyectoNoEncontradoPrecedeValidaciones(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	sembrarProyectos(t, almacen, []proyectoRegistrado{
		{ID: "p1", Proyecto: proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))},
	})
	servicio := NuevoServicioProyectos(almacen)

	modificado, err := servicio.ModificarProyecto(
		"ausente",
		"",
		fechaUS02(2026, time.June, 30),
		fechaUS02(2026, time.March, 1),
	)
	if err == nil {
		t.Fatal("ModificarProyecto() con un ID inexistente debería devolver un error")
	}
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if errors.Is(err, domain.ErrNombreVacio) {
		t.Errorf("err = %v, no debería devolverse ErrNombreVacio si el proyecto no existe", err)
	}
	if errors.Is(err, domain.ErrFechaFinAnteriorAInicio) {
		t.Errorf("err = %v, no debería devolverse ErrFechaFinAnteriorAInicio si el proyecto no existe", err)
	}
	verificarProyectoVacio(t, modificado)
}

// CA4 / D14: sin archivo de colección se interpreta como colección vacía y no
// se crea ningún archivo en ningún camino de error.
func TestUS02_Modificar_ArchivoInexistente(t *testing.T) {
	almacen, dir := nuevoAlmacenServicio(t)
	servicio := NuevoServicioProyectos(almacen)

	modificado, err := servicio.ModificarProyecto(
		"p1",
		"Proyecto actualizado",
		fechaUS02(2026, time.March, 1),
		fechaUS02(2026, time.June, 30),
	)
	if err == nil {
		t.Fatal("ModificarProyecto() sin archivo de colección debería devolver un error")
	}
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	verificarProyectoVacio(t, modificado)

	entradas, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir(%q): %v", dir, err)
	}
	if len(entradas) != 0 {
		t.Errorf("en ningún camino de error se crea el archivo, hay %d entradas", len(entradas))
	}
}

// E5 / D15: un JSON corrupto propaga storage.ErrCorrupt y no se mapea a
// ErrProyectoNoEncontrado.
func TestUS02_Modificar_JSONCorruptoNoEsNoEncontrado(t *testing.T) {
	almacen, dir := nuevoAlmacenServicio(t)
	// Mismo criterio que TestLoad_CorruptJSON de internal/storage.
	if err := os.WriteFile(filepath.Join(dir, "proyectos.json"), []byte("{no es json"), 0o644); err != nil {
		t.Fatal(err)
	}
	servicio := NuevoServicioProyectos(almacen)

	modificado, err := servicio.ModificarProyecto(
		"p1",
		"Proyecto actualizado",
		fechaUS02(2026, time.March, 1),
		fechaUS02(2026, time.June, 30),
	)
	if err == nil {
		t.Fatal("ModificarProyecto() con JSON corrupto debería devolver un error")
	}
	if !errors.Is(err, storage.ErrCorrupt) {
		t.Errorf("err = %v, se esperaba un error que envuelva %v", err, storage.ErrCorrupt)
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, no debería mapearse a ErrProyectoNoEncontrado", err)
	}
	verificarProyectoVacio(t, modificado)
}

// almacenFalloEscritura envuelve un JSONStore real: Load delega en él y Save
// siempre falla con el error de prueba.
type almacenFalloEscritura struct {
	real *storage.JSONStore
	err  error
}

func (a almacenFalloEscritura) Save(name string, v any) error {
	return a.err
}

func (a almacenFalloEscritura) Load(name string, v any) error {
	return a.real.Load(name, v)
}

// E5 / D15: un fallo de escritura se propaga envuelto y no se informa una
// modificación exitosa.
func TestUS02_Modificar_FalloDeEscrituraNoInformaExito(t *testing.T) {
	almacen, _ := nuevoAlmacenServicio(t)
	original := proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))
	sembrarProyectos(t, almacen, []proyectoRegistrado{
		{ID: "p1", Proyecto: original},
	})

	errSimulado := errors.New("fallo simulado de escritura")
	servicio := NuevoServicioProyectos(almacenFalloEscritura{real: almacen, err: errSimulado})

	modificado, err := servicio.ModificarProyecto(
		"p1",
		"Sistema de Métricas v2",
		fechaUS02(2026, time.March, 15),
		fechaUS02(2026, time.July, 31),
	)
	if err == nil {
		t.Fatal("ModificarProyecto() con fallo de escritura debería devolver un error")
	}
	if !errors.Is(err, errSimulado) {
		t.Errorf("err = %v, se esperaba un error que envuelva %v", err, errSimulado)
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, no debería ser ErrProyectoNoEncontrado", err)
	}
	if errors.Is(err, domain.ErrNombreVacio) || errors.Is(err, domain.ErrFechaFinAnteriorAInicio) {
		t.Errorf("err = %v, no debería ser un error de validación de dominio", err)
	}
	verificarProyectoVacio(t, modificado)
}
