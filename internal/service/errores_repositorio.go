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
		errors.Is(err, domain.ErrSprintDuplicado) ||
		errors.Is(err, domain.ErrSprintNoEncontrado) ||
		errors.Is(err, domain.ErrSprintCerrado) ||
		errors.Is(err, domain.ErrItemYaAsignado) ||
		errors.Is(err, domain.ErrItemNoAsignadoAlSprint) ||
		errors.Is(err, domain.ErrIntegranteNoPerteneceAlProyecto) ||
		errors.Is(err, domain.ErrIntegranteNoDelProyecto) ||
		errors.Is(err, domain.ErrValorVotoPokerInvalido) ||
		errors.Is(err, domain.ErrNumeroRondaInvalido) ||
		errors.Is(err, domain.ErrRondaNoEncontrada) ||
		errors.Is(err, domain.ErrRondaAbiertaExistente) ||
		errors.Is(err, domain.ErrRondaRevelada) ||
		errors.Is(err, domain.ErrRondaNoRevelada) ||
		errors.Is(err, domain.ErrRondaSinVotos) ||
		errors.Is(err, domain.ErrHorasEsfuerzoInvalidas) ||
		errors.Is(err, domain.ErrFechaEsfuerzoFutura) ||
		errors.Is(err, domain.ErrDescripcionDefectoVacia) ||
		errors.Is(err, domain.ErrSeveridadDefectoInvalida) ||
		errors.Is(err, domain.ErrEstadoDefectoInvalido) ||
		errors.Is(err, domain.ErrSprintDeteccionDefectoInvalido) ||
		errors.Is(err, domain.ErrSprintResolucionDefectoObligatorio) ||
		errors.Is(err, domain.ErrSprintResolucionAnteriorADeteccion) ||
		errors.Is(err, domain.ErrSprintResolucionDefectoInvalido) ||
		errors.Is(err, domain.ErrDefectoNoEncontrado) {
		return err
	}
	return fmt.Errorf("service: %s: %w", contexto, err)
}
