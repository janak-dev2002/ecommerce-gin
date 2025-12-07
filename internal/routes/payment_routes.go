package routes

import (
	"ecommerce-gin/internal/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterPaymentRoutes(r *gin.Engine, group ...*gin.RouterGroup) {

	p := group[0].Group("/payment")

	p.POST("/start/:orderId", controllers.StartPaymentHandler)
	p.POST("/webhook", controllers.PaymentWebhook)
}
