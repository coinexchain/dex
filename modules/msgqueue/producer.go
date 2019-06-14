package msgqueue

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	topicWrites map[string]*kafka.Writer
	brokers     []string
}

func NewProducer() Producer {
	return Producer{
		topicWrites: make(map[string]*kafka.Writer),
	}
}

func (k *Producer) SetParam(data GenesisState) {
	k.brokers = strings.Split(data.Brokers, ",")
	topics := strings.Split(data.Topics, ",")
	for _, topic := range topics {
		k.topicWrites[topic] = nil
	}

	if len(k.brokers) > 0 && len(k.topicWrites) > 0 {
		for topic := range k.topicWrites {
			k.topicWrites[topic] = kafka.NewWriter(kafka.WriterConfig{
				Brokers: k.brokers,
				Topic:   topic,
				Async:   true,
			})
		}
	}
}

func (k Producer) GetParam() GenesisState {
	var index int
	values := make([]string, len(k.topicWrites))
	for topic := range k.topicWrites {
		values[index] = topic
		index++
	}

	return GenesisState{
		Topics:  strings.Join(values, ","),
		Brokers: strings.Join(k.brokers, ","),
	}
}

func (k Producer) close() {
	for _, w := range k.topicWrites {
		if err := w.Close(); err != nil {
			log.Fatalln(err)
		}
	}
}

func (k Producer) IsPublishTopic(topic string) bool {
	_, ok := k.topicWrites[topic]
	return ok
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
