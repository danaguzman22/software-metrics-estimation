package service

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// fakeRepositorioSprints es el doble en memoria de RepositorioSprints para los
// tests de servicio de US-07 (D8). No depende de internal/storage.
type fakeRepositorioSprints struct {
	sprints              map[string][]domain.Sprint
	agregarErr           error
	existeNumeroErr      error
	existeActivoErr      error
	agregados            int
	existeNumeroLlamadas int
	existeActivoLlamadas int
	ultimoProyectoID     string
	ultimoSprint         domain.Sprint
}

func (f *fakeRepositorioSprints) Agregar(proyectoID string, sprint domain.Sprint) error {
	f.agregados++
	f.ultimoProyectoID = proyectoID
	f.ultimoSprint = sprint
	if f.agregarErr != nil {
		return f.agregarErr
	}
	f.sprints[proyectoID] = append(f.sprints[proyectoID], sprint)
	return nil
}

func (f *fakeRepositorioSprints) ExisteNumero(proyectoID string, numero int) (bool, error) {
	f.existeNumeroLlamadas++
	if f.existeNumeroErr != nil {
		return false, f.existeNumeroErr
	}
	for _, sprint := range f.sprints[proyectoID] {
		if sprint.Numero == numero {
			return true, nil
		}
	}
	return false, nil
}

func (f *fakeRepositorioSprints) ExisteActivo(proyectoID string) (bool, error) {
	f.existeActivoLlamadas++
	if f.existeActivoErr != nil {
		return false, f.existeActivoErr
	}
	for _, sprint := range f.sprints[proyectoID] {
		if sprint.Estado == domain.EstadoSprintActivo {
			return true, nil
		}
	}
	return false, nil
}

func nuevoFakeSprints() *fakeRepositorioSprints {
	return &fakeRepositorioSprints{sprints: map[string][]domain.Sprint{}}
}

// entornoSprints agrupa los fakes de proyectos y Sprints para sembrar datos y
// verificar el estado sin depender de internal/storage.
type entornoSprints struct {
	*fakeRepositorioProyectos
	*fakeRepositorioSprints
}

func nuevoEntornoSprints() *entornoSprints {
	return &entornoSprints{
		fakeRepositorioProyectos: nuevoFakeProyectos(),
		fakeRepositorioSprints:   nuevoFakeSprints(),
	}
}

func (e *entornoSprints) servicio() *ServicioSprints {
	return NuevoServicioSprints(e.fakeRepositorioProyectos, e.fakeRepositorioSprints)
}

// fechaUS07Servicio arma una fecha sin hora para los escenarios de US-07.
func fechaUS07Servicio(anio int, mes time.Month, dia int) time.Time {
	return time.Date(anio, mes, dia, 0, 0, 0, 0, time.UTC)
}

// sembrarProyectosSprint carga los proyectos indicados en el fake de proyectos.
func sembrarProyectosSprint(t *testing.T, entorno *entornoSprints, ids ...string) {
	t.Helper()
	for _, id := range ids {
		entorno.proyectos[id] = proyectoUS02Servicio("Sistema de Métricas", domain.Estado("En curso"))
	}
}

// sprintUS07Servicio arma un Sprint de referencia para sembrar el fake.
func sprintUS07Servicio(numero int, estado string) domain.Sprint {
	return domain.Sprint{
		Numero:      numero,
		Objetivo:    "Planificar la primera iteración",
		FechaInicio: fechaUS07Servicio(2026, time.March, 2),
		FechaFin:    fechaUS07Servicio(2026, time.March, 20),
		Estado:      estado,
	}
}

// verificarSprintServicioUS07 compara campo por campo un Sprint con su esperado.
func verificarSprintServicioUS07(t *testing.T, obtenido, esperado domain.Sprint) {
	t.Helper()
	if obtenido.Numero != esperado.Numero {
		t.Errorf("Numero = %d, se esperaba %d", obtenido.Numero, esperado.Numero)
	}
	if obtenido.Objetivo != esperado.Objetivo {
		t.Errorf("Objetivo = %q, se esperaba %q", obtenido.Objetivo, esperado.Objetivo)
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

// verificarSprintServicioUS07Vacio comprueba que un error devuelve el valor cero
// de domain.Sprint (sección 8 de la spec).
func verificarSprintServicioUS07Vacio(t *testing.T, sprint domain.Sprint) {
	t.Helper()
	if sprint.Numero != 0 || sprint.Objetivo != "" ||
		!sprint.FechaInicio.IsZero() || !sprint.FechaFin.IsZero() ||
		sprint.Estado != "" {
		t.Errorf("se esperaba domain.Sprint{} ante un error, se obtuvo %+v", sprint)
	}
}

// CA1 / R1: crear un Sprint válido lo devuelve con estado Activo y lo persiste.
func TestUS07_CrearSprint_Exitoso(t *testing.T) {
	entorno := nuevoEntornoSprints()
	sembrarProyectosSprint(t, entorno, "p1")
	servicio := entorno.servicio()

	inicio := fechaUS07Servicio(2026, time.March, 2)
	fin := fechaUS07Servicio(2026, time.March, 20)
	creado, err := servicio.CrearSprint("p1", 1, "Planificar la primera iteración", inicio, fin)
	if err != nil {
		t.Fatalf("CrearSprint() devolvió un error inesperado: %v", err)
	}

	esperado := domain.Sprint{
		Numero:      1,
		Objetivo:    "Planificar la primera iteración",
		FechaInicio: inicio,
		FechaFin:    fin,
		Estado:      "Activo",
	}
	verificarSprintServicioUS07(t, creado, esperado)
	if entorno.agregados != 1 {
		t.Errorf("Agregar fue invocado %d veces, se esperaba 1", entorno.agregados)
	}

	persistidos := entorno.sprints["p1"]
	if len(persistidos) != 1 {
		t.Fatalf("se esperaba 1 Sprint persistido, hay %d", len(persistidos))
	}
	verificarSprintServicioUS07(t, persistidos[0], esperado)
}

// CA1 / R2 / D3: el objetivo se persiste recortado.
func TestUS07_CrearSprint_ObjetivoRecortado(t *testing.T) {
	entorno := nuevoEntornoSprints()
	sembrarProyectosSprint(t, entorno, "p1")
	servicio := entorno.servicio()

	inicio := fechaUS07Servicio(2026, time.March, 2)
	fin := fechaUS07Servicio(2026, time.March, 20)
	creado, err := servicio.CrearSprint("p1", 1, "   Planificar la primera iteración   ", inicio, fin)
	if err != nil {
		t.Fatalf("CrearSprint() devolvió un error inesperado: %v", err)
	}
	if creado.Objetivo != "Planificar la primera iteración" {
		t.Errorf("Objetivo devuelto = %q, se esperaba %q", creado.Objetivo, "Planificar la primera iteración")
	}

	persistidos := entorno.sprints["p1"]
	if len(persistidos) != 1 {
		t.Fatalf("se esperaba 1 Sprint persistido, hay %d", len(persistidos))
	}
	if persistidos[0].Objetivo != "Planificar la primera iteración" {
		t.Errorf("Objetivo persistido = %q, se esperaba %q", persistidos[0].Objetivo, "Planificar la primera iteración")
	}
}

// CA1 / R3 / D4: un número menor o igual a cero no se persiste.
func TestUS07_CrearSprint_NumeroInvalido(t *testing.T) {
	for _, numero := range []int{0, -1} {
		t.Run(fmt.Sprintf("numero %d", numero), func(t *testing.T) {
			entorno := nuevoEntornoSprints()
			sembrarProyectosSprint(t, entorno, "p1")
			servicio := entorno.servicio()

			creado, err := servicio.CrearSprint(
				"p1",
				numero,
				"Planificar la primera iteración",
				fechaUS07Servicio(2026, time.March, 2),
				fechaUS07Servicio(2026, time.March, 20),
			)
			if !errors.Is(err, domain.ErrNumeroSprintInvalido) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrNumeroSprintInvalido)
			}
			verificarSprintServicioUS07Vacio(t, creado)
			if entorno.agregados != 0 {
				t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.agregados)
			}
			if len(entorno.sprints["p1"]) != 0 {
				t.Errorf("se persistió un Sprint pese al número inválido")
			}
		})
	}
}

// CA1 / R4: un número ya usado por un Sprint no activo del mismo proyecto
// devuelve ErrSprintDuplicado y no persiste.
func TestUS07_CrearSprint_NumeroDuplicadoNoPersiste(t *testing.T) {
	entorno := nuevoEntornoSprints()
	sembrarProyectosSprint(t, entorno, "p1")
	existente := sprintUS07Servicio(1, "Cerrado")
	entorno.sprints["p1"] = append(entorno.sprints["p1"], existente)
	servicio := entorno.servicio()

	creado, err := servicio.CrearSprint(
		"p1",
		1,
		"Planificar una nueva iteración",
		fechaUS07Servicio(2026, time.March, 2),
		fechaUS07Servicio(2026, time.March, 20),
	)
	if !errors.Is(err, domain.ErrSprintDuplicado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrSprintDuplicado)
	}
	verificarSprintServicioUS07Vacio(t, creado)
	if entorno.agregados != 0 {
		t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.agregados)
	}
	if len(entorno.sprints["p1"]) != 1 {
		t.Errorf("se esperaba conservar 1 Sprint con número 1, hay %d", len(entorno.sprints["p1"]))
	}
}

// CA1 / R4: el mismo número en otro proyecto es válido.
func TestUS07_CrearSprint_MismoNumeroEnOtroProyecto(t *testing.T) {
	entorno := nuevoEntornoSprints()
	sembrarProyectosSprint(t, entorno, "p1", "p2")
	existente := sprintUS07Servicio(1, "Activo")
	existente.Objetivo = "Planificar el portal"
	entorno.sprints["p1"] = append(entorno.sprints["p1"], existente)
	servicio := entorno.servicio()

	creado, err := servicio.CrearSprint(
		"p2",
		1,
		"Planificar la primera iteración",
		fechaUS07Servicio(2026, time.March, 2),
		fechaUS07Servicio(2026, time.March, 20),
	)
	if err != nil {
		t.Fatalf("CrearSprint() devolvió un error inesperado: %v", err)
	}
	if creado.Numero != 1 || creado.Estado != "Activo" {
		t.Errorf("Sprint creado = %+v, se esperaba número 1 y estado Activo", creado)
	}
	if len(entorno.sprints["p1"]) != 1 {
		t.Errorf("el proyecto p1 debería conservar 1 Sprint, hay %d", len(entorno.sprints["p1"]))
	}
	if len(entorno.sprints["p2"]) != 1 {
		t.Errorf("el proyecto p2 debería tener 1 Sprint, hay %d", len(entorno.sprints["p2"]))
	}
	if entorno.sprints["p1"][0].Objetivo != "Planificar el portal" {
		t.Errorf("el Sprint de p1 cambió su objetivo a %q", entorno.sprints["p1"][0].Objetivo)
	}
}

// CA1 / R5 / D5: una fecha de fin anterior a la de inicio no se persiste.
func TestUS07_CrearSprint_FechasInvalidasNoPersiste(t *testing.T) {
	entorno := nuevoEntornoSprints()
	sembrarProyectosSprint(t, entorno, "p1")
	servicio := entorno.servicio()

	creado, err := servicio.CrearSprint(
		"p1",
		1,
		"Planificar la primera iteración",
		fechaUS07Servicio(2026, time.March, 20),
		fechaUS07Servicio(2026, time.March, 2),
	)
	if !errors.Is(err, domain.ErrFechasSprintInvalidas) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrFechasSprintInvalidas)
	}
	verificarSprintServicioUS07Vacio(t, creado)
	if entorno.agregados != 0 {
		t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.agregados)
	}
	if len(entorno.sprints["p1"]) != 0 {
		t.Errorf("se persistió un Sprint pese a las fechas inválidas")
	}
}

// CA1 / R5 / D5: fecha de inicio igual a la de fin es válida.
func TestUS07_CrearSprint_FechasIgualesValidas(t *testing.T) {
	entorno := nuevoEntornoSprints()
	sembrarProyectosSprint(t, entorno, "p1")
	servicio := entorno.servicio()

	dia := fechaUS07Servicio(2026, time.March, 10)
	creado, err := servicio.CrearSprint("p1", 2, "Realizar una iteración de un día", dia, dia)
	if err != nil {
		t.Fatalf("CrearSprint() devolvió un error inesperado: %v", err)
	}
	if !creado.FechaInicio.Equal(dia) || !creado.FechaFin.Equal(dia) {
		t.Errorf("fechas = %v-%v, se esperaban ambas %v", creado.FechaInicio, creado.FechaFin, dia)
	}
	if len(entorno.sprints["p1"]) != 1 {
		t.Errorf("se esperaba 1 Sprint persistido, hay %d", len(entorno.sprints["p1"]))
	}
}

// CA1 / R5 / D5: dos fechas del mismo día con distinta hora son válidas.
func TestUS07_CrearSprint_MismoDiaDistintaHoraValido(t *testing.T) {
	entorno := nuevoEntornoSprints()
	sembrarProyectosSprint(t, entorno, "p1")
	servicio := entorno.servicio()

	inicio := time.Date(2026, time.March, 10, 18, 0, 0, 0, time.UTC)
	fin := time.Date(2026, time.March, 10, 9, 0, 0, 0, time.UTC)
	creado, err := servicio.CrearSprint("p1", 2, "Realizar una iteración de un día", inicio, fin)
	if err != nil {
		t.Fatalf("CrearSprint() devolvió un error inesperado: %v", err)
	}
	if !creado.FechaInicio.Equal(inicio) || !creado.FechaFin.Equal(fin) {
		t.Errorf("fechas = %v-%v, se esperaban %v-%v", creado.FechaInicio, creado.FechaFin, inicio, fin)
	}
	if len(entorno.sprints["p1"]) != 1 {
		t.Errorf("se esperaba 1 Sprint persistido, hay %d", len(entorno.sprints["p1"]))
	}
}

// CA1 / R11 / D5: un Sprint fuera de las fechas del proyecto es válido si sus
// propias fechas son válidas.
func TestUS07_CrearSprint_FueraDeFechasDelProyectoValido(t *testing.T) {
	entorno := nuevoEntornoSprints()
	sembrarProyectosSprint(t, entorno, "p1")
	servicio := entorno.servicio()

	inicio := fechaUS07Servicio(2026, time.July, 1)
	fin := fechaUS07Servicio(2026, time.July, 15)
	creado, err := servicio.CrearSprint("p1", 3, "Planificar la iteración de julio", inicio, fin)
	if err != nil {
		t.Fatalf("CrearSprint() devolvió un error inesperado: %v", err)
	}
	if !creado.FechaInicio.Equal(inicio) || !creado.FechaFin.Equal(fin) {
		t.Errorf("fechas = %v-%v, se esperaban %v-%v", creado.FechaInicio, creado.FechaFin, inicio, fin)
	}
	if len(entorno.sprints["p1"]) != 1 {
		t.Errorf("se esperaba 1 Sprint persistido, hay %d", len(entorno.sprints["p1"]))
	}
}

// CA1 / R6 / D1: un proyecto inexistente devuelve ErrProyectoNoEncontrado y no
// evalúa las demás reglas ni persiste.
func TestUS07_CrearSprint_ProyectoNoEncontrado(t *testing.T) {
	entorno := nuevoEntornoSprints()
	sembrarProyectosSprint(t, entorno, "p1")
	servicio := entorno.servicio()

	creado, err := servicio.CrearSprint(
		"ausente",
		1,
		"Planificar la primera iteración",
		fechaUS07Servicio(2026, time.March, 2),
		fechaUS07Servicio(2026, time.March, 20),
	)
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	verificarSprintServicioUS07Vacio(t, creado)
	if entorno.existeActivoLlamadas != 0 {
		t.Errorf("ExisteActivo fue invocado %d veces; el proyecto debía validarse primero", entorno.existeActivoLlamadas)
	}
	if entorno.existeNumeroLlamadas != 0 {
		t.Errorf("ExisteNumero fue invocado %d veces; el proyecto debía validarse primero", entorno.existeNumeroLlamadas)
	}
	if entorno.agregados != 0 {
		t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.agregados)
	}
}

// CA2 / R2 / D3: un objetivo vacío o solo espacios no se persiste.
func TestUS07_CrearSprint_ObjetivoVacioNoPersiste(t *testing.T) {
	casos := []struct {
		descripcion string
		objetivo    string
	}{
		{descripcion: "objetivo vacío", objetivo: ""},
		{descripcion: "objetivo compuesto solo por espacios", objetivo: "   "},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoSprints()
			sembrarProyectosSprint(t, entorno, "p1")
			servicio := entorno.servicio()

			creado, err := servicio.CrearSprint(
				"p1",
				1,
				caso.objetivo,
				fechaUS07Servicio(2026, time.March, 2),
				fechaUS07Servicio(2026, time.March, 20),
			)
			if !errors.Is(err, domain.ErrObjetivoSprintVacio) {
				t.Fatalf("err = %v, se esperaba %v", err, domain.ErrObjetivoSprintVacio)
			}
			verificarSprintServicioUS07Vacio(t, creado)
			if entorno.agregados != 0 {
				t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.agregados)
			}
			if len(entorno.sprints["p1"]) != 0 {
				t.Errorf("se persistió un Sprint pese al objetivo vacío")
			}
		})
	}
}

// CA3 / R7 / D6: si el proyecto ya tiene un Sprint activo, la creación devuelve
// ErrSprintActivoExistente y no persiste.
func TestUS07_CrearSprint_SprintActivoExistenteNoPersiste(t *testing.T) {
	entorno := nuevoEntornoSprints()
	sembrarProyectosSprint(t, entorno, "p1")
	existente := sprintUS07Servicio(1, "Activo")
	entorno.sprints["p1"] = append(entorno.sprints["p1"], existente)
	servicio := entorno.servicio()

	creado, err := servicio.CrearSprint(
		"p1",
		2,
		"Planificar la segunda iteración",
		fechaUS07Servicio(2026, time.March, 23),
		fechaUS07Servicio(2026, time.April, 10),
	)
	if !errors.Is(err, domain.ErrSprintActivoExistente) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrSprintActivoExistente)
	}
	verificarSprintServicioUS07Vacio(t, creado)
	if entorno.agregados != 0 {
		t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.agregados)
	}
	if len(entorno.sprints["p1"]) != 1 {
		t.Errorf("se esperaba conservar 1 Sprint activo, hay %d", len(entorno.sprints["p1"]))
	}
}

// CA3 / R7 / D6: un Sprint activo en otro proyecto no interfiere.
func TestUS07_CrearSprint_SprintActivoEnOtroProyectoNoInterfiere(t *testing.T) {
	entorno := nuevoEntornoSprints()
	sembrarProyectosSprint(t, entorno, "p1", "p2")
	activoAjeno := sprintUS07Servicio(7, "Activo")
	entorno.sprints["p1"] = append(entorno.sprints["p1"], activoAjeno)
	servicio := entorno.servicio()

	creado, err := servicio.CrearSprint(
		"p2",
		2,
		"Planificar la primera iteración",
		fechaUS07Servicio(2026, time.March, 2),
		fechaUS07Servicio(2026, time.March, 20),
	)
	if err != nil {
		t.Fatalf("CrearSprint() devolvió un error inesperado: %v", err)
	}
	if creado.Numero != 2 || creado.Estado != "Activo" {
		t.Errorf("Sprint creado = %+v, se esperaba número 2 y estado Activo", creado)
	}
	if len(entorno.sprints["p1"]) != 1 || entorno.sprints["p1"][0].Numero != 7 {
		t.Errorf("el Sprint activo del otro proyecto cambió: %+v", entorno.sprints["p1"])
	}
	if len(entorno.sprints["p2"]) != 1 {
		t.Errorf("se esperaba 1 Sprint persistido en p2, hay %d", len(entorno.sprints["p2"]))
	}
}

// CA1 / R6 / R10 / D7: proyecto inexistente se evalúa antes que el número
// inválido.
func TestUS07_CrearSprint_ProyectoInexistenteAntesQueNumeroInvalido(t *testing.T) {
	entorno := nuevoEntornoSprints()
	sembrarProyectosSprint(t, entorno, "p1")
	servicio := entorno.servicio()

	creado, err := servicio.CrearSprint(
		"ausente",
		0,
		"Planificar la primera iteración",
		fechaUS07Servicio(2026, time.March, 2),
		fechaUS07Servicio(2026, time.March, 20),
	)
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if errors.Is(err, domain.ErrNumeroSprintInvalido) {
		t.Errorf("err = %v, no debería devolverse ErrNumeroSprintInvalido si el proyecto no existe", err)
	}
	verificarSprintServicioUS07Vacio(t, creado)
	if entorno.agregados != 0 {
		t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.agregados)
	}
}

// CA1 / R3 / R10 / D7: número inválido se evalúa antes que objetivo vacío.
func TestUS07_CrearSprint_NumeroInvalidoAntesQueObjetivoVacio(t *testing.T) {
	entorno := nuevoEntornoSprints()
	sembrarProyectosSprint(t, entorno, "p1")
	servicio := entorno.servicio()

	creado, err := servicio.CrearSprint(
		"p1",
		0,
		"",
		fechaUS07Servicio(2026, time.March, 2),
		fechaUS07Servicio(2026, time.March, 20),
	)
	if !errors.Is(err, domain.ErrNumeroSprintInvalido) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrNumeroSprintInvalido)
	}
	if errors.Is(err, domain.ErrObjetivoSprintVacio) {
		t.Errorf("err = %v, no debería devolverse ErrObjetivoSprintVacio si el número es inválido", err)
	}
	verificarSprintServicioUS07Vacio(t, creado)
	if entorno.agregados != 0 {
		t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.agregados)
	}
}

// CA2 / R2 / R10 / D7: objetivo vacío se evalúa antes que fechas inválidas.
func TestUS07_CrearSprint_ObjetivoVacioAntesQueFechasInvalidas(t *testing.T) {
	entorno := nuevoEntornoSprints()
	sembrarProyectosSprint(t, entorno, "p1")
	servicio := entorno.servicio()

	creado, err := servicio.CrearSprint(
		"p1",
		1,
		"   ",
		fechaUS07Servicio(2026, time.March, 20),
		fechaUS07Servicio(2026, time.March, 2),
	)
	if !errors.Is(err, domain.ErrObjetivoSprintVacio) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrObjetivoSprintVacio)
	}
	if errors.Is(err, domain.ErrFechasSprintInvalidas) {
		t.Errorf("err = %v, no debería devolverse ErrFechasSprintInvalidas si el objetivo está vacío", err)
	}
	verificarSprintServicioUS07Vacio(t, creado)
	if entorno.agregados != 0 {
		t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.agregados)
	}
}

// CA3 / R7 / R10 / D7: un Sprint activo existente se evalúa antes que el número
// duplicado.
func TestUS07_CrearSprint_SprintActivoAntesQueNumeroDuplicado(t *testing.T) {
	entorno := nuevoEntornoSprints()
	sembrarProyectosSprint(t, entorno, "p1")
	existente := sprintUS07Servicio(1, "Activo")
	entorno.sprints["p1"] = append(entorno.sprints["p1"], existente)
	servicio := entorno.servicio()

	creado, err := servicio.CrearSprint(
		"p1",
		1,
		"Planificar otra iteración",
		fechaUS07Servicio(2026, time.March, 23),
		fechaUS07Servicio(2026, time.April, 10),
	)
	if !errors.Is(err, domain.ErrSprintActivoExistente) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrSprintActivoExistente)
	}
	if errors.Is(err, domain.ErrSprintDuplicado) {
		t.Errorf("err = %v, no debería devolverse ErrSprintDuplicado si ya hay un Sprint activo", err)
	}
	if entorno.existeNumeroLlamadas != 0 {
		t.Errorf("ExisteNumero fue invocado %d veces; el Sprint activo debía evaluarse primero", entorno.existeNumeroLlamadas)
	}
	verificarSprintServicioUS07Vacio(t, creado)
	if entorno.agregados != 0 {
		t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.agregados)
	}
}

// CA1 / R9: los errores de repositorio se propagan envueltos con %w y no se
// mapean a ErrProyectoNoEncontrado.
func TestUS07_CrearSprint_ErrorDeRepositorioEnvuelto(t *testing.T) {
	causa := errors.New("fallo de repositorio")

	casos := []struct {
		descripcion string
		sembrar     bool
		preparar    func(*entornoSprints)
	}{
		{
			descripcion: "Existe",
			preparar:    func(e *entornoSprints) { e.existeErr = causa },
		},
		{
			descripcion: "ExisteActivo",
			sembrar:     true,
			preparar:    func(e *entornoSprints) { e.existeActivoErr = causa },
		},
		{
			descripcion: "ExisteNumero",
			sembrar:     true,
			preparar:    func(e *entornoSprints) { e.existeNumeroErr = causa },
		},
		{
			descripcion: "Agregar",
			sembrar:     true,
			preparar:    func(e *entornoSprints) { e.agregarErr = causa },
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoSprints()
			if caso.sembrar {
				sembrarProyectosSprint(t, entorno, "p1")
			}
			caso.preparar(entorno)
			servicio := entorno.servicio()

			creado, err := servicio.CrearSprint(
				"p1",
				1,
				"Planificar la primera iteración",
				fechaUS07Servicio(2026, time.March, 2),
				fechaUS07Servicio(2026, time.March, 20),
			)
			if err == nil {
				t.Fatal("CrearSprint() debía devolver el error de repositorio")
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
			if errors.Is(err, domain.ErrProyectoNoEncontrado) {
				t.Errorf("err = %v, no debería mapearse a ErrProyectoNoEncontrado", err)
			}
			verificarSprintServicioUS07Vacio(t, creado)
		})
	}
}

// CA1 / CA2 / CA3 / R9: ningún error de repositorio persiste el Sprint ni se
// mapea a un error de dominio de no encontrado.
func TestUS07_CrearSprint_ErrorNoEscribe(t *testing.T) {
	causa := errors.New("fallo de repositorio")

	casos := []struct {
		descripcion string
		sembrar     bool
		preparar    func(*entornoSprints)
	}{
		{
			descripcion: "Existe",
			preparar:    func(e *entornoSprints) { e.existeErr = causa },
		},
		{
			descripcion: "ExisteActivo",
			sembrar:     true,
			preparar:    func(e *entornoSprints) { e.existeActivoErr = causa },
		},
		{
			descripcion: "ExisteNumero",
			sembrar:     true,
			preparar:    func(e *entornoSprints) { e.existeNumeroErr = causa },
		},
		{
			descripcion: "Agregar",
			sembrar:     true,
			preparar:    func(e *entornoSprints) { e.agregarErr = causa },
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoSprints()
			if caso.sembrar {
				sembrarProyectosSprint(t, entorno, "p1")
			}
			caso.preparar(entorno)
			servicio := entorno.servicio()

			_, err := servicio.CrearSprint(
				"p1",
				1,
				"Planificar la primera iteración",
				fechaUS07Servicio(2026, time.March, 2),
				fechaUS07Servicio(2026, time.March, 20),
			)
			if err == nil {
				t.Fatal("CrearSprint() debía devolver el error de repositorio")
			}
			if errors.Is(err, domain.ErrProyectoNoEncontrado) {
				t.Errorf("err = %v, no debería mapearse a ErrProyectoNoEncontrado", err)
			}
			if len(entorno.sprints["p1"]) != 0 {
				t.Errorf("se persistió un Sprint pese al error de repositorio: %+v", entorno.sprints["p1"])
			}
		})
	}
}
