package routes

import (
	"ecommerce-gin/internal/controllers"
	"ecommerce-gin/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Public auth routes
	r.POST("/auth/signup", controllers.Signup)
	r.POST("/auth/login", controllers.Login)
	r.POST("/auth/refresh", controllers.Refresh)

	// Protected API group
	api := r.Group("/api")
	api.Use(middleware.JWTAuth())

	// logout (requires valid access token)
	api.POST("/auth/logout", controllers.Logout)

	// example protected endpoint
	api.GET("/profile", func(c *gin.Context) {
		uid, _ := c.Get("user_id")
		role := c.GetString("role")
		c.JSON(200, gin.H{
			"user_id": uid,
			"role":    role,
		})
	})
}
