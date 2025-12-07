package routes

import (
	"github.com/gin-gonic/gin"

	"ecommerce-gin/internal/controllers"
)

func RegisterOrderRoutes(r *gin.Engine, group ...*gin.RouterGroup) {
	order := group[0].Group("/orders")
	// order.Use(middleware.JWTAuth()) // dont need to add again as group already has it

	order.POST("/checkout", controllers.Checkout)
	order.GET("/", controllers.MyOrders)
	order.GET("/:id", controllers.OrderDetails)
}
