package msgqueue

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
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

const (
	fileMode  = "file:"
	KafkaMode = "kafka:"
	StdMode   = "os:"
)

type Mode int

const (
	fMode  Mode = 1
	kaMode Mode = 2
	sMode  Mode = 3
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
	io.WriteCloser
	mode Mode
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

	if strings.HasPrefix(data.Brokers, KafkaMode) {
		p.setKafka(strings.TrimPrefix(data.Brokers, KafkaMode))
	} else if strings.HasPrefix(data.Brokers, fileMode) {
		if err := p.setFileWrite(strings.TrimPrefix(data.Brokers, fileMode)); err != nil {
			return
		}
	} else if strings.HasPrefix(data.Brokers, StdMode) {
		if err := p.setStdIO(strings.TrimPrefix(data.Brokers, StdMode)); err != nil {
			return
		}
	} else {
		return
	}

	ts := strings.Split(data.Topics, ",")
	for _, topic := range ts {
		p.subTopics[topic] = struct{}{}
	}
	p.toggle = viper.GetBool(FeatureToggle)
}

func (p *Producer) setStdIO(ioString string) error {
	if !strings.Contains(ioString, "stdout") {
		return fmt.Errorf("Unknow output identifier ")
	}
	p.WriteCloser = os.Stdout
	p.mode = sMode
	return nil
}

func (p *Producer) setFileWrite(filePath string) error {
	if s, err := os.Stat(filePath); !os.IsExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		p.WriteCloser = file
	} else {
		if s.IsDir() {
			return fmt.Errorf("Need to give the file path ")
		}
		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		p.WriteCloser = file
	}
	p.mode = fMode
	return nil
}

func (p *Producer) setKafka(brokers string) {
	bs := strings.Split(brokers, ",")
	p.Writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers: bs,
		Topic:   PubTopic,
		Async:   true,
	})
	p.mode = kaMode
}

func (p *Producer) close() {
	switch p.mode {
	case kaMode:
		if err := p.Writer.Close(); err != nil {
			log.Fatalln(err)
		}
	case fMode:
		if err := p.WriteCloser.Close(); err != nil {
			log.Fatalln(err)
		}
	case sMode:
		return
	}
}

func (p Producer) SendMsg(key []byte, v []byte) {
	switch p.mode {
	case kaMode:
		p.WriteMessages(context.Background(), kafka.Message{
			Key:   key,
			Value: v,
		})
	case sMode, fMode:
		p.WriteCloser.Write(key)
		p.WriteCloser.Write(v)
		p.WriteCloser.Write([]byte("\r\n"))
	}

}

func (p Producer) IsSubScribe(topic string) bool {
	_, ok := p.subTopics[topic]
	return ok
}

func (p Producer) IsOpenToggle() bool {
	return p.toggle
}
