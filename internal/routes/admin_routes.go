package routes

import (
	"ecommerce-gin/internal/controllers"
	"ecommerce-gin/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAdminOrderRoutes(r *gin.Engine, group ...*gin.RouterGroup) {
	admin := group[0].Group("/admin")
	admin.Use(middleware.AdminOnly())

	// Orders
	admin.GET("/orders", controllers.AdminListOrdersHandler)
	admin.GET("/orders/:id", controllers.AdminGetOrderHandler)
	admin.PUT("/orders/:id/status", controllers.AdminUpdateOrderStatusHandler)
	admin.GET("/orders/stats", controllers.AdminOrderStatsHandler)
}
