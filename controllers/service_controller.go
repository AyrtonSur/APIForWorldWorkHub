package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"example/APIForWorldWorkHub/database"
	"example/APIForWorldWorkHub/models"
	"github.com/go-playground/validator/v10"
)

func AddService(context *gin.Context, validate *validator.Validate) {
	var newService models.Service
	if err := context.BindJSON(&newService); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
		return
	}

	// Validate the service
	if err := validate.Struct(newService); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "errors": err.Error()})
		return
	}

	// Verificar se o usuário existe
	var user models.User
	if err := database.DB.Where("id = ?", newService.UserID).First(&user).Error; err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	if err := database.DB.Create(&newService).Error; err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create service"})
		return
	}

	context.IndentedJSON(http.StatusCreated, newService)
}