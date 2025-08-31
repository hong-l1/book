package service

import (
	"context"
	"github.com/hong-l1/project/webook/internal/domain"
	event "github.com/hong-l1/project/webook/internal/events/article"
	"github.com/hong-l1/project/webook/internal/pkg/logger"
	"github.com/hong-l1/project/webook/internal/repository/article"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
	Publishv1(ctx context.Context, article domain.Article) (int64, error)
	Withdraw(ctx context.Context, article domain.Article) error
	List(ctx context.Context, offset int, limit int, id int64) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPublishedById(ctx context.Context, artid, uid int64) (domain.Article, error)
}
type ServiceArticle struct {
	repo     article.ArticleRepository
	author   article.ArticleAuthorRepository
	reader   article.ArticleReaderRepository
	l        logger.Loggerv1
	producer event.Producer
}

func (s *ServiceArticle) GetPublishedById(ctx context.Context, artid, uid int64) (domain.Article, error) {
	art, err := s.repo.GetPublishedById(ctx, artid)
	if err == nil {
		go func() {
			er := s.producer.ProducerReadEvent(event.ReadEvent{
				Uid: uid,
				Aid: artid,
			})
			if er != nil {
				s.l.Error("发送 ReadEvent 失败",
					logger.Int64("aid", artid),
					logger.Int64("uid", uid),
					logger.Error(err))
			}
		}()
	}
	return art, err
}

func (s *ServiceArticle) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return s.repo.GetById(ctx, id)
}

func (s *ServiceArticle) List(ctx context.Context, offset int, limit int, id int64) ([]domain.Article, error) {
	return s.repo.List(ctx, offset, limit, id)
}

func (s *ServiceArticle) Withdraw(ctx context.Context, article domain.Article) error {
	return s.repo.SyncStatus(ctx, article.Id, article.Author.Id, domain.ArticleStatusPrivate)
}

func NewServiceArticle(repo article.ArticleRepository, l logger.Loggerv1, producer event.Producer) ArticleService {
	return &ServiceArticle{
		repo:     repo,
		l:        l,
		producer: producer,
	}
}
func NewServiceArticlev1(author article.ArticleAuthorRepository,
	reader article.ArticleReaderRepository, l logger.Loggerv1, producer event.Producer) ArticleService {
	return &ServiceArticle{
		author:   author,
		reader:   reader,
		l:        l,
		producer: producer,
	}
}
func (s *ServiceArticle) Save(ctx context.Context, article domain.Article) (int64, error) {
	article.Status = domain.ArticleStatusUnpublished
	if article.Id > 0 {
		err := s.repo.Update(ctx, article)
		return article.Id, err
	}
	return s.repo.Create(ctx, article)
}
func (s *ServiceArticle) Publish(ctx context.Context, article domain.Article) (int64, error) {
	article.Status = domain.ArticleStatusPublished
	return s.repo.Sync(ctx, article)
}
func (s *ServiceArticle) Publishv1(ctx context.Context, article domain.Article) (int64, error) {
	var (
		id  = article.Id
		err error
	)
	if id > 0 {
		err = s.author.Update(ctx, article)
	} else {
		id, err = s.author.Create(ctx, article)
	}
	if err != nil {
		return 0, err
	}
	article.Id = id
	for i := 0; i < 3; i++ {
		id, err = s.reader.Save(ctx, article)
		if err == nil {
			break
		}
		s.l.Error("部分失败，保存到线上库失败",
			logger.Getint64("article.id", id),
			logger.Error(err))
	}
	if err != nil {
		s.l.Error("部分失败，重试全部失败",
			logger.Getint64("article.id", id),
			logger.Error(err))
	}
	return id, err
}
