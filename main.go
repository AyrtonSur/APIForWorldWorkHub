package main

import (
	"errors"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             string     `gorm:"type:uuid;primary_key"`
	Firstname      string     `json:"firstname"`
	Lastname       string     `json:"lastname"`
	Email          string     `json:"email"`
	PasswordDigest string     `json:"password_digest"`
	Role           string     `json:"role"` // Adicionado campo Role
	Contact        string     `json:"contact"`
	Region         *string    `json:"region"`
	ServiceDesc    *string		`json:"description"`
	Services       []Service  `gorm:"foreignkey:UserID"` // Relacionamento um-para-muitos
}

type Service struct {
	ID     string `gorm:"type:uuid;primary_key"`
	UserID string `gorm:"type:uuid"` // Campo para associar ao usuário
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.New().String()
	return scope.SetColumn("ID", uuid)
}

func (service *Service) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.New().String()
	return scope.SetColumn("ID", uuid)
}

func initialMigration() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&User{}, &Service{})
}

func getUsers(context *gin.Context) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to connect database"})
		return
	}
	defer db.Close()

	var users []User
	if err := db.Find(&users).Error; err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve users"})
		return
	}

	context.IndentedJSON(http.StatusOK, users)
}

func addUser(context *gin.Context) {
	fmt.Println("New User Endpoint Hit")

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to connect database"})
		return
	}
	defer db.Close()

	var newUser User
	if err := context.BindJSON(&newUser); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
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

func getUserById(id string) (*User, error) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var user User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, errors.New("User not found")
	}

	return &user, nil
}

func getUser(context *gin.Context) {
	id := context.Param("id")
	user, err := getUserById(id)
	
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	context.IndentedJSON(http.StatusOK, user)
}

func handleRequests() {
	router := gin.Default()
	router.GET("/users", getUsers)
	router.POST("/users", addUser)
	router.GET("/users/:id", getUser)
	router.Run("localhost:9090")
}

func main() {
	fmt.Println("Go ORM Tutorial")
	initialMigration()

	handleRequests()

}