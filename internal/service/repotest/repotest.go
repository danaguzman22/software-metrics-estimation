package repotest

import (
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/service"
)

// Entorno define la API mínima de la suite de contrato compartida (T-08).
// Superficie pendiente de implementar: sin lógica de contrato en fase RED.
type Entorno struct {
	Proyectos       service.RepositorioProyectos
	Integrantes     service.RepositorioIntegrantes
	Backlog         service.RepositorioBacklog
	SembrarProyecto func(t *testing.T, p service.ProyectoRegistrado)
	InyectarFallo   func(t *testing.T) (causa error, ejecutar func() error)
}

// Ejecutar ejecuta la suite de contrato compartida para un adaptador.
// Implementación pendiente (RED): la suite/contrato se implementará en GREEN.
func Ejecutar(t *testing.T, nuevo func(t *testing.T) Entorno) {
	t.Helper()
	t.Fatal("T-08 RED: implementación de la suite de contrato repotest pendiente")
}
