package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID          string       `gorm:"primaryKey"`
	Name        string       `gorm:"unique;not null"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

func (role *Role) BeforeCreate(_ *gorm.DB) (err error) {
	role.ID = uuid.New().String()
	return
}
