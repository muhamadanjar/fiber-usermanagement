package interactors

import (
	"errors"

	"fiber-usermanagement/internal/domain/entities"
	"fiber-usermanagement/internal/domain/repositories"

	"gorm.io/gorm"
)

// AccessControlInteractor adalah use case untuk operasi terkait entitas AccessControl.
// Ini mengimplementasikan logika bisnis yang berinteraksi dengan AccessControlRepository.
type AccessControlInteractor struct {
	acRepo repositories.AccessControlRepository // Dependensi ke interface AccessControlRepository
}

// NewAccessControlInteractor membuat instance baru dari AccessControlInteractor.
// Menerima implementasi AccessControlRepository untuk dipasangkan.
func NewAccessControlInteractor(acr repositories.AccessControlRepository) *AccessControlInteractor {
	return &AccessControlInteractor{acRepo: acr}
}

// CreateAccessControl adalah use case untuk membuat access control baru.
// Ini menangani validasi input dasar dan memanggil repository untuk persistensi.
func (i *AccessControlInteractor) CreateAccessControl(ac *entities.AccessControl) (*entities.AccessControl, error) {
	// Validasi: field required tidak boleh kosong
	if ac.RoleID == 0 {
		return nil, errors.New("role ID tidak boleh kosong")
	}
	if ac.PermissionID == 0 {
		return nil, errors.New("permission ID tidak boleh kosong")
	}
	if ac.ResourceType == "" {
		return nil, errors.New("resource type tidak boleh kosong")
	}
	if ac.ResourceID == 0 {
		return nil, errors.New("resource ID tidak boleh kosong")
	}

	// Cek apakah access control dengan kombinasi yang sama sudah ada
	existingAC, err := i.acRepo.FindByRoleAndResource(ac.RoleID, ac.ResourceType, ac.ResourceID)
	if err == nil && existingAC != nil {
		return nil, errors.New("access control dengan kombinasi role, resource type, dan resource ID tersebut sudah ada")
	}

	// Panggil repository untuk menyimpan data
	return i.acRepo.Create(ac)
}

// GetAccessControlByID adalah use case untuk mendapatkan access control berdasarkan ID.
func (i *AccessControlInteractor) GetAccessControlByID(id uint) (*entities.AccessControl, error) {
	// Panggil repository untuk mengambil data
	ac, err := i.acRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("access control tidak ditemukan")
		}
		return nil, err
	}
	return ac, nil
}

// GetAllAccessControls adalah use case untuk mendapatkan semua access control.
func (i *AccessControlInteractor) GetAllAccessControls() ([]entities.AccessControl, error) {
	// Panggil repository untuk mengambil semua data
	return i.acRepo.FindAll()
}

// DeleteAccessControl adalah use case untuk menghapus access control.
func (i *AccessControlInteractor) DeleteAccessControl(id uint) error {
	// Cek apakah access control ada
	_, err := i.acRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("access control tidak ditemukan untuk dihapus")
		}
		return err
	}

	// Panggil repository untuk menghapus data
	return i.acRepo.Delete(id)
}
