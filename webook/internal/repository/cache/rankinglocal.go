package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/hong-l1/project/webook/internal/domain"
	"go.uber.org/atomic"
	"time"
)

type RankingLocalCache struct {
	topN       atomic.String
	ddl        atomic.Time
	expiration atomic.Duration
}

func (r *RankingLocalCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	data, err := json.Marshal(val)
	r.ddl.Store(time.Now().Add(expiration))
	r.topN.Store(string(data))
	return err
}

func (r *RankingLocalCache) Get(ctx context.Context, key string) ([]domain.Article, error) {
	arts := r.topN.Load()
	if len(arts) == 0 || r.ddl.Load().Before(time.Now()) {
		return nil, errors.New("cache expired")
	}
	var data []domain.Article
	err := json.Unmarshal([]byte(arts), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *RankingLocalCache) Foreget(ctx context.Context) ([]domain.Article, error) {
	panic("implement me")
}
