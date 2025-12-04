package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"ecommerce-gin/internal/config"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateAccessToken creates a signed JWT with short expiry
func GenerateAccessToken(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(time.Minute * time.Duration(config.Cfg.AccessTokenMinutes)).Unix(),
		"iat":  time.Now().Unix(),
		"iss":  "ecommerce-gin",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Cfg.JWTSecret))
}

// GenerateRefreshToken creates a cryptographically random string (not JWT)
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
