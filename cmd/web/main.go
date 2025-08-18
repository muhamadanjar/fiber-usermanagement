package main

import (
	"fiber-usermanagement/internal/config"
	"fiber-usermanagement/internal/container"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	logger, err := config.NewZapLogger(cfg)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	_, err = config.NewRabbitMQ(cfg)
	if err != nil {
		log.Fatal("Failed to initialize RabbitMQ", zap.Error(err))
	}

	// Inisialisasi Fiber app
	app := fiber.New()

	// Convert zap.Logger to log.Logger
	_, err = container.NewContainer(cfg, app, logger)
	if err != nil {
		log.Fatalf("Gagal menginisialisasi container aplikasi: %v", zap.Error(err))
	}

	// Mulai server Fiber
	err = app.Listen(fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		log.Fatalf("Gagal menginisialisasi server: %v", zap.Error(err))
	}
}
