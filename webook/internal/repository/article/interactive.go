package article

import (
	"github.com/gin-gonic/gin"
	"github.com/hong-l1/project/webook/internal/repository/dao/article"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx *gin.Context, biz string, id int64) error
}
type CachedInteractiveRepository struct {
	dao article.InteractiveDAO
}

func NewCachedInteractiveRepository(dao article.InteractiveDAO) *CachedInteractiveRepository {
	return &CachedInteractiveRepository{
		dao: dao,
	}
}
func (r *CachedInteractiveRepository) IncrReadCnt(ctx *gin.Context, biz string, id int64) error {
	return r.dao.IncrReadCnt(ctx, biz, id)
}
