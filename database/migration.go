package database

import (
	"fmt"
	"os"
	"path/filepath"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"example/APIForWorldWorkHub/models"
	"example/APIForWorldWorkHub/data"
)

var DB *gorm.DB

func InitialMigration() {
	dataPath := "data"
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		os.Mkdir(dataPath, os.ModePerm)
	}

	dbPath := filepath.Join(dataPath, "test.db")

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}

	DB.AutoMigrate(&models.User{}, &models.Service{}, &models.Language{}, &models.Occupation{}, &models.Region{}, &models.Role{})

	for _, occupation := range data.Occupations {
		DB.FirstOrCreate(&occupation, models.Occupation{Name: occupation.Name})
}

	// Adicionar estados dos EUA
	for _, state := range data.States {
		DB.FirstOrCreate(&state, models.Region{Name: state.Name, Abbreviation: state.Abbreviation})
	}

	for _, role := range data.Roles {
		DB.FirstOrCreate(&role, models.Role{Name: role.Name})
	}
}