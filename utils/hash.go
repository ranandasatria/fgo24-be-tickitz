package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashString(s string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash string: %v", err)
	}
	return string(hashed), nil
}

func CompareHash(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
