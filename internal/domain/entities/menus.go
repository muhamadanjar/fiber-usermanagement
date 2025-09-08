package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Menu struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string         `gorm:"unique;not null" json:"name"`
	Icon      string         `json:"icon"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	ParentID  *uuid.UUID     `gorm:"type:uuid" json:"parent_id"`
	Parent    *Menu          `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children  []Menu         `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	URL       string         `json:"url"`
	Order     int            `json:"order"`
	Roles     []*Role        `gorm:"many2many:role_menus;" json:"roles,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
