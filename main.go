package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"example/APIForWorldWorkHub/database"
	"example/APIForWorldWorkHub/routes"
)

var validate *validator.Validate

func main() {
	fmt.Println("Go ORM Tutorial")
	validate = validator.New()
	validate.RegisterValidation("password", passwordValidator) // Registrar a validação personalizada
	database.InitialMigration()
	router := gin.Default()
	routes.SetupRoutes(router, validate)
	router.Run("localhost:9090")
}

func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return len(password) >= 4
}