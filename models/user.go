package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type User struct {
	ID               string     `gorm:"type:uuid;primary_key"`
	Firstname        string     `json:"firstname" validate:"required"`
	Lastname         string     `json:"lastname" validate:"required"`
	Email            string     `json:"email" gorm:"not null" validate:"required,email"`
	PasswordDigest   string     `json:"password_digest" gorm:"not null" validate:"required,password"`
	CPF              *string    `json:"CPF"`
	Role             string     `json:"role" validate:"required"`
	Contact          string     `json:"contact"`
	Occupation       Occupation `json:"occupation"`
	Phone            string     `json:"phone"`
	Education        string     `json:"education"`
	Region           *string    `json:"region"`
	ServiceDesc      *string    `json:"description"`
	Services         []Service  `gorm:"foreignkey:UserID"`
	SpokenLanguages  []Language `gorm:"foreignkey:UserID"`
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.New().String()
	return scope.SetColumn("ID", uuid)
}