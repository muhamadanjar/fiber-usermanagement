package repositories

import (
	"fiber-usermanagement/internal/domain/entities"

	"github.com/google/uuid"
)

// AccessControlRepository mendefinisikan interface untuk operasi data AccessControl.
// Interface ini akan diimplementasikan oleh layer infrastructure.
type AccessControlRepository interface {
	// Create membuat access control baru di database
	Create(ac *entities.AccessControl) (*entities.AccessControl, error)

	// FindByID mencari access control berdasarkan ID
	FindByID(id uuid.UUID) (*entities.AccessControl, error)

	// FindByRoleAndResource mencari access control berdasarkan role dan resource
	FindByRoleAndResource(roleID uuid.UUID, resourceType string, resourceID uuid.UUID) (*entities.AccessControl, error)

	// FindAll mengambil semua access control dari database
	FindAll() ([]entities.AccessControl, error)

	// Delete menghapus access control berdasarkan ID
	Delete(id uuid.UUID) error
}
