package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Occupation struct {
	ID   string `gorm:"type:uuid;primary_key"`
	Name string `gorm:"unique"`
}

func (occupation *Occupation) BeforeCreate(_ *gorm.DB) (err error) {
	occupation.ID = uuid.New().String()
	return
}
