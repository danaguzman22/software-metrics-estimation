package service

import (
	"errors"
	"reflect"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// D36 (US-13 FASE 2, opción A): HistoriaActual devuelve la historia del
// flujo proyecto → historia + backlog.ObtenerItem (sin métodos nuevos de
// repositorio), para que la página de poker informe los Story Points
// vigentes en el encabezado.
func TestUS13_Fase2_ServicioPlanningPoker_HistoriaActualDevuelveLaHistoria(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	servicio := entorno.servicio()

	historia, err := servicio.HistoriaActual("p1", "HU-01")
	if err != nil {
		t.Fatalf("HistoriaActual() devolvió un error inesperado: %v", err)
	}
	esperada := domain.ItemBacklog{
		ID:          "HU-01",
		Titulo:      "Historia HU-01",
		Descripcion: "Descripción HU-01",
		Prioridad:   "Alta",
		Estado:      "Pendiente",
		StoryPoints: 5,
	}
	if !reflect.DeepEqual(historia, esperada) {
		t.Errorf("historia = %+v, se esperaba %+v", historia, esperada)
	}
	if entorno.fakeRepositorioBacklog.obtenerLlamadas != 1 {
		t.Errorf("ObtenerItem fue invocado %d veces, se esperaba 1", entorno.fakeRepositorioBacklog.obtenerLlamadas)
	}
}

// D36 (US-13 FASE 2): proyecto inexistente -> domain.ErrProyectoNoEncontrado,
// sin consultar la historia.
func TestUS13_Fase2_ServicioPlanningPoker_HistoriaActualProyectoNoEncontrado(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	servicio := entorno.servicio()

	historia, err := servicio.HistoriaActual("ausente", "HU-01")
	if err == nil {
		t.Fatal("HistoriaActual() con un proyecto inexistente debería devolver un error")
	}
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if !reflect.DeepEqual(historia, domain.ItemBacklog{}) {
		t.Errorf("historia = %+v, se esperaba domain.ItemBacklog{}", historia)
	}
	if entorno.fakeRepositorioBacklog.obtenerLlamadas != 0 {
		t.Errorf("ObtenerItem fue invocado %d veces; no debía consultarse la historia", entorno.fakeRepositorioBacklog.obtenerLlamadas)
	}
}

// D36 (US-13 FASE 2): historia inexistente -> domain.ErrItemNoEncontrado
// (pasa sin envolver); el proyecto sí se verifica antes.
func TestUS13_Fase2_ServicioPlanningPoker_HistoriaActualHistoriaNoEncontrada(t *testing.T) {
	entorno := nuevoEntornoPokerSembrado(t)
	servicio := entorno.servicio()

	historia, err := servicio.HistoriaActual("p1", "HU-99")
	if err == nil {
		t.Fatal("HistoriaActual() con una historia inexistente debería devolver un error")
	}
	if !errors.Is(err, domain.ErrItemNoEncontrado) {
		t.Errorf("err = %v, se esperaba %v", err, domain.ErrItemNoEncontrado)
	}
	if err != domain.ErrItemNoEncontrado {
		t.Errorf("err = %v, se esperaba exactamente %v sin envolver", err, domain.ErrItemNoEncontrado)
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("err = %v, no debería mapearse a ErrProyectoNoEncontrado", err)
	}
	if !reflect.DeepEqual(historia, domain.ItemBacklog{}) {
		t.Errorf("historia = %+v, se esperaba domain.ItemBacklog{}", historia)
	}
}
