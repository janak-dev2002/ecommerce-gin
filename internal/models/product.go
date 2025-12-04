package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name        string    `json:"name" gorm:"type:varchar(255);uniqueIndex"`
	Description string    `json:"description"`
	Products    []Product `json:"products"`
}

type Product struct {
	gorm.Model
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Discount    float64 `json:"discount"`    // 0–100 (%)
	FinalPrice  float64 `json:"final_price"` // computed: Price * (1 - Discount/100)
	SKU         string  `json:"sku" gorm:"type:varchar(255);uniqueIndex"`
	Stock       int     `json:"stock"` // inventory
	ImageURL    string  `json:"image_url"`

	CategoryID *uint    `json:"category_id"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID"`

	// audit fields
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Product) BeforeSave(tx *gorm.DB) (err error) {
	p.FinalPrice = p.Price * (1 - p.Discount/100)
	return
}
