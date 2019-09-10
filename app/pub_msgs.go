package app

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/coinexchain/dex/msgqueue"
)

type PubMsg struct {
	Key   []byte
	Value []byte
}

func collectKafkaEvents(events []abci.Event, app *CetChainApp) []abci.Event {
	nonKafkaEvents := make([]abci.Event, 0, len(events))
	for _, event := range events {
		if event.Type == msgqueue.EventTypeMsgQueue {
			for _, attr := range event.Attributes {
				app.appendPubMsg(PubMsg{Key: attr.Key, Value: attr.Value})
			}
		} else {
			nonKafkaEvents = append(nonKafkaEvents, event)
		}
	}
	return nonKafkaEvents
}

func discardKafkaEvents(events []abci.Event) []abci.Event {
	nonKafkaEvents := make([]abci.Event, 0, len(events)) // TODO: no need to make new slice
	for _, event := range events {
		if event.Type != msgqueue.EventTypeMsgQueue {
			nonKafkaEvents = append(nonKafkaEvents, event)
		}
	}
	return nonKafkaEvents
}
