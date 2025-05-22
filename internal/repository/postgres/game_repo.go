package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pamanaleph/chessforge-backend/internal/domain/game"
)

type gameRepo struct {
	dbConn *pgx.Conn
}

func NewGameRepository(conn *pgx.Conn) game.GameRepository {
	return &gameRepo{dbConn: conn}
}

// CreateGame inserts a new game session
func (r *gameRepo) CreateGame(session *game.GameSession) error {
	query := `
		INSERT INTO game_sessions (id, bot_level, result, started_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.dbConn.Exec(context.Background(), query,
		session.ID, session.BotLevel, session.Result, session.StartedAt)

	return err
}

// SaveMove inserts a move made in the game
func (r *gameRepo) SaveMove(move *game.Move) error {
	query := `
		INSERT INTO moves (game_id, move_number, color, from_square, to_square, san, fen, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	return r.dbConn.QueryRow(context.Background(), query,
		move.GameID, move.MoveNumber, move.Color,
		move.From, move.To, move.SAN, move.FEN, move.CreatedAt,
	).Scan(&move.ID)
}

// GetMoves retrieves all moves for a game
func (r *gameRepo) GetMoves(gameID string) ([]game.Move, error) {
	query := `
		SELECT id, game_id, move_number, color, from_square, to_square, san, fen, created_at
		FROM moves
		WHERE game_id = $1
		ORDER BY move_number ASC
	`

	rows, err := r.dbConn.Query(context.Background(), query, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var moves []game.Move
	for rows.Next() {
		var m game.Move
		if err := rows.Scan(
			&m.ID, &m.GameID, &m.MoveNumber, &m.Color,
			&m.From, &m.To, &m.SAN, &m.FEN, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		moves = append(moves, m)
	}

	return moves, nil
}

// EndGame sets the result and ended_at of the game
func (r *gameRepo) EndGame(gameID string, result string) error {
	query := `
		UPDATE game_sessions
		SET result = $1, ended_at = $2
		WHERE id = $3
	`

	_, err := r.dbConn.Exec(context.Background(), query, result, time.Now(), gameID)
	return err
}
