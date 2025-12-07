package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"ecommerce-gin/internal/database"
	"ecommerce-gin/internal/models"

	"github.com/gin-gonic/gin"
)

// Helper to create slugs
func generateSlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}

// --- ADMIN ONLY ---
// Create Product
func CreateProduct(c *gin.Context) {
	var body struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"required"`
		Quantity    int     `json:"quantity" binding:"required"`
		Category    string  `json:"category"`
		ImageURL    string  `json:"image_url"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slug := generateSlug(body.Name)

	product := &models.Product{
		Name:        body.Name,
		Slug:        slug,
		Description: body.Description,
		Price:       body.Price,
		Quantity:    body.Quantity,
		Category:    body.Category,
		ImageURL:    body.ImageURL,
		IsActive:    true,
	}

	if err := database.CreateProduct(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// Update Product (admin)
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body["updated_at"] = time.Now()

	// regenerate slug if name changed
	if name, exists := body["name"]; exists {
		body["slug"] = generateSlug(name.(string))
	}

	if err := database.UpdateProduct(parseUint(id), body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// Delete (soft delete)
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := database.DeleteProduct(parseUint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// --- PUBLIC ---
// Get one product by slug
func GetProduct(c *gin.Context) {
	slug := c.Param("slug")

	product, err := database.GetProductBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// List products with filtering & pagination
func ListProducts(c *gin.Context) {
	
	page := parseIntQuery(c, "page", 1)
	limit := parseIntQuery(c, "limit", 10)
	search := c.Query("search")
	category := c.Query("category")

	offset := (page - 1) * limit

	products, total, err := database.ListProducts(limit, offset, search, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":      products,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
}

// helpers
func parseUint(s string) uint {
	var id uint
	fmt.Sscanf(s, "%d", &id)
	return id
}

func parseIntQuery(c *gin.Context, key string, def int) int {
	val := c.Query(key)
	if val == "" {
		return def
	}
	var num int
	fmt.Sscanf(val, "%d", &num)
	return num
}
