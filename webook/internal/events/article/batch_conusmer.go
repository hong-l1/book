package article

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/hong-l1/project/webook/internal/pkg/logger"
	"github.com/hong-l1/project/webook/internal/pkg/saramax"
	"github.com/hong-l1/project/webook/internal/repository/article"
	"time"
)

type BatchConusmer struct {
	l      logger.Loggerv1
	repo   article.InteractiveRepository
	client sarama.Client
}

func (c *BatchConusmer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", c.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(), []string{"article_read"}, saramax.NewBatchConsumerHandle[ReadEvent](c.l, c.Consume))
		if err != nil {
			c.l.Error("退出消费循环异常", logger.Error(err))
		}
	}()
	return err
}
func (c *BatchConusmer) Consume(msg []*sarama.ConsumerMessage, t []ReadEvent) error {
	ids := make([]int64, 0, len(t))
	bizs := make([]string, 0, len(t))
	for _, evt := range t {
		ids = append(ids, evt.Aid)
		bizs = append(bizs, "article")
	}
	ctx, canel := context.WithTimeout(context.Background(), time.Second*1)
	defer canel()
	err := c.repo.BatchIncrReadCnt(ctx, bizs, ids)
	if err != nil {
		c.l.Error("<UNK>",
			logger.Field{Key: "ids", Value: ids},
			logger.Error(err))
	}
	return nil
}
func NewBatchConusmer(l logger.Loggerv1, repo article.InteractiveRepository, client sarama.Client) *BatchConusmer {
	return &BatchConusmer{
		l:      l,
		repo:   repo,
		client: client,
	}
}
