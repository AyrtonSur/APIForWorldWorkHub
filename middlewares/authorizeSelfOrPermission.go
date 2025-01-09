package middlewares

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "example/APIForWorldWorkHub/models"
  "example/APIForWorldWorkHub/database"
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