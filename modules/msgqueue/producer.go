package msgqueue

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

const (
	brokers = "brokers"
	topics  = "topics"
)

type MsgSender interface {
	SendMsg(topic string, key string, v interface{}) error
}

type Producer struct {
	topicWrites map[string]*kafka.Writer
	brokers     []string
}

type config struct {
	Brokers string
	Topics  string
}

func NewProducer() Producer {
	p := Producer{
		topicWrites: make(map[string]*kafka.Writer),
	}

	data := config{
		Brokers: viper.GetString(brokers),
		Topics:  viper.GetString(topics),
	}
	p.setParam(data)

	return p
}

func (k *Producer) setParam(data config) {
	if len(data.Brokers) == 0 || len(data.Topics) == 0 {
		return
	}
	k.brokers = strings.Split(data.Brokers, ",")
	topics := strings.Split(data.Topics, ",")

	for _, topic := range topics {
		k.topicWrites[topic] = kafka.NewWriter(kafka.WriterConfig{
			Brokers: k.brokers,
			Topic:   topic,
			Async:   true,
		})
	}
}

func (k Producer) close() {
	for _, w := range k.topicWrites {
		if err := w.Close(); err != nil {
			log.Fatalln(err)
		}
	}
}

func (k Producer) SendMsg(topic string, key string, v interface{}) error {
	if w, ok := k.topicWrites[topic]; ok {
		bytes, err := json.Marshal(v)
		if err != nil {
			return err
		}

		w.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(key),
			Value: bytes,
		})
	}

	return nil
}
