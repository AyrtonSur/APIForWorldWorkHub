package database

import (
	"fmt"
	"os"
	"path/filepath"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"example/APIForWorldWorkHub/models"
	"example/APIForWorldWorkHub/seed"
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

	DB.AutoMigrate(
		&models.User{},
		&models.Service{},
		&models.Language{},
		&models.Occupation{},
		&models.Region{},
		&models.Role{},
		&models.Permission{},
	)

	seed.InitializeOccupations(DB)
	seed.InitializeStates(DB)
	seed.InitializeRoles(DB)
	seed.InitializePermissions(DB)
}