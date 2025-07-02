package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func GenerateJWT(purpose string, userId int, expiry time.Duration) (string, error) {
	godotenv.Load()
	secret := os.Getenv("APP_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":  userId,
		"purpose": purpose,
		"exp":     time.Now().Add(expiry).Unix(),
		"iat":     time.Now().Unix(),
	})

	return token.SignedString([]byte(secret))
}

func ParseJWT(tokenStr string) (jwt.MapClaims, error) {
	godotenv.Load()
	secret := os.Getenv("APP_SECRET")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
