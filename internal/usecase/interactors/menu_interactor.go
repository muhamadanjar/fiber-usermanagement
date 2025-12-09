package interactors

import (
	"errors"

	"fiber-usermanagement/internal/domain/entities"
	"fiber-usermanagement/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MenuInteractor adalah use case untuk operasi terkait entitas Menu.
// Ini mengimplementasikan logika bisnis yang berinteraksi dengan MenuRepository.
type MenuInteractor struct {
	menuRepo repositories.MenuRepository // Dependensi ke interface MenuRepository
}

// NewMenuInteractor membuat instance baru dari MenuInteractor.
// Menerima implementasi MenuRepository untuk dipasangkan.
func NewMenuInteractor(mr repositories.MenuRepository) *MenuInteractor {
	return &MenuInteractor{menuRepo: mr}
}

// CreateMenu adalah use case untuk membuat menu baru.
// Ini menangani validasi input dasar dan memanggil repository untuk persistensi.
func (i *MenuInteractor) CreateMenu(menu *entities.Menu) (*entities.Menu, error) {
	// Validasi: name dan URL tidak boleh kosong
	if menu.Name == "" {
		return nil, errors.New("nama menu tidak boleh kosong")
	}
	if menu.URL == "" {
		return nil, errors.New("URL menu tidak boleh kosong")
	}

	// Validasi parent-child relationship: cek apakah parent ada jika ParentID diset
	if menu.ParentID != nil {
		parent, err := i.menuRepo.FindByID(*menu.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("parent menu tidak ditemukan")
			}
			return nil, err
		}

		// Validasi: parent tidak boleh memiliki parent yang sama dengan menu yang akan dibuat
		// (mencegah circular reference sederhana)
		if parent.ParentID != nil && *parent.ParentID == menu.ID {
			return nil, errors.New("circular reference terdeteksi pada parent menu")
		}
	}

	// Panggil repository untuk menyimpan data
	return i.menuRepo.Create(menu)
}

// GetMenuByID adalah use case untuk mendapatkan menu berdasarkan ID.
func (i *MenuInteractor) GetMenuByID(id uuid.UUID) (*entities.Menu, error) {
	// Panggil repository untuk mengambil data
	menu, err := i.menuRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("menu tidak ditemukan")
		}
		return nil, err
	}
	return menu, nil
}

// GetAllMenus adalah use case untuk mendapatkan semua menu.
func (i *MenuInteractor) GetAllMenus() ([]entities.Menu, error) {
	// Panggil repository untuk mengambil semua data
	return i.menuRepo.FindAll()
}

// GetMenusByParent adalah use case untuk mendapatkan menu berdasarkan parent ID.
func (i *MenuInteractor) GetMenusByParent(parentID *uuid.UUID) ([]entities.Menu, error) {
	// Panggil repository untuk mengambil data berdasarkan parent
	return i.menuRepo.FindByParentID(parentID)
}

// UpdateMenu adalah use case untuk memperbarui menu.
// Ini mengambil menu yang ada, memperbarui bidang yang diizinkan, dan menyimpan perubahan.
func (i *MenuInteractor) UpdateMenu(id uuid.UUID, menu *entities.Menu) (*entities.Menu, error) {
	// Validasi: name dan URL tidak boleh kosong
	if menu.Name == "" {
		return nil, errors.New("nama menu tidak boleh kosong")
	}
	if menu.URL == "" {
		return nil, errors.New("URL menu tidak boleh kosong")
	}

	// Ambil menu yang ada terlebih dahulu
	existingMenu, err := i.menuRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("menu tidak ditemukan untuk diperbarui")
		}
		return nil, err
	}

	// Validasi parent-child relationship jika ParentID diubah
	if menu.ParentID != nil {
		// Cek apakah parent ada
		parent, err := i.menuRepo.FindByID(*menu.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("parent menu tidak ditemukan")
			}
			return nil, err
		}

		// Validasi: tidak boleh set parent ke diri sendiri
		if *menu.ParentID == id {
			return nil, errors.New("menu tidak boleh menjadi parent dari dirinya sendiri")
		}

		// Validasi: parent tidak boleh memiliki parent yang sama dengan menu yang akan diupdate
		if parent.ParentID != nil && *parent.ParentID == id {
			return nil, errors.New("circular reference terdeteksi pada parent menu")
		}
	}

	// Perbarui field yang diizinkan oleh logika bisnis
	existingMenu.Name = menu.Name
	existingMenu.Icon = menu.Icon
	existingMenu.URL = menu.URL
	existingMenu.Order = menu.Order
	existingMenu.IsActive = menu.IsActive
	existingMenu.ParentID = menu.ParentID

	// Panggil repository untuk menyimpan pembaruan
	return i.menuRepo.Update(existingMenu)
}

// DeleteMenu adalah use case untuk menghapus menu.
// Ini dapat mencakup logika bisnis pra-penghapusan.
func (i *MenuInteractor) DeleteMenu(id uuid.UUID) error {
	// Cek apakah menu ada
	menu, err := i.menuRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("menu tidak ditemukan untuk dihapus")
		}
		return err
	}

	// Logika bisnis: cek apakah menu memiliki children
	// Jika ada, tidak boleh dihapus (atau bisa juga cascade delete tergantung requirement)
	children, err := i.menuRepo.FindByParentID(&menu.ID)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return errors.New("menu memiliki submenu, tidak dapat dihapus")
	}

	// Panggil repository untuk menghapus data
	return i.menuRepo.Delete(id)
}
