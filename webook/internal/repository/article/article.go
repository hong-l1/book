package article

import (
	"context"
	"github.com/hong-l1/project/webook/internal/domain"
	"github.com/hong-l1/project/webook/internal/repository/dao/article"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Syncv1(ctx context.Context, art domain.Article) (int64, error)
}
type CacheArticle struct {
	dao    article.ArticleDao
	reader article.ReaderDao
	author article.AuthorDao
	db     *gorm.DB
}

func (r *CacheArticle) Syncv2(ctx context.Context, art domain.Article) (int64, error) {

}
func (r *CacheArticle) Syncv1(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	templ := r.toEntity(art)
	if art.Id > 0 {
		err = r.author.UpdateById(ctx, templ)
	} else {
		id, err = r.author.Insert(ctx, templ)
	}
	if err != nil {
		return id, err
	}
	err = r.reader.Upsert(ctx, templ)
	return id, err
}

func (r *CacheArticle) Create(ctx context.Context, art domain.Article) (int64, error) {
	return r.dao.Insert(ctx, article.Article{
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}
func NewCacheArticle(dao article.ArticleDao) ArticleRepository {
	return &CacheArticle{
		dao: dao,
	}
}
func (r *CacheArticle) Update(ctx context.Context, art domain.Article) error {
	return r.dao.UpdateById(ctx, article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}
func (r *CacheArticle) toEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}
