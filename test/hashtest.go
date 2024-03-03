package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) []byte {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return []byte("")
	}
	return hash
}

func validatePassword(hash []byte, password string) bool {
	return bcrypt.CompareHashAndPassword(hash, []byte(password)) == nil
}

func main() {
	password := "password123"
	hash := hashPassword(password)
	fmt.Print(hash)
	fmt.Print(validatePassword(hash, password))
}
