package web

import "net/http"

// paginaIntegrantes atiende GET /proyectos/{id}/integrantes.
// RED (T-14): stub mínimo para compilar; se implementa en la fase GREEN.
func (a *aplicacion) paginaIntegrantes(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// crearIntegrante atiende POST /proyectos/{id}/integrantes.
// RED (T-14): stub mínimo para compilar; se implementa en la fase GREEN.
func (a *aplicacion) crearIntegrante(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
