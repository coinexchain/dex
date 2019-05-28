package market

import (
	"bytes"
	"fmt"
	"github.com/coinexchain/dex/modules/market/match"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//nolint
var (
	OrderBookKeyPrefix         = []byte{0x11}
	BidListKeyPrefix           = []byte{0x12}
	BidListKeyPrefixPlusOne    = []byte{0x13}
	AskListKeyPrefix           = []byte{0x13}
	AskListKeyPrefixPlusOne    = []byte{0x14}
	OrderQueueKeyPrefix        = []byte{0x14}
	OrderQueueKeyPrefixPlusOne = []byte{0x15}
	IocListKeyPrefix           = []byte{0x15}
)

const (
	DecByteCount = 40
	GTE          = 3
	IOC          = 4
	Buy          = match.BUY
	Sell         = match.SELL
	LIMIT        = 2
)

type OrderKeeper interface {
	Add(order *Order) sdk.Error
	Exists(orderID string) bool
	Remove(order *Order) sdk.Error
	GetOlderThan(height int64) []*Order
	GetOrdersAtHeight(height int64) []*Order
	QueryOrder(orderID string) *Order
	GetOrdersFromUser(user string) []string
	GetMatchingCandidates() []*Order
}

type PersistentOrderKeeper struct {
	store  sdk.KVStore
	symbol string
	codec  *codec.Codec
}

func concatCopyPreAllocate(slices [][]byte) []byte {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]byte, totalLen)
	var i int
	for _, s := range slices {
		i += copy(tmp[i:], s)
	}
	return tmp
}

func (keeper *PersistentOrderKeeper) orderBookKey(orderID string) []byte {
	return concatCopyPreAllocate([][]byte{
		[]byte(keeper.symbol),
		{0x0},
		OrderBookKeyPrefix,
		[]byte(orderID),
	})
}

func (keeper *PersistentOrderKeeper) bidListKey(order *Order) []byte {
	return concatCopyPreAllocate([][]byte{
		[]byte(keeper.symbol),
		{0x0},
		BidListKeyPrefix,
		decToBigEndianBytes(order.Price),
		[]byte(order.OrderID()),
	})
}

func (keeper *PersistentOrderKeeper) askListKey(order *Order) []byte {
	return concatCopyPreAllocate([][]byte{
		[]byte(keeper.symbol),
		{0x0},
		AskListKeyPrefix,
		decToBigEndianBytes(order.Price),
		[]byte(order.OrderID()),
	})
}

func (keeper *PersistentOrderKeeper) orderQueueKey(order *Order) []byte {
	return concatCopyPreAllocate([][]byte{
		[]byte(keeper.symbol),
		{0x0},
		OrderQueueKeyPrefix,
		int64ToBigEndianBytes(order.Height),
		[]byte(order.OrderID()),
	})
}

func NewOrderKeeper(store sdk.KVStore, symbol string, codec *codec.Codec) OrderKeeper {
	return &PersistentOrderKeeper{
		store:  store,
		symbol: symbol,
		codec:  codec,
	}
}

func decToBigEndianBytes(d sdk.Dec) []byte {
	var result [DecByteCount]byte
	bytes := d.Int.Bytes()
	for i := 1; i <= len(bytes); i++ {
		result[DecByteCount-i] = bytes[len(bytes)-i]
	}
	return result[:]
}

func int64ToBigEndianBytes(n int64) []byte {
	var result [8]byte
	for i := 0; i < 8; i++ {
		result[i] = byte(n >> (8 * uint(i)))
	}
	return result[:]
}

func (keeper *PersistentOrderKeeper) Add(order *Order) sdk.Error {
	key := keeper.orderBookKey(order.OrderID())
	value := keeper.codec.MustMarshalBinaryBare(order)
	keeper.store.Set(key, value)

	if order.TimeInForce == GTE {
		key = keeper.orderQueueKey(order)
		keeper.store.Set(key, nil)
	}
	if order.Side == match.BID {
		key = keeper.bidListKey(order)
		keeper.store.Set(key, nil)
	}
	if order.Side == match.ASK {
		key = keeper.askListKey(order)
		keeper.store.Set(key, nil)
	}
	return nil
}

func (keeper *PersistentOrderKeeper) Exists(orderID string) bool {
	key := keeper.orderBookKey(orderID)
	return keeper.store.Has(key)
}

func (keeper *PersistentOrderKeeper) Remove(order *Order) sdk.Error {
	if !keeper.Exists(order.OrderID()) {
		return ErrNoExistKeyInStore()
	}
	key := keeper.orderBookKey(order.OrderID())
	keeper.store.Delete(key)

	if order.TimeInForce == GTE {
		key = keeper.orderQueueKey(order)
		keeper.store.Delete(key)
	}
	if order.Side == match.BID {
		key = keeper.bidListKey(order)
		keeper.store.Delete(key)
	}
	if order.Side == match.ASK {
		key = keeper.askListKey(order)
		keeper.store.Delete(key)
	}
	return nil
}

func (keeper *PersistentOrderKeeper) GetOlderThan(height int64) []*Order {
	var result []*Order
	start := concatCopyPreAllocate([][]byte{
		[]byte(keeper.symbol),
		{0x0},
		OrderQueueKeyPrefix,
	})
	end := concatCopyPreAllocate([][]byte{
		[]byte(keeper.symbol),
		{0x0},
		OrderQueueKeyPrefix,
		int64ToBigEndianBytes(height),
	})
	for iter := keeper.store.ReverseIterator(start, end); iter.Valid(); iter.Next() {
		ikey := iter.Key()
		orderID := string(ikey[len(end):])
		result = append(result, keeper.QueryOrder(orderID))
	}
	return result
}

func (keeper *PersistentOrderKeeper) GetOrdersAtHeight(height int64) []*Order {
	var result []*Order
	start := concatCopyPreAllocate([][]byte{
		[]byte(keeper.symbol),
		{0x0},
		OrderQueueKeyPrefix,
		int64ToBigEndianBytes(height),
	})
	end := concatCopyPreAllocate([][]byte{
		[]byte(keeper.symbol),
		{0x0},
		OrderQueueKeyPrefix,
		int64ToBigEndianBytes(height + 1),
	})
	for iter := keeper.store.Iterator(start, end); iter.Valid(); iter.Next() {
		ikey := iter.Key()
		orderID := string(ikey[len(end):])
		result = append(result, keeper.QueryOrder(orderID))
	}
	return result
}

func (keeper *PersistentOrderKeeper) QueryOrder(orderID string) *Order {
	key := keeper.orderBookKey(orderID)
	orderBytes := keeper.store.Get(key)
	order := &Order{}
	keeper.codec.MustUnmarshalBinaryBare(orderBytes, order)
	return order
}

func (keeper *PersistentOrderKeeper) GetOrdersFromUser(user string) []string {
	key := keeper.orderBookKey(user + "-")
	nextKey := keeper.orderBookKey(user + string([]byte{0xFF}))
	startPos := len(key) - len(user) - 1
	var result []string
	for iter := keeper.store.ReverseIterator(key, nextKey); iter.Valid(); iter.Next() {
		k := iter.Key()
		result = append(result, string(k[startPos:]))
	}
	return result
}

func (keeper *PersistentOrderKeeper) GetMatchingCandidates() []*Order {
	priceStartPos := len(keeper.symbol) + 2
	priceEndPos := priceStartPos + DecByteCount
	bidListStart := concatCopyPreAllocate([][]byte{
		[]byte(keeper.symbol),
		{0x0},
		BidListKeyPrefix,
	})
	bidListEnd := concatCopyPreAllocate([][]byte{
		[]byte(keeper.symbol),
		{0x0},
		BidListKeyPrefixPlusOne,
	})
	askListStart := concatCopyPreAllocate([][]byte{
		[]byte(keeper.symbol),
		{0x0},
		AskListKeyPrefix,
	})
	askListEnd := concatCopyPreAllocate([][]byte{
		[]byte(keeper.symbol),
		{0x0},
		AskListKeyPrefixPlusOne,
	})
	bidIter := keeper.store.ReverseIterator(bidListStart, bidListEnd)
	askIter := keeper.store.Iterator(askListStart, askListEnd)
	if !bidIter.Valid() || !askIter.Valid() {
		return nil
	}
	firstBidKey := bidIter.Key()
	firstAskKey := askIter.Key()
	firstBidPrice := firstBidKey[priceStartPos:priceEndPos]
	firstAskPrice := firstAskKey[priceStartPos:priceEndPos]
	if bytes.Compare(firstAskPrice, firstBidPrice) > 0 {
		return nil
	}
	orderIDList := []string{string(firstBidKey[priceEndPos:]), string(firstAskKey[priceEndPos:])}
	for _, s := range orderIDList {
		fmt.Printf("here! %s\n", s)
	}
	for askIter.Next(); askIter.Valid(); askIter.Next() {
		askKey := askIter.Key()
		askPrice := askKey[priceStartPos:priceEndPos]
		if bytes.Compare(askPrice, firstBidPrice) > 0 {
			break
		} else {
			orderIDList = append(orderIDList, string(askKey[priceEndPos:]))
		}
	}
	for bidIter.Next(); bidIter.Valid(); bidIter.Next() {
		bidKey := bidIter.Key()
		bidPrice := bidKey[priceStartPos:priceEndPos]
		if bytes.Compare(firstAskPrice, bidPrice) > 0 {
			break
		} else {
			orderIDList = append(orderIDList, string(bidKey[priceEndPos:]))
		}
	}
	result := make([]*Order, 0, len(orderIDList))
	for _, orderID := range orderIDList {
		result = append(result, keeper.QueryOrder(orderID))
	}
	return result
}
