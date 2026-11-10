package domain

import "errors"

var (
	ErrItemNoEncontrado     = errors.New("el ítem no existe en el proyecto")
	ErrStoryPointsInvalidos = errors.New("los story points no son válidos: deben ser 1, 2, 3, 5, 8, 13 o 21")
)

func EsStoryPointsValido(storyPoints int) bool {
	switch storyPoints {
	case 1, 2, 3, 5, 8, 13, 21:
		return true
	default:
		return false
	}
}

func (i ItemBacklog) Estimar(storyPoints int) (ItemBacklog, error) {
	if !EsStoryPointsValido(storyPoints) {
		return ItemBacklog{}, ErrStoryPointsInvalidos
	}
	i.StoryPoints = storyPoints
	return i, nil
}
