package database

import (
	"errors"
	"time"

	"ecommerce-gin/internal/models"
	"ecommerce-gin/internal/utils"

	"gorm.io/gorm"
)

// CreateUser saves a new user (expects hashed password)
func CreateUser(user *models.User) error {
	return DB.Create(user).Error
}

// GetUserByEmail returns user pointer or nil (if no user)
func GetUserByEmail(email string) (*models.User, error) {
	var u models.User
	if err := DB.Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// GetUserByID returns user by id
func GetUserByID(id uint) (*models.User, error) {
	var u models.User
	if err := DB.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// SaveRefreshToken stores hashed refresh token + expiry for a user
func SaveRefreshToken(userID uint, hashedToken string, expiry time.Time) error {
	return DB.Model(&models.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"refresh_token": hashedToken,
			"token_expiry":  expiry,
		}).Error
}

// ClearRefreshToken clears refresh token fields for user
func ClearRefreshToken(userID uint) error {
	return DB.Model(&models.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"refresh_token": "",
			"token_expiry":  nil,
		}).Error
}

// FindUserByActiveRefreshToken tries to find user(s) with token_expiry > now and checks hashed tokens.
// This is less efficient for huge user base but works for MVP / single-token-per-user model.
func FindUserByActiveRefreshToken(refreshToken string) (*models.User, error) {
	var users []models.User
	if err := DB.Where("token_expiry > ?", time.Now()).Find(&users).Error; err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.RefreshToken != "" {
			// compare hash
			if utils.CheckTokenHash(u.RefreshToken, refreshToken) {
				// must return pointer to record found in DB
				user := u
				return &user, nil
			}
		}
	}
	return nil, nil
}
