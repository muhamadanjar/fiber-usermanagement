package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User merepresentasikan entitas pengguna dalam domain aplikasi.
// Ini adalah definisi struktur data inti yang digunakan di seluruh aplikasi.
// Tag `gorm` digunakan untuk pemetaan ORM GORM ke kolom database.
// Tag `json` digunakan untuk serialisasi/deserialisasi JSON saat berinteraksi dengan API.
type User struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username    string         `gorm:"unique;not null" json:"username"`
	Email       string         `gorm:"unique;not null" json:"email"`
	Password    string         `gorm:"not null" json:"-"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	IsSuperuser bool           `gorm:"not null;default:false" json:"is_superuser"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	Roles       []*Role        `gorm:"many2many:user_roles;" json:"roles"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
