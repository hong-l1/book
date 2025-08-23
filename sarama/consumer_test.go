package sarama

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	start := time.Now()
	cfg := sarama.NewConfig()
	g, err := sarama.NewConsumerGroup(addres, "test_topic", cfg)
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
	for msg := range msgs {
		fmt.Println(string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}
