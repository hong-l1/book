package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hong-l1/project/webook/internal/domain"
	"github.com/redis/go-redis/v9"

	"time"
)

var ErrKeyNotFound = redis.Nil

type UserCache interface {
	GetUserCache(ctx context.Context, user domain.User) (domain.User, error)
	SetUserCache(ctx context.Context, user domain.User) error
}
type RedisUserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

// A用到了B,B一定是接口
// A用到了B，B一定是A的字段
// A用到了B,A绝对不初始化，而是另外注入
func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

// 如果没有数据返回一个特定的err
// /err如果为空，user一定存在
func (c *RedisUserCache) GetUserCache(ctx context.Context, user domain.User) (domain.User, error) {
	key := c.Key(user.Id)
	val, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	u := domain.User{}
	err = json.Unmarshal(val, &u)
	return u, nil
}
func (c *RedisUserCache) SetUserCache(ctx context.Context, user domain.User) error {
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := c.Key(user.Id)
	return c.client.Set(ctx, key, val, c.expiration).Err()
}
func (c *RedisUserCache) Key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
