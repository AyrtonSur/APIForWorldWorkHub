package seed

import (
	"example/APIForWorldWorkHub/models"
	"gorm.io/gorm"
)

var Occupations = []models.Occupation{
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

func InitializeOccupations(db *gorm.DB) {
	for _, occupation := range Occupations {
		db.FirstOrCreate(&occupation, models.Occupation{Name: occupation.Name})
	}
}