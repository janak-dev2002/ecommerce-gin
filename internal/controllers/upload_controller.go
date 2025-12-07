package controllers

import (
	"ecommerce-gin/internal/config"
	"ecommerce-gin/internal/services"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadProductImage(c *gin.Context) {
	// Single file upload
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing image"})
		return
	}

	// Validate extension
	ext := filepath.Ext(file.Filename)
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}

	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file type"})
		return
	}

	// Generate unique file name
	newName := time.Now().Format("20060102150405") + ext
	savePath := "./uploads/" + newName

	// Save file
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "uploaded",
		"url":     "/uploads/" + newName,
	})
}

func UploadProductImageS3(c *gin.Context) {

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "image missing"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "cannot open"})
		return
	}
	defer f.Close()

	url, err := services.UploadToS3(f, file)
	if err != nil {
		c.JSON(500, gin.H{"error": "upload failed"})
		return
	}

	c.JSON(200, gin.H{
		"url": config.Cfg.S3PublicURL + "/" + url,
	})
}
