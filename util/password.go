package util

import (
	"fmt"

	bycript "golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	hashedPassword, err := bycript.GenerateFromPassword([]byte(password), bycript.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}

	return string(hashedPassword), nil
}

func CheckPassword(password, hashedPassword string) error {

	return bycript.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

}
