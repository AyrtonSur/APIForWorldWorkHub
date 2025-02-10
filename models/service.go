package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Service struct {
	ID          string    `gorm:"type:uuid;primary_key"`
	UserID      string    `json:"userId" gorm:"type:uuid;constraint:OnDelete:CASCADE;" validate:"required"`
	Date        time.Time `json:"date" validate:"required"`
	Pay         int64     `json:"pay" validate:"required,min=0"`
	Description string    `json:"description" validate:"required"`
}

func (service *Service) BeforeCreate(_ *gorm.DB) (err error) {
	service.ID = uuid.New().String()
	return
}
