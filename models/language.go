package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Language struct {
	ID     string `gorm:"type:uuid;primary_key"`
	Name   string `json:"name"`
	UserID string `json:"userId" gorm:"type:uuid" validate:"required"`
}

func (language *Language) BeforeCreate(_ *gorm.DB) (err error) {
	language.ID = uuid.New().String()
	return
}
