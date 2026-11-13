package domain

import "errors"

// ErrItemCompletadoNoEditable indica que el ítem está completado y no admite
// edición.
var ErrItemCompletadoNoEditable = errors.New("el ítem está completado y no se puede editar")

// EstadoCompletada es el estado de un ítem cerrado; no admite edición.
const EstadoCompletada = "Completada"

// Editar devuelve una copia del ítem con los campos editables actualizados y
// el ID y el estado originales. Reutiliza NuevoItemBacklog para validar y
// recortar título y prioridad. Para story points admite 0 ("sin estimar") o
// valida y asigna un valor de la escala con Estimar. Devuelve ItemBacklog{} y
// el error correspondiente si el ítem está en EstadoCompletada o si algún
// campo editable no es válido.
func (i ItemBacklog) Editar(
	titulo, descripcion, prioridad string,
	criterios []string,
	storyPoints int,
) (ItemBacklog, error) {
	if i.Estado == EstadoCompletada {
		return ItemBacklog{}, ErrItemCompletadoNoEditable
	}

	editado, err := NuevoItemBacklog(titulo, descripcion, prioridad, criterios)
	if err != nil {
		return ItemBacklog{}, err
	}

	if storyPoints != 0 {
		editado, err = editado.Estimar(storyPoints)
		if err != nil {
			return ItemBacklog{}, err
		}
	}

	editado.ID = i.ID
	editado.Estado = i.Estado
	return editado, nil
}
