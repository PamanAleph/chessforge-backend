package game

import (
	"time"

	"github.com/google/uuid"
	"github.com/pamanaleph/chessforge-backend/internal/domain/game"
)

type service struct {
	repo game.GameRepository
}

func NewService(repo game.GameRepository) game.GameService {
	return &service{repo: repo}
}

// StartGame initializes a new game session against a bot
func (s *service) StartGame(botLevel int) (*game.GameSession, error) {
	session := &game.GameSession{
		ID:        uuid.New().String(),
		BotLevel:  botLevel,
		Result:    "ongoing",
		StartedAt: time.Now(),
	}

	err := s.repo.CreateGame(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// SubmitMove stores a player or bot move
func (s *service) SubmitMove(gameID string, move game.Move) (*game.Move, error) {
	move.GameID = gameID
	move.CreatedAt = time.Now()

	err := s.repo.SaveMove(&move)
	if err != nil {
		return nil, err
	}

	return &move, nil
}

// EndGame updates the result and time
func (s *service) EndGame(gameID, result string) error {
	return s.repo.EndGame(gameID, result)
}

// GetMoves returns all moves of a game
func (s *service) GetMoves(gameID string) ([]game.Move, error) {
	return s.repo.GetMoves(gameID)
}
