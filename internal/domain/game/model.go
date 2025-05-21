package game

import "time"

type GameSession struct {
	ID        string
	BotLevel  int
	Result    string
	StartedAt time.Time
	EndedAt   *time.Time
}

type Move struct {
	ID         int
	GameID     string
	MoveNumber int
	Color      string
	From       string
	To         string
	SAN        string
	FEN        string
	CreatedAt  time.Time
}
