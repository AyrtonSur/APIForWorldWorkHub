package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"example/APIForWorldWorkHub/models"
)

var DB *gorm.DB

func InitialMigration() {
	var err error
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}

	DB.AutoMigrate(&models.User{}, &models.Service{}, &models.Language{}, &models.Occupation{})
}