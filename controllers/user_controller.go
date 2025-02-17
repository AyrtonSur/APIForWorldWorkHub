package controllers

import (
	"errors"
	"example/APIForWorldWorkHub/database"
	"example/APIForWorldWorkHub/models"
	"example/APIForWorldWorkHub/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
)

type UserResponse struct {
	ID              string            `json:"id"`
	Firstname       string            `json:"firstname"`
	Lastname        string            `json:"lastname"`
	Email           string            `json:"email"`
	Role            string            `json:"role"`
	Occupation      string            `json:"occupation"`
	Phone           string            `json:"phone"`
	Education       string            `json:"education"`
	Region          string            `json:"region"`
	City            string            `json:"city"`
	ZipCode         string            `json:"zipcode"`
	Services        []models.Service  `json:"services"`
	SpokenLanguages []models.Language `json:"languages"`
}

func mapUserToResponse(user models.User) UserResponse {
	return UserResponse{
		ID:              user.ID,
		Firstname:       user.Firstname,
		Lastname:        user.Lastname,
		Email:           user.Email,
		Role:            user.Role.Name,
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
	if err := database.DB.Preload("Services").Preload("SpokenLanguages").Preload("Region").Preload("Occupation").Preload("Role").Find(&users).Error; err != nil {
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

	// Verificar se já existe um usuário com o mesmo email
	var existingUser models.User
	if err := database.DB.Where("email = ?", newUser.Email).First(&existingUser).Error; err == nil {
		context.IndentedJSON(http.StatusConflict, gin.H{"message": "Email already in use"})
		return
	}

	// Verificar se já existe um usuário com o mesmo CPF
	if newUser.CPF != nil {
		if err := database.DB.Where("cpf = ?", *newUser.CPF).First(&existingUser).Error; err == nil {
			context.IndentedJSON(http.StatusConflict, gin.H{"message": "CPF already in use"})
			return
		}
	}

	var occupation models.Occupation
	if err := database.DB.First(&occupation, "name = ?", newUser.OccupationName).Error; err != nil {
		context.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "Occupation Not Found", "errors": err.Error()})
		return
	}

	var region models.Region
	if err := database.DB.Where("name = ? OR abbreviation = ?", newUser.Region, newUser.Region).First(&region).Error; err != nil {
		context.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "Region Not Found", "errors": err.Error()})
		return
	}

	var role models.Role
	if err := database.DB.Where("name = ?", "User").First(&role).Error; err != nil {
		context.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "Role Not Found", "errors": err.Error()})
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
		RoleID:         role.ID,
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
	if err := database.DB.Preload("Services").Preload("SpokenLanguages").Preload("Region").Preload("Occupation").Preload("Role").Where("id = ?", newUserModel.ID).First(&newUserModel).Error; err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to load user data"})
		return
	}

	userResponse := mapUserToResponse(newUserModel)
	context.IndentedJSON(http.StatusCreated, userResponse)
}

func GetUserByID(id string) (*models.User, error) {
	var user models.User
	if err := database.DB.Preload("Services").Preload("SpokenLanguages").Preload("Region").Preload("Occupation").Preload("Role").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func GetUser(context *gin.Context) {
	id := context.Param("id")
	user, err := GetUserByID(id)

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
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(input.Password)); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Credentials"})
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate access token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate refresh token"})
		return
	}

	user.RefreshToken = &refreshToken
	if err := database.DB.Save(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save refresh token"})
		return
	}

	if err := database.DB.Preload("Services").Preload("SpokenLanguages").Preload("Region").Preload("Occupation").Preload("Role").Where("id = ?", user.ID).First(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load user data"})
		return
	}

	fmt.Println(user)

	secure := os.Getenv("GIN_MODE") == "release"

	context.SetCookie("refresh_token", refreshToken, 7*24*60*60, "/", "", secure, true)

	userResponse := mapUserToResponse(user)
	context.JSON(http.StatusOK, gin.H{"user": userResponse, "access_token": accessToken})
}

func RefreshToken(context *gin.Context) {
	refreshToken, err := context.Cookie("refresh_token")
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "No refresh token provided"})
		return
	}

	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid refresh token"})
		return
	}

	var user models.User
	if err := database.DB.Where("id = ? AND refresh_token = ?", claims.ID, refreshToken).First(&user).Error; err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid refresh token"})
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate access token"})
		return
	}

	newRefreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate refresh token"})
		return
	}

	user.RefreshToken = &newRefreshToken
	if err := database.DB.Save(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save refresh token"})
		return
	}

	secure := os.Getenv("GIN_MODE") == "release"

	context.SetCookie("refresh_token", newRefreshToken, 7*24*60*60, "/", "", secure, true)

	context.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

func Logout(context *gin.Context) {
	userID := context.GetString("userID")

	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid user"})
		return
	}

	user.RefreshToken = nil
	if err := database.DB.Save(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not log out"})
		return
	}

	secure := os.Getenv("GIN_MODE") == "release"

	// Remover o refresh token do cookie
	context.SetCookie("refresh_token", "", -1, "/", "", secure, true)

	context.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

type UpdateUserInput struct {
	Firstname      *string `json:"firstname"`
	Lastname       *string `json:"lastname"`
	CPF            *string `json:"CPF" validate:"omitempty,cpf"`
	Role           *string `json:"role"`
	OccupationName *string `json:"occupation"`
	Phone          *string `json:"phone" validate:"omitempty,phone"`
	Education      *string `json:"education"`
	Region         *string `json:"region"`
	City           *string `json:"city"`
	ZipCode        *string `json:"zipcode" validate:"omitempty,zipcode"`
}

func UpdateUser(context *gin.Context) {
	id := context.Param("id")
	var input UpdateUserInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "errors": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	if input.Firstname != nil {
		user.Firstname = *input.Firstname
	}

	if input.Lastname != nil {
		user.Lastname = *input.Lastname
	}

	if input.CPF != nil {
		var existingUser models.User
		if err := database.DB.Where("cpf = ? AND id != ?", *input.CPF, id).First(&existingUser).Error; err == nil {
			context.JSON(http.StatusConflict, gin.H{"message": "CPF already in use"})
			return
		}
	}

	if input.Role != nil {
		var role models.Role
		if err := database.DB.First(&role, "name = ?", *input.Role).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Role Not Found"})
			return
		}
	}

	if input.OccupationName != nil {
		var occupation models.Occupation
		if err := database.DB.First(&occupation, "name = ?", *input.OccupationName).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Occupation Not Found"})
			return
		}

		user.OccupationID = &occupation.ID
	}

	if input.Phone != nil {
		user.Phone = *input.Phone
	}

	if input.Education != nil {
		user.Education = *input.Education
	}

	if input.Region != nil {
		var region models.Region
		if err := database.DB.Where("name = ? OR abbreviation = ?", *input.Region, *input.Region).First(&region).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Region Not Found"})
			return
		}

		user.RegionID = region.ID
	}

	if input.City != nil {
		user.City = *input.City
	}

	if input.ZipCode != nil {
		user.ZipCode = *input.ZipCode
	}

	if err := database.DB.Save(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user"})
		return
	}

	// Recarregar o usuário com os relacionamentos
	if err := database.DB.Preload("Services").Preload("SpokenLanguages").Preload("Region").Preload("Occupation").Preload("Role").Where("id = ?", user.ID).First(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load user data"})
		return
	}

	userResponse := mapUserToResponse(user)
	context.IndentedJSON(http.StatusOK, userResponse)
}

func DeleteUser(context *gin.Context) {
	id := context.Param("id")

	var user models.User
	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete user"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func CreateUser(context *gin.Context) {
	var newUser struct {
		Firstname      string  `json:"firstname" validate:"required"`
		Lastname       string  `json:"lastname" validate:"required"`
		Email          string  `json:"email" validate:"required,email"`
		Password       string  `json:"password" validate:"required,password"`
		CPF            *string `json:"CPF" validate:"omitempty,cpf"`
		OccupationName string  `json:"occupation" validate:"required"`
		Phone          string  `json:"phone" validate:"required,phone"`
		Role           string  `json:"role" validate:"required"`
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

	// Verificar se já existe um usuário com o mesmo email
	var existingUser models.User
	if err := database.DB.Where("email = ?", newUser.Email).First(&existingUser).Error; err == nil {
		context.IndentedJSON(http.StatusConflict, gin.H{"message": "Email already in use"})
		return
	}

	// Verificar se já existe um usuário com o mesmo CPF
	if newUser.CPF != nil {
		if err := database.DB.Where("cpf = ?", *newUser.CPF).First(&existingUser).Error; err == nil {
			context.IndentedJSON(http.StatusConflict, gin.H{"message": "CPF already in use"})
			return
		}
	}

	var occupation models.Occupation
	if err := database.DB.First(&occupation, "name = ?", newUser.OccupationName).Error; err != nil {
		context.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "Occupation Not Found", "errors": err.Error()})
		return
	}

	var region models.Region
	if err := database.DB.Where("name = ? OR abbreviation = ?", newUser.Region, newUser.Region).First(&region).Error; err != nil {
		context.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "Region Not Found", "errors": err.Error()})
		return
	}

	var role models.Role
	if err := database.DB.Where("name = ?", newUser.Role).First(&role).Error; err != nil {
		context.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "Role Not Found", "errors": err.Error()})
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
		RoleID:         role.ID,
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
	if err := database.DB.Preload("Services").Preload("SpokenLanguages").Preload("Region").Preload("Occupation").Preload("Role").Where("id = ?", newUserModel.ID).First(&newUserModel).Error; err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to load user data"})
		return
	}

	userResponse := mapUserToResponse(newUserModel)
	context.IndentedJSON(http.StatusCreated, userResponse)
}

func currentUser(c *gin.Context) (*models.User, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return nil, errors.New("user not found in context")
	}

	var user models.User
	if err := database.DB.Preload("Services").Preload("SpokenLanguages").Preload("Region").Preload("Occupation").Preload("Role").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func GetCurrentUser(context *gin.Context) {
	user, err := currentUser(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	userResponse := mapUserToResponse(*user)
	context.JSON(http.StatusOK, userResponse)
}

func CheckEmailExists(context *gin.Context) {
	var input struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err := context.BindJSON(&input); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "errors": err.Error()})
		return
	}

	var existingUser models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		context.IndentedJSON(http.StatusOK, gin.H{"exists": true})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{"exists": false})
}
