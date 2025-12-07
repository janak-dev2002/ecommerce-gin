package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model

	UserID    uint `json:"user_id"`    // This is the place to link to User model, it automatically creates foreign key relation
	ProductID uint `json:"product_id"` // This is the place to link to User model, it automatically creates foreign key relation
	Quantity  int  `json:"quantity"`

	// User    User    `json:"user" gorm:"foreignKey:UserID"`  // This helps to preload user details in cart items, not create the foreign key
	Product Product `json:"product,omitzero" gorm:"foreignKey:ProductID"` // This is helps to preload product details in cart items  // not create the foreign key
}
