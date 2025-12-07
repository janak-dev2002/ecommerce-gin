package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model

	UserID     uint    `json:"user_id"` // This is the place to link to User model, it automatically creates foreign key relation
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"` // pending, confirmed, shipped, delivered

	// User       User `json:"user"` // optional: to preload user details
	OrderItems []OrderItem
}
