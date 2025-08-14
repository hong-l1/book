package web

import (
	"github.com/gin-gonic/gin"
	"github.com/hong-l1/project/webook/internal/domain"
	logger2 "github.com/hong-l1/project/webook/internal/pkg/logger"
	"github.com/hong-l1/project/webook/internal/service"
	ijwt "github.com/hong-l1/project/webook/internal/web/jwt"
	"net/http"
)

var _ Handler = (*ArticleHandle)(nil)

type ArticleHandle struct {
	svc service.ArticleService
	l   logger2.Loggerv1
}

func NewArticleHandle(l logger2.Loggerv1, svc service.ArticleService) *ArticleHandle {
	return &ArticleHandle{
		l:   l,
		svc: svc,
	}
}
func (u *ArticleHandle) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", u.Edit)
	g.POST("/withdraw", u.Withdraw)
	g.POST("/publish", u.Publish)
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
func (req Req) todomain(userid int64) domain.Article {
	return domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: userid,
		},
	}
}

type Req struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	id      int64
}
