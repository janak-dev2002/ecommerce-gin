package middleware

import (
	"ecommerce-gin/internal/cache"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimit(maxPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "rl:" + ip

		count, _ := cache.Rdb.Incr(cache.Ctx, key).Result()
		if count == 1 {
			cache.Rdb.Expire(cache.Ctx, key, time.Minute)
		}

		if count > int64(maxPerMinute) {
			c.JSON(429, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}
