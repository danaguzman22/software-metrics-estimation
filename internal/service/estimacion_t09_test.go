package service

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func nuevoEntornoEstimacionConItem(storyPoints int) *entornoBacklog {
	entorno := nuevoEntornoBacklog()
	entorno.proyectos["p1"] = proyectoFakeValido()
	entorno.items["p1"] = []domain.ItemBacklog{itemUS11Backlog(storyPoints)}
	return entorno
}

func TestT09_EstimarItem_PropagaErrorRepositorioConWrap(t *testing.T) {
	causa := errors.New("fallo de repositorio")

	t.Run("Existe", func(t *testing.T) {
		entorno := nuevoEntornoBacklog()
		entorno.existeErr = causa
		servicio := entorno.servicio()

		item, err := servicio.EstimarItem("p1", "HU-01", 5)
		if !errors.Is(err, causa) {
			t.Fatalf("errors.Is(err, causa) = false; err = %v", err)
		}
		if err == causa {
			t.Fatalf("el error debe estar envuelto; err == causa")
		}
		if !strings.HasPrefix(err.Error(), "service:") {
			t.Fatalf("mensaje = %q; se esperaba prefijo service:", err.Error())
		}
		verificarItemBacklogVacio(t, item)
		if entorno.fakeRepositorioBacklog.obtenerLlamadas != 0 || entorno.fakeRepositorioBacklog.actualizaciones != 0 {
			t.Fatalf("se accedió al backlog tras fallar la verificación del proyecto")
		}
	})

	t.Run("ObtenerItem", func(t *testing.T) {
		entorno := nuevoEntornoEstimacionConItem(0)
		entorno.fakeRepositorioBacklog.obtenerErr = causa
		servicio := entorno.servicio()

		item, err := servicio.EstimarItem("p1", "HU-01", 5)
		if !errors.Is(err, causa) {
			t.Fatalf("errors.Is(err, causa) = false; err = %v", err)
		}
		if err == causa {
			t.Fatalf("el error debe estar envuelto; err == causa")
		}
		if !strings.HasPrefix(err.Error(), "service:") {
			t.Fatalf("mensaje = %q; se esperaba prefijo service:", err.Error())
		}
		verificarItemBacklogVacio(t, item)
		if entorno.fakeRepositorioBacklog.actualizaciones != 0 {
			t.Fatalf("ActualizarItem fue invocado %d veces", entorno.fakeRepositorioBacklog.actualizaciones)
		}
	})

	t.Run("ActualizarItem", func(t *testing.T) {
		entorno := nuevoEntornoEstimacionConItem(0)
		entorno.fakeRepositorioBacklog.actualizarErr = causa
		servicio := entorno.servicio()

		item, err := servicio.EstimarItem("p1", "HU-01", 5)
		if !errors.Is(err, causa) {
			t.Fatalf("errors.Is(err, causa) = false; err = %v", err)
		}
		if err == causa {
			t.Fatalf("el error debe estar envuelto; err == causa")
		}
		if !strings.HasPrefix(err.Error(), "service:") {
			t.Fatalf("mensaje = %q; se esperaba prefijo service:", err.Error())
		}
		verificarItemBacklogVacio(t, item)
	})
}

func TestT09_EstimarItem_NoEscribeSiProyectoNoExiste(t *testing.T) {
	entorno := nuevoEntornoBacklog()
	servicio := entorno.servicio()

	estimado, err := servicio.EstimarItem("ausente", "HU-01", 5)
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("error = %v; se esperaba ErrProyectoNoEncontrado", err)
	}
	verificarItemBacklogVacio(t, estimado)
	if entorno.fakeRepositorioBacklog.obtenerLlamadas != 0 || entorno.fakeRepositorioBacklog.actualizaciones != 0 {
		t.Fatalf("se accedió al backlog tras fallar la verificación del proyecto")
	}
}

func TestT09_EstimarItem_NoEscribeSiItemNoExiste(t *testing.T) {
	entorno := nuevoEntornoBacklog()
	entorno.proyectos["p1"] = proyectoFakeValido()
	servicio := entorno.servicio()

	estimado, err := servicio.EstimarItem("p1", "HU-02", 5)
	if !errors.Is(err, domain.ErrItemNoEncontrado) {
		t.Fatalf("error = %v; se esperaba ErrItemNoEncontrado", err)
	}
	verificarItemBacklogVacio(t, estimado)
	if entorno.fakeRepositorioBacklog.actualizaciones != 0 {
		t.Fatalf("ActualizarItem fue invocado %d veces", entorno.fakeRepositorioBacklog.actualizaciones)
	}
}

func TestT09_EstimarItem_NoEscribeSiValidacionFalla(t *testing.T) {
	for _, storyPoints := range []int{0, -1, 4, 22} {
		t.Run(strconv.Itoa(storyPoints), func(t *testing.T) {
			entorno := nuevoEntornoEstimacionConItem(3)
			servicio := entorno.servicio()

			estimado, err := servicio.EstimarItem("p1", "HU-01", storyPoints)
			if !errors.Is(err, domain.ErrStoryPointsInvalidos) {
				t.Fatalf("error = %v; se esperaba ErrStoryPointsInvalidos", err)
			}
			verificarItemBacklogVacio(t, estimado)
			if entorno.fakeRepositorioBacklog.actualizaciones != 0 {
				t.Fatalf("ActualizarItem fue invocado %d veces", entorno.fakeRepositorioBacklog.actualizaciones)
			}
		})
	}
}

func TestT09_EstimarItem_ErrorDominioSinEnvolver(t *testing.T) {
	casos := []struct {
		descripcion   string
		entorno       func() *entornoBacklog
		proyectoID    string
		itemID        string
		storyPoints   int
		errorEsperado error
	}{
		{
			descripcion:   "proyecto inexistente",
			entorno:       nuevoEntornoBacklogConProyecto,
			proyectoID:    "ausente",
			itemID:        "HU-01",
			storyPoints:   5,
			errorEsperado: domain.ErrProyectoNoEncontrado,
		},
		{
			descripcion: "ítem inexistente",
			entorno: func() *entornoBacklog {
				entorno := nuevoEntornoBacklog()
				entorno.proyectos["p1"] = proyectoFakeValido()
				return entorno
			},
			proyectoID:    "p1",
			itemID:        "HU-02",
			storyPoints:   5,
			errorEsperado: domain.ErrItemNoEncontrado,
		},
		{
			descripcion:   "story points inválidos",
			entorno:       func() *entornoBacklog { return nuevoEntornoEstimacionConItem(3) },
			proyectoID:    "p1",
			itemID:        "HU-01",
			storyPoints:   4,
			errorEsperado: domain.ErrStoryPointsInvalidos,
		},
	}

	for _, caso := range casos {
		t.Run(caso.descripcion, func(t *testing.T) {
			servicio := caso.entorno().servicio()

			_, err := servicio.EstimarItem(caso.proyectoID, caso.itemID, caso.storyPoints)
			if !errors.Is(err, caso.errorEsperado) {
				t.Fatalf("error = %v; se esperaba %v", err, caso.errorEsperado)
			}
			if errors.Unwrap(err) != nil {
				t.Fatalf("error de dominio envuelto: %v", err)
			}
		})
	}
}
