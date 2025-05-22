package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/pamanaleph/chessforge-backend/internal/config"
	"github.com/pamanaleph/chessforge-backend/internal/router"
	"github.com/pamanaleph/chessforge-backend/internal/repository/postgres"
	gameService "github.com/pamanaleph/chessforge-backend/internal/service/game"
	"github.com/pamanaleph/chessforge-backend/internal/engine"
)

func main() {
	_ = godotenv.Load()

	app := fiber.New()
	cfg := config.Load()

	db := postgres.NewDB()
	repo := postgres.NewGameRepository(db)

	stockfishEngine, err := engine.NewStockfish("./bin/stockfish.exe")
	if err != nil {
		log.Fatal("Failed to start Stockfish:", err)
	}

	service := gameService.NewService(repo, stockfishEngine)

	router.Register(app, service)

	log.Println("Server is running on port", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
