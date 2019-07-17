package market

import (
	"bytes"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/market/internal/keepers"
	types2 "github.com/coinexchain/dex/modules/market/internal/types"

	"github.com/coinexchain/dex/modules/msgqueue"
	"github.com/coinexchain/dex/types"

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
	order := keepers.newTO("00001", 1, 11051, 50, keepers.Buy, keepers.GTE, 998)
	order.Freeze = 50
	order.FrozenFee = 10
	order.DealStock = 20
	ctx, _ := keepers.newContextAndMarketKey(testNetSubString)
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

func TestRemoveOrders(t *testing.T) {
	axk := &mocAssertStatusKeeper{}
	bnk := &mocBankxKeeper{}
	ctx, keys := keepers.newContextAndMarketKey(testNetSubString)
	subspace := params.NewKeeper(msgCdc, keys.keyParams, keys.tkeyParams).Subspace(types2.StoreKey)
	keeper := keepers.NewKeeper(keys.marketKey, axk, bnk, msgCdc, msgqueue.NewProducer(), subspace)
	keeper.orderClean.SetUnixTime(ctx, time.Now().Unix())
	ctx = ctx.WithBlockTime(time.Unix(time.Now().Unix()+int64(25*60*60), 0))
	parameters := keepers.Params{}
	parameters.GTEOrderLifetime = 1
	parameters.MaxExecutedPriceChangeRatio = keepers.DefaultMaxExecutedPriceChangeRatio
	keeper.SetParams(ctx, parameters)

	keeper.SetMarket(ctx, types2.MarketInfo{
		Stock:             "cet",
		Money:             "usdt",
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(0),
	})
	keeper.SetMarket(ctx, types2.MarketInfo{
		Stock:             "btc",
		Money:             "usdt",
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(0),
	})

	cetKeeper := keepers.NewOrderKeeper(keeper.marketKey, "cet/usdt", msgCdc)
	btcKeeper := keepers.NewOrderKeeper(keeper.marketKey, "btc/usdt", msgCdc)
	order := keepers.newTO("00001", 1, 11051, 50, keepers.Buy, keepers.GTE, 98)
	order.TradingPair = "btc/usdt"
	btcKeeper.Add(ctx, order)
	order = keepers.newTO("00005", 5, 12039, 120, keepers.Sell, keepers.GTE, 96)
	order.TradingPair = "btc/usdt"
	btcKeeper.Add(ctx, order)
	order = keepers.newTO("00002", 2, 11080, 50, keepers.Buy, keepers.GTE, 98)
	cetKeeper.Add(ctx, order)
	order = keepers.newTO("00002", 3, 10900, 50, keepers.Buy, keepers.GTE, 92)
	cetKeeper.Add(ctx, order)
	order = keepers.newTO("00004", 4, 11032, 60, keepers.Sell, keepers.GTE, 90)
	cetKeeper.Add(ctx, order)

	EndBlocker(ctx, keeper)
	gKeeper := keepers.NewGlobalOrderKeeper(keeper.marketKey, msgCdc)
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
	}
}

func TestDelist(t *testing.T) {
	axk := &mocAssertStatusKeeper{}
	axk.forbiddenDenomList = []string{"bch"}
	axk.globalForbiddenDenomList = []string{"bsv"}
	addr01, _ := keepers.simpleAddr("00001")
	axk.forbiddenAddrList = []sdk.AccAddress{addr01}
	bnk := &mocBankxKeeper{}
	ctx, keys := keepers.newContextAndMarketKey(testNetSubString)
	subspace := params.NewKeeper(msgCdc, keys.keyParams, keys.tkeyParams).Subspace(types2.StoreKey)
	keeper := keepers.NewKeeper(keys.marketKey, axk, bnk, msgCdc, msgqueue.NewProducer(), subspace)
	delistKeeper := keepers.NewDelistKeeper(keys.marketKey)
	delistKeeper.AddDelistRequest(ctx, ctx.BlockHeight(), "btc/usdt")
	// currDay := ctx.BlockHeader().Time.Unix()
	// keeper.orderClean.SetUnixTime(ctx, currDay)
	keeper.orderClean.SetUnixTime(ctx, time.Now().Unix())
	ctx = ctx.WithBlockTime(time.Now())
	parameters := keepers.Params{}
	parameters.GTEOrderLifetime = 1
	parameters.MaxExecutedPriceChangeRatio = keepers.DefaultMaxExecutedPriceChangeRatio
	keeper.SetParams(ctx, parameters)

	keeper.SetMarket(ctx, types2.MarketInfo{
		Stock:             "cet",
		Money:             "usdt",
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(0),
	})
	keeper.SetMarket(ctx, types2.MarketInfo{
		Stock:             "btc",
		Money:             "usdt",
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(0),
	})
	keeper.SetMarket(ctx, types2.MarketInfo{
		Stock:             "bch",
		Money:             "usdt",
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(0),
	})
	keeper.SetMarket(ctx, types2.MarketInfo{
		Stock:             "bsv",
		Money:             "cet",
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(0),
	})

	cetKeeper := keepers.NewOrderKeeper(keeper.marketKey, "cet/usdt", msgCdc)
	btcKeeper := keepers.NewOrderKeeper(keeper.marketKey, "btc/usdt", msgCdc)
	orders := make([]*types2.Order, 10)
	orders[0] = keepers.newTO("00001", 1, 11051, 60, keepers.Buy, keepers.GTE, 98)
	orders[0].TradingPair = "btc/usdt"
	btcKeeper.Add(ctx, orders[0])
	orders[1] = keepers.newTO("00005", 5, 12039, 120, keepers.Sell, keepers.IOC, 1000)
	orders[1].TradingPair = "btc/usdt"
	btcKeeper.Add(ctx, orders[1])
	orders[2] = keepers.newTO("00020", 6, 11039, 100, keepers.Sell, keepers.IOC, 1000)
	orders[2].TradingPair = "btc/usdt"
	btcKeeper.Add(ctx, orders[2])

	orders[3] = keepers.newTO("00202", 2, 11080, 50, keepers.Buy, keepers.GTE, 98)
	cetKeeper.Add(ctx, orders[3])
	orders[4] = keepers.newTO("00102", 3, 10900, 50, keepers.Buy, keepers.GTE, 92)
	cetKeeper.Add(ctx, orders[4])
	orders[5] = keepers.newTO("00004", 4, 11032, 30, keepers.Sell, keepers.GTE, 90)
	cetKeeper.Add(ctx, orders[5])
	orders[6] = keepers.newTO("00009", 9, 11032, 30, keepers.Sell, keepers.GTE, 90)
	cetKeeper.Add(ctx, orders[6])
	orders[7] = keepers.newTO("00002", 8, 11085, 5, keepers.Buy, keepers.GTE, 98)
	cetKeeper.Add(ctx, orders[7])

	orders[8] = keepers.newTO("00001", 10, 11000, 15, keepers.Buy, keepers.GTE, 998)
	orders[8].TradingPair = "bch/usdt"
	cetKeeper.Add(ctx, orders[8])

	orders[9] = keepers.newTO("00001", 7, 11000, 15, keepers.Buy, keepers.GTE, 998)
	orders[9].TradingPair = "bsv/usdt"
	cetKeeper.Add(ctx, orders[9])

	EndBlocker(ctx, keeper)
	gKeeper := keepers.NewGlobalOrderKeeper(keeper.marketKey, msgCdc)
	allOrders := gKeeper.GetAllOrders(ctx)
	subList := []int{4, 6, 8, 9}
	if len(allOrders) != 4 {
		t.Errorf("Incorrect remain orders.")
	}
	for i, sub := range subList {
		if !keepers.sameTO(orders[sub], allOrders[i]) {
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
		"unfreeze 5 usdt at cosmos1qy352eufqy352eufqy352eufqy35qqqz9ayrkz",
		"send 5 usdt from cosmos1qy352eufqy352eufqy352eufqy35qqqz9ayrkz to cosmos1qy352eufqy352eufqy352eufqy35qqqyej8x24",
		"unfreeze 25 cet at cosmos1qy352eufqy352eufqy352eufqy35qqqyej8x24",
		"send 25 cet from cosmos1qy352eufqy352eufqy352eufqy35qqqyej8x24 to cosmos1qy352eufqy352eufqy352eufqy35qqszrgzaze",
		"unfreeze 27 usdt at cosmos1qy352eufqy352eufqy352eufqy35qqszrgzaze",
		"send 27 usdt from cosmos1qy352eufqy352eufqy352eufqy35qqszrgzaze to cosmos1qy352eufqy352eufqy352eufqy35qqqyej8x24",
		"unfreeze 25 cet at cosmos1qy352eufqy352eufqy352eufqy35qqqf464exq",
		"send 25 cet from cosmos1qy352eufqy352eufqy352eufqy35qqqf464exq to cosmos1qy352eufqy352eufqy352eufqy35qqszrgzaze",
		"unfreeze 27 usdt at cosmos1qy352eufqy352eufqy352eufqy35qqszrgzaze",
		"send 27 usdt from cosmos1qy352eufqy352eufqy352eufqy35qqszrgzaze to cosmos1qy352eufqy352eufqy352eufqy35qqqf464exq",
	}
	unfreezeRecords := []string{
		"unfreeze 40 btc at cosmos1qy352eufqy352eufqy352eufqy35qqpqe926wf",
		"unfreeze 120 btc at cosmos1qy352eufqy352eufqy352eufqy35qqq9yynnh8",
		"unfreeze 1 usdt at cosmos1qy352eufqy352eufqy352eufqy35qqqz9ayrkz",
		"unfreeze 1 usdt at cosmos1qy352eufqy352eufqy352eufqy35qqszrgzaze",
	}
	for i, rec := range bnk.records {
		fmt.Printf("bnk.records[%d] %s\n", i, rec)
	}
	for i, rec := range records {
		if rec != bnk.records[i] {
			t.Errorf("Incorrect Records, actual : %s, expect : %s", bnk.records[i], rec)
		}
	}
	unf := bnk.records[len(records):]
	sort.Strings(unf)
	sort.Strings(unfreezeRecords)
	for i, rec := range unfreezeRecords {
		if rec != unf[i] {
			t.Errorf("Incorrect Records, actual : %s, expect : %s", unf[i], rec)
		}
	}
}
