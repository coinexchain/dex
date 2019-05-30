package market

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

var (
	MarketIdentifierPrefix = []byte{0x15}
)

type Keeper struct {
	paramSubspace params.Subspace
	marketKey     sdk.StoreKey
	axk           ExpectedAssertStatusKeeper
	bnk           ExpectedBankxKeeper
	cdc           *codec.Codec
	ork           OrderKeeper
	orderClean    *OrderCleanUpDayKeeper
	//fek       incentive.FeeCollectionKeeper
}

func NewKeeper(key sdk.StoreKey, axkVal ExpectedAssertStatusKeeper,
	bnkVal ExpectedBankxKeeper, cdcVal *codec.Codec, paramstore params.Subspace) Keeper {

	return Keeper{marketKey: key, axk: axkVal, bnk: bnkVal,
		paramSubspace: paramstore.WithKeyTable(ParamKeyTable()),
		cdc:           cdcVal, ork: NewOrderKeeper(key, MarketKey, cdcVal),
		orderClean: NewOrderCleanUpDayKeeper(key)}
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
	return k.ork.Add(ctx, order)
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
	return k.ork.GetAllOrders(ctx)
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

func (k Keeper) decodeMarket(bz []byte) (info MarketInfo) {
	if err := k.cdc.UnmarshalBinaryBare(bz, &info); err != nil {
		panic(err)
	}
	return
}
