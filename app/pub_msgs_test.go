package app

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/coinexchain/cet-sdk/msgqueue"
)

func TestCollectKafkaEvents(t *testing.T) {
	events := []abci.Event{
		{
			Type: msgqueue.EventTypeMsgQueue,
			Attributes: []common.KVPair{
				{Key: []byte("k1"), Value: []byte("v1")},
				{Key: []byte("k2"), Value: []byte("v2")},
			},
		},
		{Type: "other"},
		{
			Type: msgqueue.EventTypeMsgQueue,
			Attributes: []common.KVPair{
				{Key: []byte("k3"), Value: []byte("v3")},
			},
		},
		{Type: "other"},
	}

	fakeApp := &CetChainApp{pubMsgs: nil}
	events = collectKafkaEvents(events, fakeApp)
	require.Equal(t, 2, len(events))
	require.Equal(t, "other", events[0].Type)
	require.Equal(t, "other", events[1].Type)
	require.Equal(t, 3, len(fakeApp.pubMsgs))
	require.Equal(t, "k1", string(fakeApp.pubMsgs[0].Key))
	require.Equal(t, "k2", string(fakeApp.pubMsgs[1].Key))
	require.Equal(t, "k3", string(fakeApp.pubMsgs[2].Key))
}

func TestDiscardKafkaEvents(t *testing.T) {
	events := []abci.Event{
		{Type: msgqueue.EventTypeMsgQueue},
		{Type: "other"},
		{Type: msgqueue.EventTypeMsgQueue},
		{Type: "other"},
	}
	events = discardKafkaEvents(events)
	require.Equal(t, 2, len(events))
	require.Equal(t, "other", events[0].Type)
	require.Equal(t, "other", events[1].Type)
}
