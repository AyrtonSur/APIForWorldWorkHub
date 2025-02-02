package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "example/APIForWorldWorkHub/database"
    "example/APIForWorldWorkHub/routes"
    "example/APIForWorldWorkHub/utils"
    "example/APIForWorldWorkHub/middlewares" // Importe o pacote
)

func main() {
    utils.InitValidator()
    database.InitialMigration()
    
    router := gin.Default()
    
    // Aplique o middleware CORS globalmente
    router.Use(middlewares.EnableCORS()) // <--- Adição do middleware
    
    routes.SetupRoutes(router)
    
    router.Run("localhost:9090")
    fmt.Println("Server Running in localhost:9090")
}
