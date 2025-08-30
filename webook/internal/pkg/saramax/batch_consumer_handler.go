package saramax

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	logger2 "github.com/hong-l1/project/webook/internal/pkg/logger"
	"time"
)

type BatchConsumerHandle[T any] struct {
	l  logger2.Loggerv1
	fn func(msg []*sarama.ConsumerMessage, t []T) error
}

func (b BatchConsumerHandle[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b BatchConsumerHandle[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}
func NewBatchConsumerHandle[T any](l logger2.Loggerv1, fn func(msg []*sarama.ConsumerMessage, t []T) error) BatchConsumerHandle[T] {
	return BatchConsumerHandle[T]{
		l:  l,
		fn: fn,
	}
}
func (b BatchConsumerHandle[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	const batchSize = 10
	for {
		batch := make([]*sarama.ConsumerMessage, 0, batchSize)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		ts := make([]T, 0, batchSize)
		var done = false
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				done = true
			case msg, ok := <-msgs:
				if !ok {
					cancel()
					return nil
				}
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					b.l.Error("反序列化失败",
						logger2.String("topic", msg.Topic),
						logger2.Int32("partition", msg.Partition),
						logger2.Int64("offset", msg.Offset),
						logger2.Error(err))
					continue
				}
				batch = append(batch, msg)
				ts = append(ts, t)
			}
		}
		cancel()
		err := b.fn(batch, ts)
		if err != nil {
			b.l.Error("处理消息失败",
				logger2.Error(err))
		}
		for _, msg := range batch {
			session.MarkMessage(msg, "")
		}
	}
}
