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

type config struct {
	Brokers string
	Topics  string
}

func NewProducer() MsgSender {
	p := producer{
		subTopics: make(map[string]struct{}),
		msgWriter: NewNopMsgWriter(),
	}

	p.setParam(config{
		Brokers: viper.GetString(FlagBrokers),
		Topics:  viper.GetString(FlagTopics),
	})
	return p
}

func (p *producer) setParam(cfg config) {
	if len(cfg.Brokers) == 0 || len(cfg.Topics) == 0 {
		return
	}

	msgWriter, err := createMsgWriter(cfg.Brokers)
	if err != nil {
		return // TODO log?
	}
	p.msgWriter = msgWriter

	ts := strings.Split(cfg.Topics, ",")
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
