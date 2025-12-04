package database

import (
	"errors"

	"ecommerce-gin/internal/models"

	"gorm.io/gorm"
)

/*
Create Product
*/
func CreateProduct(p *models.Product) error {
	return DB.Create(p).Error
}

/*
Get Product by ID
*/
func GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := DB.Preload("Category").First(&product, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &product, err
}

/*
Update Product
*/
func UpdateProduct(id uint, data map[string]interface{}) error {
	return DB.Model(&models.Product{}).
		Where("id = ?", id).
		Updates(data).Error
}

/*
Delete Product
*/
func DeleteProduct(id uint) error {
	return DB.Delete(&models.Product{}, id).Error
}

/*
List Products with Filters
Filters:
  - category_id
  - min_price
  - max_price
  - search (title match)
  - sort = price_asc, price_desc, newest
  - pagination: page, limit
*/
func ListProducts(filters map[string]interface{}) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := DB.Model(&models.Product{}).Preload("Category")

	// Search by name
	if s, ok := filters["search"]; ok {
		query = query.Where("title LIKE ?", "%"+s.(string)+"%")
	}

	// Filter by category
	if c, ok := filters["category_id"]; ok {
		query = query.Where("category_id = ?", c)
	}

	// Price range
	if min, ok := filters["min_price"]; ok {
		query = query.Where("final_price >= ?", min)
	}
	if max, ok := filters["max_price"]; ok {
		query = query.Where("final_price <= ?", max)
	}

	// Sorting
	if sort, ok := filters["sort"]; ok {
		switch sort {
		case "price_asc":
			query = query.Order("final_price ASC")
		case "price_desc":
			query = query.Order("final_price DESC")
		case "newest":
			query = query.Order("created_at DESC")
		}
	}

	// Pagination
	page := filters["page"].(int)
	limit := filters["limit"].(int)

	offset := (page - 1) * limit

	// Count total
	query.Count(&total)

	// Fetch page
	err := query.Limit(limit).Offset(offset).Find(&products).Error

	return products, total, err
}
