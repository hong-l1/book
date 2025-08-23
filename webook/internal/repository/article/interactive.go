package article

import (
	"context"
	cache2 "github.com/hong-l1/project/webook/internal/repository/cache"
	"github.com/hong-l1/project/webook/internal/repository/dao/article"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, id int64) error
	IncrLike(ctx context.Context, biz string, id int64, uid int64) error
	DecrLike(ctx context.Context, biz string, bizId int64, uid int64) error
	AddCollectItem(ctx context.Context, biz string, bizId int64, cid int64, uid int64) error
}
type CachedInteractiveRepository struct {
	dao   article.InteractiveDAO
	cache cache2.InteractiveCache
}

func (r *CachedInteractiveRepository) AddCollectItem(ctx context.Context, biz string, bizId int64, cid int64, uid int64) error {
	err := r.dao.InsertCollectionBiz(ctx, article.UserCollectionBiz{
		Biz:   biz,
		BizId: bizId,
		Cid:   cid,
		Uid:   uid,
	})
	if err != nil {
		return err
	}
	return r.cache.IncrCollectCntIfPresent(ctx, biz, bizId)
}

func (r *CachedInteractiveRepository) DecrLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	err := r.dao.DeleteLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return r.cache.DeleteLikeCntIfPresent(ctx, biz, bizId)
}

func (r *CachedInteractiveRepository) IncrLike(ctx context.Context, biz string, bizid int64, uid int64) error {
	err := r.dao.InsertLikeInfo(ctx, biz, bizid, uid)
	if err != nil {
		return err
	}
	return r.cache.IncrLikeCntIfPresent(ctx, biz, bizid)
}

func NewCachedInteractiveRepository(dao article.InteractiveDAO) InteractiveRepository {
	return &CachedInteractiveRepository{
		dao: dao,
	}
}
func (r *CachedInteractiveRepository) IncrReadCnt(ctx context.Context, biz string, id int64) error {
	err := r.dao.IncrReadCnt(ctx, biz, id)
	if err != nil {
		return err
	}
	return r.cache.IncrReadCntIfPresent(ctx, biz, id)
}
