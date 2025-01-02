package controllers

import (
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"example/APIForWorldWorkHub/database"
	"example/APIForWorldWorkHub/models"
	"example/APIForWorldWorkHub/utils"
)

func GetUsers(context *gin.Context) {
	var users []models.User
	if err := database.DB.Preload("Services").Preload("SpokenLanguages").Find(&users).Error; err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve users"})
		return
	}

	context.IndentedJSON(http.StatusOK, users)
}

func Register(context *gin.Context, validate *validator.Validate) {
	var newUser struct {
		Firstname      string  `json:"firstname" validate:"required"`
		Lastname       string  `json:"lastname" validate:"required"`
		Email          string  `json:"email" validate:"required,email"`
		Password       string  `json:"password" validate:"required,password"`
		CPF            *string `json:"CPF"`
		Role           string  `json:"role" validate:"required"`
		OccupationName string  `json:"occupation" validate:"required"`
		Phone          string  `json:"phone"`
		Education      string  `json:"education"`
		Region         *string `json:"region"`
	}

	if err := context.BindJSON(&newUser); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
		return
	}


	// Validate the user
	if err := validate.Struct(newUser); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "errors": err.Error()})
		return
	}

	var occupation models.Occupation
	if err := database.DB.First(&occupation, "name = ?", newUser.OccupationName).Error; err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Role Not Found", "errors": err.Error()})
		return
	}
	
	// Hash the password before saving the user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
		return
	}
	newUser.Password = string(hashedPassword)

	if err := database.DB.Create(&newUser).Error; err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}

	context.IndentedJSON(http.StatusCreated, newUser)
}

func GetUserById(id string) (*models.User, error) {
	var user models.User
	if err := database.DB.Preload("Services").Preload("SpokenLanguages").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
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

func Login(context *gin.Context, validate *validator.Validate) {
	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := validate.Struct(input); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "errors": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(input.Password)); err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Credentials"})
		return
	}

	token, err := utils.GenerateJWT(user.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"token": token})
}