package game

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router) {
	router.Post("/start", StartGame)
}

func StartGame(c *fiber.Ctx) error {
	// Temporary dummy response
	return c.JSON(fiber.Map{
		"message": "Game vs Bot started (MVP 1)",
	})
}
