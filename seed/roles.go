package seed

import (
	"example/APIForWorldWorkHub/models"
	"gorm.io/gorm"
)

var Roles = []models.Role{
	{Name: "Admin"},
	{Name: "Employee"},
	{Name: "User"},
}

func InitializeRoles(db *gorm.DB) {
	for _, role := range Roles {
		db.FirstOrCreate(&role, models.Role{Name: role.Name})
	}
}
