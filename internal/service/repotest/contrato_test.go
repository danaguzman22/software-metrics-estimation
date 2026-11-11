package repotest

import "testing"

// Tests de contrato T-08: nombres conservados y firmas ejecutables válidas.
// La implementación de los casos de contrato está pendiente (fase RED).

func TestT08_Proyectos_RepositorioVacioYRoundTrip(t *testing.T) {
	t.Fatal("T-08 RED: caso de contrato pendiente de implementar")
}

func TestT08_Integrantes_OrdenYAislamiento(t *testing.T) {
	t.Fatal("T-08 RED: caso de contrato pendiente de implementar")
}

func TestT08_Backlog_NumeracionOrdenYAislamiento(t *testing.T) {
	t.Fatal("T-08 RED: caso de contrato pendiente de implementar")
}

func TestT08_Backlog_RoundTripYNoEncontrado(t *testing.T) {
	t.Fatal("T-08 RED: caso de contrato pendiente de implementar")
}

func TestT08_ErroresInfraestructura_EnvuelvenConW(t *testing.T) {
	t.Fatal("T-08 RED: caso de contrato pendiente de implementar")
}

// Subtests contractuales D10 (nombres conservados). Implementación pendiente.
func TestUS11_EstimarItem_BacklogSinArchivo(t *testing.T) {
	t.Fatal("T-08 RED: subtest contractual pendiente de implementar")
}

func TestUS11_EstimarItem_MismoIDEnOtroProyecto(t *testing.T) {
	t.Fatal("T-08 RED: subtest contractual pendiente de implementar")
}

// Helpers internos previstos para recibir entorno (nombres sin prefijo Test).
func ejecutarProyectos(t *testing.T, entorno Entorno)           { _ = t; _ = entorno }
func ejecutarIntegrantes(t *testing.T, entorno Entorno)         { _ = t; _ = entorno }
func ejecutarBacklogNumeracion(t *testing.T, entorno Entorno)   { _ = t; _ = entorno }
func ejecutarBacklogRoundTrip(t *testing.T, entorno Entorno)    { _ = t; _ = entorno }
func ejecutarErroresInfra(t *testing.T, entorno Entorno)        { _ = t; _ = entorno }
func ejecutarUS11BacklogSinArchivo(t *testing.T, entorno Entorno)     { _ = t; _ = entorno }
func ejecutarUS11MismoIDEnOtroProyecto(t *testing.T, entorno Entorno) { _ = t; _ = entorno }
