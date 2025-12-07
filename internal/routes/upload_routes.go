package routes

import (
	"ecommerce-gin/internal/controllers"
	"ecommerce-gin/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUploadRoutes(r *gin.Engine, group ...*gin.RouterGroup) {
	// Serve files from /uploads
	r.Static("/uploads", "./uploads")

	admin := group[0].Group("/upload")
	admin.Use(middleware.AdminOnly())

	// admin.POST("/product", controllers.UploadProductImage)
	admin.POST("/product", controllers.UploadProductImageS3)
}
