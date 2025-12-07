package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"ecommerce-gin/internal/cache"
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

// CreateProduct godoc
// @Summary Create a new product (Admin only)
// @Description Creates a new product in the system
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param product body CreateProductInput true "Product Data"
// @Success 201 {object} Product
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/admin/products [post]
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

// UpdateProduct godoc
// @Summary Update a product (Admin only)
// @Description Updates product details
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body UpdateProductInput true "Product Update Data"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/admin/products/{id} [put]
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	cache.Delete("product:" + id) // invalidate cache

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

// DeleteProduct godoc
// @Summary Delete a product (Admin only)
// @Description Soft deletes a product from the system
// @Tags Products
// @Security BearerAuth
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} MessageResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	cache.Delete("product:" + id) // invalidate cache

	if err := database.DeleteProduct(parseUint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// GetProduct godoc
// @Summary Get product by slug
// @Description Retrieves a single product by its slug (cached)
// @Tags Products
// @Produce json
// @Param slug path string true "Product Slug"
// @Success 200 {object} Product
// @Failure 404 {object} ErrorResponse
// @Router /products/{slug} [get]
func GetProduct(c *gin.Context) {
	slug := c.Param("slug")
	cacheKey := "product:" + slug

	// Check cache
	if cached, err := cache.Get(cacheKey); err == nil {
		// debgug print
		fmt.Println("Cache hit for product:", slug)
		c.Data(200, "application/json", []byte(cached))
		return
	}

	// Load from DB
	product, err := database.GetProductBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	jsonData, _ := json.Marshal(product)

	// Save in Redis for 5 minutes
	cache.Set(cacheKey, string(jsonData), 5*time.Minute)

	// Return response
	c.JSON(http.StatusOK, product)
}

// ListProducts godoc
// @Summary List all products
// @Description Gets a paginated list of products with optional filtering
// @Tags Products
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Param category query string false "Filter by category"
// @Success 200 {object} ProductListResponse
// @Router /products [get]
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
