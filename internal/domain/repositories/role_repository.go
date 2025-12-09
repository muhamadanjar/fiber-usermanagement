package repositories

import (
	"fiber-usermanagement/internal/domain/entities"

	"github.com/google/uuid"
)

// RoleRepository mendefinisikan interface untuk operasi data Role.
// Interface ini akan diimplementasikan oleh layer infrastructure.
type RoleRepository interface {
	// Create membuat role baru di database
	Create(role *entities.Role) (*entities.Role, error)

	// FindByID mencari role berdasarkan ID
	FindByID(id uuid.UUID) (*entities.Role, error)

	// FindByName mencari role berdasarkan nama
	FindByName(name string) (*entities.Role, error)

	// FindAll mengambil semua role dari database
	FindAll() ([]entities.Role, error)

	// Update memperbarui role yang sudah ada
	Update(role *entities.Role) (*entities.Role, error)

	// Delete menghapus role berdasarkan ID
	Delete(id uuid.UUID) error

	// AssignPermissions menambahkan permissions ke role
	AssignPermissions(roleID uuid.UUID, permissionIDs []uuid.UUID) error
}
