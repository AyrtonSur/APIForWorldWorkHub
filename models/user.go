package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID               string     `gorm:"type:uuid;primary_key"`
	Firstname        string     `json:"firstname" validate:"required"`
	Lastname         string     `json:"lastname" validate:"required"`
	Email            string     `json:"email" gorm:"uniqueIndex" validate:"required,email"`
	PasswordDigest   string     `json:"password_digest" validate:"required,password"`
	CPF              *string    `json:"CPF"`
	Role             string     `json:"role" validate:"required"`
	OccupationID     *string       
	Occupation       *Occupation `gorm:"foreignKey:OccupationID"`
	Phone            string     `json:"phone"`
	ZipCode          string     `json:"zipcode"`
	Education        string     `json:"education"`
	Region           string     `json:"region"`
	City             string     `json:"city"`
	Services         []Service  `gorm:"foreignkey:UserID"`
	SpokenLanguages  []Language `gorm:"foreignkey:UserID"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New().String()
	return
}