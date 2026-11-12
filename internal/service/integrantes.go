package service

import (
	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

// AgregarIntegrante registra un integrante con su nombre y su rol en el
// proyecto identificado por proyectoID.
func (s *ServicioProyectos) AgregarIntegrante(
	proyectoID, nombre, rol string,
) (domain.Integrante, error) {
	existe, err := s.proyectos.Existe(proyectoID)
	if err != nil {
		return domain.Integrante{}, envolverErrorRepositorio("verificar proyecto", err)
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
		return domain.Integrante{}, envolverErrorRepositorio("listar integrantes", err)
	}
	if domain.EsDuplicado(existentes, nuevo.Nombre) {
		return domain.Integrante{}, domain.ErrIntegranteDuplicado
	}

	if err := s.integrantes.Agregar(proyectoID, nuevo); err != nil {
		return domain.Integrante{}, envolverErrorRepositorio("agregar integrante", err)
	}

	return nuevo, nil
}
