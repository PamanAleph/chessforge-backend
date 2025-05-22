package game

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pamanaleph/chessforge-backend/internal/domain/game"
	"github.com/pamanaleph/chessforge-backend/internal/engine"
	chess "github.com/notnil/chess"
)

type service struct {
	repo game.GameRepository
	engine *engine.Stockfish
}

func NewService(repo game.GameRepository, sfEngine *engine.Stockfish) game.GameService {
	return &service{
		repo:   repo,
		engine: sfEngine,
	}
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
func (s *service) SubmitMove(gameID string, playerMove game.Move) ([]game.Move, error) {
	playerMove.GameID = gameID
	playerMove.CreatedAt = time.Now()

	gameState := chess.NewGame()
	if err := gameState.Position().UnmarshalText([]byte(playerMove.FEN)); err != nil {
		return nil, fmt.Errorf("invalid FEN from player: %w", err)
	}

	// Build and apply move from UCI
	uci := fmt.Sprintf("%s%s", playerMove.From, playerMove.To)
	legalMoves := gameState.ValidMoves()
	var found *chess.Move
	for _, m := range legalMoves {
		if m.S1().String() == playerMove.From && m.S2().String() == playerMove.To {
			found = m
			break
		}
	}
	if found == nil {
		return nil, fmt.Errorf("invalid player move: %s", uci)
	}
	gameState.Move(found)

	// Set SAN & FEN for player move
	playerMove.SAN = gameState.Moves()[len(gameState.Moves())-1].String()
	playerMove.FEN = gameState.FEN()

	// Simpan player move
	if err := s.repo.SaveMove(&playerMove); err != nil {
		return nil, fmt.Errorf("failed to save player move: %w", err)
	}

	// Get bot move from Stockfish
	bestMove, err := s.engine.GetBestMove(playerMove.FEN)
	if err != nil {
		return nil, fmt.Errorf("stockfish error: %w", err)
	}
	if len(bestMove) != 4 {
		return nil, fmt.Errorf("invalid bestmove from stockfish: %s", bestMove)
	}
	botFrom := bestMove[:2]
	botTo := bestMove[2:]

	// Apply bot move
	legalMoves = gameState.ValidMoves()
	var botFound *chess.Move
	for _, m := range legalMoves {
		if m.S1().String() == botFrom && m.S2().String() == botTo {
			botFound = m
			break
		}
	}
	if botFound == nil {
		return nil, fmt.Errorf("invalid bot move: %s", bestMove)
	}
	gameState.Move(botFound)

	botMove := game.Move{
		GameID:     gameID,
		MoveNumber: playerMove.MoveNumber + 1,
		Color:      "black",
		From:       botFrom,
		To:         botTo,
		SAN:        botFound.String(),
		FEN:        gameState.FEN(),
		CreatedAt:  time.Now(),
	}
	if err := s.repo.SaveMove(&botMove); err != nil {
		return nil, fmt.Errorf("failed to save bot move: %w", err)
	}

	return []game.Move{playerMove, botMove}, nil
}

// EndGame updates the result and time
func (s *service) EndGame(gameID, result string) error {
	return s.repo.EndGame(gameID, result)
}

// GetMoves returns all moves of a game
func (s *service) GetMoves(gameID string) ([]game.Move, error) {
	return s.repo.GetMoves(gameID)
}
