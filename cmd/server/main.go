package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/service"
	"github.com/danaguzman22/software-metrics-estimation/internal/storage/postgres"
	"github.com/danaguzman22/software-metrics-estimation/internal/web"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if err := arrancar(context.Background(), dependenciasPorDefecto()); err != nil {
		log.Fatal(err)
	}
}

func dependenciasPorDefecto() dependenciasArranque {
	return dependenciasArranque{
		conectar: func(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
			bootstrapCtx, cancelar := context.WithTimeout(ctx, 30*time.Second)
			defer cancelar()
			return postgres.ConnectAndMigrate(bootstrapCtx, databaseURL)
		},
		escuchar: net.Listen,
	}
}

func requiredDatabaseURL() (string, error) {
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		return "", errors.New("la variable DATABASE_URL es obligatoria")
	}
	return databaseURL, nil
}

// direccionHTTP devuelve la dirección de escucha configurada por entorno.
func direccionHTTP() string {
	direccion := strings.TrimSpace(os.Getenv("HTTP_ADDR"))
	if direccion == "" {
		return ":8080"
	}
	return direccion
}

// dependenciasArranque agrupa las funciones y dependencias inyectables del
// arranque para poder probarlo sin PostgreSQL real.
type dependenciasArranque struct {
	conectar func(context.Context, string) (*pgxpool.Pool, error)
	escuchar func(string, string) (net.Listener, error)
	ping     web.Pinger
}

type errorConexionPostgres struct {
	err error
}

func (e errorConexionPostgres) Error() string {
	return "conectar con PostgreSQL"
}

func (e errorConexionPostgres) Unwrap() error {
	return e.err
}

// servidorApagable es la parte del servidor que usa el apagado ordenado.
type servidorApagable interface {
	Shutdown(context.Context) error
}

// poolCerrable es la parte del pool que usa el apagado ordenado.
type poolCerrable interface {
	Close()
}

// construirServidor arma la aplicación web completa: repositorios PostgreSQL,
// servicios y handlers. Conserva la firma de T-03 para que construirServidor
// con un pool nil siga funcionando en los tests de salud.
func construirServidor(pool *pgxpool.Pool, ping web.Pinger, addr string) *http.Server {
	proyectos := postgres.NuevoProyectos(pool)
	integrantes := postgres.NuevoIntegrantes(pool)
	backlog := postgres.NuevoBacklog(pool)
	sprints := postgres.NuevoSprints(pool)
	asignaciones := postgres.NuevoAsignaciones(pool)
	rondasPoker := postgres.NuevoRondasPoker(pool)

	servicioProyectos := service.NuevoServicioProyectos(proyectos, integrantes)
	servicioBacklog := service.NuevoServicioBacklog(proyectos, backlog)
	servicioSprints := service.NuevoServicioSprints(proyectos, sprints, backlog, asignaciones)
	servicioPoker := service.NuevoServicioPlanningPoker(proyectos, integrantes, backlog, rondasPoker)

	handler := web.NuevoHandlerAplicacion(ping, web.Dependencias{
		Proyectos: servicioProyectos,
		Backlog:   servicioBacklog,
		Sprints:   servicioSprints,
		Poker:     servicioPoker,
		Ahora:     time.Now,
	})
	return web.NuevoServidor(addr, handler)
}

// arrancar conecta con la base, levanta el servidor y espera una señal o un
// error de escucha. El pool se cierra siempre al terminar.
func arrancar(ctx context.Context, deps dependenciasArranque) error {
	databaseURL, err := requiredDatabaseURL()
	if err != nil {
		return err
	}

	pool, err := deps.conectar(ctx, databaseURL)
	if err != nil {
		return errorConexionPostgres{err: err}
	}

	ping := deps.ping
	if ping == nil {
		ping = pool
	}

	direccion := direccionHTTP()
	servidor := construirServidor(pool, ping, direccion)

	ctxSenales, detener := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer detener()

	listener, err := deps.escuchar("tcp", direccion)
	if err != nil {
		pool.Close()
		return fmt.Errorf("escuchar en %q: %w", direccion, err)
	}

	errServidor := make(chan error, 1)
	go func() {
		errServidor <- servidor.Serve(listener)
	}()

	select {
	case err := <-errServidor:
		pool.Close()
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("servir HTTP: %w", err)
	case <-ctxSenales.Done():
		return apagar(servidor, pool)
	}
}

// apagar ejecuta el apagado ordenado con un contexto nuevo y cierra el pool
// aunque Shutdown falle.
func apagar(servidor servidorApagable, pool poolCerrable) error {
	ctxApagado, cancelar := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelar()

	errShutdown := servidor.Shutdown(ctxApagado)
	pool.Close()

	if errShutdown != nil {
		return fmt.Errorf("apagar servidor: %w", errShutdown)
	}
	return nil
}
