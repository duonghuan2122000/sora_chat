package util

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func VerifyPassword(password string, passwordHashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHashed), []byte(password))
	return err == nil
}
