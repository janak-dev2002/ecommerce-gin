package database

import (
	"ecommerce-gin/internal/models"
)

func CreatePaymentIntent(intent *models.PaymentIntent) error {
	return DB.Create(intent).Error
}

func UpdatePaymentIntentStatus(id uint, status string) error {
	return DB.Model(&models.PaymentIntent{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func GetPaymentIntentByID(id uint) (*models.PaymentIntent, error) {
	var p models.PaymentIntent
	if err := DB.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}
