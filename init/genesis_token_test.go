package init

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/server"

	"github.com/coinexchain/dex/app"
	"github.com/coinexchain/dex/modules/asset"
)

func TestAddGenesisToken(t *testing.T) {
	token := &asset.BaseToken{}
	_ = token.SetName("aaa")
	_ = token.SetSymbol("aaa")

	genesis := app.GenesisState{
		AssetData: asset.GenesisState{
			Tokens: []asset.Token{token},
		},
	}
	_, err := addGenesisToken(genesis, token)
	assert.Error(t, err)

	token = &asset.BaseToken{}
	_ = token.SetName("bbb")
	_ = token.SetSymbol("bbb")
	state, _ := addGenesisToken(genesis, token)
	require.Equal(t, token.GetSymbol(), state.AssetData.Tokens[1].GetSymbol())
}

func TestParseTokenInfo(t *testing.T) {

	defer os.RemoveAll("./keys")
	_, err := parseTokenInfo()
	assert.Error(t, err)

	viper.Set(flagOwner, "owner")
	_, err = parseTokenInfo()
	assert.Error(t, err)

	viper.Set(flagOwner, "coinex1paehyhx9sxdfwc3rjf85vwn6kjnmzjemtedpnl")
	viper.Set(flagName, "1")
	_, err = parseTokenInfo()
	assert.Error(t, err)

	viper.Set(flagName, "aaa")
	viper.Set(flagSymbol, "1")
	_, err = parseTokenInfo()
	assert.Error(t, err)

	viper.Set(flagSymbol, "aaa")
	viper.Set(flagTotalSupply, int64(100))
	viper.Set(flagTotalBurn, int64(100))
	viper.Set(flagTotalMint, int64(100))
	token, _ := parseTokenInfo()
	require.Equal(t, "aaa", token.GetSymbol())
}

func TestAddGenesisTokenCmd(t *testing.T) {
	ctx := server.NewDefaultContext()
	cdc := app.MakeCodec()
	ctx.Config.Genesis = "genesis.json"

	viper.Set(flagOwner, "coinex1paehyhx9sxdfwc3rjf85vwn6kjnmzjemtedpnl")
	viper.Set(flagSymbol, "abbbc")
	viper.Set(flagTotalSupply, int64(100))
	viper.Set(flagTotalBurn, int64(100))
	viper.Set(flagTotalMint, int64(100))
	viper.Set("home", "./")

	defer os.Remove("./genesis.json")
	_, _, _ = initializeGenesisFile(cdc, "./genesis.json")
	cmd := AddGenesisTokenCmd(ctx, cdc)
	require.Equal(t, nil, cmd.RunE(nil, []string{}))
}
