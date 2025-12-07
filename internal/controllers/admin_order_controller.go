package controllers

import (
	"net/http"
	"time"

	"ecommerce-gin/internal/database"

	"github.com/gin-gonic/gin"
)

// AdminListOrdersHandler godoc
// @Summary List all orders (Admin only)
// @Description Gets a paginated list of all orders with optional filters
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Param status query string false "Filter by status"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} AdminOrderListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/orders [get]
func AdminListOrdersHandler(c *gin.Context) {
	limit := parseIntQuery(c, "limit", 20)
	page := parseIntQuery(c, "page", 1)
	offset := (page - 1) * limit
	status := c.Query("status")
	var from, to *time.Time
	if v := c.Query("from"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			from = &t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			// end of day
			t = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			to = &t
		}
	}

	orders, total, err := database.AdminListOrders(limit, offset, status, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": orders,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// AdminGetOrderHandler godoc
// @Summary Get order by ID (Admin only)
// @Description Retrieves detailed information about a specific order
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} Order
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/admin/orders/{id} [get]
func AdminGetOrderHandler(c *gin.Context) {
	id := parseUint(c.Param("id"))
	order, err := database.AdminGetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

// AdminUpdateOrderStatusHandler godoc
// @Summary Update order status (Admin only)
// @Description Updates the status of an order with validation
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param status body UpdateOrderStatusInput true "New Status"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/admin/orders/{id}/status [put]
func AdminUpdateOrderStatusHandler(c *gin.Context) {
	id := parseUint(c.Param("id"))
	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	// Load order
	order, err := database.AdminGetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	current := order.Status
	next := body.Status

	// Basic allowed transitions (extend as needed)
	allowed := map[string][]string{
		"pending":   {"confirmed", "cancelled"},
		"confirmed": {"shipped", "cancelled"},
		"shipped":   {"delivered"},
		"delivered": {},
		"cancelled": {},
	}

	ok := false
	for _, a := range allowed[current] {
		if a == next {
			ok = true
			break
		}
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status transition"})
		return
	}

	// If cancelling, restore stock inside a transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start tx"})
		return
	}

	if next == "cancelled" {
		if err := database.RestoreStockForOrder(tx, id); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to restore stock"})
			return
		}
	}

	// Update order status
	if err := database.UpdateOrderStatus(tx, id, next); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
		return
	}

	// (Optional) send notification to user here (email / push)

	c.JSON(http.StatusOK, gin.H{"message": "status updated", "status": next})
}

// AdminOrderStatsHandler godoc
// @Summary Get order statistics (Admin only)
// @Description Retrieves order counts by status and revenue summary
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} OrderStatsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/orders/stats [get]
func AdminOrderStatsHandler(c *gin.Context) {
	counts, err := database.GetOrderCountsByStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get counts"})
		return
	}
	// optional revenue
	var from, to *time.Time
	if v := c.Query("from"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			from = &t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			t = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			to = &t
		}
	}
	revenue, _ := database.GetRevenueSummary(from, to)

	c.JSON(http.StatusOK, gin.H{
		"counts":  counts,
		"revenue": revenue,
	})
}
