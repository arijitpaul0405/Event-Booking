package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	byteHashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(byteHashPassword), err
}

func ValidatePassword(hashPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}