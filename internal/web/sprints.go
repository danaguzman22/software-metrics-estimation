package web

import "net/http"

// sprintsProximamente muestra la página "Disponible próximamente" de Sprints
// con la navegación del proyecto. Devuelve 404 si el proyecto no existe.
func (a *aplicacion) sprintsProximamente(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	registro, err := a.deps.Proyectos.ObtenerProyecto(id)
	if err != nil {
		a.responderErrorProyecto(w, id, err)
		return
	}

	a.renderizar(w, http.StatusOK, "sprints", datosPagina{
		Titulo: "Sprints",
		Nav:    navConProyecto(id),
		Proyecto: &proyectoVista{
			ID:     id,
			Nombre: registro.Proyecto.Nombre,
		},
		Mensaje: "Disponible próximamente",
	})
}
