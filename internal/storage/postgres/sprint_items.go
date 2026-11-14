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

var _ service.RepositorioAsignaciones = (*Asignaciones)(nil)

// Asignaciones implementa el repositorio PostgreSQL de la relación
// historia–Sprint.
type Asignaciones struct {
	pool *pgxpool.Pool
}

// NuevoAsignaciones crea un repositorio de asignaciones sobre el pool
// indicado.
func NuevoAsignaciones(pool *pgxpool.Pool) *Asignaciones {
	return &Asignaciones{pool: pool}
}

// Agregar registra que la historia itemNumero pertenece a sprintNumero en el
// proyecto. Una violación de la clave primaria se traduce a
// domain.ErrItemYaAsignado y las violaciones de las claves foráneas a
// domain.ErrSprintNoEncontrado y domain.ErrItemNoEncontrado (D7).
func (r *Asignaciones) Agregar(proyectoID string, sprintNumero, itemNumero int) error {
	_, err := r.pool.Exec(context.Background(), `
		INSERT INTO sprint_items (proyecto_id, sprint_numero, item_numero)
		VALUES ($1, $2, $3)
	`, proyectoID, sprintNumero, itemNumero)
	if err != nil {
		var postgresError *pgconn.PgError
		if errors.As(err, &postgresError) {
			switch {
			case postgresError.Code == "23505" && postgresError.ConstraintName == "sprint_items_pkey":
				return domain.ErrItemYaAsignado
			case postgresError.Code == "23503" && postgresError.ConstraintName == "sprint_items_sprint_fk":
				return domain.ErrSprintNoEncontrado
			case postgresError.Code == "23503" && postgresError.ConstraintName == "sprint_items_item_fk":
				return domain.ErrItemNoEncontrado
			}
		}
		return fmt.Errorf("agregar asignación de ítem: %w", err)
	}
	return nil
}

// ExisteItem informa si la historia ya está asignada a algún Sprint del
// proyecto.
func (r *Asignaciones) ExisteItem(proyectoID string, itemNumero int) (bool, error) {
	var existe bool
	if err := r.pool.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1 FROM sprint_items WHERE proyecto_id = $1 AND item_numero = $2
		)
	`, proyectoID, itemNumero).Scan(&existe); err != nil {
		return false, fmt.Errorf("verificar asignación de ítem: %w", err)
	}
	return existe, nil
}

// ExisteItemEnSprint informa si la historia está asignada al Sprint indicado.
func (r *Asignaciones) ExisteItemEnSprint(
	proyectoID string,
	sprintNumero, itemNumero int,
) (bool, error) {
	var existe bool
	if err := r.pool.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1
			FROM sprint_items
			WHERE proyecto_id = $1 AND sprint_numero = $2 AND item_numero = $3
		)
	`, proyectoID, sprintNumero, itemNumero).Scan(&existe); err != nil {
		return false, fmt.Errorf("verificar asignación al sprint: %w", err)
	}
	return existe, nil
}

// ListarNumerosItem devuelve los números de historia asignados al Sprint,
// ordenados por número ascendente.
func (r *Asignaciones) ListarNumerosItem(proyectoID string, sprintNumero int) ([]int, error) {
	rows, err := r.pool.Query(context.Background(), `
		SELECT item_numero
		FROM sprint_items
		WHERE proyecto_id = $1 AND sprint_numero = $2
		ORDER BY item_numero
	`, proyectoID, sprintNumero)
	if err != nil {
		return nil, fmt.Errorf("listar ítems del sprint: %w", err)
	}
	defer rows.Close()

	numeros := make([]int, 0)
	for rows.Next() {
		var numero int
		if err := rows.Scan(&numero); err != nil {
			return nil, fmt.Errorf("leer ítem del sprint: %w", err)
		}
		numeros = append(numeros, numero)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("leer ítems del sprint: %w", err)
	}
	return numeros, nil
}
