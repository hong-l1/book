package article

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

type Producer interface {
	ProducerReadEvent(evt ReadEvent) error
}
type ReadEvent struct {
	Uid int64
	Aid int64
}
type KafkaProducer struct {
	SyncProducer sarama.SyncProducer
}

func (k *KafkaProducer) ProducerReadEvent(evt ReadEvent) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = k.SyncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: "read_article",
		Value: sarama.StringEncoder(data),
	})
	return err
}

func NewKafkaProducer(SyncProducer sarama.SyncProducer) Producer {
	return &KafkaProducer{
		SyncProducer: SyncProducer,
	}

}
