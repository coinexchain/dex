package market

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/coinexchain/dex/modules/msgqueue"
)

var (
	MarketIdentifierPrefix = []byte{0x15}
	DelistKey              = []byte{0x40}
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

type Keeper struct {
	paramSubspace params.Subspace
	marketKey     sdk.StoreKey
	axk           ExpectedAssetStatusKeeper
	bnk           ExpectedBankxKeeper
	cdc           *codec.Codec
	orderClean    *OrderCleanUpDayKeeper
	gmk           GlobalMarketInfoKeeper
	msgProducer   msgqueue.MsgSender
}

func NewKeeper(key sdk.StoreKey, axkVal ExpectedAssetStatusKeeper,
	bnkVal ExpectedBankxKeeper, feeK ExpectFeeKeeper,
	cdcVal *codec.Codec, msgKeeperVal msgqueue.MsgSender,
	paramstore params.Subspace) Keeper {

	return Keeper{
		marketKey:     key,
		axk:           axkVal,
		bnk:           bnkVal,
		paramSubspace: paramstore.WithKeyTable(ParamKeyTable()),
		cdc:           cdcVal,
		orderClean:    NewOrderCleanUpDayKeeper(key),
		gmk:           NewGlobalMarketInfoKeeper(key, cdcVal),
		msgProducer:   msgKeeperVal,
	}
}

func (k Keeper) SendMsg(key string, val interface{}) {
	k.msgProducer.SendMsg(Topic, key, val)
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the asset module's parameters.
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the asset module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params Params) {
	k.paramSubspace.GetParamSet(ctx, &params)
	return
}

// -----------------------------------------------------------------------------
// Order

// SetOrder implements token Keeper.
func (k Keeper) SetOrder(ctx sdk.Context, order *Order) sdk.Error {
	return NewOrderKeeper(k.marketKey, order.Symbol, k.cdc).Add(ctx, order)
}

func (k Keeper) GetAllOrders(ctx sdk.Context) []*Order {
	return NewGlobalOrderKeeper(k.marketKey, k.cdc).GetAllOrders(ctx)
}

// -----------------------------------------------
// market info

func (k Keeper) SetMarket(ctx sdk.Context, info MarketInfo) sdk.Error {
	return k.gmk.SetMarket(ctx, info)
}

func (k Keeper) RemoveMarket(ctx sdk.Context, symbol string) sdk.Error {
	return k.gmk.RemoveMarket(ctx, symbol)
}

func (k Keeper) GetAllMarketInfos(ctx sdk.Context) []MarketInfo {
	return k.gmk.GetAllMarketInfos(ctx)
}

func (k Keeper) GetMarketInfo(ctx sdk.Context, symbol string) (MarketInfo, error) {
	return k.gmk.GetMarketInfo(ctx, symbol)
}

func (k Keeper) SubtractFeeAndCollectFee(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	return k.bnk.DeductFee(ctx, addr, amt)
}

func (k Keeper) MarketOwner(ctx sdk.Context, info MarketInfo) sdk.AccAddress {
	return k.axk.GetToken(ctx, info.Stock).GetOwner()
}

// -----------------------------------------------------------------------------

type GlobalMarketInfoKeeper interface {
	SetMarket(ctx sdk.Context, info MarketInfo) sdk.Error
	RemoveMarket(ctx sdk.Context, symbol string) sdk.Error
	GetAllMarketInfos(ctx sdk.Context) []MarketInfo
	GetMarketInfo(ctx sdk.Context, symbol string) (MarketInfo, error)
}

type PersistentMarketInfoKeeper struct {
	marketKey sdk.StoreKey
	cdc       *codec.Codec
}

func NewGlobalMarketInfoKeeper(key sdk.StoreKey, cdcVal *codec.Codec) GlobalMarketInfoKeeper {
	return PersistentMarketInfoKeeper{
		marketKey: key,
		cdc:       cdcVal,
	}
}

// SetMarket implements token Keeper.
func (k PersistentMarketInfoKeeper) SetMarket(ctx sdk.Context, info MarketInfo) sdk.Error {
	store := ctx.KVStore(k.marketKey)
	bz, err := k.cdc.MarshalBinaryBare(info)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}
	store.Set(marketStoreKey(MarketIdentifierPrefix, info.Stock+SymbolSeparator+info.Money), bz)
	return nil
}

func (k PersistentMarketInfoKeeper) RemoveMarket(ctx sdk.Context, symbol string) sdk.Error {
	store := ctx.KVStore(k.marketKey)
	key := marketStoreKey(MarketIdentifierPrefix, symbol)
	value := store.Get(key)
	if value != nil {
		store.Delete(key)
	}
	return nil
}

func (k PersistentMarketInfoKeeper) GetAllMarketInfos(ctx sdk.Context) []MarketInfo {
	var infos []MarketInfo
	appendMarket := func(order MarketInfo) (stop bool) {
		infos = append(infos, order)
		return false
	}
	k.iterateMarket(ctx, appendMarket)
	return infos
}

func (k PersistentMarketInfoKeeper) iterateMarket(ctx sdk.Context, process func(info MarketInfo) bool) {
	store := ctx.KVStore(k.marketKey)
	iter := sdk.KVStorePrefixIterator(store, MarketIdentifierPrefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		if process(k.decodeMarket(val)) {
			return
		}
		iter.Next()
	}
}

func (k PersistentMarketInfoKeeper) GetMarketInfo(ctx sdk.Context, symbol string) (info MarketInfo, err error) {
	store := ctx.KVStore(k.marketKey)
	value := store.Get(marketStoreKey(MarketIdentifierPrefix, symbol))
	err = k.cdc.UnmarshalBinaryBare(value, &info)
	return
}

func (k PersistentMarketInfoKeeper) decodeMarket(bz []byte) (info MarketInfo) {
	if err := k.cdc.UnmarshalBinaryBare(bz, &info); err != nil {
		panic(err)
	}
	return
}
