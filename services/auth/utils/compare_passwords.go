package utils

import "golang.org/x/crypto/bcrypt"

func ComparePasswords(storagePassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(storagePassword), []byte(password))
}
