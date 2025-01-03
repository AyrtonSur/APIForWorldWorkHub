package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"example/APIForWorldWorkHub/database"
	"example/APIForWorldWorkHub/routes"
	"example/APIForWorldWorkHub/utils"
)

func main() {
	utils.InitValidator()
	database.InitialMigration()
	router := gin.Default()
	routes.SetupRoutes(router)
	router.Run("localhost:9090")
	fmt.Println("Server Running in localhost:9090")
}