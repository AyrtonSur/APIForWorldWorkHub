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
	Validate.RegisterValidation("phone", validatePhone)
	Validate.RegisterValidation("cpf", validateCPF)
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
	return len(password) >= 8
}

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return len(phone) == 11 || len(phone) == 13
}

func validateCPF(fl validator.FieldLevel) bool {
	cpfPtr := fl.Field().Interface().(*string)
	if cpfPtr == nil {
		return true // CPF é opcional, então é válido se estiver vazio
	}
	// Verifica se o CPF tem exatamente 11 dígitos
	return len(*cpfPtr) == 11
}