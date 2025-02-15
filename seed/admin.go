package seed

import (
	"example/APIForWorldWorkHub/models"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) {
	if os.Getenv("GIN_MODE") != "dev" {
			fmt.Println("Seed admin: ambiente não é de desenvolvimento, seed ignorado.")
			return
	}

	fmt.Println("Executando seed de admin para ambiente de desenvolvimento...")

	var occupation models.Occupation
	if err := db.First(&occupation, "name = ?", "Cleaner").Error; err != nil {
		fmt.Println("Ocupação não encontrada.")
		return
	}

	var region models.Region
	if err := db.Where("abbreviation = ?", "LA").First(&region).Error; err != nil {
		fmt.Println("Região não encontrada.")
		return
	}

	var role models.Role
	if err := db.Where("name = ?", "Admin").First(&role).Error; err != nil {
		fmt.Println("Role Admin não encontrada")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Erro a gerar senha")
		return
	}

	admin := models.User{
		Firstname:      "Admin",
		Lastname:       "User",
		Email:          "admin@example.com",
		PasswordDigest: string(hashedPassword), // lembre-se de armazenar a senha já hasheada
		CPF:            nil,               // ou um valor válido se necessário
		RoleID:         "admin-role-id",   // substitua pelo ID real do role admin, se necessário
		OccupationID:   &occupation.ID,               // ou um valor válido se necessário
		Phone:          "0000000000000",
		ZipCode:        "12345",
		Education:      "Admin Education",
		RegionID:       region.ID, // substitua pelo ID real da região, se necessário
		City:           "Admin City",
		RefreshToken:   nil,               // ou um valor válido se necessário
			// Services e SpokenLanguages podem ser preenchidos conforme necessário
	}
		
	if err := db.FirstOrCreate(&admin, admin).Error; err != nil {
		fmt.Printf("Erro ao criar admin: %v\n", err)
		return
	}

	// Mensagem de sucesso (ou outras operações necessárias)
	fmt.Println("Seed admin executado com sucesso!")
}