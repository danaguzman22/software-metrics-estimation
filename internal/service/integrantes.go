package service

import (
	"errors"
	"fmt"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage"
)

// integranteRegistrado es el registro persistido de un integrante: el ID del
// proyecto junto al integrante del dominio (D3). El ID no forma parte de
// domain.Integrante.
type integranteRegistrado struct {
	ProyectoID string            `json:"proyecto_id"`
	Integrante domain.Integrante `json:"integrante"`
}

// AgregarIntegrante registra un integrante con su nombre y su rol en el
// proyecto identificado por proyectoID (D5).
//
// Carga la colección "proyectos" (D13); un archivo inexistente se interpreta
// como colección vacía (D1). Busca el proyecto por ID: si no existe devuelve
// domain.ErrProyectoNoEncontrado, antes de validar nombre y rol (C1, C7).
// Delega las validaciones en domain.NuevoIntegrante; si falla, devuelve ese
// error sin escribir nada. Carga la colección "integrantes" (D3); un archivo
// inexistente se interpreta como colección vacía (primera alta) y cualquier
// otro error se propaga envuelto (D1). Comprueba el duplicado con
// domain.EsDuplicado contra los integrantes de ese proyecto y, si existe,
// devuelve domain.ErrIntegranteDuplicado (C2, D4, D6). En caso de éxito agrega
// el registro {ProyectoID, Integrante} y guarda la colección completa (D3).
//
// Ante cualquier error devuelve domain.Integrante{} sin invocar Save. Si la
// persistencia falla durante Save, también devuelve domain.Integrante{} junto
// con el error envuelto. Los errores de almacenamiento distintos de "no
// encontrado" se propagan envueltos con %w y no se mapean a
// domain.ErrProyectoNoEncontrado (D1).
func (s *ServicioProyectos) AgregarIntegrante(
	proyectoID, nombre, rol string,
) (domain.Integrante, error) {
	if err := verificarProyectoExistente(s.almacen, proyectoID); err != nil {
		return domain.Integrante{}, err
	}

	nuevo, err := domain.NuevoIntegrante(nombre, rol)
	if err != nil {
		return domain.Integrante{}, err
	}

	var integrantes []integranteRegistrado
	if err := s.almacen.Load("integrantes", &integrantes); err != nil {
		if !errors.Is(err, storage.ErrNotFound) {
			return domain.Integrante{}, fmt.Errorf("service: cargar integrantes: %w", err)
		}
		// D1: un archivo inexistente se interpreta como colección vacía.
		integrantes = nil
	}

	existentes := make([]domain.Integrante, 0, len(integrantes))
	for _, registro := range integrantes {
		if registro.ProyectoID == proyectoID {
			existentes = append(existentes, registro.Integrante)
		}
	}
	if domain.EsDuplicado(existentes, nuevo.Nombre) {
		return domain.Integrante{}, domain.ErrIntegranteDuplicado
	}

	integrantes = append(integrantes, integranteRegistrado{
		ProyectoID: proyectoID,
		Integrante: nuevo,
	})
	if err := s.almacen.Save("integrantes", integrantes); err != nil {
		return domain.Integrante{}, fmt.Errorf("service: guardar integrantes: %w", err)
	}

	return nuevo, nil
}
