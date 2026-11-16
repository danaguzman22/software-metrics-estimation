package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestUS19_Fase2_ConstruirServidorConRutaDashboard verifica que el armado de
// cmd/server registra la ruta GET /proyectos/{id}/dashboard de US-19 en la
// composición real (el servicio de dashboard queda cableado detrás del
// handler). Se usa un pool nil a propósito y solo se ejercitan rutas que no
// consultan la base: un método no permitido sobre la ruta registrada responde
// 405 desde el mux sin ejecutar el handler (R21). Si la ruta no estuviera
// registrada, el mux respondería 404.
func TestUS19_Fase2_ConstruirServidorConRutaDashboard(t *testing.T) {
	ping := &pingFalso{}

	servidor := construirServidor(nil, ping, "127.0.0.1:0")
	if servidor == nil {
		t.Fatal("construirServidor devolvió nil")
	}
	if servidor.Handler == nil {
		t.Fatal("el servidor no tiene handler configurado")
	}

	peticion := httptest.NewRequest(http.MethodPost, "/proyectos/PRJ-001/dashboard", nil)
	grabador := httptest.NewRecorder()
	servidor.Handler.ServeHTTP(grabador, peticion)

	if grabador.Code != http.StatusMethodNotAllowed {
		t.Fatalf("código de POST /proyectos/PRJ-001/dashboard = %d, se esperaba %d (la ruta debe estar registrada); cuerpo:\n%s",
			grabador.Code, http.StatusMethodNotAllowed, grabador.Body.String())
	}

	// Una ruta inexistente sobre el mismo recurso sigue respondiendo 404 (R21).
	inexistente := solicitarRuta(servidor.Handler, "/proyectos/PRJ-001/dashboard/inexistente")
	if inexistente.Code != http.StatusNotFound {
		t.Fatalf("código de GET /proyectos/PRJ-001/dashboard/inexistente = %d, se esperaba %d",
			inexistente.Code, http.StatusNotFound)
	}
}
