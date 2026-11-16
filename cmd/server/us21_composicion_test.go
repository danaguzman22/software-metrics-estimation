package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestUS21_ConstruirServidorConRutasDePDF verifica que el servidor construido
// registra las rutas de descarga PDF de US-21. Se usa pool nil y solo se
// ejercitan rutas que no consultan la base: un método no permitido sobre la
// ruta registrada responde 405 desde el mux sin ejecutar el handler (R7, D15).
func TestUS21_ConstruirServidorConRutasDePDF(t *testing.T) {
	ping := &pingFalso{}
	servidor := construirServidor(nil, ping, "127.0.0.1:0")
	if servidor == nil || servidor.Handler == nil {
		t.Fatal("construirServidor no devolvió un servidor con handler")
	}

	salud := solicitarRuta(servidor.Handler, "/healthz")
	if salud.Code != http.StatusOK || salud.Body.String() != "ok" {
		t.Fatalf("/healthz: código=%d cuerpo=%q", salud.Code, salud.Body.String())
	}

	for _, ruta := range []string{
		"/proyectos/PRJ-001/reporte.pdf",
		"/proyectos/PRJ-001/sprints/1/reporte.pdf",
	} {
		peticion := httptest.NewRequest(http.MethodPost, ruta, nil)
		grabador := httptest.NewRecorder()
		servidor.Handler.ServeHTTP(grabador, peticion)
		if grabador.Code != http.StatusMethodNotAllowed {
			t.Fatalf("POST %s: código=%d, se esperaba 405 (ruta registrada)", ruta, grabador.Code)
		}
	}
}
