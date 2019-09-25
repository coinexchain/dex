package keepers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	return concatCopyPreAllocate([][]byte{
		DelistKey,
		int64ToBigEndianBytes(time),
		{0x0},
		[]byte(symbol),
	})
}

func (keeper *DelistKeeper) AddDelistRequest(ctx sdk.Context, time int64, symbol string) {
	store := ctx.KVStore(keeper.marketKey)
	store.Set(getDelistKey(time, symbol), []byte{})
}

func (keeper *DelistKeeper) GetDelistSymbolsBeforeTime(ctx sdk.Context, time int64) []string {
	store := ctx.KVStore(keeper.marketKey)
	start := concatCopyPreAllocate([][]byte{
		DelistKey,
		int64ToBigEndianBytes(0),
	})
	end := concatCopyPreAllocate([][]byte{
		DelistKey,
		int64ToBigEndianBytes(time),
	})
	var result []string
	iter := store.Iterator(start, end)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		result = append(result, string(key[len(start)+1:]))
	}
	return result
}

func (keeper *DelistKeeper) RemoveDelistRequestsBeforeTime(ctx sdk.Context, time int64) {
	store := ctx.KVStore(keeper.marketKey)
	start := concatCopyPreAllocate([][]byte{
		DelistKey,
		int64ToBigEndianBytes(0),
		{0x0},
	})
	end := concatCopyPreAllocate([][]byte{
		DelistKey,
		int64ToBigEndianBytes(time),
		{0x1},
	})
	var keys [][]byte
	iter := store.Iterator(start, end)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		keys = append(keys, iter.Key())
	}
	for _, key := range keys {
		store.Delete(key)
	}
}
