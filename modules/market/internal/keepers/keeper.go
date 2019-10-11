package keepers

import (
	"bytes"
	"errors"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/market/internal/types"
	"github.com/coinexchain/dex/msgqueue"
)

type Keeper struct {
	paramSubspace params.Subspace
	marketKey     sdk.StoreKey
	cdc           *codec.Codec
	axk           types.ExpectedAssetStatusKeeper
	bnk           types.ExpectedBankxKeeper
	ock           *OrderCleanUpDayKeeper
	gmk           GlobalMarketInfoKeeper
	msgProducer   msgqueue.MsgSender
	bancorK       types.ExpectedBancorKeeper
	ak            auth.AccountKeeper
}

func NewKeeper(key sdk.StoreKey, axkVal types.ExpectedAssetStatusKeeper,
	bnkVal types.ExpectedBankxKeeper, cdcVal *codec.Codec,
	msgKeeperVal msgqueue.MsgSender,
	paramstore params.Subspace,
	bancor types.ExpectedBancorKeeper,
	ak auth.AccountKeeper) Keeper {

	return Keeper{
		paramSubspace: paramstore.WithKeyTable(types.ParamKeyTable()),
		marketKey:     key,
		cdc:           cdcVal,
		axk:           axkVal,
		bnk:           bnkVal,
		ock:           NewOrderCleanUpDayKeeper(key),
		gmk:           NewGlobalMarketInfoKeeper(key, cdcVal),
		msgProducer:   msgKeeperVal,
		bancorK:       bancor,
		ak:            ak,
	}
}

func (k Keeper) QuerySeqWithAddr(ctx sdk.Context, addr sdk.AccAddress) (uint64, sdk.Error) {
	bz, err := k.cdc.MarshalJSON(auth.QueryAccountParams{Address: addr})
	if err != nil {
		return 0, types.ErrInvalidAddress()
	}
	res, sdkErr := auth.NewQuerier(k.ak)(ctx, []string{auth.QueryAccount}, abci.RequestQuery{
		Data: bz,
	})
	if sdkErr != nil {
		return 0, sdkErr
	}

	var acc auth.Account
	if err := k.cdc.UnmarshalJSON(res, &acc); err != nil {
		return 0, types.ErrFailedUnmarshal()
	}
	return acc.GetSequence(), nil
}

func (k Keeper) IsBancorExist(ctx sdk.Context, stock string) bool {
	return k.bancorK.IsBancorExist(ctx, stock)
}

func (k Keeper) GetToken(ctx sdk.Context, symbol string) asset.Token {
	return k.axk.GetToken(ctx, symbol)
}

func (k Keeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return k.bnk.HasCoins(ctx, addr, amt)
}

func (k Keeper) FreezeCoins(ctx sdk.Context, acc sdk.AccAddress, amt sdk.Coins) sdk.Error {
	return k.bnk.FreezeCoins(ctx, acc, amt)
}

func (k Keeper) IsTokenIssuer(ctx sdk.Context, denom string, addr sdk.AccAddress) bool {
	return k.axk.IsTokenIssuer(ctx, denom, addr)
}

func (k Keeper) IsTokenExists(ctx sdk.Context, symbol string) bool {
	return k.axk.IsTokenExists(ctx, symbol)
}

func (k Keeper) IsSubScribed(topic string) bool {
	return k.msgProducer.IsSubscribed(topic)
}

func (k Keeper) IsForbiddenByTokenIssuer(ctx sdk.Context, denom string, addr sdk.AccAddress) bool {
	return k.axk.IsForbiddenByTokenIssuer(ctx, denom, addr)
}

func (k Keeper) IsTokenForbidden(ctx sdk.Context, symbol string) bool {
	return k.axk.IsTokenForbidden(ctx, symbol)
}
func (k Keeper) GetOrderCleanTime(ctx sdk.Context) int64 {
	return k.ock.GetUnixTime(ctx)
}

func (k Keeper) SetOrderCleanTime(ctx sdk.Context, t int64) {
	k.ock.SetUnixTime(ctx, t)
}

func (k Keeper) GetBankxKeeper() types.ExpectedBankxKeeper {
	return k.bnk
}

func (k Keeper) GetAssetKeeper() types.ExpectedAssetStatusKeeper {
	return k.axk
}

func (k Keeper) GetMarketKey() sdk.StoreKey {
	return k.marketKey
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the asset module's parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the asset module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSubspace.GetParamSet(ctx, &params)
	return
}

func (k Keeper) GetMarketFeeMin(ctx sdk.Context) int64 {
	return k.GetParams(ctx).MarketFeeMin
}

// -----------------------------------------------------------------------------
// Order

// SetOrder implements token Keeper.
func (k Keeper) SetOrder(ctx sdk.Context, order *types.Order) sdk.Error {
	return NewOrderKeeper(k.marketKey, order.TradingPair, k.cdc).Add(ctx, order)
}

func (k Keeper) GetAllOrders(ctx sdk.Context) []*types.Order {
	return NewGlobalOrderKeeper(k.marketKey, k.cdc).GetAllOrders(ctx)
}

// -----------------------------------------------
// market info

func (k Keeper) SetMarket(ctx sdk.Context, info types.MarketInfo) sdk.Error {
	return k.gmk.SetMarket(ctx, info)
}

func (k Keeper) RemoveMarket(ctx sdk.Context, symbol string) sdk.Error {
	return k.gmk.RemoveMarket(ctx, symbol)
}

func (k Keeper) GetAllMarketInfos(ctx sdk.Context) []types.MarketInfo {
	return k.gmk.GetAllMarketInfos(ctx)
}

func (k Keeper) GetMarketInfo(ctx sdk.Context, symbol string) (types.MarketInfo, error) {
	return k.gmk.GetMarketInfo(ctx, symbol)
}

func (k Keeper) SubtractFeeAndCollectFee(ctx sdk.Context, addr sdk.AccAddress, amt int64) sdk.Error {
	return k.bnk.DeductInt64CetFee(ctx, addr, amt)
}

func (k Keeper) MarketOwner(ctx sdk.Context, info types.MarketInfo) sdk.AccAddress {
	return k.axk.GetToken(ctx, info.Stock).GetOwner()
}

func (k *Keeper) GetMarketLastExePrice(ctx sdk.Context, symbol string) (sdk.Dec, error) {
	mi, err := k.GetMarketInfo(ctx, symbol)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	return mi.LastExecutedPrice, err
}

func (k *Keeper) IsMarketExist(ctx sdk.Context, symbol string) bool {
	_, err := k.GetMarketInfo(ctx, symbol)
	return err == nil
}

// -----------------------------------------------------------------------------

type GlobalMarketInfoKeeper interface {
	SetMarket(ctx sdk.Context, info types.MarketInfo) sdk.Error
	RemoveMarket(ctx sdk.Context, symbol string) sdk.Error
	GetAllMarketInfos(ctx sdk.Context) []types.MarketInfo
	GetMarketInfo(ctx sdk.Context, symbol string) (types.MarketInfo, error)
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
func (k PersistentMarketInfoKeeper) SetMarket(ctx sdk.Context, info types.MarketInfo) sdk.Error {
	store := ctx.KVStore(k.marketKey)
	bz, err := k.cdc.MarshalBinaryBare(info)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}
	store.Set(marketStoreKey(MarketIdentifierPrefix, info.GetSymbol()), bz)
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

func (k PersistentMarketInfoKeeper) GetAllMarketInfos(ctx sdk.Context) []types.MarketInfo {
	var infos []types.MarketInfo
	appendMarket := func(order types.MarketInfo) (stop bool) {
		infos = append(infos, order)
		return false
	}
	k.iterateMarket(ctx, appendMarket)
	return infos
}

func (k PersistentMarketInfoKeeper) iterateMarket(ctx sdk.Context, process func(info types.MarketInfo) bool) {
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

func (k PersistentMarketInfoKeeper) GetMarketInfo(ctx sdk.Context, symbol string) (info types.MarketInfo, err error) {
	store := ctx.KVStore(k.marketKey)
	value := store.Get(marketStoreKey(MarketIdentifierPrefix, symbol))
	if len(value) == 0 {
		err = errors.New("No such market exist: " + symbol)
		return
	}
	err = k.cdc.UnmarshalBinaryBare(value, &info)
	return
}

func (k PersistentMarketInfoKeeper) decodeMarket(bz []byte) (info types.MarketInfo) {
	if err := k.cdc.UnmarshalBinaryBare(bz, &info); err != nil {
		panic(err)
	}
	return
}

func marketStoreKey(prefix []byte, params ...string) []byte {
	buf := bytes.NewBuffer(prefix)
	for _, param := range params {
		if _, err := buf.Write([]byte(param)); err != nil {
			panic(err)
		}
	}
	return buf.Bytes()
}
