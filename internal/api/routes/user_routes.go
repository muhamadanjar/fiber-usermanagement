package routes

import (
	"fiber-usermanagement/internal/api/handlers"
	"fiber-usermanagement/internal/api/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SetupUserRoutes menginisialisasi rute-rute yang berkaitan dengan User.
// Router Fiber dan UserInteractor disediakan untuk menghubungkan permintaan dengan logika bisnis.
func SetupUserRoutes(router fiber.Router) {
	// Buat instance UserHandler, yang akan menggunakan UserInteractor
	userHandler := handlers.NewUserHandler(routesConfig.UserInteractor)

	// Buat grup rute khusus untuk resource 'users'
	userRoutes := router.Group("/users")

	// Rute publik (tidak memerlukan autentikasi)
	userRoutes.Post("/", userHandler.CreateUser)    // POST /api/v1/users untuk membuat pengguna baru
	userRoutes.Get("/:id", userHandler.GetUserByID) // GET /api/v1/users/:id untuk mendapatkan pengguna berdasarkan ID

	// Rute yang memerlukan autentikasi
	// Terapkan middleware autentikasi ke rute-rute berikut
	userRoutes.Use(middlewares.AuthMiddleware)
	userRoutes.Put("/:id", userHandler.UpdateUser)    // PUT /api/v1/users/:id untuk memperbarui pengguna
	userRoutes.Delete("/:id", userHandler.DeleteUser) // DELETE /api/v1/users/:id untuk menghapus pengguna
	userRoutes.Get("/", userHandler.GetAllUsers)      // GET /api/v1/users untuk mendapatkan semua pengguna
}
