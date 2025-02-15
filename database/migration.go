package database

import (
	"example/APIForWorldWorkHub/models"
	"example/APIForWorldWorkHub/seed"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
)

var DB *gorm.DB

func InitialMigration() {
	dataPath := "data"
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		err := os.Mkdir(dataPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Erro ao criar diretório: %v", err)
		}
	}

	dbPath := filepath.Join(dataPath, "test.db")

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}

	if err := DB.AutoMigrate(
		&models.User{},
		&models.Service{},
		&models.Language{},
		&models.Occupation{},
		&models.Region{},
		&models.Role{},
		&models.Permission{},
	); err != nil {
		log.Fatalf("Erro ao rodar migrações: %v", err)
	}

	seed.InitializeOccupations(DB)
	seed.InitializeStates(DB)
	seed.InitializeRoles(DB)
	if err := seed.InitializePermissions(DB); err != nil {
		log.Fatalf("Erro ao inicializar permissões: %v", err)
	}
	seed.SeedAdmin(DB)
}
