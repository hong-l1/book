//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/hong-l1/project/webook/internal/repository"
	"github.com/hong-l1/project/webook/internal/repository/article"
	"github.com/hong-l1/project/webook/internal/repository/cache"
	"github.com/hong-l1/project/webook/internal/repository/dao"
	article2 "github.com/hong-l1/project/webook/internal/repository/dao/article"
	"github.com/hong-l1/project/webook/internal/service"
	"github.com/hong-l1/project/webook/internal/web"
	ijwt "github.com/hong-l1/project/webook/internal/web/jwt"
	"github.com/hong-l1/project/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(ioc.InitRedis, InitDB,
		dao.NewUserDao,
		cache.NewUserCache,
		cache.NewCodeCache,
		ioc.InitLogger,
		repository.NewCodeRepository,
		repository.NewRepository,
		article2.NewGORMArticleDao,
		article.NewCacheArticle,
		web.NewArticleHandle,
		service.NewUserService,
		service.NewCodeService,
		service.NewServiceArticle,
		ioc.InitSmsService,
		ijwt.NewRedisJWT,
		ioc.InitOauth2WechatService,
		web.NewUserHandle,
		web.NewOAuth2WeChatHandle,
		ioc.InitGin,
		ioc.InitMiddlewares,
	)
	return gin.Default()
}
func InitArticleHandle() *web.ArticleHandle {
	wire.Build(InitDB, service.NewServiceArticle, web.NewArticleHandle, ioc.InitLogger, article.NewCacheArticle, article2.NewGORMArticleDao)
	return &web.ArticleHandle{}
}
