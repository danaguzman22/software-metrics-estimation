package service

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

const (
	proyectoEsfuerzoUS14   = "p1"
	itemEsfuerzoUS14       = "HU-01"
	integranteEsfuerzoUS14 = "Ana Pérez"
)

func hoyEsfuerzoUS14() time.Time {
	return time.Date(2026, time.November, 15, 9, 0, 0, 0, time.UTC)
}

func fechaEsfuerzoUS14(anio int, mes time.Month, dia, hora, minuto int) time.Time {
	return time.Date(anio, mes, dia, hora, minuto, 0, 0, time.UTC)
}

// entornoEsfuerzo agrupa los dobles que necesita ServicioEsfuerzo.
type entornoEsfuerzo struct {
	proyectos    *fakeRepositorioProyectos
	backlog      *fakeRepositorioBacklog
	integrantes  *fakeRepositorioIntegrantes
	sprints      *fakeRepositorioSprints
	asignaciones *fakeRepositorioAsignaciones
	esfuerzo     *fakeRepositorioEsfuerzo
}

func nuevoEntornoEsfuerzo() *entornoEsfuerzo {
	return &entornoEsfuerzo{
		proyectos:    nuevoFakeProyectos(),
		backlog:      nuevoFakeBacklog(),
		integrantes:  nuevoFakeIntegrantes(),
		sprints:      nuevoFakeSprints(),
		asignaciones: nuevoFakeAsignaciones(),
		esfuerzo:     nuevoFakeEsfuerzo(),
	}
}

func (e *entornoEsfuerzo) servicio() *ServicioEsfuerzo {
	return NuevoServicioEsfuerzo(
		e.proyectos,
		e.backlog,
		e.integrantes,
		e.sprints,
		e.asignaciones,
		e.esfuerzo,
	)
}

func sembrarProyectoEsfuerzo(t *testing.T, entorno *entornoEsfuerzo, proyectoID string) {
	t.Helper()
	entorno.proyectos.proyectos[proyectoID] = proyectoFakeValido()
}

func sembrarItemEsfuerzo(
	t *testing.T,
	entorno *entornoEsfuerzo,
	proyectoID, itemID, estado string,
) {
	t.Helper()
	entorno.backlog.items[proyectoID] = append(
		entorno.backlog.items[proyectoID],
		domain.ItemBacklog{
			ID:          itemID,
			Titulo:      "Historia " + itemID,
			Descripcion: "Historia de prueba de US-14",
			Prioridad:   "Alta",
			Estado:      estado,
			StoryPoints: 3,
		},
	)
}

func sembrarIntegranteEsfuerzo(
	t *testing.T,
	entorno *entornoEsfuerzo,
	proyectoID, nombre string,
) {
	t.Helper()
	entorno.integrantes.integrantes[proyectoID] = append(
		entorno.integrantes.integrantes[proyectoID],
		domain.Integrante{Nombre: nombre, Rol: "Desarrolladora"},
	)
}

func sembrarSprintEsfuerzo(
	t *testing.T,
	entorno *entornoEsfuerzo,
	proyectoID string,
	numero int,
	estado string,
) {
	t.Helper()
	entorno.sprints.sprints[proyectoID] = append(
		entorno.sprints.sprints[proyectoID],
		domain.Sprint{
			Numero:      numero,
			Objetivo:    "Sprint de prueba de US-14",
			FechaInicio: fechaEsfuerzoUS14(2026, time.November, 1, 0, 0),
			FechaFin:    fechaEsfuerzoUS14(2026, time.November, 30, 0, 0),
			Estado:      estado,
		},
	)
}

func sembrarAsignacionEsfuerzo(
	t *testing.T,
	entorno *entornoEsfuerzo,
	proyectoID string,
	sprintNumero, itemNumero int,
) {
	t.Helper()
	entorno.asignaciones.sembrar(proyectoID, sprintNumero, itemNumero)
}

func entornoEsfuerzoConItem(t *testing.T, estado string) *entornoEsfuerzo {
	t.Helper()
	entorno := nuevoEntornoEsfuerzo()
	sembrarProyectoEsfuerzo(t, entorno, proyectoEsfuerzoUS14)
	sembrarItemEsfuerzo(t, entorno, proyectoEsfuerzoUS14, itemEsfuerzoUS14, estado)
	sembrarIntegranteEsfuerzo(t, entorno, proyectoEsfuerzoUS14, integranteEsfuerzoUS14)
	return entorno
}

func entornoEsfuerzoValido(t *testing.T) *entornoEsfuerzo {
	t.Helper()
	return entornoEsfuerzoConItem(t, domain.EstadoItemPendiente)
}

func verificarRegistroEsfuerzoUS14(
	t *testing.T,
	obtenido, esperado domain.RegistroEsfuerzo,
) {
	t.Helper()
	if obtenido.Integrante != esperado.Integrante {
		t.Errorf("Integrante = %q, se esperaba %q", obtenido.Integrante, esperado.Integrante)
	}
	if !obtenido.Fecha.Equal(esperado.Fecha) {
		t.Errorf("Fecha = %v, se esperaba %v", obtenido.Fecha, esperado.Fecha)
	}
	if obtenido.Actividad != esperado.Actividad {
		t.Errorf("Actividad = %q, se esperaba %q", obtenido.Actividad, esperado.Actividad)
	}
	if obtenido.Horas != esperado.Horas {
		t.Errorf("Horas = %v, se esperaba %v", obtenido.Horas, esperado.Horas)
	}
}

func verificarRegistroVacioUS14(t *testing.T, registro domain.RegistroEsfuerzo) {
	t.Helper()
	if registro.Integrante != "" || !registro.Fecha.IsZero() ||
		registro.Actividad != "" || registro.Horas != 0 {
		t.Errorf("se esperaba domain.RegistroEsfuerzo{} ante un error, se obtuvo %+v", registro)
	}
}

func verificarPersistidoEsfuerzoUS14(
	t *testing.T,
	entorno *entornoEsfuerzo,
	proyectoID string,
	itemNumero int,
	esperado domain.RegistroEsfuerzo,
) {
	t.Helper()
	if entorno.esfuerzo.agregarLlamadas != 1 {
		t.Fatalf("Agregar fue invocado %d veces, se esperaba 1", entorno.esfuerzo.agregarLlamadas)
	}
	if len(entorno.esfuerzo.registros) != 1 {
		t.Fatalf("se esperaba 1 registro persistido, hay %d", len(entorno.esfuerzo.registros))
	}
	persistido := entorno.esfuerzo.registros[0]
	if persistido.proyectoID != proyectoID {
		t.Errorf("proyecto persistido = %q, se esperaba %q", persistido.proyectoID, proyectoID)
	}
	if persistido.itemNumero != itemNumero {
		t.Errorf("item persistido = %d, se esperaba %d", persistido.itemNumero, itemNumero)
	}
	verificarRegistroEsfuerzoUS14(t, persistido.registro, esperado)
}

func verificarSinRegistroEsfuerzoUS14(t *testing.T, entorno *entornoEsfuerzo) {
	t.Helper()
	if entorno.esfuerzo.agregarLlamadas != 0 {
		t.Errorf("Agregar fue invocado %d veces; no debía escribirse nada", entorno.esfuerzo.agregarLlamadas)
	}
	if len(entorno.esfuerzo.registros) != 0 {
		t.Errorf("se persistió un registro pese al error: %+v", entorno.esfuerzo.registros)
	}
}

func verificarErrorEsfuerzoUS14(t *testing.T, err, esperado error) {
	t.Helper()
	if !errors.Is(err, esperado) {
		t.Fatalf("err = %v, se esperaba %v", err, esperado)
	}
}

func verificarErrorRepositorioEsfuerzoUS14(t *testing.T, err, causa error) {
	t.Helper()
	if err == nil {
		t.Fatal("se esperaba un error de repositorio y se obtuvo nil")
	}
	if !errors.Is(err, causa) {
		t.Fatalf("errors.Is(err, causa) = false; err = %v", err)
	}
	if err == causa {
		t.Fatal("el error de repositorio debe estar envuelto; err == causa")
	}
	if !strings.HasPrefix(err.Error(), "service:") {
		t.Errorf("mensaje = %q; se esperaba prefijo service:", err.Error())
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) ||
		errors.Is(err, domain.ErrItemNoEncontrado) ||
		errors.Is(err, domain.ErrSprintNoEncontrado) ||
		errors.Is(err, domain.ErrIntegranteNoPerteneceAlProyecto) ||
		errors.Is(err, domain.ErrHorasEsfuerzoInvalidas) ||
		errors.Is(err, domain.ErrFechaEsfuerzoFutura) {
		t.Errorf("err = %v no debe mapearse a un error de dominio", err)
	}
}

// CA1 / R1: un alta válida persiste integrante, fecha, actividad y horas
// asociados al proyecto y a la historia.
func TestUS14_RegistrarEsfuerzo_AltaValida(t *testing.T) {
	entorno := entornoEsfuerzoValido(t)
	fecha := fechaEsfuerzoUS14(2026, time.November, 15, 9, 0)
	esperado := domain.RegistroEsfuerzo{
		Integrante: integranteEsfuerzoUS14,
		Fecha:      fecha,
		Actividad:  "Análisis",
		Horas:      1.5,
	}

	registro, err := entorno.servicio().RegistrarEsfuerzo(
		proyectoEsfuerzoUS14,
		itemEsfuerzoUS14,
		integranteEsfuerzoUS14,
		fecha,
		"Análisis",
		1.5,
		hoyEsfuerzoUS14(),
	)
	if err != nil {
		t.Fatalf("RegistrarEsfuerzo() devolvió un error inesperado: %v", err)
	}

	verificarRegistroEsfuerzoUS14(t, registro, esperado)
	verificarPersistidoEsfuerzoUS14(t, entorno, proyectoEsfuerzoUS14, 1, esperado)
}

// CA1 / R9: una actividad vacía o solo con espacios se guarda vacía.
func TestUS14_RegistrarEsfuerzo_ActividadVacia(t *testing.T) {
	casos := []struct {
		descripcion string
		actividad   string
	}{
		{descripcion: "actividad vacía", actividad: ""},
		{descripcion: "actividad solo con espacios", actividad: "   "},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := entornoEsfuerzoValido(t)
			fecha := fechaEsfuerzoUS14(2026, time.November, 15, 9, 0)
			esperado := domain.RegistroEsfuerzo{
				Integrante: integranteEsfuerzoUS14,
				Fecha:      fecha,
				Actividad:  "",
				Horas:      1.5,
			}

			registro, err := entorno.servicio().RegistrarEsfuerzo(
				proyectoEsfuerzoUS14,
				itemEsfuerzoUS14,
				integranteEsfuerzoUS14,
				fecha,
				caso.actividad,
				1.5,
				hoyEsfuerzoUS14(),
			)
			if err != nil {
				t.Fatalf("RegistrarEsfuerzo() devolvió un error inesperado: %v", err)
			}

			verificarRegistroEsfuerzoUS14(t, registro, esperado)
			verificarPersistidoEsfuerzoUS14(t, entorno, proyectoEsfuerzoUS14, 1, esperado)
		})
	}
}

// CA1 / R9: una actividad con contenido y espacios externos se guarda recortada.
func TestUS14_RegistrarEsfuerzo_ActividadRecortada(t *testing.T) {
	entorno := entornoEsfuerzoValido(t)
	fecha := fechaEsfuerzoUS14(2026, time.November, 15, 9, 0)
	esperado := domain.RegistroEsfuerzo{
		Integrante: integranteEsfuerzoUS14,
		Fecha:      fecha,
		Actividad:  "Análisis",
		Horas:      1.5,
	}

	registro, err := entorno.servicio().RegistrarEsfuerzo(
		proyectoEsfuerzoUS14,
		itemEsfuerzoUS14,
		integranteEsfuerzoUS14,
		fecha,
		"  Análisis  ",
		1.5,
		hoyEsfuerzoUS14(),
	)
	if err != nil {
		t.Fatalf("RegistrarEsfuerzo() devolvió un error inesperado: %v", err)
	}

	verificarRegistroEsfuerzoUS14(t, registro, esperado)
	verificarPersistidoEsfuerzoUS14(t, entorno, proyectoEsfuerzoUS14, 1, esperado)
}

// CA1 / R4 / D13: el nombre se guarda tal como está registrado, aunque la
// entrada difiera en mayúsculas y espacios externos.
func TestUS14_RegistrarEsfuerzo_IntegranteTalComoEstaRegistrado(t *testing.T) {
	entorno := entornoEsfuerzoValido(t)
	fecha := fechaEsfuerzoUS14(2026, time.November, 15, 9, 0)
	esperado := domain.RegistroEsfuerzo{
		Integrante: integranteEsfuerzoUS14,
		Fecha:      fecha,
		Actividad:  "Análisis",
		Horas:      1.5,
	}

	registro, err := entorno.servicio().RegistrarEsfuerzo(
		proyectoEsfuerzoUS14,
		itemEsfuerzoUS14,
		"  ana pérez  ",
		fecha,
		"Análisis",
		1.5,
		hoyEsfuerzoUS14(),
	)
	if err != nil {
		t.Fatalf("RegistrarEsfuerzo() devolvió un error inesperado: %v", err)
	}

	verificarRegistroEsfuerzoUS14(t, registro, esperado)
	verificarPersistidoEsfuerzoUS14(t, entorno, proyectoEsfuerzoUS14, 1, esperado)
}

// CA3 / R2 / V1 / D18: el proyecto inexistente prevalece sobre cualquier otro
// error y no se persiste nada.
func TestUS14_RegistrarEsfuerzo_ProyectoNoEncontrado(t *testing.T) {
	entorno := nuevoEntornoEsfuerzo()

	registro, err := entorno.servicio().RegistrarEsfuerzo(
		"p999",
		itemEsfuerzoUS14,
		integranteEsfuerzoUS14,
		fechaEsfuerzoUS14(2026, time.November, 16, 12, 30),
		"Análisis",
		0,
		hoyEsfuerzoUS14(),
	)
	verificarErrorEsfuerzoUS14(t, err, domain.ErrProyectoNoEncontrado)
	if errors.Is(err, domain.ErrHorasEsfuerzoInvalidas) {
		t.Errorf("err = %v, no debería devolverse ErrHorasEsfuerzoInvalidas antes que el proyecto", err)
	}
	if errors.Is(err, domain.ErrFechaEsfuerzoFutura) {
		t.Errorf("err = %v, no debería devolverse ErrFechaEsfuerzoFutura antes que el proyecto", err)
	}
	verificarRegistroVacioUS14(t, registro)
	verificarSinRegistroEsfuerzoUS14(t, entorno)
}

// CA3 / R3 / R7 / V4: una historia inexistente devuelve ErrItemNoEncontrado y
// no se persiste nada.
func TestUS14_RegistrarEsfuerzo_ItemNoEncontrado(t *testing.T) {
	entorno := entornoEsfuerzoValido(t)

	registro, err := entorno.servicio().RegistrarEsfuerzo(
		proyectoEsfuerzoUS14,
		"HU-99",
		integranteEsfuerzoUS14,
		fechaEsfuerzoUS14(2026, time.November, 15, 9, 0),
		"Análisis",
		1.5,
		hoyEsfuerzoUS14(),
	)
	verificarErrorEsfuerzoUS14(t, err, domain.ErrItemNoEncontrado)
	verificarRegistroVacioUS14(t, registro)
	verificarSinRegistroEsfuerzoUS14(t, entorno)
}

// CA3 / R3: el mismo itemID en otro proyecto no cuenta como historia existente.
func TestUS14_RegistrarEsfuerzo_ItemDeOtroProyecto(t *testing.T) {
	entorno := nuevoEntornoEsfuerzo()
	sembrarProyectoEsfuerzo(t, entorno, "p1")
	sembrarProyectoEsfuerzo(t, entorno, "p2")
	sembrarItemEsfuerzo(t, entorno, "p2", itemEsfuerzoUS14, domain.EstadoItemPendiente)
	sembrarIntegranteEsfuerzo(t, entorno, "p1", integranteEsfuerzoUS14)

	registro, err := entorno.servicio().RegistrarEsfuerzo(
		"p1",
		itemEsfuerzoUS14,
		integranteEsfuerzoUS14,
		fechaEsfuerzoUS14(2026, time.November, 15, 9, 0),
		"Análisis",
		1.5,
		hoyEsfuerzoUS14(),
	)
	verificarErrorEsfuerzoUS14(t, err, domain.ErrItemNoEncontrado)
	verificarRegistroVacioUS14(t, registro)
	verificarSinRegistroEsfuerzoUS14(t, entorno)
}

// CA3 / R7 / V4: un ID con formato inválido se trata como historia inexistente.
func TestUS14_RegistrarEsfuerzo_IDMalformado(t *testing.T) {
	casos := []string{"historia-uno", "HU-", "HU-XX", "XX-01", "HU-00"}

	for _, itemID := range casos {
		t.Run(itemID, func(t *testing.T) {
			entorno := entornoEsfuerzoValido(t)

			registro, err := entorno.servicio().RegistrarEsfuerzo(
				proyectoEsfuerzoUS14,
				itemID,
				integranteEsfuerzoUS14,
				fechaEsfuerzoUS14(2026, time.November, 15, 9, 0),
				"Análisis",
				1.5,
				hoyEsfuerzoUS14(),
			)
			verificarErrorEsfuerzoUS14(t, err, domain.ErrItemNoEncontrado)
			verificarRegistroVacioUS14(t, registro)
			verificarSinRegistroEsfuerzoUS14(t, entorno)
		})
	}
}

// CA3 / R4 / V6: un integrante que no pertenece al proyecto no registra
// esfuerzo, incluido el caso de existir solo en otro proyecto.
func TestUS14_RegistrarEsfuerzo_IntegranteAjeno(t *testing.T) {
	casos := []struct {
		descripcion string
		preparar    func(t *testing.T, entorno *entornoEsfuerzo)
	}{
		{
			descripcion: "integrante no registrado en el proyecto",
			preparar:    func(t *testing.T, entorno *entornoEsfuerzo) {},
		},
		{
			descripcion: "integrante registrado solo en otro proyecto",
			preparar: func(t *testing.T, entorno *entornoEsfuerzo) {
				sembrarIntegranteEsfuerzo(t, entorno, "p2", "Luis Gómez")
			},
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := entornoEsfuerzoValido(t)
			caso.preparar(t, entorno)

			registro, err := entorno.servicio().RegistrarEsfuerzo(
				proyectoEsfuerzoUS14,
				itemEsfuerzoUS14,
				"Luis Gómez",
				fechaEsfuerzoUS14(2026, time.November, 15, 9, 0),
				"Análisis",
				1.5,
				hoyEsfuerzoUS14(),
			)
			verificarErrorEsfuerzoUS14(t, err, domain.ErrIntegranteNoPerteneceAlProyecto)
			verificarRegistroVacioUS14(t, registro)
			verificarSinRegistroEsfuerzoUS14(t, entorno)
		})
	}
}

// CA3 / R10 / V5: una historia de un Sprint cerrado no admite esfuerzo.
func TestUS14_RegistrarEsfuerzo_SprintCerrado(t *testing.T) {
	entorno := entornoEsfuerzoValido(t)
	sembrarSprintEsfuerzo(t, entorno, proyectoEsfuerzoUS14, 1, domain.EstadoSprintCerrado)
	sembrarAsignacionEsfuerzo(t, entorno, proyectoEsfuerzoUS14, 1, 1)

	registro, err := entorno.servicio().RegistrarEsfuerzo(
		proyectoEsfuerzoUS14,
		itemEsfuerzoUS14,
		integranteEsfuerzoUS14,
		fechaEsfuerzoUS14(2026, time.November, 15, 9, 0),
		"Análisis",
		1.5,
		hoyEsfuerzoUS14(),
	)
	verificarErrorEsfuerzoUS14(t, err, domain.ErrSprintCerrado)
	verificarRegistroVacioUS14(t, registro)
	verificarSinRegistroEsfuerzoUS14(t, entorno)
}

// CA3 / R10 / D14: una historia completada de un Sprint activo admite el
// registro.
func TestUS14_RegistrarEsfuerzo_ItemCompletadoSprintActivo(t *testing.T) {
	entorno := entornoEsfuerzoConItem(t, domain.EstadoCompletada)
	sembrarSprintEsfuerzo(t, entorno, proyectoEsfuerzoUS14, 1, domain.EstadoSprintActivo)
	sembrarAsignacionEsfuerzo(t, entorno, proyectoEsfuerzoUS14, 1, 1)
	fecha := fechaEsfuerzoUS14(2026, time.November, 15, 9, 0)
	esperado := domain.RegistroEsfuerzo{
		Integrante: integranteEsfuerzoUS14,
		Fecha:      fecha,
		Actividad:  "Análisis",
		Horas:      1.5,
	}

	registro, err := entorno.servicio().RegistrarEsfuerzo(
		proyectoEsfuerzoUS14,
		itemEsfuerzoUS14,
		integranteEsfuerzoUS14,
		fecha,
		"Análisis",
		1.5,
		hoyEsfuerzoUS14(),
	)
	if err != nil {
		t.Fatalf("RegistrarEsfuerzo() devolvió un error inesperado: %v", err)
	}
	verificarRegistroEsfuerzoUS14(t, registro, esperado)
	verificarPersistidoEsfuerzoUS14(t, entorno, proyectoEsfuerzoUS14, 1, esperado)
}

// CA3 / R10 / D19: una historia sin Sprint asignado admite el registro.
func TestUS14_RegistrarEsfuerzo_ItemSinSprint(t *testing.T) {
	entorno := entornoEsfuerzoValido(t)
	fecha := fechaEsfuerzoUS14(2026, time.November, 15, 9, 0)
	esperado := domain.RegistroEsfuerzo{
		Integrante: integranteEsfuerzoUS14,
		Fecha:      fecha,
		Actividad:  "Análisis",
		Horas:      1.5,
	}

	registro, err := entorno.servicio().RegistrarEsfuerzo(
		proyectoEsfuerzoUS14,
		itemEsfuerzoUS14,
		integranteEsfuerzoUS14,
		fecha,
		"Análisis",
		1.5,
		hoyEsfuerzoUS14(),
	)
	if err != nil {
		t.Fatalf("RegistrarEsfuerzo() devolvió un error inesperado: %v", err)
	}
	verificarRegistroEsfuerzoUS14(t, registro, esperado)
	verificarPersistidoEsfuerzoUS14(t, entorno, proyectoEsfuerzoUS14, 1, esperado)
	if entorno.asignaciones.sprintDeLlamadas != 1 {
		t.Errorf("SprintDeItem fue invocado %d veces, se esperaba 1", entorno.asignaciones.sprintDeLlamadas)
	}
}

// CA2 / R5: las horas válidas en (0, 24] con hasta dos decimales no se
// rechazan.
func TestUS14_RegistrarEsfuerzo_HorasValidas(t *testing.T) {
	casos := []struct {
		descripcion string
		horas       float64
	}{
		{descripcion: "0.01 horas", horas: 0.01},
		{descripcion: "1.5 horas", horas: 1.5},
		{descripcion: "24 horas", horas: 24},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := entornoEsfuerzoValido(t)
			fecha := fechaEsfuerzoUS14(2026, time.November, 15, 9, 0)
			esperado := domain.RegistroEsfuerzo{
				Integrante: integranteEsfuerzoUS14,
				Fecha:      fecha,
				Actividad:  "Análisis",
				Horas:      caso.horas,
			}

			registro, err := entorno.servicio().RegistrarEsfuerzo(
				proyectoEsfuerzoUS14,
				itemEsfuerzoUS14,
				integranteEsfuerzoUS14,
				fecha,
				"Análisis",
				caso.horas,
				hoyEsfuerzoUS14(),
			)
			if err != nil {
				t.Fatalf("RegistrarEsfuerzo() devolvió un error inesperado: %v", err)
			}
			verificarRegistroEsfuerzoUS14(t, registro, esperado)
			verificarPersistidoEsfuerzoUS14(t, entorno, proyectoEsfuerzoUS14, 1, esperado)
		})
	}
}

// CA2 / R5 / V2 / V8: las horas fuera de (0, 24] o con más de dos decimales
// devuelven ErrHorasEsfuerzoInvalidas sin persistir ni redondear.
func TestUS14_RegistrarEsfuerzo_HorasInvalidas(t *testing.T) {
	casos := []struct {
		descripcion string
		horas       float64
	}{
		{descripcion: "0 horas", horas: 0},
		{descripcion: "horas negativas", horas: -1},
		{descripcion: "24.01 horas", horas: 24.01},
		{descripcion: "25 horas", horas: 25},
		{descripcion: "1.234 horas con más de dos decimales", horas: 1.234},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := entornoEsfuerzoValido(t)

			registro, err := entorno.servicio().RegistrarEsfuerzo(
				proyectoEsfuerzoUS14,
				itemEsfuerzoUS14,
				integranteEsfuerzoUS14,
				fechaEsfuerzoUS14(2026, time.November, 15, 9, 0),
				"Análisis",
				caso.horas,
				hoyEsfuerzoUS14(),
			)
			verificarErrorEsfuerzoUS14(t, err, domain.ErrHorasEsfuerzoInvalidas)
			verificarRegistroVacioUS14(t, registro)
			verificarSinRegistroEsfuerzoUS14(t, entorno)
		})
	}
}

// CA1 / R8: el día de hoy se registra sin importar el horario.
func TestUS14_RegistrarEsfuerzo_FechaHoyCualquierHora(t *testing.T) {
	casos := []struct {
		descripcion string
		fecha       time.Time
	}{
		{descripcion: "hoy a las 00:00", fecha: fechaEsfuerzoUS14(2026, time.November, 15, 0, 0)},
		{descripcion: "hoy a las 12:30", fecha: fechaEsfuerzoUS14(2026, time.November, 15, 12, 30)},
		{descripcion: "hoy a las 23:59", fecha: fechaEsfuerzoUS14(2026, time.November, 15, 23, 59)},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := entornoEsfuerzoValido(t)
			esperado := domain.RegistroEsfuerzo{
				Integrante: integranteEsfuerzoUS14,
				Fecha:      caso.fecha,
				Actividad:  "Análisis",
				Horas:      1.5,
			}

			registro, err := entorno.servicio().RegistrarEsfuerzo(
				proyectoEsfuerzoUS14,
				itemEsfuerzoUS14,
				integranteEsfuerzoUS14,
				caso.fecha,
				"Análisis",
				1.5,
				hoyEsfuerzoUS14(),
			)
			if err != nil {
				t.Fatalf("RegistrarEsfuerzo() devolvió un error inesperado: %v", err)
			}
			verificarRegistroEsfuerzoUS14(t, registro, esperado)
			verificarPersistidoEsfuerzoUS14(t, entorno, proyectoEsfuerzoUS14, 1, esperado)
		})
	}
}

// CA1 / R8 / V3: la fecha de mañana devuelve ErrFechaEsfuerzoFutura sin
// importar el horario y no persiste nada.
func TestUS14_RegistrarEsfuerzo_FechaFutura(t *testing.T) {
	casos := []struct {
		descripcion string
		fecha       time.Time
	}{
		{descripcion: "mañana a las 00:00", fecha: fechaEsfuerzoUS14(2026, time.November, 16, 0, 0)},
		{descripcion: "mañana a las 12:30", fecha: fechaEsfuerzoUS14(2026, time.November, 16, 12, 30)},
		{descripcion: "mañana a las 23:59", fecha: fechaEsfuerzoUS14(2026, time.November, 16, 23, 59)},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := entornoEsfuerzoValido(t)

			registro, err := entorno.servicio().RegistrarEsfuerzo(
				proyectoEsfuerzoUS14,
				itemEsfuerzoUS14,
				integranteEsfuerzoUS14,
				caso.fecha,
				"Análisis",
				1.5,
				hoyEsfuerzoUS14(),
			)
			verificarErrorEsfuerzoUS14(t, err, domain.ErrFechaEsfuerzoFutura)
			verificarRegistroVacioUS14(t, registro)
			verificarSinRegistroEsfuerzoUS14(t, entorno)
		})
	}
}

// CA3 / D18: las horas se evalúan antes que la fecha.
func TestUS14_RegistrarEsfuerzo_OrdenHorasAntesQueFecha(t *testing.T) {
	entorno := entornoEsfuerzoValido(t)

	_, err := entorno.servicio().RegistrarEsfuerzo(
		proyectoEsfuerzoUS14,
		itemEsfuerzoUS14,
		integranteEsfuerzoUS14,
		fechaEsfuerzoUS14(2026, time.November, 16, 12, 30),
		"Análisis",
		0,
		hoyEsfuerzoUS14(),
	)
	verificarErrorEsfuerzoUS14(t, err, domain.ErrHorasEsfuerzoInvalidas)
	if errors.Is(err, domain.ErrFechaEsfuerzoFutura) {
		t.Errorf("err = %v, no debería devolverse ErrFechaEsfuerzoFutura antes que las horas", err)
	}
	verificarSinRegistroEsfuerzoUS14(t, entorno)
}

// CA3 / D18: los datos (horas y fecha) se evalúan antes que la historia.
func TestUS14_RegistrarEsfuerzo_OrdenDatosAntesQueHistoria(t *testing.T) {
	casos := []struct {
		descripcion string
		fecha       time.Time
		horas       float64
		esperado    error
	}{
		{
			descripcion: "horas inválidas con historia inexistente",
			fecha:       fechaEsfuerzoUS14(2026, time.November, 15, 9, 0),
			horas:       0,
			esperado:    domain.ErrHorasEsfuerzoInvalidas,
		},
		{
			descripcion: "fecha futura con historia inexistente",
			fecha:       fechaEsfuerzoUS14(2026, time.November, 16, 12, 30),
			horas:       1.5,
			esperado:    domain.ErrFechaEsfuerzoFutura,
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoEsfuerzo()
			sembrarProyectoEsfuerzo(t, entorno, proyectoEsfuerzoUS14)
			sembrarIntegranteEsfuerzo(t, entorno, proyectoEsfuerzoUS14, integranteEsfuerzoUS14)

			_, err := entorno.servicio().RegistrarEsfuerzo(
				proyectoEsfuerzoUS14,
				"HU-99",
				integranteEsfuerzoUS14,
				caso.fecha,
				"Análisis",
				caso.horas,
				hoyEsfuerzoUS14(),
			)
			verificarErrorEsfuerzoUS14(t, err, caso.esperado)
			if errors.Is(err, domain.ErrItemNoEncontrado) {
				t.Errorf("err = %v, no debería devolverse ErrItemNoEncontrado antes que los datos", err)
			}
			if entorno.backlog.obtenerLlamadas != 0 {
				t.Errorf("ObtenerItem fue invocado %d veces; los datos debían evaluarse primero", entorno.backlog.obtenerLlamadas)
			}
			verificarSinRegistroEsfuerzoUS14(t, entorno)
		})
	}
}

// CA3 / D18: la historia se evalúa antes que el Sprint.
func TestUS14_RegistrarEsfuerzo_OrdenHistoriaAntesQueSprint(t *testing.T) {
	entorno := nuevoEntornoEsfuerzo()
	sembrarProyectoEsfuerzo(t, entorno, proyectoEsfuerzoUS14)
	sembrarIntegranteEsfuerzo(t, entorno, proyectoEsfuerzoUS14, integranteEsfuerzoUS14)

	_, err := entorno.servicio().RegistrarEsfuerzo(
		proyectoEsfuerzoUS14,
		"HU-99",
		integranteEsfuerzoUS14,
		fechaEsfuerzoUS14(2026, time.November, 15, 9, 0),
		"Análisis",
		1.5,
		hoyEsfuerzoUS14(),
	)
	verificarErrorEsfuerzoUS14(t, err, domain.ErrItemNoEncontrado)
	if entorno.asignaciones.sprintDeLlamadas != 0 {
		t.Errorf("SprintDeItem fue invocado %d veces; la historia debía evaluarse primero", entorno.asignaciones.sprintDeLlamadas)
	}
	verificarSinRegistroEsfuerzoUS14(t, entorno)
}

// CA3 / D18: el Sprint se evalúa antes que el integrante.
func TestUS14_RegistrarEsfuerzo_OrdenSprintAntesQueIntegrante(t *testing.T) {
	entorno := entornoEsfuerzoValido(t)
	sembrarSprintEsfuerzo(t, entorno, proyectoEsfuerzoUS14, 1, domain.EstadoSprintCerrado)
	sembrarAsignacionEsfuerzo(t, entorno, proyectoEsfuerzoUS14, 1, 1)

	_, err := entorno.servicio().RegistrarEsfuerzo(
		proyectoEsfuerzoUS14,
		itemEsfuerzoUS14,
		"Luis Gómez",
		fechaEsfuerzoUS14(2026, time.November, 15, 9, 0),
		"Análisis",
		1.5,
		hoyEsfuerzoUS14(),
	)
	verificarErrorEsfuerzoUS14(t, err, domain.ErrSprintCerrado)
	if errors.Is(err, domain.ErrIntegranteNoPerteneceAlProyecto) {
		t.Errorf("err = %v, no debería devolverse ErrIntegranteNoPerteneceAlProyecto antes que el Sprint", err)
	}
	if entorno.integrantes.listarLlamadas != 0 {
		t.Errorf("ListarPorProyecto fue invocado %d veces; el Sprint debía evaluarse primero", entorno.integrantes.listarLlamadas)
	}
	verificarSinRegistroEsfuerzoUS14(t, entorno)
}

// CA1 / CA2 / CA3 / R6 / E6: los errores de repositorio no conocidos se
// propagan envueltos con %w y no se mapean a errores de dominio.
func TestUS14_RegistrarEsfuerzo_ErrorRepositorioEnvuelto(t *testing.T) {
	causa := errors.New("fallo de repositorio")

	casos := []struct {
		descripcion     string
		preparar        func(t *testing.T, entorno *entornoEsfuerzo, causa error)
		agregarLlamadas int
	}{
		{
			descripcion: "verificar proyecto",
			preparar: func(t *testing.T, entorno *entornoEsfuerzo, causa error) {
				entorno.proyectos.existeErr = causa
			},
		},
		{
			descripcion: "obtener ítem",
			preparar: func(t *testing.T, entorno *entornoEsfuerzo, causa error) {
				sembrarProyectoEsfuerzo(t, entorno, proyectoEsfuerzoUS14)
				entorno.backlog.obtenerErr = causa
			},
		},
		{
			descripcion: "obtener Sprint de la historia",
			preparar: func(t *testing.T, entorno *entornoEsfuerzo, causa error) {
				sembrarProyectoEsfuerzo(t, entorno, proyectoEsfuerzoUS14)
				sembrarItemEsfuerzo(t, entorno, proyectoEsfuerzoUS14, itemEsfuerzoUS14, domain.EstadoItemPendiente)
				sembrarIntegranteEsfuerzo(t, entorno, proyectoEsfuerzoUS14, integranteEsfuerzoUS14)
				entorno.asignaciones.sprintDeErr = causa
			},
		},
		{
			descripcion: "obtener Sprint",
			preparar: func(t *testing.T, entorno *entornoEsfuerzo, causa error) {
				sembrarProyectoEsfuerzo(t, entorno, proyectoEsfuerzoUS14)
				sembrarItemEsfuerzo(t, entorno, proyectoEsfuerzoUS14, itemEsfuerzoUS14, domain.EstadoItemPendiente)
				sembrarIntegranteEsfuerzo(t, entorno, proyectoEsfuerzoUS14, integranteEsfuerzoUS14)
				sembrarSprintEsfuerzo(t, entorno, proyectoEsfuerzoUS14, 1, domain.EstadoSprintActivo)
				sembrarAsignacionEsfuerzo(t, entorno, proyectoEsfuerzoUS14, 1, 1)
				entorno.sprints.obtenerErr = causa
			},
		},
		{
			descripcion: "listar integrantes",
			preparar: func(t *testing.T, entorno *entornoEsfuerzo, causa error) {
				sembrarProyectoEsfuerzo(t, entorno, proyectoEsfuerzoUS14)
				sembrarItemEsfuerzo(t, entorno, proyectoEsfuerzoUS14, itemEsfuerzoUS14, domain.EstadoItemPendiente)
				entorno.integrantes.listarErr = causa
			},
		},
		{
			descripcion:     "agregar esfuerzo",
			agregarLlamadas: 1,
			preparar: func(t *testing.T, entorno *entornoEsfuerzo, causa error) {
				sembrarProyectoEsfuerzo(t, entorno, proyectoEsfuerzoUS14)
				sembrarItemEsfuerzo(t, entorno, proyectoEsfuerzoUS14, itemEsfuerzoUS14, domain.EstadoItemPendiente)
				sembrarIntegranteEsfuerzo(t, entorno, proyectoEsfuerzoUS14, integranteEsfuerzoUS14)
				entorno.esfuerzo.agregarErr = causa
			},
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			entorno := nuevoEntornoEsfuerzo()
			caso.preparar(t, entorno, causa)

			_, err := entorno.servicio().RegistrarEsfuerzo(
				proyectoEsfuerzoUS14,
				itemEsfuerzoUS14,
				integranteEsfuerzoUS14,
				fechaEsfuerzoUS14(2026, time.November, 15, 9, 0),
				"Análisis",
				1.5,
				hoyEsfuerzoUS14(),
			)
			verificarErrorRepositorioEsfuerzoUS14(t, err, causa)
			if entorno.esfuerzo.agregarLlamadas != caso.agregarLlamadas {
				t.Errorf("Agregar fue invocado %d veces, se esperaba %d", entorno.esfuerzo.agregarLlamadas, caso.agregarLlamadas)
			}
			if len(entorno.esfuerzo.registros) != 0 {
				t.Errorf("se persistió un registro pese al error de repositorio: %+v", entorno.esfuerzo.registros)
			}
		})
	}
}

// R6 / D24: los errores de esfuerzo pasan sin envolver por
// envolverErrorRepositorio.
func TestUS14_EnvolverErrorRepositorio_ErroresDeEsfuerzoSinEnvolver(t *testing.T) {
	casos := []struct {
		descripcion string
		err         error
	}{
		{descripcion: "integrante ajeno", err: domain.ErrIntegranteNoPerteneceAlProyecto},
		{descripcion: "horas inválidas", err: domain.ErrHorasEsfuerzoInvalidas},
		{descripcion: "fecha futura", err: domain.ErrFechaEsfuerzoFutura},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			obtenido := envolverErrorRepositorio("registrar esfuerzo", caso.err)
			if obtenido != caso.err {
				t.Errorf("envolverErrorRepositorio() = %v, se esperaba el mismo error sin envolver", obtenido)
			}
			if !errors.Is(obtenido, caso.err) {
				t.Errorf("errors.Is(obtenido, caso.err) = false; obtenido = %v", obtenido)
			}
		})
	}
}
