package middlewares

import (
	"github.com/gin-gonic/gin"
	"example/APIForWorldWorkHub/utils"
	"net/http"
)

func Authenticate() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Token necessário"})
			context.Abort()
			return
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil || claims == nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			context.Abort()
			return
		}

		context.Set("userID", claims.ID)
		context.Next()
	}
}