package keepers

import (
	"bytes"
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/market/internal/types"
)

//nolint
var (
	OrderBookKeyPrefix     = []byte{0x11}
	BidListKeyPrefix       = []byte{0x12}
	AskListKeyPrefix       = []byte{0x13}
	OrderQueueKeyPrefix    = []byte{0x14}
	LastOrderCleanUpDayKey = []byte{0x20}
)

// This keeper records at which day the last GTE-order-clean-up action was performed
type OrderCleanUpDayKeeper struct {
	marketKey sdk.StoreKey
}

func NewOrderCleanUpDayKeeper(key sdk.StoreKey) *OrderCleanUpDayKeeper {
	return &OrderCleanUpDayKeeper{
		marketKey: key,
	}
}

func (keeper *OrderCleanUpDayKeeper) GetUnixTime(ctx sdk.Context) int64 {
	store := ctx.KVStore(keeper.marketKey)
	value := store.Get(LastOrderCleanUpDayKey)
	if len(value) == 0 {
		return ctx.BlockHeader().Time.Unix()
	}
	return int64(binary.BigEndian.Uint64(value))
}

func (keeper *OrderCleanUpDayKeeper) SetUnixTime(ctx sdk.Context, unixTime int64) {
	value := make([]byte, 8)
	binary.BigEndian.PutUint64(value, uint64(unixTime))
	store := ctx.KVStore(keeper.marketKey)
	store.Set(LastOrderCleanUpDayKey, value[:])
}

/////////////////////////////////////////////////////////////////////

// OrderKeeper manages the order book of one market
type OrderKeeper interface {
	Add(ctx sdk.Context, order *types.Order) sdk.Error
	Remove(ctx sdk.Context, order *types.Order) sdk.Error
	GetOlderThan(ctx sdk.Context, height int64) []*types.Order
	GetOrdersAtHeight(ctx sdk.Context, height int64) []*types.Order
	GetMatchingCandidates(ctx sdk.Context) []*types.Order
	GetSymbol() string
}

// PersistentOrderKeeper implements OrderKeeper interface with a KVStore
type PersistentOrderKeeper struct {
	marketKey sdk.StoreKey
	symbol    string
	codec     *codec.Codec
}

// Merge several byte slices into one
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

// build the key for global order book
func orderBookKey(orderID string) []byte {
	return concatCopyPreAllocate([][]byte{
		OrderBookKeyPrefix,
		{0x0},
		[]byte(orderID),
	})
}

// build the key for bid list
func (keeper *PersistentOrderKeeper) bidListKey(order *types.Order) []byte {
	return concatCopyPreAllocate([][]byte{
		BidListKeyPrefix,
		[]byte(keeper.symbol),
		{0x0},
		types.DecToBigEndianBytes(order.Price),
		[]byte(order.OrderID()),
	})
}

// build the key for ask list
func (keeper *PersistentOrderKeeper) askListKey(order *types.Order) []byte {
	return concatCopyPreAllocate([][]byte{
		AskListKeyPrefix,
		[]byte(keeper.symbol),
		{0x0},
		types.DecToBigEndianBytes(order.Price),
		[]byte(order.OrderID()),
	})
}

// build the key for order queue
func (keeper *PersistentOrderKeeper) orderQueueKey(order *types.Order) []byte {
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

func int64ToBigEndianBytes(n int64) []byte {
	var result [8]byte
	for i := 0; i < 8; i++ {
		result[i] = byte(n >> (8 * uint(i)))
	}
	return result[:]
}

func (keeper *PersistentOrderKeeper) Add(ctx sdk.Context, order *types.Order) sdk.Error {
	// add it to the global order book
	store := ctx.KVStore(keeper.marketKey)
	key := orderBookKey(order.OrderID())
	value := keeper.codec.MustMarshalBinaryBare(order)
	store.Set(key, value)

	// add it to the local order queue
	key = keeper.orderQueueKey(order)
	store.Set(key, []byte{})

	// add it to the local bidList and askList
	if order.Side == types.BID {
		key = keeper.bidListKey(order)
		store.Set(key, []byte{})
	}
	if order.Side == types.ASK {
		key = keeper.askListKey(order)
		store.Set(key, []byte{})
	}
	return nil
}

func (keeper *PersistentOrderKeeper) Remove(ctx sdk.Context, order *types.Order) sdk.Error {
	// remove it from the global order book
	store := ctx.KVStore(keeper.marketKey)
	if keeper.getOrder(ctx, order.OrderID()) == nil {
		return types.ErrNoExistKeyInStore()
	}
	key := orderBookKey(order.OrderID())
	store.Delete(key)

	// remove it from the local order queue
	key = keeper.orderQueueKey(order)
	store.Delete(key)

	// remove it from the local bidList and askList
	if order.Side == types.BID {
		key = keeper.bidListKey(order)
		store.Delete(key)
	}
	if order.Side == types.ASK {
		key = keeper.askListKey(order)
		store.Delete(key)
	}
	return nil
}

// using the order queue, find orders which are older than a particular height
func (keeper *PersistentOrderKeeper) GetOlderThan(ctx sdk.Context, height int64) []*types.Order {
	store := ctx.KVStore(keeper.marketKey)
	var result []*types.Order
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
	iter := store.ReverseIterator(start, end)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		ikey := iter.Key()
		orderID := string(ikey[len(end):])
		result = append(result, keeper.getOrder(ctx, orderID))
	}
	return result
}

// using the order queue, find orders which locate at a particular height
func (keeper *PersistentOrderKeeper) GetOrdersAtHeight(ctx sdk.Context, height int64) []*types.Order {
	store := ctx.KVStore(keeper.marketKey)
	var result []*types.Order
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
	iter := store.Iterator(start, end)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		ikey := iter.Key()
		orderID := string(ikey[len(end):])
		result = append(result, keeper.getOrder(ctx, orderID))
	}
	return result
}

// Get a order's pointer, given its ID
func (keeper *PersistentOrderKeeper) getOrder(ctx sdk.Context, orderID string) *types.Order {
	store := ctx.KVStore(keeper.marketKey)
	key := orderBookKey(orderID)
	orderBytes := store.Get(key)
	if len(orderBytes) == 0 {
		return nil
	}
	order := &types.Order{}
	keeper.codec.MustUnmarshalBinaryBare(orderBytes, order)
	return order
}

// Return the bid orders and ask orders which have proper prices and have possibilities for deal
func (keeper *PersistentOrderKeeper) GetMatchingCandidates(ctx sdk.Context) []*types.Order {
	store := ctx.KVStore(keeper.marketKey)
	priceStartPos := len(keeper.symbol) + 2
	priceEndPos := priceStartPos + types.DecByteCount
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
	defer func() {
		bidIter.Close()
		askIter.Close()
	}()
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
	result := make([]*types.Order, 0, len(orderIDList))
	for _, orderID := range orderIDList {
		result = append(result, keeper.getOrder(ctx, orderID))
	}
	return result
}

////////////////////////////////////////////////

// Global order keep can lookup a order, given its ID or the prefix of its ID, i.e. the sender's address
type GlobalOrderKeeper interface {
	GetAllOrders(ctx sdk.Context) []*types.Order
	QueryOrder(ctx sdk.Context, orderID string) *types.Order
	GetOrdersFromUser(ctx sdk.Context, user string) []string
}

type PersistentGlobalOrderKeeper struct {
	marketKey sdk.StoreKey
	codec     *codec.Codec
}

func NewGlobalOrderKeeper(key sdk.StoreKey, codec *codec.Codec) GlobalOrderKeeper {
	return &PersistentGlobalOrderKeeper{
		marketKey: key,
		codec:     codec,
	}
}

func (keeper *PersistentGlobalOrderKeeper) GetOrdersFromUser(ctx sdk.Context, user string) []string {
	store := ctx.KVStore(keeper.marketKey)
	key := orderBookKey(user + "-")
	nextKey := orderBookKey(user + string([]byte{0xFF}))
	startPos := len(key) - len(user) - 1
	var result []string
	iter := store.Iterator(key, nextKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		k := iter.Key()
		result = append(result, string(k[startPos:]))
	}
	return result
}

// Get all the orders out. It is an expensive operation. Only use it for dumping state.
func (keeper *PersistentGlobalOrderKeeper) GetAllOrders(ctx sdk.Context) []*types.Order {
	store := ctx.KVStore(keeper.marketKey)
	var result []*types.Order
	start := concatCopyPreAllocate([][]byte{
		OrderBookKeyPrefix,
		{0x0},
	})
	end := concatCopyPreAllocate([][]byte{
		OrderBookKeyPrefix,
		{0x1},
	})

	iter := store.Iterator(start, end)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		order := &types.Order{}
		keeper.codec.MustUnmarshalBinaryBare(iter.Value(), order)
		result = append(result, order)
	}
	return result
}

func (keeper *PersistentGlobalOrderKeeper) QueryOrder(ctx sdk.Context, orderID string) *types.Order {
	store := ctx.KVStore(keeper.marketKey)
	key := orderBookKey(orderID)
	orderBytes := store.Get(key)
	if len(orderBytes) == 0 {
		return nil
	}
	order := &types.Order{}
	keeper.codec.MustUnmarshalBinaryBare(orderBytes, order)
	return order
}
