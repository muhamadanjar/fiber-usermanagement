package routes

import (
	"fiber-usermanagement/internal/api/handlers"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                  *fiber.App
	UserHandler          *handlers.UserHandler
	RoleHandler          *handlers.RoleHandler
	PermissionHandler    *handlers.PermissionHandler
	MenuHandler          *handlers.MenuHandler
	AccessControlHandler *handlers.AccessControlHandler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	// User routes (guest access)
	c.App.Post("/users", c.UserHandler.CreateUser)     // POST /api/v1/users untuk membuat pengguna baru
	c.App.Get("/users/:id", c.UserHandler.GetUserByID) // GET /api/v1/users/:id untuk mendapatkan pengguna berdasarkan ID
}

func (c *RouteConfig) SetupAuthRoute() {
	// User routes (authenticated access)
	c.App.Put("/users/:id", c.UserHandler.UpdateUser)    // PUT /api/v1/users/:id untuk memperbarui pengguna
	c.App.Delete("/users/:id", c.UserHandler.DeleteUser) // DELETE /api/v1/users/:id untuk menghapus pengguna
	c.App.Get("/users", c.UserHandler.GetAllUsers)       // GET /api/v1/users untuk mendapatkan semua pengguna

	// Role routes
	c.App.Post("/roles", c.RoleHandler.CreateRole)                        // POST /api/v1/roles
	c.App.Get("/roles", c.RoleHandler.GetAllRoles)                        // GET /api/v1/roles
	c.App.Get("/roles/:id", c.RoleHandler.GetRoleByID)                    // GET /api/v1/roles/:id
	c.App.Put("/roles/:id", c.RoleHandler.UpdateRole)                     // PUT /roles/:id
	c.App.Delete("/roles/:id", c.RoleHandler.DeleteRole)                  // DELETE /roles/:id
	c.App.Post("/roles/:id/permissions", c.RoleHandler.AssignPermissions) // POST /roles/:id/permissions

	// Permission routes
	c.App.Post("/permissions", c.PermissionHandler.CreatePermission)       // POST /permissions
	c.App.Get("/permissions", c.PermissionHandler.GetAllPermissions)       // GET /permissions
	c.App.Get("/permissions/:id", c.PermissionHandler.GetPermissionByID)   // GET /permissions/:id
	c.App.Put("/permissions/:id", c.PermissionHandler.UpdatePermission)    // PUT /permissions/:id
	c.App.Delete("/permissions/:id", c.PermissionHandler.DeletePermission) // DELETE /permissions/:id

	// Menu routes
	c.App.Post("/menus", c.MenuHandler.CreateMenu)                 // POST /menus
	c.App.Get("/menus", c.MenuHandler.GetAllMenus)                 // GET /menus
	c.App.Get("/menus/:id", c.MenuHandler.GetMenuByID)             // GET /menus/:id
	c.App.Get("/menus/parent/:id", c.MenuHandler.GetMenusByParent) // GET /menus/parent/:id atau /menus/parent/root
	c.App.Put("/menus/:id", c.MenuHandler.UpdateMenu)              // PUT /menus/:id
	c.App.Delete("/menus/:id", c.MenuHandler.DeleteMenu)           // DELETE /menus/:id

	// Access Control routes
	c.App.Post("/access-controls", c.AccessControlHandler.CreateAccessControl)       // POST /access-controls
	c.App.Get("/access-controls", c.AccessControlHandler.GetAllAccessControls)       // GET /access-controls
	c.App.Get("/access-controls/:id", c.AccessControlHandler.GetAccessControlByID)   // GET /access-controls/:id
	c.App.Delete("/access-controls/:id", c.AccessControlHandler.DeleteAccessControl) // DELETE /access-controls/:id
}
