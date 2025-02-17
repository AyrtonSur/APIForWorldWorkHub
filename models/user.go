package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              string  `gorm:"type:uuid;primary_key"`
	Firstname       string  `json:"firstname" validate:"required"`
	Lastname        string  `json:"lastname" validate:"required"`
	Email           string  `json:"email" gorm:"uniqueIndex" validate:"required,email"`
	PasswordDigest  string  `json:"password_digest" validate:"required,password"`
	CPF             *string `json:"CPF" gorm:"uniqueIndex"`
	RoleID          string
	Role            Role `gorm:"foreignKey:RoleID"`
	OccupationID    *string
	Occupation      *Occupation `gorm:"foreignKey:OccupationID"`
	Phone           string      `json:"phone"`
	ZipCode         string      `json:"zipcode"`
	Education       string      `json:"education"`
	RegionID        string
	Region          Region     `gorm:"foreignKey:RegionID"`
	City            string     `json:"city"`
	RefreshToken    *string    `gorm:"unique"`
	Services        []Service  `gorm:"foreignkey:UserID"`
	SpokenLanguages []Language `gorm:"foreignkey:UserID"`
}

func (user *User) BeforeCreate(_ *gorm.DB) (err error) {
	user.ID = uuid.New().String()
	return
}
