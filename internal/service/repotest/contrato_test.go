package repotest

import "testing"

// Tests de contrato T-08: conservan nombres; delegan en la suite compartida.
// Se omiten en ejecución directa porque requieren un Entorno de adaptador.
func skipSuiteDirecta(t *testing.T) {
	t.Helper()
	t.Skip("T-08: usar repotest.Ejecutar desde el paquete del adaptador")
}

func TestT08_Proyectos_RepositorioVacioYRoundTrip(t *testing.T) {
	skipSuiteDirecta(t)
}

func TestT08_Integrantes_OrdenYAislamiento(t *testing.T) {
	skipSuiteDirecta(t)
}

func TestT08_Backlog_NumeracionOrdenYAislamiento(t *testing.T) {
	skipSuiteDirecta(t)
}

func TestT08_Backlog_RoundTripYNoEncontrado(t *testing.T) {
	skipSuiteDirecta(t)
}

func TestT08_ErroresInfraestructura_EnvuelvenConW(t *testing.T) {
	skipSuiteDirecta(t)
}

// Subtests contractuales D10 (nombres conservados); cobertura real vía Ejecutar.
func TestUS11_EstimarItem_BacklogSinArchivo(t *testing.T) {
	skipSuiteDirecta(t)
}

func TestUS11_EstimarItem_MismoIDEnOtroProyecto(t *testing.T) {
	skipSuiteDirecta(t)
}
