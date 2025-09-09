package saramax

import (
	"encoding/json"
	"github.com/IBM/sarama"
	logger2 "github.com/hong-l1/project/webook/internal/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
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
	//在这里打监控,
	counter := newCounter()
	msgs := claim.Messages()
	for msg := range msgs {
		counter.WithLabelValues(msg.Topic, strconv.Itoa(int(msg.Partition)), strconv.Itoa(int(msg.Offset))).Add(1)
		//time := time2.Now()
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
func newvector() *prometheus.SummaryVec {
	vector := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "gowebook",
			Subsystem: "webook",
			Name:      "kafka消费",
			Help:      "kafka异步消费监控",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		}, []string{})
	prometheus.MustRegister(vector)
	return vector
}
func newCounter() *prometheus.CounterVec {
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gowebook",
		Subsystem: "webook",
		Name:      "consumer_cnt",
		Help:      "kafka消息数量总数",
	}, []string{"topic", "partition", "offset"})
	return counter
}
