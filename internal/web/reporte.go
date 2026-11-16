package web

import "net/http"

// reporteProyecto muestra el reporte imprimible del proyecto (stub de la
// fase RED: responde 501 hasta el GREEN).
func (a *aplicacion) reporteProyecto(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "no implementado", http.StatusNotImplemented)
}

// reporteSprint muestra el reporte imprimible de un Sprint (stub de la fase
// RED: responde 501 hasta el GREEN).
func (a *aplicacion) reporteSprint(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "no implementado", http.StatusNotImplemented)
}
