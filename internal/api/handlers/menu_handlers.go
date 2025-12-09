package handlers

import (
	"log"

	"fiber-usermanagement/internal/domain/entities"
	"fiber-usermanagement/internal/usecase/interactors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// MenuHandler menangani permintaan HTTP terkait entitas Menu.
type MenuHandler struct {
	menuInteractor *interactors.MenuInteractor
}

// NewMenuHandler membuat instance baru dari MenuHandler.
func NewMenuHandler(mi *interactors.MenuInteractor) *MenuHandler {
	return &MenuHandler{menuInteractor: mi}
}

// CreateMenu menangani pembuatan menu baru dari permintaan HTTP POST.
func (h *MenuHandler) CreateMenu(c *fiber.Ctx) error {
	menu := new(entities.Menu)
	// Parse body permintaan ke struct Menu
	if err := c.BodyParser(menu); err != nil {
		log.Printf("Kesalahan BodyParser saat membuat menu: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Permintaan tidak valid"})
	}

	// Panggil use case untuk membuat menu
	createdMenu, err := h.menuInteractor.CreateMenu(menu)
	if err != nil {
		log.Printf("Kesalahan CreateMenu di handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Kembalikan menu yang dibuat dengan status 201 Created
	return c.Status(fiber.StatusCreated).JSON(createdMenu)
}

// GetMenuByID menangani pengambilan menu berdasarkan ID dari permintaan HTTP GET.
func (h *MenuHandler) GetMenuByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	// Parse UUID dari string parameter
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID menu tidak valid"})
	}

	// Panggil use case untuk mendapatkan menu
	menu, err := h.menuInteractor.GetMenuByID(id)
	if err != nil {
		log.Printf("Kesalahan GetMenuByID di handler: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Menu tidak ditemukan"})
	}
	// Kembalikan menu yang ditemukan
	return c.JSON(menu)
}

// GetAllMenus menangani pengambilan semua menu dari permintaan HTTP GET.
func (h *MenuHandler) GetAllMenus(c *fiber.Ctx) error {
	// Panggil use case untuk mendapatkan semua menu
	menus, err := h.menuInteractor.GetAllMenus()
	if err != nil {
		log.Printf("Kesalahan GetAllMenus di handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil daftar menu"})
	}
	// Kembalikan daftar menu
	return c.JSON(menus)
}

// GetMenusByParent menangani pengambilan menu berdasarkan parent ID dari permintaan HTTP GET.
func (h *MenuHandler) GetMenusByParent(c *fiber.Ctx) error {
	idStr := c.Params("id")

	var parentID *uuid.UUID
	// Jika id adalah "root", maka parentID = nil (untuk menu root)
	if idStr == "root" {
		parentID = nil
	} else {
		// Parse UUID dari string parameter
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID parent tidak valid"})
		}
		parentID = &id
	}

	// Panggil use case untuk mendapatkan menu berdasarkan parent
	menus, err := h.menuInteractor.GetMenusByParent(parentID)
	if err != nil {
		log.Printf("Kesalahan GetMenusByParent di handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil menu berdasarkan parent"})
	}
	// Kembalikan daftar menu
	return c.JSON(menus)
}

// UpdateMenu menangani pembaruan menu yang ada dari permintaan HTTP PUT.
func (h *MenuHandler) UpdateMenu(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID menu tidak valid"})
	}

	menu := new(entities.Menu)
	// Parse body permintaan ke struct Menu
	if err := c.BodyParser(menu); err != nil {
		log.Printf("Kesalahan BodyParser saat memperbarui menu: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Permintaan tidak valid"})
	}

	// Panggil use case untuk memperbarui menu
	updatedMenu, err := h.menuInteractor.UpdateMenu(id, menu)
	if err != nil {
		log.Printf("Kesalahan UpdateMenu di handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Kembalikan menu yang diperbarui
	return c.JSON(updatedMenu)
}

// DeleteMenu menangani penghapusan menu berdasarkan ID dari permintaan HTTP DELETE.
func (h *MenuHandler) DeleteMenu(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID menu tidak valid"})
	}

	// Panggil use case untuk menghapus menu
	err = h.menuInteractor.DeleteMenu(id)
	if err != nil {
		log.Printf("Kesalahan DeleteMenu di handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Kembalikan status 204 No Content untuk penghapusan yang berhasil
	return c.Status(fiber.StatusNoContent).SendString("")
}
