package repository

import (
	"context"
	"github.com/hong-l1/project/webook/internal/domain"
	"github.com/hong-l1/project/webook/internal/repository/cache"
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
	return c.redis.Get(ctx)
}

func (c *CacheRankingRepository) ReplaceTopN(ctx context.Context, art []domain.Article) error {
	return c.redis.Set(ctx, art)
}
func NewCacheRankingRepository(cache cache.RankingCache) RankingRepository {
	return &CacheRankingRepository{redis: cache}
}
