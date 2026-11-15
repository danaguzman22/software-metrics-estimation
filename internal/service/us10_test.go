package service

import (
	"errors"
	"reflect"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func TestUS10_ListarSprintsCerradosFiltraYConservaOrden(t *testing.T) {
	entorno := nuevoEntornoSprints()
	sembrarProyectosSprint(t, entorno, "p1")
	entorno.sprints["p1"] = []domain.Sprint{
		sprintUS07Servicio(1, domain.EstadoSprintCerrado),
		sprintUS07Servicio(2, domain.EstadoSprintActivo),
		sprintUS07Servicio(3, domain.EstadoSprintCerrado),
	}

	obtenidos, err := entorno.servicio().ListarSprintsCerrados("p1")
	if err != nil {
		t.Fatalf("ListarSprintsCerrados() devolvió un error: %v", err)
	}
	if len(obtenidos) != 2 {
		t.Fatalf("se obtuvieron %d Sprints, se esperaban 2", len(obtenidos))
	}
	if numeros := []int{obtenidos[0].Numero, obtenidos[1].Numero}; !reflect.DeepEqual(numeros, []int{1, 3}) {
		t.Errorf("números = %v, se esperaban [1 3]", numeros)
	}
	for _, sprint := range obtenidos {
		if sprint.Estado != domain.EstadoSprintCerrado {
			t.Errorf("se devolvió Sprint con estado %q", sprint.Estado)
		}
	}
	if entorno.fakeRepositorioSprints.listarLlamadas != 1 {
		t.Errorf("ListarPorProyecto se llamó %d veces, se esperaba 1", entorno.fakeRepositorioSprints.listarLlamadas)
	}
}

func TestUS10_ListarSprintsCerradosSinCerradosDevuelveListaVacia(t *testing.T) {
	entorno := nuevoEntornoSprints()
	sembrarProyectosSprint(t, entorno, "p1")
	entorno.sprints["p1"] = []domain.Sprint{
		sprintUS07Servicio(1, domain.EstadoSprintActivo),
	}

	obtenidos, err := entorno.servicio().ListarSprintsCerrados("p1")
	if err != nil {
		t.Fatalf("ListarSprintsCerrados() devolvió un error: %v", err)
	}
	if obtenidos == nil || len(obtenidos) != 0 {
		t.Fatalf("se esperaba una lista vacía no nula, se obtuvo %#v", obtenidos)
	}
}

func TestUS10_ListarSprintsCerradosPropagaErrores(t *testing.T) {
	t.Run("proyecto inexistente", func(t *testing.T) {
		entorno := nuevoEntornoSprints()
		_, err := entorno.servicio().ListarSprintsCerrados("p1")
		if err != domain.ErrProyectoNoEncontrado {
			t.Fatalf("err = %v, se esperaba %v sin envolver", err, domain.ErrProyectoNoEncontrado)
		}
		if entorno.fakeRepositorioSprints.listarLlamadas != 0 {
			t.Error("se consultaron los Sprints de un proyecto inexistente")
		}
	})

	t.Run("error de repositorio", func(t *testing.T) {
		entorno := nuevoEntornoSprints()
		sembrarProyectosSprint(t, entorno, "p1")
		causa := errors.New("fallo de persistencia")
		entorno.fakeRepositorioSprints.listarErr = causa

		_, err := entorno.servicio().ListarSprintsCerrados("p1")
		if !errors.Is(err, causa) {
			t.Fatalf("err = %v, se esperaba que envolviera %v", err, causa)
		}
	})
}
