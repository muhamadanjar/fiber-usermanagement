package repositories

import (
	"fiber-usermanagement/internal/domain/entities"

	"github.com/google/uuid"
)

// MenuRepository mendefinisikan interface untuk operasi data Menu.
// Interface ini akan diimplementasikan oleh layer infrastructure.
type MenuRepository interface {
	// Create membuat menu baru di database
	Create(menu *entities.Menu) (*entities.Menu, error)

	// FindByID mencari menu berdasarkan ID
	FindByID(id uuid.UUID) (*entities.Menu, error)

	// FindAll mengambil semua menu dari database
	FindAll() ([]entities.Menu, error)

	// FindByParentID mengambil menu berdasarkan parent ID (untuk hierarchical menu)
	FindByParentID(parentID *uuid.UUID) ([]entities.Menu, error)

	// Update memperbarui menu yang sudah ada
	Update(menu *entities.Menu) (*entities.Menu, error)

	// Delete menghapus menu berdasarkan ID
	Delete(id uuid.UUID) error
}
