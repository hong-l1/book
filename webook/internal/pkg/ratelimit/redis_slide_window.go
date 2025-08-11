package ratelimit

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed slide_window.lua
var luaslidewindow string

type redisSlideWindowLimit struct {
	cmd      redis.Cmdable
	interval time.Duration
	// 阈值
	rate int
}

func NewredisSlideWindowLimit(cmd redis.Cmdable, interval time.Duration, rate int) *redisSlideWindowLimit {
	return &redisSlideWindowLimit{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

func (r *redisSlideWindowLimit) Limited(ctx context.Context, key string) (bool, error) {
	return r.cmd.Eval(ctx, luaslidewindow, []string{key},
		r.interval.Milliseconds(), r.rate, time.Now().UnixMilli()).Bool()
}
