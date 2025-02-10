package middlewares

import (
	"example/APIForWorldWorkHub/database"
	"example/APIForWorldWorkHub/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authorize(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID") // Assumindo que o ID do usuário está armazenado no contexto
		var user models.User
		if err := database.DB.Preload("Role.Permissions").Where("id = ?", userID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
			c.Abort()
			return
		}

		for _, p := range user.Role.Permissions {
			if p.Name == permission {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have permission to access this resource"})
		c.Abort()
	}
}
