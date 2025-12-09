package persistence

import (
	"fiber-usermanagement/internal/domain/entities"
	"fiber-usermanagement/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PermissionRepositoryImpl adalah implementasi konkret dari interface repositories.PermissionRepository.
// Ini menggunakan GORM untuk berinteraksi dengan database.
type PermissionRepositoryImpl struct {
	db *gorm.DB // Kumpulan koneksi database GORM
}

// NewPermissionRepository membuat instance baru dari PermissionRepositoryImpl.
// Ini menerima instance *gorm.DB untuk melakukan operasi database.
func NewPermissionRepository(db *gorm.DB) repositories.PermissionRepository {
	return &PermissionRepositoryImpl{db: db}
}

// Create mengimplementasikan metode Create dari PermissionRepository.
// Ini membuat record permission baru di database.
func (r *PermissionRepositoryImpl) Create(permission *entities.Permission) (*entities.Permission, error) {
	result := r.db.Create(permission) // GORM akan mengisi ID setelah pembuatan berhasil
	return permission, result.Error
}

// FindByID mengimplementasikan metode FindByID dari PermissionRepository.
// Ini mencari record permission berdasarkan ID dengan preload Roles.
func (r *PermissionRepositoryImpl) FindByID(id uuid.UUID) (*entities.Permission, error) {
	var permission entities.Permission
	result := r.db.Preload("Roles").First(&permission, "id = ?", id)
	return &permission, result.Error
}

// FindByName mengimplementasikan metode FindByName dari PermissionRepository.
// Ini mencari record permission berdasarkan nama.
func (r *PermissionRepositoryImpl) FindByName(name string) (*entities.Permission, error) {
	var permission entities.Permission
	result := r.db.Where("name = ?", name).First(&permission)
	return &permission, result.Error
}

// FindAll mengimplementasikan metode FindAll dari PermissionRepository.
// Ini mengembalikan semua record permission dari database dengan preload Roles.
func (r *PermissionRepositoryImpl) FindAll() ([]entities.Permission, error) {
	var permissions []entities.Permission
	result := r.db.Preload("Roles").Find(&permissions)
	return permissions, result.Error
}

// Update mengimplementasikan metode Update dari PermissionRepository.
// Ini memperbarui record permission yang sudah ada di database.
func (r *PermissionRepositoryImpl) Update(permission *entities.Permission) (*entities.Permission, error) {
	// `Save` akan melakukan operasi update jika record dengan ID tersebut sudah ada
	result := r.db.Save(permission)
	return permission, result.Error
}

// Delete mengimplementasikan metode Delete dari PermissionRepository.
// Ini menghapus record permission berdasarkan ID.
func (r *PermissionRepositoryImpl) Delete(id uuid.UUID) error {
	// Menghapus record Permission berdasarkan ID
	result := r.db.Delete(&entities.Permission{}, "id = ?", id)
	return result.Error
}
