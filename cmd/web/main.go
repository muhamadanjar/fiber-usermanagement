package main

import (
	"fiber-usermanagement/internal/config"
	"fiber-usermanagement/internal/container"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	log, err := config.NewZapLogger()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
	}

	viperConfig := config.NewViper()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Gagal Mengambil Config")
	}

	_, err = config.NewRabbitMQ(viperConfig)
	if err != nil {
		log.Fatal("Failed to initialize RabbitMQ", zap.Error(err))
	}

	// Inisialisasi Fiber app
	app := fiber.New()

	// Convert zap.Logger to log.Logger
	_, err = container.NewContainer(cfg, app, log)
	if err != nil {
		log.Fatal("Gagal menginisialisasi container aplikasi: %v", zap.Error(err))
	}

	// Mulai server Fiber
	err = app.Listen(fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		log.Fatal("Gagal menginisialisasi server: %v", zap.Error(err))
	}
}
