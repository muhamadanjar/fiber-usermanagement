package interactors

import (
	"errors"

	"fiber-usermanagement/internal/domain/entities"
	"fiber-usermanagement/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleInteractor adalah use case untuk operasi terkait entitas Role.
// Ini mengimplementasikan logika bisnis yang berinteraksi dengan RoleRepository.
type RoleInteractor struct {
	roleRepo repositories.RoleRepository // Dependensi ke interface RoleRepository
}

// NewRoleInteractor membuat instance baru dari RoleInteractor.
// Menerima implementasi RoleRepository untuk dipasangkan.
func NewRoleInteractor(rr repositories.RoleRepository) *RoleInteractor {
	return &RoleInteractor{roleRepo: rr}
}

// CreateRole adalah use case untuk membuat role baru.
// Ini menangani validasi input dasar dan memanggil repository untuk persistensi.
func (i *RoleInteractor) CreateRole(role *entities.Role) (*entities.Role, error) {
	// Validasi: name tidak boleh kosong
	if role.Name == "" {
		return nil, errors.New("nama role tidak boleh kosong")
	}

	// Cek apakah role dengan nama yang sama sudah ada
	existingRole, err := i.roleRepo.FindByName(role.Name)
	if err == nil && existingRole != nil {
		return nil, errors.New("role dengan nama tersebut sudah ada")
	}

	// Panggil repository untuk menyimpan data
	return i.roleRepo.Create(role)
}

// GetRoleByID adalah use case untuk mendapatkan role berdasarkan ID.
func (i *RoleInteractor) GetRoleByID(id uuid.UUID) (*entities.Role, error) {
	// Panggil repository untuk mengambil data
	role, err := i.roleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role tidak ditemukan")
		}
		return nil, err
	}
	return role, nil
}

// GetAllRoles adalah use case untuk mendapatkan semua role.
func (i *RoleInteractor) GetAllRoles() ([]entities.Role, error) {
	// Panggil repository untuk mengambil semua data
	return i.roleRepo.FindAll()
}

// UpdateRole adalah use case untuk memperbarui role.
// Ini mengambil role yang ada, memperbarui bidang yang diizinkan, dan menyimpan perubahan.
func (i *RoleInteractor) UpdateRole(id uuid.UUID, role *entities.Role) (*entities.Role, error) {
	// Validasi: name tidak boleh kosong
	if role.Name == "" {
		return nil, errors.New("nama role tidak boleh kosong")
	}

	// Ambil role yang ada terlebih dahulu
	existingRole, err := i.roleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role tidak ditemukan untuk diperbarui")
		}
		return nil, err
	}

	// Perbarui field yang diizinkan oleh logika bisnis
	existingRole.Name = role.Name
	existingRole.Description = role.Description

	// Panggil repository untuk menyimpan pembaruan
	return i.roleRepo.Update(existingRole)
}

// DeleteRole adalah use case untuk menghapus role.
// Ini dapat mencakup logika bisnis pra-penghapusan.
func (i *RoleInteractor) DeleteRole(id uuid.UUID) error {
	// Cek apakah role ada
	_, err := i.roleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role tidak ditemukan untuk dihapus")
		}
		return err
	}

	// Panggil repository untuk menghapus data
	return i.roleRepo.Delete(id)
}

// AssignPermissions adalah use case untuk menambahkan permissions ke role.
func (i *RoleInteractor) AssignPermissions(roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	// Validasi: permissionIDs tidak boleh kosong
	if len(permissionIDs) == 0 {
		return errors.New("permission IDs tidak boleh kosong")
	}

	// Cek apakah role ada
	_, err := i.roleRepo.FindByID(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role tidak ditemukan")
		}
		return err
	}

	// Panggil repository untuk assign permissions
	return i.roleRepo.AssignPermissions(roleID, permissionIDs)
}
