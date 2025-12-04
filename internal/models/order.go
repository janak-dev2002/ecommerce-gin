package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model

	UserID     uint    `json:"user_id"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"` // pending, confirmed, shipped, delivered

	User       User
	OrderItems []OrderItem
}
