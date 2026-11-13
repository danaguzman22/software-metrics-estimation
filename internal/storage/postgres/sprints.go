package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
	"github.com/jackc/pgx/v5/pgconn"
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
