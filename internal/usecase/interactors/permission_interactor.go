package interactors

import (
	"errors"

	"fiber-usermanagement/internal/domain/entities"
	"fiber-usermanagement/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PermissionInteractor adalah use case untuk operasi terkait entitas Permission.
// Ini mengimplementasikan logika bisnis yang berinteraksi dengan PermissionRepository.
type PermissionInteractor struct {
	permissionRepo repositories.PermissionRepository // Dependensi ke interface PermissionRepository
}

// NewPermissionInteractor membuat instance baru dari PermissionInteractor.
// Menerima implementasi PermissionRepository untuk dipasangkan.
func NewPermissionInteractor(pr repositories.PermissionRepository) *PermissionInteractor {
	return &PermissionInteractor{permissionRepo: pr}
}

// CreatePermission adalah use case untuk membuat permission baru.
// Ini menangani validasi input dasar dan memanggil repository untuk persistensi.
func (i *PermissionInteractor) CreatePermission(permission *entities.Permission) (*entities.Permission, error) {
	// Validasi: name tidak boleh kosong
	if permission.Name == "" {
		return nil, errors.New("nama permission tidak boleh kosong")
	}

	// Cek apakah permission dengan nama yang sama sudah ada
	existingPermission, err := i.permissionRepo.FindByName(permission.Name)
	if err == nil && existingPermission != nil {
		return nil, errors.New("permission dengan nama tersebut sudah ada")
	}

	// Panggil repository untuk menyimpan data
	return i.permissionRepo.Create(permission)
}

// GetPermissionByID adalah use case untuk mendapatkan permission berdasarkan ID.
func (i *PermissionInteractor) GetPermissionByID(id uuid.UUID) (*entities.Permission, error) {
	// Panggil repository untuk mengambil data
	permission, err := i.permissionRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission tidak ditemukan")
		}
		return nil, err
	}
	return permission, nil
}

// GetAllPermissions adalah use case untuk mendapatkan semua permission.
func (i *PermissionInteractor) GetAllPermissions() ([]entities.Permission, error) {
	// Panggil repository untuk mengambil semua data
	return i.permissionRepo.FindAll()
}

// UpdatePermission adalah use case untuk memperbarui permission.
// Ini mengambil permission yang ada, memperbarui bidang yang diizinkan, dan menyimpan perubahan.
func (i *PermissionInteractor) UpdatePermission(id uuid.UUID, permission *entities.Permission) (*entities.Permission, error) {
	// Validasi: name tidak boleh kosong
	if permission.Name == "" {
		return nil, errors.New("nama permission tidak boleh kosong")
	}

	// Ambil permission yang ada terlebih dahulu
	existingPermission, err := i.permissionRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission tidak ditemukan untuk diperbarui")
		}
		return nil, err
	}

	// Perbarui field yang diizinkan oleh logika bisnis
	existingPermission.Name = permission.Name
	existingPermission.Description = permission.Description

	// Panggil repository untuk menyimpan pembaruan
	return i.permissionRepo.Update(existingPermission)
}

// DeletePermission adalah use case untuk menghapus permission.
// Ini dapat mencakup logika bisnis pra-penghapusan.
func (i *PermissionInteractor) DeletePermission(id uuid.UUID) error {
	// Cek apakah permission ada
	_, err := i.permissionRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("permission tidak ditemukan untuk dihapus")
		}
		return err
	}

	// Panggil repository untuk menghapus data
	return i.permissionRepo.Delete(id)
}
