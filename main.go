package main

import (
	"regexp"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"example/APIForWorldWorkHub/utils"
	"example/APIForWorldWorkHub/database"
	"example/APIForWorldWorkHub/routes"
)

var validate *validator.Validate

func main() {
	fmt.Println("Go ORM Tutorial")
	validate = validator.New()
	validate.RegisterValidation("zipcode", validateZipCode)
	validate.RegisterValidation("password", passwordValidator) // Registrar a validação personalizada
	database.InitialMigration()
	router := gin.Default()
	routes.SetupRoutes(router, validate)
	router.Run("localhost:9090")
}

func validateZipCode(fl validator.FieldLevel) bool {
	zipCode := fl.Field().String()

	regex := regexp.MustCompile(`^\d{5}(?:-\d{4})?$`)
	if !regex.MatchString(zipCode) {
		return false
	}

	// Verifica se o ZipCode realmente existe
	exists, err := utils.ValidateZipCode(zipCode)
	if err != nil || !exists {
		return false
	}

	return true
}

func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return len(password) >= 4
}