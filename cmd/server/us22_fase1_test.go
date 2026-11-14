package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// solicitarRuta ejecuta una petición GET contra el handler y devuelve la
// respuesta grabada, sin abrir un servidor real.
func solicitarRuta(handler http.Handler, ruta string) *httptest.ResponseRecorder {
	peticion := httptest.NewRequest(http.MethodGet, ruta, nil)
	grabador := httptest.NewRecorder()
	handler.ServeHTTP(grabador, peticion)
	return grabador
}

// TestUS22_Fase1_ConstruirServidorConRutasDeNegocio verifica que el armado de
// cmd/server conserva /healthz de T-03 y expone la página de inicio de US-22.
// Se usa un pool nil a propósito: solo se ejercitan rutas que no consultan la base.
func TestUS22_Fase1_ConstruirServidorConRutasDeNegocio(t *testing.T) {
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

	// Ruta de negocio de US-22: la página de inicio muestra el título.
	inicio := solicitarRuta(servidor.Handler, "/")
	if inicio.Code != http.StatusOK {
		t.Fatalf("código de / = %d, se esperaba %d; cuerpo:\n%s", inicio.Code, http.StatusOK, inicio.Body.String())
	}
	if cuerpo := inicio.Body.String(); !strings.Contains(cuerpo, "Software Metrics & Estimation") {
		t.Fatalf("el cuerpo de / no contiene el título %q; cuerpo:\n%s", "Software Metrics & Estimation", cuerpo)
	}
	if ping.llamadas != 1 {
		t.Fatalf("Ping se llamó %d veces en total, se esperaba 1: la página de inicio no debe consultar la base", ping.llamadas)
	}
}
