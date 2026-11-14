package service

import (
	"errors"
	"sort"
	"strings"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// asignacionUS08 representa una fila del doble de asignaciones.
type asignacionUS08 struct {
	proyectoID   string
	sprintNumero int
	itemNumero   int
}

// fakeRepositorioAsignaciones es el doble en memoria de RepositorioAsignaciones
// para los tests de servicio de US-08. Cuenta las llamadas a Agregar y permite
// inyectar errores en cada operación.
type fakeRepositorioAsignaciones struct {
	asignaciones       []asignacionUS08
	agregarErr         error
	existeItemErr      error
	listarErr          error
	agregarLlamadas    int
	existeItemLlamadas int
	listarLlamadas     int
}

func nuevoFakeAsignaciones() *fakeRepositorioAsignaciones {
	return &fakeRepositorioAsignaciones{asignaciones: []asignacionUS08{}}
}

func (f *fakeRepositorioAsignaciones) Agregar(proyectoID string, sprintNumero, itemNumero int) error {
	f.agregarLlamadas++
	if f.agregarErr != nil {
		return f.agregarErr
	}
	for _, asignacion := range f.asignaciones {
		if asignacion.proyectoID == proyectoID && asignacion.itemNumero == itemNumero {
			return domain.ErrItemYaAsignado
		}
	}
	f.asignaciones = append(f.asignaciones, asignacionUS08{
		proyectoID:   proyectoID,
		sprintNumero: sprintNumero,
		itemNumero:   itemNumero,
	})
	return nil
}

func (f *fakeRepositorioAsignaciones) ExisteItem(proyectoID string, itemNumero int) (bool, error) {
	f.existeItemLlamadas++
	if f.existeItemErr != nil {
		return false, f.existeItemErr
	}
	_, ok := f.sprintDe(proyectoID, itemNumero)
	return ok, nil
}

func (f *fakeRepositorioAsignaciones) ListarNumerosItem(proyectoID string, sprintNumero int) ([]int, error) {
	f.listarLlamadas++
	if f.listarErr != nil {
		return nil, f.listarErr
	}
	numeros := make([]int, 0)
	for _, asignacion := range f.asignaciones {
		if asignacion.proyectoID == proyectoID && asignacion.sprintNumero == sprintNumero {
			numeros = append(numeros, asignacion.itemNumero)
		}
	}
	sort.Ints(numeros)
	return numeros, nil
}

// sembrar agrega una asignación directamente al doble sin contar la llamada a
// Agregar ni pasar por sus validaciones.
func (f *fakeRepositorioAsignaciones) sembrar(proyectoID string, sprintNumero, itemNumero int) {
	f.asignaciones = append(f.asignaciones, asignacionUS08{
		proyectoID:   proyectoID,
		sprintNumero: sprintNumero,
		itemNumero:   itemNumero,
	})
}

func (f *fakeRepositorioAsignaciones) contiene(proyectoID string, sprintNumero, itemNumero int) bool {
	for _, asignacion := range f.asignaciones {
		if asignacion.proyectoID == proyectoID &&
			asignacion.sprintNumero == sprintNumero &&
			asignacion.itemNumero == itemNumero {
			return true
		}
	}
	return false
}

func (f *fakeRepositorioAsignaciones) sprintDe(proyectoID string, itemNumero int) (int, bool) {
	for _, asignacion := range f.asignaciones {
		if asignacion.proyectoID == proyectoID && asignacion.itemNumero == itemNumero {
			return asignacion.sprintNumero, true
		}
	}
	return 0, false
}

func (f *fakeRepositorioAsignaciones) total() int {
	return len(f.asignaciones)
}

// fakeBacklogUS08 envuelve al doble de backlog de US-05/US-11 para contar las
// llamadas a ListarPorProyecto sin modificar repos_fake_test.go.
type fakeBacklogUS08 struct {
	*fakeRepositorioBacklog
	listarLlamadas int
}

func nuevoFakeBacklogUS08() *fakeBacklogUS08 {
	return &fakeBacklogUS08{fakeRepositorioBacklog: nuevoFakeBacklog()}
}

func (f *fakeBacklogUS08) ListarPorProyecto(proyectoID string) ([]domain.ItemBacklog, error) {
	f.listarLlamadas++
	return f.fakeRepositorioBacklog.ListarPorProyecto(proyectoID)
}

// entornoAsignaciones agrupa los dobles de proyectos, Sprints, backlog y
// asignaciones para sembrar datos y verificar el estado sin depender de
// internal/storage.
type entornoAsignaciones struct {
	proyectos    *fakeRepositorioProyectos
	sprints      *fakeRepositorioSprints
	backlog      *fakeBacklogUS08
	asignaciones *fakeRepositorioAsignaciones
}

func nuevoEntornoAsignaciones() *entornoAsignaciones {
	return &entornoAsignaciones{
		proyectos:    nuevoFakeProyectos(),
		sprints:      nuevoFakeSprints(),
		backlog:      nuevoFakeBacklogUS08(),
		asignaciones: nuevoFakeAsignaciones(),
	}
}

func (e *entornoAsignaciones) servicio() *ServicioSprints {
	return NuevoServicioSprints(e.proyectos, e.sprints, e.backlog, e.asignaciones)
}

// itemUS08Servicio arma una historia del backlog para los tests de servicio de
// US-08.
func itemUS08Servicio(id string, storyPoints int, estado string) domain.ItemBacklog {
	return domain.ItemBacklog{
		ID:          id,
		Titulo:      "Historia " + id,
		Descripcion: "Historia de prueba de US-08",
		Prioridad:   "Alta",
		Estado:      estado,
		StoryPoints: storyPoints,
	}
}

func sembrarProyectosUS08(t *testing.T, entorno *entornoAsignaciones, ids ...string) {
	t.Helper()
	for _, id := range ids {
		entorno.proyectos.proyectos[id] = proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))
	}
}

func sembrarSprintUS08(t *testing.T, entorno *entornoAsignaciones, proyectoID string, numero int) {
	t.Helper()
	entorno.sprints.sprints[proyectoID] = append(
		entorno.sprints.sprints[proyectoID],
		sprintUS07Servicio(numero, domain.EstadoSprintActivo),
	)
}

func sembrarItemUS08(t *testing.T, entorno *entornoAsignaciones, proyectoID string, item domain.ItemBacklog) {
	t.Helper()
	entorno.backlog.items[proyectoID] = append(entorno.backlog.items[proyectoID], item)
}

// entornoConHistoriaUS08 siembra el proyecto "p1", su Sprint número 1 y la
// historia indicada.
func entornoConHistoriaUS08(t *testing.T, item domain.ItemBacklog) (*entornoAsignaciones, *ServicioSprints) {
	t.Helper()
	entorno := nuevoEntornoAsignaciones()
	sembrarProyectosUS08(t, entorno, "p1")
	sembrarSprintUS08(t, entorno, "p1", 1)
	sembrarItemUS08(t, entorno, "p1", item)
	return entorno, entorno.servicio()
}

// verificarErrorRepositorioUS08 comprueba que el error de repositorio se
// propague envuelto y que no se mapee a un error de dominio de no encontrado.
func verificarErrorRepositorioUS08(t *testing.T, err, causa error) {
	t.Helper()
	if err == nil {
		t.Fatal("se esperaba un error de repositorio y se obtuvo nil")
	}
	if !errors.Is(err, causa) {
		t.Fatalf("errors.Is(err, causa) = false; err = %v", err)
	}
	if err == causa {
		t.Fatalf("el error debe estar envuelto; err == causa")
	}
	if !strings.HasPrefix(err.Error(), "service:") {
		t.Errorf("mensaje = %q; se esperaba prefijo service:", err.Error())
	}
	for _, errorDominio := range []error{
		domain.ErrProyectoNoEncontrado,
		domain.ErrSprintNoEncontrado,
		domain.ErrItemNoEncontrado,
	} {
		if errors.Is(err, errorDominio) {
			t.Errorf("err = %v, no debería mapearse a %v", err, errorDominio)
		}
	}
}

// CA1 / R1: una historia estimada y pendiente se asigna y se persiste una sola
// vez.
func TestUS08_AsignarItem_Exitoso(t *testing.T) {
	entorno, servicio := entornoConHistoriaUS08(t, itemUS08Servicio("HU-01", 3, domain.EstadoItemPendiente))

	err := servicio.AsignarItem("p1", 1, "HU-01")
	if err != nil {
		t.Fatalf("AsignarItem() devolvió un error inesperado: %v", err)
	}
	if entorno.asignaciones.agregarLlamadas != 1 {
		t.Errorf("Agregar fue invocado %d veces, se esperaba 1", entorno.asignaciones.agregarLlamadas)
	}
	if !entorno.asignaciones.contiene("p1", 1, 1) {
		t.Errorf("la historia HU-01 no quedó asignada al Sprint 1 del proyecto p1: %+v", entorno.asignaciones.asignaciones)
	}
	if entorno.asignaciones.total() != 1 {
		t.Errorf("se esperaba 1 asignación persistida, hay %d", entorno.asignaciones.total())
	}
}

// CA1 / R2: una historia sin estimar (story points 0 o 4) no se asigna.
func TestUS08_AsignarItem_NoEstimadoNoAsigna(t *testing.T) {
	casos := []struct {
		nombre      string
		storyPoints int
	}{
		{nombre: "story points 0", storyPoints: 0},
		{nombre: "story points 4 fuera de la serie", storyPoints: 4},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			entorno, servicio := entornoConHistoriaUS08(
				t,
				itemUS08Servicio("HU-04", caso.storyPoints, domain.EstadoItemPendiente),
			)

			err := servicio.AsignarItem("p1", 1, "HU-04")
			if !errors.Is(err, domain.ErrItemNoEstimado) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrItemNoEstimado)
			}
			if errors.Is(err, domain.ErrItemNoPendiente) {
				t.Errorf("err = %v, no debería devolverse ErrItemNoPendiente", err)
			}
			if entorno.asignaciones.agregarLlamadas != 0 {
				t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.asignaciones.agregarLlamadas)
			}
			if entorno.asignaciones.total() != 0 {
				t.Errorf("se persistió una asignación pese a la historia sin estimar: %+v", entorno.asignaciones.asignaciones)
			}
		})
	}
}

// CA1 / R3 / R13: una historia que no está Pendiente no se asigna.
func TestUS08_AsignarItem_NoPendienteNoAsigna(t *testing.T) {
	entorno, servicio := entornoConHistoriaUS08(t, itemUS08Servicio("HU-05", 5, "Completada"))

	err := servicio.AsignarItem("p1", 1, "HU-05")
	if !errors.Is(err, domain.ErrItemNoPendiente) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrItemNoPendiente)
	}
	if errors.Is(err, domain.ErrItemNoEstimado) {
		t.Errorf("err = %v, no debería devolverse ErrItemNoEstimado", err)
	}
	if entorno.asignaciones.agregarLlamadas != 0 {
		t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.asignaciones.agregarLlamadas)
	}
	if entorno.asignaciones.total() != 0 {
		t.Errorf("se persistió una asignación pese a la historia no pendiente: %+v", entorno.asignaciones.asignaciones)
	}
}

// CA1 / R6: un proyecto inexistente devuelve ErrProyectoNoEncontrado y no
// evalúa ninguna otra regla.
func TestUS08_AsignarItem_ProyectoNoEncontrado(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	servicio := entorno.servicio()

	err := servicio.AsignarItem("p999", 1, "HU-01")
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if entorno.sprints.existeNumeroLlamadas != 0 {
		t.Errorf("ExisteNumero fue invocado %d veces; el proyecto debía validarse primero", entorno.sprints.existeNumeroLlamadas)
	}
	if entorno.backlog.obtenerLlamadas != 0 {
		t.Errorf("ObtenerItem fue invocado %d veces; el proyecto debía validarse primero", entorno.backlog.obtenerLlamadas)
	}
	if entorno.asignaciones.existeItemLlamadas != 0 {
		t.Errorf("ExisteItem fue invocado %d veces; el proyecto debía validarse primero", entorno.asignaciones.existeItemLlamadas)
	}
	if entorno.asignaciones.agregarLlamadas != 0 {
		t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.asignaciones.agregarLlamadas)
	}
	if entorno.asignaciones.total() != 0 {
		t.Errorf("se persistió una asignación con el proyecto inexistente: %+v", entorno.asignaciones.asignaciones)
	}
}

// CA1 / R7: un Sprint inexistente en el proyecto devuelve
// ErrSprintNoEncontrado y no evalúa las reglas del ítem.
func TestUS08_AsignarItem_SprintNoEncontrado(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	sembrarProyectosUS08(t, entorno, "p1")
	sembrarSprintUS08(t, entorno, "p1", 1)
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-01", 3, domain.EstadoItemPendiente))
	servicio := entorno.servicio()

	err := servicio.AsignarItem("p1", 3, "HU-01")
	if !errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrSprintNoEncontrado)
	}
	if entorno.backlog.obtenerLlamadas != 0 {
		t.Errorf("ObtenerItem fue invocado %d veces; el Sprint debía validarse primero", entorno.backlog.obtenerLlamadas)
	}
	if entorno.asignaciones.existeItemLlamadas != 0 {
		t.Errorf("ExisteItem fue invocado %d veces; el Sprint debía validarse primero", entorno.asignaciones.existeItemLlamadas)
	}
	if entorno.asignaciones.agregarLlamadas != 0 {
		t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.asignaciones.agregarLlamadas)
	}
	if entorno.asignaciones.total() != 0 {
		t.Errorf("se persistió una asignación con el Sprint inexistente: %+v", entorno.asignaciones.asignaciones)
	}
}

// CA1 / R8: una historia inexistente devuelve ErrItemNoEncontrado.
func TestUS08_AsignarItem_ItemNoEncontrado(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	sembrarProyectosUS08(t, entorno, "p1")
	sembrarSprintUS08(t, entorno, "p1", 1)
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-01", 3, domain.EstadoItemPendiente))
	servicio := entorno.servicio()

	err := servicio.AsignarItem("p1", 1, "HU-99")
	if !errors.Is(err, domain.ErrItemNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrItemNoEncontrado)
	}
	if entorno.asignaciones.existeItemLlamadas != 0 {
		t.Errorf("ExisteItem fue invocado %d veces; el ítem debía validarse primero", entorno.asignaciones.existeItemLlamadas)
	}
	if entorno.asignaciones.agregarLlamadas != 0 {
		t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.asignaciones.agregarLlamadas)
	}
	if entorno.asignaciones.total() != 0 {
		t.Errorf("se persistió una asignación con el ítem inexistente: %+v", entorno.asignaciones.asignaciones)
	}
}

// CA1 / R8 / D6: un ID malformado se trata como historia inexistente.
func TestUS08_AsignarItem_IDMalformadoNoEncontrado(t *testing.T) {
	casos := []string{"HU-", "HU-X", "XX-01", "HU-00"}

	for _, itemID := range casos {
		t.Run(itemID, func(t *testing.T) {
			entorno := nuevoEntornoAsignaciones()
			sembrarProyectosUS08(t, entorno, "p1")
			sembrarSprintUS08(t, entorno, "p1", 1)
			sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-01", 3, domain.EstadoItemPendiente))
			servicio := entorno.servicio()

			err := servicio.AsignarItem("p1", 1, itemID)
			if !errors.Is(err, domain.ErrItemNoEncontrado) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrItemNoEncontrado)
			}
			if entorno.asignaciones.existeItemLlamadas != 0 {
				t.Errorf("ExisteItem fue invocado %d veces; el ítem debía validarse primero", entorno.asignaciones.existeItemLlamadas)
			}
			if entorno.asignaciones.agregarLlamadas != 0 {
				t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.asignaciones.agregarLlamadas)
			}
			if entorno.asignaciones.total() != 0 {
				t.Errorf("se persistió una asignación con el ID malformado %q: %+v", itemID, entorno.asignaciones.asignaciones)
			}
		})
	}
}

// CA1 / R6 / R11: el proyecto se evalúa antes que el Sprint.
func TestUS08_AsignarItem_OrdenProyectoAntesQueSprint(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	servicio := entorno.servicio()

	err := servicio.AsignarItem("p999", 9, "HU-01")
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Errorf("err = %v, no debería devolverse ErrSprintNoEncontrado si el proyecto no existe", err)
	}
	if entorno.sprints.existeNumeroLlamadas != 0 {
		t.Errorf("ExisteNumero fue invocado %d veces; el proyecto debía evaluarse primero", entorno.sprints.existeNumeroLlamadas)
	}
	if entorno.asignaciones.agregarLlamadas != 0 {
		t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.asignaciones.agregarLlamadas)
	}
	if entorno.asignaciones.total() != 0 {
		t.Errorf("se persistió una asignación pese al orden de evaluación: %+v", entorno.asignaciones.asignaciones)
	}
}

// CA1 / R7 / R11: el Sprint se evalúa antes que el ítem.
func TestUS08_AsignarItem_OrdenSprintAntesQueItem(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	sembrarProyectosUS08(t, entorno, "p1")
	sembrarSprintUS08(t, entorno, "p1", 1)
	servicio := entorno.servicio()

	err := servicio.AsignarItem("p1", 9, "HU-99")
	if !errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrSprintNoEncontrado)
	}
	if errors.Is(err, domain.ErrItemNoEncontrado) {
		t.Errorf("err = %v, no debería devolverse ErrItemNoEncontrado si el Sprint no existe", err)
	}
	if entorno.backlog.obtenerLlamadas != 0 {
		t.Errorf("ObtenerItem fue invocado %d veces; el Sprint debía evaluarse primero", entorno.backlog.obtenerLlamadas)
	}
	if entorno.asignaciones.agregarLlamadas != 0 {
		t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.asignaciones.agregarLlamadas)
	}
	if entorno.asignaciones.total() != 0 {
		t.Errorf("se persistió una asignación pese al orden de evaluación: %+v", entorno.asignaciones.asignaciones)
	}
}

// CA2 / R4: una historia ya asignada al mismo Sprint no se vuelve a asignar.
// Incluye el caso sin estimar para comprobar que "ya asignado" se evalúa antes
// que las validaciones de CA1.
func TestUS08_AsignarItem_YaAsignadoMismoSprint(t *testing.T) {
	casos := []struct {
		nombre      string
		storyPoints int
	}{
		{nombre: "historia ya asignada", storyPoints: 3},
		{nombre: "historia ya asignada y sin estimar", storyPoints: 0},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			entorno := nuevoEntornoAsignaciones()
			sembrarProyectosUS08(t, entorno, "p1")
			sembrarSprintUS08(t, entorno, "p1", 1)
			sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-06", caso.storyPoints, domain.EstadoItemPendiente))
			entorno.asignaciones.sembrar("p1", 1, 6)
			servicio := entorno.servicio()

			err := servicio.AsignarItem("p1", 1, "HU-06")
			if !errors.Is(err, domain.ErrItemYaAsignado) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrItemYaAsignado)
			}
			if errors.Is(err, domain.ErrItemNoEstimado) {
				t.Errorf("err = %v, no debería devolverse ErrItemNoEstimado antes que ErrItemYaAsignado", err)
			}
			if entorno.asignaciones.agregarLlamadas != 0 {
				t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.asignaciones.agregarLlamadas)
			}
			if entorno.asignaciones.total() != 1 {
				t.Errorf("se esperaba conservar 1 asignación, hay %d", entorno.asignaciones.total())
			}
			if !entorno.asignaciones.contiene("p1", 1, 6) {
				t.Errorf("la asignación existente de HU-06 al Sprint 1 se perdió: %+v", entorno.asignaciones.asignaciones)
			}
		})
	}
}

// CA2 / R4: una historia ya asignada a otro Sprint no se mueve ni se duplica.
func TestUS08_AsignarItem_YaAsignadoOtroSprint(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	sembrarProyectosUS08(t, entorno, "p1")
	sembrarSprintUS08(t, entorno, "p1", 1)
	sembrarSprintUS08(t, entorno, "p1", 2)
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-06", 3, domain.EstadoItemPendiente))
	entorno.asignaciones.sembrar("p1", 1, 6)
	servicio := entorno.servicio()

	err := servicio.AsignarItem("p1", 2, "HU-06")
	if !errors.Is(err, domain.ErrItemYaAsignado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrItemYaAsignado)
	}
	if entorno.asignaciones.agregarLlamadas != 0 {
		t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.asignaciones.agregarLlamadas)
	}
	if entorno.asignaciones.total() != 1 {
		t.Errorf("se esperaba conservar 1 asignación, hay %d", entorno.asignaciones.total())
	}
	if !entorno.asignaciones.contiene("p1", 1, 6) {
		t.Errorf("la asignación existente de HU-06 al Sprint 1 se perdió: %+v", entorno.asignaciones.asignaciones)
	}
	if entorno.asignaciones.contiene("p1", 2, 6) {
		t.Errorf("HU-06 se duplicó en el Sprint 2: %+v", entorno.asignaciones.asignaciones)
	}
}

// CA2 / R5: el mismo número de historia en otro proyecto no interfiere.
func TestUS08_AsignarItem_MismoNumeroOtroProyectoValido(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	sembrarProyectosUS08(t, entorno, "p1", "p2")
	sembrarSprintUS08(t, entorno, "p1", 1)
	sembrarSprintUS08(t, entorno, "p2", 1)
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-01", 3, domain.EstadoItemPendiente))
	sembrarItemUS08(t, entorno, "p2", itemUS08Servicio("HU-01", 8, domain.EstadoItemPendiente))
	entorno.asignaciones.sembrar("p2", 1, 1)
	servicio := entorno.servicio()

	err := servicio.AsignarItem("p1", 1, "HU-01")
	if err != nil {
		t.Fatalf("AsignarItem() devolvió un error inesperado: %v", err)
	}
	if entorno.asignaciones.agregarLlamadas != 1 {
		t.Errorf("Agregar fue invocado %d veces, se esperaba 1", entorno.asignaciones.agregarLlamadas)
	}
	if !entorno.asignaciones.contiene("p1", 1, 1) {
		t.Errorf("HU-01 no quedó asignada al Sprint 1 de p1: %+v", entorno.asignaciones.asignaciones)
	}
	if !entorno.asignaciones.contiene("p2", 1, 1) {
		t.Errorf("la asignación de HU-01 en p2 se perdió: %+v", entorno.asignaciones.asignaciones)
	}
	if entorno.asignaciones.total() != 2 {
		t.Errorf("se esperaban 2 asignaciones, hay %d", entorno.asignaciones.total())
	}
	if sprint, ok := entorno.asignaciones.sprintDe("p2", 1); !ok || sprint != 1 {
		t.Errorf("HU-01 de p2 = Sprint %d (ok=%v), se esperaba conservar el Sprint 1", sprint, ok)
	}
}

// CA3 / R9: el total se actualiza al asignar una nueva historia. Usa IDs sin
// relleno de ceros ("HU-1", "HU-10") para no asumir el formato.
func TestUS08_StoryPointsPlanificados_TrasAsignar(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	sembrarProyectosUS08(t, entorno, "p1")
	sembrarSprintUS08(t, entorno, "p1", 1)
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-1", 3, domain.EstadoItemPendiente))
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-10", 2, domain.EstadoItemPendiente))
	entorno.asignaciones.sembrar("p1", 1, 1)
	servicio := entorno.servicio()

	if err := servicio.AsignarItem("p1", 1, "HU-10"); err != nil {
		t.Fatalf("AsignarItem() devolvió un error inesperado: %v", err)
	}

	total, err := servicio.StoryPointsPlanificados("p1", 1)
	if err != nil {
		t.Fatalf("StoryPointsPlanificados() devolvió un error inesperado: %v", err)
	}
	if total != 5 {
		t.Errorf("StoryPointsPlanificados() = %d, se esperaba 5", total)
	}
	if entorno.backlog.listarLlamadas != 1 {
		t.Errorf("ListarPorProyecto fue invocado %d veces, se esperaba 1", entorno.backlog.listarLlamadas)
	}
	if entorno.asignaciones.listarLlamadas != 1 {
		t.Errorf("ListarNumerosItem fue invocado %d veces, se esperaba 1", entorno.asignaciones.listarLlamadas)
	}
}

// CA3 / R9: un Sprint sin historias tiene 0 Story Points planificados.
func TestUS08_StoryPointsPlanificados_SinItems(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	sembrarProyectosUS08(t, entorno, "p1")
	sembrarSprintUS08(t, entorno, "p1", 1)
	servicio := entorno.servicio()

	total, err := servicio.StoryPointsPlanificados("p1", 1)
	if err != nil {
		t.Fatalf("StoryPointsPlanificados() devolvió un error inesperado: %v", err)
	}
	if total != 0 {
		t.Errorf("StoryPointsPlanificados() = %d, se esperaba 0", total)
	}
	if entorno.asignaciones.listarLlamadas != 1 {
		t.Errorf("ListarNumerosItem fue invocado %d veces, se esperaba 1", entorno.asignaciones.listarLlamadas)
	}
	if entorno.backlog.listarLlamadas != 1 {
		t.Errorf("ListarPorProyecto fue invocado %d veces, se esperaba 1", entorno.backlog.listarLlamadas)
	}
}

// CA1 / CA2 / CA3 / R10: los errores de repositorio se propagan envueltos y no
// se mapean a errores de dominio de no encontrado.
func TestUS08_AsignarItem_ErrorDeRepositorioEnvuelto(t *testing.T) {
	causa := errors.New("fallo de repositorio")

	casos := []struct {
		descripcion     string
		preparar        func(t *testing.T, entorno *entornoAsignaciones, causa error)
		agregarLlamadas int
	}{
		{
			descripcion: "Existe",
			preparar: func(t *testing.T, entorno *entornoAsignaciones, causa error) {
				entorno.proyectos.existeErr = causa
			},
		},
		{
			descripcion: "ExisteNumero",
			preparar: func(t *testing.T, entorno *entornoAsignaciones, causa error) {
				sembrarProyectosUS08(t, entorno, "p1")
				entorno.sprints.existeNumeroErr = causa
			},
		},
		{
			descripcion: "ObtenerItem",
			preparar: func(t *testing.T, entorno *entornoAsignaciones, causa error) {
				sembrarProyectosUS08(t, entorno, "p1")
				sembrarSprintUS08(t, entorno, "p1", 1)
				entorno.backlog.obtenerErr = causa
			},
		},
		{
			descripcion: "ExisteItem",
			preparar: func(t *testing.T, entorno *entornoAsignaciones, causa error) {
				sembrarProyectosUS08(t, entorno, "p1")
				sembrarSprintUS08(t, entorno, "p1", 1)
				sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-01", 3, domain.EstadoItemPendiente))
				entorno.asignaciones.existeItemErr = causa
			},
		},
		{
			descripcion:     "Agregar",
			agregarLlamadas: 1,
			preparar: func(t *testing.T, entorno *entornoAsignaciones, causa error) {
				sembrarProyectosUS08(t, entorno, "p1")
				sembrarSprintUS08(t, entorno, "p1", 1)
				sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-01", 3, domain.EstadoItemPendiente))
				entorno.asignaciones.agregarErr = causa
			},
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoAsignaciones()
			caso.preparar(t, entorno, causa)
			servicio := entorno.servicio()

			err := servicio.AsignarItem("p1", 1, "HU-01")
			verificarErrorRepositorioUS08(t, err, causa)

			if entorno.asignaciones.agregarLlamadas != caso.agregarLlamadas {
				t.Errorf("Agregar fue invocado %d veces, se esperaba %d", entorno.asignaciones.agregarLlamadas, caso.agregarLlamadas)
			}
			if entorno.asignaciones.total() != 0 {
				t.Errorf("se persistió una asignación pese al error de repositorio: %+v", entorno.asignaciones.asignaciones)
			}
		})
	}
}

// CA1 / CA2 / CA3 / R10: ningún error de repositorio previo a la escritura
// asigna la historia.
func TestUS08_AsignarItem_ErrorNoAsigna(t *testing.T) {
	causa := errors.New("fallo de repositorio")

	casos := []struct {
		descripcion string
		preparar    func(t *testing.T, entorno *entornoAsignaciones, causa error)
	}{
		{
			descripcion: "Existe",
			preparar: func(t *testing.T, entorno *entornoAsignaciones, causa error) {
				entorno.proyectos.existeErr = causa
			},
		},
		{
			descripcion: "ExisteNumero",
			preparar: func(t *testing.T, entorno *entornoAsignaciones, causa error) {
				sembrarProyectosUS08(t, entorno, "p1")
				entorno.sprints.existeNumeroErr = causa
			},
		},
		{
			descripcion: "ObtenerItem",
			preparar: func(t *testing.T, entorno *entornoAsignaciones, causa error) {
				sembrarProyectosUS08(t, entorno, "p1")
				sembrarSprintUS08(t, entorno, "p1", 1)
				entorno.backlog.obtenerErr = causa
			},
		},
		{
			descripcion: "ExisteItem",
			preparar: func(t *testing.T, entorno *entornoAsignaciones, causa error) {
				sembrarProyectosUS08(t, entorno, "p1")
				sembrarSprintUS08(t, entorno, "p1", 1)
				sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-01", 3, domain.EstadoItemPendiente))
				entorno.asignaciones.existeItemErr = causa
			},
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoAsignaciones()
			caso.preparar(t, entorno, causa)
			servicio := entorno.servicio()

			err := servicio.AsignarItem("p1", 1, "HU-01")
			if err == nil {
				t.Fatal("AsignarItem() debía devolver el error de repositorio")
			}
			if !errors.Is(err, causa) {
				t.Fatalf("errors.Is(err, causa) = false; err = %v", err)
			}
			if entorno.asignaciones.agregarLlamadas != 0 {
				t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.asignaciones.agregarLlamadas)
			}
			if entorno.asignaciones.total() != 0 {
				t.Errorf("se persistió una asignación pese al error de repositorio: %+v", entorno.asignaciones.asignaciones)
			}
		})
	}
}

// CA1 / CA3 / R10: los errores de repositorio de StoryPointsPlanificados se
// propagan envueltos y no devuelven un total.
func TestUS08_StoryPointsPlanificados_ErrorDeRepositorioEnvuelto(t *testing.T) {
	causa := errors.New("fallo de repositorio")

	casos := []struct {
		descripcion string
		preparar    func(t *testing.T, entorno *entornoAsignaciones, causa error)
	}{
		{
			descripcion: "Existe",
			preparar: func(t *testing.T, entorno *entornoAsignaciones, causa error) {
				entorno.proyectos.existeErr = causa
			},
		},
		{
			descripcion: "ExisteNumero",
			preparar: func(t *testing.T, entorno *entornoAsignaciones, causa error) {
				sembrarProyectosUS08(t, entorno, "p1")
				entorno.sprints.existeNumeroErr = causa
			},
		},
		{
			descripcion: "ListarNumerosItem",
			preparar: func(t *testing.T, entorno *entornoAsignaciones, causa error) {
				sembrarProyectosUS08(t, entorno, "p1")
				sembrarSprintUS08(t, entorno, "p1", 1)
				entorno.asignaciones.listarErr = causa
			},
		},
		{
			descripcion: "ListarPorProyecto",
			preparar: func(t *testing.T, entorno *entornoAsignaciones, causa error) {
				sembrarProyectosUS08(t, entorno, "p1")
				sembrarSprintUS08(t, entorno, "p1", 1)
				entorno.backlog.listarErr = causa
			},
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoAsignaciones()
			caso.preparar(t, entorno, causa)
			servicio := entorno.servicio()

			total, err := servicio.StoryPointsPlanificados("p1", 1)
			if total != 0 {
				t.Errorf("StoryPointsPlanificados() = %d, se esperaba 0 ante el error", total)
			}
			verificarErrorRepositorioUS08(t, err, causa)
		})
	}
}
