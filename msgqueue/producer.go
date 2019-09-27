package msgqueue

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	FlagBrokers       = "brokers"
	FlagTopics        = "subscribe-modules"
	FlagFeatureToggle = "feature-toggle"
	KafkaPubTopic     = "coinex-dex"
)

const (
	CfgPrefixFile  = "file:"
	CfgPrefixKafka = "kafka:"
	CfgPrefixOS    = "os:"
)

type MsgSender interface {
	SendMsg(key []byte, v []byte)
	IsSubscribed(topic string) bool
	IsOpenToggle() bool
	GetMode() []string
	Close()
}

type producer struct {
	toggle     bool
	subTopics  map[string]struct{}
	msgWriters []MsgWriter
}

func NewProducer() MsgSender {
	brokers := viper.GetStringSlice(FlagBrokers)
	topics := viper.GetString(FlagTopics)
	featureToggle := viper.GetBool(FlagFeatureToggle)
	return NewProducerFromConfig(brokers, topics, featureToggle)
}

func NewProducerFromConfig(brokers []string, topics string, featureToggle bool) MsgSender {
	p := producer{
		subTopics:  make(map[string]struct{}),
		msgWriters: nil,
	}

	p.init(brokers, topics, featureToggle)
	return p
}

func (p *producer) init(brokers []string, topics string, featureToggle bool) {
	if len(brokers) == 0 || len(topics) == 0 {
		return
	}

	for _, broker := range brokers {
		msgWriter, err := createMsgWriter(broker)
		if err != nil {
			return // TODO log?
		}
		p.msgWriters = append(p.msgWriters, msgWriter)
	}
	ts := strings.Split(topics, ",")
	for _, topic := range ts {
		p.subTopics[topic] = struct{}{}
	}
	p.toggle = featureToggle
}

func (p producer) Close() {
	for _, w := range p.msgWriters {
		_ = w.Close()
	}
}

func (p producer) SendMsg(k []byte, v []byte) {
	for _, w := range p.msgWriters {
		err := Retry(100, time.Microsecond, func() error {
			return w.WriteKV(k, v)
		})
		// TODO. will log
		_ = err
	}
}

func (p producer) IsSubscribed(topic string) bool {
	if !p.toggle {
		return false
	}
	_, ok := p.subTopics[topic]
	return ok
}

func (p producer) IsOpenToggle() bool {
	return p.toggle
}

func (p producer) GetMode() []string {
	tags := make([]string, 0, len(p.msgWriters))
	for _, w := range p.msgWriters {
		tags = append(tags, w.String())
	}
	return tags
}

func Retry(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return Retry(attempts, 2*sleep, fn)
		}
		return err
	}
	return nil
}
