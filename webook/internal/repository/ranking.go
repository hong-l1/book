package repository

import (
	"context"
	"github.com/hong-l1/project/webook/internal/domain"
	"github.com/hong-l1/project/webook/internal/repository/cache"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, art []domain.Article) error
}
type CacheRankingRepository struct {
	cache cache.RankingCache
}

func (c *CacheRankingRepository) ReplaceTopN(ctx context.Context, art []domain.Article) error {
	return c.cache.Set(ctx, art)
}
