package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name     string `json:"name"`
	Email    string `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Password string `json:"-"`                            // don't expose
	Role     string `json:"role" gorm:"default:customer"` // "admin" or "customer"

	RefreshToken string     `json:"-" gorm:"type:text"` // hashed
	TokenExpiry  *time.Time `json:"-"`                  // refresh token expiry

	CartItems []CartItem
	Orders    []Order
}
