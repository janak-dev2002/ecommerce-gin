package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"ecommerce-gin/internal/database"
	"ecommerce-gin/internal/models"
)

// Checkout godoc
// @Summary Checkout cart
// @Description Creates an order from the user's cart
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Success 200 {object} CheckoutResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/orders/checkout [post]
func Checkout(c *gin.Context) {
	uidRaw, _ := c.Get("user_id")
	userID := uint(uidRaw.(float64))

	// 1. Load cart items
	cartItems, err := database.GetCartItems(userID)
	if err != nil || len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart is empty"})
		return
	}

	// 2. Begin DB transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start transaction"})
		return
	}

	var total float64
	var orderItems []models.OrderItem

	// 3. Validate stock + prepare order items
	for _, ci := range cartItems {
		if ci.Quantity > ci.Product.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "not enough stock for product " + ci.Product.Name,
			})
			return
		}

		price := ci.Product.Price
		subtotal := float64(ci.Quantity) * price

		orderItems = append(orderItems, models.OrderItem{
			ProductID: ci.ProductID,
			Quantity:  ci.Quantity,
			Price:     price,
			Subtotal:  subtotal,
		})

		total += subtotal

		// Deduct stock
		err := database.DeductProductStock(tx, ci.ProductID, ci.Quantity)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "stock update failed"})
			return
		}
	}

	// 4. Create Order
	order := models.Order{
		UserID:     userID,
		TotalPrice: total,
		Status:     "pending",
	}

	if err := database.CreateOrder(tx, &order); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	// assign order_id to orderItems
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
	}

	// 5. Save Order Items
	if err := database.CreateOrderItems(tx, orderItems); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order items"})
		return
	}

	// 6. Clear cart
	if err := database.ClearCartInTransaction(tx, userID); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear cart"})
		return
	}

	// 7. Commit
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "checkout commit failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "order placed",
		"order_id": order.ID,
		"total":    total,
	})
}

// MyOrders godoc
// @Summary Get user's orders
// @Description Retrieves all orders for the authenticated user
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Success 200 {array} Order
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/orders/ [get]
func MyOrders(c *gin.Context) {
	uidRaw, _ := c.Get("user_id")
	userID := uint(uidRaw.(float64))

	orders, err := database.GetOrdersForUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// OrderDetails godoc
// @Summary Get order details
// @Description Retrieves details of a specific order
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} Order
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/orders/{id} [get]
func OrderDetails(c *gin.Context) {
	
	uidRaw, _ := c.Get("user_id")
	userID := uint(uidRaw.(float64))

	orderIDStr := c.Param("id")
	orderID := uint(0)
	if _, err := fmt.Sscan(orderIDStr, &orderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	order, err := database.GetOrderByID(userID, orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}
