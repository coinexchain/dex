package keepers_test

import (
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/assert"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/testapp"
	"github.com/coinexchain/dex/testutil"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"testing"
)

var (
	app    = testapp.NewTestApp()
	keeper = app.MarketKeeper
	alice  = testutil.ToAccAddress("aliceaddr")
	bob    = testutil.ToAccAddress("abbowneraddr")
	abc    = sdk.NewCoins(sdk.NewCoin("abc", sdk.NewInt(10e8)))
	abb    = sdk.NewCoins(sdk.NewCoin("abb", sdk.NewInt(10e8)))
)

func IssueAbbToBob(ctx sdk.Context, t *testing.T) {
	app.SupplyKeeper.SetSupply(ctx, supply.Supply{Total: sdk.Coins{}})
	app.AssetKeeper.SetParams(ctx, asset.DefaultParams())
	err := app.AssetKeeper.IssueToken(ctx, "abb", "abb", sdk.NewInt(10e8), bob,
		false, false, false, false, "", "", "123")
	assert.Nil(t, err)
	bobAcc := app.AccountKeeper.NewAccountWithAddress(ctx, bob)
	_ = bobAcc.SetCoins(abb)
	app.AccountKeeper.SetAccount(ctx, bobAcc)
}

func BobSendAbbToAlice(ctx sdk.Context, amount sdk.Coins, t *testing.T) {
	err := app.BankxKeeper.SendCoins(ctx, bob, alice, amount)
	assert.Nil(t, err)
}

func TestKeeper_QuerySeqWithAddr(t *testing.T) {
	ctx := app.NewCtx()
	aliceAcc := app.AccountKeeper.NewAccountWithAddress(ctx, alice)
	app.AccountKeeper.SetAccount(ctx, aliceAcc)
	s, _ := keeper.QuerySeqWithAddr(ctx, alice)
	assert.Equal(t, s, uint64(0))
}

func TestKeeper_SetOrderCleanTime(t *testing.T) {
	ctx := app.NewCtx()
	keeper.SetOrderCleanTime(ctx, 100)
	time := keeper.GetOrderCleanTime(ctx)
	assert.Equal(t, time, int64(100))
}

func TestKeeper_FreezeCoins(t *testing.T) {
	ctx := app.NewCtx()
	aliceAcc := app.AccountKeeper.NewAccountWithAddress(ctx, alice)
	app.AccountKeeper.SetAccount(ctx, aliceAcc)

	token := keeper.GetToken(ctx, "abb")
	assert.Nil(t, token)
	IssueAbbToBob(ctx, t)
	yes := keeper.IsTokenExists(ctx, "abb")
	assert.True(t, yes)
	token = keeper.GetToken(ctx, "abb")
	assert.Equal(t, token.GetSymbol(), "abb")
	no := keeper.IsForbiddenByTokenIssuer(ctx, "abb", alice)
	assert.False(t, no)
	no = keeper.IsTokenForbidden(ctx, "abb")
	assert.False(t, no)

	coin := sdk.Coins{sdk.Coin{Denom: "abb", Amount: sdk.NewInt(10)}}
	BobSendAbbToAlice(ctx, coin, t)
	no = keeper.HasCoins(ctx, alice, abc)
	assert.False(t, no)
	yes = keeper.HasCoins(ctx, alice, coin)
	assert.True(t, yes)
	yes = keeper.IsTokenIssuer(ctx, "abb", bob)
	assert.True(t, yes)
	err := keeper.FreezeCoins(ctx, alice, coin)
	assert.Nil(t, err)

	no = keeper.IsSubScribed("msg")
	assert.False(t, no)
	exist := keeper.IsBancorExist(ctx, "abc")
	assert.False(t, exist)

	err = keeper.GetBankxKeeper().UnFreezeCoins(ctx, alice, coin)
	assert.Nil(t, err)
	no = keeper.GetAssetKeeper().IsTokenExists(ctx, "abc")
	assert.False(t, no)
}

func TestKeeper_GetParams(t *testing.T) {
	ctx := app.NewCtx()
	keeper.SetParams(ctx, market.DefaultParams())
	params := keeper.GetParams(ctx)
	assert.Equal(t, params, market.DefaultParams())
	fee := keeper.GetMarketFeeMin(ctx)
	assert.Equal(t, fee, market.DefaultParams().MarketFeeMin)
}

func TestKeeper_SetMarket(t *testing.T) {
	ctx := app.NewCtx()
	info1 := market.MarketInfo{
		Stock:             "abc",
		Money:             "abb",
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(10),
	}
	info2 := market.MarketInfo{
		Stock:             "abd",
		Money:             "abb",
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(100),
	}
	err := keeper.SetMarket(ctx, info1)
	assert.Nil(t, err)
	err = keeper.SetMarket(ctx, info1)
	assert.Nil(t, err)

	err = keeper.SetMarket(ctx, info2)
	assert.Nil(t, err)

	infos := keeper.GetAllMarketInfos(ctx)
	assert.Equal(t, len(infos), 2)

	assert.Equal(t, int64(0), keeper.MarketCountOfStock(ctx, "abb"))
	assert.Equal(t, int64(1), keeper.MarketCountOfStock(ctx, "abc"))

	err = keeper.RemoveMarket(ctx, "abc/abb")
	assert.Nil(t, err)

	assert.Equal(t, int64(0), keeper.MarketCountOfStock(ctx, "abc"))

	info, e := keeper.GetMarketInfo(ctx, "abd/abb")
	assert.Nil(t, e)
	assert.Equal(t, info.Stock, "abd")

	p, e := keeper.GetMarketLastExePrice(ctx, "abd/abb")
	assert.Nil(t, e)
	assert.Equal(t, p.String(), "100.000000000000000000")

	yes := keeper.IsMarketExist(ctx, "abd/abb")
	assert.True(t, yes)

	no := keeper.IsMarketExist(ctx, "abc/abb")
	assert.False(t, no)

	info3 := market.MarketInfo{
		Stock:             "abc",
		Money:             "xyz",
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(10),
	}
	err = keeper.SetMarket(ctx, info3)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), keeper.MarketCountOfStock(ctx, "abc"))

	info3 = market.MarketInfo{
		Stock:             "abc",
		Money:             "btc",
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(10),
	}
	err = keeper.SetMarket(ctx, info3)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), keeper.MarketCountOfStock(ctx, "abc"))

	info3 = market.MarketInfo{
		Stock:             "abc",
		Money:             "usdt",
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(10),
	}
	err = keeper.SetMarket(ctx, info3)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), keeper.MarketCountOfStock(ctx, "abc"))
}
