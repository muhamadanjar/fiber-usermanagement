package persistence

import (
	"fiber-usermanagement/internal/domain/entities"
	"fiber-usermanagement/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccessControlRepositoryImpl adalah implementasi konkret dari interface repositories.AccessControlRepository.
// Ini menggunakan GORM untuk berinteraksi dengan database.
type AccessControlRepositoryImpl struct {
	db *gorm.DB // Kumpulan koneksi database GORM
}

// NewAccessControlRepository membuat instance baru dari AccessControlRepositoryImpl.
// Ini menerima instance *gorm.DB untuk melakukan operasi database.
func NewAccessControlRepository(db *gorm.DB) repositories.AccessControlRepository {
	return &AccessControlRepositoryImpl{db: db}
}

// Create mengimplementasikan metode Create dari AccessControlRepository.
// Ini membuat record access control baru di database.
func (r *AccessControlRepositoryImpl) Create(ac *entities.AccessControl) (*entities.AccessControl, error) {
	result := r.db.Create(ac) // GORM akan mengisi ID setelah pembuatan berhasil
	return ac, result.Error
}

// FindByID mengimplementasikan metode FindByID dari AccessControlRepository.
// Ini mencari record access control berdasarkan ID dengan preload Role dan Permission.
func (r *AccessControlRepositoryImpl) FindByID(id uuid.UUID) (*entities.AccessControl, error) {
	var ac entities.AccessControl
	result := r.db.Preload("Role").Preload("Permission").Where("id = ?", id).First(&ac)
	return &ac, result.Error
}

// FindByRoleAndResource mengimplementasikan metode FindByRoleAndResource dari AccessControlRepository.
// Ini mencari access control berdasarkan role dan resource.
func (r *AccessControlRepositoryImpl) FindByRoleAndResource(roleID uuid.UUID, resourceType string, resourceID uuid.UUID) (*entities.AccessControl, error) {
	var ac entities.AccessControl
	result := r.db.Where("role_id = ? AND resource_type = ? AND resource_id = ?", roleID, resourceType, resourceID).
		Preload("Role").
		Preload("Permission").
		First(&ac)
	return &ac, result.Error
}

// FindAll mengimplementasikan metode FindAll dari AccessControlRepository.
// Ini mengembalikan semua record access control dari database dengan preload Role dan Permission.
func (r *AccessControlRepositoryImpl) FindAll() ([]entities.AccessControl, error) {
	var acs []entities.AccessControl
	result := r.db.Preload("Role").Preload("Permission").Find(&acs)
	return acs, result.Error
}

// Delete mengimplementasikan metode Delete dari AccessControlRepository.
// Ini menghapus record access control berdasarkan ID.
func (r *AccessControlRepositoryImpl) Delete(id uuid.UUID) error {
	// Menghapus record AccessControl berdasarkan ID
	result := r.db.Where("id = ?", id).Delete(&entities.AccessControl{})
	return result.Error
}
