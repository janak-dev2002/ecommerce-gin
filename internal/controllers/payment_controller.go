package controllers

import (
	"net/http"
	"strconv"

	"ecommerce-gin/internal/database"
	"ecommerce-gin/internal/services"

	"github.com/gin-gonic/gin"
)

func StartPaymentHandler(c *gin.Context) {
	
	orderID, _ := strconv.ParseUint(c.Param("orderId"), 10, 64)

	// Get order
	order, err := database.AdminGetOrderByID(uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Create payment intent
	intent, err := services.CreatePaymentIntent(order.ID, order.TotalPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create payment intent"})
		return
	}

	// Fake payment gateway checkout URL
	redirectURL := "http://localhost:8080/pay-gateway?intent=" + strconv.Itoa(int(intent.ID))

	c.JSON(http.StatusOK, gin.H{
		"payment_intent": intent.ID,
		"amount":         intent.Amount,
		"redirect_url":   redirectURL,
	})
}
