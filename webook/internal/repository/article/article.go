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
	//存储并同步
	Syncv1(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, articleId int64, AuthorId int64, Status domain.ArticleStatus) error
}
type CacheArticle struct {
	dao    article.ArticleDao
	reader article.ReaderDao
	author article.AuthorDao
	db     *gorm.DB
}

func (r *CacheArticle) SyncStatus(ctx context.Context, articleId int64, AuthorId int64, Status domain.ArticleStatus) error {
	return r.dao.SyncStatus(ctx, articleId, AuthorId, Status.ToUint8())
}

func (r *CacheArticle) Sync(ctx context.Context, art domain.Article) (int64, error) {
	return r.dao.Sync(ctx, r.toEntity(art))
}
func (r *CacheArticle) Syncv2(ctx context.Context, art domain.Article) (int64, error) {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return 0, err
	}
	defer tx.Rollback()
	author := article.NewAuthorDAO(tx)
	reader := article.NewReaderDAO(tx)
	var (
		id  = art.Id
		err error
	)
	templ := r.toEntity(art)
	if art.Id > 0 {
		err = author.UpdateById(ctx, templ)
	} else {
		id, err = author.Insert(ctx, templ)
	}
	if err != nil {

		return id, err
	}
	err = reader.Upsert(ctx, templ)
	tx.Commit()
	return id, err
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
	return r.dao.Insert(ctx, r.toEntity(art))
}
func NewCacheArticle(dao article.ArticleDao) ArticleRepository {
	return &CacheArticle{
		dao: dao,
	}
}
func (r *CacheArticle) Update(ctx context.Context, art domain.Article) error {
	return r.dao.UpdateById(ctx, r.toEntity(art))
}
func (r *CacheArticle) toEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	}
}
