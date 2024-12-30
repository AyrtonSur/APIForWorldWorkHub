package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Language struct {
	ID     string `gorm:"type:uuid;primary_key"`
	Name   string `json:"name"`
	UserID string `gorm:"type:uuid"`
}

func (service *Language) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.New().String()
	return scope.SetColumn("ID", uuid)
}