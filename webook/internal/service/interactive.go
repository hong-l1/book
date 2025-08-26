package service

import (
	"context"
	"github.com/hong-l1/project/webook/internal/domain"
	"github.com/hong-l1/project/webook/internal/repository/article"
	"golang.org/x/sync/errgroup"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	Like(ctx context.Context, biz string, bizId int64, uid int64) error
	CancelLike(ctx context.Context, biz string, id int64, uid int64) error
	Collect(ctx context.Context, biz string, bizId int64, cid int64, uid int64) error
	Get(ctx context.Context, biz string, artid int64, uid int64) (domain.Interactive, error)
}
type InteractiveServiceImpl struct {
	repo article.InteractiveRepository
}

func (i *InteractiveServiceImpl) Get(ctx context.Context, biz string, artid int64, uid int64) (domain.Interactive, error) {
	var (
		err       error
		inter     domain.Interactive
		liked     bool
		collected bool
	)
	var eg errgroup.Group
	eg.Go(func() error {
		var err error
		inter, err = i.repo.Get(ctx, biz, artid)
		return err
	})
	eg.Go(func() error {
		var err error
		liked, err = i.repo.Liked(ctx, biz, artid, uid)
		return err
	})
	eg.Go(func() error {
		var err error
		collected, err = i.repo.Collected(ctx, biz, artid, uid)
		return err
	})
	err = eg.Wait()
	if err != nil {
		return domain.Interactive{}, err
	}
	inter.Liked = liked
	inter.Collected = collected
	return inter, nil
}

func (i *InteractiveServiceImpl) Collect(ctx context.Context, biz string, bizId int64, cid int64, uid int64) error {
	return i.repo.AddCollectItem(ctx, biz, bizId, cid, uid)
}

func NewInteractiveServiceImpl(repo article.InteractiveRepository) InteractiveService {
	return &InteractiveServiceImpl{repo: repo}
}
func (i *InteractiveServiceImpl) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return i.repo.IncrReadCnt(ctx, biz, bizId)
}
func (i *InteractiveServiceImpl) Like(ctx context.Context, biz string, bizId int64, uid int64) error {
	return i.repo.IncrLike(ctx, biz, bizId, uid)
}

func (i *InteractiveServiceImpl) CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	return i.repo.DecrLike(ctx, biz, bizId, uid)
}
