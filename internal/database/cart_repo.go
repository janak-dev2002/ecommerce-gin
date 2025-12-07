package database

import (
	"ecommerce-gin/internal/models"
)

// Get all cart items for a user
func GetCartItems(userID uint) ([]models.CartItem, error) {
	var items []models.CartItem
	err := DB.Preload("Product").Where("user_id = ?", userID).Find(&items).Error
	return items, err
}

// Find specific cart item (user + product)
func GetCartItem(userID, productID uint) (*models.CartItem, error) {
	var item models.CartItem
	err := DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func GetCartItemByID(itemID uint) (*models.CartItem, error) {
	var item models.CartItem
	err := DB.First(&item, itemID).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Create item
func CreateCartItem(item *models.CartItem) error {
	return DB.Create(item).Error
}

// Update quantity
func UpdateCartItem(itemID uint, qty int) error {
	return DB.Model(&models.CartItem{}).Where("id = ?", itemID).Update("quantity", qty).Error
}

// Remove one item
func RemoveCartItem(itemID uint) error {
	return DB.Delete(&models.CartItem{}, itemID).Error
}

// Clear cart
func ClearCart(userID uint) error {
	return DB.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error
}
