package controllers

import (
	"strconv"

	"ecommerce-gin/internal/database"

	"github.com/gin-gonic/gin"
)

// PaymentWebhook godoc
// @Summary Payment webhook
// @Description Handles payment gateway callbacks
// @Tags Payment
// @Produce json
// @Param intent_id query string true "Payment Intent ID"
// @Success 200 {object} PaymentWebhookResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/payment/webhook [post]
func PaymentWebhook(c *gin.Context) {
	// in real world gateway sends JSON
	intentID, _ := strconv.ParseUint(c.Query("intent_id"), 10, 64)

	intent, err := database.GetPaymentIntentByID(uint(intentID))
	if err != nil {
		c.JSON(404, gin.H{"error": "intent not found"})
		return
	}

	// mark payment as paid
	database.UpdatePaymentIntentStatus(uint(intentID), "paid")

	// update order to confirmed
	database.UpdateOrderStatus(nil, intent.OrderID, "confirmed")

	c.JSON(200, gin.H{
		"message": "payment successful",
		"order":   intent.OrderID,
	})
}
