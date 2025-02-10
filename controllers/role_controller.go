package controllers

import (
	"example/APIForWorldWorkHub/database"
	"example/APIForWorldWorkHub/models"
	"example/APIForWorldWorkHub/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateRole(context *gin.Context) {
	var input struct {
		Name string `json:"name" validate:"required,min=2"`
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

	if err := database.DB.Create(&input).Error; err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create role"})
		return
	}

	context.IndentedJSON(http.StatusCreated, input)
}

func GetRoles(context *gin.Context) {
	var roles []models.Role
	if err := database.DB.Find(&roles).Error; err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve roles"})
		return
	}

	context.IndentedJSON(http.StatusOK, roles)
}

func DeleteRole(context *gin.Context) {
	name := context.Param("name")

	var role models.Role
	if err := database.DB.Where("Name = ?", name).First(&role).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Role not found"})
		return
	}

	if err := database.DB.Delete(&role).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete role"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}
