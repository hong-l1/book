package article

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/hong-l1/project/webook/internal/pkg/logger"
	"github.com/hong-l1/project/webook/internal/pkg/saramax"
	"github.com/hong-l1/project/webook/internal/repository/article"
	"time"
)

type Consumer interface {
	Start() error
}
type KafkaConsumer struct {
	l      logger.Loggerv1
	repo   article.InteractiveRepository
	client sarama.Client
}

func (c *KafkaConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", c.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(), []string{"article_read"}, saramax.NewHandler(c.l, c.Consume))
		if err != nil {
			c.l.Error("退出消费循环异常", logger.Error(err))
		}
	}()
	return err
}
func (c *KafkaConsumer) Consume(msg *sarama.ConsumerMessage, t ReadEvent) error {
	ctx, canel := context.WithTimeout(context.Background(), time.Second)
	defer canel()
	return c.repo.IncrReadCnt(ctx, "article", t.Uid)
}
