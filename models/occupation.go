package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Occupation struct {
	ID     string `gorm:"type:uuid;primary_key"`
	Name   string `json:"name"`
	UserID string `json:"userId" gorm:"type:uuid" validate:"required"`
}

func (occupation *Occupation) BeforeCreate(tx *gorm.DB) (err error) {
	occupation.ID = uuid.New().String()
	return
}