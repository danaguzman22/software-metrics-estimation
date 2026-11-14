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
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ service.RepositorioSprints = (*Sprints)(nil)

// Sprints implementa el repositorio PostgreSQL de Sprints.
type Sprints struct {
	pool *pgxpool.Pool
}

// NuevoSprints crea un repositorio de Sprints sobre el pool indicado.
func NuevoSprints(pool *pgxpool.Pool) *Sprints {
	return &Sprints{pool: pool}
}

// Agregar registra un Sprint en el proyecto.
// Una violación de la clave primaria se traduce a domain.ErrSprintDuplicado y
// una violación del índice único parcial a domain.ErrSprintActivoExistente.
func (r *Sprints) Agregar(proyectoID string, sprint domain.Sprint) error {
	_, err := r.pool.Exec(context.Background(), `
		INSERT INTO sprints (proyecto_id, numero, objetivo, fecha_inicio, fecha_fin, estado)
		VALUES ($1, $2, $3, $4::date, $5::date, $6)
	`, proyectoID, sprint.Numero, sprint.Objetivo,
		fechaPostgres(sprint.FechaInicio), fechaPostgres(sprint.FechaFin), sprint.Estado)
	if err != nil {
		var postgresError *pgconn.PgError
		if errors.As(err, &postgresError) && postgresError.Code == "23505" {
			switch postgresError.ConstraintName {
			case "sprints_proyecto_activo_idx":
				return domain.ErrSprintActivoExistente
			case "sprints_pkey":
				return domain.ErrSprintDuplicado
			}
		}
		return fmt.Errorf("agregar sprint: %w", err)
	}
	return nil
}

// ExisteNumero informa si el proyecto ya tiene un Sprint con ese número.
func (r *Sprints) ExisteNumero(proyectoID string, numero int) (bool, error) {
	var existe bool
	if err := r.pool.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1 FROM sprints WHERE proyecto_id = $1 AND numero = $2
		)
	`, proyectoID, numero).Scan(&existe); err != nil {
		return false, fmt.Errorf("verificar número de sprint: %w", err)
	}
	return existe, nil
}

// ExisteActivo informa si el proyecto tiene un Sprint en estado Activo.
func (r *Sprints) ExisteActivo(proyectoID string) (bool, error) {
	var existe bool
	if err := r.pool.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1 FROM sprints WHERE proyecto_id = $1 AND estado = $2
		)
	`, proyectoID, domain.EstadoSprintActivo).Scan(&existe); err != nil {
		return false, fmt.Errorf("verificar sprint activo: %w", err)
	}
	return existe, nil
}

// Obtener devuelve el Sprint identificado dentro del proyecto.
func (r *Sprints) Obtener(proyectoID string, numero int) (domain.Sprint, error) {
	sprint, err := scanSprint(r.pool.QueryRow(context.Background(), `
		SELECT numero, objetivo, fecha_inicio, fecha_fin, estado, sp_planificados
		FROM sprints
		WHERE proyecto_id = $1 AND numero = $2
	`, proyectoID, numero))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Sprint{}, domain.ErrSprintNoEncontrado
	}
	if err != nil {
		return domain.Sprint{}, fmt.Errorf("obtener sprint: %w", err)
	}
	return sprint, nil
}

// ListarPorProyecto devuelve los Sprints del proyecto ordenados por número.
func (r *Sprints) ListarPorProyecto(proyectoID string) ([]domain.Sprint, error) {
	rows, err := r.pool.Query(context.Background(), `
		SELECT numero, objetivo, fecha_inicio, fecha_fin, estado, sp_planificados
		FROM sprints
		WHERE proyecto_id = $1
		ORDER BY numero
	`, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar sprints del proyecto: %w", err)
	}
	defer rows.Close()

	sprints := make([]domain.Sprint, 0)
	for rows.Next() {
		sprint, err := scanSprint(rows)
		if err != nil {
			return nil, fmt.Errorf("leer sprint: %w", err)
		}
		sprints = append(sprints, sprint)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recorrer sprints del proyecto: %w", err)
	}
	return sprints, nil
}

// Cerrar actualiza el Sprint activo y elimina las asignaciones no completadas
// en una única transacción.
func (r *Sprints) Cerrar(proyectoID string, numero, spPlanificados int) error {
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("iniciar transacción para cerrar sprint: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	tag, err := tx.Exec(ctx, `
		UPDATE sprints
		SET estado = $3, sp_planificados = $4
		WHERE proyecto_id = $1 AND numero = $2 AND estado = $5
	`, proyectoID, numero, domain.EstadoSprintCerrado, spPlanificados, domain.EstadoSprintActivo)
	if err != nil {
		return fmt.Errorf("actualizar estado del sprint: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrSprintCerrado
	}

	if _, err := tx.Exec(ctx, `
		DELETE FROM sprint_items AS si
		USING items_backlog AS ib
		WHERE si.proyecto_id = $1
			AND si.sprint_numero = $2
			AND ib.proyecto_id = si.proyecto_id
			AND ib.numero = si.item_numero
			AND ib.estado <> $3
	`, proyectoID, numero, domain.EstadoCompletada); err != nil {
		return fmt.Errorf("eliminar asignaciones no completadas: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("confirmar transacción de cierre del sprint: %w", err)
	}
	return nil
}

func scanSprint(row rowScanner) (domain.Sprint, error) {
	var sprint domain.Sprint
	var fechaInicio, fechaFin time.Time
	var spPlanificados pgtype.Int8
	if err := row.Scan(
		&sprint.Numero,
		&sprint.Objetivo,
		&fechaInicio,
		&fechaFin,
		&sprint.Estado,
		&spPlanificados,
	); err != nil {
		return domain.Sprint{}, err
	}
	sprint.FechaInicio = fechaUTC(fechaInicio)
	sprint.FechaFin = fechaUTC(fechaFin)
	if spPlanificados.Valid {
		valor := int(spPlanificados.Int64)
		sprint.SPPlanificados = &valor
	}
	return sprint, nil
}
