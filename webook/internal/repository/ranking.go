package repository

import (
	"context"
	"github.com/hong-l1/project/webook/internal/domain"
	"github.com/hong-l1/project/webook/internal/repository/cache"
	"time"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, art []domain.Article) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}
type CacheRankingRepository struct {
	redis cache.RankingCache
	local cache.RankingLocalCache
}

func (c *CacheRankingRepository) GetTopN(ctx context.Context) ([]domain.Article, error) {
	data, err := c.local.Get(ctx, "")
	if err == nil {
		return data, nil
	}
	data, err = c.redis.Get(ctx)
	if err == nil {
		_ = c.local.Set(ctx, "", data, time.Minute*5)
	} else {
		return c.local.Foreget(ctx)
	}
	return data, err
}

func (c *CacheRankingRepository) ReplaceTopN(ctx context.Context, art []domain.Article) error {
	_ = c.local.Set(ctx, "", art, time.Minute*5)
	return c.redis.Set(ctx, art)
}
func NewCacheRankingRepository(cache cache.RankingCache, local cache.RankingLocalCache) RankingRepository {
	return &CacheRankingRepository{
		redis: cache,
		local: local,
	}
}
