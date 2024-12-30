package models

import (
	"time"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Service struct {
	ID          string    `gorm:"type:uuid;primary_key"`
	UserID      string    `json:"userId" gorm:"type:uuid" validate:"required"`
	Date        time.Time `json:"date" validate:"required"`
	Pay         int64     `json:"pay" validate:"required,min=0"`
	Description string    `json:"description" validate:"required"`
}

func (service *Service) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.New().String()
	return scope.SetColumn("ID", uuid)
}