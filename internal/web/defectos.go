package web

import "net/http"

// defectos muestra la página "Disponible próximamente" de Defectos sin
// ejecutar lógica de negocio.
func (a *aplicacion) defectos(w http.ResponseWriter, r *http.Request) {
	a.renderizar(w, http.StatusOK, "defectos", datosPagina{
		Titulo:  "Defectos",
		Nav:     navSinProyecto(false),
		Mensaje: "Disponible próximamente",
	})
}
