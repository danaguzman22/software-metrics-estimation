package service

import (
	"errors"
	"math"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func TestUS16_VelocidadPromedioSoloSprintsCerrados(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	sembrarProyectosUS08(t, entorno, "p1")
	entorno.sprints.sprints["p1"] = []domain.Sprint{
		{Numero: 1, Estado: domain.EstadoSprintCerrado},
		{Numero: 2, Estado: domain.EstadoSprintCerrado},
		{Numero: 3, Estado: domain.EstadoSprintCerrado},
		{Numero: 4, Estado: domain.EstadoSprintActivo},
	}
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-01", 0, domain.EstadoCompletada))
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-02", 0, domain.EstadoCompletada))
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-03", 1, domain.EstadoCompletada))
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-04", 9, domain.EstadoCompletada))
	entorno.asignaciones.sembrar("p1", 1, 1)
	entorno.asignaciones.sembrar("p1", 2, 2)
	entorno.asignaciones.sembrar("p1", 3, 3)
	entorno.asignaciones.sembrar("p1", 4, 4)

	promedio, err := entorno.servicio().VelocidadPromedio("p1")
	if err != nil {
		t.Fatalf("VelocidadPromedio() devolvió un error inesperado: %v", err)
	}
	if math.Abs(promedio-(1.0/3.0)) > 1e-12 {
		t.Errorf("VelocidadPromedio() = %.12f, se esperaba %.12f", promedio, 1.0/3.0)
	}
}

func TestUS16_VelocidadPromedioSinSprintsCerradosEsCeroSinError(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	sembrarProyectosUS08(t, entorno, "p1")
	entorno.sprints.sprints["p1"] = []domain.Sprint{
		{Numero: 1, Estado: domain.EstadoSprintActivo},
	}

	promedio, err := entorno.servicio().VelocidadPromedio("p1")
	if err != nil {
		t.Fatalf("VelocidadPromedio() devolvió un error inesperado: %v", err)
	}
	if promedio != 0 {
		t.Errorf("VelocidadPromedio() = %v, se esperaba 0", promedio)
	}
}

func TestUS16_VelocidadPromedioValidaProyectoAntesDeListarSprints(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()

	_, err := entorno.servicio().VelocidadPromedio("p999")
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("err = %v, se esperaba %v", err, domain.ErrProyectoNoEncontrado)
	}
	if entorno.sprints.listarLlamadas != 0 {
		t.Error("se listaron Sprints antes de validar el proyecto")
	}
}

func TestUS16_VelocidadPromedioPropagaErrorAlListarSprints(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	sembrarProyectosUS08(t, entorno, "p1")
	entorno.sprints.listarErr = errors.New("fallo de repositorio")

	_, err := entorno.servicio().VelocidadPromedio("p1")
	if err == nil {
		t.Fatal("VelocidadPromedio() no propagó el error de ListarPorProyecto")
	}
}

func TestUS16_PorcentajeHistoriasCompletadasBacklogVacioEsCeroSinError(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	sembrarProyectosUS08(t, entorno, "p1")

	porcentaje, err := entorno.servicio().PorcentajeHistoriasCompletadas("p1")
	if err != nil {
		t.Fatalf("PorcentajeHistoriasCompletadas() devolvió un error inesperado: %v", err)
	}
	if porcentaje != 0 {
		t.Errorf("PorcentajeHistoriasCompletadas() = %v, se esperaba 0", porcentaje)
	}
}

func TestUS16_PorcentajeHistoriasCompletadasUsaTodoElBacklog(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	sembrarProyectosUS08(t, entorno, "p1")
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-01", 3, domain.EstadoCompletada))
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-02", 5, domain.EstadoItemPendiente))
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-03", 8, domain.EstadoItemPendiente))
	sembrarItemUS08(t, entorno, "p1", itemUS08Servicio("HU-04", 13, domain.EstadoItemPendiente))

	porcentaje, err := entorno.servicio().PorcentajeHistoriasCompletadas("p1")
	if err != nil {
		t.Fatalf("PorcentajeHistoriasCompletadas() devolvió un error inesperado: %v", err)
	}
	if porcentaje != 25 {
		t.Errorf("PorcentajeHistoriasCompletadas() = %v, se esperaba 25", porcentaje)
	}
}

func TestUS16_StoryPointsPlanificadosSprintCerradoSinValorGuardadoDevuelveError(t *testing.T) {
	entorno := nuevoEntornoAsignaciones()
	sembrarProyectosUS08(t, entorno, "p1")
	entorno.sprints.sprints["p1"] = []domain.Sprint{
		{Numero: 1, Estado: domain.EstadoSprintCerrado},
	}

	_, err := entorno.servicio().StoryPointsPlanificados("p1", 1)
	if err == nil {
		t.Fatal("StoryPointsPlanificados() no informó que faltan SPPlanificados en el Sprint cerrado")
	}
}
