package game

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	chess "github.com/notnil/chess"
	"github.com/pamanaleph/chessforge-backend/internal/domain/game"
	"github.com/pamanaleph/chessforge-backend/internal/engine"
)

type service struct {
	repo   game.GameRepository
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

	// Ambil semua langkah sebelumnya
	moves, err := s.repo.GetMoves(gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve existing moves: %w", err)
	}
	playerMove.MoveNumber = len(moves) + 1

	// Hitung giliran jalan
	isWhiteTurn := len(moves)%2 == 0
	if (isWhiteTurn && playerMove.Color != "white") || (!isWhiteTurn && playerMove.Color != "black") {
		return nil, fmt.Errorf("not your turn: it's %s's move", map[bool]string{true: "white", false: "black"}[isWhiteTurn])
	}

	// Build ulang state papan dari histori
	gameState := chess.NewGame()
	uci := chess.UCINotation{}
	for _, m := range moves {
		moveStr := fmt.Sprintf("%s%s", m.From, m.To)
		move, err := uci.Decode(gameState.Position(), moveStr)
		if err != nil {
			return nil, fmt.Errorf("invalid historical move: %w", err)
		}
		gameState.Move(move)
	}

	// Apply langkah player
	playerMoveStr := fmt.Sprintf("%s%s", playerMove.From, playerMove.To)
	playerMoveDecoded, err := uci.Decode(gameState.Position(), playerMoveStr)
	if err != nil {
		return nil, fmt.Errorf("invalid player move: %w", err)
	}
	gameState.Move(playerMoveDecoded)

	playerMove.SAN = gameState.Moves()[len(gameState.Moves())-1].String()
	playerMove.FEN = gameState.FEN()

	// Simpan langkah player
	if err := s.repo.SaveMove(&playerMove); err != nil {
		return nil, fmt.Errorf("failed to save player move: %w", err)
	}

	// Get bot move dari Stockfish
	bestMoveStr, err := s.engine.GetBestMove(playerMove.FEN)
	if err != nil || len(bestMoveStr) != 4 {
		return nil, fmt.Errorf("stockfish error: %v", err)
	}
	botFrom, botTo := bestMoveStr[:2], bestMoveStr[2:]

	// Apply bot move
	botMoveDecoded, err := uci.Decode(gameState.Position(), bestMoveStr)
	if err != nil {
		return nil, fmt.Errorf("invalid bot move: %w", err)
	}
	gameState.Move(botMoveDecoded)

	// Simpan langkah bot
	botMove := game.Move{
		GameID:     gameID,
		MoveNumber: playerMove.MoveNumber + 1,
		Color:      "black",
		From:       botFrom,
		To:         botTo,
		SAN:        gameState.Moves()[len(gameState.Moves())-1].String(),
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
