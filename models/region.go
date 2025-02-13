package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Region struct {
	ID           string `gorm:"type:uuid;primary_key"`
	Name         string `gorm:"unique"`
	Abbreviation string `gorm:"unique"`
}

func (region *Region) BeforeCreate(_ *gorm.DB) (err error) {
	region.ID = uuid.New().String()
	return
}
