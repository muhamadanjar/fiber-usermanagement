package handlers

import (
	"log"

	"fiber-usermanagement/internal/domain/entities"
	"fiber-usermanagement/internal/usecase/interactors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// PermissionHandler menangani permintaan HTTP terkait entitas Permission.
type PermissionHandler struct {
	permissionInteractor *interactors.PermissionInteractor
}

// NewPermissionHandler membuat instance baru dari PermissionHandler.
func NewPermissionHandler(pi *interactors.PermissionInteractor) *PermissionHandler {
	return &PermissionHandler{permissionInteractor: pi}
}

// CreatePermission menangani pembuatan permission baru dari permintaan HTTP POST.
func (h *PermissionHandler) CreatePermission(c *fiber.Ctx) error {
	permission := new(entities.Permission)
	// Parse body permintaan ke struct Permission
	if err := c.BodyParser(permission); err != nil {
		log.Printf("Kesalahan BodyParser saat membuat permission: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Permintaan tidak valid"})
	}

	// Panggil use case untuk membuat permission
	createdPermission, err := h.permissionInteractor.CreatePermission(permission)
	if err != nil {
		log.Printf("Kesalahan CreatePermission di handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Kembalikan permission yang dibuat dengan status 201 Created
	return c.Status(fiber.StatusCreated).JSON(createdPermission)
}

// GetPermissionByID menangani pengambilan permission berdasarkan ID dari permintaan HTTP GET.
func (h *PermissionHandler) GetPermissionByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	// Parse UUID dari string parameter
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID permission tidak valid"})
	}

	// Panggil use case untuk mendapatkan permission
	permission, err := h.permissionInteractor.GetPermissionByID(id)
	if err != nil {
		log.Printf("Kesalahan GetPermissionByID di handler: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Permission tidak ditemukan"})
	}
	// Kembalikan permission yang ditemukan
	return c.JSON(permission)
}

// GetAllPermissions menangani pengambilan semua permission dari permintaan HTTP GET.
func (h *PermissionHandler) GetAllPermissions(c *fiber.Ctx) error {
	// Panggil use case untuk mendapatkan semua permission
	permissions, err := h.permissionInteractor.GetAllPermissions()
	if err != nil {
		log.Printf("Kesalahan GetAllPermissions di handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil daftar permission"})
	}
	// Kembalikan daftar permission
	return c.JSON(permissions)
}

// UpdatePermission menangani pembaruan permission yang ada dari permintaan HTTP PUT.
func (h *PermissionHandler) UpdatePermission(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID permission tidak valid"})
	}

	permission := new(entities.Permission)
	// Parse body permintaan ke struct Permission
	if err := c.BodyParser(permission); err != nil {
		log.Printf("Kesalahan BodyParser saat memperbarui permission: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Permintaan tidak valid"})
	}

	// Panggil use case untuk memperbarui permission
	updatedPermission, err := h.permissionInteractor.UpdatePermission(id, permission)
	if err != nil {
		log.Printf("Kesalahan UpdatePermission di handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Kembalikan permission yang diperbarui
	return c.JSON(updatedPermission)
}

// DeletePermission menangani penghapusan permission berdasarkan ID dari permintaan HTTP DELETE.
func (h *PermissionHandler) DeletePermission(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID permission tidak valid"})
	}

	// Panggil use case untuk menghapus permission
	err = h.permissionInteractor.DeletePermission(id)
	if err != nil {
		log.Printf("Kesalahan DeletePermission di handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Kembalikan status 204 No Content untuk penghapusan yang berhasil
	return c.Status(fiber.StatusNoContent).SendString("")
}
