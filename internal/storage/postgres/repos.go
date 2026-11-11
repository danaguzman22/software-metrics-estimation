package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	_ service.RepositorioProyectos   = (*Proyectos)(nil)
	_ service.RepositorioIntegrantes = (*Integrantes)(nil)
	_ service.RepositorioBacklog     = (*Backlog)(nil)
)

type Proyectos struct {
	pool *pgxpool.Pool
}

func NuevoProyectos(pool *pgxpool.Pool) *Proyectos {
	return &Proyectos{pool: pool}
}

func (r *Proyectos) Obtener(id string) (service.ProyectoRegistrado, error) {
	proyecto, err := scanProyecto(r.pool.QueryRow(context.Background(), `
		SELECT id, nombre, fecha_inicio, fecha_fin, estado
		FROM proyectos
		WHERE id = $1
	`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return service.ProyectoRegistrado{}, domain.ErrProyectoNoEncontrado
	}
	if err != nil {
		return service.ProyectoRegistrado{}, fmt.Errorf("obtener proyecto: %w", err)
	}
	return proyecto, nil
}

func (r *Proyectos) Existe(id string) (bool, error) {
	var existe bool
	if err := r.pool.QueryRow(context.Background(),
		"SELECT EXISTS (SELECT 1 FROM proyectos WHERE id = $1)", id,
	).Scan(&existe); err != nil {
		return false, fmt.Errorf("verificar existencia de proyecto: %w", err)
	}
	return existe, nil
}

func (r *Proyectos) Listar() ([]service.ProyectoRegistrado, error) {
	rows, err := r.pool.Query(context.Background(), `
		SELECT id, nombre, fecha_inicio, fecha_fin, estado
		FROM proyectos
		ORDER BY id
	`)
	if err != nil {
		return nil, fmt.Errorf("listar proyectos: %w", err)
	}
	defer rows.Close()

	proyectos := make([]service.ProyectoRegistrado, 0)
	for rows.Next() {
		proyecto, err := scanProyecto(rows)
		if err != nil {
			return nil, fmt.Errorf("leer proyecto: %w", err)
		}
		proyectos = append(proyectos, proyecto)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recorrer proyectos: %w", err)
	}
	return proyectos, nil
}

func (r *Proyectos) Actualizar(id string, proyecto domain.Proyecto) error {
	tag, err := r.pool.Exec(context.Background(), `
		UPDATE proyectos
		SET nombre = $2, fecha_inicio = $3::date, fecha_fin = $4::date, estado = $5
		WHERE id = $1
	`, id, proyecto.Nombre, fechaPostgres(proyecto.FechaInicio),
		fechaPostgres(proyecto.FechaFin), proyecto.Estado)
	if err != nil {
		return fmt.Errorf("actualizar proyecto: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrProyectoNoEncontrado
	}
	return nil
}

type Integrantes struct {
	pool *pgxpool.Pool
}

func NuevoIntegrantes(pool *pgxpool.Pool) *Integrantes {
	return &Integrantes{pool: pool}
}

func (r *Integrantes) ListarPorProyecto(proyectoID string) ([]domain.Integrante, error) {
	rows, err := r.pool.Query(context.Background(), `
		SELECT nombre, rol
		FROM integrantes
		WHERE proyecto_id = $1
		ORDER BY id
	`, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar integrantes: %w", err)
	}
	defer rows.Close()

	integrantes := make([]domain.Integrante, 0)
	for rows.Next() {
		var integrante domain.Integrante
		if err := rows.Scan(&integrante.Nombre, &integrante.Rol); err != nil {
			return nil, fmt.Errorf("leer integrante: %w", err)
		}
		integrantes = append(integrantes, integrante)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recorrer integrantes: %w", err)
	}
	return integrantes, nil
}

func (r *Integrantes) Agregar(proyectoID string, integrante domain.Integrante) error {
	_, err := r.pool.Exec(context.Background(), `
		INSERT INTO integrantes (proyecto_id, nombre, rol)
		VALUES ($1, $2, $3)
	`, proyectoID, integrante.Nombre, integrante.Rol)
	if err != nil {
		var postgresError *pgconn.PgError
		if errors.As(err, &postgresError) && postgresError.Code == "23505" {
			return domain.ErrIntegranteDuplicado
		}
		return fmt.Errorf("agregar integrante: %w", err)
	}
	return nil
}

type Backlog struct {
	pool *pgxpool.Pool
}

func NuevoBacklog(pool *pgxpool.Pool) *Backlog {
	return &Backlog{pool: pool}
}

func (r *Backlog) AgregarItem(proyectoID string, item domain.ItemBacklog) (domain.ItemBacklog, error) {
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.ItemBacklog{}, fmt.Errorf("iniciar transacción para agregar ítem: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var lockedProjectID string
	if err := tx.QueryRow(ctx, `
		SELECT id
		FROM proyectos
		WHERE id = $1
		FOR UPDATE
	`, proyectoID).Scan(&lockedProjectID); err != nil {
		return domain.ItemBacklog{}, fmt.Errorf("bloquear proyecto para agregar ítem: %w", err)
	}

	rows, err := tx.Query(ctx, `
		SELECT numero
		FROM items_backlog
		WHERE proyecto_id = $1
	`, proyectoID)
	if err != nil {
		return domain.ItemBacklog{}, fmt.Errorf("listar números de ítems: %w", err)
	}
	var ids []string
	for rows.Next() {
		var numero int
		if err := rows.Scan(&numero); err != nil {
			rows.Close()
			return domain.ItemBacklog{}, fmt.Errorf("leer número de ítem: %w", err)
		}
		ids = append(ids, fmt.Sprintf("HU-%02d", numero))
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return domain.ItemBacklog{}, fmt.Errorf("recorrer números de ítems: %w", err)
	}
	rows.Close()

	item.ID = domain.SiguienteIDItem(ids)
	numero, ok := domain.NumeroDeItem(item.ID)
	if !ok {
		return domain.ItemBacklog{}, fmt.Errorf("el ID calculado para el ítem no es válido: %q", item.ID)
	}
	criterios := item.Criterios
	if criterios == nil {
		criterios = []string{}
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO items_backlog
			(proyecto_id, numero, titulo, descripcion, prioridad, estado, story_points, criterios_aceptacion)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, proyectoID, numero, item.Titulo, item.Descripcion, item.Prioridad,
		item.Estado, item.StoryPoints, criterios)
	if err != nil {
		return domain.ItemBacklog{}, fmt.Errorf("insertar ítem: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.ItemBacklog{}, fmt.Errorf("confirmar transacción de ítem: %w", err)
	}
	if item.Criterios == nil {
		item.Criterios = []string{}
	}
	return item, nil
}

func (r *Backlog) ObtenerItem(proyectoID, itemID string) (domain.ItemBacklog, error) {
	numero, ok := domain.NumeroDeItem(itemID)
	if !ok {
		return domain.ItemBacklog{}, domain.ErrItemNoEncontrado
	}
	item, err := scanItem(r.pool.QueryRow(context.Background(), `
		SELECT numero, titulo, descripcion, prioridad, estado, story_points, criterios_aceptacion
		FROM items_backlog
		WHERE proyecto_id = $1 AND numero = $2
	`, proyectoID, numero))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ItemBacklog{}, domain.ErrItemNoEncontrado
	}
	if err != nil {
		return domain.ItemBacklog{}, fmt.Errorf("obtener ítem: %w", err)
	}
	return item, nil
}

func (r *Backlog) ListarPorProyecto(proyectoID string) ([]domain.ItemBacklog, error) {
	rows, err := r.pool.Query(context.Background(), `
		SELECT numero, titulo, descripcion, prioridad, estado, story_points, criterios_aceptacion
		FROM items_backlog
		WHERE proyecto_id = $1
		ORDER BY numero
	`, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar ítems: %w", err)
	}
	defer rows.Close()

	items := make([]domain.ItemBacklog, 0)
	for rows.Next() {
		item, err := scanItem(rows)
		if err != nil {
			return nil, fmt.Errorf("leer ítem: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recorrer ítems: %w", err)
	}
	return items, nil
}

func (r *Backlog) ActualizarItem(proyectoID string, item domain.ItemBacklog) error {
	numero, ok := domain.NumeroDeItem(item.ID)
	if !ok {
		return domain.ErrItemNoEncontrado
	}
	criterios := item.Criterios
	if criterios == nil {
		criterios = []string{}
	}
	tag, err := r.pool.Exec(context.Background(), `
		UPDATE items_backlog
		SET titulo = $3, descripcion = $4, prioridad = $5, estado = $6,
			story_points = $7, criterios_aceptacion = $8
		WHERE proyecto_id = $1 AND numero = $2
	`, proyectoID, numero, item.Titulo, item.Descripcion, item.Prioridad,
		item.Estado, item.StoryPoints, criterios)
	if err != nil {
		return fmt.Errorf("actualizar ítem: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrItemNoEncontrado
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanProyecto(row rowScanner) (service.ProyectoRegistrado, error) {
	var proyecto service.ProyectoRegistrado
	var fechaInicio, fechaFin time.Time
	err := row.Scan(
		&proyecto.ID,
		&proyecto.Proyecto.Nombre,
		&fechaInicio,
		&fechaFin,
		&proyecto.Proyecto.Estado,
	)
	if err != nil {
		return service.ProyectoRegistrado{}, err
	}
	proyecto.Proyecto.FechaInicio = fechaUTC(fechaInicio)
	proyecto.Proyecto.FechaFin = fechaUTC(fechaFin)
	return proyecto, nil
}

func scanItem(row rowScanner) (domain.ItemBacklog, error) {
	var item domain.ItemBacklog
	var numero int
	if err := row.Scan(
		&numero,
		&item.Titulo,
		&item.Descripcion,
		&item.Prioridad,
		&item.Estado,
		&item.StoryPoints,
		&item.Criterios,
	); err != nil {
		return domain.ItemBacklog{}, err
	}
	item.ID = fmt.Sprintf("HU-%02d", numero)
	if item.Criterios == nil {
		item.Criterios = []string{}
	}
	return item, nil
}

func fechaPostgres(fecha time.Time) string {
	return fecha.UTC().Format("2006-01-02")
}

func fechaUTC(fecha time.Time) time.Time {
	anio, mes, dia := fecha.Date()
	return time.Date(anio, mes, dia, 0, 0, 0, 0, time.UTC)
}
