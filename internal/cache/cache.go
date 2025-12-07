package cache

import (
	"time"
)

func Set(key string, value string, ttl time.Duration) error {
	return Rdb.Set(Ctx, key, value, ttl).Err()
}

func Get(key string) (string, error) {
	return Rdb.Get(Ctx, key).Result()
}

func Delete(key string) error {
	return Rdb.Del(Ctx, key).Err()
}
