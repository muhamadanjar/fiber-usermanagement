package handlers

import (
	"log"
	"strconv"

	"fiber-usermanagement/internal/domain/entities"
	"fiber-usermanagement/internal/usecase/interactors"

	"github.com/gofiber/fiber/v2"
)

// AccessControlHandler menangani permintaan HTTP terkait entitas AccessControl.
type AccessControlHandler struct {
	acInteractor *interactors.AccessControlInteractor
}

// NewAccessControlHandler membuat instance baru dari AccessControlHandler.
func NewAccessControlHandler(aci *interactors.AccessControlInteractor) *AccessControlHandler {
	return &AccessControlHandler{acInteractor: aci}
}

// CreateAccessControl menangani pembuatan access control baru dari permintaan HTTP POST.
func (h *AccessControlHandler) CreateAccessControl(c *fiber.Ctx) error {
	ac := new(entities.AccessControl)
	// Parse body permintaan ke struct AccessControl
	if err := c.BodyParser(ac); err != nil {
		log.Printf("Kesalahan BodyParser saat membuat access control: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Permintaan tidak valid"})
	}

	// Panggil use case untuk membuat access control
	createdAC, err := h.acInteractor.CreateAccessControl(ac)
	if err != nil {
		log.Printf("Kesalahan CreateAccessControl di handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Kembalikan access control yang dibuat dengan status 201 Created
	return c.Status(fiber.StatusCreated).JSON(createdAC)
}

// GetAccessControlByID menangani pengambilan access control berdasarkan ID dari permintaan HTTP GET.
func (h *AccessControlHandler) GetAccessControlByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	// Konversi ID dari string parameter ke uint
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID access control tidak valid"})
	}

	// Panggil use case untuk mendapatkan access control
	ac, err := h.acInteractor.GetAccessControlByID(uint(id))
	if err != nil {
		log.Printf("Kesalahan GetAccessControlByID di handler: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Access control tidak ditemukan"})
	}
	// Kembalikan access control yang ditemukan
	return c.JSON(ac)
}

// GetAllAccessControls menangani pengambilan semua access control dari permintaan HTTP GET.
func (h *AccessControlHandler) GetAllAccessControls(c *fiber.Ctx) error {
	// Panggil use case untuk mendapatkan semua access control
	acs, err := h.acInteractor.GetAllAccessControls()
	if err != nil {
		log.Printf("Kesalahan GetAllAccessControls di handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil daftar access control"})
	}
	// Kembalikan daftar access control
	return c.JSON(acs)
}

// DeleteAccessControl menangani penghapusan access control berdasarkan ID dari permintaan HTTP DELETE.
func (h *AccessControlHandler) DeleteAccessControl(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID access control tidak valid"})
	}

	// Panggil use case untuk menghapus access control
	err = h.acInteractor.DeleteAccessControl(uint(id))
	if err != nil {
		log.Printf("Kesalahan DeleteAccessControl di handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Kembalikan status 204 No Content untuk penghapusan yang berhasil
	return c.Status(fiber.StatusNoContent).SendString("")
}
