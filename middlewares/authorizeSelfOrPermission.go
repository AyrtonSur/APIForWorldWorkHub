package middlewares

import (
	"example/APIForWorldWorkHub/database"
	"example/APIForWorldWorkHub/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthorizeSelfOrPermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		targetID := c.Param("id")

		// Se for o próprio usuário, libera
		if userID == targetID {
			c.Next()
			return
		}

		// Senão, checar permissão
		var user models.User
		if err := database.DB.Preload("Role.Permissions").Where("id = ?", userID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
			c.Abort()
			return
		}

		for _, p := range user.Role.Permissions {
			if p.Name == permission {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have permission to perform this action"})
		c.Abort()
	}
}
