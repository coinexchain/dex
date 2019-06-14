package msgqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type Producer struct {
	topicWrites map[string]*kafka.Writer
	brokers     []string
}

const (
	brokers = "brokers"
	topics  = "topics"
)

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
	p.SetParam(data)

	return p
}

func (k *Producer) SetParam(data config) {
	k.brokers = strings.Split(data.Brokers, ",")
	topics := strings.Split(data.Topics, ",")
	fmt.Println(k.brokers, len(k.brokers))
	fmt.Println(topics, len(topics))
	if len(k.brokers) <= 1 || len(topics) <= 1 {
		return
	}

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
