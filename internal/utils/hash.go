package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashToken stores refresh tokens securely (SHA-256)
func HashToken(token string) (string, error) {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:]), nil
}

func CheckTokenHash(hashed, token string) bool {
	hash := sha256.Sum256([]byte(token))
	return hashed == hex.EncodeToString(hash[:])
}
