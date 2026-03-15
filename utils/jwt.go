package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "secretKey" // In reality it should be more secret and hard to guess

func GenerateToken(email string, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"userId": userId, // Saves it in float despite of type int
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string) (int64, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		// Check if token's Singing Method is same when it was created
		_, ok := t.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("Unexpected signing method!")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		err_msg := fmt.Sprintf("Could not parse token! %v", err)
		return 0, errors.New(err_msg)
	}

	if !parsedToken.Valid {
		return 0, errors.New("Invalid Token!")
	}

	// Way to access tokens info
	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("Invalid Token Claims!")
	}

	// email := claims["email"].(string)
	userId := int64(claims["userId"].(float64))

	return userId, nil
}