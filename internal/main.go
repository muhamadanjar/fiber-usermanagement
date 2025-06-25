package main

import (
	"fiber-usermanagement/internal/api/routes"
	"fiber-usermanagement/internal/config"
	"fiber-usermanagement/internal/container"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Gagal Mengambil COnfig")
	}

	appContainer, err := container.NewContainer(cfg)
	if err != nil {
		log.Fatalf("Gagal menginisialisasi container aplikasi: %v", err)
	}

	// Inisialisasi Fiber app
	app := fiber.New()

	routes.InitRoutesModule(appContainer.UserInteractor)

	// API Layer (Routes)
	routes.SetupRoutes(app)

	// Mulai server Fiber
	log.Fatal(app.Listen(":" + cfg.Port))
}
