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
	OrderBookKeyPrefix     = []byte{0x11}
	BidListKeyPrefix       = []byte{0x12}
	AskListKeyPrefix       = []byte{0x13}
	OrderQueueKeyPrefix    = []byte{0x14}
	LastOrderCleanUpDayKey = []byte{0x20}
)

const (
	DecByteCount = 40
	GTE          = 3
	IOC          = 4
	Buy          = match.BUY
	Sell         = match.SELL
	LIMIT        = 2
)

type OrderCleanUpDayKeeper struct {
	marketKey sdk.StoreKey
}

func NewOrderCleanUpDayKeeper(key sdk.StoreKey) *OrderCleanUpDayKeeper {
	return &OrderCleanUpDayKeeper{
		marketKey: key,
	}
}

func (keeper *OrderCleanUpDayKeeper) GetDay(ctx sdk.Context) int {
	store := ctx.KVStore(keeper.marketKey)
	value := store.Get(LastOrderCleanUpDayKey)
	return int(value[0])
}

func (keeper *OrderCleanUpDayKeeper) SetDay(ctx sdk.Context, day int) {
	var value [1]byte
	value[0] = byte(day)
	store := ctx.KVStore(keeper.marketKey)
	store.Set(LastOrderCleanUpDayKey, value[:])
}

type OrderKeeper interface {
	Add(ctx sdk.Context, order *Order) sdk.Error
	Exists(ctx sdk.Context, orderID string) bool
	Remove(ctx sdk.Context, order *Order) sdk.Error
	GetOlderThan(ctx sdk.Context, height int64) []*Order
	GetAllOrders(ctx sdk.Context) []*Order
	RemoveAllOrders(ctx sdk.Context)
	GetOrdersAtHeight(ctx sdk.Context, height int64) []*Order
	QueryOrder(ctx sdk.Context, orderID string) *Order
	GetOrdersFromUser(ctx sdk.Context, user string) []string
	GetMatchingCandidates(ctx sdk.Context) []*Order
	GetSymbol() string
}

type PersistentOrderKeeper struct {
	marketKey sdk.StoreKey
	symbol    string
	codec     *codec.Codec
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

func (keeper *PersistentOrderKeeper) GetSymbol() string {
	return keeper.symbol
}

func (keeper *PersistentOrderKeeper) orderBookKey(orderID string) []byte {
	return concatCopyPreAllocate([][]byte{
		OrderBookKeyPrefix,
		[]byte(keeper.symbol),
		{0x0},
		[]byte(orderID),
	})
}

func (keeper *PersistentOrderKeeper) bidListKey(order *Order) []byte {
	return concatCopyPreAllocate([][]byte{
		BidListKeyPrefix,
		[]byte(keeper.symbol),
		{0x0},
		decToBigEndianBytes(order.Price),
		[]byte(order.OrderID()),
	})
}

func (keeper *PersistentOrderKeeper) askListKey(order *Order) []byte {
	return concatCopyPreAllocate([][]byte{
		AskListKeyPrefix,
		[]byte(keeper.symbol),
		{0x0},
		decToBigEndianBytes(order.Price),
		[]byte(order.OrderID()),
	})
}

func (keeper *PersistentOrderKeeper) orderQueueKey(order *Order) []byte {
	return concatCopyPreAllocate([][]byte{
		OrderQueueKeyPrefix,
		[]byte(keeper.symbol),
		{0x0},
		int64ToBigEndianBytes(order.Height),
		[]byte(order.OrderID()),
	})
}

func NewOrderKeeper(key sdk.StoreKey, symbol string, codec *codec.Codec) OrderKeeper {
	return &PersistentOrderKeeper{
		marketKey: key,
		symbol:    symbol,
		codec:     codec,
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

func (keeper *PersistentOrderKeeper) Add(ctx sdk.Context, order *Order) sdk.Error {
	store := ctx.KVStore(keeper.marketKey)
	key := keeper.orderBookKey(order.OrderID())
	value := keeper.codec.MustMarshalBinaryBare(order)
	store.Set(key, value)

	if order.TimeInForce == GTE {
		key = keeper.orderQueueKey(order)
		store.Set(key, []byte{})
	}
	if order.Side == match.BID {
		key = keeper.bidListKey(order)
		store.Set(key, []byte{})
	}
	if order.Side == match.ASK {
		key = keeper.askListKey(order)
		store.Set(key, []byte{})
	}
	return nil
}

func (keeper *PersistentOrderKeeper) Exists(ctx sdk.Context, orderID string) bool {
	store := ctx.KVStore(keeper.marketKey)
	key := keeper.orderBookKey(orderID)
	return store.Has(key)
}

func (keeper *PersistentOrderKeeper) Remove(ctx sdk.Context, order *Order) sdk.Error {
	store := ctx.KVStore(keeper.marketKey)
	if !keeper.Exists(ctx, order.OrderID()) {
		return ErrNoExistKeyInStore()
	}
	key := keeper.orderBookKey(order.OrderID())
	store.Delete(key)

	if order.TimeInForce == GTE {
		key = keeper.orderQueueKey(order)
		store.Delete(key)
	}
	if order.Side == match.BID {
		key = keeper.bidListKey(order)
		store.Delete(key)
	}
	if order.Side == match.ASK {
		key = keeper.askListKey(order)
		store.Delete(key)
	}
	return nil
}

func (keeper *PersistentOrderKeeper) GetOlderThan(ctx sdk.Context, height int64) []*Order {
	store := ctx.KVStore(keeper.marketKey)
	var result []*Order
	start := concatCopyPreAllocate([][]byte{
		OrderQueueKeyPrefix,
		[]byte(keeper.symbol),
		{0x0},
	})
	end := concatCopyPreAllocate([][]byte{
		OrderQueueKeyPrefix,
		[]byte(keeper.symbol),
		{0x0},
		int64ToBigEndianBytes(height),
	})
	for iter := store.ReverseIterator(start, end); iter.Valid(); iter.Next() {
		ikey := iter.Key()
		orderID := string(ikey[len(end):])
		result = append(result, keeper.QueryOrder(ctx, orderID))
	}
	return result
}

func (keeper *PersistentOrderKeeper) GetAllOrders(ctx sdk.Context) []*Order {
	store := ctx.KVStore(keeper.marketKey)
	var result []*Order
	start := concatCopyPreAllocate([][]byte{
		OrderBookKeyPrefix,
		[]byte(keeper.symbol),
		{0x0},
	})
	end := concatCopyPreAllocate([][]byte{
		OrderBookKeyPrefix,
		[]byte(keeper.symbol),
		{0x1},
	})
	for iter := store.Iterator(start, end); iter.Valid(); iter.Next() {
		order := &Order{}
		keeper.codec.MustUnmarshalBinaryBare(iter.Value(), order)
		result = append(result, order)
	}
	return result
}

func (keeper *PersistentOrderKeeper) RemoveAllOrders(ctx sdk.Context) {
	store := ctx.KVStore(keeper.marketKey)
	keys := make([][]byte, 0, 10)
	keeper.fillKeys(store, keys, OrderBookKeyPrefix)
	keeper.fillKeys(store, keys, BidListKeyPrefix)
	keeper.fillKeys(store, keys, AskListKeyPrefix)
	keeper.fillKeys(store, keys, OrderQueueKeyPrefix)
	for _, key := range keys {
		store.Delete(key)
	}
}

func (keeper *PersistentOrderKeeper) fillKeys(store sdk.KVStore, keys [][]byte, keyPrefix []byte) {
	start := concatCopyPreAllocate([][]byte{
		keyPrefix,
		[]byte(keeper.symbol),
		{0x0},
	})
	end := concatCopyPreAllocate([][]byte{
		keyPrefix,
		[]byte(keeper.symbol),
		{0x1},
	})
	for iter := store.Iterator(start, end); iter.Valid(); iter.Next() {
		keys = append(keys, iter.Key())
	}
}

func (keeper *PersistentOrderKeeper) GetOrdersAtHeight(ctx sdk.Context, height int64) []*Order {
	store := ctx.KVStore(keeper.marketKey)
	var result []*Order
	start := concatCopyPreAllocate([][]byte{
		OrderQueueKeyPrefix,
		[]byte(keeper.symbol),
		{0x0},
		int64ToBigEndianBytes(height),
	})
	end := concatCopyPreAllocate([][]byte{
		OrderQueueKeyPrefix,
		[]byte(keeper.symbol),
		{0x0},
		int64ToBigEndianBytes(height + 1),
	})
	for iter := store.Iterator(start, end); iter.Valid(); iter.Next() {
		ikey := iter.Key()
		orderID := string(ikey[len(end):])
		result = append(result, keeper.QueryOrder(ctx, orderID))
	}
	return result
}

func (keeper *PersistentOrderKeeper) QueryOrder(ctx sdk.Context, orderID string) *Order {
	store := ctx.KVStore(keeper.marketKey)
	key := keeper.orderBookKey(orderID)
	orderBytes := store.Get(key)
	if len(orderBytes) == 0 {
		return nil
	}
	order := &Order{}
	keeper.codec.MustUnmarshalBinaryBare(orderBytes, order)
	return order
}

func (keeper *PersistentOrderKeeper) GetOrdersFromUser(ctx sdk.Context, user string) []string {
	store := ctx.KVStore(keeper.marketKey)
	key := keeper.orderBookKey(user + "-")
	nextKey := keeper.orderBookKey(user + string([]byte{0xFF}))
	startPos := len(key) - len(user) - 1
	var result []string
	for iter := store.ReverseIterator(key, nextKey); iter.Valid(); iter.Next() {
		k := iter.Key()
		result = append(result, string(k[startPos:]))
	}
	return result
}

func (keeper *PersistentOrderKeeper) GetMatchingCandidates(ctx sdk.Context) []*Order {
	store := ctx.KVStore(keeper.marketKey)
	priceStartPos := len(keeper.symbol) + 2
	priceEndPos := priceStartPos + DecByteCount
	bidListStart := concatCopyPreAllocate([][]byte{
		BidListKeyPrefix,
		[]byte(keeper.symbol),
		{0x0},
	})
	bidListEnd := concatCopyPreAllocate([][]byte{
		BidListKeyPrefix,
		[]byte(keeper.symbol),
		{0x1},
	})
	askListStart := concatCopyPreAllocate([][]byte{
		AskListKeyPrefix,
		[]byte(keeper.symbol),
		{0x0},
	})
	askListEnd := concatCopyPreAllocate([][]byte{
		AskListKeyPrefix,
		[]byte(keeper.symbol),
		{0x1},
	})
	bidIter := store.ReverseIterator(bidListStart, bidListEnd)
	askIter := store.Iterator(askListStart, askListEnd)
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
		result = append(result, keeper.QueryOrder(ctx, orderID))
	}
	return result
}
