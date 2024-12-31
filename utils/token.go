package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtKey = []byte("sua-chave-secreta") // Mantenha esta chave privada

func GenerateJWT(email string) (string, error) {
	claims := &jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}