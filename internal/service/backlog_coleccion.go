package service

import (
	"errors"
	"fmt"

	"github.com/danaguzman22/software-metrics-estimation/internal/storage"
)

// cargarColeccionBacklog carga la colección "backlog" desde el almacenamiento
// (D2). Un archivo inexistente se interpreta como colección vacía; cualquier
// otro error de almacenamiento se devuelve envuelto con %w.
func cargarColeccionBacklog(almacen almacen) ([]itemRegistrado, error) {
	var backlog []itemRegistrado
	if err := almacen.Load("backlog", &backlog); err != nil {
		if !errors.Is(err, storage.ErrNotFound) {
			return nil, fmt.Errorf("service: cargar backlog: %w", err)
		}
		// D2: un archivo inexistente se interpreta como backlog vacío.
		return nil, nil
	}
	return backlog, nil
}
