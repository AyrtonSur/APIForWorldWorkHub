package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Language struct {
	ID     string `gorm:"type:uuid;primary_key"`
	Name   string `json:"name"`
	UserID string `json:"userId" gorm:"type:uuid" validate:"required"`
}

func (service *Language) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.New().String()
	return scope.SetColumn("ID", uuid)
}