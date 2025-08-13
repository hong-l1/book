package article

import (
	"context"
	"gorm.io/gorm"
)

type ReaderDao interface {
	Upsert(ctx context.Context, art Article) error
}

func NewReaderDAO(db *gorm.DB) ReaderDao {
	panic("implement me")
}

type PublishArticleDAO struct {
	Article
}
