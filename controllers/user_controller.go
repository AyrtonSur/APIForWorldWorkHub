package controllers

import (
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"example/APIForWorldWorkHub/database"
	"example/APIForWorldWorkHub/models"
	"example/APIForWorldWorkHub/utils"
)

type UserResponse struct {
	ID              string     `json:"id"`
	Firstname       string     `json:"firstname"`
	Lastname        string     `json:"lastname"`
	Email           string     `json:"email"`
	Role            string     `json:"role"`
	Occupation      string     `json:"occupation"`
	Phone           string     `json:"phone"`
	Education       string     `json:"education"`
	Region          string     `json:"region"`
	City            string     `json:"city"`
	ZipCode         string     `json:"zipcode"`
	Services        []models.Service  `json:"services"`
	SpokenLanguages []models.Language `json:"languages"`
}

func mapUserToResponse(user models.User) UserResponse {
	return UserResponse{
		ID:              user.ID,
		Firstname:       user.Firstname,
		Lastname:        user.Lastname,
		Email:           user.Email,
		Role:            user.Role,
		Occupation:      user.Occupation.Name,
		Phone:           user.Phone,
		Education:       user.Education,
		Region:          user.Region.Abbreviation,
		City:            user.City,
		ZipCode:         user.ZipCode,
		Services:        user.Services,
		SpokenLanguages: user.SpokenLanguages,
	}
}

func GetUsers(context *gin.Context) {
	var users []models.User
	if err := database.DB.Preload("Services").Preload("SpokenLanguages").Preload("Region").Preload("Occupation").Find(&users).Error; err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve users"})
		return
	}

	var userResponses []UserResponse
	for _, user := range users {
			userResponse := mapUserToResponse(user)
			userResponses = append(userResponses, userResponse)
	}

	context.IndentedJSON(http.StatusOK, userResponses)
}

func Register(context *gin.Context) {
	var newUser struct {
		Firstname      string  `json:"firstname" validate:"required"`
		Lastname       string  `json:"lastname" validate:"required"`
		Email          string  `json:"email" validate:"required,email"`
		Password       string  `json:"password" validate:"required,password"`
		CPF            *string `json:"CPF" validate:"omitempty,cpf"`
		Role           string  `json:"role" validate:"required"`
		OccupationName string  `json:"occupation" validate:"required"`
		Phone          string  `json:"phone" validate:"required,phone"`
		Education      string  `json:"education" validate:"required"`
		Region         string  `json:"region" validate:"required"`
		City           string  `json:"city" validate:"required"`
		ZipCode        string  `json:"zipcode" validate:"required,zipcode"`
	}

	if err := context.BindJSON(&newUser); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
		return
	}

	if err := utils.Validate.Struct(newUser); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "errors": err.Error()})
		return
	}

	var occupation models.Occupation
	if err := database.DB.First(&occupation, "name = ?", newUser.OccupationName).Error; err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Role Not Found", "errors": err.Error()})
		return
	}

	var region models.Region
	if err := database.DB.Where("name = ? OR abbreviation = ?", newUser.Region, newUser.Region).First(&region).Error; err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Region Not Found", "errors": err.Error()})
		return
	}
	
	// Hash the password before saving the user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
		return
	}

	// Cria o usuário com a Occupation associada e Region associada
	newUserModel := models.User{
		Firstname:      newUser.Firstname,
		Lastname:       newUser.Lastname,
		Email:          newUser.Email,
		PasswordDigest: string(hashedPassword),
		CPF:            newUser.CPF,
		Role:           newUser.Role,
		OccupationID:   &occupation.ID,
		Phone:          newUser.Phone,
		Education:      newUser.Education,
		RegionID:       region.ID,
		City:           newUser.City,
		ZipCode:        newUser.ZipCode,
	}

	if err := database.DB.Create(&newUserModel).Error; err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}
	
	    // Recarregar o usuário com os relacionamentos
	if err := database.DB.Preload("Services").Preload("SpokenLanguages").Preload("Region").Preload("Occupation").Where("id = ?", newUserModel.ID).First(&newUserModel).Error; err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to load user data"})
		return
	}
	
	
	userResponse := mapUserToResponse(newUserModel)
	context.IndentedJSON(http.StatusCreated, userResponse)
}

func GetUserById(id string) (*models.User, error) {
	var user models.User
	if err := database.DB.Preload("Services").Preload("SpokenLanguages").Preload("Region").Preload("Occupation").Where("id = ?", id).First(&user).Error; err != nil {
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

	userResponse := mapUserToResponse(*user)
	context.IndentedJSON(http.StatusOK, userResponse)
}

func Login(context *gin.Context) {
	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
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