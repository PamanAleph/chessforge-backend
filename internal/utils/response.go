package utils

import "github.com/gofiber/fiber/v2"

type successResponse struct {
	Success  bool        `json:"success"`
	Messages string      `json:"messages"`
	Data     interface{} `json:"data"`
}

type errorResponse struct {
	Success  bool   `json:"success"`
	Messages string `json:"messages"`
}

func Success(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(successResponse{
		Success:  true,
		Messages: message,
		Data:     data,
	})
}

func Error(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(errorResponse{
		Success:  false,
		Messages: message,
	})
}
