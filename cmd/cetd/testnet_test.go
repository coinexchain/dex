package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	cfg "github.com/tendermint/tendermint/config"

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
	//viper.Set(flagNumValidators, 2)
	//cdc := app.MakeCodec()
	//ctx := server.NewDefaultContext()
	//err := initTestnet(ctx.Config, cdc)
	//_ = os.RemoveAll("./0")
	//_ = os.RemoveAll("./1")
	//_ = os.RemoveAll("./keys")
	//_ = os.RemoveAll("./gentxs")
	//require.Equal(t, nil, err)
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
