package ioc

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hong-l1/project/webook/internal/pkg/ginx/middlewares/logger"
	logger2 "github.com/hong-l1/project/webook/internal/pkg/logger"
	"github.com/hong-l1/project/webook/internal/web"
	"strings"
	"time"
)

func InitGin(mdls []gin.HandlerFunc, hdl *web.UserHandle,
	OAuth2WeChatHandle *web.OAuth2WeChatHandle, article *web.ArticleHandle) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	hdl.RegisterUsersRoutes(server)
	article.RegisterRoutes(server)
	OAuth2WeChatHandle.RegisterRoutes(server)
	return server
}
func InitMiddlewares(l logger2.Loggerv1) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			AllowMethods:     []string{"PUT", "POST"},
			AllowHeaders:     []string{"content-type", "authorization"},
			ExposeHeaders:    []string{"jwt-token", "refresh-token"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					return true
				}
				return strings.Contains(origin, "company.com")
			},
			MaxAge: 12 * time.Second,
		}),
		logger.NewBuilder(func(ctx context.Context, al *logger.AccessLogger) {
			l.Debug("HTTP请求", logger2.Field{
				Key:   "al",
				Value: al,
			})
		}).AllowRespBody().AllowReBody().Build(),
	}
}
