package wrapper

import (
	"github.com/gin-gonic/gin"
	"github.com/hong-l1/project/webook/internal/pkg/logger"
	ijwt "github.com/hong-l1/project/webook/internal/web/jwt"
	"net/http"
)

var l logger.Logger

func Wrapper[T any](fn func(ctx *gin.Context, req T) (Result, error), l logger.Loggerv1) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		err := c.Bind(req)
		if err != nil {
			return
		}
		result, err := fn(c, req)
		if err != nil {
			l.Error("处理业务逻辑出错",
				logger.Error(err),
				logger.String("route", c.FullPath()))
		}
		c.JSON(http.StatusOK, result)
	}
}
func WrapBodyAndToken[T any, C ijwt.Claim](fn func(ctx *gin.Context, req T, uc C) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		err := ctx.Bind(req)
		if err != nil {
			return
		}
		val, ok := ctx.Get("claim")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c := val.(*C)
		result, err := fn(ctx, req, *c)
		if err != nil {
			l.Error("处理业务逻辑出错",
				logger.Error(err),
				logger.String("route", ctx.FullPath()))
		}
		ctx.JSON(http.StatusOK, result)
	}
}

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
