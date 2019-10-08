package msgqueue

import (
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

var _ MsgWriter = kafkaMsgWriter{}

type kafkaMsgWriter struct {
	sarama.SyncProducer
}

func NewKafkaMsgWriter(brokers string) (MsgWriter, error) {
	bs := strings.Split(brokers, ",")
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	producer, err := sarama.NewSyncProducer(bs, config)
	return kafkaMsgWriter{producer}, err
}

func (w kafkaMsgWriter) WriteKV(k, v []byte) error {
	_, _, err := w.SyncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: KafkaPubTopic,
		Key:   sarama.ByteEncoder(k),
		Value: sarama.ByteEncoder(v),
	})
	return err
}

func (w kafkaMsgWriter) Close() error {
	return w.SyncProducer.Close()
}

func (w kafkaMsgWriter) String() string {
	return "kafka"
}
