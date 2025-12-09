package persistence

import (
	"fiber-usermanagement/internal/domain/entities"
	"fiber-usermanagement/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleRepositoryImpl adalah implementasi konkret dari interface repositories.RoleRepository.
// Ini menggunakan GORM untuk berinteraksi dengan database.
type RoleRepositoryImpl struct {
	db *gorm.DB // Kumpulan koneksi database GORM
}

// NewRoleRepository membuat instance baru dari RoleRepositoryImpl.
// Ini menerima instance *gorm.DB untuk melakukan operasi database.
func NewRoleRepository(db *gorm.DB) repositories.RoleRepository {
	return &RoleRepositoryImpl{db: db}
}

// Create mengimplementasikan metode Create dari RoleRepository.
// Ini membuat record role baru di database.
func (r *RoleRepositoryImpl) Create(role *entities.Role) (*entities.Role, error) {
	result := r.db.Create(role) // GORM akan mengisi ID setelah pembuatan berhasil
	return role, result.Error
}

// FindByID mengimplementasikan metode FindByID dari RoleRepository.
// Ini mencari record role berdasarkan ID dengan preload Users.
func (r *RoleRepositoryImpl) FindByID(id uuid.UUID) (*entities.Role, error) {
	var role entities.Role
	result := r.db.Preload("Users").First(&role, "id = ?", id)
	return &role, result.Error
}

// FindByName mengimplementasikan metode FindByName dari RoleRepository.
// Ini mencari record role berdasarkan nama.
func (r *RoleRepositoryImpl) FindByName(name string) (*entities.Role, error) {
	var role entities.Role
	result := r.db.Where("name = ?", name).First(&role)
	return &role, result.Error
}

// FindAll mengimplementasikan metode FindAll dari RoleRepository.
// Ini mengembalikan semua record role dari database dengan preload Users.
func (r *RoleRepositoryImpl) FindAll() ([]entities.Role, error) {
	var roles []entities.Role
	result := r.db.Preload("Users").Find(&roles)
	return roles, result.Error
}

// Update mengimplementasikan metode Update dari RoleRepository.
// Ini memperbarui record role yang sudah ada di database.
func (r *RoleRepositoryImpl) Update(role *entities.Role) (*entities.Role, error) {
	// `Save` akan melakukan operasi update jika record dengan ID tersebut sudah ada
	result := r.db.Save(role)
	return role, result.Error
}

// Delete mengimplementasikan metode Delete dari RoleRepository.
// Ini menghapus record role berdasarkan ID.
func (r *RoleRepositoryImpl) Delete(id uuid.UUID) error {
	// Menghapus record Role berdasarkan ID
	result := r.db.Delete(&entities.Role{}, "id = ?", id)
	return result.Error
}

// AssignPermissions mengimplementasikan metode AssignPermissions dari RoleRepository.
// Ini menambahkan permissions ke role menggunakan GORM associations.
func (r *RoleRepositoryImpl) AssignPermissions(roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	// Cari role terlebih dahulu
	var role entities.Role
	if err := r.db.First(&role, "id = ?", roleID).Error; err != nil {
		return err
	}

	// Cari permissions berdasarkan IDs
	var permissions []*entities.Permission
	if err := r.db.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return err
	}

	// Replace associations (akan menghapus yang lama dan menambahkan yang baru)
	return r.db.Model(&role).Association("Permissions").Replace(permissions)
}
