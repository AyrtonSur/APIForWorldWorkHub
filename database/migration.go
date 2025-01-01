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

	occupations := []models.Occupation{
		{Name: "Cleaner"},
		{Name: "Lawyer"},
		{Name: "Doctor"},
		{Name: "Engineer"},
		{Name: "Teacher"},
		{Name: "Nurse"},
		{Name: "Accountant"},
		{Name: "Architect"},
		{Name: "Chef"},
		{Name: "Dentist"},
		{Name: "Electrician"},
		{Name: "Mechanic"},
		{Name: "Pharmacist"},
		{Name: "Plumber"},
		{Name: "Software Developer"},
		{Name: "Graphic Designer"},
		{Name: "Journalist"},
		{Name: "Photographer"},
		{Name: "Pilot"},
		{Name: "Police Officer"},
		{Name: "Scientist"},
		{Name: "Veterinarian"},
		{Name: "Writer"},
		{Name: "Musician"},
		{Name: "Artist"},
		{Name: "Therapist"},
		{Name: "Translator"},
		{Name: "Web Developer"},
		{Name: "Data Analyst"},
		{Name: "Marketing Specialist"},
	}

	for _, occupation := range occupations {
		DB.FirstOrCreate(&occupation, models.Occupation{Name: occupation.Name})
	}
}