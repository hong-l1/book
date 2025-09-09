package cache

import (
	"context"
	"encoding/json"
	"github.com/hong-l1/project/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type RankingCache interface {
	Set(ctx context.Context, val []domain.Article) error
	Get(ctx context.Context) ([]domain.Article, error)
}
type RedisRanking struct {
	client redis.Cmdable
	key    string
}

func (r *RedisRanking) Set(ctx context.Context, val []domain.Article) error {
	for k := range val {
		val[k].Content = ""
	}
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key, data, time.Minute*10).Err()
}

func (r *RedisRanking) Get(ctx context.Context) ([]domain.Article, error) {
	data, err := r.client.Get(ctx, r.key).Bytes()
	if err != nil {
		return nil, err
	}
	var val []domain.Article
	if err := json.Unmarshal(data, &val); err != nil {
		return nil, err
	}
	return val, nil
}
