package article

import (
	"context"
	"github.com/hong-l1/project/webook/internal/domain"
	logger2 "github.com/hong-l1/project/webook/internal/pkg/logger"
	cache2 "github.com/hong-l1/project/webook/internal/repository/cache"
	"github.com/hong-l1/project/webook/internal/repository/dao/article"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, id int64) error
	IncrLike(ctx context.Context, biz string, id int64, uid int64) error
	DecrLike(ctx context.Context, biz string, bizId int64, uid int64) error
	AddCollectItem(ctx context.Context, biz string, bizId int64, cid int64, uid int64) error
	Get(ctx context.Context, biz string, artid int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, artid int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, artid int64, uid int64) (bool, error)
}
type CachedInteractiveRepository struct {
	dao   article.InteractiveDAO
	cache cache2.InteractiveCache
	l     logger2.Loggerv1
}

func (r *CachedInteractiveRepository) Get(ctx context.Context, biz string, artid int64) (domain.Interactive, error) {
	inter, err := r.cache.Get(ctx, biz, artid)
	if err == nil {
		return inter, nil
	}
	daointer, err := r.dao.Get(ctx, biz, artid)
	if err != nil {
		return domain.Interactive{}, err
	}
	inter = r.toDomain(daointer)
	go func() {
		er := r.cache.Set(ctx, biz, artid, inter)
		r.l.Error("回写缓存失败",
			logger2.String("biz", biz),
			logger2.Int64("artid", artid),
			logger2.Error(er))
	}()
	return inter, nil
}

func (r *CachedInteractiveRepository) Liked(ctx context.Context, biz string, artid int64, uid int64) (bool, error) {
	_, err := r.dao.GetLikeInfo(ctx, biz, artid, uid)
	switch err {
	case nil:
		return true, nil
	case article.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (r *CachedInteractiveRepository) Collected(ctx context.Context, biz string, artid int64, uid int64) (bool, error) {
	_, err := r.dao.GetCollectInfo(ctx, biz, artid, uid)
	switch err {
	case nil:
		return true, nil
	case article.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
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

func NewCachedInteractiveRepository(dao article.InteractiveDAO, l logger2.Loggerv1) InteractiveRepository {
	return &CachedInteractiveRepository{
		dao: dao,
		l:   l,
	}
}
func (r *CachedInteractiveRepository) IncrReadCnt(ctx context.Context, biz string, id int64) error {
	err := r.dao.IncrReadCnt(ctx, biz, id)
	if err != nil {
		return err
	}
	return r.cache.IncrReadCntIfPresent(ctx, biz, id)
}
func (r *CachedInteractiveRepository) toDomain(dao article.Interactive) domain.Interactive {
	return domain.Interactive{
		ReadCnt:    dao.ReadCnt,
		LikeCnt:    dao.LikeCnt,
		CollectCnt: dao.CollectCnt,
	}
}
