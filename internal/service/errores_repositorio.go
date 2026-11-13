package service

import (
	"errors"
	"fmt"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func envolverErrorRepositorio(contexto string, err error) error {
	if errors.Is(err, domain.ErrProyectoNoEncontrado) ||
		errors.Is(err, domain.ErrIntegranteDuplicado) ||
		errors.Is(err, domain.ErrItemNoEncontrado) ||
		errors.Is(err, domain.ErrSprintActivoExistente) ||
		errors.Is(err, domain.ErrSprintDuplicado) {
		return err
	}
	return fmt.Errorf("service: %s: %w", contexto, err)
}
