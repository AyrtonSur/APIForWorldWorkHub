package seed

import (
	"example/APIForWorldWorkHub/models"
	"gorm.io/gorm"
)

func InitializePermissions(db *gorm.DB) {
	permissions := []models.Permission{
		{Name: "view_users"},
		{Name: "view_user"},
		{Name: "create_service"},
		{Name: "update_user"},
		{Name: "delete_user"},
	}

	for _, permission := range permissions {
		db.FirstOrCreate(&permission, models.Permission{Name: permission.Name})
	}

	var adminRole models.Role
	db.Where("name = ?", "Admin").First(&adminRole)
	db.Model(&adminRole).Association("Permissions").Replace(permissions)

	var userRole models.Role
	db.Where("name = ?", "User").First(&userRole)
	db.Model(&userRole).Association("Permissions").Replace([]models.Permission{
		{Name: "view_user"},
		{Name: "create_service"},
	})
}