package controllers

import (
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"example/APIForWorldWorkHub/database"
	"example/APIForWorldWorkHub/models"
	"example/APIForWorldWorkHub/utils"
)

func AddService(context *gin.Context) {
	var input struct {
		UserID      string    `json:"userId" validate:"required"`
		Date        time.Time `json:"date" validate:"required"`
		Pay         float64   `json:"pay" validate:"required,min=0"`
		Description string    `json:"description" validate:"required"`
	}

	if err := context.BindJSON(&input); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
		return
	}

	// Validate the service
	if err := utils.Validate.Struct(input); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "errors": err.Error()})
		return
	}

	// Verificar se o usuário existe
	var user models.User
	if err := database.DB.Where("id = ?", input.UserID).First(&user).Error; err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	newService := models.Service{
		UserID:      input.UserID,
		Date:        input.Date,
		Pay:         int64(input.Pay * 100),
		Description: input.Description,
	}

	if err := database.DB.Create(&newService).Error; err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create service"})
		return
	}

	context.IndentedJSON(http.StatusCreated, newService)
}

func GetServices(context *gin.Context) {
	var services []models.Service
	if err := database.DB.Find(&services).Error; err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve services"})
		return
	}

	context.IndentedJSON(http.StatusOK, services)
}