package routes

import (
	"fiber-usermanagement/usecase/interactors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes menginisialisasi semua rute API aplikasi.
// Ini menerima instance Fiber app dan interactor yang diperlukan.
func SetupRoutes(app *fiber.App, userInteractor *interactors.UserInteractor) {
	// Middleware global Fiber
	app.Use(logger.New()) // Aktifkan middleware logging untuk setiap permintaan

	// Grouping API versi 1 untuk prefix '/api/v1'
	v1 := app.Group("/api/v1")

	// Setup rute-rute spesifik untuk entitas User
	SetupUserRoutes(v1, userInteractor)

	// Setup rute-rute spesifik untuk entitas Product

	// Rute default atau root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Selamat datang di Go Fiber Clean Architecture!")
	})
}
