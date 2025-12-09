package persistence

import (
	"fiber-usermanagement/internal/domain/entities"
	"fiber-usermanagement/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MenuRepositoryImpl adalah implementasi konkret dari interface repositories.MenuRepository.
// Ini menggunakan GORM untuk berinteraksi dengan database.
type MenuRepositoryImpl struct {
	db *gorm.DB // Kumpulan koneksi database GORM
}

// NewMenuRepository membuat instance baru dari MenuRepositoryImpl.
// Ini menerima instance *gorm.DB untuk melakukan operasi database.
func NewMenuRepository(db *gorm.DB) repositories.MenuRepository {
	return &MenuRepositoryImpl{db: db}
}

// Create mengimplementasikan metode Create dari MenuRepository.
// Ini membuat record menu baru di database.
func (r *MenuRepositoryImpl) Create(menu *entities.Menu) (*entities.Menu, error) {
	result := r.db.Create(menu) // GORM akan mengisi ID setelah pembuatan berhasil
	return menu, result.Error
}

// FindByID mengimplementasikan metode FindByID dari MenuRepository.
// Ini mencari record menu berdasarkan ID dengan preload Parent, Children, dan Roles.
func (r *MenuRepositoryImpl) FindByID(id uuid.UUID) (*entities.Menu, error) {
	var menu entities.Menu
	result := r.db.Preload("Parent").Preload("Children").Preload("Roles").First(&menu, "id = ?", id)
	return &menu, result.Error
}

// FindAll mengimplementasikan metode FindAll dari MenuRepository.
// Ini mengembalikan semua record menu dari database dengan preload Parent, Children, dan Roles.
func (r *MenuRepositoryImpl) FindAll() ([]entities.Menu, error) {
	var menus []entities.Menu
	result := r.db.Preload("Parent").Preload("Children").Preload("Roles").Order("\"order\" ASC").Find(&menus)
	return menus, result.Error
}

// FindByParentID mengimplementasikan metode FindByParentID dari MenuRepository.
// Ini mencari menu berdasarkan parent ID untuk mendukung hierarchical menu.
// Jika parentID nil, akan mengembalikan menu root (menu tanpa parent).
func (r *MenuRepositoryImpl) FindByParentID(parentID *uuid.UUID) ([]entities.Menu, error) {
	var menus []entities.Menu
	var result *gorm.DB

	if parentID == nil {
		// Cari menu root (menu tanpa parent)
		result = r.db.Where("parent_id IS NULL").Preload("Children").Preload("Roles").Order("\"order\" ASC").Find(&menus)
	} else {
		// Cari menu dengan parent ID tertentu
		result = r.db.Where("parent_id = ?", parentID).Preload("Children").Preload("Roles").Order("\"order\" ASC").Find(&menus)
	}

	return menus, result.Error
}

// Update mengimplementasikan metode Update dari MenuRepository.
// Ini memperbarui record menu yang sudah ada di database.
func (r *MenuRepositoryImpl) Update(menu *entities.Menu) (*entities.Menu, error) {
	// `Save` akan melakukan operasi update jika record dengan ID tersebut sudah ada
	result := r.db.Save(menu)
	return menu, result.Error
}

// Delete mengimplementasikan metode Delete dari MenuRepository.
// Ini menghapus record menu berdasarkan ID.
func (r *MenuRepositoryImpl) Delete(id uuid.UUID) error {
	// Menghapus record Menu berdasarkan ID
	result := r.db.Delete(&entities.Menu{}, "id = ?", id)
	return result.Error
}
