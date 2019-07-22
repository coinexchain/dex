package app

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/coinexchain/dex/modules/msgqueue"
)

var PubMsgs []PubMsg

type PubMsg struct {
	Key   []byte
	Value []byte
}

func FilterMsgsOnlyKafka(events []abci.Event) []abci.Event {
	evs := make([]abci.Event, 0, len(events))
	for _, event := range events {
		if event.Type == msgqueue.EventTypeMsgQueue {
			for _, attr := range event.Attributes {
				PubMsgs = append(PubMsgs, PubMsg{Key: attr.Key, Value: attr.Value})
			}
		} else {
			evs = append(evs, event)
		}
	}
	return evs
}

func RemoveMsgsOnlyKafka(events []abci.Event) []abci.Event {
	evs := make([]abci.Event, 0, len(events))
	for _, event := range events {
		if event.Type != msgqueue.EventTypeMsgQueue {
			evs = append(evs, event)
		}
	}
	return evs
}
