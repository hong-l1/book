package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hong-l1/project/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, author int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, author int64, arts []domain.Article) error
	DelFirstPage(ctx context.Context, author int64) error
	Set(ctx context.Context, art domain.Article, id int64) error
}
type RedisArticleCache struct {
	client redis.Cmdable
}

func (r RedisArticleCache) GetFirstPage(ctx context.Context, author int64) ([]domain.Article, error) {
	panic("implement me")
}

func (r RedisArticleCache) SetFirstPage(ctx context.Context, author int64, arts []domain.Article) error {
	for k := range arts {
		arts[k].Content = arts[k].Abstract()
	}
	data, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.FirstPageKey(author), data, time.Minute*10).Err()
}
func (r RedisArticleCache) Key(uid int64) string {
	return fmt.Sprintf("article:first_page:%d", uid)
}
func (r RedisArticleCache) FirstPageKey(uid int64) string {
	return fmt.Sprintf("article:first_page:%d", uid)
}
func (r RedisArticleCache) DelFirstPage(ctx context.Context, author int64) error {
	panic("implement me")
}
func (r RedisArticleCache) Set(ctx context.Context, art domain.Article, id int64) error {
	data, err := json.Marshal(art)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.Key(id), data, time.Second*30).Err()
}
