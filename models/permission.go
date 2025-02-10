package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	ID   string `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

func (permission *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	permission.ID = uuid.New().String()
	return
}
