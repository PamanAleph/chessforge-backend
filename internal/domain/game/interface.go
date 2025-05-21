package game

type GameRepository interface {
	CreateGame(session *GameSession) error
	SaveMove(move *Move) error
	GetMoves(gameID string) ([]Move, error)
	EndGame(gameID string, result string) error
}

type GameService interface {
	StartGame(botLevel int) (*GameSession, error)
	SubmitMove(gameID string, move Move) (*Move, error)
	EndGame(gameID, result string) error
	GetMoves(gameID string) ([]Move, error)
}
