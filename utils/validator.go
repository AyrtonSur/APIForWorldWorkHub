package utils

import (
	"regexp"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()
	Validate.RegisterValidation("zipcode", validateZipCode)
	Validate.RegisterValidation("password", passwordValidator)
}

func validateZipCode(fl validator.FieldLevel) bool {
	zipCode := fl.Field().String()
	// Regex para validar ZipCode americano (XXXXX ou XXXXX-XXXX)
	regex := regexp.MustCompile(`^\d{5}(?:-\d{4})?$`)
	if !regex.MatchString(zipCode) {
		return false
	}

	// Verifica se o ZipCode realmente existe
	exists, err := ValidateZipCode(zipCode)
	if err != nil || !exists {
		return false
	}

	return true
}

func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return len(password) >= 4
}