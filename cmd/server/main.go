package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/pamanaleph/chessforge-backend/internal/config"
	"github.com/pamanaleph/chessforge-backend/internal/router"
)

func main() {
	_ = godotenv.Load()

	app := fiber.New()

	cfg := config.Load()

	router.Register(app)

	log.Println("Server is running on port", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
