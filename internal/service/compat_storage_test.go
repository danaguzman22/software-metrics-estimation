package service

// Helpers de armado de los tests de servicio sobre los fakes en memoria de
// repos_fake_test.go. La antigua compatibilidad con storage.JSONStore se retiró
// al migrar ServicioBacklog en la Fase 2 de T-09.

import "testing"

func sembrarProyectos(t *testing.T, entorno *entornoBacklog, registros []ProyectoRegistrado) {
	t.Helper()
	for _, registro := range registros {
		entorno.proyectos[registro.ID] = registro.Proyecto
	}
}
