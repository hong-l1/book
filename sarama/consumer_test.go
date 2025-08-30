package sarama

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"log"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	start := time.Now()
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	g, err := sarama.NewConsumerGroup(addres, "test_group", cfg)
	assert.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = g.Consume(ctx, []string{"test_topic"}, testConsumerHandler{})
	t.Log(err, time.Since(start))
}

type testConsumerHandler struct {
}

func (t testConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	fmt.Println("testConsumerHandler.Setup")
	partions := session.Claims()["test_topic"]
	for _, part := range partions {
		session.ResetOffset("test_topic", part, sarama.OffsetOldest, "")
	}
	return nil
}

func (t testConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	fmt.Println("testConsumerHandler.Cleanup")
	return nil
}
func (t testConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	const batchSize = 10
	for {
		batch := make([]*sarama.ConsumerMessage, 0, batchSize)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		var done = false
		var eg errgroup.Group
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				done = true
			case msg, ok := <-msgs:
				if !ok {
					cancel()
					return nil
				}
				batch = append(batch, msg)
				eg.Go(func() error {
					log.Println(string(msg.Value))
					return nil
				})
			}
		}
		cancel()
		err := eg.Wait()
		if err != nil {
			continue
		}
		for _, msg := range batch {
			session.MarkMessage(msg, "")
		}
	}
}
func (t testConsumerHandler) ConsumeClaimV1(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		fmt.Println(string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}
