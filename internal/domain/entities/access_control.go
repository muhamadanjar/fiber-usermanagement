package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccessControl struct {
	ID           uuid.UUID `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleID       uuid.UUID `gorm:"not null;index;uniqueIndex:idx_role_permission_resource" json:"role_id"`
	PermissionID uuid.UUID `gorm:"not null;index;uniqueIndex:idx_role_permission_resource" json:"permission_id"`
	ResourceType string    `gorm:"not null;uniqueIndex:idx_resource_type" json:"resource_type"`
	ResourceID   uuid.UUID `gorm:"not null;uniqueIndex:idx_resource_id" json:"resource_id"`

	Role       *Role          `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Permission *Permission    `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ac *AccessControl) GetResource() string {
	return ac.ResourceType
}
