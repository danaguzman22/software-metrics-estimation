package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
	"github.com/danaguzman22/software-metrics-estimation/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ service.RepositorioRondas = (*RondasPoker)(nil)

// RondasPoker implementa el repositorio PostgreSQL de rondas de Planning Poker.
type RondasPoker struct {
	pool *pgxpool.Pool
}

// NuevoRondasPoker crea un repositorio de rondas sobre el pool indicado.
func NuevoRondasPoker(pool *pgxpool.Pool) *RondasPoker {
	return &RondasPoker{pool: pool}
}

// Iniciar inserta la ronda numero para la historia en estado Abierta.
// Toma SELECT ... FOR UPDATE sobre la fila de items_backlog para serializar
// las rondas de una misma historia. domain.ErrItemNoEncontrado si la historia
// no existe o el itemID está mal formado; domain.ErrRondaAbiertaExistente si
// ya hay una ronda Abierta (defensa de concurrencia; también traducida desde
// cualquier SQLSTATE 23505, sin importar la constraint: el índice parcial
// rondas_poker_historia_abierta_idx o la PK).
func (r *RondasPoker) Iniciar(proyectoID, itemID string, numero int) error {
	numeroItem, ok := domain.NumeroDeItem(itemID)
	if !ok {
		return domain.ErrItemNoEncontrado
	}
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("iniciar transacción para iniciar ronda: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var itemBloqueado int
	if err := tx.QueryRow(ctx, `
		SELECT numero
		FROM items_backlog
		WHERE proyecto_id = $1 AND numero = $2
		FOR UPDATE
	`, proyectoID, numeroItem).Scan(&itemBloqueado); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrItemNoEncontrado
		}
		return fmt.Errorf("bloquear historia para iniciar ronda: %w", err)
	}

	var hayAbierta bool
	if err := tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM rondas_poker
			WHERE proyecto_id = $1 AND item_numero = $2 AND estado = $3
		)
	`, proyectoID, numeroItem, domain.EstadoRondaAbierta).Scan(&hayAbierta); err != nil {
		return fmt.Errorf("verificar rondas abiertas de la historia: %w", err)
	}
	if hayAbierta {
		return domain.ErrRondaAbiertaExistente
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO rondas_poker (proyecto_id, item_numero, numero, estado)
		VALUES ($1, $2, $3, $4)
	`, proyectoID, numeroItem, numero, domain.EstadoRondaAbierta)
	if err != nil {
		var postgresError *pgconn.PgError
		if errors.As(err, &postgresError) && postgresError.Code == "23505" {
			return domain.ErrRondaAbiertaExistente
		}
		return fmt.Errorf("insertar ronda: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("confirmar transacción de inicio de ronda: %w", err)
	}
	return nil
}

// Obtener devuelve la ronda numero de la historia.
// domain.ErrItemNoEncontrado si el itemID está mal formado;
// domain.ErrRondaNoEncontrada si no existe.
func (r *RondasPoker) Obtener(proyectoID, itemID string, numero int) (domain.RondaPoker, error) {
	numeroItem, ok := domain.NumeroDeItem(itemID)
	if !ok {
		return domain.RondaPoker{}, domain.ErrItemNoEncontrado
	}
	var ronda domain.RondaPoker
	if err := r.pool.QueryRow(context.Background(), `
		SELECT numero, estado
		FROM rondas_poker
		WHERE proyecto_id = $1 AND item_numero = $2 AND numero = $3
	`, proyectoID, numeroItem, numero).Scan(&ronda.Numero, &ronda.Estado); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.RondaPoker{}, domain.ErrRondaNoEncontrada
		}
		return domain.RondaPoker{}, fmt.Errorf("obtener ronda: %w", err)
	}
	return ronda, nil
}

// ListarPorHistoria devuelve las rondas de la historia ordenadas por numero.
// domain.ErrItemNoEncontrado si el itemID está mal formado.
func (r *RondasPoker) ListarPorHistoria(proyectoID, itemID string) ([]domain.RondaPoker, error) {
	numeroItem, ok := domain.NumeroDeItem(itemID)
	if !ok {
		return nil, domain.ErrItemNoEncontrado
	}
	rows, err := r.pool.Query(context.Background(), `
		SELECT numero, estado
		FROM rondas_poker
		WHERE proyecto_id = $1 AND item_numero = $2
		ORDER BY numero ASC
	`, proyectoID, numeroItem)
	if err != nil {
		return nil, fmt.Errorf("listar rondas de la historia: %w", err)
	}
	defer rows.Close()

	rondas := make([]domain.RondaPoker, 0)
	for rows.Next() {
		var ronda domain.RondaPoker
		if err := rows.Scan(&ronda.Numero, &ronda.Estado); err != nil {
			return nil, fmt.Errorf("leer ronda: %w", err)
		}
		rondas = append(rondas, ronda)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recorrer rondas de la historia: %w", err)
	}
	return rondas, nil
}

// RegistrarVoto hace upsert del voto en la ronda dentro de una transacción,
// solo si está Abierta. domain.ErrItemNoEncontrado si el itemID está mal
// formado; domain.ErrRondaNoEncontrada si la ronda no existe;
// domain.ErrRondaRevelada si ya está Revelada.
func (r *RondasPoker) RegistrarVoto(proyectoID, itemID string, numero int, voto domain.VotoPoker) error {
	numeroItem, ok := domain.NumeroDeItem(itemID)
	if !ok {
		return domain.ErrItemNoEncontrado
	}
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("iniciar transacción para registrar voto: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var estado string
	if err := tx.QueryRow(ctx, `
		SELECT estado
		FROM rondas_poker
		WHERE proyecto_id = $1 AND item_numero = $2 AND numero = $3
		FOR UPDATE
	`, proyectoID, numeroItem, numero).Scan(&estado); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrRondaNoEncontrada
		}
		return fmt.Errorf("bloquear ronda para registrar voto: %w", err)
	}
	if estado != domain.EstadoRondaAbierta {
		return domain.ErrRondaRevelada
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO votos_poker (proyecto_id, item_numero, ronda_numero, integrante_nombre, valor)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (proyecto_id, item_numero, ronda_numero, integrante_nombre)
		DO UPDATE SET valor = EXCLUDED.valor
	`, proyectoID, numeroItem, numero, voto.Integrante, voto.Valor)
	if err != nil {
		return fmt.Errorf("registrar voto: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("confirmar transacción de voto: %w", err)
	}
	return nil
}

// YaVotaron devuelve los nombres de quienes votaron en la ronda, ordenados
// alfabéticamente sin distinguir mayúsculas y sin los valores.
// El orden usa COLLATE "C" para comparar por bytes, independiente de la
// colación de la base; coincide con el orden del lado de servicio/dominio.
// domain.ErrItemNoEncontrado si el itemID está mal formado.
func (r *RondasPoker) YaVotaron(proyectoID, itemID string, numero int) ([]string, error) {
	numeroItem, ok := domain.NumeroDeItem(itemID)
	if !ok {
		return nil, domain.ErrItemNoEncontrado
	}
	rows, err := r.pool.Query(context.Background(), `
		SELECT integrante_nombre
		FROM votos_poker
		WHERE proyecto_id = $1 AND item_numero = $2 AND ronda_numero = $3
		ORDER BY lower(integrante_nombre), integrante_nombre COLLATE "C"
	`, proyectoID, numeroItem, numero)
	if err != nil {
		return nil, fmt.Errorf("listar quiénes votaron: %w", err)
	}
	defer rows.Close()

	nombres := make([]string, 0)
	for rows.Next() {
		var nombre string
		if err := rows.Scan(&nombre); err != nil {
			return nil, fmt.Errorf("leer integrante de voto: %w", err)
		}
		nombres = append(nombres, nombre)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recorrer integrantes de la ronda: %w", err)
	}
	return nombres, nil
}

// Votos devuelve los votos de la ronda con sus valores, ordenados
// alfabéticamente sin distinguir mayúsculas. Solo se invoca con la ronda
// Revelada; el servicio lo decide.
// El orden usa COLLATE "C" para comparar por bytes, independiente de la
// colación de la base; coincide con el orden del lado de servicio/dominio.
// domain.ErrItemNoEncontrado si el itemID está mal formado.
func (r *RondasPoker) Votos(proyectoID, itemID string, numero int) ([]domain.VotoPoker, error) {
	numeroItem, ok := domain.NumeroDeItem(itemID)
	if !ok {
		return nil, domain.ErrItemNoEncontrado
	}
	rows, err := r.pool.Query(context.Background(), `
		SELECT integrante_nombre, valor
		FROM votos_poker
		WHERE proyecto_id = $1 AND item_numero = $2 AND ronda_numero = $3
		ORDER BY lower(integrante_nombre), integrante_nombre COLLATE "C"
	`, proyectoID, numeroItem, numero)
	if err != nil {
		return nil, fmt.Errorf("listar votos de la ronda: %w", err)
	}
	defer rows.Close()

	votos := make([]domain.VotoPoker, 0)
	for rows.Next() {
		var voto domain.VotoPoker
		if err := rows.Scan(&voto.Integrante, &voto.Valor); err != nil {
			return nil, fmt.Errorf("leer voto: %w", err)
		}
		votos = append(votos, voto)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recorrer votos de la ronda: %w", err)
	}
	return votos, nil
}

// Revelar pasa la ronda Abierta a Revelada; es irreversible.
// domain.ErrItemNoEncontrado si el itemID está mal formado;
// domain.ErrRondaNoEncontrada si no existe; domain.ErrRondaRevelada si ya
// estaba Revelada.
func (r *RondasPoker) Revelar(proyectoID, itemID string, numero int) error {
	numeroItem, ok := domain.NumeroDeItem(itemID)
	if !ok {
		return domain.ErrItemNoEncontrado
	}
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("iniciar transacción para revelar ronda: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	tag, err := tx.Exec(ctx, `
		UPDATE rondas_poker
		SET estado = $4
		WHERE proyecto_id = $1 AND item_numero = $2 AND numero = $3 AND estado = $5
	`, proyectoID, numeroItem, numero, domain.EstadoRondaRevelada, domain.EstadoRondaAbierta)
	if err != nil {
		return fmt.Errorf("actualizar estado de la ronda: %w", err)
	}
	if tag.RowsAffected() == 0 {
		var estado string
		if err := tx.QueryRow(ctx, `
			SELECT estado
			FROM rondas_poker
			WHERE proyecto_id = $1 AND item_numero = $2 AND numero = $3
		`, proyectoID, numeroItem, numero).Scan(&estado); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domain.ErrRondaNoEncontrada
			}
			return fmt.Errorf("consultar estado de la ronda: %w", err)
		}
		return domain.ErrRondaRevelada
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("confirmar transacción de revelado: %w", err)
	}
	return nil
}
