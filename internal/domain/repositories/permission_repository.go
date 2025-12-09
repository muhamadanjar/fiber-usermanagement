package repositories

import (
	"fiber-usermanagement/internal/domain/entities"

	"github.com/google/uuid"
)

// PermissionRepository mendefinisikan interface untuk operasi data Permission.
// Interface ini akan diimplementasikan oleh layer infrastructure.
type PermissionRepository interface {
	// Create membuat permission baru di database
	Create(permission *entities.Permission) (*entities.Permission, error)

	// FindByID mencari permission berdasarkan ID
	FindByID(id uuid.UUID) (*entities.Permission, error)

	// FindByName mencari permission berdasarkan nama
	FindByName(name string) (*entities.Permission, error)

	// FindAll mengambil semua permission dari database
	FindAll() ([]entities.Permission, error)

	// Update memperbarui permission yang sudah ada
	Update(permission *entities.Permission) (*entities.Permission, error)

	// Delete menghapus permission berdasarkan ID
	Delete(id uuid.UUID) error
}
