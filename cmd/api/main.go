package main

import (
	"ecommerce-gin/internal/config"
	"ecommerce-gin/internal/database"
	"ecommerce-gin/internal/routes"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect DB and automigrate
	database.Connect()

	// Gin setup
	if config.Cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r)

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	fmt.Println("Server running on port", config.Cfg.Port)
	r.Run(":" + config.Cfg.Port)
}
