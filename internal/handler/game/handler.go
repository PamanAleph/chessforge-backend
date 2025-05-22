package game

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pamanaleph/chessforge-backend/internal/domain/game"
	"github.com/pamanaleph/chessforge-backend/internal/utils"
)

type Handler struct {
	service game.GameService
}

func NewHandler(service game.GameService) *Handler {
	return &Handler{service: service}
}

func RegisterRoutes(router fiber.Router, handler *Handler) {
	router.Post("/start", handler.StartGame)
	router.Post("/:id/move", handler.SubmitMove)
}

type startGameRequest struct {
	BotLevel int `json:"bot_level"`
}

func (h *Handler) StartGame(c *fiber.Ctx) error {
	var req startGameRequest
	if err := c.BodyParser(&req); err != nil || req.BotLevel < 1 || req.BotLevel > 10 {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid bot_level")
	}

	session, err := h.service.StartGame(req.BotLevel)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "Failed to start game")
	}

	return utils.Success(c, "Game started successfully", fiber.Map{
		"game_id":    session.ID,
		"bot_level":  session.BotLevel,
		"started_at": session.StartedAt.Format("2006-01-02 15:04:05"),
	})
}

type moveRequest struct {
	MoveNumber int    `json:"move_number"`
	Color      string `json:"color"` // "white" or "black"
	From       string `json:"from"`
	To         string `json:"to"`
	SAN        string `json:"san"`
	FEN        string `json:"fen"`
}

func (h *Handler) SubmitMove(c *fiber.Ctx) error {
	gameID := c.Params("id")
	if gameID == "" {
		return utils.Error(c, fiber.StatusBadRequest, "Missing game ID")
	}

	var req moveRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid move payload")
	}

	move := game.Move{
		GameID:     gameID,
		MoveNumber: req.MoveNumber,
		Color:      req.Color,
		From:       req.From,
		To:         req.To,
		SAN:        req.SAN,
		FEN:        req.FEN,
	}

	savedMove, err := h.service.SubmitMove(gameID, move)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "Failed to save move")
	}

	return utils.Success(c, "Move saved successfully", savedMove)
}
