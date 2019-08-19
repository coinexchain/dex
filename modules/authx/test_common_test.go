package authx_test

import (
	"github.com/coinexchain/dex/app"
	"time"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/cosmos/cosmos-sdk/x/supply"

	dex "github.com/coinexchain/dex/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

type testInput struct {
	ctx sdk.Context
	axk authx.AccountXKeeper
	ak  auth.AccountKeeper
	sk  supply.Keeper
	cdc *codec.Codec
	tk  asset.Keeper
}

func setupTestInput() testInput {
	testApp := app.NewTestApp()
	ctx := sdk.NewContext(testApp.Cms, abci.Header{ChainID: "test-chain-id", Time: time.Unix(1560334620, 0)}, false, log.NewNopLogger())
	initSupply := dex.NewCetCoinsE8(10000)
	testApp.SupplyKeeper.SetSupply(ctx, supply.NewSupply(initSupply))

	return testInput{ctx: ctx, axk: testApp.AccountXKeeper, ak: testApp.AccountKeeper,
		sk: testApp.SupplyKeeper, cdc: testApp.Cdc, tk: testApp.AssetKeeper}
}
