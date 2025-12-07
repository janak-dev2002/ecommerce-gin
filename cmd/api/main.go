package main

import (
	"ecommerce-gin/internal/config"
	"ecommerce-gin/internal/database"
	"ecommerce-gin/internal/routes"
	"ecommerce-gin/internal/services"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect DB and automigrate
	database.Connect()

	// Initialize S3 client
	services.InitS3()

	// Gin setup
	if config.Cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.GET("/pay-gateway", func(c *gin.Context) {
		intent := c.Query("intent")

		c.JSON(200, gin.H{
			"message":      "This simulates a real payment page.",
			"instructions": "POST /webhook/payment to simulate success",
			"intent_id":    intent,
		})
	})

	// Setup routes
	api := routes.SetupRoutes(r)
	routes.RegisterProductRoutes(r, api)
	routes.RegisterCartRoutes(r, api)
	routes.RegisterOrderRoutes(r, api)
	routes.RegisterAdminOrderRoutes(r, api)
	routes.RegisterPaymentRoutes(r, api)
	routes.RegisterUploadRoutes(r, api)

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	fmt.Println("Server running on port", config.Cfg.Port)
	r.Run(":" + config.Cfg.Port)
}
