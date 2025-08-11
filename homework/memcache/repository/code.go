package repository

import "github.com/hong-l1/project/homework/memcache/repository/cache"

type CodeRepository struct {
	CodeCache *cache.CodeCache
}

func NewCodeRepository(CodeCache *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		CodeCache: CodeCache,
	}
}
func (c *CodeRepository) Set(biz, phone, code string) error {
	return c.CodeCache.SetCode(biz, phone, code)
}
func (c *CodeRepository) Verify(biz, phone, code string) error {
	return c.CodeCache.VerifyCode(biz, phone, code)
}
