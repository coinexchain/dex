package market

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
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
func getDelistKey(height int64, symbol string) []byte {
	return concatCopyPreAllocate([][]byte{
		DelistKey,
		int64ToBigEndianBytes(height),
		{0x0},
		[]byte(symbol),
	})
}
func (keeper *DelistKeeper) AddDelistRequest(ctx sdk.Context, height int64, symbol string) {
	store := ctx.KVStore(keeper.marketKey)
	store.Set(getDelistKey(height, symbol), []byte{})
}
func (keeper *DelistKeeper) GetDelistSymbolsAtHeight(ctx sdk.Context, height int64) []string {
	store := ctx.KVStore(keeper.marketKey)
	start := concatCopyPreAllocate([][]byte{
		DelistKey,
		int64ToBigEndianBytes(height),
		{0x0},
	})
	end := concatCopyPreAllocate([][]byte{
		DelistKey,
		int64ToBigEndianBytes(height),
		{0x1},
	})
	var result []string
	for iter := store.Iterator(start, end); iter.Valid(); iter.Next() {
		key := iter.Key()
		result = append(result, string(key[len(start):]))
	}
	return result
}

func (keeper *DelistKeeper) RemoveDelistRequestsAtHeight(ctx sdk.Context, height int64) {
	store := ctx.KVStore(keeper.marketKey)
	start := concatCopyPreAllocate([][]byte{
		DelistKey,
		int64ToBigEndianBytes(height),
		{0x0},
	})
	end := concatCopyPreAllocate([][]byte{
		DelistKey,
		int64ToBigEndianBytes(height),
		{0x1},
	})
	var keys [][]byte
	for iter := store.Iterator(start, end); iter.Valid(); iter.Next() {
		keys = append(keys, iter.Key())
	}
	for _, key := range keys {
		store.Delete(key)
	}
}

type Keeper struct {
	paramSubspace params.Subspace
	marketKey     sdk.StoreKey
	axk           ExpectedAssertStatusKeeper
	bnk           ExpectedBankxKeeper
	cdc           *codec.Codec
	orderClean    *OrderCleanUpDayKeeper
	//fek       incentive.FeeCollectionKeeper
}

func NewKeeper(key sdk.StoreKey, axkVal ExpectedAssertStatusKeeper,
	bnkVal ExpectedBankxKeeper, cdcVal *codec.Codec, paramstore params.Subspace) Keeper {

	return Keeper{marketKey: key, axk: axkVal, bnk: bnkVal,
		paramSubspace: paramstore.WithKeyTable(ParamKeyTable()),
		cdc:           cdcVal,
		orderClean:    NewOrderCleanUpDayKeeper(key)}
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

// SetOrder implements token Keeper.
func (k Keeper) SetOrder(ctx sdk.Context, order *Order) sdk.Error {
	return NewOrderKeeper(k.marketKey, order.Symbol, k.cdc).Add(ctx, order)
}

// SetMarket implements token Keeper.
func (k Keeper) SetMarket(ctx sdk.Context, info MarketInfo) sdk.Error {
	store := ctx.KVStore(k.marketKey)
	bz, err := k.cdc.MarshalBinaryBare(info)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}
	store.Set(marketStoreKey(MarketIdentifierPrefix, info.Stock+SymbolSeparator+info.Money), bz)
	return nil
}

func (k Keeper) RemoveMarket(ctx sdk.Context, symbol string) sdk.Error {
	store := ctx.KVStore(k.marketKey)
	key := marketStoreKey(MarketIdentifierPrefix, symbol)
	value := store.Get(key)
	if value != nil {
		store.Delete(key)
	}
	return nil
}

// RegisterCodec registers concrete types on the codec
func (k Keeper) RegisterCodec() {
	RegisterCodec(k.cdc)
}

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(Order{}, "market/order", nil)
	cdc.RegisterConcrete(MarketInfo{}, "market/market", nil)
	cdc.RegisterConcrete(MsgCreateMarketInfo{}, "market/market-info", nil)
	cdc.RegisterConcrete(MsgCreateGTEOrder{}, "market/order-info", nil)
}

func (k Keeper) GetAllOrders(ctx sdk.Context) []*Order {
	var orders []*Order
	appendOrder := func(order *Order) (stop bool) {
		orders = append(orders, order)
		return false
	}
	k.IterateOrder(ctx, appendOrder)
	return orders
}

func (k Keeper) IterateOrder(ctx sdk.Context, process func(*Order) bool) {
	store := ctx.KVStore(k.marketKey)
	iter := sdk.KVStorePrefixIterator(store, OrderBookKeyPrefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		or := k.decodeOrder(iter.Value())
		if process(&or) {
			return
		}
		iter.Next()
	}
}

func (k Keeper) decodeOrder(bz []byte) (order Order) {
	if err := k.cdc.UnmarshalBinaryBare(bz, &order); err != nil {
		panic(err)
	}
	return
}

func (k Keeper) GetAllMarketInfos(ctx sdk.Context) []MarketInfo {
	var infos []MarketInfo
	appendMarket := func(order MarketInfo) (stop bool) {
		infos = append(infos, order)
		return false
	}
	k.IterateMarket(ctx, appendMarket)
	return infos
}

func (k Keeper) IterateMarket(ctx sdk.Context, process func(info MarketInfo) bool) {
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

func (k Keeper) GetMarketInfo(ctx sdk.Context, symbol string) (info MarketInfo, err error) {
	store := ctx.KVStore(k.marketKey)
	value := store.Get(marketStoreKey(MarketIdentifierPrefix, symbol))

	//TODO. will modify, because the function maybe panic in case of error .
	err = k.cdc.UnmarshalBinaryBare(value, &info)
	return
}

func (k Keeper) decodeMarket(bz []byte) (info MarketInfo) {
	if err := k.cdc.UnmarshalBinaryBare(bz, &info); err != nil {
		panic(err)
	}
	return
}
