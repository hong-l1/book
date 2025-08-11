package startup

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	addr := viper.GetString("redis.addr")
	fmt.Printf("Connecting to redis at %s\n", addr)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: "default",
		Password: "123456",
	})
	return redisClient
}
