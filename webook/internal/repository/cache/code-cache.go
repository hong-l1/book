package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

//go:embed Lua/set_code.lua
var LuaSetCode string

//go:embed Lua/verify.lua
var VerifyCode string
var (
	ErrSendTooMany   = errors.New("send too many")
	ErrVerifyTooMany = errors.New("verify too many")
	ErrCodeWrong     = errors.New("wrong code") // 新增：验证码错误
	ErrCodeExpire    = errors.New("code expired or not exists")
	ErrUnknown       = errors.New("unknown")
)

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) error
}
type RedisCodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) CodeCache {
	return &RedisCodeCache{client: client}
}
func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.client.Eval(ctx, LuaSetCode, []string{c.Key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		zap.L().Warn("短信发送太频繁", zap.String("biz", biz))
		return ErrSendTooMany
	default:
		return errors.New("system error")
	}
}
func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) error {
	res, err := c.client.Eval(ctx, VerifyCode, []string{c.Key(biz, phone)}, inputCode).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil // 验证成功
	case -1:
		return ErrVerifyTooMany // 超过最大验证次数
	case -2:
		return ErrCodeWrong // 验证码错误
	case -3:
		return ErrCodeExpire // 验证码不存在或过期
	default:
		return ErrUnknown
	}
}

func (c *RedisCodeCache) Key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
