package market

import (
	"bytes"
	"fmt"
	"sort"
	"testing"

	"github.com/coinexchain/dex/modules/msgqueue"
	"github.com/coinexchain/dex/types"

	"github.com/coinexchain/dex/modules/asset"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

type mocBankxKeeper struct {
	records []string
}

func (k *mocBankxKeeper) SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	return nil
}
func (k *mocBankxKeeper) DeductFee(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	return nil
}
func (k *mocBankxKeeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return true
}
func (k *mocBankxKeeper) SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amt sdk.Coins) sdk.Error {
	k.records = append(k.records, fmt.Sprintf("send %s %s from %s to %s",
		amt[0].Amount.String(), amt[0].Denom, from.String(), to.String()))
	return nil
}
func (k *mocBankxKeeper) FreezeCoins(ctx sdk.Context, acc sdk.AccAddress, amt sdk.Coins) sdk.Error {
	k.records = append(k.records, fmt.Sprintf("freeze %s %s at %s",
		amt[0].Amount.String(), amt[0].Denom, string(acc)))
	return nil
}
func (k *mocBankxKeeper) UnFreezeCoins(ctx sdk.Context, acc sdk.AccAddress, amt sdk.Coins) sdk.Error {
	k.records = append(k.records, fmt.Sprintf("unfreeze %s %s at %s",
		amt[0].Amount.String(), amt[0].Denom, acc.String()))
	return nil
}

type mocAssertStatusKeeper struct {
	forbiddenDenomList       []string
	globalForbiddenDenomList []string
	forbiddenAddrList        []sdk.AccAddress
}

func (k *mocAssertStatusKeeper) IsTokenForbidden(ctx sdk.Context, denom string) bool {
	for _, d := range k.globalForbiddenDenomList {
		if denom == d {
			return true
		}
	}
	return false
}
func (k *mocAssertStatusKeeper) IsTokenExists(ctx sdk.Context, denom string) bool {
	return true
}
func (k *mocAssertStatusKeeper) IsTokenIssuer(ctx sdk.Context, denom string, addr sdk.AccAddress) bool {
	return false
}
func (k *mocAssertStatusKeeper) IsForbiddenByTokenIssuer(ctx sdk.Context, denom string, addr sdk.AccAddress) bool {
	for i := 0; i < len(k.forbiddenDenomList); i++ {
		if denom == k.forbiddenDenomList[i] && bytes.Equal(addr, k.forbiddenAddrList[i]) {
			return true
		}
	}
	return false
}
func (k *mocAssertStatusKeeper) GetToken(ctx sdk.Context, symbol string) asset.Token {
	return nil
}

type mockFeeColletKeeper struct {
	records []string
}

func (k *mockFeeColletKeeper) SubtractFeeAndCollectFee(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	fee := fmt.Sprintf("addr : %s, fee : %s", addr, amt.String())
	k.records = append(k.records, fee)
	return nil
}

func TestUnfreezeCoinsForOrder(t *testing.T) {
	bxKeeper := &mocBankxKeeper{records: make([]string, 0, 10)}
	mockFeeK := &mockFeeColletKeeper{}
	order := newTO("00001", 1, 11051, 50, Buy, GTE, 998)
	order.Freeze = 50
	order.FrozenFee = 10
	order.DealStock = 20
	ctx, _ := newContextAndMarketKey()
	unfreezeCoinsForOrder(ctx, bxKeeper, order, 0, mockFeeK)
	refout := "unfreeze 50 usdt at cosmos1qy352eufqy352eufqy352eufqy35qqqptw34ca"
	if refout != bxKeeper.records[0] {
		t.Errorf("Error in unfreezeCoinsForOrder")
	}

	coinFee := sdk.NewDec(order.DealStock).Mul(sdk.NewDec(order.FrozenFee)).Quo(sdk.NewDec(order.Quantity)).RoundInt64()
	refout = fmt.Sprintf("addr : %s, fee : %s", order.Sender, types.NewCetCoin(coinFee).String())
	if refout != mockFeeK.records[0] {
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
	keeper := NewKeeper(keys.marketKey, axk, bnk, mockFeeKeeper{}, msgCdc, msgqueue.NewProducer(), subspace)
	currDay := ctx.BlockHeader().Time.Day()
	keeper.orderClean.SetDay(ctx, currDay-1)
	parameters := Params{}
	parameters.GTEOrderLifetime = 1
	parameters.MaxExecutedPriceChangeRatio = DefaultMaxExecutedPriceChangeRatio
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
		"unfreeze 55 usdt at cosmos1qy352eufqy352eufqy352eufqy35qqqptw34ca",
		"unfreeze 120 btc at cosmos1qy352eufqy352eufqy352eufqy35qqq9yynnh8",
		"unfreeze 55 usdt at cosmos1qy352eufqy352eufqy352eufqy35qqqz9ayrkz",
		"unfreeze 54 usdt at cosmos1qy352eufqy352eufqy352eufqy35qqqz9ayrkz",
		"unfreeze 60 cet at cosmos1qy352eufqy352eufqy352eufqy35qqqyej8x24",
	}
	for i, rec := range bnk.records {
		if records[i] != rec {
			t.Errorf("Error in Removing Old Orders!")
		}
		fmt.Printf("%s\n", rec)
	}
}

func TestDelist(t *testing.T) {
	axk := &mocAssertStatusKeeper{}
	axk.forbiddenDenomList = []string{"bch"}
	axk.globalForbiddenDenomList = []string{"bsv"}
	addr01, _ := simpleAddr("00001")
	axk.forbiddenAddrList = []sdk.AccAddress{addr01}
	bnk := &mocBankxKeeper{}
	ctx, keys := newContextAndMarketKey()
	subspace := params.NewKeeper(msgCdc, keys.keyParams, keys.tkeyParams).Subspace(StoreKey)
	keeper := NewKeeper(keys.marketKey, axk, bnk, mockFeeKeeper{}, msgCdc, msgqueue.NewProducer(), subspace)
	delistKeeper := NewDelistKeeper(keys.marketKey)
	delistKeeper.AddDelistRequest(ctx, ctx.BlockHeight(), "btc/usdt")
	currDay := ctx.BlockHeader().Time.Day()
	keeper.orderClean.SetDay(ctx, currDay)
	parameters := Params{}
	parameters.GTEOrderLifetime = 1
	parameters.MaxExecutedPriceChangeRatio = DefaultMaxExecutedPriceChangeRatio
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
	keeper.SetMarket(ctx, MarketInfo{
		Stock:             "bch",
		Money:             "usdt",
		Creator:           creater2,
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(0),
	})
	keeper.SetMarket(ctx, MarketInfo{
		Stock:             "bsv",
		Money:             "cet",
		Creator:           creater1,
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(0),
	})

	cetKeeper := NewOrderKeeper(keeper.marketKey, "cet/usdt", msgCdc)
	btcKeeper := NewOrderKeeper(keeper.marketKey, "btc/usdt", msgCdc)
	orders := make([]*Order, 10)
	orders[0] = newTO("00001", 1, 11051, 60, Buy, GTE, 98)
	orders[0].Symbol = "btc/usdt"
	btcKeeper.Add(ctx, orders[0])
	orders[1] = newTO("00005", 5, 12039, 120, Sell, IOC, 1000)
	orders[1].Symbol = "btc/usdt"
	btcKeeper.Add(ctx, orders[1])
	orders[2] = newTO("00020", 6, 11039, 100, Sell, IOC, 1000)
	orders[2].Symbol = "btc/usdt"
	btcKeeper.Add(ctx, orders[2])

	orders[3] = newTO("00002", 2, 11080, 50, Buy, GTE, 98)
	cetKeeper.Add(ctx, orders[3])
	orders[4] = newTO("00002", 3, 10900, 50, Buy, GTE, 92)
	cetKeeper.Add(ctx, orders[4])
	orders[5] = newTO("00004", 4, 11032, 30, Sell, GTE, 90)
	cetKeeper.Add(ctx, orders[5])
	orders[6] = newTO("00009", 9, 11032, 30, Sell, GTE, 90)
	cetKeeper.Add(ctx, orders[6])
	orders[7] = newTO("00002", 8, 11085, 5, Buy, GTE, 98)
	cetKeeper.Add(ctx, orders[7])

	orders[8] = newTO("00001", 10, 11000, 15, Buy, GTE, 998)
	orders[8].Symbol = "bch/usdt"
	cetKeeper.Add(ctx, orders[8])

	orders[9] = newTO("00001", 7, 11000, 15, Buy, GTE, 998)
	orders[9].Symbol = "bsv/usdt"
	cetKeeper.Add(ctx, orders[9])

	EndBlocker(ctx, keeper)
	gKeeper := NewGlobalOrderKeeper(keeper.marketKey, msgCdc)
	allOrders := gKeeper.GetAllOrders(ctx)
	subList := []int{6, 8, 9, 4}
	if len(allOrders) != 4 {
		t.Errorf("Incorrect remain orders.")
	}
	for i, sub := range subList {
		if !sameTO(orders[sub], allOrders[i]) {
			t.Errorf("Incorrect remain orders.")
		}
	}
	for _, order := range allOrders {
		fmt.Printf("Remain: %s\n", order.OrderID())
	}
	records := []string{
		"unfreeze 60 btc at cosmos1qy352eufqy352eufqy352eufqy35qqpqe926wf",
		"send 60 btc from cosmos1qy352eufqy352eufqy352eufqy35qqpqe926wf to cosmos1qy352eufqy352eufqy352eufqy35qqqptw34ca",
		"unfreeze 66 usdt at cosmos1qy352eufqy352eufqy352eufqy35qqqptw34ca",
		"send 66 usdt from cosmos1qy352eufqy352eufqy352eufqy35qqqptw34ca to cosmos1qy352eufqy352eufqy352eufqy35qqpqe926wf",
		"unfreeze 5 cet at cosmos1qy352eufqy352eufqy352eufqy35qqqyej8x24",
		"send 5 cet from cosmos1qy352eufqy352eufqy352eufqy35qqqyej8x24 to cosmos1qy352eufqy352eufqy352eufqy35qqqz9ayrkz",
		"unfreeze 6 usdt at cosmos1qy352eufqy352eufqy352eufqy35qqqz9ayrkz",
		"send 6 usdt from cosmos1qy352eufqy352eufqy352eufqy35qqqz9ayrkz to cosmos1qy352eufqy352eufqy352eufqy35qqqyej8x24",
		"unfreeze 25 cet at cosmos1qy352eufqy352eufqy352eufqy35qqqyej8x24",
		"send 25 cet from cosmos1qy352eufqy352eufqy352eufqy35qqqyej8x24 to cosmos1qy352eufqy352eufqy352eufqy35qqqz9ayrkz",
		"unfreeze 28 usdt at cosmos1qy352eufqy352eufqy352eufqy35qqqz9ayrkz",
		"send 28 usdt from cosmos1qy352eufqy352eufqy352eufqy35qqqz9ayrkz to cosmos1qy352eufqy352eufqy352eufqy35qqqyej8x24",
	}
	unfreezeRecords := []string{
		"unfreeze 40 btc at cosmos1qy352eufqy352eufqy352eufqy35qqpqe926wf",
		"unfreeze 120 btc at cosmos1qy352eufqy352eufqy352eufqy35qqq9yynnh8",
		"unfreeze 27 usdt at cosmos1qy352eufqy352eufqy352eufqy35qqqz9ayrkz",
	}
	for i, rec := range records {
		if rec != bnk.records[i] {
			t.Errorf("Incorrect Records, actual : %s, expect : %s", bnk.records[i], rec)
		}
		fmt.Printf("%s\n", rec)
	}
	unf := bnk.records[len(records):]
	sort.Strings(unf)
	sort.Strings(unfreezeRecords)
	for i, rec := range unfreezeRecords {
		if rec != unf[i] {
			t.Errorf("Incorrect Records, actual : %s, expect : %s", unf[i], rec)
		}
		fmt.Printf("%s\n", rec)
	}
}
