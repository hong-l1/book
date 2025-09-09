package article

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/hong-l1/project/webook/internal/domain"
	"github.com/hong-l1/project/webook/internal/pkg/logger"
	"github.com/hong-l1/project/webook/internal/repository"
	"github.com/hong-l1/project/webook/internal/repository/cache"
	"github.com/hong-l1/project/webook/internal/repository/dao/article"
	"gorm.io/gorm"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	//Syncv1存储并同步
	Syncv1(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, articleId int64, AuthorId int64, Status domain.ArticleStatus) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
	Syncv2(ctx context.Context, art domain.Article) (int64, error)
	List(ctx context.Context, offset int, limit int, id int64) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	PreCache(ctx context.Context, data []domain.Article)
	GetPublishedById(ctx context.Context, artid int64) (domain.Article, error)
	ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]domain.Article, error)
}
type CacheArticle struct {
	dao article.ArticleDao
	//reader article.ReaderDao
	//author article.AuthorDao
	userrepo repository.UserRepository
	db       *gorm.DB
	cache    cache.ArticleCache
	l        logger.Loggerv1
}

func (r *CacheArticle) ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]domain.Article, error) {
	res, err := r.dao.ListPub(ctx, start, offset, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map(res, func(idx int, dao article.Article) domain.Article {
		return r.todomain(dao)
	}), nil
}

func (r *CacheArticle) GetPublishedById(ctx context.Context, artid int64) (domain.Article, error) {
	art, err := r.dao.GetById(ctx, artid)
	if err != nil {
		return domain.Article{}, err
	}
	user, err := r.userrepo.FindById(ctx, domain.User{Id: artid})
	if err != nil {
		return domain.Article{}, err
	}
	res := domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id:   user.Id,
			Name: user.Nickname,
		},
		Ctime: time.UnixMilli(art.Ctime),
		Utime: time.UnixMilli(art.Utime),
	}
	return res, nil
}

func (r *CacheArticle) GetById(ctx context.Context, id int64) (domain.Article, error) {
	data, err := r.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return r.todomain(data), nil
}

func (r *CacheArticle) List(ctx context.Context, offset int, limit int, id int64) ([]domain.Article, error) {
	if offset == 0 && limit <= 100 {
		data, err := r.cache.GetFirstPage(ctx, id)
		if err == nil {
			go func() {
				r.PreCache(ctx, data)
			}()
			return data, nil
		}
	}
	res, err := r.dao.GetbyAuthor(ctx, id, offset, limit)
	if err != nil {
		return nil, err
	}
	data := slice.Map[article.Article, domain.Article](res, func(idx int, src article.Article) domain.Article {
		return r.todomain(src)
	})
	go func() {
		err := r.cache.SetFirstPage(ctx, id, data)
		r.l.Error("缓存回写失败", logger.Error(err))
		r.PreCache(ctx, data)
	}()
	return data, nil
}
func (r *CacheArticle) SyncStatus(ctx context.Context, articleId int64, AuthorId int64, Status domain.ArticleStatus) error {
	return r.dao.SyncStatus(ctx, articleId, AuthorId, Status.ToUint8())
}
func (r *CacheArticle) Sync(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	templ := r.toEntity(art)
	if art.Id > 0 {
		err = r.dao.UpdateById(ctx, templ)
	} else {
		id, err = r.dao.Insert(ctx, templ)
	}
	if err != nil {
		return id, err
	}
	id, err = r.dao.Sync(ctx, templ)
	return id, err
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
	//var (
	//	id  = art.Id
	//	err error
	//)
	//templ := r.toEntity(art)
	//if art.Id > 0 {
	//	err = r.author.UpdateById(ctx, templ)
	//} else {
	//	id, err = r.author.Insert(ctx, templ)
	//}
	//if err != nil {
	//
	//	return id, err
	//}
	//err = r.reader.Upsert(ctx, templ)
	//return id, err
	panic("implement me")
}
func (r *CacheArticle) Create(ctx context.Context, art domain.Article) (int64, error) {
	defer func() {
		//情况缓存
		err := r.cache.DelFirstPage(ctx, art.Id)
		r.l.Error("缓存删除失败", logger.Error(err))
	}()
	return r.dao.Insert(ctx, r.toEntity(art))
}
func (r *CacheArticle) Update(ctx context.Context, art domain.Article) error {
	defer func() {
		//情况缓存
		err := r.cache.DelFirstPage(ctx, art.Id)
		r.l.Error("缓存删除失败", logger.Error(err))
	}()
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
func (r *CacheArticle) todomain(art article.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
		},
		Ctime: time.UnixMilli(art.Ctime),
		Utime: time.UnixMilli(art.Utime),
	}
}

func (r *CacheArticle) PreCache(ctx context.Context, data []domain.Article) {
	if len(data) > 0 && len(data[0].Content) < 1024*1024 {
		err := r.cache.Set(ctx, data[0], data[0].Id)
		if err != nil {
			r.l.Error("预加载缓存失败", logger.Error(err))
		}
	}
}
func NewCacheArticle(dao article.ArticleDao, db *gorm.DB, cache cache.ArticleCache, l logger.Loggerv1) ArticleRepository {
	return &CacheArticle{
		dao:   dao,
		db:    db,
		cache: cache,
		l:     l,
	}
}
