package msgqueue

import (
	"strings"

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
	GetMode() string
}

type producer struct {
	toggle    bool
	subTopics map[string]struct{}
	msgWriter MsgWriter
}

func NewProducer() MsgSender {
	brokers := viper.GetString(FlagBrokers)
	topics := viper.GetString(FlagTopics)
	return newProducerFromConfig(brokers, topics)
}

func newProducerFromConfig(brokers, topics string) MsgSender {
	p := producer{
		subTopics: make(map[string]struct{}),
		msgWriter: NewNopMsgWriter(),
	}

	p.init(brokers, topics)
	return p
}

func (p *producer) init(brokers, topics string) {
	if len(brokers) == 0 || len(topics) == 0 {
		return
	}

	msgWriter, err := createMsgWriter(brokers)
	if err != nil {
		return // TODO log?
	}
	p.msgWriter = msgWriter

	ts := strings.Split(topics, ",")
	for _, topic := range ts {
		p.subTopics[topic] = struct{}{}
	}
	p.toggle = viper.GetBool(FlagFeatureToggle)
}

func (p producer) Close() {
	_ = p.msgWriter.Close() // TODO: handle error
}

func (p producer) SendMsg(k []byte, v []byte) {
	_ = p.msgWriter.WriteKV(k, v) // TODO: handle error
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

func (p producer) GetMode() string {
	return p.msgWriter.String()
}
