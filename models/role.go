package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID     string `gorm:"type:uuid;primary_key"`
	Name   string `gorm:"unique"`
}

func (role *Role) BeforeCreate(tx *gorm.DB) (err error) {
	role.ID = uuid.New().String()
	return
}