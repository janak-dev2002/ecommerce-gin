package controllers

import (
	"net/http"

	"ecommerce-gin/internal/database"
	"ecommerce-gin/internal/models"

	"github.com/gin-gonic/gin"
)

// AddToCart godoc
// @Summary Add item to cart
// @Description Adds a product to the user's shopping cart
// @Tags Cart
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param item body AddToCartInput true "Cart Item"
// @Success 200 {object} MessageResponse
// @Success 201 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/cart/add [post]
func AddToCart(c *gin.Context) {
	var body struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uidRaw, _ := c.Get("user_id")
	userID := uint(uidRaw.(float64))

	// Product exists?
	product, err := database.GetProductByID(body.ProductID)
	if err != nil || product == nil || !product.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product not found or inactive"})
		return
	}

	// Enough stock?
	if body.Quantity > product.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not enough stock"})
		return
	}

	// Already in cart?
	existing, _ := database.GetCartItem(userID, body.ProductID)
	if existing != nil {
		// update quantity
		newQty := existing.Quantity + body.Quantity

		if newQty > product.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "exceeds available stock"})
			return
		}

		_ = database.UpdateCartItem(existing.ID, newQty)
		c.JSON(http.StatusOK, gin.H{"message": "quantity updated"})
		return
	}

	// Create new cart item
	item := &models.CartItem{
		UserID:    userID,
		ProductID: body.ProductID,
		Quantity:  body.Quantity,
	}

	if err := database.CreateCartItem(item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add to cart"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "added to cart"})
}

// GetCart godoc
// @Summary Get user's cart
// @Description Retrieves the user's shopping cart with all items
// @Tags Cart
// @Security BearerAuth
// @Produce json
// @Success 200 {object} CartResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/cart [get]
func GetCart(c *gin.Context) {
	
	uidRaw, _ := c.Get("user_id")
	userID := uint(uidRaw.(float64))

	items, err := database.GetCartItems(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load cart"})
		return
	}

	// Calculate totals
	var total float64 = 0
	for _, item := range items {
		total += float64(item.Quantity) * item.Product.Price
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"total": total,
	})
}

// UpdateCartQuantity godoc
// @Summary Update cart item quantity
// @Description Updates the quantity of a specific cart item
// @Tags Cart
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Cart Item ID"
// @Param quantity body UpdateCartQuantityInput true "New Quantity"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/cart/{id} [put]
func UpdateCartQuantity(c *gin.Context) {

	itemID := parseUint(c.Param("id"))

	var body struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check item
	item, err := database.GetCartItemByID(itemID)
	if err != nil || item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cart item not found"})
		return
	}

	// Check product stock
	product, _ := database.GetProductByID(item.ProductID)
	if body.Quantity > product.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not enough stock"})
		return
	}

	database.UpdateCartItem(itemID, body.Quantity)
	c.JSON(http.StatusOK, gin.H{"message": "cart updated"})
}

// RemoveFromCart godoc
// @Summary Remove item from cart
// @Description Removes a specific item from the cart
// @Tags Cart
// @Security BearerAuth
// @Produce json
// @Param id path string true "Cart Item ID"
// @Success 200 {object} MessageResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/cart/{id} [delete]
func RemoveFromCart(c *gin.Context) {
	id := parseUint(c.Param("id"))
	if err := database.RemoveCartItem(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "removed"})
}

// ClearCart godoc
// @Summary Clear entire cart
// @Description Removes all items from the user's cart
// @Tags Cart
// @Security BearerAuth
// @Produce json
// @Success 200 {object} MessageResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/cart [delete]
func ClearCart(c *gin.Context) {
	uidRaw, _ := c.Get("user_id")
	userID := uint(uidRaw.(float64))

	if err := database.ClearCart(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cart cleared"})
}
