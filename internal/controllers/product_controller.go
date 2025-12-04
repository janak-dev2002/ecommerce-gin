package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"ecommerce-gin/internal/database"
	"ecommerce-gin/internal/models"

	"github.com/gin-gonic/gin"
)

// Request DTOs -------------------------------------------------------

type CreateProductReq struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Discount    float64 `json:"discount" binding:"gte=0,lte=100"`
	SKU         string  `json:"sku" binding:"required"`
	Stock       int     `json:"stock" binding:"gte=0"`
	CategoryID  *uint   `json:"category_id" binding:"required"`
	ImageURL    string  `json:"image_url"` // optional, or set by upload endpoint
}

type UpdateProductReq struct {
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price" binding:"omitempty,gt=0"`
	Discount    *float64 `json:"discount" binding:"omitempty,gte=0,lte=100"`
	Stock       *int     `json:"stock" binding:"omitempty,gte=0"`
	CategoryID  *uint    `json:"category_id"`
	ImageURL    *string  `json:"image_url"`
}

// Helpers ------------------------------------------------------------

// isAdmin checks role from context (set by JWT middleware)
// returns true if role == "admin"
func isAdmin(c *gin.Context) bool {
	role, _ := c.Get("role")
	rs, _ := role.(string)
	return strings.ToLower(rs) == "admin"
}

// Controllers --------------------------------------------------------

// CreateProduct - Admin only
func CreateProduct(c *gin.Context) {

	if !isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
		return
	}

	var req CreateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Build model
	p := &models.Product{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Discount:    req.Discount,
		SKU:         req.SKU,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		CategoryID:  req.CategoryID,
	}

	if err := database.CreateProduct(p); err != nil {
		// unique constraint on SKU possible
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": p})
}

// UpdateProduct - Admin only (partial updates)
func UpdateProduct(c *gin.Context) {
	if !isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
		return
	}

	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	id := uint(id64)

	var req UpdateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Build update map
	data := map[string]interface{}{}
	if req.Title != nil {
		data["title"] = *req.Title
	}
	if req.Description != nil {
		data["description"] = *req.Description
	}
	if req.Price != nil {
		data["price"] = *req.Price
	}
	if req.Discount != nil {
		data["discount"] = *req.Discount
	}
	if req.Stock != nil {
		data["stock"] = *req.Stock
	}
	if req.CategoryID != nil {
		data["category_id"] = *req.CategoryID
	}
	if req.ImageURL != nil {
		data["image_url"] = *req.ImageURL
	}

	if len(data) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	if err := database.UpdateProduct(id, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product", "details": err.Error()})
		return
	}

	updated, _ := database.GetProductByID(id) // reload
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

// DeleteProduct - Admin only
func DeleteProduct(c *gin.Context) {
	if !isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
		return
	}

	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	id := uint(id64)

	if err := database.DeleteProduct(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetProduct - Public
func GetProduct(c *gin.Context) {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	id := uint(id64)

	product, err := database.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch product"})
		return
	}
	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": product})
}

// ListProducts - Public with filters & pagination
func ListProducts(c *gin.Context) {
	// read query params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "12"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 12
	}

	filters := map[string]interface{}{
		"page":  page,
		"limit": limit,
	}

	// optional filters
	if s := c.Query("search"); s != "" {
		filters["search"] = s
	}
	if cat := c.Query("category_id"); cat != "" {
		if v, err := strconv.Atoi(cat); err == nil {
			filters["category_id"] = v
		}
	}
	if min := c.Query("min_price"); min != "" {
		if v, err := strconv.ParseFloat(min, 64); err == nil {
			filters["min_price"] = v
		}
	}
	if max := c.Query("max_price"); max != "" {
		if v, err := strconv.ParseFloat(max, 64); err == nil {
			filters["max_price"] = v
		}
	}
	if sort := c.Query("sort"); sort != "" {
		filters["sort"] = sort
	}

	products, total, err := database.ListProducts(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list products", "details": err.Error()})
		return
	}

	meta := gin.H{
		"page":  page,
		"limit": limit,
		"total": total,
		"next":  fmt.Sprintf("/api/products?page=%d&limit=%d", page+1, limit),
	}

	c.JSON(http.StatusOK, gin.H{"data": products, "meta": meta})
}

// UploadProductImage - Admin only, saves locally to ./uploads/ and returns URL
func UploadProductImage(c *gin.Context) {
	if !isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
		return
	}

	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	id := uint(id64)

	// get file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image is required"})
		return
	}

	// ensure uploads directory exists
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload dir"})
		return
	}

	// build filename: product-<id>-<origname>
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("product-%d-%d%s", id, timeNowUnix(), ext)
	dst := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file", "details": err.Error()})
		return
	}

	// You may want to store a full URL (e.g., https://cdn.example.com/uploads/...)
	// For now store local path
	imageURL := "/uploads/" + filename

	// update product record
	_ = database.UpdateProduct(id, map[string]interface{}{"image_url": imageURL})

	updated, _ := database.GetProductByID(id)

	c.JSON(http.StatusOK, gin.H{"data": updated})
}

// helper to return unix timestamp int for filename
func timeNowUnix() int64 {
	return (int64)(float64(time.Now().Unix()))
}
