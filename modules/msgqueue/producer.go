package msgqueue

import (
	"context"
	"log"
	"strings"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

const (
	brokers       = "brokers"
	topics        = "subscribe-modules"
	PubTopic      = "coinex-dex"
	FeatureToggle = "feature-toggle"
)

type MsgSender interface {
	SendMsg(key []byte, v []byte)
	IsSubScribe(topic string) bool
	IsOpenToggle() bool
}

type Producer struct {
	toggle    bool
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

func (p *Producer) setParam(data config) {
	if len(data.Brokers) == 0 || len(data.Topics) == 0 {
		return
	}
	brokers := strings.Split(data.Brokers, ",")
	topics := strings.Split(data.Topics, ",")

	for _, topic := range topics {
		p.subTopics[topic] = struct{}{}
	}
	p.Writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   PubTopic,
		Async:   true,
	})

	p.toggle = viper.GetBool(FeatureToggle)
}

func (p Producer) close() {
	if err := p.Close(); err != nil {
		log.Fatalln(err)
	}
}

func (p Producer) SendMsg(key []byte, v []byte) {
	p.WriteMessages(context.Background(), kafka.Message{
		Key:   key,
		Value: v,
	})
}

func (p Producer) IsSubScribe(topic string) bool {
	_, ok := p.subTopics[topic]
	return ok
}

func (p Producer) IsOpenToggle() bool {
	return p.toggle
}
