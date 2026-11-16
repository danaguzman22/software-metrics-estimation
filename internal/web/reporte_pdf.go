package web

import "net/http"

// reporteProyectoPDF es el stub de descarga del reporte del proyecto (US-21, RED).
func (a *aplicacion) reporteProyectoPDF(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "no implementado", http.StatusNotImplemented)
}

// reporteSprintPDF es el stub de descarga del reporte de un Sprint (US-21, RED).
func (a *aplicacion) reporteSprintPDF(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "no implementado", http.StatusNotImplemented)
}
