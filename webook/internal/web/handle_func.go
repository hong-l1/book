package web

import (
	"github.com/gin-gonic/gin"
	logger2 "github.com/hong-l1/project/webook/internal/pkg/logger"
	"net/http"
	"time"
)

func WrapHandle[T any](fn func(ctx *gin.Context, req T) (Result, error), l logger2.Loggerv1) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		now := time.Now()
		var req T
		if err := ctx.ShouldBind(&req); err != nil {
			l.Error("参数绑定错误",
				logger2.Field{Key: "path", Value: ctx.FullPath()},
				logger2.Field{Key: "method", Value: ctx.Request.Method},
				logger2.Field{Key: "error", Value: err})
			ctx.JSON(http.StatusBadRequest, Result{
				Code: 4,
				Msg:  "系统异常",
			})
			return
		}
		res, err := fn(ctx, req)
		cost := time.Since(now)
		if err != nil {
			l.Error("业务处理失败",
				logger2.Field{Key: "path", Value: ctx.FullPath()},
				logger2.Field{Key: "method", Value: ctx.Request.Method},
				logger2.Field{Key: "req", Value: req},
				logger2.Field{Key: "error", Value: err},
				logger2.Field{Key: "cost", Value: cost})
			ctx.JSON(http.StatusInternalServerError, Result{
				Code: 5,
				Msg:  "请求失败",
			})
			return
		}
		l.Info("请求成功",
			logger2.Field{Key: "path", Value: ctx.FullPath()},
			logger2.Field{Key: "method", Value: ctx.Request.Method},
			logger2.Field{Key: "req", Value: req},
			logger2.Field{Key: "cost", Value: cost})
		ctx.JSON(http.StatusOK, res)
		return
	}
}
