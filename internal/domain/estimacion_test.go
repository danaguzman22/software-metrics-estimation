package domain_test

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func itemUS11(storyPoints int) domain.ItemBacklog {
	return domain.ItemBacklog{
		ID:          "HU-01",
		Titulo:      "Registrar proyecto",
		Descripcion: "Permitir dar de alta proyectos en el sistema",
		Prioridad:   "Alta",
		Estado:      "Pendiente",
		StoryPoints: storyPoints,
		Criterios: []string{
			"El nombre del proyecto no puede estar vacío",
			"El proyecto queda disponible para asignar integrantes",
		},
	}
}

func verificarItemUS11(t *testing.T, obtenido, esperado domain.ItemBacklog) {
	t.Helper()
	if obtenido.ID != esperado.ID {
		t.Errorf("ID = %q, se esperaba %q", obtenido.ID, esperado.ID)
	}
	if obtenido.Titulo != esperado.Titulo {
		t.Errorf("Titulo = %q, se esperaba %q", obtenido.Titulo, esperado.Titulo)
	}
	if obtenido.Descripcion != esperado.Descripcion {
		t.Errorf("Descripcion = %q, se esperaba %q", obtenido.Descripcion, esperado.Descripcion)
	}
	if obtenido.Prioridad != esperado.Prioridad {
		t.Errorf("Prioridad = %q, se esperaba %q", obtenido.Prioridad, esperado.Prioridad)
	}
	if obtenido.Estado != esperado.Estado {
		t.Errorf("Estado = %q, se esperaba %q", obtenido.Estado, esperado.Estado)
	}
	if obtenido.StoryPoints != esperado.StoryPoints {
		t.Errorf("StoryPoints = %d, se esperaba %d", obtenido.StoryPoints, esperado.StoryPoints)
	}
	criteriosObtenidos := obtenido.Criterios
	if criteriosObtenidos == nil {
		criteriosObtenidos = []string{}
	}
	criteriosEsperados := esperado.Criterios
	if criteriosEsperados == nil {
		criteriosEsperados = []string{}
	}
	if !reflect.DeepEqual(criteriosObtenidos, criteriosEsperados) {
		t.Errorf("Criterios = %v, se esperaba %v", obtenido.Criterios, esperado.Criterios)
	}
}

// CA1 / R1 / R3: Estimar con un valor válido devuelve una copia con ese valor.
func TestUS11_Estimar_AsignaStoryPoints(t *testing.T) {
	for _, storyPoints := range []int{1, 2, 3, 5, 8, 13, 21} {
		t.Run(strconv.Itoa(storyPoints), func(t *testing.T) {
			estimado, err := itemUS11(0).Estimar(storyPoints)
			if err != nil {
				t.Fatalf("Estimar(%d) devolvió un error inesperado: %v", storyPoints, err)
			}
			if estimado.StoryPoints != storyPoints {
				t.Errorf("StoryPoints = %d, se esperaba %d", estimado.StoryPoints, storyPoints)
			}
		})
	}
}

// CA1 / R1 / R3: Estimar conserva ID, título, descripción, prioridad, estado y
// criterios del ítem.
func TestUS11_Estimar_PreservaRestoDelItem(t *testing.T) {
	original := itemUS11(0)
	esperado := original
	esperado.StoryPoints = 5

	estimado, err := original.Estimar(5)
	if err != nil {
		t.Fatalf("Estimar(5) devolvió un error inesperado: %v", err)
	}
	verificarItemUS11(t, estimado, esperado)
}

// CA1 / R3: Estimar no modifica el ítem receptor; el original conserva sus
// Story Points previos.
func TestUS11_Estimar_NoModificaElItemOriginal(t *testing.T) {
	original := itemUS11(3)
	esperadoOriginal := original

	_, err := original.Estimar(8)
	if err != nil {
		t.Fatalf("Estimar(8) devolvió un error inesperado: %v", err)
	}
	verificarItemUS11(t, original, esperadoOriginal)
}

// CA1 / CA2 / R2: los siete valores de la escala Fibonacci son válidos.
func TestUS11_EsStoryPointsValido_EscalaValida(t *testing.T) {
	for _, storyPoints := range []int{1, 2, 3, 5, 8, 13, 21} {
		if !domain.EsStoryPointsValido(storyPoints) {
			t.Errorf("EsStoryPointsValido(%d) = false, se esperaba true", storyPoints)
		}
	}
}

// CA2 / R2: los valores fuera de la escala no son válidos.
func TestUS11_EsStoryPointsValido_ValoresInvalidos(t *testing.T) {
	for _, storyPoints := range []int{0, -1, -5, 4, 6, 7, 9, 10, 11, 12, 14, 20, 22, 100} {
		if domain.EsStoryPointsValido(storyPoints) {
			t.Errorf("EsStoryPointsValido(%d) = true, se esperaba false", storyPoints)
		}
	}
}

// CA2 / R3: Estimar con un valor inválido devuelve ErrStoryPointsInvalidos y
// domain.ItemBacklog{}.
func TestUS11_Estimar_ValorInvalidoDevuelveError(t *testing.T) {
	for _, storyPoints := range []int{0, -1, 4, 22} {
		t.Run(strconv.Itoa(storyPoints), func(t *testing.T) {
			estimado, err := itemUS11(3).Estimar(storyPoints)
			if !errors.Is(err, domain.ErrStoryPointsInvalidos) {
				t.Errorf("err = %v, se esperaba %v", err, domain.ErrStoryPointsInvalidos)
			}
			verificarItemUS11(t, estimado, domain.ItemBacklog{})
		})
	}
}

// CA1 / R4 / R3: re-estimar un ítem con 3 usando 8 devuelve una copia con 8.
func TestUS11_Estimar_ReestimacionReemplazaValor(t *testing.T) {
	original := itemUS11(3)
	esperado := original
	esperado.StoryPoints = 8

	estimado, err := original.Estimar(8)
	if err != nil {
		t.Fatalf("Estimar(8) devolvió un error inesperado: %v", err)
	}
	verificarItemUS11(t, estimado, esperado)
}
