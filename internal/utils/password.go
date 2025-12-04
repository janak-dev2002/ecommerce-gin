package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes a plain-text password
func HashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}

// CheckPassword compares hashed password and plain text
func CheckPassword(hashed, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pw)) == nil
}
