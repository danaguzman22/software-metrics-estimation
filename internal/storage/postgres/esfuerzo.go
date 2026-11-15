package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ service.RepositorioEsfuerzo = (*Esfuerzo)(nil)

// Esfuerzo implementa el repositorio PostgreSQL de registros de esfuerzo.
type Esfuerzo struct {
	pool *pgxpool.Pool
}

// NuevoEsfuerzo crea un repositorio de esfuerzo sobre el pool indicado.
func NuevoEsfuerzo(pool *pgxpool.Pool) *Esfuerzo {
	return &Esfuerzo{pool: pool}
}

// Agregar persiste un registro de esfuerzo asociado al proyecto y a la historia.
func (r *Esfuerzo) Agregar(
	proyectoID string,
	itemNumero int,
	sprintNumero *int,
	registro domain.RegistroEsfuerzo,
) error {
	_, err := r.pool.Exec(context.Background(), `
		INSERT INTO esfuerzos (proyecto_id, item_numero, sprint_numero, integrante, fecha, actividad, horas)
		VALUES ($1, $2, $3, $4, $5::date, $6, $7)
	`, proyectoID, itemNumero, sprintNumero, registro.Integrante, registro.Fecha.Format("2006-01-02"),
		registro.Actividad, registro.Horas)
	if err != nil {
		return fmt.Errorf("agregar esfuerzo: %w", err)
	}
	return nil
}

// HorasRealesPorItem suma los registros de esfuerzo de una historia.
func (r *Esfuerzo) HorasRealesPorItem(proyectoID string, itemNumero int) (float64, error) {
	return r.sumarHoras(`
		SELECT COALESCE(SUM(horas), 0)::double precision
		FROM esfuerzos WHERE proyecto_id = $1 AND item_numero = $2
	`, proyectoID, itemNumero)
}

// HorasRealesPorSprint suma los registros atribuidos explícitamente al Sprint.
func (r *Esfuerzo) HorasRealesPorSprint(proyectoID string, sprintNumero int) (float64, error) {
	return r.sumarHoras(`
		SELECT COALESCE(SUM(horas), 0)::double precision
		FROM esfuerzos WHERE proyecto_id = $1 AND sprint_numero = $2
	`, proyectoID, sprintNumero)
}

// HorasRealesPorProyecto suma todos los registros del proyecto, incluidos los
// registros históricos que no tenían atribución a un Sprint.
func (r *Esfuerzo) HorasRealesPorProyecto(proyectoID string) (float64, error) {
	return r.sumarHoras(`
		SELECT COALESCE(SUM(horas), 0)::double precision
		FROM esfuerzos WHERE proyecto_id = $1
	`, proyectoID)
}

func (r *Esfuerzo) sumarHoras(query string, argumentos ...any) (float64, error) {
	var total float64
	if err := r.pool.QueryRow(context.Background(), query, argumentos...).Scan(&total); err != nil {
		return 0, fmt.Errorf("sumar horas reales: %w", err)
	}
	return total, nil
}

// SnapshotSprint devuelve el snapshot persistido al cerrar el Sprint. El
// indicador distingue un snapshot vacío de un Sprint cerrado antes de US-17.
func (r *Esfuerzo) SnapshotSprint(
	proyectoID string,
	sprintNumero int,
) ([]domain.ItemSprintEsfuerzo, bool, error) {
	var disponible bool
	err := r.pool.QueryRow(context.Background(), `
		SELECT TRUE
		FROM sprint_esfuerzo_snapshots
		WHERE proyecto_id = $1 AND sprint_numero = $2
	`, proyectoID, sprintNumero).Scan(&disponible)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("verificar snapshot del sprint: %w", err)
	}

	rows, err := r.pool.Query(context.Background(), `
		SELECT item_numero, horas_estimadas::double precision
		FROM sprint_esfuerzo_items
		WHERE proyecto_id = $1 AND sprint_numero = $2
		ORDER BY item_numero
	`, proyectoID, sprintNumero)
	if err != nil {
		return nil, false, fmt.Errorf("listar snapshot del sprint: %w", err)
	}
	defer rows.Close()

	items := make([]domain.ItemSprintEsfuerzo, 0)
	for rows.Next() {
		var item domain.ItemSprintEsfuerzo
		if err := rows.Scan(&item.ItemNumero, &item.HorasEstimadas); err != nil {
			return nil, false, fmt.Errorf("leer snapshot del sprint: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, false, fmt.Errorf("recorrer snapshot del sprint: %w", err)
	}
	return items, true, nil
}
