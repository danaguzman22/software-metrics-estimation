package web_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/web"
)

// pingFalso es un doble de Pinger que devuelve un error configurable y cuenta
// cuántas veces se invocó.
type pingFalso struct {
	err      error
	llamadas int
}

func (p *pingFalso) Ping(context.Context) error {
	p.llamadas++
	return p.err
}

// pingBloqueante simula una base que no responde y espera a que el contexto
// se cancele o venza.
type pingBloqueante struct {
	observoDeadline    bool
	observoVencimiento bool
}

func (p *pingBloqueante) Ping(ctx context.Context) error {
	_, p.observoDeadline = ctx.Deadline()
	<-ctx.Done()
	p.observoVencimiento = errors.Is(ctx.Err(), context.DeadlineExceeded)
	return ctx.Err()
}

func TestT03_HealthzResponde200ConPingOK(t *testing.T) {
	ping := &pingFalso{}
	handler := web.NuevoHandler(ping)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba %d", rec.Code, http.StatusOK)
	}
	if got := rec.Body.String(); got != "ok" {
		t.Fatalf("cuerpo = %q, se esperaba %q", got, "ok")
	}
	if got := rec.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
		t.Fatalf("Content-Type = %q, se esperaba %q", got, "text/plain; charset=utf-8")
	}
	if ping.llamadas != 1 {
		t.Fatalf("Ping se llamó %d veces, se esperaba 1", ping.llamadas)
	}
}

func TestT03_HealthzResponde503ConPingFallido(t *testing.T) {
	ping := &pingFalso{err: errors.New("base caída")}
	handler := web.NuevoHandler(ping)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("código = %d, se esperaba %d", rec.Code, http.StatusServiceUnavailable)
	}
	if rec.Code == http.StatusInternalServerError {
		t.Fatal("no se esperaba un 500 ante una base caída")
	}
	if got := rec.Body.String(); got != "unavailable" {
		t.Fatalf("cuerpo = %q, se esperaba %q", got, "unavailable")
	}
	if ping.llamadas != 1 {
		t.Fatalf("Ping se llamó %d veces, se esperaba 1", ping.llamadas)
	}
}

func TestT03_HealthzRespetaTimeoutDePing(t *testing.T) {
	ping := &pingBloqueante{}
	handler := web.NuevoHandler(ping)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("código = %d, se esperaba %d", rec.Code, http.StatusServiceUnavailable)
	}
	if !ping.observoDeadline {
		t.Fatal("Ping no recibió un contexto con deadline")
	}
	if !ping.observoVencimiento {
		t.Fatal("el contexto de Ping no venció por deadline")
	}
}

func TestT03_HealthzRechazaMetodoNoPermitido(t *testing.T) {
	ping := &pingFalso{}
	handler := web.NuevoHandler(ping)

	req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("código = %d, se esperaba %d", rec.Code, http.StatusMethodNotAllowed)
	}
	if ping.llamadas != 0 {
		t.Fatalf("Ping se llamó %d veces y no debía ejecutarse", ping.llamadas)
	}
}

func TestT03_SoloHealthzExpuesta(t *testing.T) {
	handler := web.NuevoHandler(&pingFalso{})

	for _, ruta := range []string{"/", "/proyectos", "/api/proyectos"} {
		t.Run(ruta, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, ruta, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusNotFound {
				t.Fatalf("código = %d, se esperaba %d para %s", rec.Code, http.StatusNotFound, ruta)
			}
		})
	}
}

func TestT03_ServidorConfiguraTimeouts(t *testing.T) {
	handler := web.NuevoHandler(&pingFalso{})
	servidor := web.NuevoServidor("127.0.0.1:0", handler)

	if servidor == nil {
		t.Fatal("NuevoServidor devolvió nil")
	}
	if servidor.Handler == nil {
		t.Fatal("el servidor no tiene handler configurado")
	}
	if servidor.Addr != "127.0.0.1:0" {
		t.Fatalf("dirección = %q, se esperaba %q", servidor.Addr, "127.0.0.1:0")
	}

	casos := []struct {
		nombre string
		valor  time.Duration
	}{
		{"ReadHeaderTimeout", servidor.ReadHeaderTimeout},
		{"ReadTimeout", servidor.ReadTimeout},
		{"WriteTimeout", servidor.WriteTimeout},
		{"IdleTimeout", servidor.IdleTimeout},
	}
	for _, caso := range casos {
		if caso.valor <= 0 {
			t.Errorf("%s = %v, se esperaba mayor que cero", caso.nombre, caso.valor)
		}
	}
}
