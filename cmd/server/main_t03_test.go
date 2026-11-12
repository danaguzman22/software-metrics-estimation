package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/web"
	"github.com/jackc/pgx/v5/pgxpool"
)

// pingFalso es un doble de Pinger para el armado del servidor.
type pingFalso struct {
	err      error
	llamadas int
}

func (p *pingFalso) Ping(context.Context) error {
	p.llamadas++
	return p.err
}

// servidorApagableFalso registra las llamadas a Shutdown y el deadline que
// recibió, sin abrir un servidor real.
type servidorApagableFalso struct {
	llamadas      int
	tieneDeadline bool
	restante      time.Duration
	err           error
	alApagar      func()
}

func (s *servidorApagableFalso) Shutdown(ctx context.Context) error {
	s.llamadas++
	if deadline, ok := ctx.Deadline(); ok {
		s.tieneDeadline = true
		s.restante = time.Until(deadline)
	}
	if s.alApagar != nil {
		s.alApagar()
	}
	return s.err
}

// poolFalso registra los cierres para verificar el orden del apagado.
type poolFalso struct {
	cierres int
	orden   *[]string
}

func (p *poolFalso) Close() {
	p.cierres++
	if p.orden != nil {
		*p.orden = append(*p.orden, "close")
	}
}

func TestT03_DatabaseURLObligatoria(t *testing.T) {
	casos := []struct {
		nombre string
		valor  string
	}{
		{"vacía", ""},
		{"solo espacios", "   "},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			t.Setenv("DATABASE_URL", caso.valor)

			if _, err := requiredDatabaseURL(); err == nil {
				t.Fatal("se esperaba un error de configuración")
			}
		})
	}
}

func TestT03_HTTPAddrPorDefecto(t *testing.T) {
	casos := []struct {
		nombre string
		valor  string
	}{
		{"vacía", ""},
		{"solo espacios", "   "},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			t.Setenv("HTTP_ADDR", caso.valor)

			if got := direccionHTTP(); got != ":8080" {
				t.Fatalf("dirección = %q, se esperaba %q", got, ":8080")
			}
		})
	}
}

func TestT03_HTTPAddrPersonalizado(t *testing.T) {
	t.Setenv("HTTP_ADDR", "127.0.0.1:9090")

	if got := direccionHTTP(); got != "127.0.0.1:9090" {
		t.Fatalf("dirección = %q, se esperaba %q", got, "127.0.0.1:9090")
	}
}

func TestT03_ConstruirServidorConDependenciasFalsas(t *testing.T) {
	ping := &pingFalso{}

	servidor := construirServidor(nil, ping, "127.0.0.1:0")

	if servidor == nil {
		t.Fatal("construirServidor devolvió nil")
	}
	if servidor.Addr != "127.0.0.1:0" {
		t.Fatalf("dirección = %q, se esperaba %q", servidor.Addr, "127.0.0.1:0")
	}
	if servidor.Handler == nil {
		t.Fatal("el servidor no tiene handler configurado")
	}

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	servidor.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba %d", rec.Code, http.StatusOK)
	}
	if ping.llamadas != 1 {
		t.Fatalf("Ping se llamó %d veces, se esperaba 1", ping.llamadas)
	}
}

func TestT03_ArranqueConstruyeServidorYNoTermina(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://usuario:clave@127.0.0.1:1/base")
	t.Setenv("HTTP_ADDR", "127.0.0.1:0")

	pool := nuevoPoolInalcanzable(t)
	ping := &pingFalso{}
	listener := nuevoListener(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	deps := dependenciasArranque{
		conectar: func(context.Context, string) (*pgxpool.Pool, error) { return pool, nil },
		escuchar: func(string, string) (net.Listener, error) { return listener, nil },
		ping:     ping,
	}

	terminado := make(chan error, 1)
	go func() { terminado <- arrancar(ctx, deps) }()

	esperarHealthz(t, listener.Addr().String())

	select {
	case err := <-terminado:
		t.Fatalf("el arranque terminó antes de cancelar: %v", err)
	case <-time.After(100 * time.Millisecond):
	}

	cancel()

	select {
	case err := <-terminado:
		if err != nil {
			t.Fatalf("el arranque devolvió error al cancelar: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("el arranque no terminó tras cancelar el contexto")
	}

	if err := pool.Ping(context.Background()); err == nil || !strings.Contains(err.Error(), "closed pool") {
		t.Fatalf("el pool debía estar cerrado; Ping devolvió: %v", err)
	}
}

func TestT03_ApagadoOrdenado(t *testing.T) {
	var orden []string
	servidor := &servidorApagableFalso{
		alApagar: func() { orden = append(orden, "shutdown") },
	}
	pool := &poolFalso{orden: &orden}

	if err := apagar(servidor, pool); err != nil {
		t.Fatalf("apagar devolvió error: %v", err)
	}
	if servidor.llamadas != 1 {
		t.Fatalf("Shutdown se llamó %d veces, se esperaba 1", servidor.llamadas)
	}
	if pool.cierres != 1 {
		t.Fatalf("Close se llamó %d veces, se esperaba 1", pool.cierres)
	}
	if !servidor.tieneDeadline {
		t.Fatal("Shutdown no recibió un contexto con deadline")
	}
	if servidor.restante < 9*time.Second || servidor.restante > 11*time.Second {
		t.Fatalf("deadline de Shutdown = %v, se esperaba cerca de 10s", servidor.restante)
	}

	esperado := []string{"shutdown", "close"}
	if len(orden) != len(esperado) || orden[0] != esperado[0] || orden[1] != esperado[1] {
		t.Fatalf("orden de apagado = %v, se esperaba %v", orden, esperado)
	}
}

func TestT03_ApagadoCierraPoolAunqueShutdownFalle(t *testing.T) {
	fallo := errors.New("agotado el tiempo de apagado")
	servidor := &servidorApagableFalso{err: fallo}
	pool := &poolFalso{}

	err := apagar(servidor, pool)

	if err == nil {
		t.Fatal("se esperaba el error de Shutdown")
	}
	if !errors.Is(err, fallo) {
		t.Fatalf("error = %v, se esperaba envolver %v", err, fallo)
	}
	if pool.cierres != 1 {
		t.Fatalf("Close se llamó %d veces, se esperaba 1", pool.cierres)
	}
}

func TestT03_RunTerminaAlCancelarContexto(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://usuario:clave@127.0.0.1:1/base")
	t.Setenv("HTTP_ADDR", "127.0.0.1:0")

	pool := nuevoPoolInalcanzable(t)
	listener := nuevoListener(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	deps := dependenciasArranque{
		conectar: func(context.Context, string) (*pgxpool.Pool, error) { return pool, nil },
		escuchar: func(string, string) (net.Listener, error) { return listener, nil },
		ping:     &pingFalso{},
	}

	terminado := make(chan error, 1)
	go func() { terminado <- arrancar(ctx, deps) }()

	esperarHealthz(t, listener.Addr().String())
	cancel()

	select {
	case err := <-terminado:
		if err != nil {
			t.Fatalf("arrancar devolvió error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("arrancar no terminó al cancelar el contexto")
	}
}

func TestT03_ArranqueNoIniciaServidorSiFallaConexion(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://usuario:clave@127.0.0.1:1/base")
	t.Setenv("HTTP_ADDR", "127.0.0.1:0")

	falloConexion := errors.New("conexión rechazada")
	escucho := false

	deps := dependenciasArranque{
		conectar: func(context.Context, string) (*pgxpool.Pool, error) {
			return nil, falloConexion
		},
		escuchar: func(string, string) (net.Listener, error) {
			escucho = true
			return nil, errors.New("no debía escuchar")
		},
		ping: &pingFalso{},
	}

	err := arrancar(context.Background(), deps)

	if err == nil {
		t.Fatal("se esperaba un error de conexión")
	}
	if !errors.Is(err, falloConexion) {
		t.Fatalf("error = %v, se esperaba envolver %v", err, falloConexion)
	}
	if escucho {
		t.Fatal("no debía abrirse el listener si falló la conexión")
	}
}

func TestT03_ArranqueDevuelveErrorSiNoPuedeEscuchar(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://usuario:clave@127.0.0.1:1/base")
	t.Setenv("HTTP_ADDR", "127.0.0.1:0")

	pool := nuevoPoolInalcanzable(t)
	falloEscucha := errors.New("puerto ocupado")

	deps := dependenciasArranque{
		conectar: func(context.Context, string) (*pgxpool.Pool, error) { return pool, nil },
		escuchar: func(string, string) (net.Listener, error) { return nil, falloEscucha },
		ping:     &pingFalso{},
	}

	err := arrancar(context.Background(), deps)

	if err == nil {
		t.Fatal("se esperaba un error de escucha")
	}
	if !errors.Is(err, falloEscucha) {
		t.Fatalf("error = %v, se esperaba envolver %v", err, falloEscucha)
	}
	if err := pool.Ping(context.Background()); err == nil || !strings.Contains(err.Error(), "closed pool") {
		t.Fatalf("el pool debía estar cerrado; Ping devolvió: %v", err)
	}
}

// nuevoPoolInalcanzable crea un pool que no abre conexiones reales: pgxpool
// valida la cadena de conexión de forma perezosa.
func nuevoPoolInalcanzable(t *testing.T) *pgxpool.Pool {
	t.Helper()

	pool, err := pgxpool.New(context.Background(), "postgres://usuario:clave@127.0.0.1:1/base")
	if err != nil {
		t.Fatalf("crear pool de prueba: %v", err)
	}
	return pool
}

func nuevoListener(t *testing.T) net.Listener {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("abrir listener de prueba: %v", err)
	}
	return listener
}

func esperarHealthz(t *testing.T, direccion string) {
	t.Helper()

	cliente := &http.Client{
		Timeout:   2 * time.Second,
		Transport: &http.Transport{DisableKeepAlives: true},
	}
	url := "http://" + direccion + "/healthz"
	limite := time.Now().Add(5 * time.Second)

	for time.Now().Before(limite) {
		resp, err := cliente.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("el servidor no respondió 200 en %s", url)
}
