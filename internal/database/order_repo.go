package database

import (
	"ecommerce-gin/internal/models"

	"gorm.io/gorm"
)

func CreateOrder(tx *gorm.DB, order *models.Order) error {
	return tx.Create(order).Error
}

func CreateOrderItems(tx *gorm.DB, items []models.OrderItem) error {
	return tx.Create(&items).Error
}

func GetOrdersForUser(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := DB.Preload("OrderItems.Product").Preload("OrderItems.Order").
		Where("user_id = ?", userID).Order("id DESC").
		Find(&orders).Error
	return orders, err
}

func GetOrderByID(userID, orderID uint) (*models.Order, error) {
	var order models.Order
	err := DB.Preload("OrderItems.Product").
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func DeductProductStock(tx *gorm.DB, productID uint, quantity int) error {
	return tx.Model(&models.Product{}).
		Where("id = ?", productID).
		Update("quantity", gorm.Expr("quantity - ?", quantity)).Error
}

func ClearCartInTransaction(tx *gorm.DB, userID uint) error {
	return tx.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error
}
