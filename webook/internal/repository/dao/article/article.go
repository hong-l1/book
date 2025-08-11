package article

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ArticleDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
}
type GORMArticleDao struct {
	db *gorm.DB
}

func (g *GORMArticleDao) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.Utime = now
	res := g.db.WithContext(ctx).Model(&art).
		Where("id = ? AND author_id = ?", art.Id, art.AuthorId).Updates(map[string]any{
		"title":   art.Title,
		"content": art.Content,
		"utime":   art.Utime,
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
	Id       int64  `gorm:"primary_key;auto_increment"`
	Title    string `gorm:"type=varchar(1024)"`
	Content  string `gorm:"type:BLOB"`
	AuthorId int64  `gorm:"index"`
	//Author  int64  `gorm:"index=aid_ctime"`
	//Ctime   int64  `gorm:"index=aid_ctime"`
	Ctime int64
	Utime int64
}
