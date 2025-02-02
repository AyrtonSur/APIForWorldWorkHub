package middlewares

import "github.com/gin-gonic/gin"

// Nova versão adaptada para Gin
func EnableCORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Configurar cabeçalhos CORS
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Lidar com pré-voo OPTIONS
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(200)
            return
        }

        // Continuar o fluxo de requisição
        c.Next()
    }
}
