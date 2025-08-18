package article

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
	Upsert(ctx context.Context, art PublishArticleDAO) error
	SyncStatus(ctx context.Context, articleId int64, AuthorId int64, status uint8) error
	GetbyAuthor(ctx context.Context, id int64, offset int, limit int) ([]Article, error)
}
type GORMArticleDao struct {
	db *gorm.DB
}

func (g *GORMArticleDao) GetbyAuthor(ctx context.Context, id int64, offset int, limit int) ([]Article, error) {
	//id是作者id
	arts := []Article{}
	err := g.db.WithContext(ctx).Model(&Article{}).
		Where("author_id = ?", id).
		Offset(offset).
		Limit(limit).
		Order("utime DESC").
		Find(&arts).Error
	return arts, err
}

func (g *GORMArticleDao) SyncStatus(ctx context.Context, articleId int64, AuthorId int64, status uint8) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).Where("id = ? AND author_id = ?", articleId, AuthorId).Updates(map[string]interface{}{
			"status": status,
			"utime":  now,
		})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return fmt.Errorf("误操作非自己的文章 uid:%d,author_id:%d", articleId, AuthorId)
		}
		return tx.Model(&Article{}).Where("id = ?", articleId).Updates(map[string]interface{}{
			"status": status,
			"utime":  now,
		}).Error
	})
}

func (g *GORMArticleDao) Upsert(ctx context.Context, art PublishArticleDAO) error {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := g.db.Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   art.Title,
			"content": art.Content,
			"utime":   now,
			"status":  art.Status,
		}),
	}).Create(&art).Error
	return err
}
func (g *GORMArticleDao) Sync(ctx context.Context, art Article) (int64, error) {
	var id = art.Id
	err := g.db.Transaction(func(tx *gorm.DB) error {
		var err error
		txdao := NewGORMArticleDao(tx)
		if id > 0 {
			err = txdao.UpdateById(ctx, art)
		} else {
			id, err = txdao.Insert(ctx, art)
		}
		if err != nil {
			return err
		}
		return txdao.Upsert(ctx, PublishArticleDAO{Article: art})
	})
	return id, err
}
func (g *GORMArticleDao) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.Utime = now
	res := g.db.WithContext(ctx).Model(&art).
		Where("id = ? AND author_id = ?", art.Id, art.AuthorId).Updates(map[string]any{
		"title":   art.Title,
		"content": art.Content,
		"utime":   art.Utime,
		"status":  art.Status,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("ID 不对或者创作者不对 id %d,author_id ", art.Id, art.AuthorId)
	}
	return nil
}
func (g *GORMArticleDao) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := g.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func NewGORMArticleDao(db *gorm.DB) ArticleDao {
	return &GORMArticleDao{
		db: db,
	}
}

type Article struct {
	Id       int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Title    string `gorm:"type=varchar(4096)" bson:"title,omitempty"`
	Content  string `gorm:"type=BLOB" bson:"content,omitempty"`
	AuthorId int64  `gorm:"index" bson:"author_id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}
