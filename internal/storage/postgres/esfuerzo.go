package postgres

import (
	"context"
	"fmt"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
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
	registro domain.RegistroEsfuerzo,
) error {
	_, err := r.pool.Exec(context.Background(), `
		INSERT INTO esfuerzos (proyecto_id, item_numero, integrante, fecha, actividad, horas)
		VALUES ($1, $2, $3, $4::date, $5, $6)
	`, proyectoID, itemNumero, registro.Integrante, registro.Fecha.Format("2006-01-02"),
		registro.Actividad, registro.Horas)
	if err != nil {
		return fmt.Errorf("agregar esfuerzo: %w", err)
	}
	return nil
}
