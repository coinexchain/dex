package market

import (
	"bytes"
	"fmt"
	"sort"
	"testing"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	sdkstore "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/market/internal/keepers"
	"github.com/coinexchain/dex/modules/market/internal/types"
	"github.com/coinexchain/dex/modules/msgqueue"
	dex "github.com/coinexchain/dex/types"
)

var msgCdc = types.ModuleCdc

var unitTestChainID = "coinex-test"

// TODO: duplicated code, copied from order_keeper_test.go
func newTO(sender string, seq uint64, price int64, qua int64, side byte, tif int, h int64) *types.Order {
	addr, _ := simpleAddr(sender)
	decPrice := sdk.NewDec(price).QuoInt(sdk.NewInt(10000))
	freeze := qua
	if side == types.BUY {
		freeze = decPrice.Mul(sdk.NewDec(qua)).RoundInt64()
	}
	return &types.Order{
		Sender:      addr,
		Sequence:    seq,
		TradingPair: "cet/usdt",
		OrderType:   types.LIMIT,
		Price:       decPrice,
		Quantity:    qua,
		Side:        side,
		TimeInForce: tif,
		Height:      h,
		Freeze:      freeze,
		LeftStock:   qua,
	}
}
func simpleAddr(s string) (sdk.AccAddress, error) {
	return sdk.AccAddressFromHex("01234567890123456789012345678901234" + s)
}
func newContextAndMarketKey(chainid string) (sdk.Context, storeKeys) {
	db := dbm.NewMemDB()
	ms := sdkstore.NewCommitMultiStore(db)

	keys := storeKeys{}
	keys.marketKey = sdk.NewKVStoreKey(types.StoreKey)
	keys.keyParams = sdk.NewKVStoreKey(params.StoreKey)
	keys.tkeyParams = sdk.NewTransientStoreKey(params.TStoreKey)
	ms.MountStoreWithDB(keys.keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keys.marketKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: chainid, Height: 1000}, false, log.NewNopLogger())
	return ctx, keys
}
func sameTO(a, order *types.Order) bool {
	res := bytes.Equal(order.Sender, order.Sender) && order.Sequence == order.Sequence &&
		order.TradingPair == order.TradingPair && order.OrderType == order.OrderType && a.Price.Equal(order.Price) &&
		order.Quantity == order.Quantity && order.Side == order.Side && order.TimeInForce == order.TimeInForce &&
		order.Height == order.Height
	return res
}

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
	order := newTO("00001", 1, 11051, 50, types.BUY, types.GTE, 998)
	order.Freeze = 50
	order.FrozenFee = 10
	order.DealStock = 20
	ctx, _ := newContextAndMarketKey(unitTestChainID)
	unfreezeCoinsForOrder(ctx, bxKeeper, order, 0, mockFeeK)
	refout := "unfreeze 50 usdt at cosmos1qy352eufqy352eufqy352eufqy35qqqptw34ca"
	if refout != bxKeeper.records[0] {
		t.Errorf("Error in unfreezeCoinsForOrder")
	}

	coinFee := sdk.NewDec(order.DealStock).Mul(sdk.NewDec(order.FrozenFee)).Quo(sdk.NewDec(order.Quantity)).RoundInt64()
	refout = fmt.Sprintf("addr : %s, fee : %s", order.Sender, dex.NewCetCoin(coinFee).String())
	if refout != mockFeeK.records[0] {
		t.Errorf("Error in unfreezeCoinsForOrder")
	}

}

func TestRemoveOrders(t *testing.T) {
	axk := &mocAssertStatusKeeper{}
	bnk := &mocBankxKeeper{}
	ctx, keys := newContextAndMarketKey(unitTestChainID)
	subspace := params.NewKeeper(msgCdc, keys.keyParams, keys.tkeyParams, params.DefaultCodespace).Subspace(types.StoreKey)
	keeper := keepers.NewKeeper(keys.marketKey, axk, bnk, msgCdc, msgqueue.NewProducer(), subspace, mockBancorKeeper{})
	keeper.SetUnixTime(ctx, time.Now().Unix())
	ctx = ctx.WithBlockTime(time.Unix(time.Now().Unix()+int64(25*60*60), 0))
	parameters := keepers.Params{}
	parameters.GTEOrderLifetime = 1
	parameters.MaxExecutedPriceChangeRatio = keepers.DefaultMaxExecutedPriceChangeRatio
	keeper.SetParams(ctx, parameters)

	keeper.SetMarket(ctx, types.MarketInfo{
		Stock:             "cet",
		Money:             "usdt",
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(0),
	})
	keeper.SetMarket(ctx, types.MarketInfo{
		Stock:             "btc",
		Money:             "usdt",
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(0),
	})

	cetKeeper := keepers.NewOrderKeeper(keys.marketKey, "cet/usdt", msgCdc)
	btcKeeper := keepers.NewOrderKeeper(keys.marketKey, "btc/usdt", msgCdc)
	order := newTO("00001", 1, 11051, 50, types.BUY, types.GTE, 98)
	order.TradingPair = "btc/usdt"
	btcKeeper.Add(ctx, order)
	order = newTO("00005", 5, 12039, 120, types.SELL, types.GTE, 96)
	order.TradingPair = "btc/usdt"
	btcKeeper.Add(ctx, order)
	order = newTO("00002", 2, 11080, 50, types.BUY, types.GTE, 98)
	cetKeeper.Add(ctx, order)
	order = newTO("00002", 3, 10900, 50, types.BUY, types.GTE, 92)
	cetKeeper.Add(ctx, order)
	order = newTO("00004", 4, 11032, 60, types.SELL, types.GTE, 90)
	cetKeeper.Add(ctx, order)

	EndBlocker(ctx, keeper)
	gKeeper := keepers.NewGlobalOrderKeeper(keys.marketKey, msgCdc)
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
	addr01, _ := simpleAddr("00001")
	axk.forbiddenAddrList = []sdk.AccAddress{addr01}
	bnk := &mocBankxKeeper{}
	ctx, keys := newContextAndMarketKey(unitTestChainID)
	subspace := params.NewKeeper(msgCdc, keys.keyParams, keys.tkeyParams, params.DefaultCodespace).Subspace(types.StoreKey)
	keeper := keepers.NewKeeper(keys.marketKey, axk, bnk, msgCdc, msgqueue.NewProducer(), subspace, mockBancorKeeper{})
	delistKeeper := keepers.NewDelistKeeper(keys.marketKey)
	delistKeeper.AddDelistRequest(ctx, ctx.BlockHeight(), "btc/usdt")
	// currDay := ctx.BlockHeader().Time.Unix()
	// keeper.orderClean.SetUnixTime(ctx, currDay)
	keeper.SetUnixTime(ctx, time.Now().Unix())
	ctx = ctx.WithBlockTime(time.Now())
	parameters := keepers.Params{}
	parameters.GTEOrderLifetime = 1
	parameters.MaxExecutedPriceChangeRatio = keepers.DefaultMaxExecutedPriceChangeRatio
	keeper.SetParams(ctx, parameters)

	keeper.SetMarket(ctx, types.MarketInfo{
		Stock:             "cet",
		Money:             "usdt",
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(0),
	})
	keeper.SetMarket(ctx, types.MarketInfo{
		Stock:             "btc",
		Money:             "usdt",
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(0),
	})
	keeper.SetMarket(ctx, types.MarketInfo{
		Stock:             "bch",
		Money:             "usdt",
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(0),
	})
	keeper.SetMarket(ctx, types.MarketInfo{
		Stock:             "bsv",
		Money:             "cet",
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(0),
	})

	cetKeeper := keepers.NewOrderKeeper(keys.marketKey, "cet/usdt", msgCdc)
	btcKeeper := keepers.NewOrderKeeper(keys.marketKey, "btc/usdt", msgCdc)
	orders := make([]*types.Order, 10)
	orders[0] = newTO("00001", 1, 11051, 60, types.BUY, types.GTE, 98)
	orders[0].TradingPair = "btc/usdt"
	btcKeeper.Add(ctx, orders[0])
	orders[1] = newTO("00005", 5, 12039, 120, types.SELL, types.IOC, 1000)
	orders[1].TradingPair = "btc/usdt"
	btcKeeper.Add(ctx, orders[1])
	orders[2] = newTO("00020", 6, 11039, 100, types.SELL, types.IOC, 1000)
	orders[2].TradingPair = "btc/usdt"
	btcKeeper.Add(ctx, orders[2])

	orders[3] = newTO("00202", 2, 11080, 50, types.BUY, types.GTE, 98)
	cetKeeper.Add(ctx, orders[3])
	orders[4] = newTO("00102", 3, 10900, 50, types.BUY, types.GTE, 92)
	cetKeeper.Add(ctx, orders[4])
	orders[5] = newTO("00004", 4, 11032, 30, types.SELL, types.GTE, 90)
	cetKeeper.Add(ctx, orders[5])
	orders[6] = newTO("00009", 9, 11032, 30, types.SELL, types.GTE, 90)
	cetKeeper.Add(ctx, orders[6])
	orders[7] = newTO("00002", 8, 11085, 5, types.BUY, types.GTE, 98)
	cetKeeper.Add(ctx, orders[7])

	orders[8] = newTO("00001", 10, 11000, 15, types.BUY, types.GTE, 998)
	orders[8].TradingPair = "bch/usdt"
	cetKeeper.Add(ctx, orders[8])

	orders[9] = newTO("00001", 7, 11000, 15, types.BUY, types.GTE, 998)
	orders[9].TradingPair = "bsv/usdt"
	cetKeeper.Add(ctx, orders[9])

	EndBlocker(ctx, keeper)
	gKeeper := keepers.NewGlobalOrderKeeper(keys.marketKey, msgCdc)
	allOrders := gKeeper.GetAllOrders(ctx)
	subList := []int{4, 6, 8, 9}
	if len(allOrders) != 4 {
		t.Errorf("Incorrect remain orders.")
	}
	for i, sub := range subList {
		if !sameTO(orders[sub], allOrders[i]) {
			t.Errorf("Incorrect remain orders.")
		}
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
