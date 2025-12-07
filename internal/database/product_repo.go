package database

import (
	"ecommerce-gin/internal/models"
)

func CreateProduct(p *models.Product) error {
	return DB.Create(p).Error
}

func UpdateProduct(id uint, data map[string]any) error {
	return DB.Model(&models.Product{}).Where("id = ?", id).Updates(data).Error
}

func DeleteProduct(id uint) error {
	return DB.Delete(&models.Product{}, id).Error
}

func GetProductByID(id uint) (*models.Product, error) {
	var p models.Product
	err := DB.First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func GetProductBySlug(slug string) (*models.Product, error) {
	var p models.Product
	err := DB.Where("slug = ?", slug).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func ListProducts(limit, offset int, search, category string) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := DB.Model(&models.Product{}).Where("is_active = ?", true)

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	query.Count(&total)

	err := query.Limit(limit).Offset(offset).Find(&products).Error
	return products, total, err
}
