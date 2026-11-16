package web

import (
	"net/http"
)

// dashboardProyecto muestra el dashboard del proyecto. En la fase RED responde
// 501; la implementación llega en la fase GREEN.
func (a *aplicacion) dashboardProyecto(w http.ResponseWriter, r *http.Request) {
	a.renderizar(w, http.StatusNotImplemented, "dashboard", datosPagina{
		Titulo: "Dashboard",
		Nav:    navSinProyecto(false),
	})
}

// escalaConMinimoUno devuelve el tope del eje Y de un gráfico: el máximo de las
// barras o 1 cuando todas valen cero (D15).
func escalaConMinimoUno(valores ...float64) float64 {
	return 0
}
