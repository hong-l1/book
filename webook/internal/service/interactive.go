package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hong-l1/project/webook/internal/repository/article"
)

type InteractiveService interface {
	IncrReadCnt(ctx *gin.Context, biz string, bizId int64) error
}
type InteractiveServiceImpl struct {
	repo article.InteractiveRepository
}

func NewInteractiveServiceImpl(repo article.InteractiveRepository) *InteractiveServiceImpl {
	return &InteractiveServiceImpl{repo: repo}
}
func (i *InteractiveServiceImpl) IncrReadCnt(ctx *gin.Context, biz string, bizId int64) error {
	return i.repo.IncrReadCnt(ctx, biz, bizId)
}
