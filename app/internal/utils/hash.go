package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(input string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	return string(bytes), err
}
