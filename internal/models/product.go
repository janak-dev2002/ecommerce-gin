package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	Name        string  `json:"name"`
	Slug        string  `json:"slug" gorm:"type:varchar(255);uniqueIndex;not null"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
	IsActive    bool    `json:"is_active" gorm:"default:true"`
}
