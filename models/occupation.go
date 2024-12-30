package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Occupation struct {
	ID     string `gorm:"type:uuid;primary_key"`
	Name   string `json:"name"`
	UserID string `json:"userId" gorm:"type:uuid" validate:"required"`
}

func (occupation *Occupation) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.New().String()
	return scope.SetColumn("ID", uuid)
}