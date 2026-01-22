package utils

import "golang.org/x/crypto/bcrypt"

func VerfiyPassword(hashedPassword, candidatePassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))

	return err == nil
}
