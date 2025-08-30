package saramax

import (
	"encoding/json"
	"github.com/IBM/sarama"
	logger2 "github.com/hong-l1/project/webook/internal/pkg/logger"
)

type Handler[T any] struct {
	l  logger2.Loggerv1
	fn func(msg *sarama.ConsumerMessage, t T) error
}

func (h Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}
func (h Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}
func NewHandler[T any](l logger2.Loggerv1, fn func(msg *sarama.ConsumerMessage, t T) error) *Handler[T] {
	return &Handler[T]{
		l:  l,
		fn: fn,
	}
}
func (h Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			h.l.Error("反序列化失败",
				logger2.Int32("partition", msg.Partition),
				logger2.Int64("offset", msg.Offset),
				logger2.String("topic", msg.Topic),
				logger2.Error(err))
			continue
		}

		for i := 0; i < 3; i++ {
			err = h.fn(msg, t)
			if err == nil {
				break
			}
			h.l.Error("处理消息失败",
				logger2.Int32("partition", msg.Partition),
				logger2.Int64("offset", msg.Offset),
				logger2.String("topic", msg.Topic),
				logger2.Error(err))
		}
		if err != nil {
			h.l.Error("处理消息失败-重试上限",
				logger2.Int32("partition", msg.Partition),
				logger2.Int64("offset", msg.Offset),
				logger2.String("topic", msg.Topic),
				logger2.Error(err))
		} else {
			session.MarkMessage(msg, "")
		}
	}
	return nil
}
