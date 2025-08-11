package repository

import (
	"context"
	"github.com/hong-l1/project/webook/internal/repository/cache"
)

var ErrSendTooMany = cache.ErrSendTooMany
var ErrVerifyTooMany = cache.ErrVerifyTooMany

type CodeRepository interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) error
}
type CacheCodeRepository struct {
	CodeCache cache.CodeCache
}

func NewCodeRepository(CodeCache cache.CodeCache) CodeRepository {
	return &CacheCodeRepository{
		CodeCache: CodeCache,
	}
}
func (c *CacheCodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return c.CodeCache.Set(ctx, biz, phone, code)
}
func (c *CacheCodeRepository) Verify(ctx context.Context, biz, phone, inputCode string) error {
	return c.CodeCache.Verify(ctx, biz, phone, inputCode)
}
