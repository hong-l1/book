package web

import (
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"github.com/hong-l1/project/webook/internal/domain"
	logger2 "github.com/hong-l1/project/webook/internal/pkg/logger"
	"github.com/hong-l1/project/webook/internal/pkg/wrapper"
	"github.com/hong-l1/project/webook/internal/service"
	ijwt "github.com/hong-l1/project/webook/internal/web/jwt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"time"
)

var _ Handler = (*ArticleHandle)(nil)

type ArticleHandle struct {
	svc     service.ArticleService
	l       logger2.Loggerv1
	intrsvc service.InteractiveService
	biz     string
}

func NewArticleHandle(l logger2.Loggerv1, svc service.ArticleService, intrsvc service.InteractiveService, biz string) *ArticleHandle {
	return &ArticleHandle{
		l:       l,
		svc:     svc,
		intrsvc: intrsvc,
		biz:     biz,
	}
}
func (u *ArticleHandle) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", u.Edit)
	g.POST("/withdraw", u.Withdraw)
	g.POST("/publish", u.Publish)
	g.POST("/list", wrapper.WrapBodyAndToken[ListReq, ijwt.Claim](u.List))
	g.GET("/detail/:id", wrapper.WrapToken[ijwt.Claim](u.Detail))
	pub := g.Group("/pub")
	pub.GET("/:id", wrapper.WrapToken[ijwt.Claim](u.PubDetail))
	pub.POST("/like", wrapper.WrapBodyAndToken[LikeReq, ijwt.Claim](u.Like))
	pub.POST("/collect", wrapper.WrapBodyAndToken[CollectReq, ijwt.Claim](u.Collect))
}
func (a *ArticleHandle) Withdraw(ctx *gin.Context) {
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	c := ctx.MustGet("claim")
	claim, ok := c.(*ijwt.Claim)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("未发现用户的session信息")
		return
	}
	err := a.svc.Withdraw(ctx, domain.Article{
		Id: req.Id,
		Author: domain.Author{
			Id: claim.UserId,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("撤销帖子失败", logger2.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "ok",
	})
	return
}
func (a *ArticleHandle) Edit(ctx *gin.Context) {
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//跳过输入检测
	c := ctx.MustGet("claim")
	claim, ok := c.(*ijwt.Claim)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("未发现用户的session信息")
		return
	}
	id, err := a.svc.Save(ctx, req.todomain(claim.UserId))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("保存帖子失败", logger2.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "ok",
		Data: id,
	})
	return
}
func (a *ArticleHandle) Publish(ctx *gin.Context) {
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	c := ctx.MustGet("claim")
	claim, ok := c.(*ijwt.Claim)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("未发现用户的session信息")
		return
	}
	id, err := a.svc.Publish(ctx, req.todomain(claim.UserId))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("发表帖子失败", logger2.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "ok",
		Data: id,
	})
	return
}
func (u *ArticleHandle) List(ctx *gin.Context, req ListReq, uc ijwt.Claim) (wrapper.Result, error) {
	res, err := u.svc.List(ctx, req.Offset, req.Limit, uc.UserId)
	if err != nil {
		return Result{Msg: "系统错误", Code: 5}, nil
	}
	return Result{
		Data: slice.Map[domain.Article, ArticleVO](res, func(idx int, src domain.Article) ArticleVO {
			return ArticleVO{
				Id:       src.Id,
				Title:    src.Title,
				Abstract: src.Abstract(),
				Status:   src.Status.ToUint8(),
				Ctime:    src.Ctime.Format(time.DateTime),
				Utime:    src.Utime.Format(time.DateTime),
			}
		}),
	}, nil
}
func (u *ArticleHandle) Detail(ctx *gin.Context, uc ijwt.Claim) (wrapper.Result, error) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		return Result{
			Code: 4,
			Msg:  "参数错误",
		}, err
	}
	art, err := u.svc.GetById(ctx, id)
	if err != nil {
		return Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	if art.Author.Id != uc.UserId {
		return Result{
			Code: 4,
			Msg:  "输入错误",
		}, fmt.Errorf("非法访问文章，创作者 ID 不匹配 %d", id)
	}
	return Result{
		Data: ArticleVO{
			Id:      art.Id,
			Title:   art.Title,
			Content: art.Content,
			Ctime:   art.Ctime.Format(time.DateTime),
			Utime:   art.Utime.Format(time.DateTime),
			Status:  art.Status.ToUint8(),
		},
	}, nil
}
func (u *ArticleHandle) PubDetail(ctx *gin.Context, uc ijwt.Claim) (Result, error) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		return Result{
			Code: 4,
			Msg:  "参数错误",
		}, err
	}
	var art domain.Article
	var eg errgroup.Group
	eg.Go(func() error {
		art, err = u.svc.GetPublishedById(ctx, id)
		return err
	})
	var inter domain.Interactive
	eg.Go(func() error {
		inter, err = u.intrsvc.Get(ctx, u.biz, id, uc.UserId)
		return err
	})
	err = eg.Wait()
	if err != nil {
		return Result{
			Msg:  "系统错误",
			Code: 5,
		}, err
	}
	go func() {
		er := u.intrsvc.IncrReadCnt(ctx, u.biz, art.Id)
		if er != nil {
			u.l.Error("增加阅读计数失败",
				logger2.Int64("article id", art.Id),
				logger2.Error(er))
		}
	}()
	return Result{
		Data: ArticleVO{
			Id:         art.Id,
			Title:      art.Title,
			Content:    art.Content,
			Ctime:      art.Ctime.Format(time.DateTime),
			Utime:      art.Utime.Format(time.DateTime),
			Status:     art.Status.ToUint8(),
			AuthorId:   art.Author.Id,
			AuthorName: art.Author.Name,
			ReadCnt:    inter.ReadCnt,
			LikeCnt:    inter.LikeCnt,
			CollectCnt: inter.CollectCnt,
			Liked:      inter.Liked,
			Collected:  inter.Collected,
		},
	}, nil
}
func (u *ArticleHandle) Like(ctx *gin.Context, req LikeReq, uc ijwt.Claim) (Result, error) {
	var err error
	if req.Like {
		err = u.intrsvc.Like(ctx, u.biz, req.Id, uc.UserId)
	} else {
		err = u.intrsvc.CancelLike(ctx, u.biz, req.Id, uc.UserId)
	}
	if err != nil {
		return Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return Result{Msg: "OK"}, nil
}
func (u *ArticleHandle) Collect(context *gin.Context, req CollectReq, uc ijwt.Claim) (Result, error) {
	panic("implement me")
}
