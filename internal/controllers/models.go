package controllers

import "time"

// Swagger Model Definitions

// SignupInput represents the signup request body
type SignupInput struct {
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
	Role     string `json:"role" example:"customer"`
}

// LoginInput represents the login request body
type LoginInput struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn   int    `json:"expires_in" example:"3600"`
	User        UserResponse `json:"user"`
}

// UserResponse represents user data in responses
type UserResponse struct {
	ID    uint   `json:"id" example:"1"`
	Email string `json:"email" example:"john@example.com"`
	Role  string `json:"role" example:"customer"`
	Name  string `json:"name" example:"John Doe"`
}

// ResponseUser represents signup response
type ResponseUser struct {
	Message string `json:"message" example:"user created"`
}

// RefreshResponse represents refresh token response
type RefreshResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn   int    `json:"expires_in" example:"3600"`
}

// MessageResponse represents a generic message response
type MessageResponse struct {
	Message string `json:"message" example:"success"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

// CreateProductInput represents product creation request
type CreateProductInput struct {
	Name        string  `json:"name" binding:"required" example:"Laptop"`
	Description string  `json:"description" example:"High-performance laptop"`
	Price       float64 `json:"price" binding:"required" example:"999.99"`
	Quantity    int     `json:"quantity" binding:"required" example:"10"`
	Category    string  `json:"category" example:"Electronics"`
	ImageURL    string  `json:"image_url" example:"https://example.com/image.jpg"`
}

// UpdateProductInput represents product update request
type UpdateProductInput struct {
	Name        string  `json:"name" example:"Laptop"`
	Description string  `json:"description" example:"High-performance laptop"`
	Price       float64 `json:"price" example:"999.99"`
	Quantity    int     `json:"quantity" example:"10"`
	Category    string  `json:"category" example:"Electronics"`
	ImageURL    string  `json:"image_url" example:"https://example.com/image.jpg"`
	IsActive    bool    `json:"is_active" example:"true"`
}

// Product represents a product model
type Product struct {
	ID          uint      `json:"id" example:"1"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	Name        string    `json:"name" example:"Laptop"`
	Slug        string    `json:"slug" example:"laptop"`
	Description string    `json:"description" example:"High-performance laptop"`
	Price       float64   `json:"price" example:"999.99"`
	Quantity    int       `json:"quantity" example:"10"`
	Category    string    `json:"category" example:"Electronics"`
	ImageURL    string    `json:"image_url" example:"https://example.com/image.jpg"`
	IsActive    bool      `json:"is_active" example:"true"`
}

// ProductListResponse represents paginated product list
type ProductListResponse struct {
	Items      []Product `json:"items"`
	Total      int64     `json:"total" example:"100"`
	Page       int       `json:"page" example:"1"`
	Limit      int       `json:"limit" example:"10"`
	TotalPages int64     `json:"totalPages" example:"10"`
}

// AddToCartInput represents add to cart request
type AddToCartInput struct {
	ProductID uint `json:"product_id" binding:"required" example:"1"`
	Quantity  int  `json:"quantity" binding:"required,min=1" example:"2"`
}

// UpdateCartQuantityInput represents update cart quantity request
type UpdateCartQuantityInput struct {
	Quantity int `json:"quantity" binding:"required,min=1" example:"3"`
}

// CartItem represents a cart item
type CartItem struct {
	ID        uint      `json:"id" example:"1"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	UserID    uint      `json:"user_id" example:"1"`
	ProductID uint      `json:"product_id" example:"1"`
	Quantity  int       `json:"quantity" example:"2"`
	Product   Product   `json:"product"`
}

// CartResponse represents the cart response
type CartResponse struct {
	Items []CartItem `json:"items"`
	Total float64    `json:"total" example:"1999.98"`
}

// Order represents an order
type Order struct {
	ID         uint        `json:"id" example:"1"`
	CreatedAt  time.Time   `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt  time.Time   `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	UserID     uint        `json:"user_id" example:"1"`
	TotalPrice float64     `json:"total_price" example:"1999.98"`
	Status     string      `json:"status" example:"pending"`
	OrderItems []OrderItem `json:"order_items"`
}

// OrderItem represents an order item
type OrderItem struct {
	ID        uint      `json:"id" example:"1"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	OrderID   uint      `json:"order_id" example:"1"`
	ProductID uint      `json:"product_id" example:"1"`
	Quantity  int       `json:"quantity" example:"2"`
	Price     float64   `json:"price" example:"999.99"`
	Subtotal  float64   `json:"subtotal" example:"1999.98"`
	Product   Product   `json:"product"`
}

// CheckoutResponse represents checkout response
type CheckoutResponse struct {
	Message string  `json:"message" example:"order placed"`
	OrderID uint    `json:"order_id" example:"1"`
	Total   float64 `json:"total" example:"1999.98"`
}

// PaymentIntentResponse represents payment intent response
type PaymentIntentResponse struct {
	PaymentIntent uint    `json:"payment_intent" example:"1"`
	Amount        float64 `json:"amount" example:"1999.98"`
	RedirectURL   string  `json:"redirect_url" example:"http://localhost:8080/pay-gateway?intent=1"`
}

// PaymentWebhookResponse represents payment webhook response
type PaymentWebhookResponse struct {
	Message string `json:"message" example:"payment successful"`
	Order   uint   `json:"order" example:"1"`
}

// UploadResponse represents upload response
type UploadResponse struct {
	Message string `json:"message" example:"uploaded"`
	URL     string `json:"url" example:"/uploads/20240101120000.jpg"`
}

// UpdateOrderStatusInput represents order status update request
type UpdateOrderStatusInput struct {
	Status string `json:"status" binding:"required" example:"confirmed"`
}

// AdminOrderListResponse represents admin order list response
type AdminOrderListResponse struct {
	Items []Order  `json:"items"`
	Meta  MetaInfo `json:"meta"`
}

// MetaInfo represents pagination metadata
type MetaInfo struct {
	Page  int   `json:"page" example:"1"`
	Limit int   `json:"limit" example:"20"`
	Total int64 `json:"total" example:"100"`
}

// OrderStatsResponse represents order statistics response
type OrderStatsResponse struct {
	Counts  map[string]int64 `json:"counts" example:"pending:10,confirmed:20"`
	Revenue float64          `json:"revenue" example:"50000.00"`
}
