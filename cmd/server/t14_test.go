package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestT14_ConstruirServidorConRutasDeIntegrantes verifica que la composición
// real registra GET y POST de /proyectos/{id}/integrantes. Se usa un pool nil y
// solo se ejercitan rutas que no consultan la base: el método no permitido lo
// resuelve el ServeMux sin ejecutar el handler (R1, R12, D12).
func TestT14_ConstruirServidorConRutasDeIntegrantes(t *testing.T) {
	ping := &pingFalso{}
	servidor := construirServidor(nil, ping, "127.0.0.1:0")
	if servidor == nil || servidor.Handler == nil {
		t.Fatal("construirServidor no devolvió un servidor con handler")
	}

	salud := solicitarRuta(servidor.Handler, "/healthz")
	if salud.Code != http.StatusOK || salud.Body.String() != "ok" {
		t.Fatalf("/healthz: código=%d cuerpo=%q", salud.Code, salud.Body.String())
	}

	for _, metodo := range []string{http.MethodPut, http.MethodPatch, http.MethodDelete} {
		peticion := httptest.NewRequest(metodo, "/proyectos/PRJ-001/integrantes", nil)
		grabador := httptest.NewRecorder()
		servidor.Handler.ServeHTTP(grabador, peticion)
		if grabador.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s /proyectos/PRJ-001/integrantes: código=%d, se esperaba 405 (ruta registrada)",
				metodo, grabador.Code)
		}
	}

	inexistente := solicitarRuta(servidor.Handler, "/proyectos/PRJ-001/integrantes/inexistente")
	if inexistente.Code != http.StatusNotFound {
		t.Errorf("ruta inexistente de integrantes: código=%d, se esperaba 404",
			inexistente.Code)
	}
}
