package domain

import "errors"

var (
	ErrItemNoEncontrado     = errors.New("el ítem no existe en el proyecto")
	ErrStoryPointsInvalidos = errors.New("los story points no son válidos: deben ser 1, 2, 3, 5, 8, 13 o 21")
)

func EsStoryPointsValido(storyPoints int) bool {
	return false
}

func (i ItemBacklog) Estimar(storyPoints int) (ItemBacklog, error) {
	return ItemBacklog{}, errors.New("no implementado")
}
