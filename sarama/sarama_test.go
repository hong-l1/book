package sarama

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"testing"
)

var addres = []string{"127.0.0.1:9094"}

func TestSyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Partitioner = sarama.NewCustomPartitioner()
	producer, err := sarama.NewSyncProducer(addres, cfg)
	assert.NoError(t, err)
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.StringEncoder("hello,这是一条消息"),
	})
	assert.NoError(t, err)
}
func TestASyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Errors = true
	producer, err := sarama.NewAsyncProducer(addres, cfg)
	assert.NoError(t, err)
	msgch := producer.Input()
	msgch <- &sarama.ProducerMessage{
		Topic: "test_topic",
		Key:   sarama.StringEncoder("id-123"),
		Value: sarama.StringEncoder("hello,这是一条消息"),
	}
	errch := producer.Errors()
	successch := producer.Successes()
	select {
	case err := <-errch:
		t.Log("发送出了问题", err.Err)
	case msg := <-successch:
		t.Log("发送成功", msg.Value)
	}
}
