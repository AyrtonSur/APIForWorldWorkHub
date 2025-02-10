package utils

import (
	"github.com/go-playground/validator/v10"
	"log"
	"regexp"
)

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()
	if err := Validate.RegisterValidation("zipcode", validateZipCode); err != nil {
		log.Fatalf("Erro ao registrar validação: %v", err)
	}
	if err := Validate.RegisterValidation("password", passwordValidator); err != nil {
		log.Fatalf("Erro ao registrar validação: %v", err)
	}
	if err := Validate.RegisterValidation("phone", validatePhone); err != nil {
		log.Fatalf("Erro ao registrar validação: %v", err)
	}
	if err := Validate.RegisterValidation("cpf", validateCPF); err != nil {
		log.Fatalf("Erro ao registrar validação: %v", err)
	}
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
	field := fl.Field().Interface()

	var cpf string

	switch v := field.(type) {
	case *string: // Caso o campo seja um ponteiro para string
		if v == nil {
			return true // Válido se for nil (omitempty)
		}
		cpf = *v
	case string: // Caso o campo seja uma string
		cpf = v
	default:
		return false // Tipo inválido
	}

	// Considera válido se vazio devido ao omitempty
	if cpf == "" {
		return true
	}

	// Verifica se o CPF tem exatamente 11 dígitos
	return len(cpf) == 11
}
