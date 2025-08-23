package service

import (
	"context"
	"github.com/hong-l1/project/webook/internal/repository/article"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	Like(ctx context.Context, biz string, bizId int64, uid int64) error
	CancelLike(ctx context.Context, biz string, id int64, uid int64) error
	Collect(ctx context.Context, biz string, bizId int64, cid int64, uid int64) error
}
type InteractiveServiceImpl struct {
	repo article.InteractiveRepository
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
