package keepers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	dex "github.com/coinexchain/dex/types"
)

type DelistKeeper struct {
	marketKey sdk.StoreKey
}

func NewDelistKeeper(key sdk.StoreKey) *DelistKeeper {
	return &DelistKeeper{
		marketKey: key,
	}
}

func getDelistKey(time int64, symbol string) []byte {
	return dex.ConcatKeys(
		DelistKey,
		int64ToBigEndianBytes(time),
		[]byte{0x0},
		[]byte(symbol),
	)
}

func getDelistKeyRangeByTime(time int64) (start, end []byte) {
	start = dex.ConcatKeys(DelistKey, int64ToBigEndianBytes(0), []byte{0x0})
	end = dex.ConcatKeys(DelistKey, int64ToBigEndianBytes(time), []byte{0x1})
	return
}

func (keeper *DelistKeeper) AddDelistRequest(ctx sdk.Context, time int64, symbol string) {
	store := ctx.KVStore(keeper.marketKey)
	store.Set(getDelistKey(time, symbol), []byte(symbol))
	store.Set(append(DelistRevKey, []byte(symbol)...), int64ToBigEndianBytes(time))
}

func (keeper *DelistKeeper) HasDelistRequest(ctx sdk.Context, symbol string) bool {
	store := ctx.KVStore(keeper.marketKey)
	return store.Has(append(DelistRevKey, []byte(symbol)...))
}

//include the specific time
func (keeper *DelistKeeper) GetDelistSymbolsBeforeTime(ctx sdk.Context, time int64) []string {
	store := ctx.KVStore(keeper.marketKey)
	start, end := getDelistKeyRangeByTime(time)
	var result []string
	iter := store.Iterator(start, end)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		result = append(result, string(key[len(start):]))
	}
	return result
}

func (keeper *DelistKeeper) RemoveDelistRequestsBeforeTime(ctx sdk.Context, time int64) {
	store := ctx.KVStore(keeper.marketKey)
	start, end := getDelistKeyRangeByTime(time)
	keys := make([][]byte, 0, 100)
	symbols := make([][]byte, 0, 100)
	iter := store.Iterator(start, end)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		keys = append(keys, iter.Key())
		symbols = append(symbols, iter.Value())
	}
	for _, key := range keys {
		store.Delete(key)
	}
	for _, symbol := range symbols {
		store.Delete(append(DelistRevKey, symbol...))
	}
}
