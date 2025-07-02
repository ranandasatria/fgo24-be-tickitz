package utils

import (
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
