package models

import "gorm.io/gorm"

type PaymentIntent struct {
	gorm.Model

	OrderID uint    `json:"order_id"`
	Amount  float64 `json:"amount"`

	Status string `json:"status"` // pending, paid, failed

	// Simulated gateway reference
	GatewayRef string `json:"gateway_ref"`
}
