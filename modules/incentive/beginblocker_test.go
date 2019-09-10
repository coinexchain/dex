package incentive_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	cetapp "github.com/coinexchain/dex/app"
	"github.com/coinexchain/dex/modules/incentive"
	dex "github.com/coinexchain/dex/types"
)

type TestInput struct {
	ctx    sdk.Context
	cdc    *codec.Codec
	keeper incentive.Keeper
	ak     auth.AccountKeeper
}

func SetupTestInput() TestInput {
	app := cetapp.NewTestApp()
	ctx := sdk.NewContext(app.Cms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	return TestInput{ctx: ctx, cdc: app.Cdc(), keeper: app.IncentiveKeeper(), ak: app.AccountKeeper()}
}

func TestBeginBlockerInvalidCoin(t *testing.T) {
	input := SetupTestInput()
	_ = input.keeper.SetState(input.ctx, incentive.State{HeightAdjustment: 10})
	input.keeper.SetParams(input.ctx, incentive.DefaultParams())
	err := incentive.BeginBlocker(input.ctx, input.keeper)
	require.Equal(t, 0xa, int(err.Result().Code))
}

func TestBeginBlocker(t *testing.T) {
	input := SetupTestInput()
	_ = input.keeper.SetState(input.ctx, incentive.State{HeightAdjustment: 10})
	input.keeper.SetParams(input.ctx, incentive.DefaultParams())
	acc := input.ak.NewAccountWithAddress(input.ctx, incentive.PoolAddr)
	_ = acc.SetCoins(dex.NewCetCoins(10000 * 1e8))
	input.ak.SetAccount(input.ctx, acc)
	err := incentive.BeginBlocker(input.ctx, input.keeper)
	require.Equal(t, nil, err)
}

func TestIncentiveCoinsAddress(t *testing.T) {
	require.Equal(t, "coinex1gc5t98jap4zyhmhmyq5af5s7pyv57w5694el97", incentive.PoolAddr.String())
}

func TestIncentiveCoinsAddressInTestNet(t *testing.T) {
	config := sdk.GetConfig()
	testnetAddrPrefix := "cettest"
	config.SetBech32PrefixForAccount(testnetAddrPrefix, testnetAddrPrefix+sdk.PrefixPublic)
	require.Equal(t, "cettest1gc5t98jap4zyhmhmyq5af5s7pyv57w566ewmx0", incentive.PoolAddr.String())
}

func TestMain(m *testing.M) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(dex.Bech32MainPrefix, dex.Bech32MainPrefix+sdk.PrefixPublic)
	os.Exit(m.Run())
}
