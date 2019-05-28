package market

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

var (
	marketIdetifierPrefix     = []byte{0x01}
	orderBookIdetifierPrefix  = []byte{0x02}
	orderQueueIdetifierPrefix = []byte{0x03}
	askListIdetifierPrefix    = []byte{0x04}
	bidListIdetifierPrefix    = []byte{0x05}
)

type Keeper struct {
	paramSubspace params.Subspace
	marketKey     sdk.StoreKey
	axk           ExpectedAssertStatusKeeper
	bnk           ExpectedBankxKeeper
	cdc           *codec.Codec
	//fek       incentive.FeeCollectionKeeper
}

func NewKeeper(key sdk.StoreKey, axkVal ExpectedAssertStatusKeeper,
	bnkVal ExpectedBankxKeeper, cdcVal *codec.Codec, paramstore params.Subspace) Keeper {
	return Keeper{marketKey: key, axk: axkVal, bnk: bnkVal,
		cdc: cdcVal, paramSubspace: paramstore.WithKeyTable(ParamKeyTable())}
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
func (k Keeper) SetOrder(ctx sdk.Context, order Order) sdk.Error {
	store := ctx.KVStore(k.marketKey)
	bz, err := k.cdc.MarshalBinaryBare(order)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}
	store.Set(marketStoreKey(orderBookIdetifierPrefix, order.Symbol, order.OrderID()), bz)
	return nil
}

// SetMarket implements token Keeper.
func (k Keeper) SetMarket(ctx sdk.Context, info MarketInfo) sdk.Error {
	store := ctx.KVStore(k.marketKey)
	bz, err := k.cdc.MarshalBinaryBare(info)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}
	store.Set(marketStoreKey(marketIdetifierPrefix, info.Stock+SymbolSeparator+info.Money), bz)
	return nil
}

// RegisterCodec registers concrete types on the codec
func (k Keeper) RegisterCodec() {
	registerCodec(k.cdc)
}

func registerCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(Order{}, "cet-chain/Order", nil)
	cdc.RegisterConcrete(MarketInfo{}, "cet-chain/MarketInfo", nil)
}

func (k Keeper) GetAllTokens(ctx sdk.Context) []Order {
	return nil
}

func (k Keeper) GetAllMarketInfos(ctx sdk.Context) []MarketInfo {
	return nil
}
