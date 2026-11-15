package domain

import (
	"errors"
	"strings"
)

const (
	SeveridadDefectoAlta  = "Alta"
	SeveridadDefectoMedia = "Media"
	SeveridadDefectoBaja  = "Baja"

	EstadoDefectoAbierto  = "Abierto"
	EstadoDefectoEnCurso  = "En curso"
	EstadoDefectoResuelto = "Resuelto"
)

var (
	ErrDescripcionDefectoVacia            = errors.New("la descripción del defecto no puede estar vacía")
	ErrSeveridadDefectoInvalida           = errors.New("la severidad del defecto no es válida")
	ErrEstadoDefectoInvalido              = errors.New("el estado del defecto no es válido")
	ErrSprintDeteccionDefectoInvalido     = errors.New("el Sprint de detección debe ser mayor que cero")
	ErrSprintResolucionDefectoObligatorio = errors.New("un defecto resuelto requiere Sprint de resolución")
	ErrSprintResolucionAnteriorADeteccion = errors.New("el Sprint de resolución no puede ser anterior al de detección")
	ErrSprintResolucionDefectoInvalido    = errors.New("el Sprint de resolución debe ser mayor que cero")
	ErrDefectoNoEncontrado                = errors.New("no se encontró el defecto")
)

// Defecto registra un defecto y sus referencias a una historia y dos Sprints.
type Defecto struct {
	ID               int64
	Descripcion      string
	Severidad        string
	Estado           string
	ItemID           string
	SprintDeteccion  int
	SprintResolucion *int
}

// NuevoDefecto valida y normaliza los datos iniciales de un defecto.
func NuevoDefecto(
	descripcion, severidad, estado, itemID string,
	sprintDeteccion int,
	sprintResolucion *int,
) (Defecto, error) {
	descripcion = strings.TrimSpace(descripcion)
	if descripcion == "" {
		return Defecto{}, ErrDescripcionDefectoVacia
	}
	switch severidad {
	case SeveridadDefectoAlta, SeveridadDefectoMedia, SeveridadDefectoBaja:
	default:
		return Defecto{}, ErrSeveridadDefectoInvalida
	}
	if err := validarEstadoDefecto(estado); err != nil {
		return Defecto{}, err
	}
	if sprintDeteccion < 1 {
		return Defecto{}, ErrSprintDeteccionDefectoInvalido
	}
	if estado == EstadoDefectoResuelto && sprintResolucion == nil {
		return Defecto{}, ErrSprintResolucionDefectoObligatorio
	}
	if sprintResolucion != nil {
		if *sprintResolucion < 1 {
			return Defecto{}, ErrSprintResolucionDefectoInvalido
		}
		if *sprintResolucion < sprintDeteccion {
			return Defecto{}, ErrSprintResolucionAnteriorADeteccion
		}
	}
	return Defecto{
		Descripcion:      descripcion,
		Severidad:        severidad,
		Estado:           estado,
		ItemID:           itemID,
		SprintDeteccion:  sprintDeteccion,
		SprintResolucion: copiarSprintDefecto(sprintResolucion),
	}, nil
}

// CambiarEstado devuelve una copia del defecto con su nuevo estado y Sprint de
// resolución, si corresponde.
func (d Defecto) CambiarEstado(estado string, sprintResolucion *int) (Defecto, error) {
	actualizado, err := NuevoDefecto(
		d.Descripcion,
		d.Severidad,
		estado,
		d.ItemID,
		d.SprintDeteccion,
		sprintResolucion,
	)
	if err != nil {
		return Defecto{}, err
	}
	actualizado.ID = d.ID
	return actualizado, nil
}

func validarEstadoDefecto(estado string) error {
	switch estado {
	case EstadoDefectoAbierto, EstadoDefectoEnCurso, EstadoDefectoResuelto:
		return nil
	default:
		return ErrEstadoDefectoInvalido
	}
}

func copiarSprintDefecto(sprint *int) *int {
	if sprint == nil {
		return nil
	}
	copia := *sprint
	return &copia
}
