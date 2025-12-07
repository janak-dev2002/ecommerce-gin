package main

// @title E-Commerce API
// @version 1.0
// @description This is a backend service for the Go e-commerce project with comprehensive authentication, product management, cart, orders, and payment functionality.
// @contact.name API Support
// @contact.email support@ecommerce.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

import (
	"ecommerce-gin/internal/cache"
	"ecommerce-gin/internal/config"
	"ecommerce-gin/internal/database"
	"ecommerce-gin/internal/routes"
	"ecommerce-gin/internal/services"
	"fmt"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "ecommerce-gin/docs" // Import generated docs
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect DB and automigrate
	database.Connect()

	// Initialize S3 client
	services.InitS3()

	// Connect Redis
	cache.Connect()

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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	fmt.Println("Server running on port", config.Cfg.Port)
	r.Run(":" + config.Cfg.Port)
}
