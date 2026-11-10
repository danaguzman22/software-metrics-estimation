package service

import (
	"errors"
	"fmt"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage"
)

// verificarProyectoExistente carga la colección "proyectos" y comprueba que
// exista el proyecto identificado por proyectoID.
//
// Un archivo inexistente se interpreta como colección vacía, de modo que el
// proyecto no se encuentra. Devuelve domain.ErrProyectoNoEncontrado si el
// proyecto no existe. Los errores de almacenamiento distintos de "no
// encontrado" se propagan envueltos con %w y no se mapean a
// domain.ErrProyectoNoEncontrado.
func verificarProyectoExistente(almacen almacen, proyectoID string) error {
	var proyectos []proyectoRegistrado
	if err := almacen.Load("proyectos", &proyectos); err != nil {
		if !errors.Is(err, storage.ErrNotFound) {
			return fmt.Errorf("service: cargar proyectos: %w", err)
		}
		// Un archivo inexistente se interpreta como colección vacía.
		proyectos = nil
	}

	for _, registro := range proyectos {
		if registro.ID == proyectoID {
			return nil
		}
	}
	return domain.ErrProyectoNoEncontrado
}
