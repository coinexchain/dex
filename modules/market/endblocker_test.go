package market

import (
	"bytes"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

type mocBankxKeeper struct {
	records []string
}

func (k *mocBankxKeeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return true
}
func (k *mocBankxKeeper) SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amt sdk.Coins) sdk.Error {
	k.records = append(k.records, fmt.Sprintf("send %s from %s to %s",
		amt[0].Amount.String(), string(from), string(to)))
	return nil
}
func (k *mocBankxKeeper) FreezeCoins(ctx sdk.Context, acc sdk.AccAddress, amt sdk.Coins) sdk.Error {
	k.records = append(k.records, fmt.Sprintf("freeze %s at %s",
		amt[0].Amount.String(), string(acc)))
	return nil
}
func (k *mocBankxKeeper) UnFreezeCoins(ctx sdk.Context, acc sdk.AccAddress, amt sdk.Coins) sdk.Error {
	k.records = append(k.records, fmt.Sprintf("unfreeze %s at %s",
		amt[0].Amount.String(), acc.String()))
	return nil
}

type mocAssertStatusKeeper struct {
	frzDenomList []string
	frzAddrList  []sdk.AccAddress
}

func (k *mocAssertStatusKeeper) IsTokenFrozen(ctx sdk.Context, denom string) bool {
	return false
}
func (k *mocAssertStatusKeeper) IsTokenExists(ctx sdk.Context, denom string) bool {
	return true
}
func (k *mocAssertStatusKeeper) IsTokenIssuer(ctx sdk.Context, denom string, addr sdk.AccAddress) bool {
	return false
}
func (k *mocAssertStatusKeeper) IsForbiddenByTokenIssuer(ctx sdk.Context, denom string, addr sdk.AccAddress) bool {
	for i := 0; i < len(k.frzDenomList); i++ {
		if denom == k.frzDenomList[i] && bytes.Equal(addr, k.frzAddrList[i]) {
			return true
		}
	}
	return false
}

func TestUnfreezeCoinsForOrder(t *testing.T) {
	bxKeeper := &mocBankxKeeper{records: make([]string, 0, 10)}
	order := newTO("00001", 1, 11051, 50, Buy, GTE, 998)
	order.Freeze = 50
	ctx, _ := newContextAndMarketKey()
	unfreezeCoinsForOrder(ctx, bxKeeper, order)
	refout := "unfreeze 50 at cosmos1qy352eufqy352eufqy352eufqy35qqqptw34ca"
	if refout != bxKeeper.records[0] {
		t.Errorf("Error in unfreezeCoinsForOrder")
	}
}

//func NewKeeper(key sdk.StoreKey, axkVal ExpectedAssertStatusKeeper,
//	bnkVal ExpectedBankxKeeper, cdcVal *codec.Codec, paramstore params.Subspace) Keeper {

func TestRemoveOrders(t *testing.T) {
	axk := &mocAssertStatusKeeper{}
	bnk := &mocBankxKeeper{}
	ctx, keys := newContextAndMarketKey()
	subspace := params.NewKeeper(msgCdc, keys.keyParams, keys.tkeyParams).Subspace(StoreKey)
	keeper := NewKeeper(keys.marketKey, axk, bnk, msgCdc, subspace)
	currDay := ctx.BlockHeader().Time.Day()
	keeper.orderClean.SetDay(ctx, currDay-1)
	parameters := Params{}
	parameters.GTEOrderLifetime = 1
	keeper.SetParams(ctx, parameters)

	creater1, _ := simpleAddr("10000")
	creater2, _ := simpleAddr("11000")
	keeper.SetMarket(ctx, MarketInfo{
		Stock:             "cet",
		Money:             "usdt",
		Creator:           creater1,
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(0),
	})
	keeper.SetMarket(ctx, MarketInfo{
		Stock:             "btc",
		Money:             "usdt",
		Creator:           creater2,
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(0),
	})

	cetKeeper := NewOrderKeeper(keeper.marketKey, "cet/usdt", msgCdc)
	btcKeeper := NewOrderKeeper(keeper.marketKey, "btc/usdt", msgCdc)
	order := newTO("00001", 1, 11051, 50, Buy, GTE, 98)
	order.Symbol = "btc/usdt"
	btcKeeper.Add(ctx, order)
	order = newTO("00005", 5, 12039, 120, Sell, GTE, 96)
	order.Symbol = "btc/usdt"
	btcKeeper.Add(ctx, order)
	order = newTO("00002", 2, 11080, 50, Buy, GTE, 98)
	cetKeeper.Add(ctx, order)
	order = newTO("00002", 3, 10900, 50, Buy, GTE, 92)
	cetKeeper.Add(ctx, order)
	order = newTO("00004", 4, 11032, 60, Sell, GTE, 90)
	cetKeeper.Add(ctx, order)

	EndBlocker(ctx, keeper)
	gKeeper := NewGlobalOrderKeeper(keeper.marketKey, msgCdc)
	allOrders := gKeeper.GetAllOrders(ctx)
	if len(allOrders) != 0 {
		t.Errorf("Error in Removing Old Orders!")
	}
	records := []string{
		"unfreeze 50 at cosmos1qy352eufqy352eufqy352eufqy35qqqptw34ca",
		"unfreeze 120 at cosmos1qy352eufqy352eufqy352eufqy35qqq9yynnh8",
		"unfreeze 50 at cosmos1qy352eufqy352eufqy352eufqy35qqqz9ayrkz",
		"unfreeze 50 at cosmos1qy352eufqy352eufqy352eufqy35qqqz9ayrkz",
		"unfreeze 60 at cosmos1qy352eufqy352eufqy352eufqy35qqqyej8x24",
	}
	for i, rec := range bnk.records {
		if records[i] != rec {
			t.Errorf("Error in Removing Old Orders!")
		}
	}
}

//func (keeper *DelistKeeper) AddDelistRequest(ctx sdk.Context, height int64, symbol string) {
//func (keeper *DelistKeeper) GetDelistSymbolsAtHeight(ctx sdk.Context, height int64) []string {
//func NewDelistKeeper(key sdk.StoreKey) *DelistKeeper {
func TestDelist(t *testing.T) {
	axk := &mocAssertStatusKeeper{}
	bnk := &mocBankxKeeper{}
	ctx, keys := newContextAndMarketKey()
	subspace := params.NewKeeper(msgCdc, keys.keyParams, keys.tkeyParams).Subspace(StoreKey)
	keeper := NewKeeper(keys.marketKey, axk, bnk, msgCdc, subspace)
	delistKeeper := NewDelistKeeper(keys.marketKey)
	delistKeeper.AddDelistRequest(ctx, ctx.BlockHeight(), "btc/usdt")
	currDay := ctx.BlockHeader().Time.Day()
	keeper.orderClean.SetDay(ctx, currDay)

	creater1, _ := simpleAddr("10000")
	creater2, _ := simpleAddr("11000")
	keeper.SetMarket(ctx, MarketInfo{
		Stock:             "cet",
		Money:             "usdt",
		Creator:           creater1,
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(0),
	})
	keeper.SetMarket(ctx, MarketInfo{
		Stock:             "btc",
		Money:             "usdt",
		Creator:           creater2,
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(0),
	})

	cetKeeper := NewOrderKeeper(keeper.marketKey, "cet/usdt", msgCdc)
	btcKeeper := NewOrderKeeper(keeper.marketKey, "btc/usdt", msgCdc)
	order := newTO("00001", 1, 11051, 50, Buy, GTE, 98)
	order.Symbol = "btc/usdt"
	btcKeeper.Add(ctx, order)
	order = newTO("00005", 5, 12039, 120, Sell, GTE, 96)
	order.Symbol = "btc/usdt"
	btcKeeper.Add(ctx, order)
	order = newTO("00002", 2, 11080, 50, Buy, GTE, 98)
	cetKeeper.Add(ctx, order)
	order = newTO("00002", 3, 10900, 50, Buy, GTE, 92)
	cetKeeper.Add(ctx, order)
	order = newTO("00004", 4, 11032, 60, Sell, GTE, 90)
	cetKeeper.Add(ctx, order)

	EndBlocker(ctx, keeper)
	gKeeper := NewGlobalOrderKeeper(keeper.marketKey, msgCdc)
	allOrders := gKeeper.GetAllOrders(ctx)
	for _, order := range allOrders {
		fmt.Printf("Remain: %s\n", order.OrderID())
	}
}
