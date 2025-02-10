package main

import (
	"example/APIForWorldWorkHub/database"
	"example/APIForWorldWorkHub/routes"
	"example/APIForWorldWorkHub/utils"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
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
	if err := router.Run("localhost:9090"); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
	fmt.Println("Server Running in localhost:9090")
}
