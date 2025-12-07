package database

import (
	"time"

	"ecommerce-gin/internal/models"

	"gorm.io/gorm"
)

// AdminListOrders returns orders with pagination, status filter, and date range
func AdminListOrders(limit, offset int, status string, from, to *time.Time) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := DB.Model(&models.Order{}).Preload("OrderItems.Product").Order("created_at desc")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if from != nil {
		query = query.Where("created_at >= ?", *from)
	}
	if to != nil {
		query = query.Where("created_at <= ?", *to)
	}

	if err := query.Count(&total).Limit(limit).Offset(offset).Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

// AdminGetOrderByID returns an order by id (admin)
func AdminGetOrderByID(orderID uint) (*models.Order, error) {
	var order models.Order
	if err := DB.Preload("OrderItems.Product").First(&order, orderID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// UpdateOrderStatus changes status inside given tx (if tx==nil it uses DB)
func UpdateOrderStatus(tx *gorm.DB, orderID uint, newStatus string) error {
	exec := tx
	if exec == nil {
		exec = DB
	}
	return exec.Model(&models.Order{}).Where("id = ?", orderID).Update("status", newStatus).Error
}

// RefundStockOnCancel decrements nothing here — instead restore stock when cancelling
func RestoreStockForOrder(tx *gorm.DB, orderID uint) error {
	// For each order item, add quantity back to product
	var items []models.OrderItem
	if err := tx.Where("order_id = ?", orderID).Find(&items).Error; err != nil {
		return err
	}
	for _, it := range items {
		if err := tx.Model(&models.Product{}).
			Where("id = ?", it.ProductID).
			Update("quantity", gorm.Expr("quantity + ?", it.Quantity)).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetOrderCountsByStatus simple stats
func GetOrderCountsByStatus() (map[string]int64, error) {
	type row struct {
		Status string
		Count  int64
	}
	var rows []row
	result := make(map[string]int64)
	if err := DB.Model(&models.Order{}).Select("status, count(*) as count").Group("status").Scan(&rows).Error; err != nil {
		return nil, err
	}
	for _, r := range rows {
		result[r.Status] = r.Count
	}
	return result, nil
}

// GetRevenueSummary between dates
func GetRevenueSummary(from, to *time.Time) (float64, error) {
	query := DB.Model(&models.Order{}).Select("COALESCE(SUM(total_amount),0) as total")
	if from != nil {
		query = query.Where("created_at >= ?", *from)
	}
	if to != nil {
		query = query.Where("created_at <= ?", *to)
	}
	var total float64
	if err := query.Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}
