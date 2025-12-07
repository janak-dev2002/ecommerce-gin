package controllers

import (
	"net/http"
	"time"

	"ecommerce-gin/internal/database"

	"github.com/gin-gonic/gin"
)

// AdminListOrdersHandler GET /admin/orders
// query: status, page, limit, from (YYYY-MM-DD), to (YYYY-MM-DD)
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

// AdminGetOrderHandler GET /admin/orders/:id
func AdminGetOrderHandler(c *gin.Context) {
	id := parseUint(c.Param("id"))
	order, err := database.AdminGetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

// AdminUpdateOrderStatusHandler PUT /admin/orders/:id/status
// Body: { "status": "confirmed" } — allowed transitions handled here
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

// AdminOrderStatsHandler GET /admin/orders/stats
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
