package handlers

import (
	"log"

	"fiber-usermanagement/internal/domain/entities"
	"fiber-usermanagement/internal/usecase/interactors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RoleHandler menangani permintaan HTTP terkait entitas Role.
type RoleHandler struct {
	roleInteractor *interactors.RoleInteractor
}

// NewRoleHandler membuat instance baru dari RoleHandler.
func NewRoleHandler(ri *interactors.RoleInteractor) *RoleHandler {
	return &RoleHandler{roleInteractor: ri}
}

// CreateRole menangani pembuatan role baru dari permintaan HTTP POST.
func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	role := new(entities.Role)
	// Parse body permintaan ke struct Role
	if err := c.BodyParser(role); err != nil {
		log.Printf("Kesalahan BodyParser saat membuat role: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Permintaan tidak valid"})
	}

	// Panggil use case untuk membuat role
	createdRole, err := h.roleInteractor.CreateRole(role)
	if err != nil {
		log.Printf("Kesalahan CreateRole di handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Kembalikan role yang dibuat dengan status 201 Created
	return c.Status(fiber.StatusCreated).JSON(createdRole)
}

// GetRoleByID menangani pengambilan role berdasarkan ID dari permintaan HTTP GET.
func (h *RoleHandler) GetRoleByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	// Parse UUID dari string parameter
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID role tidak valid"})
	}

	// Panggil use case untuk mendapatkan role
	role, err := h.roleInteractor.GetRoleByID(id)
	if err != nil {
		log.Printf("Kesalahan GetRoleByID di handler: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Role tidak ditemukan"})
	}
	// Kembalikan role yang ditemukan
	return c.JSON(role)
}

// GetAllRoles menangani pengambilan semua role dari permintaan HTTP GET.
func (h *RoleHandler) GetAllRoles(c *fiber.Ctx) error {
	// Panggil use case untuk mendapatkan semua role
	roles, err := h.roleInteractor.GetAllRoles()
	if err != nil {
		log.Printf("Kesalahan GetAllRoles di handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil daftar role"})
	}
	// Kembalikan daftar role
	return c.JSON(roles)
}

// UpdateRole menangani pembaruan role yang ada dari permintaan HTTP PUT.
func (h *RoleHandler) UpdateRole(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID role tidak valid"})
	}

	role := new(entities.Role)
	// Parse body permintaan ke struct Role
	if err := c.BodyParser(role); err != nil {
		log.Printf("Kesalahan BodyParser saat memperbarui role: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Permintaan tidak valid"})
	}

	// Panggil use case untuk memperbarui role
	updatedRole, err := h.roleInteractor.UpdateRole(id, role)
	if err != nil {
		log.Printf("Kesalahan UpdateRole di handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Kembalikan role yang diperbarui
	return c.JSON(updatedRole)
}

// DeleteRole menangani penghapusan role berdasarkan ID dari permintaan HTTP DELETE.
func (h *RoleHandler) DeleteRole(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID role tidak valid"})
	}

	// Panggil use case untuk menghapus role
	err = h.roleInteractor.DeleteRole(id)
	if err != nil {
		log.Printf("Kesalahan DeleteRole di handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Kembalikan status 204 No Content untuk penghapusan yang berhasil
	return c.Status(fiber.StatusNoContent).SendString("")
}

// AssignPermissions menangani penambahan permissions ke role dari permintaan HTTP POST.
func (h *RoleHandler) AssignPermissions(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID role tidak valid"})
	}

	// Parse body untuk mendapatkan array permission IDs
	var requestBody struct {
		PermissionIDs []uuid.UUID `json:"permission_ids"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		log.Printf("Kesalahan BodyParser saat assign permissions: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Permintaan tidak valid"})
	}

	// Panggil use case untuk assign permissions
	err = h.roleInteractor.AssignPermissions(id, requestBody.PermissionIDs)
	if err != nil {
		log.Printf("Kesalahan AssignPermissions di handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Kembalikan status 200 OK
	return c.JSON(fiber.Map{"message": "Permissions berhasil ditambahkan ke role"})
}
