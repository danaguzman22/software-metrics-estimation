package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ service.RepositorioDefectos = (*Defectos)(nil)

// Defectos implementa el repositorio PostgreSQL de defectos.
type Defectos struct {
	pool *pgxpool.Pool
}

// NuevoDefectos crea un repositorio de defectos sobre el pool indicado.
func NuevoDefectos(pool *pgxpool.Pool) *Defectos {
	return &Defectos{pool: pool}
}

// Crear registra el defecto y devuelve su identidad persistida.
func (r *Defectos) Crear(proyectoID string, defecto domain.Defecto) (domain.Defecto, error) {
	existe, err := proyectoExiste(r.pool, proyectoID)
	if err != nil {
		return domain.Defecto{}, fmt.Errorf("verificar proyecto del defecto: %w", err)
	}
	if !existe {
		return domain.Defecto{}, domain.ErrProyectoNoEncontrado
	}
	itemNumero, ok := domain.NumeroDeItem(defecto.ItemID)
	if !ok {
		return domain.Defecto{}, domain.ErrItemNoEncontrado
	}

	var sprintResolucion any
	if defecto.SprintResolucion != nil {
		sprintResolucion = *defecto.SprintResolucion
	}
	err = r.pool.QueryRow(context.Background(), `
		INSERT INTO defectos (
			proyecto_id, descripcion, severidad, estado, item_numero,
			sprint_deteccion_numero, sprint_resolucion_numero
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, proyectoID, defecto.Descripcion, defecto.Severidad, defecto.Estado,
		itemNumero, defecto.SprintDeteccion, sprintResolucion,
	).Scan(&defecto.ID)
	if err != nil {
		if dominioErr := traducirErrorDefectoPostgres(err); dominioErr != nil {
			return domain.Defecto{}, dominioErr
		}
		return domain.Defecto{}, fmt.Errorf("crear defecto: %w", err)
	}
	return defecto, nil
}

// ListarPorProyecto devuelve los defectos del proyecto en orden de creación.
func (r *Defectos) ListarPorProyecto(proyectoID string) ([]domain.Defecto, error) {
	rows, err := r.pool.Query(context.Background(), `
		SELECT id, descripcion, severidad, estado, item_numero,
			sprint_deteccion_numero, sprint_resolucion_numero
		FROM defectos
		WHERE proyecto_id = $1
		ORDER BY id
	`, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar defectos: %w", err)
	}
	defer rows.Close()

	defectos := make([]domain.Defecto, 0)
	for rows.Next() {
		defecto, err := scanDefecto(rows)
		if err != nil {
			return nil, fmt.Errorf("leer defecto: %w", err)
		}
		defectos = append(defectos, defecto)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recorrer defectos: %w", err)
	}
	return defectos, nil
}

// Obtener devuelve el defecto identificado dentro del proyecto.
func (r *Defectos) Obtener(proyectoID string, defectoID int64) (domain.Defecto, error) {
	defecto, err := scanDefecto(r.pool.QueryRow(context.Background(), `
		SELECT id, descripcion, severidad, estado, item_numero,
			sprint_deteccion_numero, sprint_resolucion_numero
		FROM defectos
		WHERE proyecto_id = $1 AND id = $2
	`, proyectoID, defectoID))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Defecto{}, domain.ErrDefectoNoEncontrado
	}
	if err != nil {
		return domain.Defecto{}, fmt.Errorf("obtener defecto: %w", err)
	}
	return defecto, nil
}

// Actualizar persiste el estado y Sprint de resolución del defecto.
func (r *Defectos) Actualizar(proyectoID string, defecto domain.Defecto) error {
	if _, ok := domain.NumeroDeItem(defecto.ItemID); !ok {
		return domain.ErrItemNoEncontrado
	}
	var sprintResolucion any
	if defecto.SprintResolucion != nil {
		sprintResolucion = *defecto.SprintResolucion
	}
	tag, err := r.pool.Exec(context.Background(), `
		UPDATE defectos
		SET estado = $3, sprint_resolucion_numero = $4
		WHERE proyecto_id = $1 AND id = $2
	`, proyectoID, defecto.ID, defecto.Estado, sprintResolucion)
	if err != nil {
		if dominioErr := traducirErrorDefectoPostgres(err); dominioErr != nil {
			return dominioErr
		}
		return fmt.Errorf("actualizar defecto: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrDefectoNoEncontrado
	}
	return nil
}

type filaDefecto interface {
	Scan(dest ...any) error
}

func scanDefecto(fila filaDefecto) (domain.Defecto, error) {
	var defecto domain.Defecto
	var itemNumero int
	var sprintResolucion pgtype.Int4
	if err := fila.Scan(
		&defecto.ID,
		&defecto.Descripcion,
		&defecto.Severidad,
		&defecto.Estado,
		&itemNumero,
		&defecto.SprintDeteccion,
		&sprintResolucion,
	); err != nil {
		return domain.Defecto{}, err
	}
	defecto.ItemID = "HU-" + fmt.Sprintf("%02d", itemNumero)
	if sprintResolucion.Valid {
		numero := int(sprintResolucion.Int32)
		defecto.SprintResolucion = &numero
	}
	return defecto, nil
}

func traducirErrorDefectoPostgres(err error) error {
	var postgresError *pgconn.PgError
	if !errors.As(err, &postgresError) {
		return nil
	}
	switch postgresError.ConstraintName {
	case "defectos_item_fk":
		return domain.ErrItemNoEncontrado
	case "defectos_sprint_deteccion_fk", "defectos_sprint_resolucion_fk":
		return domain.ErrSprintNoEncontrado
	case "defectos_resuelto_sprint_check":
		return domain.ErrSprintResolucionDefectoObligatorio
	case "defectos_cronologia_check":
		return domain.ErrSprintResolucionAnteriorADeteccion
	}
	return nil
}

func proyectoExiste(pool *pgxpool.Pool, proyectoID string) (bool, error) {
	var existe bool
	if err := pool.QueryRow(context.Background(),
		"SELECT EXISTS (SELECT 1 FROM proyectos WHERE id = $1)", proyectoID,
	).Scan(&existe); err != nil {
		return false, fmt.Errorf("verificar proyecto: %w", err)
	}
	return existe, nil
}
