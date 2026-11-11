package service

import (
	"errors"
	"fmt"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// AgregarIntegrante registra un integrante con su nombre y su rol en el
// proyecto identificado por proyectoID.
func (s *ServicioProyectos) AgregarIntegrante(
	proyectoID, nombre, rol string,
) (domain.Integrante, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		if errors.Is(err, domain.ErrProyectoNoEncontrado) || errors.Unwrap(err) == nil {
			return domain.Integrante{}, err
		}
		return domain.Integrante{}, fmt.Errorf("service: verificar proyecto: %w", err)
	}
	if !existe {
		return domain.Integrante{}, domain.ErrProyectoNoEncontrado
	}

	nuevo, err := domain.NuevoIntegrante(nombre, rol)
	if err != nil {
		return domain.Integrante{}, err
	}

	existentes, err := s.integrantes.ListarPorProyecto(proyectoID)
	if err != nil {
		if errors.Is(err, domain.ErrProyectoNoEncontrado) || errors.Unwrap(err) == nil {
			return domain.Integrante{}, err
		}
		return domain.Integrante{}, fmt.Errorf("service: listar integrantes: %w", err)
	}
	if domain.EsDuplicado(existentes, nuevo.Nombre) {
		return domain.Integrante{}, domain.ErrIntegranteDuplicado
	}

	if err := s.integrantes.Agregar(proyectoID, nuevo); err != nil {
		if errors.Is(err, domain.ErrProyectoNoEncontrado) || errors.Is(err, domain.ErrIntegranteDuplicado) || errors.Unwrap(err) == nil {
			return domain.Integrante{}, err
		}
		return domain.Integrante{}, fmt.Errorf("service: agregar integrante: %w", err)
	}

	return nuevo, nil
}
