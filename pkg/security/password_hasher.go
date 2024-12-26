package security

import "golang.org/x/crypto/bcrypt"

// Hash user password with bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return *new(string), err
	}
	password = string(hashedPassword)
	return password, err
}

// Compare user password and payload
func ComparePasswords(hashedPassword string, password string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false, err
	}
	return true, nil
}
