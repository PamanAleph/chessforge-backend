package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pamanaleph/chessforge-backend/internal/handler/game"
)

func Register(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Chess Backend API running"})
	})

	gameGroup := app.Group("/game")
	game.RegisterRoutes(gameGroup)
}
