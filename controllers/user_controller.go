package controllers

import (
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"example/APIForWorldWorkHub/models"
)

func GetUsers(context *gin.Context) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to connect database"})
		return
	}

	var users []models.User
	if err := db.Preload("Services").Preload("SpokenLanguages").Find(&users).Error; err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve users"})
		return
	}

	context.IndentedJSON(http.StatusOK, users)
}

func AddUser(context *gin.Context, validate *validator.Validate) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to connect database"})
		return
	}

	var newUser models.User
	if err := context.BindJSON(&newUser); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
		return
	}

	// Validate the user
	if err := validate.Struct(newUser); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "errors": err.Error()})
		return
	}

	// Hash the password before saving the user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.PasswordDigest), bcrypt.DefaultCost)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
		return
	}
	newUser.PasswordDigest = string(hashedPassword)

	if err := db.Create(&newUser).Error; err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}

	context.IndentedJSON(http.StatusCreated, newUser)
}

func GetUserById(id string) (*models.User, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := db.Preload("Services").Preload("SpokenLanguages").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, errors.New("User not found")
	}

	return &user, nil
}

func GetUser(context *gin.Context) {
	id := context.Param("id")
	user, err := GetUserById(id)
	
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	context.IndentedJSON(http.StatusOK, user)
}