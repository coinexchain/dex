package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"

	"github.com/coinexchain/dex/app"
)

func TestTestnetCmd(t *testing.T) {
	ctx := server.Context{
		Config: &cfg.Config{},
	}
	cdc := app.MakeCodec()
	cmd := testnetCmd(&ctx, cdc, app.ModuleBasics, genaccounts.AppModuleBasic{})
	require.Equal(t, "Initialize files for a Cetd testnet", cmd.Short)
}

func TestInitTestnet(t *testing.T) {
	testHmoe := "./testhome"
	testDataDir := "./testnetdata"
	defer os.RemoveAll(testHmoe)
	defer os.RemoveAll(testDataDir)

	os.Args = []string{"cetd", "testnet", "--v", "2", "-o", testDataDir}
	cetdCmd := createCetdCmd()
	executor := cli.PrepareBaseCmd(cetdCmd, "GA", testHmoe)

	err := executor.Execute()
	require.NoError(t, err)
	// TODO: more asserts
}

func TestInitGenFiles(t *testing.T) {
	//cdc := app.MakeCodec()
	//coins := sdk.Coins{
	//	sdk.Coin{Denom: "abc", Amount: sdk.NewInt(100)},
	//}
	//addr, _ := sdk.AccAddressFromBech32("coinex1paehyhx9sxdfwc3rjf85vwn6kjnmzjemtedpnl")
	//accInfo := &accountInfo{
	//	addr,
	//	coins,
	//	coins,
	//	100,
	//	200,
	//}
	//acc, _ := newGenesisAccount(accInfo)
	//err := initGenFiles(cdc, "chain", []app.GenesisAccount{acc}, []string{"./genesis.json"}, 1)
	//require.Equal(t, nil, err)
}
