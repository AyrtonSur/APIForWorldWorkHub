package main

import (
	"fmt"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"example/APIForWorldWorkHub/database"
	"example/APIForWorldWorkHub/routes"
	"example/APIForWorldWorkHub/utils"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Definir a variável de ambiente GIN_MODE se não estiver definida
	if os.Getenv("GIN_MODE") == "" {
		os.Setenv("GIN_MODE", "debug") // Altere para "release" em produção
	}

	utils.InitValidator()
	database.InitialMigration()
	router := gin.Default()

	// Configurar o middleware CORS
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // Altere para os domínios permitidos
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	routes.SetupRoutes(router)
	router.Run("localhost:9090")
	fmt.Println("Server Running in localhost:9090")
}