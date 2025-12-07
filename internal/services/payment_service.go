package services

import (
	"fmt"

	"ecommerce-gin/internal/database"
	"ecommerce-gin/internal/models"
)

func CreatePaymentIntent(orderID uint, amount float64) (*models.PaymentIntent, error) {
	intent := &models.PaymentIntent{
		OrderID: orderID,
		Amount:  amount,
		Status:  "pending",
		GatewayRef: fmt.Sprintf("PAY-%d", orderID),
	}

	if err := database.CreatePaymentIntent(intent); err != nil {
		return nil, err
	}
	return intent, nil
}
