package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestUS20_Fase2_ConstruirServidorConRutasDeReporte verifica que el armado de
// cmd/server conserva /healthz de T-03 y registra las rutas de reporte de
// US-20 en la composición real (el servicio de reportes queda cableado
// detrás del handler). Se usa un pool nil a propósito y solo se ejercitan
// rutas que no consultan la base: un método no permitido sobre la ruta
// registrada responde 405 desde el mux sin ejecutar el handler. Si la ruta
// no estuviera registrada, el mux respondería 404.
func TestUS20_Fase2_ConstruirServidorConRutasDeReporte(t *testing.T) {
	ping := &pingFalso{}

	servidor := construirServidor(nil, ping, "127.0.0.1:0")
	if servidor == nil {
		t.Fatal("construirServidor devolvió nil")
	}
	if servidor.Handler == nil {
		t.Fatal("el servidor no tiene handler configurado")
	}

	// Comportamiento de T-03: /healthz responde 200 "ok" y consulta el ping.
	salud := solicitarRuta(servidor.Handler, "/healthz")
	if salud.Code != http.StatusOK {
		t.Fatalf("código de /healthz = %d, se esperaba %d", salud.Code, http.StatusOK)
	}
	if cuerpo := salud.Body.String(); cuerpo != "ok" {
		t.Fatalf("cuerpo de /healthz = %q, se esperaba %q", cuerpo, "ok")
	}
	if ping.llamadas != 1 {
		t.Fatalf("Ping se llamó %d veces en /healthz, se esperaba 1", ping.llamadas)
	}

	// Ruta de negocio de US-20: un método no permitido sobre la ruta registrada
	// responde 405 desde el mux sin ejecutar el handler (R20). Si la ruta no
	// estuviera registrada, el mux respondería 404.
	for _, ruta := range []string{
		"/proyectos/PRJ-001/reporte",
		"/proyectos/PRJ-001/sprints/1/reporte",
	} {
		t.Run(ruta, func(t *testing.T) {
			peticion := httptest.NewRequest(http.MethodPost, ruta, nil)
			grabador := httptest.NewRecorder()
			servidor.Handler.ServeHTTP(grabador, peticion)

			if grabador.Code != http.StatusMethodNotAllowed {
				t.Fatalf("código de POST %s = %d, se esperaba %d (la ruta debe estar registrada); cuerpo:\n%s",
					ruta, grabador.Code, http.StatusMethodNotAllowed, grabador.Body.String())
			}
		})
	}

	// Una ruta inexistente sigue respondiendo 404 (R20).
	inexistente := solicitarRuta(servidor.Handler, "/proyectos/PRJ-001/reporte/inexistente")
	if inexistente.Code != http.StatusNotFound {
		t.Fatalf("código de GET /proyectos/PRJ-001/reporte/inexistente = %d, se esperaba %d",
			inexistente.Code, http.StatusNotFound)
	}
}
