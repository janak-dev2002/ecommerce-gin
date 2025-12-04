package controllers

import (
	"net/http"
	"time"

	"ecommerce-gin/internal/config"
	"ecommerce-gin/internal/database"
	"ecommerce-gin/internal/models"
	"ecommerce-gin/internal/utils"

	"github.com/gin-gonic/gin"
)

// Signup creates a new user (password hashed)
func Signup(c *gin.Context) {
	var body struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role"` // optional
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Check exists
	existing, err := database.GetUserByEmail(body.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	hashed, err := utils.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	role := body.Role
	if role == "" {
		role = "customer"
	}

	user := &models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: hashed,
		Role:     role,
	}

	if err := database.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

// Login authenticates and issues access token + sets refresh cookie
func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, err := database.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	hashedRT, err := utils.HashToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash refresh token"})
		return
	}

	expiry := time.Now().Add(time.Hour * 24 * time.Duration(config.Cfg.RefreshTokenDays))

	if err := database.SaveRefreshToken(user.ID, hashedRT, expiry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save refresh token"})
		return
	}

	// set cookie; in production set Secure=true
	c.SetCookie("refresh_token", refreshToken, int(time.Until(expiry).Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"expires_in":   config.Cfg.AccessTokenMinutes * 60, // seconds
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
			"name":  user.Name,
		},
	})
}

// Refresh rotates the refresh token and issues a new access token
func Refresh(c *gin.Context) {
	rt, err := c.Cookie("refresh_token")
	if err != nil || rt == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token"})
		return
	}

	user, err := database.FindUserByActiveRefreshToken(rt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	// rotate refresh token
	newRT, err := utils.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh"})
		return
	}
	newHashed, _ := utils.HashToken(newRT)
	newExpiry := time.Now().Add(time.Hour * 24 * time.Duration(config.Cfg.RefreshTokenDays))

	if err := database.SaveRefreshToken(user.ID, newHashed, newExpiry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save refresh token"})
		return
	}

	// set cookie
	c.SetCookie("refresh_token", newRT, int(time.Until(newExpiry).Seconds()), "/", "", false, true)

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"expires_in":   config.Cfg.AccessTokenMinutes * 60,
	})
}

// Logout clears stored refresh token and removes cookie
func Logout(c *gin.Context) {
	uidRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// tokens from jwt are numeric (float64) — handle accordingly
	var uid uint
	switch v := uidRaw.(type) {
	case float64:
		uid = uint(v)
	case float32:
		uid = uint(v)
	case int:
		uid = uint(v)
	case uint:
		uid = v
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return
	}

	if err := database.ClearRefreshToken(uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear refresh token"})
		return
	}

	// clear cookie
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
