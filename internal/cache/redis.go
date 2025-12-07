package cache

import (
	"context"
	"fmt"

	"ecommerce-gin/internal/config"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var Ctx = context.Background()

func Connect() {
	cfg := config.Cfg

	Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		panic("Redis connection failed: " + err.Error())
	}

	fmt.Println("Redis connected")
}
