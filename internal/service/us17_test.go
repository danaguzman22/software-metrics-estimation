package service

import (
	"errors"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func horasPuntero(valor float64) *float64 {
	return &valor
}

func metricasIguales(t *testing.T, obtenidas, esperadas domain.MetricasEsfuerzo) {
	t.Helper()
	if obtenidas.HorasEstimadas != esperadas.HorasEstimadas {
		t.Errorf("HorasEstimadas = %v, se esperaba %v", obtenidas.HorasEstimadas, esperadas.HorasEstimadas)
	}
	if obtenidas.HorasReales != esperadas.HorasReales {
		t.Errorf("HorasReales = %v, se esperaba %v", obtenidas.HorasReales, esperadas.HorasReales)
	}
	if obtenidas.EstimacionesCompletas != esperadas.EstimacionesCompletas {
		t.Errorf("EstimacionesCompletas = %v, se esperaba %v", obtenidas.EstimacionesCompletas, esperadas.EstimacionesCompletas)
	}
	if (obtenidas.Desviacion == nil) != (esperadas.Desviacion == nil) {
		t.Fatalf("Desviacion = %v, se esperaba %v", obtenidas.Desviacion, esperadas.Desviacion)
	}
	if obtenidas.Desviacion != nil && *obtenidas.Desviacion != *esperadas.Desviacion {
		t.Errorf("Desviacion = %v, se esperaba %v", *obtenidas.Desviacion, *esperadas.Desviacion)
	}
}

func TestUS17_EstimarHorasReemplazaYPermiteCero(t *testing.T) {
	entorno := nuevoEntornoEsfuerzo()
	sembrarProyectoEsfuerzo(t, entorno, "p1")
	sembrarItemEsfuerzo(t, entorno, "p1", "HU-01", domain.EstadoItemPendiente)
	servicio := entorno.servicio()

	item, err := servicio.EstimarHoras("p1", "HU-01", 4.25)
	if err != nil {
		t.Fatalf("EstimarHoras() inesperadamente falló: %v", err)
	}
	if item.HorasEstimadas == nil || *item.HorasEstimadas != 4.25 {
		t.Fatalf("HorasEstimadas = %v, se esperaba 4.25", item.HorasEstimadas)
	}

	item, err = servicio.EstimarHoras("p1", "HU-01", 0)
	if err != nil {
		t.Fatalf("reemplazar la estimación por cero falló: %v", err)
	}
	if item.HorasEstimadas == nil || *item.HorasEstimadas != 0 {
		t.Fatalf("HorasEstimadas = %v, se esperaba cero explícito", item.HorasEstimadas)
	}
}

func TestUS17_EstimarHorasValidaProyectoHorasEHistoria(t *testing.T) {
	casos := []struct {
		nombre   string
		proyecto string
		itemID   string
		horas    float64
		preparar func(*testing.T, *entornoEsfuerzo)
		esperado error
	}{
		{
			nombre:   "proyecto inexistente",
			proyecto: "p999", itemID: "HU-01", horas: -1,
			esperado: domain.ErrProyectoNoEncontrado,
		},
		{
			nombre:   "horas inválidas",
			proyecto: "p1", itemID: "HU-01", horas: 1.234,
			preparar: func(t *testing.T, e *entornoEsfuerzo) { sembrarProyectoEsfuerzo(t, e, "p1") },
			esperado: domain.ErrHorasEstimadasInvalidas,
		},
		{
			nombre:   "historia inexistente",
			proyecto: "p1", itemID: "HU-99", horas: 1,
			preparar: func(t *testing.T, e *entornoEsfuerzo) { sembrarProyectoEsfuerzo(t, e, "p1") },
			esperado: domain.ErrItemNoEncontrado,
		},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			entorno := nuevoEntornoEsfuerzo()
			if caso.preparar != nil {
				caso.preparar(t, entorno)
			}
			_, err := entorno.servicio().EstimarHoras(caso.proyecto, caso.itemID, caso.horas)
			if err != caso.esperado {
				t.Fatalf("err = %v, se esperaba el error de dominio sin envolver %v", err, caso.esperado)
			}
			if entorno.backlog.actualizaciones != 0 {
				t.Errorf("no debía actualizarse la historia, actualizaciones=%d", entorno.backlog.actualizaciones)
			}
		})
	}
}

func TestUS17_ConsultasPropaganErroresDeDominioSinEnvolver(t *testing.T) {
	casos := []struct {
		nombre   string
		consulta func(*ServicioEsfuerzo) error
		preparar func(*testing.T, *entornoEsfuerzo)
		esperado error
	}{
		{
			nombre: "proyecto inexistente",
			consulta: func(servicio *ServicioEsfuerzo) error {
				_, err := servicio.MetricasProyecto("p999")
				return err
			},
			esperado: domain.ErrProyectoNoEncontrado,
		},
		{
			nombre: "historia inexistente",
			preparar: func(t *testing.T, entorno *entornoEsfuerzo) {
				sembrarProyectoEsfuerzo(t, entorno, "p1")
			},
			consulta: func(servicio *ServicioEsfuerzo) error {
				_, err := servicio.MetricasHistoria("p1", "HU-99")
				return err
			},
			esperado: domain.ErrItemNoEncontrado,
		},
		{
			nombre: "sprint inexistente",
			preparar: func(t *testing.T, entorno *entornoEsfuerzo) {
				sembrarProyectoEsfuerzo(t, entorno, "p1")
			},
			consulta: func(servicio *ServicioEsfuerzo) error {
				_, err := servicio.MetricasSprint("p1", 99)
				return err
			},
			esperado: domain.ErrSprintNoEncontrado,
		},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			entorno := nuevoEntornoEsfuerzo()
			if caso.preparar != nil {
				caso.preparar(t, entorno)
			}
			if err := caso.consulta(entorno.servicio()); err != caso.esperado {
				t.Fatalf("err = %v, se esperaba el error de dominio sin envolver %v", err, caso.esperado)
			}
		})
	}
}

func TestUS17_ConsultasEnvuelvenErroresDeInfraestructura(t *testing.T) {
	causa := errors.New("fallo de base de datos")
	casos := []struct {
		nombre   string
		preparar func(*testing.T, *entornoEsfuerzo)
		consulta func(*ServicioEsfuerzo) error
	}{
		{
			nombre: "horas reales de historia",
			preparar: func(t *testing.T, entorno *entornoEsfuerzo) {
				sembrarProyectoEsfuerzo(t, entorno, "p1")
				sembrarItemEsfuerzo(t, entorno, "p1", "HU-01", domain.EstadoItemPendiente)
				entorno.esfuerzo.horasItemErr = causa
			},
			consulta: func(servicio *ServicioEsfuerzo) error {
				_, err := servicio.MetricasHistoria("p1", "HU-01")
				return err
			},
		},
		{
			nombre: "horas reales de sprint",
			preparar: func(t *testing.T, entorno *entornoEsfuerzo) {
				sembrarProyectoEsfuerzo(t, entorno, "p1")
				sembrarSprintEsfuerzo(t, entorno, "p1", 1, domain.EstadoSprintActivo)
				entorno.esfuerzo.horasSprintErr = causa
			},
			consulta: func(servicio *ServicioEsfuerzo) error {
				_, err := servicio.MetricasSprint("p1", 1)
				return err
			},
		},
		{
			nombre: "horas reales de proyecto",
			preparar: func(t *testing.T, entorno *entornoEsfuerzo) {
				sembrarProyectoEsfuerzo(t, entorno, "p1")
				entorno.esfuerzo.horasProyectoErr = causa
			},
			consulta: func(servicio *ServicioEsfuerzo) error {
				_, err := servicio.MetricasProyecto("p1")
				return err
			},
		},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			entorno := nuevoEntornoEsfuerzo()
			caso.preparar(t, entorno)
			if err := caso.consulta(entorno.servicio()); !errors.Is(err, causa) {
				t.Fatalf("err = %v, se esperaba que envolviera la causa %v", err, causa)
			}
		})
	}
}

func TestUS17_MetricasHistoriaSumaRealesYCalculaDesviacion(t *testing.T) {
	entorno := nuevoEntornoEsfuerzo()
	sembrarProyectoEsfuerzo(t, entorno, "p1")
	sembrarItemEsfuerzo(t, entorno, "p1", "HU-01", domain.EstadoItemPendiente)
	entorno.backlog.items["p1"][0].HorasEstimadas = horasPuntero(4)
	entorno.esfuerzo.registros = []registroEsfuerzoFake{
		{proyectoID: "p1", itemNumero: 1, registro: domain.RegistroEsfuerzo{Horas: 1.5}},
		{proyectoID: "p1", itemNumero: 1, registro: domain.RegistroEsfuerzo{Horas: 3.5}},
	}

	obtenidas, err := entorno.servicio().MetricasHistoria("p1", "HU-01")
	if err != nil {
		t.Fatalf("MetricasHistoria() devolvió un error: %v", err)
	}
	metricasIguales(t, obtenidas, domain.MetricasEsfuerzo{
		HorasEstimadas: 4, HorasReales: 5, EstimacionesCompletas: true, Desviacion: horasPuntero(25),
	})
}

func TestUS17_MetricasHistoriaSinEstimacionDevuelveRealesSinDesviacion(t *testing.T) {
	entorno := nuevoEntornoEsfuerzo()
	sembrarProyectoEsfuerzo(t, entorno, "p1")
	sembrarItemEsfuerzo(t, entorno, "p1", "HU-01", domain.EstadoItemPendiente)
	entorno.esfuerzo.registros = []registroEsfuerzoFake{
		{proyectoID: "p1", itemNumero: 1, registro: domain.RegistroEsfuerzo{Horas: 2}},
	}

	obtenidas, err := entorno.servicio().MetricasHistoria("p1", "HU-01")
	if err != nil {
		t.Fatalf("MetricasHistoria() devolvió un error: %v", err)
	}
	metricasIguales(t, obtenidas, domain.MetricasEsfuerzo{
		HorasReales: 2, EstimacionesCompletas: false,
	})
}

func TestUS17_MetricasProyectoCalculaDesviacionSobreTotales(t *testing.T) {
	entorno := nuevoEntornoEsfuerzo()
	sembrarProyectoEsfuerzo(t, entorno, "p1")
	sembrarItemEsfuerzo(t, entorno, "p1", "HU-01", domain.EstadoItemPendiente)
	sembrarItemEsfuerzo(t, entorno, "p1", "HU-02", domain.EstadoItemPendiente)
	entorno.backlog.items["p1"][0].HorasEstimadas = horasPuntero(1)
	entorno.backlog.items["p1"][1].HorasEstimadas = horasPuntero(9)
	entorno.esfuerzo.registros = []registroEsfuerzoFake{
		{proyectoID: "p1", itemNumero: 1, registro: domain.RegistroEsfuerzo{Horas: 2}},
		{proyectoID: "p1", itemNumero: 2, registro: domain.RegistroEsfuerzo{Horas: 9}},
	}

	obtenidas, err := entorno.servicio().MetricasProyecto("p1")
	if err != nil {
		t.Fatalf("MetricasProyecto() devolvió un error: %v", err)
	}
	metricasIguales(t, obtenidas, domain.MetricasEsfuerzo{
		HorasEstimadas: 10, HorasReales: 11, EstimacionesCompletas: true, Desviacion: horasPuntero(10),
	})
}

func TestUS17_MetricasSprintUsaSnapshotCerrado(t *testing.T) {
	entorno := nuevoEntornoEsfuerzo()
	sembrarProyectoEsfuerzo(t, entorno, "p1")
	sembrarSprintEsfuerzo(t, entorno, "p1", 1, domain.EstadoSprintCerrado)
	entorno.esfuerzo.snapshots["p1:1"] = snapshotEsfuerzoFake{
		disponible: true,
		items: []domain.ItemSprintEsfuerzo{
			{ItemNumero: 1, HorasEstimadas: horasPuntero(10)},
		},
	}
	entorno.esfuerzo.registros = []registroEsfuerzoFake{
		{proyectoID: "p1", itemNumero: 1, sprintNum: intPuntero(1), registro: domain.RegistroEsfuerzo{Horas: 12}},
		{proyectoID: "p1", itemNumero: 2, sprintNum: nil, registro: domain.RegistroEsfuerzo{Horas: 100}},
	}

	obtenidas, err := entorno.servicio().MetricasSprint("p1", 1)
	if err != nil {
		t.Fatalf("MetricasSprint() devolvió un error: %v", err)
	}
	metricasIguales(t, obtenidas, domain.MetricasEsfuerzo{
		HorasEstimadas: 10, HorasReales: 12, EstimacionesCompletas: true, Desviacion: horasPuntero(20),
	})
}

func TestUS17_MetricasSprintCerradoSinSnapshotNoReconstruye(t *testing.T) {
	entorno := nuevoEntornoEsfuerzo()
	sembrarProyectoEsfuerzo(t, entorno, "p1")
	sembrarSprintEsfuerzo(t, entorno, "p1", 1, domain.EstadoSprintCerrado)
	sembrarItemEsfuerzo(t, entorno, "p1", "HU-01", domain.EstadoItemPendiente)
	sembrarAsignacionEsfuerzo(t, entorno, "p1", 1, 1)

	_, err := entorno.servicio().MetricasSprint("p1", 1)
	if !errors.Is(err, domain.ErrMetricasSprintNoDisponibles) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrMetricasSprintNoDisponibles)
	}
}

func TestUS17_RegistrarEsfuerzoAtribuyeSprintVigente(t *testing.T) {
	entorno := entornoEsfuerzoValido(t)
	sembrarSprintEsfuerzo(t, entorno, "p1", 1, domain.EstadoSprintActivo)
	sembrarAsignacionEsfuerzo(t, entorno, "p1", 1, 1)
	_, err := entorno.servicio().RegistrarEsfuerzo(
		"p1", "HU-01", "Ana Pérez",
		time.Date(2026, time.November, 15, 9, 0, 0, 0, time.UTC),
		"Análisis", 1.5,
		time.Date(2026, time.November, 15, 9, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("RegistrarEsfuerzo() devolvió un error: %v", err)
	}
	if entorno.esfuerzo.registros[0].sprintNum == nil || *entorno.esfuerzo.registros[0].sprintNum != 1 {
		t.Fatalf("Sprint atribuido = %v, se esperaba 1", entorno.esfuerzo.registros[0].sprintNum)
	}
}

func intPuntero(valor int) *int {
	return &valor
}
