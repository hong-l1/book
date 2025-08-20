package cache

import (
	"github.com/redis/go-redis/v9"
)

type InteractiveCache interface {
}
type RedisInteractiveCache struct {
	client *redis.Cmdable
}

func NewRedisInteractiveCache(client *redis.Cmdable) *RedisInteractiveCache {
	return &RedisInteractiveCache{
		client: client,
	}
}
