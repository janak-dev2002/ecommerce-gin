package routes

import (
	"ecommerce-gin/internal/controllers"
	"ecommerce-gin/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r *gin.Engine, group ...*gin.RouterGroup) {

	// Public
	r.GET("/products", controllers.ListProducts) // /api/products
	r.GET("/products/:slug", middleware.RateLimit(2), controllers.GetProduct)

	// Admin only
	admin := group[0].Group("/admin/products") // /api/admin/products
	admin.Use(middleware.AdminOnly())

	admin.POST("/", controllers.CreateProduct)
	admin.PUT("/:id", controllers.UpdateProduct)
	admin.DELETE("/:id", controllers.DeleteProduct)
}
