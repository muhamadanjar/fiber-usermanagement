package routes

import (
	"fiber-usermanagement/internal/api/handlers"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App         *fiber.App
	UserHandler *handlers.UserHandler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Post("/", c.UserHandler.CreateUser)    // POST /api/v1/users untuk membuat pengguna baru
	c.App.Get("/:id", c.UserHandler.GetUserByID) // GET /api/v1/users/:id untuk mendapatkan pengguna berdasarkan ID
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Put("/:id", c.UserHandler.UpdateUser)    // PUT /api/v1/users/:id untuk memperbarui pengguna
	c.App.Delete("/:id", c.UserHandler.DeleteUser) // DELETE /api/v1/users/:id untuk menghapus pengguna
	c.App.Get("/", c.UserHandler.GetAllUsers)      // GET /api/v1/users untuk mendapatkan semua pengguna

}
