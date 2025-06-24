package handlers

import (
	"log"
	"strconv" // Untuk konversi string ke uint

	"fiber-usermanagement/domain/entities"
	"fiber-usermanagement/usecase/interactors"

	"github.com/gofiber/fiber/v2"
)

// UserHandler menangani permintaan HTTP terkait entitas User.
type UserHandler struct {
	userInteractor *interactors.UserInteractor
}

// NewUserHandler membuat instance baru dari UserHandler.
func NewUserHandler(ui *interactors.UserInteractor) *UserHandler {
	return &UserHandler{userInteractor: ui}
}

// CreateUser menangani pembuatan pengguna baru dari permintaan HTTP POST.
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	user := new(entities.User)
	// Parse body permintaan ke struct User
	if err := c.BodyParser(user); err != nil {
		log.Printf("Kesalahan BodyParser saat membuat pengguna: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Permintaan tidak valid"})
	}

	// Panggil use case untuk membuat pengguna
	createdUser, err := h.userInteractor.CreateUser(user)
	if err != nil {
		log.Printf("Kesalahan CreateUser di handler: %v", err)
		// Sesuaikan status error berdasarkan jenis error dari use case
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat pengguna"})
	}
	// Kembalikan pengguna yang dibuat dengan status 201 Created
	return c.Status(fiber.StatusCreated).JSON(createdUser)
}

// GetUserByID menangani pengambilan pengguna berdasarkan ID dari permintaan HTTP GET.
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	// Konversi ID dari string parameter ke uint
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID pengguna tidak valid"})
	}

	// Panggil use case untuk mendapatkan pengguna
	user, err := h.userInteractor.GetUserByID(uint(id))
	if err != nil {
		log.Printf("Kesalahan GetUserByID di handler: %v", err)
		// Jika pengguna tidak ditemukan, kembalikan 404 Not Found
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pengguna tidak ditemukan"})
	}
	// Kembalikan pengguna yang ditemukan
	return c.JSON(user)
}

// GetAllUsers menangani pengambilan semua pengguna dari permintaan HTTP GET.
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	// Panggil use case untuk mendapatkan semua pengguna
	users, err := h.userInteractor.GetAllUsers()
	if err != nil {
		log.Printf("Kesalahan GetAllUsers di handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil daftar pengguna"})
	}
	// Kembalikan daftar pengguna
	return c.JSON(users)
}

// UpdateUser menangani pembaruan pengguna yang ada dari permintaan HTTP PUT.
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID pengguna tidak valid"})
	}

	user := new(entities.User)
	// Parse body permintaan ke struct User
	if err := c.BodyParser(user); err != nil {
		log.Printf("Kesalahan BodyParser saat memperbarui pengguna: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Permintaan tidak valid"})
	}

	// Panggil use case untuk memperbarui pengguna
	updatedUser, err := h.userInteractor.UpdateUser(uint(id), user)
	if err != nil {
		log.Printf("Kesalahan UpdateUser di handler: %v", err)
		// Sesuaikan status error berdasarkan jenis error dari use case
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui pengguna"})
	}
	// Kembalikan pengguna yang diperbarui
	return c.JSON(updatedUser)
}

// DeleteUser menangani penghapusan pengguna berdasarkan ID dari permintaan HTTP DELETE.
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID pengguna tidak valid"})
	}

	// Panggil use case untuk menghapus pengguna
	err = h.userInteractor.DeleteUser(uint(id))
	if err != nil {
		log.Printf("Kesalahan DeleteUser di handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghapus pengguna"})
	}
	// Kembalikan status 204 No Content untuk penghapusan yang berhasil
	return c.Status(fiber.StatusNoContent).SendString("")
}
