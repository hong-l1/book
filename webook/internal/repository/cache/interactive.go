package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/hong-l1/project/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"strconv"
)

//go:embed Lua/interactive_incr.lua
var luaIncrCnt string

const fieldReadCnt = "read_cnt"
const fieldLikeCnt = "like_cnt"
const fieldCollectCnt = "collect_cnt"

type InteractiveCache interface {
	IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	DeleteLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Set(ctx context.Context, biz string, bizId int64, inter domain.Interactive) error
}
type RedisInteractiveCache struct {
	client redis.Cmdable
}

func (r *RedisInteractiveCache) Set(ctx context.Context, biz string, bizId int64, inter domain.Interactive) error {
	panic("implement me")
}

func (r *RedisInteractiveCache) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	key := r.key(biz, bizId)
	data, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
	}
	if len(data) == 0 {
		return domain.Interactive{}, err
	}
	readCnt, _ := strconv.ParseInt(data[fieldReadCnt], 10, 64)
	likeCnt, _ := strconv.ParseInt(data[fieldLikeCnt], 10, 64)
	collectCnt, _ := strconv.ParseInt(data[fieldCollectCnt], 10, 64)
	return domain.Interactive{
		ReadCnt:    readCnt,
		LikeCnt:    likeCnt,
		CollectCnt: collectCnt,
	}, nil
}

func (r *RedisInteractiveCache) IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	key := r.key(biz, bizId)
	return r.client.Eval(ctx, luaIncrCnt, []string{key}, fieldCollectCnt, 1).Err()
}

func (r *RedisInteractiveCache) DeleteLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	key := r.key(biz, bizId)
	return r.client.Eval(ctx, luaIncrCnt, []string{key}, fieldLikeCnt, -1).Err()
}

func (r *RedisInteractiveCache) IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	key := r.key(biz, bizId)
	return r.client.Eval(ctx, luaIncrCnt, []string{key}, fieldLikeCnt, 1).Err()
}
func (r *RedisInteractiveCache) IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	key := r.key(biz, bizId)
	return r.client.Eval(ctx, luaIncrCnt, []string{key}, fieldReadCnt, 1).Err()
}
func NewRedisInteractiveCache(client redis.Cmdable) InteractiveCache {
	return &RedisInteractiveCache{client: client}
}
func (r *RedisInteractiveCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}
