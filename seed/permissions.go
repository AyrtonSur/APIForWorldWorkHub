package seed

import (
	"example/APIForWorldWorkHub/models"
	"gorm.io/gorm"
)

func InitializePermissions(db *gorm.DB) error {
	permissions := []models.Permission{
		{Name: "view_users"},
		{Name: "view_user"},
		{Name: "create_service"},
		{Name: "update_user"},
		{Name: "delete_user"},
		{Name: "view_services"},
	}

	for i, permission := range permissions {
		perm := permission
		if err := db.FirstOrCreate(&perm, models.Permission{Name: permission.Name}).Error; err != nil {
			return err
		}
		permissions[i] = perm
	}

	var adminRole models.Role
	if err := db.Where("name = ?", "Admin").First(&adminRole).Error; err != nil {
		adminRole = models.Role{Name: "Admin"}
		if err := db.Create(&adminRole).Error; err != nil {
			return err
		}
	}

	// Associa as permissões à role "Admin"
	for _, permission := range permissions {
		var count int64
		db.Table("role_permissions").Where("role_id = ? AND permission_id = ?", adminRole.ID, permission.ID).Count(&count)
		if count == 0 {
			if err := db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", adminRole.ID, permission.ID).Error; err != nil {
				return err
			}
		}
	}

	var userRole models.Role
	if err := db.Where("name = ?", "User").First(&userRole).Error; err != nil {
		userRole = models.Role{Name: "User"}
		if err := db.Create(&userRole).Error; err != nil {
			return err
		}
	}

	// Associa permissões específicas à role "User"
	userPermissions := []models.Permission{}
	if err := db.Where("name IN ?", []string{"view_user", "create_service"}).Find(&userPermissions).Error; err != nil {
		return err
	}
	for _, permission := range userPermissions {
		var count int64
		db.Table("role_permissions").Where("role_id = ? AND permission_id = ?", userRole.ID, permission.ID).Count(&count)
		if count == 0 {
			if err := db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", userRole.ID, permission.ID).Error; err != nil {
				return err
			}
		}
	}

	return nil
}