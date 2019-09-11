package msgqueue

import (
	"context"
	"strings"

	"github.com/segmentio/kafka-go"
)

var _ MsgWriter = kafkaMsgWriter{}

type kafkaMsgWriter struct {
	*kafka.Writer
}

func NewKafkaMsgWriter(brokers string) MsgWriter {
	bs := strings.Split(brokers, ",")
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: bs,
		Topic:   KafkaPubTopic,
		Async:   true,
	})
	return kafkaMsgWriter{w}
}

func (w kafkaMsgWriter) WriteKV(k, v []byte) error {
	return w.WriteMessages(context.Background(), kafka.Message{
		Key:   k,
		Value: v,
	})
}

func (w kafkaMsgWriter) Close() error {
	return w.Writer.Close()
}

func (w kafkaMsgWriter) String() string {
	return "kafka"
}
