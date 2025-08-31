package ioc

import (
	"github.com/IBM/sarama"
	"github.com/hong-l1/project/webook/internal/events/article"
)

var addres = []string{"127.0.0.1:9094"}

func Initkafka() sarama.Client {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	client, err := sarama.NewClient(addres, cfg)
	if err != nil {
		panic(err)
	}
	return client
}
func InitSyncProducer(c sarama.Client) sarama.SyncProducer {
	p, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		panic(err)
	}
	return p
}
func InitConsumers(c *article.BatchConusmer) []article.Consumer {
	return []article.Consumer{c}
}
