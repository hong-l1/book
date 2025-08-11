package article

import (
	"context"
	"github.com/hong-l1/project/webook/internal/domain"
)

type ArticleReaderRepository interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
}
