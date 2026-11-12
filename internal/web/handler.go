package web

import (
	"context"
	"io"
	"net/http"
	"time"
)

// tiempoEsperaPing es el plazo máximo que se espera a la base en /healthz.
const tiempoEsperaPing = 2 * time.Second

// Pinger es la dependencia mínima para verificar la salud de la base.
type Pinger interface {
	Ping(context.Context) error
}

// NuevoHandler arma las rutas de la capa web. Por ahora solo expone
// GET /healthz, que T-04 usará como healthcheck.
func NuevoHandler(ping Pinger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancelar := context.WithTimeout(r.Context(), tiempoEsperaPing)
		defer cancelar()

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		if err := ping.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = io.WriteString(w, "unavailable")
			return
		}

		_, _ = io.WriteString(w, "ok")
	})

	return mux
}

// NuevoServidor construye el servidor HTTP con timeouts explícitos.
func NuevoServidor(addr string, h http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
