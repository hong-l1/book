//go:build wireinject

package main

import (
	"github.com/google/wire"
	event "github.com/hong-l1/project/webook/internal/events/article"
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

func InitWebServer() *App {
	wire.Build(
		// 第三方依赖
		ioc.InitRedis,
		ioc.InitDb,
		ioc.Initkafka,
		ioc.InitSyncProducer,
		ioc.InitLogger,
		ProvideBizConfig,
		// DAO 部分
		dao.NewUserDao,
		article2.NewGORMArticleDao,
		article2.NewGORMInteractiveDAO,
		event.NewKafkaProducer,
		event.NewBatchConusmer,
		//event.NewKafkaConsumer,
		ioc.InitConsumers,
		// cache 部分
		cache.NewUserCache,
		cache.NewCodeCache,
		//cache.NewRedisInteractiveCache,
		cache.NewRedisArticleCache,
		// repository 部分
		repository.NewCodeRepository,
		repository.NewRepository,
		article.NewCacheArticle,
		article.NewCachedInteractiveRepository,
		// Service 部分
		service.NewUserService,
		service.NewCodeService,
		service.NewServiceArticle,
		ioc.InitSmsService,
		service.NewInteractiveServiceImpl,
		// handler 部分
		web.NewArticleHandle,
		ijwt.NewRedisJWT,
		ioc.InitOauth2WechatService,
		web.NewUserHandle,
		web.NewOAuth2WeChatHandle,
		ioc.InitGin,
		ioc.InitMiddlewares,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}

func ProvideBizConfig() string {
	return "article"
}

//func InitArticleHandle() *web.ArticleHandle {
//	wire.Build(service.NewServiceArticle, web.NewArticleHandle, ioc.InitLogger, article.NewCacheArticle)
//	return &web.ArticleHandle{}
//}
