package routes

import (
	"ecommerce-gin/internal/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterCartRoutes(r *gin.Engine, group ...*gin.RouterGroup) {

	cart := group[0].Group("/cart")
	// cart.Use(middleware.JWTAuth()) // dont need to add again as group already has it

	cart.POST("/add", controllers.AddToCart)
	cart.GET("/", controllers.GetCart)
	cart.PUT("/:id", controllers.UpdateCartQuantity)
	cart.DELETE("/:id", controllers.RemoveFromCart)
	cart.DELETE("/", controllers.ClearCart)
}
