package msgqueue

import (
	"context"
	"log"
	"strings"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

const (
	brokers  = "brokers"
	topics   = "topics"
	PubTopic = "coinex-dex"
)

type MsgSender interface {
	SendMsg(key []byte, v []byte)
	IsSubScribe(topic string) bool
}

type Producer struct {
	subTopics map[string]struct{}
	*kafka.Writer
}

type config struct {
	Brokers string
	Topics  string
}

func NewProducer() Producer {
	p := Producer{
		subTopics: make(map[string]struct{}),
	}

	p.setParam(config{
		Brokers: viper.GetString(brokers),
		Topics:  viper.GetString(topics),
	})
	return p
}

func (k *Producer) setParam(data config) {
	if len(data.Brokers) == 0 || len(data.Topics) == 0 {
		return
	}
	brokers := strings.Split(data.Brokers, ",")
	topics := strings.Split(data.Topics, ",")

	for _, topic := range topics {
		k.subTopics[topic] = struct{}{}
	}
	k.Writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   PubTopic,
		Async:   true,
	})
}

func (k Producer) close() {
	if err := k.Close(); err != nil {
		log.Fatalln(err)
	}
}

func (k Producer) SendMsg(key []byte, v []byte) {
	k.WriteMessages(context.Background(), kafka.Message{
		Key:   key,
		Value: v,
	})
}

func (k Producer) IsSubScribe(topic string) bool {
	_, ok := k.subTopics[topic]
	return ok
}
